package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type RebatePeriod struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Name      string     `gorm:"not null;size:100" json:"name"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

func (rp *RebatePeriod) FindRebatePeriods(db *gorm.DB) (*[]RebatePeriod, error) {
	var err error
	rps := []RebatePeriod{}
	err = db.Debug().Model(&RebatePeriod{}).
		Order("id, created_at desc").
		Find(&rps).Error

	if err != nil {
		return &[]RebatePeriod{}, err
	}

	return &rps, nil
}

func (RebateProgram) TableName() string {
	return "rebate_program"
}
