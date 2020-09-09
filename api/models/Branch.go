package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

//ShortenedBranch
type ShortenedBranch struct {
	SubCorporateID         int                             `gorm:"primary_key;auto_increment" json:"sub_corporate_id"`
	CCID                   int                             `gorm:"not null;" json:"cc_id"`
	MCMSID                 int                             `gorm:"not null;" json:"mcms_id"`
	PaddedMCMSID           string                          `json:"padded_mcms_id"`
	IsEnabledCardGroup     *bool                           `sql:"column:isEnabledCardGroup;" gorm:"not null;" json:"is_enabled_card_group"`
	ResProfileID           int                             `gorm:"not null" json:"res_profile_id"`
	GSAPCustomerMasterData ShortenedGSAPCustomerMasterData `json:"gsap_customer_master_data"`
	CreatedAt              time.Time                       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              *time.Time                      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt              *time.Time                      `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type Branch struct {
	SubCorporateID         int                    `gorm:"primary_key;auto_increment" json:"sub_corporate_id"`
	CCID                   int                    `gorm:"not null;" json:"cc_id"`
	MCMSID                 int                    `gorm:"not null;" json:"mcms_id"`
	PaddedMCMSID           string                 `json:"padded_mcms_id"`
	ResProfileID           int                    `gorm:"not null" json:"res_profile_id"`
	GSAPCustomerMasterData GSAPCustomerMasterData `json:"gsap_customer_master_data"`
	IsEnabledCardGroup     *bool                  `sql:"column:isEnabledCardGroup;" gorm:"not null;" json:"is_enabled_card_group"`
	CreatedAt              time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              *time.Time             `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt              *time.Time             `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (branch *Branch) ValidateConfiguration() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if branch.IsEnabledCardGroup == nil {
		err = errors.New("Enable card group flag field is required")
		errorMessages["required_enable_card_group"] = err.Error()
	}

	if branch.ResProfileID < 1 {
		err = errors.New("Restriction profile field is required")
		errorMessages["res_profile"] = err.Error()
	}

	return errorMessages
}

func (branch *ShortenedBranch) FindBranchByCCID(db *gorm.DB, CCID uint64) (*[]ShortenedBranch, error) {
	var err error
	branches := []ShortenedBranch{}
	err = db.Debug().Model(&Branch{}).Unscoped().Where("cc_id = ?", CCID).Order("created_at desc").Find(&branches).Error
	if err != nil {
		return &[]ShortenedBranch{}, err
	}

	if len(branches) > 0 {
		for i := range branches {
			customerDataErr := db.Debug().Model(&Payer{}).Unscoped().Where("mcms_id = ?", branches[i].MCMSID).Order("mcms_id desc").Take(&branches[i].GSAPCustomerMasterData).Error
			if customerDataErr != nil {
				return &[]ShortenedBranch{}, err
			}

			branches[i].PaddedMCMSID = fmt.Sprintf("%010v", strconv.Itoa(branches[i].MCMSID))
		}
	}

	return &branches, nil
}

func (branch *Branch) FindBranchByID(db *gorm.DB, branchID uint64) (*Branch, error) {
	var err error
	err = db.Debug().Model(&Branch{}).Unscoped().Where("sub_corporate_id = ?", branchID).Order("created_at desc").Take(&branch).Error
	if err != nil {
		return &Branch{}, err
	}

	customerDataErr := db.Debug().Model(&Payer{}).Unscoped().Where("mcms_id = ?", branch.MCMSID).Order("mcms_id desc").Take(&branch.GSAPCustomerMasterData).Error
	if customerDataErr != nil {
		return &Branch{}, err
	}

	branch.PaddedMCMSID = fmt.Sprintf("%010v", strconv.Itoa(branch.MCMSID))

	return branch, nil
}

func (branch *Branch) UpdateConfiguration(db *gorm.DB) (*Branch, error) {
	var err error
	dateTimeNow := time.Now()
	err = db.Debug().Model(&Branch{}).Where("sub_corporate_id = ?", branch.SubCorporateID).Updates(
		Branch{
			ResProfileID:       branch.ResProfileID,
			IsEnabledCardGroup: branch.IsEnabledCardGroup,
			UpdatedAt:          &dateTimeNow,
		}).Error

	if err != nil {
		return &Branch{}, err
	}

	return branch, nil
}

func (ShortenedBranch) TableName() string {
	return "mstSubCorporate"
}

func (Branch) TableName() string {
	return "mstSubCorporate"
}
