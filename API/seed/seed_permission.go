package seed

import (
	"mysql/constant/seeddata"
	"mysql/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedPermissions(db *gorm.DB) error {

	var permissions []model.Permission

	for _, p := range seeddata.Permissions {
		permissions = append(permissions, model.Permission{
			Name:        p.Name,
			DisplayName: p.DisplayName,
		})
	}

	return db.Clauses(clause.OnConflict{
		// OnConflict is Upsert = Update + Insert
		Columns: []clause.Column{
			{Name: "name"}, // unique key
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
		}),
	}).Create(&permissions).Error
}
