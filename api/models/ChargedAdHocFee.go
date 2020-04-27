package models

import (
	"errors"
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type ChargedAdHocFee struct {
	ID             uint64          `gorm:"primary_key;auto_increment;" json:"id"`
	SubCorporateID uint64          `gorm:"not null" json:"sub_corporate_id"`
	Branch         ShortenedBranch `gorm:"foreignkey:SubCorporateID;association_foreignkey:SubCorporateID;" json:"branch"`
	FeeID          uint64          `gorm:"not null;" json:"fee_id"`
	Fee            ShortenedFee    `json:"fee"`
	Value          float32         `gorm:"not null;" json:"value"`
	UnitID         uint64          `gorm:"not null;" json:"unit_id"`
	Unit           Unit            `json:"unit"`
	Description    string          `gorm:"size:50;" json:"description"`
	CreatedAt      time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt      *time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (cahf *ChargedAdHocFee) Prepare() {
	cahf.Description = html.EscapeString(strings.TrimSpace(cahf.Description))
	cahf.CreatedAt = time.Now()
}

func (cahf *ChargedAdHocFee) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if cahf.FeeID < 1 {
		err = errors.New("Fee field is required")
		errorMessages["unit"] = err.Error()
	}

	if cahf.SubCorporateID < 1 {
		err = errors.New("Branch field is required")
		errorMessages["unit"] = err.Error()
	}

	if cahf.UnitID < 1 {
		err = errors.New("Unit field is required")
		errorMessages["unit"] = err.Error()
	}

	return errorMessages
}

func (cahf *ChargedAdHocFee) FindAdHocFees(db *gorm.DB) (*[]ChargedAdHocFee, error) {
	var err error
	cahfs := []ChargedAdHocFee{}
	err = db.Debug().Model(&ChargedAdHocFee{}).Unscoped().
		Preload("Unit").
		Preload("Fee").
		Preload("Branch").
		Order("id, created_at desc").
		Find(&cahfs).Error

	if err != nil {
		return &[]ChargedAdHocFee{}, err
	}

	if len(cahfs) > 0 {
		for i, _ := range cahfs {
			customerDataErr := db.Debug().Model(&ShortenedGSAPCustomerMasterData{}).Unscoped().Where("mcms_id = ?", cahfs[i].Branch.MCMSID).Order("mcms_id desc").Take(&cahfs[i].Branch.GSAPCustomerMasterData).Error
			if customerDataErr != nil {
				return &[]ChargedAdHocFee{}, err
			}

			cahfs[i].Branch.PaddedMCMSID = fmt.Sprintf("%010v", strconv.Itoa(cahfs[i].Branch.MCMSID))
		}
	}

	return &cahfs, nil
}

func (cahf *ChargedAdHocFee) FindChargedAdHocFeeByID(db *gorm.DB, relationID uint64) (*ChargedAdHocFee, error) {
	var err error
	err = db.Debug().Model(&ChargedAdHocFee{}).Unscoped().
		Preload("Unit").
		Preload("Fee").
		Preload("Branch").
		Where("id = ?", relationID).
		Order("created_at desc").
		Take(&cahf).Error

	if err != nil {
		return &ChargedAdHocFee{}, err
	}

	customerDataErr := db.Debug().Model(&ShortenedGSAPCustomerMasterData{}).Unscoped().Where("mcms_id = ?", cahf.Branch.MCMSID).Order("mcms_id desc").Take(&cahf.Branch.GSAPCustomerMasterData).Error
	if customerDataErr != nil {
		return &ChargedAdHocFee{}, err
	}

	cahf.Branch.PaddedMCMSID = fmt.Sprintf("%010v", strconv.Itoa(cahf.Branch.MCMSID))

	return cahf, nil
}

func (cahf *ChargedAdHocFee) ChargeAdHocFee(db *gorm.DB) (*ChargedAdHocFee, error) {
	var err error
	err = db.Debug().Model(&ChargedAdHocFee{}).Create(&cahf).Error
	if err != nil {
		return &ChargedAdHocFee{}, err
	}

	//Select created fee
	_, err = cahf.FindChargedAdHocFeeByID(db, cahf.ID)
	if err != nil {
		return &ChargedAdHocFee{}, err
	}

	return cahf, nil
}

func (ChargedAdHocFee) TableName() string {
	return "fee_branch_relation"
}
