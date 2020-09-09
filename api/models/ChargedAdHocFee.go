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
	Description    string          `gorm:"size:50;" json:"description"`
	CreatedAt      time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt      *time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type ChargedAdHocFeeFromFile struct {
	MCMSID                 string                          `sql:"-" json:"mcms_id"`
	SubCorporateID         int                             `json:"sub_corporate_id"`
	SubCorporate           ShortenedBranch                 `json:"sub_corporate"`
	GSAPCustomerMasterData ShortenedGSAPCustomerMasterData `json:"gsap_customer_master_data"`
	FeeID                  uint64                          `json:"fee_id"`
	Fee                    ShortenedFee                    `json:"fee"`
	FeeName                string                          `sql:"-" json:"fee_name"`
	Value                  float32                         `json:"value"`
	Description            string                          `json:"description"`
	CreatedAt              time.Time                       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              time.Time                       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt              *time.Time                      `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type BulkChargeAdHocFee struct {
	FeeCollection []ChargedAdHocFeeFromFile `json:"fee_collection"`
}

func (cahf *ChargedAdHocFee) Prepare() {
	cahf.Description = html.EscapeString(strings.TrimSpace(cahf.Description))
	cahf.CreatedAt = time.Now()
}

func (bcadhf BulkChargeAdHocFee) Prepare() {
	if len(bcadhf.FeeCollection) > 0 {
		for i, _ := range bcadhf.FeeCollection {
			bcadhf.FeeCollection[i].MCMSID = html.EscapeString(strings.TrimSpace(bcadhf.FeeCollection[i].MCMSID))
			bcadhf.FeeCollection[i].FeeName = html.EscapeString(strings.TrimSpace(bcadhf.FeeCollection[i].FeeName))
			bcadhf.FeeCollection[i].Description = html.EscapeString(strings.TrimSpace(bcadhf.FeeCollection[i].Description))
		}
	}
}

func (bcadhf BulkChargeAdHocFee) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if len(bcadhf.FeeCollection) > 0 {
		var convertedLineNumber string

		for i, _ := range bcadhf.FeeCollection {
			convertedLineNumber = strconv.Itoa(i + 1)

			if bcadhf.FeeCollection[i].MCMSID == "" {
				err = errors.New("MCMS ID field in line " + convertedLineNumber + " is required")
				errorMessages["mcms_id_line_"+convertedLineNumber] = err.Error()
			}

			if bcadhf.FeeCollection[i].FeeName == "" {
				err = errors.New("Fee Name field in line " + convertedLineNumber + " is required")
				errorMessages["fee_name_line_"+convertedLineNumber] = err.Error()
			}

			if bcadhf.FeeCollection[i].Value < 1 {
				err = errors.New("Value field in line " + convertedLineNumber + " is required")
				errorMessages["value_line_"+convertedLineNumber] = err.Error()
			}
		}
	}

	return errorMessages
}

func (bcadhf BulkChargeAdHocFee) ChargePrepare() {
	if len(bcadhf.FeeCollection) > 0 {
		for i, _ := range bcadhf.FeeCollection {
			bcadhf.FeeCollection[i].MCMSID = html.EscapeString(strings.TrimSpace(bcadhf.FeeCollection[i].MCMSID))
			bcadhf.FeeCollection[i].FeeName = html.EscapeString(strings.TrimSpace(bcadhf.FeeCollection[i].FeeName))
			bcadhf.FeeCollection[i].Description = html.EscapeString(strings.TrimSpace(bcadhf.FeeCollection[i].Description))
			bcadhf.FeeCollection[i].CreatedAt = time.Now()
		}
	}
}

func (bcadhf BulkChargeAdHocFee) ChargeValidate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if len(bcadhf.FeeCollection) > 0 {
		var convertedLineNumber string

		for i, _ := range bcadhf.FeeCollection {
			convertedLineNumber = strconv.Itoa(i + 1)

			if bcadhf.FeeCollection[i].SubCorporateID < 1 {
				err = errors.New("Sub account field in line " + convertedLineNumber + " is required")
				errorMessages["mcms_id_line_"+convertedLineNumber] = err.Error()
			}

			if bcadhf.FeeCollection[i].FeeID < 1 {
				err = errors.New("Fee field in line " + convertedLineNumber + " is required")
				errorMessages["fee_name_line_"+convertedLineNumber] = err.Error()
			}

			if bcadhf.FeeCollection[i].Value < 1 {
				err = errors.New("Value field in line " + convertedLineNumber + " is required")
				errorMessages["value_line_"+convertedLineNumber] = err.Error()
			}
		}
	}

	return errorMessages
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

	return errorMessages
}

func (cahf *ChargedAdHocFee) FindAdHocFees(db *gorm.DB) (*[]ChargedAdHocFee, error) {
	var err error
	cahfs := []ChargedAdHocFee{}
	err = db.Debug().Model(&ChargedAdHocFee{}).Unscoped().
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

func (bcadhf *BulkChargeAdHocFee) BulkCheckAdHocFee(db *gorm.DB) (*BulkChargeAdHocFee, map[string]string) {
	var err error
	var errorMessages = make(map[string]string)

	if len(bcadhf.FeeCollection) > 0 {
		var convertedMCMSID int
		var convertedLineNumber string

		for i := range bcadhf.FeeCollection {
			convertedLineNumber = strconv.Itoa(i + 1)
			convertedMCMSID, err = strconv.Atoi(bcadhf.FeeCollection[i].MCMSID)
			if err != nil {
				err = errors.New("MCMS ID : " + bcadhf.FeeCollection[i].MCMSID + " could not be processed")
				errorMessages["unprocessed_mcms_id_line_"+convertedLineNumber] = err.Error()
			}

			customerDataErr := db.Debug().Model(&ShortenedGSAPCustomerMasterData{}).Unscoped().Where("mcms_id = ?", convertedMCMSID).Order("mcms_id desc").Take(&bcadhf.FeeCollection[i].GSAPCustomerMasterData).Error
			if customerDataErr != nil {
				err = errors.New("MCMS ID : " + bcadhf.FeeCollection[i].MCMSID + " is invalid")
				errorMessages["mcms_id_line_"+convertedLineNumber] = err.Error()
			}

			subCorporateDataErr := db.Debug().Model(&ShortenedBranch{}).Unscoped().Where("mcms_id = ?", convertedMCMSID).Order("sub_corporate_id desc").Take(&bcadhf.FeeCollection[i].SubCorporate).Error
			if subCorporateDataErr != nil {
				err = errors.New("MCMS ID : " + bcadhf.FeeCollection[i].MCMSID + " is not registered as sub account")
				errorMessages["not_registered_mcms_id_line_"+convertedLineNumber] = err.Error()
			}

			bcadhf.FeeCollection[i].SubCorporateID = bcadhf.FeeCollection[i].SubCorporate.SubCorporateID

			feeDataErr := db.Debug().Model(&ShortenedFee{}).Unscoped().Select("id").Where("name = ?", bcadhf.FeeCollection[i].FeeName).Order("id desc").Take(&bcadhf.FeeCollection[i].Fee).Error
			if feeDataErr != nil {
				err = errors.New("Fee Name : " + bcadhf.FeeCollection[i].FeeName + " is invalid")
				errorMessages["fee_name_line_"+convertedLineNumber] = err.Error()
			}

			bcadhf.FeeCollection[i].FeeID = bcadhf.FeeCollection[i].Fee.ID
		}
	}

	return bcadhf, errorMessages
}

func (bcadhf *BulkChargeAdHocFee) BulkChargeAdHocFee(db *gorm.DB) (*BulkChargeAdHocFee, map[string]string) {
	var err error
	var errorMessages = make(map[string]string)

	if len(bcadhf.FeeCollection) > 0 {
		for i, _ := range bcadhf.FeeCollection {
			var convertedLineNumber string
			convertedLineNumber = strconv.Itoa(i + 1)

			err = db.Debug().Model(&ChargedAdHocFeeFromFile{}).Create(&bcadhf.FeeCollection[i]).Error
			if err != nil {
				errorMessages["line_"+convertedLineNumber] = err.Error()
			}
		}
	}

	return bcadhf, errorMessages
}

func (ChargedAdHocFee) TableName() string {
	return "fee_branch_relation"
}

func (ChargedAdHocFeeFromFile) TableName() string {
	return "fee_branch_relation"
}
