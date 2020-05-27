package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type GSAPCustomerMasterData struct {
	MCMSID               int               `gorm:"primary_key;auto_increment" json:"id"`
	RecordType           string            `json:"record_type"`
	PayerNumber          string            `json:"payer_number"`
	BranchAccountNumber  string            `json:"branch_account_number"`
	AgentAccountNumber   string            `json:"agent_account_number"`
	SalesOrganization    string            `json:"sales_organization"`
	IndustryClassID      string            `json:"industry_class_id"`
	IndustryClass        GSAPIndustryClass `gorm:"ForeignKey:IndustryClassID;AssociationForeignKey:Code" json:"industry_class"`
	BusinessTypeID       string            `json:"business_type_id"`
	BusinessType         GSAPBusinessType  `gorm:"ForeignKey:BusinessTypeID;AssociationForeignKey:Code" json:"business_type"`
	ContactName_1        string            `json:"contact_name_1"`
	ContactName_2        string            `json:"contact_name_2"`
	ContactName_3        string            `json:"contact_name_3"`
	ContactName_4        string            `json:"contact_name_4"`
	ContactAddress_1     string            `json:"contact_address_1"`
	ContactAddress_2     string            `json:"contact_address_2"`
	ContactAddress_3     string            `json:"contact_address_3"`
	ContactAddress_4     string            `json:"contact_address_4"`
	ContactPostalCode    string            `json:"contact_postal_code"`
	ContactSecondAddress string            `json:"contact_second_address"`
	ContactThirdAddress  string            `json:"contact_third_address"`
	WorkPhoneNumber      string            `json:"work_phone_number"`
	MobilePhoneNumber    string            `json:"mobile_phone_number"`
	FaxNumber            string            `json:"fax_number"`
	CompanyRegNumber     string            `json:"company_reg_number"`
	VATRegNumber         string            `json:"vat_reg_number"`
	SalesRepID           string            `json:"sales_rep_id"`
	SalesRep             GSAPSalesRep      `gorm:"ForeignKey:SalesRepID;AssociationForeignKey:Code" json:"sales_rep"`
	InvoiceName_1        string            `json:"invoice_name_1"`
	InvoiceName_2        string            `json:"invoice_name_2"`
	InvoiceName_3        string            `json:"invoice_name_3"`
	InvoiceName_4        string            `json:"invoice_name_4"`
	InvoiceAddress_1     string            `json:"invoice_address_1"`
	InvoiceAddress_2     string            `json:"invoice_address_2"`
	InvoiceAddress_3     string            `json:"invoice_address_3"`
	InvoiceAddress_4     string            `json:"invoice_address_4"`
	InvoicePostalCode    string            `json:"invoice_postal_code"`
	InvoiceSecondAddress string            `json:"invoice_second_address"`
	InvoiceThirdAddress  string            `json:"invoice_third_address"`
	BankAccountNumber    string            `json:"bank_account_number"`
	AccountClassID       string            `json:"account_class_id"`
	AccountClass         GSAPAccountClass  `gorm:"ForeignKey:AccountClassID;AssociationForeignKey:Code" json:"account_class"`
	PaymentTermID        string            `json:"payment_term_id"`
	PaymentTerm          GSAPPaymentTerm   `gorm:"ForeignKey:PaymentTermID;AssociationForeignKey:Code" json:"payment_term"`
	BillingSchedule      string            `json:"billing_schedule"`
	PaymentMethod        string            `json:"payment_method"`
	DistributionChannel  string            `json:"distribution_channel"`
	Division             string            `json:"division"`
	BillToEmailAddress   string            `json:"bill_to_email_address"`
	CreatedAt            time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt            time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt            *time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

type ShortenedGSAPCustomerMasterData struct {
	MCMSID              int        `gorm:"primary_key;auto_increment" json:"mcms_id"`
	PayerNumber         string     `json:"payer_number"`
	AgentAccountNumber  string     `json:"agent_account_number"`
	BranchAccountNumber string     `json:"branch_account_number"`
	ContactName_1       string     `json:"contact_name_1"`
	ContactName_2       string     `json:"contact_name_2"`
	ContactName_3       string     `json:"contact_name_3"`
	ContactName_4       string     `json:"contact_name_4"`
	AccountClass        string     `json:"account_class"`
	CreatedAt           time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt           time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt           *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

func (gsapCustomerMasterDatum *GSAPCustomerMasterData) FindDataByMCMSID(db *gorm.DB, MCMSID uint64) (*GSAPCustomerMasterData, error) {
	var err error
	err = db.Debug().Model(&GSAPCustomerMasterData{}).
		Unscoped().
		Where("mcms_id = ?", MCMSID).
		Order("created_at desc").
		Take(&gsapCustomerMasterDatum).Error

	if err != nil {
		return &GSAPCustomerMasterData{}, err
	}

	fmt.Println(gsapCustomerMasterDatum.MCMSID)

	if gsapCustomerMasterDatum.MCMSID != 0 {
		err := db.Debug().Model(&GSAPBusinessType{}).Unscoped().Where("code = ?", gsapCustomerMasterDatum.BusinessTypeID).Order("id desc").Take(&gsapCustomerMasterDatum.BusinessType).Error
		if err != nil {
			return &GSAPCustomerMasterData{}, err
		}

		err = db.Debug().Model(&GSAPIndustryClass{}).Unscoped().Where("code = ?", gsapCustomerMasterDatum.IndustryClassID).Order("id desc").Take(&gsapCustomerMasterDatum.IndustryClass).Error
		if err != nil {
			return &GSAPCustomerMasterData{}, err
		}
		fmt.Println("hehe")
	}

	return gsapCustomerMasterDatum, nil
}

func (GSAPCustomerMasterData) TableName() string {
	return "gsap_customer_master_data"
}

func (ShortenedGSAPCustomerMasterData) TableName() string {
	return "gsap_customer_master_data"
}
