package request

import "mysql/model"

type ProductRequestCreate struct {
	ProductRequest []ProductRequest `json:"product"`
}

type ProductRequest struct {
	Name         string  `gorm:"type:varchar(150);not null" json:"name"`
	DefaultPrice float64 `json:"default_price" gorm:"column:default_price"`
}

type ProductRequestUpdate struct {
	Name         string              `gorm:"type:varchar(150);not null" json:"name"`
	DefaultPrice float64             `json:"default_price" gorm:"column:default_price"`
	Status       model.ProductStatus `gorm:"type:enum('ACTIVE','INACTIVE');not null;default:ACTIVE" json:"status"`
}
