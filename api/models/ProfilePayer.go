package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//ProfilePayer => Account Profile
type ProfilePayer struct {
	ResProfileID uint64        `gorm:"primary_key;" json:"res_profile_id"`
	ResProfile   ProfileMaster `gorm:"foreignkey:ResProfileID;references:ResProfileID" json:"res_profile"`
	CCID         uint64        `gorm:"primary_key;" json:"cc_id"`
	CreatedAt    time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    *time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

//Prepare => Prepare Restriction Profile
func (pp *ProfilePayer) Prepare() {
	pp.CreatedAt = time.Now()
}

//FindProfilePayerByCCID => Find Profile Payer By CC ID
func (pp *ProfilePayer) FindProfilePayerByCCID(db *gorm.DB, CCID uint64) (*ProfilePayer, error) {
	var err error
	err = db.Debug().Model(&ProfilePayer{}).
		Preload("ResProfile").
		Preload("ResProfile.Velocity").
		Order("res_profile_id, created_at desc").
		Find(&pp).Error

	if err != nil {
		return &ProfilePayer{}, err
	}

	return pp, nil
}

//CreateProfilePayer => Create Profile Payer
func (pp *ProfilePayer) CreateProfilePayer(db *gorm.DB) (*ProfilePayer, error) {
	var err error
	err = db.Debug().Model(&ProfileMaster{}).Create(&pp).Error
	if err != nil {
		return &ProfilePayer{}, err
	}

	return pp, nil
}

//TableName => Define Table
func (ProfilePayer) TableName() string {
	return "CorporateRestriction"
}
