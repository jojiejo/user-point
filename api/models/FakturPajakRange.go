package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type FakturPajakRange struct {
	ID                   uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Prefix               string     `gorm:"not null" json:"prefix"`
	StartRange           uint64     `gorm:"not null" json:"start_range"`
	EndRange             uint64     `gorm:"not null" json:"end_range"`
	TotalRange           uint64     `gorm:"not null" json:"total_range"`
	UsedFakturPajakRange uint64     `gorm:"-" json:"used_faktur_pajak_range"`
	AvailableRange       uint64     `gorm:"-" json:"available_range"`
	CreatedAt            time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt            *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type AvailableFakturPajakRange struct {
	AvailableFakturPajakRange uint64 `gorm:"not null" json:"available_faktur_pajak_range"`
}

func (fpr *FakturPajakRange) Prepare() {
	fpr.Prefix = html.EscapeString(strings.TrimSpace(fpr.Prefix))
	fpr.CreatedAt = time.Now()
}

func (fpr *FakturPajakRange) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if fpr.StartRange < 1 {
		err = errors.New("Start range field is required")
		errorMessages["start_range"] = err.Error()
	}

	if fpr.EndRange < 1 {
		err = errors.New("End range field is required")
		errorMessages["end_range"] = err.Error()
	}

	if fpr.EndRange < fpr.StartRange {
		err = errors.New("End range field cannot be lower than Start range field")
		errorMessages["lower_end_range"] = err.Error()
	}

	if fpr.TotalRange < 1 {
		err = errors.New("Total range field is required")
		errorMessages["total_range"] = err.Error()
	}

	if fpr.Prefix == "" {
		err = errors.New("Prefix field is required")
		errorMessages["prefix"] = err.Error()
	}

	return errorMessages
}

func (fpr *FakturPajakRange) FindFakturPajakRanges(db *gorm.DB) (*[]FakturPajakRange, error) {
	var err error
	fprs := []FakturPajakRange{}
	err = db.Debug().Raw("EXEC spAPI_FakturPajak_GetAllRangeWithAvailability").Scan(&fprs).Error
	if err != nil {
		return &[]FakturPajakRange{}, err
	}

	return &fprs, nil
}

func (fpr *FakturPajakRange) FindFakturPajakNextAvailableNumber(db *gorm.DB) (*FakturPajakRange, error) {
	err := db.Debug().Raw("EXEC spFakturPajak_GetNextAvailableNumber").Scan(&fpr).Error
	if err != nil {
		return &FakturPajakRange{}, err
	}

	return fpr, nil
}

func (afpr *AvailableFakturPajakRange) FindAvailableFakturPajakRange(db *gorm.DB) (*AvailableFakturPajakRange, error) {
	err := db.Debug().Raw("EXEC spFakturPajak_GetAvailableNumber").Scan(&afpr).Error
	if err != nil {
		return &AvailableFakturPajakRange{}, err
	}

	return afpr, nil
}

func (fpr *FakturPajakRange) CreateFakturPajakRange(db *gorm.DB) (*FakturPajakRange, error) {
	var err error
	err = db.Debug().Model(&FakturPajakRange{}).Create(&fpr).Error
	if err != nil {
		return &FakturPajakRange{}, err
	}

	return fpr, nil
}

func (FakturPajakRange) TableName() string {
	return "faktur_pajak_range"
}
