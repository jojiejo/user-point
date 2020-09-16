package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

//CardPerso => Card Number
type CardPerso struct {
	CCID      int        `gorm:"primary_key;" json:"cc_id"`
	CardID    string     `gorm:"primary_key;" json:"card_id"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

//Prepare => Prepare null data
func (cp *CardPerso) Prepare() {
	cp.CreatedAt = time.Now()
}

//Validate => Validate request
func (cp *CardPerso) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if cp.CardID == "" {
		err = errors.New("Card ID field is required")
		errorMessages["card_id"] = err.Error()
	}

	if cp.CCID <= 0 {
		err = errors.New("CC ID field is required")
		errorMessages["cc_id"] = err.Error()
	}

	return errorMessages
}

//FindCardPersos => Find Card Persos
func (cp *CardPerso) FindCardPersos(db *gorm.DB) (*CardPerso, error) {
	var err error

	err = db.Debug().Model(&CardPerso{}).
		Order("created_at desc").
		Find(&cp).Error

	if err != nil {
		return &CardPerso{}, err
	}

	return cp, nil
}

//CreateCardPerso => Create Card Perso
func (cp *CardPerso) CreateCardPerso(db *gorm.DB) (*CardPerso, error) {
	var err error
	err = db.Debug().Model(&CardPerso{}).Create(&cp).Error
	if err != nil {
		return &CardPerso{}, err
	}

	return cp, nil
}

//DeleteCardPerso => Delete Card Perso
func (cp *CardPerso) DeleteCardPerso(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&CardPerso{}).Where("card_id = ?", cp.CardID).Delete(&CardPerso{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

//TableName => Define Table
func (CardPerso) TableName() string {
	return "mstPersoFile"
}
