package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Terminal struct {
	ID         int       `gorm:"primary_key;auto_increment" json:"id"`
	TerminalSN string    `gorm:"not null"; json:"terminal_sn"`
	TerminalID string    `gorm:"not null"; json:"terminal_id"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (terminal *Terminal) Prepare() {
	terminal.TerminalSN = html.EscapeString(strings.TrimSpace(terminal.TerminalSN))
	terminal.TerminalID = html.EscapeString(strings.TrimSpace(terminal.TerminalID))
	terminal.CreatedAt = time.Now()
}

func (terminal *Terminal) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if terminal.TerminalSN == "" {
		err = errors.New("Terminal SN field is required")
		errorMessages["required_terminal_sn"] = err.Error()
	}

	if terminal.TerminalID == "" {
		err = errors.New("Terminal ID field is required")
		errorMessages["required_terminal_id"] = err.Error()
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

func (terminal *Terminal) FindTerminalByID(db *gorm.DB, terminalID uint64) (*Terminal, error) {
	var err error
	err = db.Debug().Model(&Terminal{}).Where("id = ?", terminalID).Order("created_at desc").Take(&terminal).Error
	if err != nil {
		return &Terminal{}, err
	}

	return terminal, nil
}

func (terminal *Terminal) CreateTerminal(db *gorm.DB) (*Terminal, error) {
	var err error
	err = db.Debug().Model(&Terminal{}).Create(&terminal).Error
	if err != nil {
		return &Terminal{}, err
	}

	return terminal, nil
}

func (terminal *Terminal) UpdateTerminal(db *gorm.DB) (*Terminal, error) {
	var err error

	err = db.Debug().Model(&Site{}).Where("id = ?", terminal.ID).Updates(
		Terminal{
			TerminalSN: terminal.TerminalSN,
			TerminalID: terminal.TerminalID,
			UpdatedAt:  time.Now(),
		}).Error

	if err != nil {
		return &Terminal{}, err
	}

	return terminal, nil
}

func (terminal *Terminal) DeleteTerminal(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Terminal{}).Where("id = ?", terminal.ID).Delete(&Terminal{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (Terminal) TableName() string {
	return "terminal"
}
