package models

import "time"

type GSAPAccountClass struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
}

func (GSAPAccountClass) TableName() string {
	return "gsap_account_class"
}
