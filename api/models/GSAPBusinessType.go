package models

import "time"

type GSAPBusinessType struct {
	ID                         uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Code                       string     `json:"code"`
	Name                       string     `json:"name"`
	LineOffBusinessDescription string     `json:"line_off_business_description"`
	CreatedAt                  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt                  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt                  *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

func (GSAPBusinessType) TableName() string {
	return "gsap_business_type"
}
