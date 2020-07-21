package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Membership struct {
	ID        uint64     `gorm:"primary_key;auto_increment;column:membership_code" json:"id"`
	Name      string     `gorm:"not null;size:100;column:membership_name" json:"name"`
	StartedAt *time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (m *Membership) Prepare() {
	m.Name = html.EscapeString(strings.TrimSpace(m.Name))
	m.CreatedAt = time.Now()
}

func (m *Membership) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if m.Name == "" {
		err = errors.New("Name field is required")
		errorMessages["name"] = err.Error()
	}

	if m.StartedAt == nil {
		err = errors.New("Started at field is required")
		errorMessages["started_at"] = err.Error()
	}

	return errorMessages
}

func (m *Membership) FindMemberships(db *gorm.DB) (*[]Membership, error) {
	var err error
	memberships := []Membership{}
	err = db.Debug().Model(&Membership{}).
		Order("membership_code, created_at desc").
		Find(&memberships).Error

	if err != nil {
		return &[]Membership{}, err
	}

	return &memberships, nil
}

func (m *Membership) FindMembershipByID(db *gorm.DB, membershipID uint64) (*Membership, error) {
	var err error
	err = db.Debug().Model(&Membership{}).Unscoped().
		Where("membership_code = ?", membershipID).
		Order("created_at desc").
		Take(&m).Error

	if err != nil {
		return &Membership{}, err
	}

	return m, nil
}

func (m *Membership) CreateMembership(db *gorm.DB) (*Membership, error) {
	var err error
	err = db.Debug().Model(&Membership{}).Create(&m).Error
	if err != nil {
		return &Membership{}, err
	}

	//Select created fee
	_, err = m.FindMembershipByID(db, m.ID)
	if err != nil {
		return &Membership{}, err
	}

	return m, nil
}

func (m *Membership) UpdateMembership(db *gorm.DB) (*Membership, error) {
	var err error
	dateTimeNow := time.Now()

	//Update the data
	err = db.Debug().Model(&m).Updates(
		map[string]interface{}{
			"membership_name": m.Name,
			"started_at":      m.StartedAt,
			"ended_at":        m.EndedAt,
			"updated_at":      dateTimeNow,
		}).Error

	if err != nil {
		return &Membership{}, err
	}

	//Select updated membership
	_, err = m.FindMembershipByID(db, m.ID)
	if err != nil {
		return &Membership{}, err
	}

	return m, nil
}

func (Membership) TableName() string {
	return "mstMembership"
}
