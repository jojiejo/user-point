package models

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
)

type TransactionInvoiceByPayer struct {
	Payer                   Payer                     `gorm:"not null" json:"payer"`
	InvoiceNumber           InvoiceNumberByPayer      `gorm:"not null" json:"invoice_number"`
	InvoiceProductSummary   []InvoiceProductSummary   `gorm:"not null" json:"invoice_product_summary"`
	InvoiceTaxSummary       []InvoiceTaxSummary       `gorm:"not null" json:"invoice_tax_summary"`
	InvoiceSubAccountDetail []InvoiceSubAccountDetail `gorm:"not null" json:"invoice_sub_account_detail"`
}

type TransactionInvoiceByBranch struct {
	Branch                  Branch                    `gorm:"not null" json:"branch"`
	InvoiceNumber           InvoiceNumberByBranch     `gorm:"not null" json:"invoice_number"`
	InvoiceProductSummary   []InvoiceProductSummary   `gorm:"not null" json:"invoice_product_summary"`
	InvoiceTaxSummary       []InvoiceTaxSummary       `gorm:"not null" json:"invoice_tax_summary"`
	InvoiceSubAccountDetail []InvoiceSubAccountDetail `gorm:"not null" json:"invoice_sub_account_detail"`
}

type InvoiceNumberByPayer struct {
	CCID              int    `json:"cc_id"`
	IssuedAt          string `json:"issued_at"`
	DueDate           string `json:"due_date"`
	FakturPajakNumber string `json:"faktur_pajak_number"`
	InvoiceNumber     int64  `json:"invoice_number"`
	InvoiceName       string `json:"invoice_name"`
}

type InvoiceNumberByBranch struct {
	SubCorporateID    int    `json:"sub_corporate_id"`
	IssuedAt          string `json:"issued_at"`
	DueDate           string `json:"due_date"`
	FakturPajakNumber string `json:"faktur_pajak_number"`
	InvoiceNumber     int64  `json:"invoice_number"`
	InvoiceName       string `json:"invoice_name"`
}

type InvoiceSubAccountSummary struct {
	CCID           int     `json:"cc_id"`
	SubCorporateID int     `json:"sub_corporate_id"`
	ContactName_3  string  `json:"contact_name_3"`
	ContactName_4  string  `json:"contact_name_4"`
	SaleAmount     float32 `json:"sale_amount"`
	IssuedAt       string  `json:"issued_at"`
	DueDate        string  `json:"due_date"`
	InvoiceNumber  int64   `json:"invoice_number"`
	InvoiceName    string  `json:"invoice_name"`
}

type InvoiceProductSummary struct {
	ProductName     string  `json:"product_name"`
	ProductCode     string  `json:"product_code"`
	TotalVolume     float32 `json:"total_volume"`
	TotalAmount     float32 `json:"total_amount"`
	PPNPercentage   float32 `json:"ppn_percentage"`
	PPNAmount       float32 `json:"ppn_amount"`
	PPHPercentage   float32 `json:"pph_percentage"`
	PPHAmount       float32 `json:"pph_amount"`
	PBBKBPercentage float32 `json:"pbbkb_percentage"`
	PBBKBAmount     float32 `json:"pbbkb_amount"`
	TotalTax        float32 `json:"total_tax"`
	Total           float32 `json:"total"`
}

type InvoiceTaxSummary struct {
	TaxName       string  `json:"tax_name"`
	TaxPercentage float32 `json:"tax_percentage"`
	NetAmount     float32 `json:"net_amoount"`
	TaxTotal      float32 `json:"tax_total"`
}

type InvoiceSubAccountDetail struct {
	SubCorporateID   int     `json:"sub_corporate_id"`
	SubCorporateName string  `json:"sub_corporate_name"`
	CardGroupCode    int     `json:"card_group_code"`
	CardGroupName    string  `json:"card_group_name"`
	CardID           string  `json:"card_id"`
	TransactionDate  string  `json:"transaction_date"`
	TransactionTime  string  `json:"transaction_time"`
	SiteID           int     `json:"site_id"`
	ShipToName       string  `json:"ship_to_name"`
	ReferenceNo      string  `json:"reference_no"`
	SaleOdometer     string  `json:"sale_odometer"`
	SaleProductCode  string  `json:"sale_product_code"`
	ProductName      string  `json:"product_name"`
	SaleTotalVolume  float32 `json:"sale_total_volume"`
	Price            float32 `json:"price"`
	NetAmount        float32 `json:"net_amount"`
	PPNPercentage    float32 `json:"ppn_percentage"`
	PPNAmount        float32 `json:"ppn_amount"`
	PPHPercentage    float32 `json:"pph_percentage"`
	PPHAmount        float32 `json:"pph_amount"`
	PBBKBPercentage  float32 `json:"pbbkb_percentage"`
	PBBKBAmount      float32 `json:"pbbkb_amount"`
	TaxTotal         float32 `json:"tax_total"`
	SaleAmount       float32 `json:"sale_amount"`
}

