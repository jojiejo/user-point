package models

import (
	"errors"
	"html"
	"strings"
	"time"
)

//CardBearerRequest => Struct to contain Card Bearer Request
type CardRequestBearer struct {
	Batch          uint64    `json:"batch"`
	ExpDate        string    `json:"exp_date"`
	CardTypeID     uint64    `json:"card_type_id"`
	CardTypePrefix string    `json:"card_type_prefix"`
	CardTypeSuffix string    `json:"card_type_suffix"`
	CardGroupID    uint64    `json:"card_group_id"`
	ResProfileID   int       `json:"res_profile_id"`
	SubCorporateID uint64    `json:"sub_corporate_id"`
	CCID           uint64    `json:"cc_id"`
	CardCount      uint64    `json:"card_count"`
	RequestedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"requested_at"`
}

//Prepare => Prepare string & datetime
func (cbr *CardRequestBearer) Prepare() {
	cbr.ExpDate = html.EscapeString(strings.TrimSpace(cbr.ExpDate))
	cbr.RequestedAt = time.Now()
}

//Validate => Validate given request body
func (cbr *CardRequestBearer) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if cbr.CardTypeID == 0 {
		err = errors.New("Card Type field is required")
		errorMessages["card_type"] = err.Error()
	}

	if cbr.ExpDate == "" {
		err = errors.New("Expiry Date field is required")
		errorMessages["exp_date"] = err.Error()
	}

	if cbr.CardGroupID == 0 {
		err = errors.New("Card Group ID field is required")
		errorMessages["card_group"] = err.Error()
	}

	if cbr.SubCorporateID == 0 {
		err = errors.New("Sub Corporate ID field is required")
		errorMessages["sub_corporate"] = err.Error()
	}

	if cbr.CCID == 0 {
		err = errors.New("CC ID field is required")
		errorMessages["cc_id"] = err.Error()
	}

	if cbr.CardCount == 0 {
		err = errors.New("Card Count field is required")
		errorMessages["card_count"] = err.Error()
	}

	return errorMessages
}

//GenerateBearerCard => Generate Bearer Card
/*func (cbr *CardRequestBearer) GenerateBearerCard() map[string]string {

}
*/
