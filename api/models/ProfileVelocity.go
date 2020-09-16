package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

//ProfileVelocity => Restriction Profile
type ProfileVelocity struct {
	ResProfileID             uint64     `gorm:"primary_key" json:"res_profile_id"`
	FreqLimitPerDay          int        `json:"freq_limit_per_day"`
	FreqLimitPerWeek         int        `json:"freq_limit_per_week"`
	FreqLimitPerMonth        int        `json:"freq_limit_per_month"`
	VolumeLimitPerDay        float32    `json:"volume_limit_per_day"`
	VolumeLimitPerWeek       float32    `json:"volume_limit_per_week"`
	VolumeLimitPerMonth      float32    `json:"volume_limit_per_month"`
	ValueLimitPerDay         float32    `json:"value_limit_per_day"`
	ValueLimitPerWeek        float32    `json:"value_limit_per_week"`
	ValueLimitPerMonth       float32    `json:"value_limit_per_month"`
	PrepaidBalancePercentage float32    `json:"prepaid_balance_percentage"`
	CreatedAt                time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt                *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

//Prepare => Prepare Restriction Profile
func (pv *ProfileVelocity) Prepare() {
	pv.CreatedAt = time.Now()
}

//Prepare => Validate Restriction Profile
func (pv *ProfileVelocity) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if pv.FreqLimitPerDay == 0 {
		err = errors.New("Freq limit per day field is required")
		errorMessages["freq_limit_per_day"] = err.Error()
	}

	if pv.FreqLimitPerWeek == 0 {
		err = errors.New("Freq limit per week field is required")
		errorMessages["freq_limit_per_week"] = err.Error()
	}

	if pv.FreqLimitPerMonth == 0 {
		err = errors.New("Freq limit per month field is required")
		errorMessages["freq_limit_per_month"] = err.Error()
	}

	return errorMessages
}

//FindProfileVelocityByProfileID => Find Restriction Profile by Restriction Profile ID
func (pv *ProfileVelocity) FindProfileVelocityByProfileID(db *gorm.DB, resProfileID uint64) (*ProfileVelocity, error) {
	var err error
	err = db.Debug().Model(&ProfileVelocity{}).Unscoped().
		Where("res_profile_id = ?", resProfileID).
		Order("created_at desc").
		Take(&pv).Error

	if err != nil {
		return &ProfileVelocity{}, err
	}

	return pv, nil
}

//CreateProfileVelocity => Create Restriction Profile by ID
func (pv *ProfileVelocity) CreateProfileVelocity(db *gorm.DB) (*ProfileVelocity, error) {
	var err error
	err = db.Debug().Model(&ProfileVelocity{}).Create(&pv).Error
	if err != nil {
		return &ProfileVelocity{}, err
	}

	//Select created product group
	_, err = pv.FindProfileVelocityByProfileID(db, pv.ResProfileID)
	if err != nil {
		return &ProfileVelocity{}, err
	}

	return pv, nil
}

//UpdateProfileMaster => Update Restriction Profile by ID
func (pv *ProfileVelocity) UpdateProfileVelocity(db *gorm.DB) (*ProfileVelocity, error) {
	var err error
	dateTimeNow := time.Now()

	//Update the data
	err = db.Debug().Model(&pv).Updates(
		map[string]interface{}{
			"freq_limit_per_day":         pv.FreqLimitPerDay,
			"freq_limit_per_week":        pv.FreqLimitPerWeek,
			"freq_limit_per_month":       pv.FreqLimitPerMonth,
			"volume_limit_per_day":       pv.VolumeLimitPerDay,
			"volume_limit_per_week":      pv.VolumeLimitPerWeek,
			"volume_limit_per_month":     pv.VolumeLimitPerMonth,
			"value_limit_per_day":        pv.ValueLimitPerDay,
			"value_limit_per_week":       pv.ValueLimitPerWeek,
			"value_limit_per_month":      pv.ValueLimitPerMonth,
			"prepaid_balance_percentage": pv.PrepaidBalancePercentage,
			"updated_at":                 dateTimeNow,
		}).Error

	if err != nil {
		return &ProfileVelocity{}, err
	}

	//Select updated profile master
	_, err = pv.FindProfileVelocityByProfileID(db, pv.ResProfileID)
	if err != nil {
		return &ProfileVelocity{}, err
	}

	return pv, nil
}

//TableName => Define Table
func (ProfileVelocity) TableName() string {
	return "mstVelocityRestriction"
}
