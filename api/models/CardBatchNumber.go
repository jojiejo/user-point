package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

//CardBatchNumber => Card Batch Number
type CardBatchNumber struct {
	CCID      int        `gorm:"primary_key;" json:"cc_id"`
	LastBatch int        `json:"last_batch"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

//Prepare => Prepare null data
func (cbn *CardBatchNumber) Prepare() {
	cbn.CreatedAt = time.Now()
}

//Validate => Validate request
func (cbn *CardBatchNumber) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if cbn.CCID <= 0 {
		err = errors.New("CC ID field is required")
		errorMessages["cc_id"] = err.Error()
	}

	return errorMessages
}

//FindBatchNumberByCCID => Find Card Number By CC ID
func (cbn *CardBatchNumber) FindBatchNumberByCCID(db *gorm.DB, ccID int) (*CardBatchNumber, error) {
	var err error

	err = db.Debug().Model(&CardBatchNumber{}).
		Order("created_at desc").
		Find(&cbn).Error

	if err != nil {
		return &CardBatchNumber{}, err
	}

	return cbn, nil
}

//CreateCardBatchNumber => Create Card Batch Number
func (cbn *CardBatchNumber) CreateCardBatchNumber(db *gorm.DB) (*CardBatchNumber, error) {
	var err error
	err = db.Debug().Model(&CardNumber{}).Create(&cbn).Error
	if err != nil {
		return &CardBatchNumber{}, err
	}

	_, err = cbn.FindBatchNumberByCCID(db, cbn.CCID)
	if err != nil {
		return &CardBatchNumber{}, err
	}

	return cbn, nil
}

//UpdateCardBatchNumber => Update Card Batch Number
func (cbn *CardBatchNumber) UpdateCardBatchNumber(db *gorm.DB) (*CardBatchNumber, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&cbn).Updates(
		map[string]interface{}{
			"last_batch": cbn.LastBatch,
			"updated_at": dateTimeNow,
		}).Error

	if err != nil {
		return &CardBatchNumber{}, err
	}

	_, err = cbn.FindBatchNumberByCCID(db, cbn.CCID)
	if err != nil {
		return &CardBatchNumber{}, err
	}

	return cbn, nil
}

//TableName => Define Table
func (CardBatchNumber) TableName() string {
	return "mstCardInventory"
}
