package models

import "github.com/jinzhu/gorm"

type GSAPCustomerMasterData struct {
	MCMSID               int               `gorm:"primary_key;auto_increment" json:"id"`
	RecordType           string            `json:"record_type"`
	PayerNumber          string            `json:"payer_number"`
	BranchAccountNumber  string            `json:"branch_account_number"`
	AgentAccountNumber   string            `json:"agent_account_number"`
	SalesOrganization    string            `json:"sales_organization"`
	IndustryClassID      int               `json:"industry_class_id"`
	IndustryClass        GSAPIndustryClass `json:"industry_class"`
	BusinessTypeID       int               `json:"business_type_id"`
	BusinessType         GSAPBusinessType  `json:"business_type"`
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
	FaxNumber            string            `json:"fax_number"`
	CompanyRegNumber     string            `json:"company_reg_number"`
	VATRegNumber         string            `json:"vat_reg_number"`
	SalesRep             string            `json:"sales_rep"`
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
	AccountClass         string            `json:"account_class"`
	PaymentTerm          string            `json:"payment_term"`
	BillingSchedule      string            `json:"billing_schedule"`
	PaymentMethod        string            `json:"payment_method"`
	DistributionChannel  string            `json:"distribution_channel"`
	Division             string            `json:"division"`
	BillToEmailAddress   string            `json:"bill_to_email_address"`
}

type ShortenedGSAPCustomerMasterData struct {
	MCMSID             int    `json:"mcms_id"`
	PayerNumber        string `json:"payer_number"`
	AgentAccountNumber string `json:"agent_account_number"`
	ContactName_1      string `json:"contact_name_1"`
	ContactName_2      string `json:"contact_name_2"`
	ContactName_3      string `json:"contact_name_3"`
	ContactName_4      string `json:"contact_name_4"`
}

type GSAPBusinessType struct {
	ID   int    `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"not null;size:50" json:"name"`
}

type GSAPIndustryClass struct {
	ID   int    `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"not null;size:50" json:"name"`
}

func (gsapCustomerMasterDatum *GSAPCustomerMasterData) FindDataByMCMSID(db *gorm.DB, MCMSID uint64) (*GSAPCustomerMasterData, error) {
	var err error
	err = db.Debug().Model(&GSAPCustomerMasterData{}).Unscoped().Where("mcms_id = ?", MCMSID).Order("created_at desc").Take(&gsapCustomerMasterDatum).Error
	if err != nil {
		return &GSAPCustomerMasterData{}, err
	}

	if gsapCustomerMasterDatum.MCMSID != 0 {
		err := db.Debug().Model(&GSAPBusinessType{}).Unscoped().Where("id = ?", gsapCustomerMasterDatum.BusinessTypeID).Order("id desc").Take(&gsapCustomerMasterDatum.BusinessType).Error
		if err != nil {
			return &GSAPCustomerMasterData{}, err
		}

		err = db.Debug().Model(&GSAPIndustryClass{}).Unscoped().Where("id = ?", gsapCustomerMasterDatum.IndustryClassID).Order("id desc").Take(&gsapCustomerMasterDatum.IndustryClass).Error
		if err != nil {
			return &GSAPCustomerMasterData{}, err
		}
	}

	return gsapCustomerMasterDatum, nil
}

func (GSAPCustomerMasterData) TableName() string {
	return "gsap_customer_master_data"
}

func (ShortenedGSAPCustomerMasterData) TableName() string {
	return "gsap_customer_master_data"
}

func (GSAPIndustryClass) TableName() string {
	return "gsap_industry_class"
}

func (GSAPBusinessType) TableName() string {
	return "gsap_business_type"
}
