package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type ShortenedFee struct {
	ID           uint64  `gorm:"primary_key;auto_increment" json:"id"`
	Name         string  `gorm:"not null;size:50" json:"name"`
	DefaultValue float64 `gorm:"not null;" json:"default_value"`
}

type Fee struct {
	ID                    uint64            `gorm:"primary_key;auto_increment" json:"id"`
	Name                  string            `gorm:"not null;size:50" json:"name"`
	DefaultValue          float64           `gorm:"not null;" json:"default_value"`
	UnitID                int               `json:"unit_id"`
	Unit                  Unit              `json:"unit"`
	FeeTypeID             int               `gorm:"not null" json:"fee_type_id"`
	FeeType               FeeType           `json:"fee_type"`
	FeeChargingCardStatus []CardStatus      `gorm:"many2many:fee_card_status;association_autoupdate:false;association_jointable_foreignkey:fee_id;association_jointable_foreignkey:card_status_id" json:"fee_charging_card_status"`
	FeeDormantDay         FeeDormantDay     `gorm:"foreignkey:FeeID;association_foreignkey:ID" json:"fee_dormant_day"`
	ChargingPeriodID      int               `json:"charging_period_id"`
	ChargingPeriod        FeeChargingPeriod `json:"charging_period"`
	CreatedAt             time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt             *time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type FeeChargingPeriod struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Name      string     `gorm:"not null;" json:"name"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

type FeeDormantDay struct {
	FeeID      uint64 `gorm:"not null;" json:"fee_id"`
	DormantDay uint64 `gorm:"not null;" json:"dormant_day"`
}

type FeeType struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Name      string     `gorm:"not null;" json:"name"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
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
		errorMessages["unit"] = err.Error()
	}

	return errorMessages
}

func (fee *Fee) FindIntialFees(db *gorm.DB) (*[]Fee, error) {
	var err error
	fees := []Fee{}
	err = db.Debug().Model(&Site{}).Unscoped().
		Preload("Unit").
		Preload("FeeType").
		Preload("FeeChargingCardStatus").
		Preload("ChargingPeriod").
		Preload("FeeDormantDay").
		Order("id, created_at desc").
		Find(&fees).Error

	if err != nil {
		return &[]Fee{}, err
	}

	return &fees, nil
}

func (fee *Fee) FindFeeByID(db *gorm.DB, feeID uint64) (*Fee, error) {
	var err error
	err = db.Debug().Model(&Fee{}).Unscoped().
		Preload("Unit").
		Preload("FeeType").
		Preload("FeeChargingCardStatus").
		Preload("ChargingPeriod").
		Preload("FeeDormantDay").
		Where("id = ?", feeID).
		Order("created_at desc").
		Take(&fee).Error

	if err != nil {
		return &Fee{}, err
	}

	return fee, nil
}

func (fee *Fee) CreateFee(db *gorm.DB) (*Fee, error) {
	var err error
	err = db.Debug().Model(&Fee{}).Create(&fee).Error
	if err != nil {
		return &Fee{}, err
	}

	//Select created fee
	_, err = fee.FindFeeByID(db, fee.ID)
	if err != nil {
		return &Fee{}, err
	}

	return fee, nil
}

func (fee *Fee) UpdateFee(db *gorm.DB) (*Fee, error) {
	var err error
	dateTimeNow := time.Now()

	//Update the data
	err = db.Debug().Model(&Fee{}).Where("id = ?", fee.ID).Updates(
		Fee{
			Name:             fee.Name,
			DefaultValue:     fee.DefaultValue,
			UnitID:           fee.UnitID,
			FeeTypeID:        fee.FeeTypeID,
			ChargingPeriodID: fee.ChargingPeriodID,
			UpdatedAt:        dateTimeNow,
		}).Error
	if err != nil {
		return &Fee{}, err
	}

	//Update dormant day
	err = db.Debug().Model(&fee).Where("id = ?", fee.ID).Association("FeeDormantDay").Append(fee.FeeDormantDay).Error
	if err != nil {
		return &Fee{}, err
	}

	//Update card status
	err = db.Debug().Model(&fee).Where("id = ?", fee.ID).Association("FeeChargingCardStatus").Replace(fee.FeeChargingCardStatus).Error
	if err != nil {
		return &Fee{}, err
	}

	//Select created fee
	_, err = fee.FindFeeByID(db, fee.ID)
	if err != nil {
		return &Fee{}, err
	}

	return fee, nil
}

func (fee *Fee) DeactivateFeeLater(db *gorm.DB) (int64, error) {
	var err error
	err = db.Debug().Model(&fee).Unscoped().Updates(
		Fee{
			DeletedAt: fee.DeletedAt,
		}).Error

	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (ShortenedFee) TableName() string {
	return "fee"
}

func (Fee) TableName() string {
	return "fee"
}

func (FeeType) TableName() string {
	return "fee_type"
}

func (FeeChargingPeriod) TableName() string {
	return "fee_charging_period"
}

func (FeeDormantDay) TableName() string {
	return "fee_dormant_day"
}
