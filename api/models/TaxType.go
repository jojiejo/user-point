package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type TaxType struct {
	ID          int        `gorm:"primary_key;auto_increment" json:"id"`
	Name        string     `gorm:"not null;" json:"name"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt   *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

func (tt *TaxType) FindAllTaxTypes(db *gorm.DB) (*[]TaxType, error) {
	var err error
	tts := []TaxType{}
	err = db.Debug().Model(&TaxType{}).
		Order("id, created_at desc").
		Find(&tts).Error

	if err != nil {
		return &[]TaxType{}, err
	}

	return &tts, nil
}

func (TaxType) TableName() string {
	return "tax_type"
}
