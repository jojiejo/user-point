package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type TransactionForManualSettlement struct {
	LastTransactionDateTime time.Time `gorm:"column:last_transaction_datetime" json:"last_transaction_datetime"`
	TerminalID              string    `json:"terminal_id"`
	MerchantID              string    `json:"merchant_id"`
	TotalTransaction        uint64    `json:"total_transaction"`
	TotalAmount             float64   `json:"total_amount"`
	TotalQuantity           float32   `json:"total_quantity"`
	BatchNumber             uint64    `json:"batch_number"`
	SettlementStatus        string    `json:"settlement_status"`
}

func (trx *TransactionForManualSettlement) Prepare() {
	trx.TerminalID = html.EscapeString(strings.TrimSpace(trx.TerminalID))
	trx.MerchantID = html.EscapeString(strings.TrimSpace(trx.MerchantID))
}

func (trx *TransactionForManualSettlement) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if trx.BatchNumber < 1 {
		err = errors.New("Batch number is required")
		errorMessages["required_batch_number"] = err.Error()
	}

	if trx.TerminalID == "" {
		err = errors.New("Terminal ID is required")
		errorMessages["required_batch_number"] = err.Error()
	}

	if trx.MerchantID == "" {
		err = errors.New("Merchant ID is required")
		errorMessages["required_batch_number"] = err.Error()
	}

	return errorMessages
}

func (trx *TransactionForManualSettlement) FindTransactionForManualSettlement(db *gorm.DB, transactionDateFrom string, transactionDateTo string) (*[]TransactionForManualSettlement, error) {
	var err error
	trxs := []TransactionForManualSettlement{}
	err = db.Debug().Raw("EXEC spAPI_ManualSettlement_GetTrxByDate ?, ?", transactionDateFrom, transactionDateTo).Scan(&trxs).Error
	if err != nil {
		return &[]TransactionForManualSettlement{}, err
	}

	return &trxs, nil
}

func (trx *TransactionForManualSettlement) ManualSettle(db *gorm.DB) (*TransactionForManualSettlement, error) {
	var err error
	dateTimeNow := time.Now()
	parsedDateTimeNow := dateTimeNow.Format("2006-01-02 15:04:05")
	dashReplacedDateTimeNow := strings.Replace(parsedDateTimeNow, "-", "", -1)
	colonReplacedDateTimeNow := strings.Replace(dashReplacedDateTimeNow, ":", "", -1)
	splitDateTimeNow := strings.Split(colonReplacedDateTimeNow, " ")
	clearDate := splitDateTimeNow[0]
	clearTime := splitDateTimeNow[1]

	saleSettleStatus := "N"
	requestMTI := "0200"
	saleStatus := "SA"
	err = db.Debug().Table("EDCTransactionMonitor").
		Where("terminal_id = ? AND merchant_id = ? AND batch_no = ? AND sale_settle_status = ? AND request_mti = ? AND sale_status = ?",
			trx.TerminalID, trx.MerchantID, trx.BatchNumber, saleSettleStatus, requestMTI, saleStatus).
		Updates(
			map[string]interface{}{
				"sale_settle_status": "Y",
				"date_update":        dateTimeNow,
				"sale_settle_date":   clearDate,
				"sale_settle_time":   clearTime,
			}).Error

	if err != nil {
		return &TransactionForManualSettlement{}, err
	}

	return trx, nil
}

func (trx *TransactionForManualSettlement) ManualSettleAllTransaction(db *gorm.DB) error {
	var err error
	dateTimeNow := time.Now()
	parsedDateTimeNow := dateTimeNow.Format("2006-01-02 15:04:05")
	dashReplacedDateTimeNow := strings.Replace(parsedDateTimeNow, "-", "", -1)
	colonReplacedDateTimeNow := strings.Replace(dashReplacedDateTimeNow, ":", "", -1)
	splitDateTimeNow := strings.Split(colonReplacedDateTimeNow, " ")
	clearDate := splitDateTimeNow[0]
	clearTime := splitDateTimeNow[1]

	saleSettleStatus := "N"
	requestMTI := "0200"
	saleStatus := "SA"
	err = db.Debug().Table("EDCTransactionMonitor").
		Where("sale_settle_status = ? AND request_mti = ? AND sale_status = ?",
			saleSettleStatus, requestMTI, saleStatus).
		Updates(
			map[string]interface{}{
				"sale_settle_status": "Y",
				"date_update":        dateTimeNow,
				"sale_settle_date":   clearDate,
				"sale_settle_time":   clearTime,
			}).Error

	if err != nil {
		return err
	}

	return nil
}

func (TransactionForManualSettlement) TableName() string {
	return "EDCTransactionMonitor"
}
