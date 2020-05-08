package models

import (
	"errors"
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Payer struct {
	CCID                     uint64                 `gorm:"primary_key;auto_increment" json:"cc_id"`
	ContractNumber           string                 `gorm:"not null" json:"contract_number"`
	Alias                    string                 `gorm:"not null" json:"alias"`
	TelematicSubscriptionFee *bool                  `gorm:"not null" json:"telematic_subscription_fee"`
	PaperInvoice             *bool                  `gorm:"not null" json:"paper_invoice"`
	UseInvoiceAddress        *bool                  `gorm:"not null" json:"use_invoice_address"`
	ShowCreditLimit          *bool                  `gorm:"not null" json:"show_credit_limit"`
	InvoiceProductionLevel   uint64                 `gorm:"not null" json:"invoice_production_level"`
	BankVirtualAccount       string                 `gorm:"not null;size:30" json:"bank_virtual_account"`
	CreditLimit              float64                `gorm:"not null;" json:"credit_limit"`
	MembershipID             *int                   `gorm:"not null" json:"membership_id"`
	MCMSID                   int                    `gorm:"not null;" json:"mcms_id"`
	PaddedMCMSID             string                 `json:"padded_mcms_id"`
	GSAPCustomerMasterData   GSAPCustomerMasterData `json:"gsap_customer_master_data"`
	LatestPayerStatus        HistoricalPayerStatus  `json:"latest_payer_status"`
	Branch                   []ShortenedBranch      `json:"branch"`
	CreatedAt                time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                *time.Time             `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt                *time.Time             `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type ShortenedPayer struct {
	CCID                     uint64                          `gorm:"primary_key;auto_increment" json:"cc_id"`
	ContractNumber           string                          `gorm:"not null" json:"contract_number"`
	Alias                    string                          `gorm:"not null" json:"alias"`
	TelematicSubscriptionFee *bool                           `gorm:"not null" json:"telematic_subscription_fee"`
	PaperInvoice             *bool                           `gorm:"not null" json:"paper_invoice"`
	UseInvoiceAddress        *bool                           `gorm:"not null" json:"use_invoice_address"`
	ShowCreditLimit          *bool                           `gorm:"not null" json:"show_credit_limit"`
	InvoiceProductionLevel   uint64                          `gorm:"not null" json:"invoice_production_level"`
	BankVirtualAccount       string                          `gorm:"not null;size:30" json:"bank_virtual_account"`
	CreditLimit              float64                         `gorm:"not null;" json:"credit_limit"`
	MembershipID             *int                            `gorm:"not null" json:"membership_id"`
	MCMSID                   int                             `gorm:"not null;" json:"mcms_id"`
	PaddedMCMSID             string                          `json:"padded_mcms_id"`
	GSAPCustomerMasterData   ShortenedGSAPCustomerMasterData `gorm:"foreignkey:MCMSID; association_foreignkey:MCMSID" json:"gsap_customer_master_data"`
	LatestPayerStatus        HistoricalPayerStatus           `gorm:"foreignkey:CCID" json:"latest_payer_status"`
	CreatedAt                time.Time                       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                *time.Time                      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt                *time.Time                      `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

type HistoricalPayerStatus struct {
	ID            int         `json:"id"`
	CCID          uint64      `json:"cc_id"`
	PayerStatusID int         `json:"payer_status_id"`
	PayerStatus   PayerStatus `json:"payer_status"`
	CreatedAt     time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     *time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt     *time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

type PayerStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (payer *ShortenedPayer) Prepare() {
	payer.ContractNumber = html.EscapeString(strings.TrimSpace(payer.ContractNumber))
	payer.Alias = html.EscapeString(strings.TrimSpace(payer.Alias))
	payer.BankVirtualAccount = html.EscapeString(strings.TrimSpace(payer.BankVirtualAccount))
}

func (payer *ShortenedPayer) ValidateConfiguration() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

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

	if payer.ShowCreditLimit == nil {
		err = errors.New("Show credit limit field is required")
		errorMessages["required_show_credit_limit"] = err.Error()
	}

	if payer.InvoiceProductionLevel < 1 {
		err = errors.New("Invoice production level field is required")
		errorMessages["required_invoice_production_level"] = err.Error()
	}

	if payer.CreditLimit < 0 {
		err = errors.New("Credit limit field is required")
		errorMessages["required_credit_limit"] = err.Error()
	}

	return errorMessages
}

func (payer *ShortenedPayer) FindAllPayers(db *gorm.DB) (*[]ShortenedPayer, error) {
	var err error
	payers := []ShortenedPayer{}
	err = db.Debug().Model(&Payer{}).Unscoped().Order("created_at desc").Find(&payers).Error
	if err != nil {
		return &[]ShortenedPayer{}, err
	}

	if len(payers) > 0 {
		for i, _ := range payers {
			customerDataErr := db.Debug().Model(&Payer{}).Unscoped().Where("mcms_id = ?", payers[i].MCMSID).Order("mcms_id desc").Take(&payers[i].GSAPCustomerMasterData).Error
			if customerDataErr != nil {
				return &[]ShortenedPayer{}, err
			}

			latestStatusErr := db.Debug().Model(&HistoricalPayerStatus{}).Where("cc_id = ?", payers[i].CCID).Order("id asc").Find(&payers[i].LatestPayerStatus).Error
			if latestStatusErr != nil {
				return &[]ShortenedPayer{}, err
			}

			statusErr := db.Debug().Model(&PayerStatus{}).Where("id = ?", payers[i].LatestPayerStatus.PayerStatusID).Order("id desc").Take(&payers[i].LatestPayerStatus.PayerStatus).Error
			if statusErr != nil {
				return &[]ShortenedPayer{}, err
			}

			payers[i].PaddedMCMSID = fmt.Sprintf("%010v", strconv.Itoa(payers[i].MCMSID))
		}
	}

	return &payers, nil
}

func (payer *Payer) FindPayerByCCID(db *gorm.DB, CCID uint64) (*Payer, error) {
	var err error
	err = db.Debug().Model(&Payer{}).Unscoped().Where("cc_id = ?", CCID).Order("created_at desc").Take(&payer).Error
	if err != nil {
		return &Payer{}, err
	}

	customerDataErr := db.Debug().Model(&Payer{}).Unscoped().Where("mcms_id = ?", payer.MCMSID).Order("mcms_id desc").Take(&payer.GSAPCustomerMasterData).Error
	if customerDataErr != nil {
		return &Payer{}, err
	}

	latestStatusErr := db.Debug().Model(&HistoricalPayerStatus{}).Where("cc_id = ?", payer.CCID).Order("id asc").Find(&payer.LatestPayerStatus).Error
	if latestStatusErr != nil {
		return &Payer{}, err
	}

	statusErr := db.Debug().Model(&PayerStatus{}).Where("id = ?", payer.LatestPayerStatus.PayerStatusID).Order("id desc").Take(&payer.LatestPayerStatus.PayerStatus).Error
	if statusErr != nil {
		return &Payer{}, err
	}

	branchErr := db.Debug().Model(&ShortenedBranch{}).Where("cc_id = ?", payer.CCID).Order("sub_corporate_id desc").Find(&payer.Branch).Error
	if branchErr != nil {
		return &Payer{}, err
	}

	payer.PaddedMCMSID = fmt.Sprintf("%010v", strconv.Itoa(payer.MCMSID))

	return payer, nil
}

func (payer *ShortenedPayer) UpdatePayerConfiguration(db *gorm.DB) (*ShortenedPayer, error) {
	var err error
	dateTimeNow := time.Now()
	err = db.Debug().Model(&ShortenedPayer{}).Where("cc_id = ?", payer.CCID).Updates(
		ShortenedPayer{
			ContractNumber:           payer.ContractNumber,
			MembershipID:             payer.MembershipID,
			TelematicSubscriptionFee: payer.TelematicSubscriptionFee,
			Alias:                    payer.Alias,
			PaperInvoice:             payer.PaperInvoice,
			ShowCreditLimit:          payer.ShowCreditLimit,
			InvoiceProductionLevel:   payer.InvoiceProductionLevel,
			CreditLimit:              payer.CreditLimit,
			BankVirtualAccount:       payer.BankVirtualAccount,
			UpdatedAt:                &dateTimeNow,
		}).Error

	if err != nil {
		return &ShortenedPayer{}, err
	}

	customStatus := HistoricalPayerStatus{
		PayerStatusID: payer.LatestPayerStatus.PayerStatusID,
		CCID:          payer.CCID,
		CreatedAt:     dateTimeNow,
	}
	err = db.Debug().Model(&HistoricalPayerStatus{}).Create(&customStatus).Error
	if err != nil {
		return &ShortenedPayer{}, err
	}

	if *payer.PaperInvoice {
		paperInvoiceFee := ShortenedFee{}
		err = db.Debug().Model(&ShortenedFee{}).Where("name = ?", "Paper Invoice Fee").Order("id desc").Take(&paperInvoiceFee).Error
		if err != nil {
			return &ShortenedPayer{}, err
		}

		if paperInvoiceFee.DefaultValue > 0 {
			chargedPaperInvoiceFee := ChargedAutomatedFee{}
			fmt.Printf("%v", chargedPaperInvoiceFee.ID)
			err = db.Debug().Model(&ChargedAutomatedFee{}).Where("fee_id = ?", paperInvoiceFee.ID).Order("id desc").Take(&chargedPaperInvoiceFee).Error
			if err != nil {
				if err.Error() == "record not found" {
					chargePaperInvoiceFee := ChargedAutomatedFee{
						CCID:  payer.CCID,
						FeeID: paperInvoiceFee.ID,
						Value: paperInvoiceFee.DefaultValue,
					}

					err = db.Debug().Create(&chargePaperInvoiceFee).Error
					if err != nil {
						return &ShortenedPayer{}, err
					}
				} else {
					return &ShortenedPayer{}, err
				}
			}
		}
	}

	if *payer.TelematicSubscriptionFee {
		telematicSubscriptionFee := ShortenedFee{}
		err = db.Debug().Model(&ShortenedFee{}).Where("name = ?", "Telematic Subscription Fee").Order("id desc").Take(&telematicSubscriptionFee).Error
		if err != nil {
			return &ShortenedPayer{}, err
		}

		if telematicSubscriptionFee.DefaultValue > 0 {
			chargedTelematicSubscriptionFee := ChargedAutomatedFee{}
			err = db.Debug().Model(&ChargedAutomatedFee{}).Where("fee_id = ?", telematicSubscriptionFee.ID).Order("id desc").Take(&chargedTelematicSubscriptionFee).Error
			if err != nil {
				if err.Error() == "record not found" {
					chargeTelematicSubscriptionFee := ChargedAutomatedFee{
						CCID:  payer.CCID,
						FeeID: telematicSubscriptionFee.ID,
						Value: telematicSubscriptionFee.DefaultValue,
					}

					err = db.Debug().Create(&chargeTelematicSubscriptionFee).Error
					if err != nil {
						return &ShortenedPayer{}, err
					}
				} else {
					return &ShortenedPayer{}, err
				}
			}
		}
	}

	return payer, nil
}

func (Payer) TableName() string {
	return "corporate_client_relation"
}

func (ShortenedPayer) TableName() string {
	return "corporate_client_relation"
}

func (HistoricalPayerStatus) TableName() string {
	return "payer_status_relation"
}

func (PayerStatus) TableName() string {
	return "payer_status"
}
