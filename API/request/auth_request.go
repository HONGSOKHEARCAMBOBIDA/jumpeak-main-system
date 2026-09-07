package request

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginQrRequest struct {
	QrToken string `json:"qr_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `cookie:"refresh_token" binding:"required"`
}

type NewPasswordRequest struct {
	NewPassword string `json:"new_password"`
}
