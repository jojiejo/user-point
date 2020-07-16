package models

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
)

type FeeInvoiceByPayer struct {
	Payer             Payer                   `gorm:"not null" json:"payer"`
	InvoiceNumber     FeeInvoiceNumberByPayer `gorm:"not null" json:"invoice_number"`
	InvoiceFeeSummary []FeeInvoiceFeeSummary  `gorm:"not null" json:"invoice_fee_summary"`
	InvoiceTaxSummary []FeeInvoiceTaxSummary  `gorm:"not null" json:"invoice_tax_summary"`
	InvoiceFeeDetail  []FeeInvoiceFeeDetail   `gorm:"not null" json:"invoice_fee_detail"`
}

type FeeInvoiceByBranch struct {
	Branch            Branch                 `gorm:"not null" json:"branch"`
	InvoiceNumber     InvoiceNumberByBranch  `gorm:"not null" json:"invoice_number"`
	InvoiceFeeSummary []FeeInvoiceFeeSummary `gorm:"not null" json:"invoice_fee_summary"`
	InvoiceTaxSummary []FeeInvoiceTaxSummary `gorm:"not null" json:"invoice_tax_summary"`
	InvoiceFeeDetail  []FeeInvoiceFeeDetail  `gorm:"not null" json:"invoice_fee_detail"`
}

type FeeInvoiceNumberByPayer struct {
	CCID              int    `json:"cc_id"`
	IssuedAt          string `json:"issued_at"`
	DueDate           string `json:"due_date"`
	FakturPajakNumber string `json:"faktur_pajak_number"`
	InvoiceNumber     int64  `json:"invoice_number"`
	InvoiceName       string `json:"invoice_name"`
}

type FeeNumberByBranch struct {
	SubCorporateID    int    `json:"sub_corporate_id"`
	IssuedAt          string `json:"issued_at"`
	DueDate           string `json:"due_date"`
	FakturPajakNumber string `json:"faktur_pajak_number"`
	InvoiceNumber     int64  `json:"invoice_number"`
	InvoiceName       string `json:"invoice_name"`
}

type FeeInvoiceFeeSummary struct {
	FeeID         int     `json:"fee_id"`
	FeeCode       string  `json:"fee_code"`
	FeeName       string  `json:"fee_name"`
	Unit          string  `json:"unit"`
	Quantity      float32 `json:"quantity"`
	FeeValue      float32 `json:"fee_value"`
	NetTotal      float32 `json:"net_total"`
	PPNPercentage float32 `json:"ppn_percentage"`
	PPN           float32 `json:"ppn"`
	Total         float32 `json:"total"`
}

type FeeInvoiceTaxSummary struct {
	TaxName       string  `json:"tax_name"`
	TaxPercentage float32 `json:"tax_percentage"`
	NetAmount     float32 `json:"net_amoount"`
	TaxTotal      float32 `json:"tax_total"`
}

type FeeInvoiceFeeDetail struct {
	SubCorporateID   int     `json:"sub_corporate_id"`
	SubCorporateName string  `json:"sub_corporate_name"`
	FeeID            int     `json:"fee_id"`
	FeeCode          string  `json:"fee_code"`
	FeeName          string  `json:"fee_name"`
	Description      string  `json:"description"`
	Unit             string  `json:"unit"`
	Quantity         float32 `json:"quantity"`
	FeeValue         float32 `json:"fee_value"`
	NetTotal         float32 `json:"net_total"`
	PPNPercentage    float32 `json:"ppn_percentage"`
	PPN              float32 `json:"ppn"`
	Total            float32 `json:"total"`
}

