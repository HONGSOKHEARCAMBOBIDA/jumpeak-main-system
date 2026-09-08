package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"mysql/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(input request.AuthRequest, c *gin.Context) (*response.AuthResponse, error)
	RefreshToken(refreshToken string, c *gin.Context) (*response.AuthResponse, error)
	GetUserData(ctx context.Context, id int) (response.UserDataResponse, error)
	GetRole(ctx context.Context, id int) ([]model.Role, error)
	Create(ctx context.Context, input request.UserRequestCreate) error
	Update(ctx context.Context, id int, input request.UserRequestUpdate) error
}

type authservice struct {
	db *gorm.DB
}

func NewAuthService() AuthService {
	return &authservice{
		db: config.DB,
	}
}

var requiredPermissions = []string{
	"add.role.has.permission",
	"add.company",
	"update.company",
	"add.Branch",
	"update.Branch",
	"add.user",
	"edit.user",
}

func (s *authservice) GetRole(ctx context.Context, id int) ([]model.Role, error) {
	var role []model.Role
	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	roleQuery := s.db.WithContext(ctx).Table("role r").
		Select(`
		 r.id AS id,
		 r.name AS name,
		 r.display_name AS display_name
	`)
	roleQuery = helper.ApplyAccessGetRole(roleQuery, s.db, user.Role, user)

	if err := roleQuery.Scan(&role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (s *authservice) Login(input request.AuthRequest, c *gin.Context) (*response.AuthResponse, error) {
	key := "login_attempt:" + input.Email
	attempts, _ := utils.Redis.Get(utils.Ctx, key).Int()
	if attempts >= 5 {
		return nil, errors.New("អ្នកព្យាយាមចូលច្រើនពេក សូមព្យាយាមម្តងទៀតក្រោយ 10 នាទី")
	}
	var user model.User
	if err := s.db.Select("id,email,password_hash,role_id,status,name").
		Where("email = ? AND status = ?", input.Email, model.UserStatusActive).
		First(&user).Error; err != nil {
		return nil, errors.New("ព័ត៌មានមិនត្រឹមត្រូវ ឬ អ្នកប្រើប្រាស់ត្រូវបានបិទគណនី")
	}

	var settings []model.Setting
	if err := s.db.Where("`key` IN ?", []string{
		"ACCESS_TOKEN_EXPIRE_HOURS",
		"REFRESH_TOKEN_EXPIRE_DAYS",
	}).Find(&settings).Error; err != nil {
		return nil, errors.New("Setting Not Found")
	}

	settingMap := make(map[string]string)
	for _, s := range settings {
		settingMap[s.Key] = s.Value
	}

	accesstoken, err := strconv.Atoi(settingMap["ACCESS_TOKEN_EXPIRE_HOURS"])
	if err != nil {
		return nil, errors.New("Bad request")
	}
	refreshtoken, err := strconv.Atoi(settingMap["REFRESH_TOKEN_EXPIRE_DAYS"])
	if err != nil {
		return nil, errors.New("Bad request")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	); err != nil {
		utils.Redis.Incr(utils.Ctx, key)
		utils.Redis.Expire(utils.Ctx, key, 10*time.Minute)
		return nil, errors.New("លេខសម្ងាត់មិនត្រឹមត្រូវ")
	}
	utils.Redis.Del(utils.Ctx, key)

	accessExpiry := time.Now().Add(time.Duration(accesstoken) * time.Hour)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role_id": user.RoleID,
		"exp":     accessExpiry.Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenStr, err := accessToken.SignedString(utils.Jwtkey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshTokenStr := hex.EncodeToString(refreshTokenBytes)
	tokenPrefix := refreshTokenStr[:16]
	hashedRefresh := utils.HashToken(refreshTokenStr)
	if err := s.db.Where("user_id = ?", user.ID).Delete(&model.Session{}).Error; err != nil {
		return nil, fmt.Errorf("failed to delete session")
	}
	refreshExpiry := time.Now().Add(time.Duration(refreshtoken) * 24 * time.Hour)
	session := model.Session{
		UserID:       uint(user.ID),
		RefreshToken: string(hashedRefresh),
		TokenPrefix:  tokenPrefix,
		ExpiresAt:    refreshExpiry,
	}
	if err := s.db.Create(&session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	maxAge := int(time.Until(refreshExpiry).Seconds())
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"refresh_token",
		refreshTokenStr,
		maxAge,
		"/",
		"",   // domain: leave empty for current host, or set "yourdomain.com" explicitly
		true, // Secure - HTTPS only
		true, // httpOnly - JS cannot read this
	)
	resp := &response.AuthResponse{
		AccessToken: accessTokenStr,
	}

	return resp, nil
}

func (s *authservice) RefreshToken(refreshToken string, c *gin.Context) (*response.AuthResponse, error) {
	if len(refreshToken) < 16 {
		return nil, errors.New("Invalid refresh token")
	}
	prefix := refreshToken[:16]

	var session model.Session
	err := s.db.Where("token_prefix = ?", prefix).First(&session).Error
	if err != nil {
		return nil, errors.New("Invalid or expired refresh token")
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("Invalid or expired refresh token")
	}

	if !utils.VerifyToken(session.RefreshToken, refreshToken) {
		return nil, errors.New("Invalid or expired refresh token")
	}

	var settings []model.Setting
	if err := s.db.Where("`key` IN ?", []string{
		"ACCESS_TOKEN_EXPIRE_HOURS",
		"REFRESH_TOKEN_EXPIRE_DAYS",
	}).Find(&settings).Error; err != nil {
		return nil, err
	}

	settingMap := make(map[string]string)
	for _, s := range settings {
		settingMap[s.Key] = s.Value
	}

	accesstoken, err := strconv.Atoi(settingMap["ACCESS_TOKEN_EXPIRE_HOURS"])
	if err != nil {
		return nil, err
	}

	refreshtoken, err := strconv.Atoi(settingMap["REFRESH_TOKEN_EXPIRE_DAYS"])
	if err != nil {
		return nil, errors.New("Bad request")
	}

	accessExpiry := time.Now().Add(time.Duration(accesstoken) * time.Minute)
	refreshExpiry := time.Now().Add(time.Duration(refreshtoken) * 24 * time.Hour)
	newRefreshBytes := make([]byte, 32)
	if _, err := rand.Read(newRefreshBytes); err != nil {
		return nil, errors.New("failed to generate refresh token")
	}
	newRefreshStr := hex.EncodeToString(newRefreshBytes)
	newHash := utils.HashToken(newRefreshStr)
	newPrefix := newRefreshStr[:16]

	if err := s.db.Model(&session).Updates(model.Session{
		RefreshToken: newHash,
		TokenPrefix:  newPrefix,
		ExpiresAt:    refreshExpiry,
	}).Error; err != nil {
		return nil, err
	}

	var user model.User
	if err := s.db.Select("id,role_id").Where("id = ?", session.UserID).First(&user).Error; err != nil {
		return nil, err
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role_id": user.RoleID,
		"exp":     accessExpiry.Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(utils.Jwtkey)
	if err != nil {
		return nil, err
	}

	maxAge := int(time.Until(refreshExpiry).Seconds())
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"refresh_token",
		newRefreshStr,
		maxAge,
		"/",
		"",   // domain: leave empty for current host, or set "yourdomain.com" explicitly
		true, // Secure - HTTPS only
		true, // httpOnly - JS cannot read this
	)

	return &response.AuthResponse{
		AccessToken: accessToken,
		//Permissions: permissions,
	}, nil
}

func (s *authservice) GetUserData(ctx context.Context, id int) (response.UserDataResponse, error) {
	var userdata response.UserDataResponse

	err := s.db.WithContext(ctx).
		Table("users u").
		Select(`
			u.id AS id,
			u.name AS name,
			u.role_id AS role_id
		`).
		Where("u.id = ?", id).First(&userdata).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return userdata, fmt.Errorf("user with id %d not found", id)
	}

	if err != nil {
		return userdata, fmt.Errorf("failed to get user data: %w", err)
	}

	if userdata.ID == 0 {
		return userdata, fmt.Errorf("user with id %d not found", id)
	}

	var permissions []model.Permission
	if err := s.db.WithContext(ctx).
		Table("permission p").
		Select("p.name AS name").
		Joins("JOIN role_permission rhp ON rhp.permission_id = p.id").
		Where("rhp.role_id = ? AND p.name IN ?", userdata.RoleID, requiredPermissions).
		Scan(&permissions).Error; err != nil {
		return userdata, fmt.Errorf("failed to get user permissions: %w", err)
	}

	userdata.Permissions = permissions

	return userdata, nil
}

func (s *authservice) Create(ctx context.Context, input request.UserRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	var branch model.Branch
	if err := s.db.WithContext(ctx).First(&branch, input.BranchID).Error; err != nil {
		return err
	}

	email := helper.GenerateEmail(input.Name, 168)
	passwordHash := utils.HasPassword("12345678")

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newdata := model.User{
			CompanyID:    branch.CompanyID,
			BranchID:     input.BranchID,
			Name:         strings.TrimSpace(input.Name),
			Email:        email,
			PasswordHash: passwordHash,
			RoleID:       input.RoleID,
			ManageBranch: input.ManageBranch,
			Status:       model.UserStatusActive,
		}

		if err := tx.Create(&newdata).Error; err != nil {
			return err
		}

		if input.ManageBranch == model.UserManageBranchMultiple &&
			input.BranchIDs != nil && len(*input.BranchIDs) != 0 {
			for i := range *input.BranchIDs {
				userbranch := model.UserBranch{
					UserID:   uint64(newdata.ID),
					BranchID: (*input.BranchIDs)[i],
				}
				if err := tx.Create(&userbranch).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}

func (s *authservice) Update(ctx context.Context, id int, input request.UserRequestUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var branch model.Branch
	if err := s.db.WithContext(ctx).First(&branch, input.BranchID).Error; err != nil {
		return err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var data model.User
		if err := tx.Where("id = ?", id).First(&data).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "classcurriculumn not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch classcurriculumn", nil)
		}
		data.CompanyID = branch.CompanyID
		data.BranchID = input.BranchID
		data.Name = input.Name
		data.RoleID = input.RoleID
		data.ManageBranch = input.ManageBranch
		data.Status = input.Status
		if err := tx.Save(&data).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update student", nil)
		}
		if input.ManageBranch != model.UserManageBranchMultiple {
			if err := tx.Where("user_id = ?", id).Delete(&model.UserBranch{}).Error; err != nil {
				return err
			}
		} else if input.ManageBranch == model.UserManageBranchMultiple {
			if err := tx.Where("user_id =?", id).Delete(&model.UserBranch{}).Error; err != nil {
				return err
			}
			for i := range *input.BranchIDs {
				userbranch := model.UserBranch{
					UserID:   uint64(id),
					BranchID: (*input.BranchIDs)[i],
				}
				if err := tx.Create(&userbranch).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}
