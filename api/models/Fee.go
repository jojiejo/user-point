package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Fee struct {
	ID        int        `gorm:"primary_key;auto_increment" json:"id"`
	Name      string     `gorm:"not null;" json:"name"`
	Value     int        `gorm:"not null;" json:"value"`
	FeeTypeID int        `gorm:"not null" json:"fee_type_id"`
	FeeType   FeeType    `json:"fee_type"`
	UnitID    int        `json:"unit_id"`
	Unit      Unit       `json:"unit"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type FeeType struct {
	ID   int    `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"not null;" json:"name"`
}

func (fee *Fee) Prepare() {
	fee.Name = html.EscapeString(strings.TrimSpace(fee.Name))
}

func (fee *Fee) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if fee.Name == "" {
		err = errors.New("Name field is required")
		errorMessages["name"] = err.Error()
	}

	if fee.UnitID < 1 {
		err = errors.New("Unit field is required")
		errorMessages["site"] = err.Error()
	}

	return errorMessages
}

func (fee *Fee) FindAllFees(db *gorm.DB) (*[]Fee, error) {
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
	}

	return fee, nil
}

func (Fee) TableName() string {
	return "fee"
}
