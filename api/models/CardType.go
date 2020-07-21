package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type CardType struct {
	ID                       uint64     `gorm:"primary_key;auto_increment;column:card_type_id" json:"id"`
	Code                     string     `gorm:"not null;column:card_type_code" json:"code"`
	Prefix                   string     `gorm:"not null;column:card_type_prefix" json:"prefix"`
	Name                     string     `gorm:"not null;size:100;column:card_type_name" json:"name"`
	OfflineFreqLimitPerDay   *uint64    `json:"offline_freq_limit_per_day"`
	OfflineValueLimitPerDay  *float32   `json:"offline_value_limit_per_day"`
	OfflineVolumeLimitPerDay *float32   `json:"offline_volume_limit_per_day"`
	StartedAt                *time.Time `json:"started_at"`
	EndedAt                  *time.Time `json:"ended_at"`
	CreatedAt                time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt                *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (ct *CardType) Prepare() {
	ct.Code = html.EscapeString(strings.TrimSpace(ct.Code))
	ct.Name = html.EscapeString(strings.TrimSpace(ct.Name))
	ct.Prefix = html.EscapeString(strings.TrimSpace(ct.Prefix))
	ct.CreatedAt = time.Now()
}

func (ct *CardType) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if ct.Code == "" {
		err = errors.New("Code field is required")
		errorMessages["code"] = err.Error()
	}

	if ct.Name == "" {
		err = errors.New("Name field is required")
		errorMessages["name"] = err.Error()
	}

	if ct.Prefix == "" {
		err = errors.New("Price field is required")
		errorMessages["prefix"] = err.Error()
	}

	if ct.OfflineFreqLimitPerDay == nil {
		err = errors.New("Offline frequency limit per day field is required")
		errorMessages["offline_freq_limit_per_day"] = err.Error()
	}

	if ct.OfflineValueLimitPerDay == nil {
		err = errors.New("Offline value limit per day field is required")
		errorMessages["offline_value_limit_per_day"] = err.Error()
	}

	if ct.OfflineVolumeLimitPerDay == nil {
		err = errors.New("Offline volume limit per day field is required")
		errorMessages["offline_volume_limit_per_day"] = err.Error()
	}

	if ct.StartedAt == nil {
		err = errors.New("Started at field is required")
		errorMessages["started_at"] = err.Error()
	}

	return errorMessages
}

func (ct *CardType) FindCardTypes(db *gorm.DB) (*[]CardType, error) {
	var err error
	ps := []CardType{}
	err = db.Debug().Model(&CardType{}).
		Order("card_type_id, created_at desc").
		Find(&ps).Error

	if err != nil {
		return &[]CardType{}, err
	}

	return &ps, nil
}

func (ct *CardType) FindCardTypeByID(db *gorm.DB, cardTypeID uint64) (*CardType, error) {
	var err error
	err = db.Debug().Model(&CardType{}).Unscoped().
		Where("card_type_id = ?", cardTypeID).
		Order("created_at desc").
		Take(&ct).Error

	if err != nil {
		return &CardType{}, err
	}

	return ct, nil
}

func (ct *CardType) CreateCardType(db *gorm.DB) (*CardType, error) {
	var err error
	err = db.Debug().Model(&CardType{}).Create(&ct).Error
	if err != nil {
		return &CardType{}, err
	}

	//Select created fee
	_, err = ct.FindCardTypeByID(db, ct.ID)
	if err != nil {
		return &CardType{}, err
	}

	return ct, nil
}

func (ct *CardType) UpdateCardType(db *gorm.DB) (*CardType, error) {
	var err error
	dateTimeNow := time.Now()

	//Update the data
	err = db.Debug().Model(&ct).Updates(
		map[string]interface{}{
			"card_type_code":               ct.Code,
			"card_type_name":               ct.Name,
			"card_type_prefix":             ct.Prefix,
			"offline_freq_limit_per_day":   ct.OfflineFreqLimitPerDay,
			"offline_volume_limit_per_day": ct.OfflineVolumeLimitPerDay,
			"offline_value_limit_per_day":  ct.OfflineValueLimitPerDay,
			"started_at":                   ct.StartedAt,
			"ended_at":                     ct.EndedAt,
			"updated_at":                   dateTimeNow,
		}).Error

	if err != nil {
		return &CardType{}, err
	}

	//Select updated product
	_, err = ct.FindCardTypeByID(db, ct.ID)
	if err != nil {
		return &CardType{}, err
	}

	return ct, nil
}

func (CardType) TableName() string {
	return "mstCardType"
}
