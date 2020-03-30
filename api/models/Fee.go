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
	UnitID    int        `json:"unit_id"`
	Unit      Unit       `json:"unit"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
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
			err := db.Debug().Model(&City{}).Unscoped().Where("id = ?", fees[i].UnitID).Order("id desc").Take(&fees[i].Unit).Error
			if err != nil {
				return &[]Fee{}, err
			}
		}
	}

	return &fees, nil
}

func (Fee) TableName() string {
	return "fee"
}
