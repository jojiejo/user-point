package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type RebatePayer struct {
	ID              uint64         `gorm:"primary_key;auto_increment" json:"id"`
	CCID            uint64         `gorm:"not null;size:100" json:"cc_id"`
	Payer           ShortenedPayer `gorm:"foreignkey:CCID;association_foreignkey:CCID;" json:"payer"`
	RebateProgramID uint64         `gorm:"not null;size:100" json:"rebate_program_id"`
	RebateProgram   RebateProgram  `gorm:"foreignkey:RebateProgramID;association_foreignkey:RebateProgramID;" json:"rebate_program"`
	StartedAt       *time.Time     `json:"started_at"`
	EndedAt         *time.Time     `json:"ended_at"`
	CreatedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       *time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (rp *RebatePayer) Prepare() {
	rp.CreatedAt = time.Now()
}

func (rp *RebatePayer) FindRebatePayerRelations(db *gorm.DB) (*[]RebatePayer, error) {
	var err error
	rps := []RebatePayer{}
	err = db.Debug().Model(&RebatePayer{}).Unscoped().
		Preload("Payer").
		Preload("RebateProgram").
		Order("id, created_at desc").
		Find(&rps).Error

	if err != nil {
		return &[]RebatePayer{}, err
	}

	return &rps, nil
}

func (RebatePayer) TableName() string {
	return "rebate_payer"
}
