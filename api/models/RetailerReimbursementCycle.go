package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type RetailerReimbursementCycle struct {
	ID        int       `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (retailerReimburstmentCycle *RetailerReimbursementCycle) FindAllRetailerReimbursementCycles(db *gorm.DB) (*[]RetailerReimbursementCycle, error) {
	var err error
	retailerReimburstmentCycles := []RetailerReimbursementCycle{}
	err = db.Debug().Model(&RetailerReimbursementCycle{}).Limit(100).Order("created_at desc").Find(&retailerReimburstmentCycles).Error
	if err != nil {
		return &[]RetailerReimbursementCycle{}, err
	}

	return &retailerReimburstmentCycles, nil
}

func (RetailerReimbursementCycle) TableName() string {
	return "retailer_reimbursement_cycle"
}