func (tibp *TransactionInvoiceByPayer) FindInvoiceByCCIDAndDate(db *gorm.DB, CCID uint64, dateFrom string, dateTo string) (*TransactionInvoiceByPayer, error) {
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
		Take(&tibp.Payer).Error
	if err != nil {
		return &TransactionInvoiceByPayer{}, err
	}

	customerDataErr := db.Debug().Model(&Payer{}).Unscoped().Where("mcms_id = ?", tibp.Payer.MCMSID).Order("mcms_id desc").Take(&tibp.Payer.GSAPCustomerMasterData).Error
	if customerDataErr != nil {
		return &TransactionInvoiceByPayer{}, err
	}

	latestStatusErr := db.Debug().Model(&HistoricalPayerStatus{}).Where("cc_id = ?", tibp.Payer.CCID).Order("id asc").Find(&tibp.Payer.LatestPayerStatus).Error
	if latestStatusErr != nil {
		return &TransactionInvoiceByPayer{}, err
	}

	statusErr := db.Debug().Model(&PayerStatus{}).Where("id = ?", tibp.Payer.LatestPayerStatus.PayerStatusID).Order("id desc").Take(&tibp.Payer.LatestPayerStatus.PayerStatus).Error
	if statusErr != nil {
		return &TransactionInvoiceByPayer{}, err
	}

	branchErr := db.Debug().Model(&ShortenedBranch{}).Where("cc_id = ?", tibp.Payer.CCID).Order("sub_corporate_id desc").Find(&tibp.Payer.Branch).Error
	if branchErr != nil {
		return &TransactionInvoiceByPayer{}, err
	}

	tibp.Payer.PaddedMCMSID = fmt.Sprintf("%010v", strconv.Itoa(tibp.Payer.MCMSID))

	/* Load Invoice Number */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Trx_GetInvoiceNumberByPayer ?, ?, ?", CCID, dateFrom, dateTo).Scan(&tibp.InvoiceNumber).Error
	if err != nil {
		return &TransactionInvoiceByPayer{}, err
	}

	/* Load Product Summary */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Trx_GetProductTrxSummaryByPayer ?, ?, ?", CCID, dateFrom, dateTo).Scan(&tibp.InvoiceProductSummary).Error
	if err != nil {
		return &TransactionInvoiceByPayer{}, err
	}

	/* Load Tax Summary */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Trx_GetTaxTrxSummaryByPayer ?, ?, ?", CCID, dateFrom, dateTo).Scan(&tibp.InvoiceTaxSummary).Error
	if err != nil {
		return &TransactionInvoiceByPayer{}, err
	}

	/* Load Details */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Trx_GetSubAccountTrxDetailByPayer ?, ?, ?", CCID, dateFrom, dateTo).Scan(&tibp.InvoiceSubAccountDetail).Error
	if err != nil {
		return &TransactionInvoiceByPayer{}, err
	}

	return tibp, nil
}

func (tibb *TransactionInvoiceByBranch) FindInvoiceBySubCorporateIDAndDate(db *gorm.DB, subCorporateID uint64, dateFrom string, dateTo string) (*TransactionInvoiceByBranch, error) {
	var err error
	err = db.Debug().Model(&Branch{}).Unscoped().Where("sub_corporate_id = ?", subCorporateID).Order("created_at desc").Take(&tibb.Branch).Error
	if err != nil {
		return &TransactionInvoiceByBranch{}, err
	}

	customerDataErr := db.Debug().Model(&Payer{}).Unscoped().Where("mcms_id = ?", tibb.Branch.MCMSID).Order("mcms_id desc").Take(&tibb.Branch.GSAPCustomerMasterData).Error
	if customerDataErr != nil {
		return &TransactionInvoiceByBranch{}, err
	}

	tibb.Branch.PaddedMCMSID = fmt.Sprintf("%010v", strconv.Itoa(tibb.Branch.MCMSID))

	/* Load Invoice Number */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Trx_GetInvoiceNumberByBranch ?, ?, ?", subCorporateID, dateFrom, dateTo).Scan(&tibb.InvoiceNumber).Error
	if err != nil {
		return &TransactionInvoiceByBranch{}, err
	}

	/* Load Product Summary */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Trx_GetProductTrxSummaryByBranch ?, ?, ?", subCorporateID, dateFrom, dateTo).Scan(&tibb.InvoiceProductSummary).Error
	if err != nil {
		return &TransactionInvoiceByBranch{}, err
	}

	/* Load Tax Summary */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Trx_GetTaxTrxSummaryByBranch ?, ?, ?", subCorporateID, dateFrom, dateTo).Scan(&tibb.InvoiceTaxSummary).Error
	if err != nil {
		return &TransactionInvoiceByBranch{}, err
	}

	/* Load Details */
	err = db.Debug().Raw("EXEC spAPI_Invoice_Trx_GetSubAccountTrxDetailByBranch ?, ?, ?", subCorporateID, dateFrom, dateTo).Scan(&tibb.InvoiceSubAccountDetail).Error
	if err != nil {
		return &TransactionInvoiceByBranch{}, err
	}

	return tibb, nil
}
