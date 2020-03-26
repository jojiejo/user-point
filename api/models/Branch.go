package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ShortenedBranch struct {
	SubCorporateID         int                             `gorm:"primary_key;auto_increment" json:"sub_corporate_id"`
	CCID                   int                             `gorm:"not null;" json:"cc_id"`
	MCMSID                 int                             `gorm:"not null;" json:"mcms_id"`
	GSAPCustomerMasterData ShortenedGSAPCustomerMasterData `json:"gsap_customer_master_data"`
	CreatedAt              time.Time                       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              *time.Time                      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt              *time.Time                      `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type Branch struct {
	SubCorporateID         int                    `gorm:"primary_key;auto_increment" json:"sub_corporate_id"`
	CCID                   int                    `gorm:"not null;" json:"cc_id"`
	MCMSID                 int                    `gorm:"not null;" json:"mcms_id"`
	GSAPCustomerMasterData GSAPCustomerMasterData `json:"gsap_customer_master_data"`
	CreatedAt              time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              *time.Time             `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt              *time.Time             `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (branch *ShortenedBranch) FindBranchByCCID(db *gorm.DB, CCID uint64) (*[]ShortenedBranch, error) {
	var err error
	branches := []ShortenedBranch{}
	err = db.Debug().Model(&Branch{}).Unscoped().Where("cc_id = ?", CCID).Order("created_at desc").Find(&branches).Error
	if err != nil {
		return &[]ShortenedBranch{}, err
	}

	if len(branches) > 0 {
		for i, _ := range branches {
			customerDataErr := db.Debug().Model(&Payer{}).Unscoped().Where("mcms_id = ?", branches[i].MCMSID).Order("mcms_id desc").Take(&branches[i].GSAPCustomerMasterData).Error
			if customerDataErr != nil {
				return &[]ShortenedBranch{}, err
			}
		}
	}

	return &branches, nil
}

func (branch *Branch) FindBranchByID(db *gorm.DB, branchID uint64) (*Branch, error) {
	var err error
	err = db.Debug().Model(&Branch{}).Unscoped().Where("sub_corporate_id = ?", branchID).Order("created_at desc").Find(&branch).Error
	if err != nil {
		return &Branch{}, err
	}

	customerDataErr := db.Debug().Model(&Payer{}).Unscoped().Where("mcms_id = ?", branch.MCMSID).Order("mcms_id desc").Take(&branch.GSAPCustomerMasterData).Error
	if customerDataErr != nil {
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
