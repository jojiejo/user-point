package models

type RebateProgramTier struct {
	RebateProgramID uint64   `gorm:"not null" json:"-"`
	Sequence        uint64   `gorm:"not null;size:100" json:"sequence"`
	BottomLimit     float64  `gorm:"not null;" json:"bottom_limit"`
	TopLimit        *float64 `gorm:"not null;" json:"top_limit"`
	Value           float64  `gorm:"not null;" json:"value"`
}

func (RebateProgramTier) TableName() string {
	return "rebate_program_tier"
}
