package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Fee struct {
	ID                    int                     `gorm:"primary_key;auto_increment" json:"id"`
	FeeNameID             int                     `gorm:"not null;" json:"fee_name_id"`
	FeeName               FeeName                 `json:"fee_name"`
	Value                 float64                 `gorm:"not null;" json:"value"`
	UnitID                int                     `json:"unit_id"`
	Unit                  Unit                    `json:"unit"`
	IsDefault             *bool                   `json:"is_default"`
	FeeChargingCardStatus []FeeChargingCardStatus `json:"fee_charging_card_status"`
	FeeChargingPeriod     FeeChargingPeriod       `json:"fee_charging_period"`
	CreatedAt             time.Time               `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             *time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt             *time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type FeeChargingPeriod struct {
	FeeChargingPeriodID int            `gorm:"primary_key;auto_increment" json:"fee_charging_period_id"`
	ChargingPeriod      ChargingPeriod `gorm:"not null;" json:"charging_period"`
}

type ChargingPeriod struct {
	ID   int    `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"not null;" json:"name"`
}

type FeeChargingCardStatus struct {
	CardStatusID int        `gorm:"not null;" json:"card_status_id"`
	CardStatus   CardStatus `json:"card_status"`
}

type FeeType struct {
	ID   int    `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"not null;" json:"name"`
}

type FeeName struct {
	ID        int        `gorm:"primary_key;auto_increment" json:"id"`
	Name      string     `gorm:"not null;" json:"name"`
	FeeTypeID int        `gorm:"not null" json:"fee_type_id"`
	FeeType   FeeType    `json:"fee_type"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (fee *Fee) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if fee.UnitID < 1 {
		err = errors.New("Unit field is required")
		errorMessages["site"] = err.Error()
	}

	return errorMessages
}

func (fn *FeeName) FindAllFeeNames(db *gorm.DB) (*[]FeeName, error) {
	var err error
	fns := []FeeName{}
	err = db.Debug().Model(&FeeName{}).Unscoped().Order("id, created_at desc").Find(&fns).Error
	if err != nil {
		return &[]FeeName{}, err
	}

	if len(fns) > 0 {
		for i, _ := range fns {
			feeTypeErr := db.Debug().Model(&FeeType{}).Unscoped().Where("id = ?", fns[i].FeeTypeID).Order("id desc").Take(&fns[i].FeeType).Error
			if feeTypeErr != nil {
				return &[]FeeName{}, err
			}
		}
	}

	return &fns, nil
}

func (fee *Fee) FindIntialFees(db *gorm.DB) (*[]Fee, error) {
	var err error
	fees := []Fee{}
	err = db.Debug().Model(&Site{}).Unscoped().Where("is_default = 1").Order("id, created_at desc").Find(&fees).Error
	if err != nil {
		return &[]Fee{}, err
	}

	if len(fees) > 0 {
		for i, _ := range fees {
			fnErr := db.Debug().Model(&FeeName{}).Unscoped().Where("id = ?", fees[i].FeeNameID).Order("id desc").Take(&fees[i].FeeName).Error
			if fnErr != nil {
				return &[]Fee{}, err
			}

			unitErr := db.Debug().Model(&Unit{}).Unscoped().Where("id = ?", fees[i].UnitID).Order("id desc").Take(&fees[i].Unit).Error
			if unitErr != nil {
				return &[]Fee{}, err
			}

			feeTypeErr := db.Debug().Model(&FeeType{}).Unscoped().Where("id = ?", fees[i].FeeName.FeeTypeID).Order("id desc").Take(&fees[i].FeeName.FeeType).Error
			if feeTypeErr != nil {
				return &[]Fee{}, err
			}
		}
	}

	return &fees, nil
}

func (fee *Fee) FindFeeByID(db *gorm.DB, feeID uint64) (*Fee, error) {
	var err error
	err = db.Debug().Model(&Fee{}).Unscoped().Where("id = ?", feeID).Order("created_at desc").Take(&fee).Error
	if err != nil {
		return &Fee{}, err
	}

	if fee.ID != 0 {
		fnErr := db.Debug().Model(&FeeName{}).Unscoped().Where("id = ?", fee.FeeNameID).Order("id desc").Take(&fee.FeeName).Error
		if fnErr != nil {
			return &Fee{}, err
		}

		unitErr := db.Debug().Model(&Unit{}).Unscoped().Where("id = ?", fee.UnitID).Order("id desc").Take(&fee.Unit).Error
		if unitErr != nil {
			return &Fee{}, err
		}

		feeTypeErr := db.Debug().Model(&FeeType{}).Unscoped().Where("id = ?", fee.FeeName.FeeTypeID).Order("id desc").Take(&fee.FeeName.FeeType).Error
		if feeTypeErr != nil {
			return &Fee{}, err
		}

		fccsErr := db.Debug().Model(&FeeChargingCardStatus{}).Where("fee_id = ?", fee.ID).Order("id desc").Find(&fee.FeeChargingCardStatus).Error
		if fccsErr != nil {
			return &Fee{}, err
		}

		if len(fee.FeeChargingCardStatus) > 0 {
			for i, _ := range fee.FeeChargingCardStatus {
				cardStatusErr := db.Debug().Model(&CardStatus{}).Where("id = ?", fee.FeeChargingCardStatus[i].CardStatusID).Order("id desc").Take(&fee.FeeChargingCardStatus[i].CardStatus).Error
				if cardStatusErr != nil {
					return &Fee{}, err
				}
			}
		}

		fcpErr := db.Debug().Model(&FeeChargingPeriod{}).Where("fee_id = ?", fee.ID).Order("id desc").Find(&fee.FeeChargingPeriod).Error
		if fcpErr != nil {
			return &Fee{}, err
		}

		cpErr := db.Debug().Model(&ChargingPeriod{}).Where("id = ?", fee.ID).Order("id desc").Find(&fee.FeeChargingPeriod.ChargingPeriod).Error
		if cpErr != nil {
			return &Fee{}, err
		}
	}

	return fee, nil
}

func (Fee) TableName() string {
	return "fee"
}

func (FeeType) TableName() string {
	return "fee_type"
}

func (FeeName) TableName() string {
	return "fee_name"
}

func (FeeChargingCardStatus) TableName() string {
	return "fee_card_status_relation"
}

func (FeeChargingPeriod) TableName() string {
	return "fee_charging_period_relation"
}

func (ChargingPeriod) TableName() string {
	return "fee_charging_period"
}
