package models

import (
	"errors"
	"html"
	"strings"
	"time"
)

//CardRequestDriver => Struct to contain Card Driver Request
type CardRequestDriver struct {
	Batch          int       `json:"batch"`
	ExpDate        string    `json:"exp_date"`
	CardTypeID     int       `json:"card_type_id"`
	CardTypePrefix string    `json:"card_type_prefix"`
	CardTypeSuffix string    `json:"card_type_suffix"`
	CardGroupID    int       `json:"card_group_id"`
	ResProfileID   int       `json:"res_profile_id"`
	SubCorporateID int       `json:"sub_corporate_id"`
	CCID           int       `json:"cc_id"`
	Drivers        []Driver  `json:"drivers"`
	RequestedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"requested_at"`
}

//Prepare => Prepare string & datetime
func (crd *CardRequestDriver) Prepare() {
	crd.ExpDate = html.EscapeString(strings.TrimSpace(crd.ExpDate))
	crd.RequestedAt = time.Now()
}

//Validate => Validate given request body
func (crd *CardRequestDriver) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if crd.CardTypeID == 0 {
		err = errors.New("Card Type field is required")
		errorMessages["card_type"] = err.Error()
	}

	if crd.ExpDate == "" {
		err = errors.New("Expiry Date field is required")
		errorMessages["exp_date"] = err.Error()
	}

	if crd.CardGroupID == 0 {
		err = errors.New("Card Group ID field is required")
		errorMessages["card_group"] = err.Error()
	}

	if crd.SubCorporateID == 0 {
		err = errors.New("Sub Corporate ID field is required")
		errorMessages["sub_corporate"] = err.Error()
	}

	if crd.CCID == 0 {
		err = errors.New("CC ID field is required")
		errorMessages["cc_id"] = err.Error()
	}

	if len(crd.Drivers) == 0 {
		err = errors.New("Driver field is required")
		errorMessages["card_count"] = err.Error()
	}

	return errorMessages
}
