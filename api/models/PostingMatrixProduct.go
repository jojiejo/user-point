package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type PostingMatrixProduct struct {
	ID              uint64     `gorm:"primary_key;auto_increment" json:"id"`
	ProductID       uint64     `gorm:"not null" json:"product_id"`
	VolumeProceedGL string     `gorm:"not null" json:"volume_proceed_gl"`
	DiscountGL      string     `gorm:"not null" json:"discount_gl"`
	SurchargeGL     string     `gorm:"not null" json:"surcharge_gl"`
	ProfitCentre    string     `gorm:"not null" json:"profit_centre"`
	MaterialCode    string     `gorm:"not null" json:"material_code"`
	VATID           string     `gorm:"column:posting_matrix_vat_id" json:"vat_id"`
	CreatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (pmp *PostingMatrixProduct) FindPostingMatrixProducts(db *gorm.DB) (*[]PostingMatrixProduct, error) {
	var err error
	pmps := []PostingMatrixProduct{}
	err = db.Debug().Model(&PostingMatrixProduct{}).
		Order("product_id, created_at desc").
		Find(&pmps).Error

	if err != nil {
		return &[]PostingMatrixProduct{}, err
	}

	return &pmps, nil
}

func (RebatePayer) TableName() string {
	return "posting_matrix_product"
}
