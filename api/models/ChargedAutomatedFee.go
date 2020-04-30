package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type ChargedAutomatedFee struct {
	ID                    uint64                        `sql:"index" gorm:"primary_key;auto_increment;" json:"id"`
	CCID                  uint64                        `gorm:"not null" json:"cc_id"`
	FeeID                 uint64                        `gorm:"not null;" json:"fee_id"`
	Fee                   ShortenedFee                  `json:"fee"`
	Value                 float32                       `gorm:"not null;" json:"value"`
	FeeDormantDay         ChargedAutomatedFeeDormantDay `gorm:"foreignkey:fee_payer_id;association_foreignkey:id;" json:"fee_dormant_day"`
	FeeChargingCardStatus []CardStatus                  `gorm:"many2many:fee_payer_card_status;association_autoupdate:false;jointable_foreignkey:fee_payer_id;association_jointable_foreignkey:fee_payer_id;association_jointable_foreignkey:card_status_id" json:"fee_charging_card_status"`
	ChargingPeriodID      int                           `json:"charging_period_id"`
	ChargingPeriod        FeeChargingPeriod             `json:"charging_period"`
	WaivedAt              *time.Time                    `gorm:"default:CURRENT_TIMESTAMP" json:"waived_at"`
	CreatedAt             time.Time                     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time                     `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt             *time.Time                    `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type ChargedAutomatedFeeChargingPeriod struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Name      string     `gorm:"not null;" json:"name"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

type ChargedAutomatedFeeDormantDay struct {
	FeePayerID uint64 `gorm:"primary_key;not null;" json:"fee_payer_id"`
	DormantDay uint64 `gorm:"not null;" json:"dormant_day"`
}

func (caf *ChargedAutomatedFee) Prepare() {
	caf.CreatedAt = time.Now()
}

func (caf *ChargedAutomatedFee) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if caf.FeeID < 1 {
		err = errors.New("Fee field is required")
		errorMessages["unit"] = err.Error()
	}

	if caf.CCID < 1 {
		err = errors.New("Payer field is required")
		errorMessages["unit"] = err.Error()
	}

	return errorMessages
}

func (caf *ChargedAutomatedFee) FindChargedAutomatedFeeByCCID(db *gorm.DB, CCID uint64) (*[]ChargedAutomatedFee, error) {
	var err error
	cafs := []ChargedAutomatedFee{}
	err = db.Debug().Model(&ChargedAutomatedFee{}).Unscoped().
		Preload("Fee").
		Where("cc_id = ?", CCID).
		Order("created_at desc").
		Find(&cafs).Error

	if err != nil {
		return &[]ChargedAutomatedFee{}, err
	}

	return &cafs, nil
}

func (caf *ChargedAutomatedFee) FindChargedAutomatedFeeByID(db *gorm.DB, relationID uint64) (*ChargedAutomatedFee, error) {
	var err error
	err = db.Debug().Model(&ChargedAutomatedFee{}).Unscoped().
		Preload("Fee").
		Where("id = ?", relationID).
		Order("created_at desc").
		Take(&caf).Error

	if err != nil {
		return &ChargedAutomatedFee{}, err
	}

	return caf, nil
}

func (caf *ChargedAutomatedFee) UpdateAutomatedFee(db *gorm.DB, relationID uint64) (*ChargedAutomatedFee, error) {
	var err error
	dateTimeNow := time.Now()
	err = db.Debug().Model(&ChargedAutomatedFee{}).Where("id = ?", relationID).Updates(
		ChargedAutomatedFee{
			Value:            caf.Value,
			ChargingPeriodID: caf.ChargingPeriodID,
			UpdatedAt:        dateTimeNow,
		}).Error

	if err != nil {
		return &ChargedAutomatedFee{}, err
	}

	//Update dormant day
	err = db.Debug().Model(&caf).Where("id = ?", caf.ID).Association("FeeDormantDay").Append(caf.FeeDormantDay).Error
	if err != nil {
		return &ChargedAutomatedFee{}, err
	}

	//Update card status
	err = db.Debug().Model(&caf).Where("id = ?", caf.ID).Association("FeeChargingCardStatus").Replace(caf.FeeChargingCardStatus).Error
	if err != nil {
		return &ChargedAutomatedFee{}, err
	}

	//Select created fee
	_, err = caf.FindChargedAutomatedFeeByID(db, caf.ID)
	if err != nil {
		return &ChargedAutomatedFee{}, err
	}

	return caf, nil
}

func (ChargedAutomatedFee) TableName() string {
	return "fee_payer_relation"
}

func (ChargedAutomatedFeeChargingPeriod) TableName() string {
	return "fee_payer_charging_period"
}

func (ChargedAutomatedFeeDormantDay) TableName() string {
	return "fee_payer_dormant_day"
}
