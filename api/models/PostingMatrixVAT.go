package models

import (
	"time"
)

type PostingMatrixVAT struct {
	ID          uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Category    string     `gorm:"not null" json:"category"`
	Description string     `gorm:"not null" json:"description"`
	Percentage  float64    `gorm:"not null" json:"percentage"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt   *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

func (PostingMatrixVAT) TableName() string {
	return "posting_matrix_vat"
}
