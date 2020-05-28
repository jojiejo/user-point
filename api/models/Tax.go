package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Tax struct {
	ID         uint64     `gorm:"primary_key;auto_increment" json:"id"`
	TaxTypeID  uint64     `gorm:"not null;" json:"tax_type_id"`
	TaxType    TaxType    `gorm:"not null;" json:"tax_type"`
	ProvinceID uint64     `json:"province_id"`
	Province   Province   `json:"province"`
	Value      float64    `gorm:"not null;" json:"value"`
	StartedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"started_at"`
	EndedAt    *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"ended_at"`
	CreatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt  *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

func (tax *Tax) Prepare() {
	tax.CreatedAt = time.Now()
}

func (tax *Tax) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if tax.TaxTypeID < 1 {
		err = errors.New("Tax type field is required")
		errorMessages["tax_type"] = err.Error()
	}

	if tax.Value <= 0 {
		err = errors.New("Value is required")
		errorMessages["value"] = err.Error()
	}

	return errorMessages
}

func (t *Tax) FindAllTaxes(db *gorm.DB) (*[]Tax, error) {
	var err error
	ts := []Tax{}
	err = db.Debug().Model(&Tax{}).
		Preload("TaxType").
		Preload("Province").
		Order("id, created_at desc").
		Find(&ts).Error

	if err != nil {
		return &[]Tax{}, err
	}

	return &ts, nil
}

func (t *Tax) FindTax(db *gorm.DB, taxID uint64) (*Tax, error) {
	var err error
	err = db.Debug().Model(&Tax{}).Unscoped().
		Preload("TaxType").
		Preload("Province").
		Where("id = ?", taxID).
		Order("id, created_at desc").
		Find(&t).Error

	if err != nil {
		return &Tax{}, err
	}

	return t, nil
}

func (t *Tax) CreateTax(db *gorm.DB) (*Tax, error) {
	var err error
	err = db.Debug().Model(&Tax{}).Create(&t).Error
	if err != nil {
		return &Tax{}, err
	}

	//Select created fee
	_, err = t.FindTax(db, t.ID)
	if err != nil {
		return &Tax{}, err
	}

	return t, nil
}

func (t *Tax) UpdateRelatedTax(db *gorm.DB) error {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&Tax{}).
		Where("tax_type_id = ? AND province_id = ? AND (ended_at IS NULL OR ended_at > ?)", t.TaxTypeID, t.ProvinceID, t.StartedAt).
		Updates(
			map[string]interface{}{
				"ended_at":   t.StartedAt.Add(-1 * time.Second),
				"updated_at": dateTimeNow,
			}).Error

	if err != nil {
		return err
	}

	return nil
}

func (t *Tax) UpdateTax(db *gorm.DB) (*Tax, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&t).Updates(
		map[string]interface{}{
			"ended_at":   t.EndedAt,
			"updated_at": dateTimeNow,
		}).Error

	if err != nil {
		return &Tax{}, err
	}

	//Select updated posting matrix
	_, err = t.FindTax(db, t.ID)
	if err != nil {
		return &Tax{}, err
	}

	return t, nil
}

func (Tax) TableName() string {
	return "tax"
}
