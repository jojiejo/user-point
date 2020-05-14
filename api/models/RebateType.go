package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type RebateType struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Name      string     `gorm:"not null;size:100" json:"name"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

func (rt *RebateType) FindRebateTypes(db *gorm.DB) (*[]RebateType, error) {
	var err error
	rts := []RebateType{}
	err = db.Debug().Model(&RebateType{}).
		Order("id, created_at desc").
		Find(&rts).Error

	if err != nil {
		return &[]RebateType{}, err
	}

	return &rts, nil
}

func (RebateType) TableName() string {
	return "rebate_type"
}
