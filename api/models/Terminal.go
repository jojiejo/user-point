package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Terminal struct {
	ID            int        `gorm:"primary_key;auto_increment" json:"id"`
	OriginalID    int        `json:"original_id"`
	TerminalSN    string     `gorm:"not null;" json:"terminal_sn"`
	TerminalID    string     `gorm:"not null;" json:"terminal_id"`
	MerchantID    string     `gorm:"not null;" json:"merchant_id"`
	SiteID        int        `json:"site_id"`
	Site          Site       `json:"site"`
	StartedAt     time.Time  `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at"`
	CreatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
	ReactivatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"reactivated_at"`
}

func (terminal *Terminal) Prepare() {
	terminal.TerminalSN = html.EscapeString(strings.TrimSpace(terminal.TerminalSN))
	terminal.TerminalID = html.EscapeString(strings.TrimSpace(terminal.TerminalID))
	terminal.MerchantID = html.EscapeString(strings.TrimSpace(terminal.MerchantID))
	terminal.CreatedAt = time.Now()
}

func (terminal *Terminal) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if len(terminal.TerminalSN) != 9 {
		err = errors.New("Terminal SN must contain 9 characters")
		errorMessages["terminal_sn"] = err.Error()
	}

	if terminal.TerminalSN == "" {
		err = errors.New("Terminal SN field is required")
		errorMessages["terminal_sn"] = err.Error()
	}

	if len(terminal.TerminalID) != 8 {
		err = errors.New("Terminal ID must contain 8 characters")
		errorMessages["terminal_id"] = err.Error()
	}

	if terminal.TerminalID == "" {
		err = errors.New("Terminal ID field is required")
		errorMessages["terminal_id"] = err.Error()
	}

	if len(terminal.MerchantID) != 15 {
		err = errors.New("Merchant ID must contain 15 characters")
		errorMessages["merchant_id"] = err.Error()
	}

	if terminal.MerchantID == "" {
		err = errors.New("Merchant ID field is required")
		errorMessages["merchant_id"] = err.Error()
	}

	if terminal.SiteID < 1 {
		err = errors.New("Site ID field is required")
		errorMessages["site_id"] = err.Error()
	}

	return errorMessages
}

func (terminal *Terminal) FindAllTerminals(db *gorm.DB) (*[]Terminal, error) {
	var err error
	terminals := []Terminal{}
	err = db.Debug().Model(&Terminal{}).Unscoped().Order("terminal_sn, terminal_id, merchant_id, created_at desc").Find(&terminals).Error
	if err != nil {
		return &[]Terminal{}, err
	}

	if len(terminals) > 0 {
		for i, _ := range terminals {
			siteErr := db.Debug().Model(&City{}).Unscoped().Where("id = ?", terminals[i].SiteID).Order("id desc").Take(&terminals[i].Site).Error
			if siteErr != nil {
				return &[]Terminal{}, err
			}

			cityErr := db.Debug().Model(&City{}).Unscoped().Where("id = ?", terminals[i].Site.CityID).Order("id desc").Take(&terminals[i].Site.City).Error
			if cityErr != nil {
				return &[]Terminal{}, err
			}

			siteTypeErr := db.Debug().Model(&SiteType{}).Unscoped().Where("id = ?", terminals[i].Site.SiteTypeID).Order("id desc").Take(&terminals[i].Site.SiteType).Error
			if siteTypeErr != nil {
				return &[]Terminal{}, err
			}
		}
	}

	return &terminals, nil
}

func (terminal *Terminal) FindAllLatestTerminals(db *gorm.DB) (*[]Terminal, error) {
	var err error
	terminals := []Terminal{}
	err = db.Debug().Raw("EXEC spAPI_Terminal_GetLatest").Scan(&terminals).Error
	if err != nil {
		return &[]Terminal{}, err
	}

	if len(terminals) > 0 {
		for i, _ := range terminals {
			siteErr := db.Debug().Model(&City{}).Unscoped().Where("id = ?", terminals[i].SiteID).Order("id desc").Take(&terminals[i].Site).Error
			if siteErr != nil {
				return &[]Terminal{}, err
			}

			cityErr := db.Debug().Model(&City{}).Unscoped().Where("id = ?", terminals[i].Site.CityID).Order("id desc").Take(&terminals[i].Site.City).Error
			if cityErr != nil {
				return &[]Terminal{}, err
			}

			siteTypeErr := db.Debug().Model(&SiteType{}).Unscoped().Where("id = ?", terminals[i].Site.SiteTypeID).Order("id desc").Take(&terminals[i].Site.SiteType).Error
			if siteTypeErr != nil {
				return &[]Terminal{}, err
			}
		}
	}

	return &terminals, nil
}

func (terminal *Terminal) FindTerminalByID(db *gorm.DB, terminalID uint64) (*Terminal, error) {
	var err error
	err = db.Debug().Model(&Terminal{}).Where("id = ?", terminalID).Order("created_at desc").Take(&terminal).Error
	if err != nil {
		return &Terminal{}, err
	}

	siteErr := db.Debug().Model(&City{}).Unscoped().Where("id = ?", terminal.SiteID).Order("id desc").Take(&terminal.Site).Error
	if siteErr != nil {
		return &Terminal{}, err
	}

	cityErr := db.Debug().Model(&City{}).Unscoped().Where("id = ?", terminal.Site.CityID).Order("id desc").Take(&terminal.Site.City).Error
	if cityErr != nil {
		return &Terminal{}, err
	}

	siteTypeErr := db.Debug().Model(&SiteType{}).Unscoped().Where("id = ?", terminal.Site.SiteTypeID).Order("id desc").Take(&terminal.Site.SiteType).Error
	if siteTypeErr != nil {
		return &Terminal{}, err
	}

	return terminal, nil
}

func (terminal *Terminal) FindAllTerminalBySiteID(db *gorm.DB, siteID uint64) (*[]Terminal, error) {
	var err error

	terminals := []Terminal{}
	err = db.Debug().Model(&Terminal{}).Unscoped().Where("site_id = ?", siteID).Order("created_at desc").Find(&terminals).Error
	if err != nil {
		return &[]Terminal{}, err
	}

	return &terminals, nil
}

func (terminal *Terminal) FindTerminalHistoryByID(db *gorm.DB, originalRetailerID uint64) (*[]Terminal, error) {
	var err error
	var terminals = []Terminal{}
	err = db.Debug().Model(&Terminal{}).Unscoped().Where("original_id = ?", originalRetailerID).Order("created_at desc").Find(&terminals).Error
	if err != nil {
		return &[]Terminal{}, err
	}

	if len(terminals) > 0 {
		for i, _ := range terminals {
			siteErr := db.Debug().Model(&City{}).Unscoped().Where("id = ?", terminals[i].SiteID).Order("id desc").Take(&terminals[i].Site).Error
			if siteErr != nil {
				return &[]Terminal{}, err
			}

			cityErr := db.Debug().Model(&City{}).Unscoped().Where("id = ?", terminals[i].Site.CityID).Order("id desc").Take(&terminals[i].Site.City).Error
			if cityErr != nil {
				return &[]Terminal{}, err
			}

			siteTypeErr := db.Debug().Model(&SiteType{}).Unscoped().Where("id = ?", terminals[i].Site.SiteTypeID).Order("id desc").Take(&terminals[i].Site.SiteType).Error
			if siteTypeErr != nil {
				return &[]Terminal{}, err
			}
		}
	}

	return &terminals, nil
}

func (terminal *Terminal) CreateTerminal(db *gorm.DB) (*Terminal, error) {
	var err error
	tx := db.Begin()
	err = db.Debug().Model(&Terminal{}).Create(&terminal).Error
	if err != nil {
		tx.Rollback()
		return &Terminal{}, err
	}

	err = db.Debug().Model(&Terminal{}).Where("id = ?", terminal.ID).Updates(
		Retailer{
			OriginalID: terminal.ID,
		}).Error
	if err != nil {
		tx.Rollback()
		return &Terminal{}, err
	}

	tx.Commit()
	return terminal, nil
}

func (terminal *Terminal) UpdateTerminal(db *gorm.DB) (*Terminal, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&Terminal{}).Where("id = ?", terminal.ID).Updates(
		Terminal{
			TerminalSN: terminal.TerminalSN,
			TerminalID: terminal.TerminalID,
			MerchantID: terminal.MerchantID,
			UpdatedAt:  &dateTimeNow,
		}).Error

	if err != nil {
		return &Terminal{}, err
	}

	return terminal, nil
}

func (terminal *Terminal) DeactivateTerminalNow(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Site{}).Where("id = ?", terminal.ID).Delete(&Terminal{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (terminal *Terminal) DeactivateTerminalLater(db *gorm.DB) (int64, error) {
	var err error
	err = db.Debug().Model(&Terminal{}).Where("id = ?", terminal.ID).Updates(
		Terminal{
			DeletedAt: terminal.DeletedAt,
		}).Error

	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (terminal *Terminal) ReactivateTerminal(db *gorm.DB) (*Terminal, error) {
	var err error
	tx := db.Begin()
	dateTimeNow := time.Now()
	err = db.Debug().Model(&Terminal{}).Unscoped().Where("id = ?", terminal.ID).Updates(
		Terminal{
			ReactivatedAt: &dateTimeNow,
		}).Error
	if err != nil {
		tx.Rollback()
		return &Terminal{}, err
	}

	terminal.ID = 0
	terminal.DeletedAt = nil
	terminal.ReactivatedAt = nil
	err = db.Debug().Model(&Terminal{}).Create(&terminal).Error
	if err != nil {
		tx.Rollback()
		return &Terminal{}, err
	}

	tx.Commit()
	return terminal, nil
}

/*func (terminal *Terminal) DeactivateTerminal(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Site{}).Where("id = ?", terminal.ID).Delete(&Terminal{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (terminal *Terminal) ReactivateTerminal(db *gorm.DB) (*Terminal, error) {
	var err error
	tx := db.Begin()
	dateTimeNow := time.Now()
	err = db.Debug().Model(&Terminal{}).Unscoped().Where("id = ?", terminal.ID).Updates(
		Retailer{
			ReactivatedAt: &dateTimeNow,
		}).Error
	if err != nil {
		tx.Rollback()
		return &Terminal{}, err
	}

	terminal.ID = 0
	terminal.DeletedAt = nil
	terminal.ReactivatedAt = nil
	err = db.Debug().Model(&Terminal{}).Create(&terminal).Error
	if err != nil {
		tx.Rollback()
		return &Terminal{}, err
	}

	tx.Commit()
	return terminal, nil
}*/

/*func (terminal *Terminal) TerminateTerminalLater(db *gorm.DB) (int64, error) {
	var err error
	err = db.Debug().Model(&Terminal{}).Where("id = ?", terminal.ID).Updates(
		Retailer{
			DeletedAt: terminal.DeletedAt,
		}).Error

	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (terminal *Terminal) TerminateTerminalNow(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Terminal{}).Where("id = ?", terminal.ID).Delete(&Terminal{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}*/

func (Terminal) TableName() string {
	return "terminal"
}
