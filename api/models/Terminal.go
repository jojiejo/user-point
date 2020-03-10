package models

import (
	"errors"
	"fmt"
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

	fmt.Println("Terminal SN : ", terminal.TerminalSN)
	if terminal.TerminalSN == "" {
		err = errors.New("Terminal SN field is required")
		errorMessages["required_terminal_sn"] = err.Error()
	}

	if terminal.TerminalID == "" {
		err = errors.New("Terminal ID field is required")
		errorMessages["required_terminal_id"] = err.Error()
	}

	if terminal.MerchantID == "" {
		err = errors.New("Merchant ID field is required")
		errorMessages["required_merchant_id"] = err.Error()
	}

	return errorMessages
}

func (terminal *Terminal) FindAllTerminals(db *gorm.DB) (*[]Terminal, error) {
	var err error
	terminals := []Terminal{}
	err = db.Debug().Model(&Terminal{}).Limit(100).Order("created_at desc").Find(&terminals).Error
	if err != nil {
		return &[]Terminal{}, err
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

	return &terminals, nil
}

func (terminal *Terminal) FindTerminalByID(db *gorm.DB, terminalID uint64) (*Terminal, error) {
	var err error
	err = db.Debug().Model(&Terminal{}).Where("id = ?", terminalID).Order("created_at desc").Take(&terminal).Error
	if err != nil {
		return &Terminal{}, err
	}

	return terminal, nil
}

func (terminal *Terminal) FindTerminalHistoryByID(db *gorm.DB, originalRetailerID uint64) (*[]Terminal, error) {
	var err error
	var terminals = []Terminal{}
	err = db.Debug().Model(&Terminal{}).Unscoped().Where("original_id = ?", originalRetailerID).Order("created_at desc").Find(&terminals).Error
	if err != nil {
		return &[]Terminal{}, err
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

func (terminal *Terminal) DeactivateTerminal(db *gorm.DB) (int64, error) {
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
}

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
