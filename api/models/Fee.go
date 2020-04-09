package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Fee struct {
	ID                    int                     `gorm:"primary_key;auto_increment" json:"id"`
	Name                  string                  `gorm:"not null;size:50" json:"name"`
	DefaultValue          float64                 `gorm:"not null;" json:"default_value"`
	UnitID                int                     `json:"unit_id"`
	Unit                  Unit                    `json:"unit"`
	FeeTypeID             int                     `gorm:"not null" json:"fee_type_id"`
	FeeType               FeeType                 `json:"fee_type"`
	FeeChargingCardStatus []FeeChargingCardStatus `json:"fee_charging_card_status"`
	FeeDormantDay         FeeDormantDay           `json:"fee_dormant_day"`
	ChargingPeriodID      int                     `json:"charging_period_id"`
	FeeChargingPeriod     FeeChargingPeriod       `json:"fee_charging_period"`
	CreatedAt             time.Time               `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             *time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt             *time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type FeeChargingPeriod struct {
	ID   int    `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"not null;" json:"name"`
}

type FeeChargingCardStatus struct {
	CardStatusID int        `gorm:"not null;" json:"card_status_id"`
	CardStatus   CardStatus `json:"card_status"`
}

type FeeDormantDay struct {
	DormantDay int `gorm:"not null;" json:"dormant_day"`
}

type FeeType struct {
	ID   int    `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"not null;" json:"name"`
}

func (fee *Fee) Prepare() {
	fee.Name = html.EscapeString(strings.TrimSpace(fee.Name))
	fee.CreatedAt = time.Now()
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

func (fee *Fee) FindIntialFees(db *gorm.DB) (*[]Fee, error) {
	var err error
	fees := []Fee{}
	err = db.Debug().Model(&Site{}).Unscoped().Order("id, created_at desc").Find(&fees).Error
	if err != nil {
		return &[]Fee{}, err
	}

	if len(fees) > 0 {
		for i, _ := range fees {
			unitErr := db.Debug().Model(&Unit{}).Unscoped().Where("id = ?", fees[i].UnitID).Order("id desc").Take(&fees[i].Unit).Error
			if unitErr != nil {
				return &[]Fee{}, err
			}

			feeTypeErr := db.Debug().Model(&FeeType{}).Unscoped().Where("id = ?", fees[i].FeeTypeID).Order("id desc").Take(&fees[i].FeeType).Error
			if feeTypeErr != nil {
				return &[]Fee{}, err
			}

			feeChargingPeriodErr := db.Debug().Model(&FeeChargingPeriod{}).Unscoped().Where("id = ?", fees[i].ChargingPeriodID).Order("id desc").Take(&fees[i].FeeChargingPeriod).Error
			if feeChargingPeriodErr != nil {
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
		unitErr := db.Debug().Model(&Unit{}).Unscoped().Where("id = ?", fee.UnitID).Order("id desc").Take(&fee.Unit).Error
		if unitErr != nil {
			return &Fee{}, err
		}

		feeTypeErr := db.Debug().Model(&FeeType{}).Unscoped().Where("id = ?", fee.FeeTypeID).Order("id desc").Take(&fee.FeeType).Error
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

		fddErr := db.Debug().Model(&FeeDormantDay{}).Where("fee_id = ?", fee.ID).Order("id desc").Take(&fee.FeeDormantDay).Error
		if fddErr != nil {
			return &Fee{}, err
		}

		feeChargingPeriodErr := db.Debug().Model(&FeeChargingPeriod{}).Unscoped().Where("id = ?", fee.ChargingPeriodID).Order("id desc").Take(&fee.FeeChargingPeriod).Error
		if feeChargingPeriodErr != nil {
			return &Fee{}, err
		}
	}

	return fee, nil
}

func (fee *Fee) CreateAdHocFee(db *gorm.DB) (*Fee, error) {
	var err error
	err = db.Debug().Model(&Fee{}).Create(&fee).Error
	if err != nil {
		return &Fee{}, err
	}

	return fee, nil
}

func (fee *Fee) UpdateAdHocFee(db *gorm.DB) (*Fee, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&Fee{}).Where("id = ?", fee.ID).Updates(
		Fee{
			Name:             fee.Name,
			DefaultValue:     fee.DefaultValue,
			UnitID:           fee.UnitID,
			FeeTypeID:        fee.FeeTypeID,
			ChargingPeriodID: fee.ChargingPeriodID,
			UpdatedAt:        &dateTimeNow,
		}).Error

	if err != nil {
		return &Fee{}, err
	}

	return fee, nil
}

func (fee *Fee) DeactivateAdHocFeeLater(db *gorm.DB) (int64, error) {
	var err error
	err = db.Debug().Model(&Fee{}).Unscoped().Where("id = ?", fee.ID).Updates(
		Fee{
			DeletedAt: fee.DeletedAt,
		}).Error

	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (Fee) TableName() string {
	return "fee"
}

func (FeeType) TableName() string {
	return "fee_type"
}

func (FeeChargingCardStatus) TableName() string {
	return "fee_card_status_relation"
}

func (FeeChargingPeriod) TableName() string {
	return "fee_charging_period"
}

func (FeeDormantDay) TableName() string {
	return "fee_dormant_day_relation"
}
