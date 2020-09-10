package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type BusinessType struct {
	ID                         uint64     `gorm:"primary_key;auto_increment;" json:"id"`
	Code                       string     `gorm:"not null;" json:"code"`
	Name                       string     `gorm:"not null;size:100;" json:"name"`
	LineOffBusinessDescription string     `gorm:"not null;" json:"line_off_business_description"`
	StartedAt                  *time.Time `json:"started_at"`
	EndedAt                    *time.Time `json:"ended_at"`
	CreatedAt                  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt                  *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (bt *BusinessType) Prepare() {
	bt.Code = html.EscapeString(strings.TrimSpace(bt.Code))
	bt.Name = html.EscapeString(strings.TrimSpace(bt.Name))
	bt.CreatedAt = time.Now()
}

func (bt *BusinessType) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if bt.Code == "" {
		err = errors.New("Code field is required")
		errorMessages["code"] = err.Error()
	}

	if bt.Name == "" {
		err = errors.New("Name field is required")
		errorMessages["name"] = err.Error()
	}

	if bt.LineOffBusinessDescription == "" {
		err = errors.New("Line off business description field is required")
		errorMessages["name"] = err.Error()
	}

	return errorMessages
}

func (bt *BusinessType) FindBusinessTypes(db *gorm.DB) (*[]BusinessType, error) {
	var err error
	bts := []BusinessType{}
	err = db.Debug().Model(&BusinessType{}).
		Order("id, created_at desc").
		Find(&bts).Error

	if err != nil {
		return &[]BusinessType{}, err
	}

	return &bts, nil
}

func (bt *BusinessType) FindBusinessTypeByID(db *gorm.DB, btID uint64) (*BusinessType, error) {
	var err error
	err = db.Debug().Model(&BusinessType{}).Unscoped().
		Where("id = ?", btID).
		Order("created_at desc").
		Take(&bt).Error

	if err != nil {
		return &BusinessType{}, err
	}

	return bt, nil
}

func (bt *BusinessType) CreateBusinessType(db *gorm.DB) (*BusinessType, error) {
	var err error
	err = db.Debug().Model(&BusinessType{}).Create(&bt).Error
	if err != nil {
		return &BusinessType{}, err
	}

	//Select created fee
	_, err = bt.FindBusinessTypeByID(db, bt.ID)
	if err != nil {
		return &BusinessType{}, err
	}

	return bt, nil
}

func (bt *BusinessType) UpdateBusinessType(db *gorm.DB) (*BusinessType, error) {
	var err error
	dateTimeNow := time.Now()

	//Update the data
	err = db.Debug().Model(&bt).Updates(
		map[string]interface{}{
			"code":                          bt.Code,
			"name":                          bt.Name,
			"updated_at":                    dateTimeNow,
			"started_at":                    bt.StartedAt,
			"ended_at":                      bt.EndedAt,
			"line_off_business_description": bt.LineOffBusinessDescription,
		}).Error

	if err != nil {
		return &BusinessType{}, err
	}

	//Select updated sales rep
	_, err = bt.FindBusinessTypeByID(db, bt.ID)
	if err != nil {
		return &BusinessType{}, err
	}

	return bt, nil
}

func (BusinessType) TableName() string {
	return "gsap_business_type"
}
