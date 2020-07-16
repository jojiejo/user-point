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
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
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

func (p *Product) FindProductByID(db *gorm.DB, productID uint64) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Unscoped().
		Where("id = ?", productID).
		Order("created_at desc").
		Take(&p).Error

	if err != nil {
		return &Product{}, err
	}

	return p, nil
}

func (p *Product) CreateProduct(db *gorm.DB) (*Product, error) {
	var err error
	err = db.Debug().Model(&Fee{}).Create(&p).Error
	if err != nil {
		return &Product{}, err
	}

	//Select created fee
	_, err = p.FindProductByID(db, p.ID)
	if err != nil {
		return &Product{}, err
	}

	return p, nil
}

func (p *Product) UpdateProduct(db *gorm.DB) (*Product, error) {
	var err error
	dateTimeNow := time.Now()

	//Update the data
	err = db.Debug().Model(&p).Updates(
		map[string]interface{}{
			"product_code": p.Code,
			"product_name": p.Name,
			"price":        p.Price,
			"material":     p.Material,
			"started_at":   p.StartedAt,
			"ended_at":     p.EndedAt,
			"updated_at":   dateTimeNow,
		}).Error

	if err != nil {
		return &Product{}, err
	}

	//Select updated product
	_, err = p.FindProductByID(db, p.ID)
	if err != nil {
		return &Product{}, err
	}

	return p, nil
}

func (Product) TableName() string {
	return "mstProduct"
}
