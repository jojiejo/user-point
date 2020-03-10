package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type RetailerPaymentTerm struct {
	ID        int        `gorm:"primary_key;auto_increment" json:"id"`
	Name      string     `gorm:"not null" json:"name"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (retailerPaymentTerm *RetailerPaymentTerm) FindAllRetailerPaymentTerms(db *gorm.DB) (*[]RetailerPaymentTerm, error) {
	var err error
	retailerPaymentTerms := []RetailerPaymentTerm{}
	err = db.Debug().Model(&RetailerPaymentTerm{}).Limit(100).Order("created_at desc").Find(&retailerPaymentTerms).Error
	if err != nil {
		return &[]RetailerPaymentTerm{}, err
	}

	return &retailerPaymentTerms, nil
}

func (RetailerPaymentTerm) TableName() string {
	return "retailer_payment_term"
}
