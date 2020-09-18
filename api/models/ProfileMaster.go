package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

//ProfileMaster => Restriction Profile
type ProfileMaster struct {
	ResProfileID    uint64          `gorm:"primary_key;auto_increment" json:"res_profile_id"`
	CardProfileFlag uint32          `json:"card_profile_flag"`
	AllSiteFlag     uint32          `json:"all_site_flag"`
	Name            string          `gorm:"column:res_profile_name" json:"name"`
	Notes           string          `json:"notes"`
	Velocity        ProfileVelocity `gorm:"foreignkey:ResProfileID;references:ResProfileID" json:"velocity"`
	CreatedAt       time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       *time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

//Prepare => Prepare Restriction Profile
func (pm *ProfileMaster) Prepare() {
	pm.Name = html.EscapeString(strings.TrimSpace(pm.Name))
	pm.CreatedAt = time.Now()
}

//Validate => Validate Restriction Profile
func (pm *ProfileMaster) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if pm.Name == "" {
		err = errors.New("Name field is required")
		errorMessages["name"] = err.Error()
	}

	if pm.CardProfileFlag != 0 && pm.CardProfileFlag != 1 {
		err = errors.New("Invalid card profile flag")
		errorMessages["card_profile_flag"] = err.Error()
	}

	if pm.AllSiteFlag != 0 && pm.AllSiteFlag != 1 {
		err = errors.New("Invalid all site flag")
		errorMessages["all_site_flag"] = err.Error()
	}

	return errorMessages
}

//FindProfileMasters => Find Restriction Profile
func (pm *ProfileMaster) FindProfileMasters(db *gorm.DB) (*[]ProfileMaster, error) {
	var err error
	pms := []ProfileMaster{}
	err = db.Debug().Model(&ProfileMaster{}).
		Order("res_profile_id, created_at desc").
		Find(&pms).Error

	if err != nil {
		return &[]ProfileMaster{}, err
	}

	return &pms, nil
}

//FindProfileMasterByID => Find Restriction Profile by ID
func (pm *ProfileMaster) FindProfileMasterByID(db *gorm.DB, resProfileID uint64) (*ProfileMaster, error) {
	var err error
	err = db.Debug().Model(&ProfileMaster{}).Unscoped().
		Where("res_profile_id = ?", resProfileID).
		Order("created_at desc").
		Take(&pm).Error

	if err != nil {
		return &ProfileMaster{}, err
	}

	return pm, nil
}

//CreateProfileMaster => Create Restriction Profile by ID
func (pm *ProfileMaster) CreateProfileMaster(db *gorm.DB) (*ProfileMaster, error) {
	var err error
	err = db.Debug().Model(&ProfileMaster{}).Create(&pm).Error
	if err != nil {
		return &ProfileMaster{}, err
	}

	//Select created product group
	_, err = pm.FindProfileMasterByID(db, pm.ResProfileID)
	if err != nil {
		return &ProfileMaster{}, err
	}

	return pm, nil
}

//UpdateProfileMaster => Update Restriction Profile by ID
func (pm *ProfileMaster) UpdateProfileMaster(db *gorm.DB) (*ProfileMaster, error) {
	var err error
	dateTimeNow := time.Now()

	//Update the data
	err = db.Debug().Model(&pm).Updates(
		map[string]interface{}{
			"res_profile_name":  pm.Name,
			"notes":             pm.Notes,
			"card_profile_flag": pm.CardProfileFlag,
			"all_site_flag":     pm.AllSiteFlag,
			"updated_at":        dateTimeNow,
		}).Error

	if err != nil {
		return &ProfileMaster{}, err
	}

	//Select updated profile master
	_, err = pm.FindProfileMasterByID(db, pm.ResProfileID)
	if err != nil {
		return &ProfileMaster{}, err
	}

	return pm, nil
}

//TableName => Define Table
func (ProfileMaster) TableName() string {
	return "mstRestrictionProfile"
}
