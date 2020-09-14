package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

//MemberCard => Member Card on SHELL FLeet
type MemberCard struct {
	CardID                  string     `gorm:"column:card_id" json:"card_id"`
	ExpDate                 string     `json:"exp_date"`
	CVV                     string     `json:"cvv"`
	BankCode                string     `json:"bank_code"`
	CountryCode             string     `json:"country_code"`
	Status                  string     `json:"status"`
	BlockReasonCode         *string    `json:"block_reason_code"`
	ActivationDateTime      *time.Time `gorm:"column:activation_datetime" json:"activation_datetime"`
	LastTransactionDateTime *time.Time `gorm:"column:last_transaction_datetime" json:"last_transaction_datetime"`
	PINOffset               string     `json:"pin_offset"`
	PINOffsetVersion        int        `json:"pin_offset_version"`
	PINGenerationDateTime   *time.Time `gorm:"column:pin_generation_datetime" json:"pin_generation_datetime"`
	PINChangeDateTime       *time.Time `gorm:"column:pin_change_datetime" json:"pin_change_datetime"`
	Batch                   int        `json:"batch"`
	FlagFleetID             *int       `json:"flag_fleet_id"`
	FlagOdometer            *int       `json:"flag_odometer"`
	CardGroupID             int        `gorm:"column:card_group_code" json:"card_group_code"`
	CardHolderTypeID        int        `gorm:"column:card_holder_type_code" json:"card_holder_type_code"`
	CardTypeID              int        `json:"card_type_id"`
	CardProfileID           int        `json:"card_profile_id"`
	TelematicDeviceID       string     `json:"telematic_device_id"`
	TelematicStartedAt      *time.Time `json:"telematic_started_at"`
	TelematicEndedAt        *time.Time `json:"telematic_ended_at"`
	CreatedAt               time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt               time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt               *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

//Prepare => Prepare Member Card
func (mc *MemberCard) Prepare() {
	mc.CreatedAt = time.Now()
}

//Validate => Validate Member Card
func (mc *MemberCard) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if mc.CardID == "" {
		err = errors.New("Card ID field is required")
		errorMessages["card_id"] = err.Error()
	}

	if mc.ExpDate == "" {
		err = errors.New("Expiry Date field is required")
		errorMessages["expiry_date"] = err.Error()
	}

	if mc.CVV == "" {
		err = errors.New("CVV field is required")
		errorMessages["cvv"] = err.Error()
	}

	if mc.BankCode == "" {
		err = errors.New("Bank Code field is required")
		errorMessages["bank_code"] = err.Error()
	}

	if mc.CountryCode == "" {
		err = errors.New("CVV field is required")
		errorMessages["country_code"] = err.Error()
	}

	if mc.Batch == 0 {
		err = errors.New("Batch field is required")
		errorMessages["batch"] = err.Error()
	}

	if mc.CardGroupID == 0 {
		err = errors.New("Card Group ID field is required")
		errorMessages["batch"] = err.Error()
	}

	if mc.CardHolderTypeID == 0 {
		err = errors.New("Card Holder Type ID field is required")
		errorMessages["batch"] = err.Error()
	}

	if mc.CardTypeID == 0 {
		err = errors.New("Card Type ID field is required")
		errorMessages["batch"] = err.Error()
	}

	if mc.CardProfileID == 0 {
		err = errors.New("Card Profile ID field is required")
		errorMessages["batch"] = err.Error()
	}

	return errorMessages
}

//FindMemberCards => Find Member Cards
func (mc *MemberCard) FindMemberCards(db *gorm.DB) (*[]MemberCard, error) {
	var err error
	mcs := []MemberCard{}
	err = db.Debug().Model(&MemberCard{}).
		Order("card_id, created_at desc").
		Find(&mcs).Error

	if err != nil {
		return &[]MemberCard{}, err
	}

	return &mcs, nil
}

//FindMemberCardByID => Find Member Card by ID
func (mc *MemberCard) FindMemberCardByID(db *gorm.DB, cardID string) (*MemberCard, error) {
	var err error
	err = db.Debug().Model(&AccountClass{}).Unscoped().
		Where("card_id = ?", cardID).
		Order("created_at desc").
		Take(&mc).Error

	if err != nil {
		return &MemberCard{}, err
	}

	return mc, nil
}

//CreateMemberCard => Create Member Card
func (mc *MemberCard) CreateMemberCard(db *gorm.DB) (*MemberCard, error) {
	var err error
	err = db.Debug().Model(&MemberCard{}).Create(&mc).Error
	if err != nil {
		return &MemberCard{}, err
	}

	_, err = mc.FindMemberCardByID(db, mc.CardID)
	if err != nil {
		return &MemberCard{}, err
	}

	return mc, nil
}

//TableName => Define Table
func (MemberCard) TableName() string {
	return "member_card"
}
