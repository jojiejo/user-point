package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

//CardNumber => Card Number
type CardNumber struct {
	CountryCode string     `gorm:"primary_key;" json:"country_code"`
	BankCode    string     `gorm:"primary_key;" json:"bank_code"`
	LastCardNo  int        `gorm:"not null;" json:"last_card_no"`
	CardTypeID  int        `gorm:"primary_key;" json:"card_type_id`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

//Prepare => Prepare null data
func (cn *CardNumber) Prepare() {
	cn.CreatedAt = time.Now()
}

//Validate => Validate request
func (cn *CardNumber) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if cn.CountryCode == "" {
		err = errors.New("Country Code field is required")
		errorMessages["country_code"] = err.Error()
	}

	if cn.BankCode == "" {
		err = errors.New("Bank Code field is required")
		errorMessages["bank_code"] = err.Error()
	}

	if cn.CardTypeID <= 0 {
		err = errors.New("Card Type field is required")
		errorMessages["card_type"] = err.Error()
	}

	return errorMessages
}

//FindCardNumberByKey => Find Card Number By Key
func (cn *CardNumber) FindCardNumberByKey(db *gorm.DB, countryCode string, bankCode string, cardTypeID int) (*CardNumber, error) {
	var err error

	err = db.Debug().Model(&CardNumber{}).
		Order("created_at desc").
		Find(&cn).Error

	if err != nil {
		return &CardNumber{}, err
	}

	return cn, nil
}

//CreateCardNumber => Create Card Number
func (cn *CardNumber) CreateCardNumber(db *gorm.DB) (*CardNumber, error) {
	var err error
	err = db.Debug().Model(&CardNumber{}).Create(&cn).Error
	if err != nil {
		return &CardNumber{}, err
	}

	_, err = cn.FindCardNumberByKey(db, cn.CountryCode, cn.BankCode, cn.CardTypeID)
	if err != nil {
		return &CardNumber{}, err
	}

	return cn, nil
}

//UpdateCardNumber => Update Card Number
func (cn *CardNumber) UpdateCardNumber(db *gorm.DB) (*CardNumber, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&cn).Updates(
		map[string]interface{}{
			"last_card_no": cn.LastCardNo,
			"updated_at":   dateTimeNow,
		}).Error

	if err != nil {
		return &CardNumber{}, err
	}

	_, err = cn.FindCardNumberByKey(db, cn.CountryCode, cn.BankCode, cn.CardTypeID)
	if err != nil {
		return &CardNumber{}, err
	}

	return cn, nil
}

//TableName => Define Table
func (CardNumber) TableName() string {
	return "mstCardNumber"
}
