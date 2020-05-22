package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type PostingMatrixTax struct {
	ID              uint64     `gorm:"primary_key;auto_increment" json:"id"`
	GLName          string     `gorm:"not null" json:"gl_name"`
	GLAccountNumber string     `gorm:"not null" json:"gl_account_number"`
	Description     string     `json:"description"`
	UsedIn          string     `gorm:"not null" json:"used_in"`
	CreatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (pmt *PostingMatrixTax) Prepare() {
	pmt.GLName = html.EscapeString(strings.TrimSpace(pmt.GLName))
	pmt.GLAccountNumber = html.EscapeString(strings.TrimSpace(pmt.GLAccountNumber))
	pmt.Description = html.EscapeString(strings.TrimSpace(pmt.Description))
	pmt.CreatedAt = time.Now()
}

func (pmt *PostingMatrixTax) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if pmt.GLName == "" {
		err = errors.New("General ledger name field is required")
		errorMessages["general_ledger_name"] = err.Error()
	}

	if pmt.GLAccountNumber == "" {
		err = errors.New("General ledger account number is required")
		errorMessages["general_ledger_account_number"] = err.Error()
	}

	if pmt.Description == "" {
		err = errors.New("Description field is required")
		errorMessages["description"] = err.Error()
	}

	if pmt.UsedIn == "" {
		err = errors.New("Used in field is required")
		errorMessages["used_in"] = err.Error()
	}

	return errorMessages
}

func (pmt *PostingMatrixTax) FindPostingMatrixTaxes(db *gorm.DB) (*[]PostingMatrixTax, error) {
	var err error
	pmts := []PostingMatrixTax{}
	err = db.Debug().Model(&PostingMatrixTax{}).
		Order("id, created_at desc").
		Find(&pmts).Error

	if err != nil {
		return &[]PostingMatrixTax{}, err
	}

	return &pmts, nil
}

func (pmt *PostingMatrixTax) FindPostingMatrixTax(db *gorm.DB, pmtID uint64) (*PostingMatrixTax, error) {
	var err error
	err = db.Debug().Model(&PostingMatrixTax{}).Unscoped().
		Where("id = ?", pmtID).
		Order("id, created_at desc").
		Find(&pmt).Error

	if err != nil {
		return &PostingMatrixTax{}, err
	}

	return pmt, nil
}

func (pmt *PostingMatrixTax) CreatePostingMatrixTax(db *gorm.DB) (*PostingMatrixTax, error) {
	var err error
	err = db.Debug().Model(&PostingMatrixProduct{}).Create(&pmt).Error
	if err != nil {
		return &PostingMatrixTax{}, err
	}

	//Select created fee
	_, err = pmt.FindPostingMatrixTax(db, pmt.ID)
	if err != nil {
		return &PostingMatrixTax{}, err
	}

	return pmt, nil
}

func (pmt *PostingMatrixTax) UpdatePostingMatrixTax(db *gorm.DB) (*PostingMatrixTax, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&pmt).Updates(
		map[string]interface{}{
			"gl_name":           pmt.GLName,
			"gl_account_number": pmt.GLAccountNumber,
			"description":       pmt.Description,
			"used_in":           pmt.UsedIn,
			"updated_at":        dateTimeNow,
		}).Error

	if err != nil {
		return &PostingMatrixTax{}, err
	}

	//Select updated posting matrix
	_, err = pmt.FindPostingMatrixTax(db, pmt.ID)
	if err != nil {
		return &PostingMatrixTax{}, err
	}

	return pmt, nil
}

func (PostingMatrixTax) TableName() string {
	return "posting_matrix_tax"
}
