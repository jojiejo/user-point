package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Site struct {
	ID            int        `gorm:"primary_key;auto_increment" json:"id"`
	OriginalID    int        `json:"original_id"`
	ShipToNumber  string     `gorm:"not null;size:50" json:"ship_to_number"`
	ShipToName    string     `gorm:"not null;size:50" json:"ship_to_name"`
	Address_1     string     `gorm:"not null;size:30" json:"address_1"`
	Address_2     string     `gorm:"not null;size:30" json:"address_2"`
	Address_3     string     `gorm:"not null;size:30" json:"address_3"`
	CityID        int        `gorm:"not null" json:"city_id"`
	City          City       `json:"city"`
	ZipCode       string     `gorm:"not null;size:5" json:"zip_code"`
	Phone         string     `gorm:"not null;size:15" json:"phone"`
	CreatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
	ReactivatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"reactivated_at"`
}

func (site *Site) Prepare() {
	site.ShipToNumber = html.EscapeString(strings.TrimSpace(site.ShipToNumber))
	site.ShipToName = html.EscapeString(strings.TrimSpace(site.ShipToName))
	site.Address_1 = html.EscapeString(strings.TrimSpace(site.Address_1))
	site.Address_2 = html.EscapeString(strings.TrimSpace(site.Address_2))
	site.Address_3 = html.EscapeString(strings.TrimSpace(site.Address_3))
	site.ZipCode = html.EscapeString(strings.TrimSpace(site.ZipCode))
	site.Phone = html.EscapeString(strings.TrimSpace(site.Phone))
	site.CreatedAt = time.Now()
}

func (site *Site) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if site.ShipToNumber == "" {
		err = errors.New("Ship to Number field is required")
		errorMessages["required_ship_to_number"] = err.Error()
	}

	if site.ShipToName == "" {
		err = errors.New("Ship to Name field is required")
		errorMessages["required_ship_to_name"] = err.Error()
	}

	if site.Address_1 == "" {
		err = errors.New("Address 1 field is required")
		errorMessages["required_address_1"] = err.Error()
	}

	if site.Address_2 == "" {
		err = errors.New("Address 2 field is required")
		errorMessages["required_address_2"] = err.Error()
	}

	if site.Address_3 == "" {
		err = errors.New("Address 3 field is required")
		errorMessages["required_address_3"] = err.Error()
	}

	if site.CityID < 1 {
		err = errors.New("City field is required")
		errorMessages["required_city"] = err.Error()
	}

	if site.ZipCode == "" {
		err = errors.New("ZIP Code field is required")
		errorMessages["required_zip_code"] = err.Error()
	}

	if site.Phone == "" {
		err = errors.New("Phone field is required")
		errorMessages["required_phone"] = err.Error()
	}

	return errorMessages
}

func (site *Site) FindAllSites(db *gorm.DB) (*[]Site, error) {
	var err error
	sites := []Site{}
	err = db.Debug().Model(&Site{}).Unscoped().Order("created_at desc").Find(&sites).Error
	if err != nil {
		return &[]Site{}, err
	}

	if len(sites) > 0 {
		for i, _ := range sites {
			err := db.Debug().Model(&City{}).Unscoped().Where("id = ?", sites[i].CityID).Order("id desc").Take(&sites[i].City).Error
			if err != nil {
				return &[]Site{}, err
			}
		}
	}

	return &sites, nil
}

func (site *Site) FindAllLatestSites(db *gorm.DB) (*[]Site, error) {
	var err error
	sites := []Site{}
	err = db.Debug().Raw("EXEC spAPI_Site_GetLatest").Scan(&sites).Error
	if err != nil {
		return &[]Site{}, err
	}

	if len(sites) > 0 {
		for i, _ := range sites {
			err := db.Debug().Model(&City{}).Unscoped().Where("id = ?", sites[i].CityID).Order("id desc").Take(&sites[i].City).Error
			if err != nil {
				return &[]Site{}, err
			}
		}
	}

	return &sites, nil
}

func (site *Site) FindSiteByID(db *gorm.DB, siteID uint64) (*Site, error) {
	var err error
	err = db.Debug().Model(&Site{}).Unscoped().Where("id = ?", siteID).Order("created_at desc").Take(&site).Error
	if err != nil {
		return &Site{}, err
	}

	if site.ID != 0 {
		err := db.Debug().Model(&City{}).Unscoped().Where("id = ?", site.CityID).Order("id desc").Take(&site.City).Error
		if err != nil {
			return &Site{}, err
		}
	}

	return site, nil
}

func (site *Site) FindSiteHistoryByID(db *gorm.DB, originalRetailerID uint64) (*[]Site, error) {
	var err error
	var sites = []Site{}
	err = db.Debug().Model(&Site{}).Unscoped().Where("original_id = ?", originalRetailerID).Order("created_at desc").Find(&sites).Error
	if err != nil {
		return &[]Site{}, err
	}

	if len(sites) > 0 {
		for i, _ := range sites {
			err := db.Debug().Model(&City{}).Unscoped().Where("id = ?", sites[i].CityID).Order("id desc").Take(&sites[i].City).Error
			if err != nil {
				return &[]Site{}, err
			}
		}
	}

	return &sites, nil
}

func (site *Site) CreateSite(db *gorm.DB) (*Site, error) {
	var err error
	tx := db.Begin()
	err = db.Debug().Model(&Site{}).Create(&site).Error
	if err != nil {
		tx.Rollback()
		return &Site{}, err
	}

	err = db.Debug().Model(&Site{}).Where("id = ?", site.ID).Updates(
		Retailer{
			OriginalID: site.ID,
		}).Error
	if err != nil {
		tx.Rollback()
		return &Site{}, err
	}

	tx.Commit()
	return site, nil
}

func (site *Site) UpdateSite(db *gorm.DB) (*Site, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&Site{}).Where("id = ?", site.ID).Updates(
		Site{
			ShipToNumber: site.ShipToNumber,
			ShipToName:   site.ShipToName,
			Address_1:    site.Address_1,
			Address_2:    site.Address_2,
			Address_3:    site.Address_3,
			CityID:       site.CityID,
			ZipCode:      site.ZipCode,
			Phone:        site.Phone,
			UpdatedAt:    &dateTimeNow,
		}).Error

	if err != nil {
		return &Site{}, err
	}

	return site, nil
}

func (site *Site) DeactivateSite(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Site{}).Where("id = ?", site.ID).Delete(&Site{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (site *Site) ReactivateSite(db *gorm.DB) (*Site, error) {
	var err error
	tx := db.Begin()
	dateTimeNow := time.Now()
	err = db.Debug().Model(&Site{}).Unscoped().Where("id = ?", site.ID).Updates(
		Retailer{
			ReactivatedAt: &dateTimeNow,
		}).Error
	if err != nil {
		tx.Rollback()
		return &Site{}, err
	}

	site.ID = 0
	site.DeletedAt = nil
	site.ReactivatedAt = nil
	err = db.Debug().Model(&Site{}).Create(&site).Error
	if err != nil {
		tx.Rollback()
		return &Site{}, err
	}

	tx.Commit()
	return site, nil
}

/*func (site *Site) TerminateSiteLater(db *gorm.DB) (int64, error) {
	var err error
	err = db.Debug().Model(&Retailer{}).Where("id = ?", site.ID).Updates(
		Retailer{
			DeletedAt: site.DeletedAt,
		}).Error

	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (site *Site) TerminateSiteNow(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Site{}).Where("id = ?", site.ID).Delete(&Site{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}*/

func (Site) TableName() string {
	return "site"
}
