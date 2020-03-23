package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type SiteTerminalRelation struct {
	ID         int                                   `gorm:"primary_key;auto_increment" json:"id"`
	SiteID     int                                   `gorm:"not null" json:"site_id"`
	Site       SiteTerminalRelationDisplayedSite     `json:"site"`
	TerminalID int                                   `gorm:"not null" json:"terminal_id"`
	Terminal   SiteTerminalRelationDisplayedTerminal `json:"terminal"`
	StartedAt  time.Time                             `json:"started_at"`
	EndedAt    *time.Time                            `json:"ended_at"`
	CreatedAt  time.Time                             `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  *time.Time                            `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt  *time.Time                            `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type SiteTerminalRelationDisplayedSite struct {
	ID           int    `gorm:"primary_key;auto_increment" json:"id"`
	ShipToNumber string `gorm:"not null;size:50" json:"ship_to_number"`
	ShipToName   string `gorm:"not null;size:50" json:"ship_to_name"`
	Address_1    string `gorm:"not null;size:30" json:"address_1"`
	Address_2    string `gorm:"not null;size:30" json:"address_2"`
	Address_3    string `gorm:"not null;size:30" json:"address_3"`
	CityID       int    `gorm:"not null" json:"city_id"`
	City         City   `json:"city"`
}

type SiteTerminalRelationDisplayedTerminal struct {
	ID         int    `gorm:"primary_key;auto_increment" json:"id"`
	TerminalSN string `gorm:"not null;size:10" json:"terminal_sn"`
	TerminalID string `gorm:"not null; size:60" json:"terminal_id"`
	MerchantID string `gorm:"not null;size:30" json:"merchant_id"`
}

func (siteTerminalRelation *SiteTerminalRelation) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if siteTerminalRelation.TerminalID < 1 {
		err = errors.New("Retailer field is required")
		errorMessages["required_retailer"] = err.Error()
	}

	if siteTerminalRelation.SiteID < 1 {
		err = errors.New("Site field is required")
		errorMessages["required_site"] = err.Error()
	}

	return errorMessages
}

func (siteTerminalRelation *SiteTerminalRelation) FindAllSiteTerminalRelationBySiteID(db *gorm.DB, siteTerminalRelationID uint64) (*[]SiteTerminalRelation, error) {
	var err error
	siteTerminalRelations := []SiteTerminalRelation{}
	err = db.Debug().Model(&SiteTerminalRelation{}).Unscoped().Where("site_id = ?", siteTerminalRelationID).Order("created_at desc").Find(&siteTerminalRelations).Error
	if err != nil {
		return &[]SiteTerminalRelation{}, err
	}

	if len(siteTerminalRelations) > 0 {
		for i, _ := range siteTerminalRelations {
			siteErr := db.Debug().Model(&SiteTerminalRelationDisplayedSite{}).Unscoped().Where("original_id = ?", siteTerminalRelations[i].SiteID).Order("id desc").Take(&siteTerminalRelations[i].Site).Error
			if siteErr != nil {
				return &[]SiteTerminalRelation{}, err
			}

			if siteTerminalRelations[i].Site.CityID != 0 {
				err := db.Debug().Model(&City{}).Unscoped().Where("id = ?", siteTerminalRelations[i].Site.CityID).Order("id desc").Take(&siteTerminalRelations[i].Site.City).Error
				if err != nil {
					return &[]SiteTerminalRelation{}, err
				}
			}

			terminalErr := db.Debug().Model(&SiteTerminalRelationDisplayedTerminal{}).Unscoped().Where("original_id = ?", siteTerminalRelations[i].TerminalID).Order("id desc").Take(&siteTerminalRelations[i].Terminal).Error
			if terminalErr != nil {
				return &[]SiteTerminalRelation{}, err
			}
		}
	}

	return &siteTerminalRelations, nil
}

func (siteTerminalRelation *SiteTerminalRelation) CreateSiteTerminalRelation(db *gorm.DB) (*SiteTerminalRelation, error) {
	var err error
	tx := db.Begin()
	dateTimeNow := time.Now()

	err = db.Debug().Model(&SiteTerminalRelation{}).Where("site_id = ? AND terminal_id = ? AND ended_at IS NULL", siteTerminalRelation.SiteID, siteTerminalRelation.TerminalID).Updates(
		SiteTerminalRelation{
			EndedAt:   &siteTerminalRelation.StartedAt,
			UpdatedAt: &dateTimeNow,
		}).Error
	if err != nil {
		tx.Rollback()
		return &SiteTerminalRelation{}, err
	}

	err = db.Debug().Model(&SiteTerminalRelation{}).Create(&siteTerminalRelation).Error
	if err != nil {
		tx.Rollback()
		return &SiteTerminalRelation{}, err
	}

	tx.Commit()
	return siteTerminalRelation, nil
}

func (siteTerminalRelation *SiteTerminalRelation) UpdateSiteTerminalRelation(db *gorm.DB) (*SiteTerminalRelation, error) {
	var err error
	dateTimeNow := time.Now()
	err = db.Debug().Model(&SiteTerminalRelation{}).Where("id = ?", siteTerminalRelation.ID).Updates(
		SiteTerminalRelation{
			EndedAt:   siteTerminalRelation.EndedAt,
			UpdatedAt: &dateTimeNow,
		}).Error

	if err != nil {
		return &SiteTerminalRelation{}, err
	}

	return siteTerminalRelation, nil
}

func (siteTerminalRelation *SiteTerminalRelation) UnlinkSiteTerminalRelation(db *gorm.DB) (*SiteTerminalRelation, error) {
	var err error
	dateTimeNow := time.Now()
	err = db.Debug().Model(&SiteTerminalRelation{}).Where("id = ?", siteTerminalRelation.ID).Updates(
		SiteTerminalRelation{
			EndedAt:   &dateTimeNow,
			UpdatedAt: &dateTimeNow,
		}).Error

	if err != nil {
		return &SiteTerminalRelation{}, err
	}

	return siteTerminalRelation, nil
}

func (SiteTerminalRelation) TableName() string {
	return "site_terminal_relation"
}

func (SiteTerminalRelationDisplayedSite) TableName() string {
	return "site"
}

func (SiteTerminalRelationDisplayedTerminal) TableName() string {
	return "terminal"
}
