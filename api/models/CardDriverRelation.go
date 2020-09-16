package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

//CardDriverRelation => Card Number
type CardDriverRelation struct {
	DriverID  uint64     `gorm:"primary_key;column:card_holder_id" json:"driver_id"`
	CardID    string     `gorm:"primary_key;" json:"card_id"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

//Prepare => Prepare Query
func (cdr *CardDriverRelation) Prepare() {
	cdr.CreatedAt = time.Now()
}

//Validate => Validate given request body
func (cdr *CardDriverRelation) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if cdr.DriverID == 0 {
		err = errors.New("CC ID field is required")
		errorMessages["cc_id"] = err.Error()
	}

	if cdr.CardID == "" {
		err = errors.New("Card ID field is required")
		errorMessages["card_id"] = err.Error()
	}

	return errorMessages
}

//CreateCardDriverRelation => Create Card Perso
func (cdr *CardDriverRelation) CreateCardDriverRelation(db *gorm.DB) (*CardDriverRelation, error) {
	var err error
	err = db.Debug().Model(&CardDriverRelation{}).Create(&cdr).Error
	if err != nil {
		return &CardDriverRelation{}, err
	}

	return cdr, nil
}

//TableName => Define Table
func (CardDriverRelation) TableName() string {
	return "CardHolder_MemberCard_Relation"
}
