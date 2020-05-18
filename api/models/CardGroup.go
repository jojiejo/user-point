package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type CardGroup struct {
	CardGroupCode  int        `gorm:"primary_key;auto_increment" json:"card_group_code"`
	CardGroupName  string     `gorm:"not null" json:"card_group_name"`
	SubCorporateID int        `gorm:"not null" json:"sub_corporate_id"`
	ResProfileID   int        `gorm:"not null" json:"res_profile_id"`
	CreatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (cg *CardGroup) Prepare() {
	cg.CardGroupName = html.EscapeString(strings.TrimSpace(cg.CardGroupName))
	cg.CreatedAt = time.Now()
}

func (cg *CardGroup) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if cg.CardGroupName == "" {
		err = errors.New("Card group name field is required")
		errorMessages["card_group_name"] = err.Error()
	}

	if cg.SubCorporateID < 1 {
		err = errors.New("Sub account field is required")
		errorMessages["sub_account"] = err.Error()
	}

	if cg.ResProfileID < 1 {
		err = errors.New("Restriction profile field is required")
		errorMessages["res_profile"] = err.Error()
	}

	return errorMessages
}

func (cg *CardGroup) FindAllCardGroupsByBranchID(db *gorm.DB, branchID uint64) (*[]CardGroup, error) {
	var err error
	cgs := []CardGroup{}
	err = db.Debug().Model(&CardGroup{}).Unscoped().Where("sub_corporate_id = ?", branchID).Order("card_group_code, created_at desc").Find(&cgs).Error
	if err != nil {
		return &[]CardGroup{}, err
	}

	return &cgs, nil
}

func (cg *CardGroup) FindCardGroupByID(db *gorm.DB, cgID uint64) (*CardGroup, error) {
	var err error
	err = db.Debug().Model(&CardGroup{}).Unscoped().Where("card_group_code = ?", cgID).Order("card_group_code, created_at desc").Take(&cg).Error
	if err != nil {
		return &CardGroup{}, err
	}

	return cg, nil
}

func (cg *CardGroup) CreateCardGroup(db *gorm.DB) (*CardGroup, error) {
	var err error
	err = db.Debug().Model(&Site{}).Create(&cg).Error
	if err != nil {
		return &CardGroup{}, err
	}

	return cg, nil
}

func (cg *CardGroup) UpdateCardGroup(db *gorm.DB) (*CardGroup, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&CardGroup{}).Where("card_group_code = ?", cg.CardGroupCode).Updates(
		CardGroup{
			CardGroupName:  cg.CardGroupName,
			SubCorporateID: cg.SubCorporateID,
			ResProfileID:   cg.ResProfileID,
			UpdatedAt:      &dateTimeNow,
		}).Error

	if err != nil {
		return &CardGroup{}, err
	}

	return cg, nil
}

func (cg *CardGroup) DeleteCardGroup(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&CardGroup{}).Where("card_group_code = ?", cg.CardGroupCode).Delete(&CardGroup{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (CardGroup) TableName() string {
	return "mstCardGroup"
}
