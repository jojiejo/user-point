package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type PostingMatrixProduct struct {
	ID              uint64           `gorm:"primary_key;auto_increment" json:"id"`
	ProductID       uint64           `gorm:"not null" json:"product_id"`
	Product         Product          `json:"product"`
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

func (pmp *PostingMatrixProduct) Prepare() {
	pmp.VolumeProceedGL = html.EscapeString(strings.TrimSpace(pmp.VolumeProceedGL))
	pmp.DiscountGL = html.EscapeString(strings.TrimSpace(pmp.DiscountGL))
	pmp.SurchargeGL = html.EscapeString(strings.TrimSpace(pmp.SurchargeGL))
	pmp.ProfitCentre = html.EscapeString(strings.TrimSpace(pmp.ProfitCentre))
	pmp.MaterialCode = html.EscapeString(strings.TrimSpace(pmp.MaterialCode))
	pmp.CreatedAt = time.Now()
}

func (pmp *PostingMatrixProduct) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if pmp.ProductID < 1 {
		err = errors.New("Product field is required")
		errorMessages["product"] = err.Error()
	}

	if pmp.VolumeProceedGL == "" {
		err = errors.New("Volume / Proceed GL field is required")
		errorMessages["volume_proceed_gl"] = err.Error()
	}

	if pmp.DiscountGL == "" {
		err = errors.New("Discount GL field is required")
		errorMessages["discount_gl"] = err.Error()
	}

	if pmp.SurchargeGL == "" {
		err = errors.New("Surcharge GL field is required")
		errorMessages["surcharge_gl"] = err.Error()
	}

	if pmp.ProfitCentre == "" {
		err = errors.New("Profit Centre field is required")
		errorMessages["profit_centre"] = err.Error()
	}

	if pmp.SurchargeGL == "" {
		err = errors.New("Material code field is required")
		errorMessages["material_code"] = err.Error()
	}

	return errorMessages
}

func (pmp *PostingMatrixProduct) FindPostingMatrixProducts(db *gorm.DB) (*[]PostingMatrixProduct, error) {
	var err error
	pmps := []PostingMatrixProduct{}
	err = db.Debug().Model(&PostingMatrixProduct{}).
		Preload("Product").
		Preload("VAT").
		Order("product_id, created_at desc").
		Find(&pmps).Error

	if err != nil {
		return &[]PostingMatrixProduct{}, err
	}

	return &pmps, nil
}

func (pmp *PostingMatrixProduct) FindPostingMatrixProduct(db *gorm.DB, pmpID uint64) (*PostingMatrixProduct, error) {
	var err error
	err = db.Debug().Model(&PostingMatrixProduct{}).Unscoped().
		Preload("Product").
		Preload("VAT").
		Where("id = ?", pmpID).
		Order("id, created_at desc").
		Find(&pmp).Error

	if err != nil {
		return &PostingMatrixProduct{}, err
	}

	return pmp, nil
}

func (pmp *PostingMatrixProduct) CreatePostingMatrixProduct(db *gorm.DB) (*PostingMatrixProduct, error) {
	var err error
	err = db.Debug().Model(&PostingMatrixProduct{}).Create(&pmp).Error
	if err != nil {
		return &PostingMatrixProduct{}, err
	}

	//Select created fee
	_, err = pmp.FindPostingMatrixProduct(db, pmp.ID)
	if err != nil {
		return &PostingMatrixProduct{}, err
	}

	return pmp, nil
}

func (pmp *PostingMatrixProduct) UpdatePostingMatrixProduct(db *gorm.DB) (*PostingMatrixProduct, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&pmp).Updates(
		map[string]interface{}{
			"volume_proceed_gl":     pmp.VolumeProceedGL,
			"discount_gl":           pmp.DiscountGL,
			"surcharge_gl":          pmp.SurchargeGL,
			"profit_centre":         pmp.ProfitCentre,
			"material_code":         pmp.MaterialCode,
			"posting_matrix_vat_id": pmp.VATID,
			"updated_at":            dateTimeNow,
		}).Error

	if err != nil {
		return &PostingMatrixProduct{}, err
	}

	//Select updated posting matrix
	_, err = pmp.FindPostingMatrixProduct(db, pmp.ID)
	if err != nil {
		return &PostingMatrixProduct{}, err
	}

	return pmp, nil
}

func (PostingMatrixProduct) TableName() string {
	return "posting_matrix_product"
}
