package models

import "github.com/jinzhu/gorm"

type Province struct {
	ID           int     `gorm:"primary_key;auto_increment" json:"id"`
	Name         string  `gorm:"not null;size:50" json:"name"`
	FuelTaxValue float32 `gorm:"not null;" json:"fuel_tax_value"`
}

func (province *Province) FindAllProvinces(db *gorm.DB) (*[]Province, error) {
	var err error
	provinces := []Province{}
	err = db.Debug().Model(&Province{}).Limit(100).Order("id desc").Find(&provinces).Error
	if err != nil {
		return &[]Province{}, err
	}

	return &provinces, nil
}

func (province *Province) FindProvinceByID(db *gorm.DB, provinceID uint64) (*Province, error) {
	var err error
	err = db.Debug().Model(&Province{}).Where("id = ?", provinceID).Order("id desc").Take(&province).Error
	if err != nil {
		return &Province{}, err
	}

	return province, nil
}

func (Province) TableName() string {
	return "province"
}
