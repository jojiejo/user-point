package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type SalesRep struct {
	ID        uint64     `gorm:"primary_key;auto_increment;" json:"id"`
	Code      string     `gorm:"not null;" json:"code"`
	Name      string     `gorm:"not null;size:100;" json:"name"`
	StartedAt *time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (sr *SalesRep) Prepare() {
	sr.Code = html.EscapeString(strings.TrimSpace(sr.Code))
	sr.Name = html.EscapeString(strings.TrimSpace(sr.Name))
	sr.CreatedAt = time.Now()
}

func (sr *SalesRep) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if sr.Code == "" {
		err = errors.New("Code field is required")
		errorMessages["code"] = err.Error()
	}

	if sr.Name == "" {
		err = errors.New("Name field is required")
		errorMessages["name"] = err.Error()
	}

	if sr.StartedAt == nil {
		err = errors.New("Started at field is required")
		errorMessages["material"] = err.Error()
	}

	return errorMessages
}

func (sr *SalesRep) FindSalesReps(db *gorm.DB) (*[]SalesRep, error) {
	var err error
	srs := []SalesRep{}
	err = db.Debug().Model(&SalesRep{}).
		Order("id, created_at desc").
		Find(&srs).Error

	if err != nil {
		return &[]SalesRep{}, err
	}

	return &srs, nil
}

func (sr *SalesRep) FindSalesRepByID(db *gorm.DB, srID uint64) (*SalesRep, error) {
	var err error
	err = db.Debug().Model(&SalesRep{}).Unscoped().
		Where("id = ?", srID).
		Order("created_at desc").
		Take(&sr).Error

	if err != nil {
		return &SalesRep{}, err
	}

	return sr, nil
}

func (sr *SalesRep) CreateSalesRep(db *gorm.DB) (*SalesRep, error) {
	var err error
	err = db.Debug().Model(&SalesRep{}).Create(&sr).Error
	if err != nil {
		return &SalesRep{}, err
	}

	//Select created fee
	_, err = sr.FindSalesRepByID(db, sr.ID)
	if err != nil {
		return &SalesRep{}, err
	}

	return sr, nil
}

func (sr *SalesRep) UpdateSalesRep(db *gorm.DB) (*SalesRep, error) {
	var err error
	dateTimeNow := time.Now()

	//Update the data
	err = db.Debug().Model(&sr).Updates(
		map[string]interface{}{
			"code":       sr.Code,
			"name":       sr.Name,
			"started_at": sr.StartedAt,
			"ended_at":   sr.EndedAt,
			"updated_at": dateTimeNow,
		}).Error

	if err != nil {
		return &SalesRep{}, err
	}

	//Select updated sales rep
	_, err = sr.FindSalesRepByID(db, sr.ID)
	if err != nil {
		return &SalesRep{}, err
	}

	return sr, nil
}

func (SalesRep) TableName() string {
	return "gsap_sales_rep"
}
