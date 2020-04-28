package models

import "github.com/jinzhu/gorm"

type Unit struct {
	ID   int    `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"not null;" json:"name"`
}

func (unit *Unit) FindAllUnits(db *gorm.DB) (*[]Unit, error) {
	var err error
	units := []Unit{}
	err = db.Debug().Model(&Unit{}).Limit(100).Order("id asc").Find(&units).Error
	if err != nil {
		return &[]Unit{}, err
	}

	return &units, nil
}

func (Unit) TableName() string {
	return "unit"
}
