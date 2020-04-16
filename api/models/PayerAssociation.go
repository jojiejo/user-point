package models

import "github.com/jinzhu/gorm"

type PayerAssociation struct {
	Number string `gorm:"column:agent_account_number;not null;" json:"payer_association_number"`
}

func (pa *PayerAssociation) FindPayerAssociations(db *gorm.DB) (*[]PayerAssociation, error) {
	var err error
	pas := []PayerAssociation{}
	err = db.Debug().Model(&PayerAssociation{}).Unscoped().
		Select("DISTINCT(agent_account_number)").
		Order("agent_account_number desc").
		Find(&pas).Error

	if err != nil {
		return &[]PayerAssociation{}, err
	}

	return &pas, nil
}

func (payer *ShortenedPayer) FindPayerByPayerAssociationNumber(db *gorm.DB, payerAssociationNumber string) (*[]ShortenedPayer, error) {
	var err error

	payers := []ShortenedPayer{}
	if payerAssociationNumber == "0" {
		err = db.Debug().Model(&ShortenedPayer{}).Unscoped().
			Preload("GSAPCustomerMasterData").
			Preload("LatestPayerStatus").
			Preload("LatestPayerStatus.PayerStatus").
			Joins("JOIN gsap_customer_master_data ON Corporate_Client_Relation.mcms_id = gsap_customer_master_data.mcms_id").
			Where("gsap_customer_master_data.agent_account_number IS NULL").
			Order("created_at desc").Find(&payers).Error

		if err != nil {
			return &[]ShortenedPayer{}, err
		}
	} else {
		err = db.Debug().Model(&ShortenedPayer{}).Unscoped().
			Preload("GSAPCustomerMasterData").
			Preload("LatestPayerStatus").
			Preload("LatestPayerStatus.PayerStatus").
			Where("agent_account_number = ?", payerAssociationNumber).
			Joins("JOIN gsap_customer_master_data ON Corporate_Client_Relation.mcms_id = gsap_customer_master_data.mcms_id").
			Where("gsap_customer_master_data.agent_account_number = ?", payerAssociationNumber).
			Order("created_at desc").Find(&payers).Error

		if err != nil {
			return &[]ShortenedPayer{}, err
		}
	}

	return &payers, nil
}

func (PayerAssociation) TableName() string {
	return "gsap_customer_master_data"
}
