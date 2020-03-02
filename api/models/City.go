package models

import "github.com/jinzhu/gorm"

type City struct {
	ID          int    `gorm:"primary_key;auto_increment" json:"id"`
	Name        string `gorm:"not null;size:50" json:"name"`
	PhonePrefix string `gorm:"not null;size:5" json:"phone_prefix"`
	ProvinceID  int    `gorm:"not null" json:"province_id"`
}

func (city *City) FindAllCities(db *gorm.DB) (*[]City, error) {
	var err error
	cities := []City{}
	err = db.Debug().Model(&City{}).Limit(100).Order("id desc").Find(&cities).Error
	if err != nil {
		return &[]City{}, err
	}

	return &cities, nil
}

func (city *City) FindCityByID(db *gorm.DB, cityID uint64) (*City, error) {
	var err error
	err = db.Debug().Model(&City{}).Where("id = ?", cityID).Order("id desc").Take(&city).Error
	if err != nil {
		return &City{}, err
	}

	return city, nil
}

func (city *City) FindCityByProvinceID(db *gorm.DB, provinceID uint64) (*[]City, error) {
	var err error
	cities := []City{}
	err = db.Debug().Model(&City{}).Where("province_id = ?", provinceID).Order("id desc").Find(&cities).Error
	if err != nil {
		return &[]City{}, err
	}

	return &cities, nil
}

func (City) TableName() string {
	return "city"
}
