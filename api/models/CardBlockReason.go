package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//CardBlockReason => Card Block Reason
type CardBlockReason struct {
	ID        uint64     `gorm:"primary_key;auto_increment;" json:"id"`
	Name      string     `gorm:"not null;size:100;" json:"name"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

//FindCardBlockReasons => Find Card Block Reasons
func (cbr *CardBlockReason) FindCardBlockReasons(db *gorm.DB) (*[]CardBlockReason, error) {
	var err error
	cbrs := []CardBlockReason{}
	err = db.Debug().Model(&CardBlockReason{}).
		Order("id, created_at desc").
		Find(&cbrs).Error

	if err != nil {
		return &[]CardBlockReason{}, err
	}

	return &cbrs, nil
}

//FindCardBlockReasonByID => Find Card Block Reason By ID
func (cbr *CardBlockReason) FindCardBlockReasonByID(db *gorm.DB, reasonID uint64) (*CardBlockReason, error) {
	var err error
	err = db.Debug().Model(&CardBlockReason{}).Unscoped().
		Where("id = ?", reasonID).
		Order("created_at desc").
		Take(&cbr).Error

	if err != nil {
		return &CardBlockReason{}, err
	}

	return cbr, nil
}

//TableName => Define table
func (CardBlockReason) TableName() string {
	return "member_card_block_reason"
}
