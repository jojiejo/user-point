package models

import (
	"time"
)

type FakturPajakInvoice struct {
	ID                 uint64     `gorm:"primary_key;auto_increment" json:"id"`
	FakturPajakRangeID string     `gorm:"not null" json:"faktur_pajak_range_id"`
	FakturPajakNumber  string     `gorm:"not null" json:"faktur_pajak_number"`
	InvoiceNumber      string     `gorm:"not null" json:"invoice_number"`
	InvoiceDate        string     `gorm:"not null" json:"invoice_date"`
	CCID               uint64     `gorm:"not null" json:"cc_id"`
	CreatedAt          time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt          *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (FakturPajakInvoice) TableName() string {
	return "faktur_pajak_invoice"
}
