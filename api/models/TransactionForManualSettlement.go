package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type TransactionForManualSettlement struct {
	LastTransactionDateTime time.Time `gorm:"column:last_transaction_datetime" json:"last_transaction_datetime"`
	TerminalID              string    `json:"terminal_id"`
	MerchantID              string    `json:"merchant_id"`
	TotalTransaction        uint64    `json:"total_transaction"`
	TotalAmount             float64   `json:"total_amount"`
	TotalQuantity           float64   `json:"total_quantity"`
	BatchNumber             uint64    `json:"batch_number"`
	SettlementStatus        string    `json:"settlement_status"`
}

func (trx *TransactionForManualSettlement) FindTransactionForManualSettlement(db *gorm.DB, transactionDate string) (*[]TransactionForManualSettlement, error) {
	var err error
	trxs := []TransactionForManualSettlement{}
	err = db.Debug().Raw("EXEC spAPI_ManualSettlement_GetTrxByDate ?", transactionDate).Scan(&trxs).Error
	if err != nil {
		return &[]TransactionForManualSettlement{}, err
	}

	return &trxs, nil
}

func (trx *TransactionForManualSettlement) ManualSettle(db *gorm.DB) (*TransactionForManualSettlement, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Table(&trx).Omit("name").Updates(
		map[string]interface{}{
			"sale_settle_status": "Y",
			"updated_at":         dateTimeNow,
		}).Error

	if err != nil {
		return &TransactionForManualSettlement{}, err
	}

	return trx, nil
}

func (TransactionForManualSettlement) TableName() string {
	return "EDCTransactionMonitor"
}
