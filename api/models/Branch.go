package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Branch struct {
	CCID                   int                             `gorm:"primary_key;auto_increment" json:"cc_id"`
	MCMSID                 int                             `gorm:"not null;" json:"mcms_id"`
	GSAPCustomerMasterData DisplayedGSAPCustomerMasterData `json:"gsap_customer_master_data"`
	CreatedAt              time.Time                       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              *time.Time                      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt              *time.Time                      `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (branch *Branch) FindBranchByCCID(db *gorm.DB, CCID uint64) (*[]Branch, error) {
	var err error
	branches := []Branch{}
	err = db.Debug().Model(&Branch{}).Unscoped().Where("cc_id = ?", CCID).Order("created_at desc").Find(&branches).Error
	if err != nil {
		return &[]Branch{}, err
	}

	if len(branches) > 0 {
		for i, _ := range branches {
			customerDataErr := db.Debug().Model(&Payer{}).Unscoped().Where("mcms_id = ?", branches[i].MCMSID).Order("mcms_id desc").Take(&branches[i].GSAPCustomerMasterData).Error
			if customerDataErr != nil {
				return &[]Branch{}, err
			}
		}
	}

	return &branches, nil
}

/*func (payer *Payer) FindPayerByCCID(db *gorm.DB, CCID uint64) (*[]Payer, error) {
	var err error
	payers := []Payer{}
	err = db.Debug().Model(&Payer{}).Unscoped().Where("cc_id = ?", CCID).Order("created_at desc").Find(&payers).Error
	if err != nil {
		return &[]Payer{}, err
	}

	if len(payers) > 0 {
		for i, _ := range payers {
			customerDataErr := db.Debug().Model(&Payer{}).Unscoped().Where("mcms_id = ?", payers[i].MCMSID).Order("mcms_id desc").Take(&payers[i].GSAPCustomerMasterData).Error
			if customerDataErr != nil {
				return &[]Payer{}, err
			}

			latestStatusErr := db.Debug().Model(&HistoricalPayerStatus{}).Where("cc_id = ?", payers[i].CCID).Order("created_at desc").Find(&payers[i].LatestPayerStatus).Error
			if latestStatusErr != nil {
				return &[]Payer{}, err
			}

			statusErr := db.Debug().Model(&PayerStatus{}).Where("id = ?", payers[i].LatestPayerStatus.PayerStatusID).Order("id desc").Take(&payers[i].LatestPayerStatus.PayerStatus).Error
			if statusErr != nil {
				return &[]Payer{}, err
			}
		}
	}

	return &payers, nil
}*/

func (Branch) TableName() string {
	return "mstSubCorporate"
}
