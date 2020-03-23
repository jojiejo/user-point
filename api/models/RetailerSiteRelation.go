package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type RetailerSiteRelation struct {
	ID                int                                   `gorm:"primary_key;auto_increment" json:"id"`
	RetailerID        int                                   `gorm:"not null" json:"retailer_id"`
	Retailer          RetailerSiteRelationDisplayedRetailer `json:"retailer"`
	SiteID            int                                   `gorm:"not null" json:"site_id"`
	Site              RetailerSiteRelationDisplayedSite     `json:"site"`
	BankAccountNumber string                                `json:"bank_account_number"`
	StartedAt         time.Time                             `json:"started_at"`
	EndedAt           *time.Time                            `json:"ended_at"`
	CreatedAt         time.Time                             `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         *time.Time                            `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt         *time.Time                            `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type RetailerSiteRelationDisplayedSite struct {
	ID           int    `gorm:"primary_key;auto_increment" json:"id"`
	ShipToNumber string `gorm:"not null;size:50" json:"ship_to_number"`
	ShipToName   string `gorm:"not null;size:50" json:"ship_to_name"`
	Address_1    string `gorm:"not null;size:30" json:"address_1"`
	Address_2    string `gorm:"not null;size:30" json:"address_2"`
	Address_3    string `gorm:"not null;size:30" json:"address_3"`
	CityID       int    `gorm:"not null" json:"city_id"`
	City         City   `json:"city"`
}

type RetailerSiteRelationDisplayedRetailer struct {
	ID           int    `gorm:"primary_key;auto_increment" json:"id"`
	SoldToNumber string `gorm:"not null;size:10" json:"sold_to_number"`
	SoldToName   string `gorm:"not null; size:60" json:"sold_to_name"`
	Address_1    string `gorm:"not null;size:30" json:"address_1"`
	Address_2    string `gorm:"not null;size:30" json:"address_2"`
	Address_3    string `gorm:"not null;size:30" json:"address_3"`
	CityID       int    `gorm:"not null" json:"city_id"`
	City         City   `json:"city"`
}

func (retailerSiteRelation *RetailerSiteRelation) Prepare() {
	retailerSiteRelation.BankAccountNumber = html.EscapeString(strings.TrimSpace(retailerSiteRelation.BankAccountNumber))
}

func (retailerSiteRelation *RetailerSiteRelation) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if len(retailerSiteRelation.BankAccountNumber) > 0 && len(retailerSiteRelation.BankAccountNumber) != 4 {
		err = errors.New("Bank account number field must contain 4 characters")
		errorMessages["bank_account_number"] = err.Error()
	}

	if retailerSiteRelation.RetailerID < 1 {
		err = errors.New("Retailer field is required")
		errorMessages["retailer"] = err.Error()
	}

	if retailerSiteRelation.SiteID < 1 {
		err = errors.New("Site field is required")
		errorMessages["site"] = err.Error()
	}

	return errorMessages
}

func (retailerSiteRelation *RetailerSiteRelation) UpdateValidate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if len(retailerSiteRelation.BankAccountNumber) > 0 && len(retailerSiteRelation.BankAccountNumber) != 4 {
		err = errors.New("Bank account number field must contain 4 characters")
		errorMessages["bank_account_number"] = err.Error()
	}

	return errorMessages
}

func (retailerSiteRelation *RetailerSiteRelation) FindAllRetailerSiteRelationByRetailerID(db *gorm.DB, retailerSiteRelationID uint64) (*[]RetailerSiteRelation, error) {
	var err error
	retailerSiteRelations := []RetailerSiteRelation{}
	err = db.Debug().Model(&RetailerSiteRelation{}).Unscoped().Where("retailer_id = ?", retailerSiteRelationID).Order("created_at desc").Find(&retailerSiteRelations).Error
	if err != nil {
		return &[]RetailerSiteRelation{}, err
	}

	if len(retailerSiteRelations) > 0 {
		for i, _ := range retailerSiteRelations {
			siteErr := db.Debug().Model(&RetailerSiteRelationDisplayedSite{}).Unscoped().Where("original_id = ?", retailerSiteRelations[i].SiteID).Order("id desc").Take(&retailerSiteRelations[i].Site).Error
			if siteErr != nil {
				return &[]RetailerSiteRelation{}, err
			}

			if retailerSiteRelations[i].Site.CityID != 0 {
				err := db.Debug().Model(&City{}).Unscoped().Where("id = ?", retailerSiteRelations[i].Site.CityID).Order("id desc").Take(&retailerSiteRelations[i].Site.City).Error
				if err != nil {
					return &[]RetailerSiteRelation{}, err
				}
			}

			retailerErr := db.Debug().Model(&RetailerSiteRelationDisplayedRetailer{}).Unscoped().Where("original_id = ?", retailerSiteRelations[i].RetailerID).Order("id desc").Take(&retailerSiteRelations[i].Retailer).Error
			if retailerErr != nil {
				return &[]RetailerSiteRelation{}, err
			}

			if retailerSiteRelations[i].Retailer.CityID != 0 {
				err := db.Debug().Model(&City{}).Unscoped().Where("id = ?", retailerSiteRelations[i].Retailer.CityID).Order("id desc").Take(&retailerSiteRelations[i].Retailer.City).Error
				if err != nil {
					return &[]RetailerSiteRelation{}, err
				}
			}
		}
	}

	return &retailerSiteRelations, nil
}

func (retailerSiteRelation *RetailerSiteRelation) CreateRetailerSiteRelation(db *gorm.DB) (*RetailerSiteRelation, error) {
	var err error
	tx := db.Begin()
	dateTimeNow := time.Now()

	err = db.Debug().Model(&RetailerSiteRelation{}).Where("retailer_id = ? AND site_id = ? AND ended_at IS NULL", retailerSiteRelation.RetailerID, retailerSiteRelation.SiteID).Updates(
		RetailerSiteRelation{
			EndedAt:   &retailerSiteRelation.StartedAt,
			UpdatedAt: &dateTimeNow,
		}).Error
	if err != nil {
		tx.Rollback()
		return &RetailerSiteRelation{}, err
	}

	err = db.Debug().Model(&RetailerSiteRelation{}).Create(&retailerSiteRelation).Error
	if err != nil {
		tx.Rollback()
		return &RetailerSiteRelation{}, err
	}

	tx.Commit()
	return retailerSiteRelation, nil
}

func (retailerSiteRelation *RetailerSiteRelation) UpdateRetailerSiteRelation(db *gorm.DB) (*RetailerSiteRelation, error) {
	var err error
	dateTimeNow := time.Now()
	err = db.Debug().Model(&RetailerSiteRelation{}).Where("id = ?", retailerSiteRelation.ID).Updates(
		RetailerSiteRelation{
			StartedAt:         retailerSiteRelation.StartedAt,
			BankAccountNumber: retailerSiteRelation.BankAccountNumber,
			UpdatedAt:         &dateTimeNow,
		}).Error

	if err != nil {
		return &RetailerSiteRelation{}, err
	}

	return retailerSiteRelation, nil
}

func (retailerSiteRelation *RetailerSiteRelation) UnlinkRetailerSiteRelation(db *gorm.DB) (*RetailerSiteRelation, error) {
	var err error
	dateTimeNow := time.Now()
	err = db.Debug().Model(&RetailerSiteRelation{}).Where("id = ?", retailerSiteRelation.ID).Updates(
		RetailerSiteRelation{
			EndedAt:   &dateTimeNow,
			UpdatedAt: &dateTimeNow,
		}).Error

	if err != nil {
		return &RetailerSiteRelation{}, err
	}

	return retailerSiteRelation, nil
}

func (RetailerSiteRelation) TableName() string {
	return "retailer_site_relation"
}

func (RetailerSiteRelationDisplayedSite) TableName() string {
	return "site"
}

func (RetailerSiteRelationDisplayedRetailer) TableName() string {
	return "retailer"
}
