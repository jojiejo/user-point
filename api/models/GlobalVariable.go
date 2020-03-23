package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type GlobalVariable struct {
	ID        int        `gorm:"primary_key;auto_increment" json:"id"`
	Name      string     `gorm:"size:50" json:"name"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type GlobalVariableDetail struct {
	ID               int        `gorm:"primary_key;auto_increment" json:"id"`
	Name             string     `gorm:"size:50" json:"name"`
	Value            string     `gorm:"size:50" json:"value"`
	GlobalVariableID int        `json:"global_variable_id"`
	CreatedAt        time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (gvd *GlobalVariableDetail) Prepare() {
	gvd.Name = html.EscapeString(strings.TrimSpace(gvd.Name))
	gvd.Value = html.EscapeString(strings.TrimSpace(gvd.Value))
	gvd.CreatedAt = time.Now()
}

func (gvd *GlobalVariableDetail) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if gvd.Name == "" {
		err = errors.New("Name field is required")
		errorMessages["name"] = err.Error()
	}

	if gvd.Value == "" {
		err = errors.New("Value field is required")
		errorMessages["value"] = err.Error()
	}

	return errorMessages
}

func (gv *GlobalVariable) FindAllGlobalVariables(db *gorm.DB) (*[]GlobalVariable, error) {
	var err error
	globalVariables := []GlobalVariable{}
	err = db.Debug().Model(&GlobalVariable{}).Unscoped().Order("created_at desc").Find(&globalVariables).Error
	if err != nil {
		return &[]GlobalVariable{}, err
	}

	return &globalVariables, nil
}

func (gvd *GlobalVariableDetail) FindGlobalVariableDetailByGlobalVariableID(db *gorm.DB, gvID uint64) (*[]GlobalVariableDetail, error) {
	var err error
	var gvds = []GlobalVariableDetail{}
	err = db.Debug().Model(&Terminal{}).Where("global_variable_id = ?", gvID).Order("created_at desc").Find(&gvds).Error
	if err != nil {
		return &[]GlobalVariableDetail{}, err
	}

	return &gvds, nil
}

func (gvd *GlobalVariableDetail) FindGlobalVariableDetailByID(db *gorm.DB, gvdID uint64) (*GlobalVariableDetail, error) {
	var err error
	err = db.Debug().Model(&Terminal{}).Where("id = ?", gvdID).Order("created_at desc").Take(&gvd).Error
	if err != nil {
		return &GlobalVariableDetail{}, err
	}

	return gvd, nil
}

func (gvd *GlobalVariableDetail) CreateGlobalVariableDetail(db *gorm.DB) (*GlobalVariableDetail, error) {
	var err error
	err = db.Debug().Model(&GlobalVariableDetail{}).Create(&gvd).Error
	if err != nil {
		return &GlobalVariableDetail{}, err
	}

	return gvd, nil
}

func (gvd *GlobalVariableDetail) UpdateGlobalVariableDetail(db *gorm.DB) (*GlobalVariableDetail, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&GlobalVariableDetail{}).Unscoped().Where("id = ?", gvd.ID).Updates(
		GlobalVariableDetail{
			Value:     gvd.Value,
			UpdatedAt: &dateTimeNow,
		}).Error

	if err != nil {
		return &GlobalVariableDetail{}, err
	}

	return gvd, nil
}

func (GlobalVariable) TableName() string {
	return "global_variable"
}

func (GlobalVariableDetail) TableName() string {
	return "global_variable_detail"
}
