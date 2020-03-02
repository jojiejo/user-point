package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type RetailerSiteRelation struct {
	ID                int               `gorm:"primary_key;auto_increment" json:"id"`
	RetailerID        int               `gorm:"not null" json:"retailer_id"`
	Retailer          DisplayedRetailer `json:"retailer"`
	SiteID            int               `gorm:"not null" json:"site_id"`
	Site              DisplayedSite     `json:"site"`
	BankAccountNumber string            `json:"bank_account_number"`
	StartedAt         time.Time         `json:"start_date"`
	EndedAt           time.Time         `json:"end_date"`
	CreatedAt         time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt         time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type DisplayedSite struct {
	ID           int    `gorm:"primary_key;auto_increment" json:"id"`
	ShipToNumber string `gorm:"not null;size:50" json:"ship_to_number"`
	ShipToName   string `gorm:"not null;size:50" json:"ship_to_name"`
}

type DisplayedRetailer struct {
	ID           int    `gorm:"primary_key;auto_increment" json:"id"`
	SoldToNumber string `gorm:"not null;size:10" json:"sold_to_number"`
	SoldToName   string `gorm:"not null; size:60" json:"sold_to_name"`
}

func (retailerSiteRelation *RetailerSiteRelation) Prepare() {
	retailerSiteRelation.BankAccountNumber = html.EscapeString(strings.TrimSpace(retailerSiteRelation.BankAccountNumber))
}

func (retailerSiteRelation *RetailerSiteRelation) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if retailerSiteRelation.BankAccountNumber == "" {
		err = errors.New("Bank account number field is required")
		errorMessages["required_bank_account_number"] = err.Error()
	}

	return errorMessages
}

func (retailerSiteRelation *RetailerSiteRelation) FindAllRetailerSiteRelationByRetailerID(db *gorm.DB, retailerSiteRelationID uint64) (*RetailerSiteRelation, error) {
	var err error
	err = db.Debug().Model(&RetailerSiteRelation{}).Where("retailer_id = ?", retailerSiteRelationID).Order("created_at desc").Take(&retailerSiteRelation).Error
	if err != nil {
		return &RetailerSiteRelation{}, err
	}

	if retailerSiteRelation.ID != 0 {
		siteErr := db.Debug().Model(&Site{}).Select("id, ship_to_number, ship_to_name").Where("id = ?", retailerSiteRelation.SiteID).Order("id desc").Take(&retailerSiteRelation.Site).Error
		if siteErr != nil {
			return &RetailerSiteRelation{}, err
		}

		retailerErr := db.Debug().Model(&Retailer{}).Select("id, sold_to_number, sold_to_name").Where("id = ?", retailerSiteRelation.RetailerID).Order("id desc").Take(&retailerSiteRelation.Retailer).Error
		if retailerErr != nil {
			return &RetailerSiteRelation{}, err
		}
	}

	return retailerSiteRelation, nil
}

func (retailerSiteRelation *RetailerSiteRelation) CreateRetailerSiteRelation(db *gorm.DB) (*RetailerSiteRelation, error) {
	var err error
	err = db.Debug().Model(&RetailerSiteRelation{}).Create(&retailerSiteRelation).Error
	if err != nil {
		return &RetailerSiteRelation{}, err
	}

	return retailerSiteRelation, nil
}

func (retailerSiteRelation *RetailerSiteRelation) UpdateRetailerSiteRelation(db *gorm.DB) (*RetailerSiteRelation, error) {
	var err error
	err = db.Debug().Model(&Site{}).Where("id = ?", retailerSiteRelation.ID).Updates(
		RetailerSiteRelation{
			BankAccountNumber: retailerSiteRelation.BankAccountNumber,
			StartedAt:         retailerSiteRelation.StartedAt,
			EndedAt:           retailerSiteRelation.EndedAt,
			UpdatedAt:         time.Now(),
		}).Error

	if err != nil {
		return &RetailerSiteRelation{}, err
	}

	return retailerSiteRelation, nil
}

func (retailerSiteRelation *RetailerSiteRelation) TerminateRetailerSiteRelation(db *gorm.DB) (*RetailerSiteRelation, error) {
	var err error
	err = db.Debug().Model(&Site{}).Where("id = ?", retailerSiteRelation.ID).Updates(
		RetailerSiteRelation{
			EndedAt:   time.Now(),
			UpdatedAt: time.Now(),
		}).Error

	if err != nil {
		return &RetailerSiteRelation{}, err
	}

	return retailerSiteRelation, nil
}

func (RetailerSiteRelation) TableName() string {
	return "retailer_site_relation"
}

func (DisplayedSite) TableName() string {
	return "site"
}

func (DisplayedRetailer) TableName() string {
	return "retailer"
}
