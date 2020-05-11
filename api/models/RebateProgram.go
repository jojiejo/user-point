package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type RebateProgram struct {
	ID                      uint64                `gorm:"primary_key;auto_increment" json:"id"`
	Name                    string                `gorm:"not null;size:100" json:"code"`
	RebateTypeID            uint64                `gorm:"not null;" json:"rebate_type_id"`
	RebateType              RebateType            `json:"rebate_type"`
	RebateCalculationTypeID uint64                `gorm:"not null;" json:"rebate_calculation_type_id"`
	RebateCalculationType   RebateCalculationType `json:"rebate_calculation_type"`
	RebatePeriodID          uint64                `gorm:"not null;" json:"rebate_period_id"`
	RebatePeriod            RebatePeriod          `json:"rebate_period"`
	Site                    []Site                `gorm:"many2many:rebate_program_site;association_autoupdate:false;association_jointable_foreignkey:rebate_program_id;association_jointable_foreignkey:site_id" json:"site"`
	Product                 []Product             `gorm:"many2many:rebate_program_product;association_autoupdate:false;association_jointable_foreignkey:rebate_program_id;association_jointable_foreignkey:product_id" json:"product"`
	StartedAt               *time.Time            `json:"started_at"`
	EndedAt                 *time.Time            `json:"ended_at"`
	CreatedAt               time.Time             `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt               time.Time             `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt               *time.Time            `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (rp *RebateProgram) FindRebatePrograms(db *gorm.DB) (*[]RebateProgram, error) {
	var err error
	rps := []RebateProgram{}
	err = db.Debug().Model(&RebateProgram{}).Unscoped().
		Preload("Site").
		Preload("Product").
		Preload("RebateType").
		Preload("RebatePeriod").
		Preload("RebateCalculationType").
		Order("id, created_at desc").
		Find(&rps).Error

	if err != nil {
		return &[]RebateProgram{}, err
	}

	return &rps, nil
}

func (rp *RebateProgram) FindRebateProgramByID(db *gorm.DB, rpID uint64) (*RebateProgram, error) {
	var err error
	err = db.Debug().Model(&Fee{}).Unscoped().
		Preload("Site").
		Preload("Product").
		Preload("RebateType").
		Preload("RebatePeriod").
		Preload("RebateCalculationType").
		Where("id = ?", rpID).
		Order("created_at desc").
		Take(&rp).Error

	if err != nil {
		return &RebateProgram{}, err
	}

	return rp, nil
}

func (RebatePeriod) TableName() string {
	return "rebate_period"
}
