package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type PayerContactPerson struct {
	ID        uint64     `gorm:"primary_key;auto_increment;" json:"id"`
	CCID      uint64     `json:"cc_id"`
	Title     string     `json:"title"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Telephone string     `json:"telephone"`
	StartedAt *time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (pcp *PayerContactPerson) Prepare() {
	pcp.Name = html.EscapeString(strings.TrimSpace(pcp.Name))
	pcp.Email = html.EscapeString(strings.TrimSpace(pcp.Email))
	pcp.Telephone = html.EscapeString(strings.TrimSpace(pcp.Telephone))
	pcp.CreatedAt = time.Now()
}

func (pcp *PayerContactPerson) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if pcp.Name == "" {
		err = errors.New("Name field is required")
		errorMessages["name"] = err.Error()
	}

	if pcp.Email == "" {
		err = errors.New("Email field is required")
		errorMessages["email"] = err.Error()
	}

	if pcp.Telephone == "" {
		err = errors.New("Telephone field is required")
		errorMessages["telephone"] = err.Error()
	}

	if pcp.StartedAt == nil {
		err = errors.New("Started at field is required")
		errorMessages["material"] = err.Error()
	}

	return errorMessages
}

func (pcp *PayerContactPerson) FindPayerContactPersons(db *gorm.DB, ccID uint64) (*[]PayerContactPerson, error) {
	var err error
	pcps := []PayerContactPerson{}
	err = db.Debug().Model(&PayerContactPerson{}).
		Where("cc_id = ?", ccID).
		Order("id, created_at desc").
		Find(&pcps).Error

	if err != nil {
		return &[]PayerContactPerson{}, err
	}

	return &pcps, nil
}

func (pcp *PayerContactPerson) FindPayerContactPersonByID(db *gorm.DB, pcpID uint64) (*PayerContactPerson, error) {
	var err error
	err = db.Debug().Model(&PayerContactPerson{}).Unscoped().
		Where("id = ?", pcpID).
		Order("created_at desc").
		Take(&pcp).Error

	if err != nil {
		return &PayerContactPerson{}, err
	}

	return pcp, nil
}

func (pcp *PayerContactPerson) CreatePayerContactPerson(db *gorm.DB) (*PayerContactPerson, error) {
	var err error
	err = db.Debug().Model(&PayerContactPerson{}).Create(&pcp).Error
	if err != nil {
		return &PayerContactPerson{}, err
	}

	_, err = pcp.FindPayerContactPersonByID(db, pcp.ID)
	if err != nil {
		return &PayerContactPerson{}, err
	}

	return pcp, nil
}

func (pcp *PayerContactPerson) UpdatePayerContactPerson(db *gorm.DB) (*PayerContactPerson, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&pcp).Updates(
		map[string]interface{}{
			"title":      pcp.Title,
			"name":       pcp.Name,
			"email":      pcp.Email,
			"telephone":  pcp.Telephone,
			"started_at": pcp.StartedAt,
			"ended_at":   pcp.EndedAt,
			"updated_at": dateTimeNow,
		}).Error

	if err != nil {
		return &PayerContactPerson{}, err
	}

	_, err = pcp.FindPayerContactPersonByID(db, pcp.ID)
	if err != nil {
		return &PayerContactPerson{}, err
	}

	return pcp, nil
}

func (PayerContactPerson) TableName() string {
	return "payer_contact_person"
}
