package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Product struct {
	ID        uint64     `gorm:"primary_key;auto_increment;column:product_id" json:"id"`
	Code      string     `gorm:"not null;column:product_code" json:"code"`
	Name      string     `gorm:"not null;size:100;column:product_name" json:"name"`
	Price     float32    `gorm:"not null;" json:"price"`
	Material  string     `json:"material"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (p *Product) FindProducts(db *gorm.DB) (*[]Product, error) {
	var err error
	ps := []Product{}
	err = db.Debug().Model(&Product{}).
		Order("product_id, created_at desc").
		Find(&ps).Error

	if err != nil {
		return &[]Product{}, err
	}

	return &ps, nil
}

func (Product) TableName() string {
	return "mstProduct"
}
