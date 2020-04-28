package models

import "github.com/jinzhu/gorm"

type CardStatus struct {
	ID   int    `gorm:"primary_key:true;auto_increment" json:"id"`
	Name string `gorm:"not null;" json:"name"`
}

func (CardStatus) TableName() string {
	return "card_status"
}

func (cs *CardStatus) FindAllCardStatus(db *gorm.DB) (*[]CardStatus, error) {
	var err error
	css := []CardStatus{}
	err = db.Debug().Model(&CardStatus{}).Limit(100).Order("id asc").Find(&css).Error
	if err != nil {
		return &[]CardStatus{}, err
	}

	return &css, nil
}
