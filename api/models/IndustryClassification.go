package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type IndustryClassification struct {
	ID        uint64     `gorm:"primary_key;auto_increment;" json:"id"`
	Code      string     `gorm:"not null;" json:"code"`
	Name      string     `gorm:"not null;size:100;" json:"name"`
	StartedAt *time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (ic *IndustryClassification) Prepare() {
	ic.Code = html.EscapeString(strings.TrimSpace(ic.Code))
	ic.Name = html.EscapeString(strings.TrimSpace(ic.Name))
	ic.CreatedAt = time.Now()
}

func (ic *IndustryClassification) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if ic.Code == "" {
		err = errors.New("Code field is required")
		errorMessages["code"] = err.Error()
	}

	if ic.Name == "" {
		err = errors.New("Name field is required")
		errorMessages["name"] = err.Error()
	}

	if ic.StartedAt == nil {
		err = errors.New("Started at field is required")
		errorMessages["started_at"] = err.Error()
	}

	return errorMessages
}

func (ic *IndustryClassification) FindIndustryClassifications(db *gorm.DB) (*[]IndustryClassification, error) {
	var err error
	ics := []IndustryClassification{}
	err = db.Debug().Model(&IndustryClassification{}).
		Order("id, created_at desc").
		Find(&ics).Error

	if err != nil {
		return &[]IndustryClassification{}, err
	}

	return &ics, nil
}

func (ic *IndustryClassification) FindIndustryClassification(db *gorm.DB, industryClassificationID uint64) (*IndustryClassification, error) {
	var err error
	err = db.Debug().Model(&IndustryClassification{}).Unscoped().
		Where("id = ?", industryClassificationID).
		Order("created_at desc").
		Take(&ic).Error

	if err != nil {
		return &IndustryClassification{}, err
	}

	return ic, nil
}

func (ic *IndustryClassification) CreateIndustryClassification(db *gorm.DB) (*IndustryClassification, error) {
	var err error
	err = db.Debug().Model(&IndustryClassification{}).Create(&ic).Error
	if err != nil {
		return &IndustryClassification{}, err
	}

	//Select created fee
	_, err = ic.FindIndustryClassification(db, ic.ID)
	if err != nil {
		return &IndustryClassification{}, err
	}

	return ic, nil
}

func (ic *IndustryClassification) UpdateIndustryClassification(db *gorm.DB) (*IndustryClassification, error) {
	var err error
	dateTimeNow := time.Now()

	//Update the data
	err = db.Debug().Model(&ic).Updates(
		map[string]interface{}{
			"code":       ic.Code,
			"name":       ic.Name,
			"started_at": ic.StartedAt,
			"ended_at":   ic.EndedAt,
			"updated_at": dateTimeNow,
		}).Error

	if err != nil {
		return &IndustryClassification{}, err
	}

	//Select updated product
	_, err = ic.FindIndustryClassification(db, ic.ID)
	if err != nil {
		return &IndustryClassification{}, err
	}

	return ic, nil
}

func (IndustryClassification) TableName() string {
	return "gsap_industry_class"
}
