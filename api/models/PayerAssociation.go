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

func (payer *ShortenedPayer) FindPayerByPayerAssociationNumber(db *gorm.DB, payerAssociationNumber string) (*ShortenedPayer, error) {
	var err error
	err = db.Debug().Model(&ShortenedPayer{}).Unscoped().
		Preload("GSAPCustomerMasterData").
		Preload("LatestPayerStatus").
		Preload("LatestPayerStatus.PayerStatus").
		Order("created_at desc").Take(&payer).Error

	if err != nil {
		return &ShortenedPayer{}, err
	}

	return payer, nil
}

func (PayerAssociation) TableName() string {
	return "gsap_customer_master_data"
}
