package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type RebatePayer struct {
	ID              uint64         `gorm:"primary_key;auto_increment" json:"id"`
	CCID            uint64         `gorm:"not null;size:100" json:"cc_id"`
	Payer           ShortenedPayer `gorm:"foreignkey:CCID;association_foreignkey:CCID;" json:"payer"`
	RebateProgramID uint64         `gorm:"not null;size:100" json:"rebate_program_id"`
	RebateProgram   RebateProgram  `gorm:"foreignkey:RebateProgramID;association_foreignkey:ID;" json:"rebate_program"`
	StartedAt       time.Time      `json:"started_at"`
	EndedAt         *time.Time     `json:"ended_at"`
	CreatedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       *time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type PostedRebatePayer struct {
	CCID            []uint64   `gorm:"not null;size:100" json:"cc_id"`
	RebateProgramID uint64     `gorm:"not null;size:100" json:"rebate_program_id"`
	StartedAt       time.Time  `json:"started_at"`
	EndedAt         *time.Time `json:"ended_at"`
	CreatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (prp *PostedRebatePayer) Prepare() {
	prp.CreatedAt = time.Now()
}

func (prp *PostedRebatePayer) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if len(prp.CCID) < 1 {
		err = errors.New("Account field is required")
		errorMessages["account"] = err.Error()
	}

	if prp.RebateProgramID < 1 {
		err = errors.New("Rebate program field is required")
		errorMessages["rebate_program"] = err.Error()
	}

	return errorMessages
}

func (rp *RebatePayer) FindRebatePayerRelations(db *gorm.DB) (*[]RebatePayer, error) {
	var err error
	rps := []RebatePayer{}
	err = db.Debug().Model(&RebatePayer{}).Unscoped().
		Preload("Payer").
		Preload("Payer.GSAPCustomerMasterData").
		Preload("Payer.LatestPayerStatus").
		Preload("Payer.LatestPayerStatus.PayerStatus").
		Preload("RebateProgram").
		Preload("RebateProgram.RebateType").
		Preload("RebateProgram.RebateCalculationType").
		Preload("RebateProgram.RebatePeriod").
		Preload("RebateProgram.Site").
		Preload("RebateProgram.Product").
		Order("id, created_at desc").
		Find(&rps).Error

	if err != nil {
		return &[]RebatePayer{}, err
	}

	return &rps, nil
}

func (prp *PostedRebatePayer) CreateRebatePayerRelation(db *gorm.DB) (*PostedRebatePayer, map[string]string) {
	var err error
	var errorMessages = make(map[string]string)
	rebatePayerRelation := RebatePayer{}

	if len(prp.CCID) > 0 {
		for i, _ := range prp.CCID {
			var convertedLineNumber string
			convertedLineNumber = strconv.Itoa(i)
			rebatePayerRelation.CCID = prp.CCID[i]
			rebatePayerRelation.RebateProgramID = prp.RebateProgramID
			rebatePayerRelation.StartedAt = prp.StartedAt
			rebatePayerRelation.EndedAt = prp.EndedAt
			rebatePayerRelation.CreatedAt = prp.CreatedAt
			rebatePayerRelation.UpdatedAt = prp.UpdatedAt

			err = db.Debug().Model(&RebatePayer{}).Create(&rebatePayerRelation).Error
			if err != nil {
				errorMessages["account_"+convertedLineNumber] = err.Error()
			}
		}
	}

	return prp, errorMessages
}

func (RebatePayer) TableName() string {
	return "rebate_program_payer"
}
