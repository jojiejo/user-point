package models

type Unit struct {
	ID   int `gorm:"primary_key;auto_increment" json:"id"`
	Name int `gorm:"not null;" json:"name"`
}

func (Unit) TableName() string {
	return "unit"
}
