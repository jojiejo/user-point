package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type RebateCalculationType struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Name      string     `gorm:"not null;size:100" json:"name"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

func (rct *RebateCalculationType) FindRebateCalculationTypes(db *gorm.DB) (*[]RebateCalculationType, error) {
	var err error
	rcts := []RebateCalculationType{}
	err = db.Debug().Model(&RebateCalculationType{}).
		Order("id, created_at desc").
		Find(&rcts).Error

	if err != nil {
		return &[]RebateCalculationType{}, err
	}

	return &rcts, nil
}

func (RebateCalculationType) TableName() string {
	return "rebate_calculation_type"
}
