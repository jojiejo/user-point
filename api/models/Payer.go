package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Payer struct {
	CCID                     int                             `gorm:"primary_key;auto_increment" json:"cc_id"`
	ContractNumber           string                          `gorm:"not null" json:"contract_number"`
	Alias                    string                          `gorm:"not null" json:"alias"`
	TelematicSubscriptionFee *bool                           `gorm:"not null" json:"telematic_subscription_fee"`
	PaperInvoice             *bool                           `gorm:"not null" json:"paper_invoice"`
	UseInvoiceAddress        *bool                           `gorm:"not null" json:"use_invoice_address"`
	ShowCreditLimit          *bool                           `gorm:"not null" json:"show_credit_limit"`
	InvoiceProductionLevel   *bool                           `gorm:"not null" json:"invoice_production_level"`
	BankVirtualAccount       string                          `gorm:"not null;size:30" json:"bank_virtual_account"`
	CreditLimit              float64                         `gorm:"not null;" json:"credit_limit"`
	MembershipID             *int                            `gorm:"not null" json:"membership_id"`
	MCMSID                   int                             `gorm:"not null;" json:"mcms_id"`
	GSAPCustomerMasterData   DisplayedGSAPCustomerMasterData `json:"gsap_customer_master_data"`
	LatestPayerStatus        HistoricalPayerStatus           `json:"latest_payer_status"`
	CreatedAt                time.Time                       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                *time.Time                      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt                *time.Time                      `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type DisplayedGSAPCustomerMasterData struct {
	MCMSID             int    `json:"mcms_id"`
	PayerNumber        string `json:"payer_number"`
	AgentAccountNumber string `json:"agent_account_number"`
	ContactName_1      string `json:"contact_name_1"`
	ContactName_2      string `json:"contact_name_2"`
	ContactName_3      string `json:"contact_name_3"`
	ContactName_4      string `json:"contact_name_4"`
}

type HistoricalPayerStatus struct {
	ID            int         `json:"id"`
	CCID          int         `json:"cc_id"`
	PayerStatusID int         `json:"payer_status_id"`
	PayerStatus   PayerStatus `json:"payer_status"`
	CreatedAt     time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

type PayerStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (payer *Payer) Prepare() {
	payer.ContractNumber = html.EscapeString(strings.TrimSpace(payer.ContractNumber))
	payer.Alias = html.EscapeString(strings.TrimSpace(payer.Alias))
	payer.BankVirtualAccount = html.EscapeString(strings.TrimSpace(payer.BankVirtualAccount))
}

func (payer *Payer) ValidateInvoiceProduction() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if payer.ShowCreditLimit == nil {
		err = errors.New("Show credit limit field is required")
		errorMessages["required_show_credit_limit"] = err.Error()
	}

	if payer.InvoiceProductionLevel == nil {
		err = errors.New("Invoice production level field is required")
		errorMessages["required_invoice_production_level"] = err.Error()
	}

	return errorMessages
}

func (payer *Payer) ValidateCredit() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if payer.BankVirtualAccount == "" {
		err = errors.New("Bank virtual account field is required")
		errorMessages["required_bank_virtual_account"] = err.Error()
	}

	if payer.CreditLimit < 0 {
		err = errors.New("Credit limit field is required")
		errorMessages["required_credit_limit"] = err.Error()
	}

	return errorMessages
}

func (payer *Payer) ValidateConfiguration() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if payer.ContractNumber == "" {
		err = errors.New("Contract number field is required")
		errorMessages["required_contract_number"] = err.Error()
	}

	if payer.Alias == "" {
		err = errors.New("Alias field is required")
		errorMessages["required_alias"] = err.Error()
	}

	if payer.MembershipID == nil {
		err = errors.New("Membership field is required")
		errorMessages["required_membership"] = err.Error()
	}

	if payer.TelematicSubscriptionFee == nil {
		err = errors.New("Telematic subscription fee field is required")
		errorMessages["required_telematic_subscription_fee"] = err.Error()
	}

	if payer.PaperInvoice == nil {
		err = errors.New("Paper invoice field is required")
		errorMessages["required_paper_invoice"] = err.Error()
	}

	if payer.UseInvoiceAddress == nil {
		err = errors.New("Use invoice address field is required")
		errorMessages["required_use_invoice_address"] = err.Error()
	}

	return errorMessages
}

func (payer *Payer) FindAllPayers(db *gorm.DB) (*[]Payer, error) {
	var err error
	payers := []Payer{}
	err = db.Debug().Model(&Payer{}).Unscoped().Order("created_at desc").Find(&payers).Error
	if err != nil {
		return &[]Payer{}, err
	}

	if len(payers) > 0 {
		for i, _ := range payers {
			customerDataErr := db.Debug().Model(&Payer{}).Unscoped().Where("mcms_id = ?", payers[i].MCMSID).Order("mcms_id desc").Take(&payers[i].GSAPCustomerMasterData).Error
			if customerDataErr != nil {
				return &[]Payer{}, err
			}

			latestStatusErr := db.Debug().Model(&HistoricalPayerStatus{}).Where("cc_id = ?", payers[i].CCID).Order("created_at desc").Find(&payers[i].LatestPayerStatus).Error
			if latestStatusErr != nil {
				return &[]Payer{}, err
			}

			statusErr := db.Debug().Model(&PayerStatus{}).Where("id = ?", payers[i].LatestPayerStatus.PayerStatusID).Order("id desc").Take(&payers[i].LatestPayerStatus.PayerStatus).Error
			if statusErr != nil {
				return &[]Payer{}, err
			}
		}
	}

	return &payers, nil
}

func (payer *Payer) FindPayerByCCID(db *gorm.DB, CCID uint64) (*[]Payer, error) {
	var err error
	payers := []Payer{}
	err = db.Debug().Model(&Payer{}).Unscoped().Where("cc_id = ?", CCID).Order("created_at desc").Find(&payers).Error
	if err != nil {
		return &[]Payer{}, err
	}

	if len(payers) > 0 {
		for i, _ := range payers {
			customerDataErr := db.Debug().Model(&Payer{}).Unscoped().Where("mcms_id = ?", payers[i].MCMSID).Order("mcms_id desc").Take(&payers[i].GSAPCustomerMasterData).Error
			if customerDataErr != nil {
				return &[]Payer{}, err
			}

			latestStatusErr := db.Debug().Model(&HistoricalPayerStatus{}).Where("cc_id = ?", payers[i].CCID).Order("created_at desc").Find(&payers[i].LatestPayerStatus).Error
			if latestStatusErr != nil {
				return &[]Payer{}, err
			}

			statusErr := db.Debug().Model(&PayerStatus{}).Where("id = ?", payers[i].LatestPayerStatus.PayerStatusID).Order("id desc").Take(&payers[i].LatestPayerStatus.PayerStatus).Error
			if statusErr != nil {
				return &[]Payer{}, err
			}
		}
	}

	return &payers, nil
}

func (payer *Payer) UpdatePayerConfiguration(db *gorm.DB) (*Payer, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&Payer{}).Where("cc_id = ?", payer.CCID).Updates(
		Payer{
			ContractNumber:           payer.ContractNumber,
			MembershipID:             payer.MembershipID,
			TelematicSubscriptionFee: payer.TelematicSubscriptionFee,
			PaperInvoice:             payer.PaperInvoice,
			InvoiceProductionLevel:   payer.InvoiceProductionLevel,
			UpdatedAt:                &dateTimeNow,
		}).Error

	if err != nil {
		return &Payer{}, err
	}

	return payer, nil
}

func (payer *Payer) UpdateInvoiceProduction(db *gorm.DB) (*Payer, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&Payer{}).Where("cc_id = ?", payer.CCID).Updates(
		Payer{
			ShowCreditLimit:        payer.ShowCreditLimit,
			InvoiceProductionLevel: payer.InvoiceProductionLevel,
			UpdatedAt:              &dateTimeNow,
		}).Error

	if err != nil {
		return &Payer{}, err
	}

	return payer, nil
}

func (payer *Payer) UpdateCredit(db *gorm.DB) (*Payer, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&Payer{}).Where("cc_id = ?", payer.CCID).Updates(
		Payer{
			CreditLimit:        payer.CreditLimit,
			BankVirtualAccount: payer.BankVirtualAccount,
			UpdatedAt:          &dateTimeNow,
		}).Error

	if err != nil {
		return &Payer{}, err
	}

	return payer, nil
}

func (Payer) TableName() string {
	return "corporate_client_relation"
}

func (DisplayedGSAPCustomerMasterData) TableName() string {
	return "gsap_customer_master_data"
}

func (HistoricalPayerStatus) TableName() string {
	return "payer_status_relation"
}

func (PayerStatus) TableName() string {
	return "payer_status"
}
