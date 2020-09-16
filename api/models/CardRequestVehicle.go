package models

import (
	"errors"
	"html"
	"strings"
	"time"
)

//CardRequestVehicle => Struct to contain Card Vehicle Request
type CardRequestVehicle struct {
	Batch          int       `json:"batch"`
	ExpDate        string    `json:"exp_date"`
	CardTypeID     int       `json:"card_type_id"`
	CardTypePrefix string    `json:"card_type_prefix"`
	CardTypeSuffix string    `json:"card_type_suffix"`
	CardGroupID    int       `json:"card_group_id"`
	ResProfileID   int       `json:"res_profile_id"`
	SubCorporateID int       `json:"sub_corporate_id"`
	CCID           int       `json:"cc_id"`
	Vehicles       []Vehicle `json:"vehicle"`
	RequestedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"requested_at"`
}

//Prepare => Prepare string & datetime
func (crv *CardRequestVehicle) Prepare() {
	crv.ExpDate = html.EscapeString(strings.TrimSpace(crv.ExpDate))
	crv.RequestedAt = time.Now()
}

//Validate => Validate given request body
func (crv *CardRequestVehicle) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if crv.CardTypeID == 0 {
		err = errors.New("Card Type field is required")
		errorMessages["card_type"] = err.Error()
	}

	if crv.ExpDate == "" {
		err = errors.New("Expiry Date field is required")
		errorMessages["exp_date"] = err.Error()
	}

	if crv.CardGroupID == 0 {
		err = errors.New("Card Group ID field is required")
		errorMessages["card_group"] = err.Error()
	}

	if crv.SubCorporateID == 0 {
		err = errors.New("Sub Corporate ID field is required")
		errorMessages["sub_corporate"] = err.Error()
	}

	if crv.CCID == 0 {
		err = errors.New("CC ID field is required")
		errorMessages["cc_id"] = err.Error()
	}

	if len(crv.Vehicles) == 0 {
		err = errors.New("Vehicle field is required")
		errorMessages["card_count"] = err.Error()
	}

	return errorMessages
}
