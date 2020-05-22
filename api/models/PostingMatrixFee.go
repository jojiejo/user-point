package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type PostingMatrixFee struct {
	ID              uint64           `gorm:"primary_key;auto_increment" json:"id"`
	FeeID           uint64           `gorm:"not null" json:"fee_id"`
	Fee             ShortenedFee     `json:"fee"`
	VolumeProceedGL string           `gorm:"not null" json:"volume_proceed_gl"`
	DiscountGL      string           `gorm:"not null" json:"discount_gl"`
	SurchargeGL     string           `gorm:"not null" json:"surcharge_gl"`
	ProfitCentre    string           `gorm:"not null" json:"profit_centre"`
	MaterialCode    string           `gorm:"not null" json:"material_code"`
	VATID           string           `gorm:"column:posting_matrix_vat_id" json:"vat_id"`
	VAT             PostingMatrixVAT `gorm:"ForeignKey:VATID;AssociationForeignKey:ID" json:"vat"`
	CreatedAt       time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       *time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (pmf *PostingMatrixFee) Prepare() {
	pmf.VolumeProceedGL = html.EscapeString(strings.TrimSpace(pmf.VolumeProceedGL))
	pmf.DiscountGL = html.EscapeString(strings.TrimSpace(pmf.DiscountGL))
	pmf.SurchargeGL = html.EscapeString(strings.TrimSpace(pmf.SurchargeGL))
	pmf.ProfitCentre = html.EscapeString(strings.TrimSpace(pmf.ProfitCentre))
	pmf.MaterialCode = html.EscapeString(strings.TrimSpace(pmf.MaterialCode))
	pmf.CreatedAt = time.Now()
}

func (pmf *PostingMatrixFee) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if pmf.FeeID < 1 {
		err = errors.New("Product field is required")
		errorMessages["product"] = err.Error()
	}

	if pmf.VolumeProceedGL == "" {
		err = errors.New("Volume / Proceed GL field is required")
		errorMessages["volume_proceed_gl"] = err.Error()
	}

	if pmf.DiscountGL == "" {
		err = errors.New("Discount GL field is required")
		errorMessages["discount_gl"] = err.Error()
	}

	if pmf.SurchargeGL == "" {
		err = errors.New("Surcharge GL field is required")
		errorMessages["surcharge_gl"] = err.Error()
	}

	if pmf.ProfitCentre == "" {
		err = errors.New("Profit Centre field is required")
		errorMessages["profit_centre"] = err.Error()
	}

	if pmf.SurchargeGL == "" {
		err = errors.New("Material code field is required")
		errorMessages["material_code"] = err.Error()
	}

	return errorMessages
}

func (pmf *PostingMatrixFee) FindPostingMatrixFees(db *gorm.DB) (*[]PostingMatrixFee, error) {
	var err error
	pmfs := []PostingMatrixFee{}
	err = db.Debug().Model(&PostingMatrixFee{}).
		Preload("Fee").
		Preload("Fee.Unit").
		Preload("VAT").
		Order("fee_id, created_at desc").
		Find(&pmfs).Error

	if err != nil {
		return &[]PostingMatrixFee{}, err
	}

	return &pmfs, nil
}

func (pmf *PostingMatrixFee) FindPostingMatrixFee(db *gorm.DB, pmfID uint64) (*PostingMatrixFee, error) {
	var err error
	err = db.Debug().Model(&PostingMatrixFee{}).Unscoped().
		Preload("Fee").
		Preload("Fee.Unit").
		Preload("VAT").
		Where("id = ?", pmfID).
		Order("id, created_at desc").
		Find(&pmf).Error

	if err != nil {
		return &PostingMatrixFee{}, err
	}

	return pmf, nil
}

func (pmf *PostingMatrixFee) CreatePostingMatrixFee(db *gorm.DB) (*PostingMatrixFee, error) {
	var err error
	err = db.Debug().Model(&PostingMatrixFee{}).Create(&pmf).Error
	if err != nil {
		return &PostingMatrixFee{}, err
	}

	//Select created fee
	_, err = pmf.FindPostingMatrixFee(db, pmf.ID)
	if err != nil {
		return &PostingMatrixFee{}, err
	}

	return pmf, nil
}

func (pmf *PostingMatrixFee) UpdatePostingMatrixFee(db *gorm.DB) (*PostingMatrixFee, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&pmf).Updates(
		map[string]interface{}{
			"volume_proceed_gl":     pmf.VolumeProceedGL,
			"discount_gl":           pmf.DiscountGL,
			"surcharge_gl":          pmf.SurchargeGL,
			"profit_centre":         pmf.ProfitCentre,
			"material_code":         pmf.MaterialCode,
			"posting_matrix_vat_id": pmf.VATID,
			"updated_at":            dateTimeNow,
		}).Error

	if err != nil {
		return &PostingMatrixFee{}, err
	}

	//Select updated posting matrix
	_, err = pmf.FindPostingMatrixFee(db, pmf.ID)
	if err != nil {
		return &PostingMatrixFee{}, err
	}

	return pmf, nil
}

func (PostingMatrixFee) TableName() string {
	return "posting_matrix_fee"
}
