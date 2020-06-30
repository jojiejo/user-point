package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type CardTelematicDevice struct {
	CardID             string     `gorm:"primary_key;" json:"card_id"`
	TelematicDeviceID  string     `json:"telematic_device_id"`
	TelematicStartedAt *time.Time `json:"telematic_started_at"`
	TelematicEndedAt   *time.Time `json:"telematic_ended_at"`
}

func (ctd *CardTelematicDevice) Prepare() {
	ctd.TelematicDeviceID = html.EscapeString(strings.TrimSpace(ctd.TelematicDeviceID))
}

func (ctd *CardTelematicDevice) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if ctd.TelematicDeviceID == "" {
		err = errors.New("Telematic device ID field is required")
		errorMessages["telematic_device_id"] = err.Error()
	}

	return errorMessages
}

func (ctd *CardTelematicDevice) FindTelematicDeviceByCardID(db *gorm.DB, cardID string) (*CardTelematicDevice, error) {
	var err error
	err = db.Debug().Model(&CardTelematicDevice{}).Unscoped().
		Where("card_id = ?", cardID).
		Order("card_id, created_at desc").
		Find(&ctd).Error

	if err != nil {
		return &CardTelematicDevice{}, err
	}

	return ctd, nil
}

func (ctd *CardTelematicDevice) UpdateTelematicDevice(db *gorm.DB) (*CardTelematicDevice, error) {
	var err error
	dateTimeNow := time.Now()

	err = db.Debug().Model(&ctd).Updates(
		map[string]interface{}{
			"telematic_device_id":  ctd.TelematicDeviceID,
			"telematic_started_at": ctd.TelematicStartedAt,
			"telematic_ended_at":   ctd.TelematicEndedAt,
			"updated_at":           dateTimeNow,
		}).Error

	if err != nil {
		return &CardTelematicDevice{}, err
	}

	//Select created fee
	_, err = ctd.FindTelematicDeviceByCardID(db, ctd.CardID)
	if err != nil {
		return &CardTelematicDevice{}, err
	}

	return ctd, nil
}

func (CardTelematicDevice) TableName() string {
	return "mstMemberCard"
}
