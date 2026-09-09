package request

import "mysql/model"

type ProductRequestCreate struct {
	ProductRequest []ProductRequest `json:"product"`
}

type ProductRequest struct {
	Name string `gorm:"type:varchar(150);not null" json:"name"`
}

type ProductRequestUpdate struct {
	Name   string              `gorm:"type:varchar(150);not null" json:"name"`
	Status model.ProductStatus `gorm:"type:enum('ACTIVE','INACTIVE');not null;default:ACTIVE" json:"status"`
}
