package models

import "time"

type RebateProgram struct {
	ID                        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Name                      string     `gorm:"not null;size:100" json:"code"`
	RebateTypeID              uint64     `gorm:"not null;" json:"name"`
	RebateCalculationPeriodID uint64     `gorm:"not null;" json:"default_value"`
	RebateCalculationTypeID   uint64     `gorm:"not null;" json:"unit_id"`
	StartedAt                 *time.Time `json:"started_at"`
	EndedAt                   *time.Time `json:"ended_at"`
	CreatedAt                 time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                 time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt                 *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (RebatePeriod) TableName() string {
	return "rebate_period"
}
