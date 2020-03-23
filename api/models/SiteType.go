package models

import "github.com/jinzhu/gorm"

type SiteType struct {
	ID           int    `gorm:"primary_key;auto_increment" json:"id"`
	Name         string `gorm:"not null;size:50" json:"name"`
	Abbreviation string `gorm:"not null;size:5" json:"abbreviation"`
}

func (siteType *SiteType) FindAllSiteTypes(db *gorm.DB) (*[]SiteType, error) {
	var err error
	siteTypes := []SiteType{}
	err = db.Debug().Model(&SiteType{}).Order("id desc").Find(&siteTypes).Error
	if err != nil {
		return &[]SiteType{}, err
	}

	return &siteTypes, nil
}

func (siteType *SiteType) FindSiteTypeByID(db *gorm.DB, cityID uint64) (*SiteType, error) {
	var err error
	err = db.Debug().Model(&SiteType{}).Where("id = ?", cityID).Order("id desc").Take(&siteType).Error
	if err != nil {
		return &SiteType{}, err
	}

	return siteType, nil
}

func (SiteType) TableName() string {
	return "site_type"
}