func (fibp *FeeInvoiceByPayer) FindFeeInvoiceByCCIDAndDate(db *gorm.DB, CCID uint64, dateFrom string, dateTo string) (*FeeInvoiceByPayer, error) {
	var err error
	err = db.Debug().Model(&Payer{}).
		Preload("GSAPCustomerMasterData.SalesRep").
		Preload("GSAPCustomerMasterData.IndustryClass").
		Preload("GSAPCustomerMasterData.BusinessType").
		Preload("GSAPCustomerMasterData.AccountClass").
		Preload("GSAPCustomerMasterData.PaymentTerm").
		//Preload("FakturPajakInvoice").
		Unscoped().Where("cc_id = ?", CCID).
		Order("cc_id desc").
		Take(&fibp.Payer).Error
	if err != nil {
		return &FeeInvoiceByPayer{}, err
	}

	customerDataErr := db.Debug().Model(&Payer{}).Unscoped().Where("mcms_id = ?", fibp.Payer.MCMSID).Order("mcms_id desc").Take(&fibp.Payer.GSAPCustomerMasterData).Error
	if customerDataErr != nil {
		return &FeeInvoiceByPayer{}, err
	}

	latestStatusErr := db.Debug().Model(&HistoricalPayerStatus{}).Where("cc_id = ?", fibp.Payer.CCID).Order("id asc").Find(&fibp.Payer.LatestPayerStatus).Error
	if latestStatusErr != nil {
		return &FeeInvoiceByPayer{}, err
	}

	statusErr := db.Debug().Model(&PayerStatus{}).Where("id = ?", fibp.Payer.LatestPayerStatus.PayerStatusID).Order("id desc").Take(&fibp.Payer.LatestPayerStatus.PayerStatus).Error
	if statusErr != nil {
		return &FeeInvoiceByPayer{}, err
	}

	branchErr := db.Debug().Model(&ShortenedBranch{}).Where("cc_id = ?", fibp.Payer.CCID).Order("sub_corporate_id desc").Find(&fibp.Payer.Branch).Error
	if branchErr != nil {
		return &FeeInvoiceByPayer{}, err
	}

	fibp.Payer.PaddedMCMSID = fmt.Sprintf("%010v", strconv.Itoa(fibp.Payer.MCMSID))

	/* Load Invoice Number */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Fee_GetInvoiceNumberByPayer ?, ?, ?", CCID, dateFrom, dateTo).Scan(&fibp.InvoiceNumber).Error
	if err != nil {
		return &FeeInvoiceByPayer{}, err
	}

	/* Load Fee Summary */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Fee_GetFeeSummaryByPayer ?, ?, ?", CCID, dateFrom, dateTo).Scan(&fibp.InvoiceFeeSummary).Error
	if err != nil {
		return &FeeInvoiceByPayer{}, err
	}

	/* Load Tax Summary */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Fee_GetTaxSummaryByPayer ?, ?, ?", CCID, dateFrom, dateTo).Scan(&fibp.InvoiceTaxSummary).Error
	if err != nil {
		return &FeeInvoiceByPayer{}, err
	}

	/* Load Details */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Fee_GetFeeDetailByPayer ?, ?, ?", CCID, dateFrom, dateTo).Scan(&fibp.InvoiceFeeDetail).Error
	if err != nil {
		return &FeeInvoiceByPayer{}, err
	}

	return fibp, nil
}

func (fibb *FeeInvoiceByBranch) FindFeeInvoiceBySubCorporateIDAndDate(db *gorm.DB, subCorporateID uint64, dateFrom string, dateTo string) (*FeeInvoiceByBranch, error) {
	var err error
	err = db.Debug().Model(&Branch{}).Unscoped().Where("sub_corporate_id = ?", subCorporateID).Order("created_at desc").Take(&fibb.Branch).Error
	if err != nil {
		return &FeeInvoiceByBranch{}, err
	}

	customerDataErr := db.Debug().Model(&Branch{}).Unscoped().Where("mcms_id = ?", fibb.Branch.MCMSID).Order("mcms_id desc").Take(&fibb.Branch.GSAPCustomerMasterData).Error
	if customerDataErr != nil {
		return &FeeInvoiceByBranch{}, err
	}

	fibb.Branch.PaddedMCMSID = fmt.Sprintf("%010v", strconv.Itoa(fibb.Branch.MCMSID))

	/* Load Invoice Number */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Fee_GetInvoiceNumberByBranch ?, ?, ?", subCorporateID, dateFrom, dateTo).Scan(&fibb.InvoiceNumber).Error
	if err != nil {
		return &FeeInvoiceByBranch{}, err
	}

	/* Load Fee Summary */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Fee_GetFeeSummaryByBranch ?, ?, ?", subCorporateID, dateFrom, dateTo).Scan(&fibb.InvoiceFeeSummary).Error
	if err != nil {
		return &FeeInvoiceByBranch{}, err
	}

	/* Load Tax Summary */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Fee_GetTaxSummaryByBranch ?, ?, ?", subCorporateID, dateFrom, dateTo).Scan(&fibb.InvoiceTaxSummary).Error
	if err != nil {
		return &FeeInvoiceByBranch{}, err
	}

	/* Load Details */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Fee_GetFeeDetailByBranch ?, ?, ?", subCorporateID, dateFrom, dateTo).Scan(&fibb.InvoiceFeeDetail).Error
	if err != nil {
		return &FeeInvoiceByBranch{}, err
	}

	return fibb, nil
}
