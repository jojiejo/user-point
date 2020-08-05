package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type AccountClass struct {
	ID        uint64     `gorm:"primary_key;auto_increment;" json:"id"`
	Code      string     `gorm:"not null;" json:"code"`
	Name      string     `gorm:"not null;size:100;" json:"name"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (ac *AccountClass) Prepare() {
	ac.Code = html.EscapeString(strings.TrimSpace(ac.Code))
	ac.Name = html.EscapeString(strings.TrimSpace(ac.Name))
	ac.CreatedAt = time.Now()
}

func (ac *AccountClass) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if ac.Code == "" {
		err = errors.New("Code field is required")
		errorMessages["code"] = err.Error()
	}

	if ac.Name == "" {
		err = errors.New("Name field is required")
		errorMessages["name"] = err.Error()
	}

	return errorMessages
}

func (ac *AccountClass) FindAccountClasses(db *gorm.DB) (*[]AccountClass, error) {
	var err error
	acs := []AccountClass{}
	err = db.Debug().Model(&AccountClass{}).
		Order("id, created_at desc").
		Find(&acs).Error

	if err != nil {
		return &[]AccountClass{}, err
	}

	return &acs, nil
}

func (ac *AccountClass) FindAccountClassByID(db *gorm.DB, srID uint64) (*AccountClass, error) {
	var err error
	err = db.Debug().Model(&AccountClass{}).Unscoped().
		Where("id = ?", srID).
		Order("created_at desc").
		Take(&ac).Error

	if err != nil {
		return &AccountClass{}, err
	}

	return ac, nil
}

func (ac *AccountClass) CreateAccountClass(db *gorm.DB) (*AccountClass, error) {
	var err error
	err = db.Debug().Model(&AccountClass{}).Create(&ac).Error
	if err != nil {
		return &AccountClass{}, err
	}

	//Select created fee
	_, err = ac.FindAccountClassByID(db, ac.ID)
	if err != nil {
		return &AccountClass{}, err
	}

	return ac, nil
}

func (ac *AccountClass) UpdateAccountClass(db *gorm.DB) (*AccountClass, error) {
	var err error
	dateTimeNow := time.Now()

	//Update the data
	err = db.Debug().Model(&ac).Updates(
		map[string]interface{}{
			"code":       ac.Code,
			"name":       ac.Name,
			"updated_at": dateTimeNow,
		}).Error

	if err != nil {
		return &AccountClass{}, err
	}

	//Select updated sales rep
	_, err = ac.FindAccountClassByID(db, ac.ID)
	if err != nil {
		return &AccountClass{}, err
	}

	return ac, nil
}

func (AccountClass) TableName() string {
	return "gsap_account_class"
}
