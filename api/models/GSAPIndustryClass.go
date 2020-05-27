package models

import "time"

type GSAPIndustryClass struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

func (GSAPIndustryClass) TableName() string {
	return "gsap_industry_class"
}
