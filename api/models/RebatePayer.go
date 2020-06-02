package models

import (
	"errors"
	"html"
	"strconv"
	"strings"
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

type AssignedRebateFromFile struct {
	MCMSID                 string                          `sql:"-" json:"mcms_id"`
	CCID                   uint64                          `json:"cc_id"`
	GSAPCustomerMasterData ShortenedGSAPCustomerMasterData `json:"gsap_customer_master_data"`
	Payer                  ShortenedPayer                  `json:"-"`
	RebateProgramID        uint64                          `json:"rebate_program_id"`
	RebateProgramName      string                          `sql:"-" json:"rebate_program_name"`
	RebateProgram          RebateProgram                   `json:"rebate_program"`
	StartedAt              time.Time                       `json:"started_at"`
	EndedAt                *time.Time                      `json:"ended_at"`
	CreatedAt              time.Time                       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              time.Time                       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt              *time.Time                      `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type BulkAssignRebate struct {
	RebateCollection []AssignedRebateFromFile `json:"rebate_collection"`
}

func (bar *BulkAssignRebate) Prepare() {
	if len(bar.RebateCollection) > 0 {
		for i, _ := range bar.RebateCollection {
			bar.RebateCollection[i].MCMSID = html.EscapeString(strings.TrimSpace(bar.RebateCollection[i].MCMSID))
			bar.RebateCollection[i].RebateProgramName = html.EscapeString(strings.TrimSpace(bar.RebateCollection[i].RebateProgramName))
		}
	}
}

func (bar *BulkAssignRebate) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if len(bar.RebateCollection) > 0 {
		var convertedLineNumber string

		for i, _ := range bar.RebateCollection {
			convertedLineNumber = strconv.Itoa(i + 1)

			if bar.RebateCollection[i].MCMSID == "" {
				err = errors.New("MCMS ID field in line " + convertedLineNumber + " is required")
				errorMessages["mcms_id_line_"+convertedLineNumber] = err.Error()
			}

			if bar.RebateCollection[i].RebateProgramName == "" {
				err = errors.New("Rebate program name field in line " + convertedLineNumber + " is required")
				errorMessages["rebate_program_name_line_"+convertedLineNumber] = err.Error()
			}
		}
	}

	return errorMessages
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

func (bar *BulkAssignRebate) ChargePrepare() {
	if len(bar.RebateCollection) > 0 {
		for i, _ := range bar.RebateCollection {
			bar.RebateCollection[i].CreatedAt = time.Now()
		}
	}
}

func (bar *BulkAssignRebate) ChargeValidate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if len(bar.RebateCollection) > 0 {
		var convertedLineNumber string

		for i, _ := range bar.RebateCollection {
			convertedLineNumber = strconv.Itoa(i + 1)

			if bar.RebateCollection[i].CCID < 1 {
				err = errors.New("Payer field in line " + convertedLineNumber + " is required")
				errorMessages["payer_line_"+convertedLineNumber] = err.Error()
			}

			if bar.RebateCollection[i].RebateProgramID < 1 {
				err = errors.New("Rebate program field in line " + convertedLineNumber + " is required")
				errorMessages["rebate_program_line_"+convertedLineNumber] = err.Error()
			}
		}
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

func (bar *BulkAssignRebate) BulkCheckAssignRebate(db *gorm.DB) (*BulkAssignRebate, map[string]string) {
	var err error
	var errorMessages = make(map[string]string)

	if len(bar.RebateCollection) > 0 {
		var convertedMCMSID int
		var convertedLineNumber string

		for i, _ := range bar.RebateCollection {
			convertedLineNumber = strconv.Itoa(i + 1)
			convertedMCMSID, err = strconv.Atoi(bar.RebateCollection[i].MCMSID)
			if err != nil {
				err = errors.New("MCMS ID : " + bar.RebateCollection[i].MCMSID + " could not be processed")
				errorMessages["unprocessed_mcms_id_line_"+convertedLineNumber] = err.Error()
			}

			customerDataErr := db.Debug().Model(&ShortenedGSAPCustomerMasterData{}).Unscoped().Where("mcms_id = ?", convertedMCMSID).Order("mcms_id desc").Take(&bar.RebateCollection[i].GSAPCustomerMasterData).Error
			if customerDataErr != nil {
				err = errors.New("MCMS ID : " + bar.RebateCollection[i].MCMSID + " is invalid")
				errorMessages["mcms_id_line_"+convertedLineNumber] = err.Error()
			}

			payerDataErr := db.Debug().Model(&ShortenedPayer{}).Unscoped().Where("mcms_id = ?", convertedMCMSID).Order("cc_id desc").Take(&bar.RebateCollection[i].Payer).Error
			if payerDataErr != nil {
				err = errors.New("MCMS ID : " + bar.RebateCollection[i].MCMSID + " is not registered as an account")
				errorMessages["not_registered_mcms_id_line_"+convertedLineNumber] = err.Error()
			}

			bar.RebateCollection[i].CCID = bar.RebateCollection[i].Payer.CCID

			rebateProgramDataErr := db.Debug().Model(&RebateProgram{}).Unscoped().Where("name = ?", bar.RebateCollection[i].RebateProgramName).Order("id desc").Take(&bar.RebateCollection[i].RebateProgram).Error
			if rebateProgramDataErr != nil {
				err = errors.New("Fee Name : " + bar.RebateCollection[i].RebateProgramName + " is invalid")
				errorMessages["fee_name_line_"+convertedLineNumber] = err.Error()
			}

			bar.RebateCollection[i].RebateProgramID = bar.RebateCollection[i].RebateProgram.ID
			bar.RebateCollection[i].StartedAt = bar.RebateCollection[i].RebateProgram.StartedAt
			bar.RebateCollection[i].EndedAt = bar.RebateCollection[i].RebateProgram.EndedAt
		}
	}

	return bar, errorMessages
}

func (bar *BulkAssignRebate) BulkAssignRebate(db *gorm.DB) (*BulkAssignRebate, map[string]string) {
	var err error
	var errorMessages = make(map[string]string)

	if len(bar.RebateCollection) > 0 {
		for i, _ := range bar.RebateCollection {
			var convertedLineNumber string
			convertedLineNumber = strconv.Itoa(i + 1)

			err = db.Debug().Model(&RebatePayer{}).Create(&bar.RebateCollection[i]).Error
			if err != nil {
				errorMessages["line_"+convertedLineNumber] = err.Error()
			}
		}
	}

	return bar, errorMessages
}

func (AssignedRebateFromFile) TableName() string {
	return "rebate_program_payer"
}

func (RebatePayer) TableName() string {
	return "rebate_program_payer"
}
