package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Retailer struct {
	ID                    int    `gorm:"primary_key;auto_increment" json:"id"`
	OriginalID            int    `json:"original_id"`
	SoldToNumber          string `gorm:"not null;size:10" json:"sold_to_number"`
	SoldToName            string `gorm:"not null; size:60" json:"sold_to_name"`
	Address_1             string `gorm:"not null;size:30" json:"address_1"`
	Address_2             string `gorm:"not null;size:30" json:"address_2"`
	Address_3             string `gorm:"not null;size:30" json:"address_3"`
	CityID                int    `gorm:"not null" json:"city_id"`
	City                  City   `json:"city"`
	RetailerPaymentTermID int    `gorm:"not null;" json:"retailer_payment_term_id"`
	//RetailerPaymentTerm          RetailerPaymentTerm        `json:"retailer_payment_term"`
	RetailerReimbursementCycleID int `gorm:"not null;" json:"retailer_reimbursement_cycle_id"`
	//RetailerReimbursementCycle   RetailerReimbursementCycle `json:"retailer_reimbursement_cycle"`
	ZipCode       string     `gorm:"not null;size:5" json:"zip_code"`
	Phone         string     `gorm:"not null;size:15" json:"phone"`
	Email         string     `gorm:"not null" json:"email"`
	CreatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP"  json:"created_at"`
	UpdatedAt     *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
	ReactivatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"reactivated_at"`
}

func (retailer *Retailer) Prepare() {
	retailer.SoldToNumber = html.EscapeString(strings.TrimSpace(retailer.SoldToNumber))
	retailer.SoldToName = html.EscapeString(strings.TrimSpace(retailer.SoldToName))
	retailer.Address_1 = html.EscapeString(strings.TrimSpace(retailer.Address_1))
	retailer.Address_2 = html.EscapeString(strings.TrimSpace(retailer.Address_2))
	retailer.Address_3 = html.EscapeString(strings.TrimSpace(retailer.Address_3))
	retailer.ZipCode = html.EscapeString(strings.TrimSpace(retailer.ZipCode))
	retailer.Phone = html.EscapeString(strings.TrimSpace(retailer.Phone))
	retailer.Email = html.EscapeString(strings.TrimSpace(retailer.Email))
	retailer.CreatedAt = time.Now()
}

func (retailer *Retailer) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if retailer.SoldToNumber == "" {
		err = errors.New("Sold to Number field is required")
		errorMessages["required_sold_to_number"] = err.Error()
	}

	if retailer.SoldToName == "" {
		err = errors.New("Sold to Name field is required")
		errorMessages["required_sold_to_name"] = err.Error()
	}

	if retailer.Address_1 == "" {
		err = errors.New("Address 1 field is required")
		errorMessages["required_address_1"] = err.Error()
	}

	if retailer.Address_2 == "" {
		err = errors.New("Address 2 field is required")
		errorMessages["required_address_2"] = err.Error()
	}

	if retailer.Address_3 == "" {
		err = errors.New("Address 3 field is required")
		errorMessages["required_address_3"] = err.Error()
	}

	if retailer.CityID < 1 {
		err = errors.New("City field is required")
		errorMessages["required_city"] = err.Error()
	}

	if retailer.RetailerPaymentTermID < 1 {
		err = errors.New("Retailer payment term field is required")
		errorMessages["required_retailer_payment_term"] = err.Error()
	}

	if retailer.RetailerReimbursementCycleID < 1 {
		err = errors.New("Retailer reimbursement cycle field is required")
		errorMessages["required_retailer_reimbursement_cycle"] = err.Error()
	}

	if retailer.ZipCode == "" {
		err = errors.New("ZIP Code field is required")
		errorMessages["required_zip_code"] = err.Error()
	}

	/*if retailer.Phone == "" {
		err = errors.New("Phone field is required")
		errorMessages["required_phone"] = err.Error()
	}*/

	if retailer.Email == "" {
		err = errors.New("Email field is required")
		errorMessages["required_email"] = err.Error()
	}

	return errorMessages
}

func (retailer *Retailer) FindAllRetailers(db *gorm.DB) (*[]Retailer, error) {
	var err error
	retailers := []Retailer{}
	err = db.Debug().Model(&Retailer{}).Unscoped().Order("created_at desc").Find(&retailers).Error
	if err != nil {
		return &[]Retailer{}, err
	}

	if len(retailers) > 0 {
		for i, _ := range retailers {
			err := db.Debug().Model(&City{}).Unscoped().Where("id = ?", retailers[i].CityID).Order("id desc").Take(&retailers[i].City).Error
			if err != nil {
				return &[]Retailer{}, err
			}
		}
	}

	return &retailers, nil
}

func (retailer *Retailer) FindAllLatestRetailers(db *gorm.DB) (*[]Retailer, error) {
	var err error
	retailers := []Retailer{}
	err = db.Debug().Raw("EXEC spAPI_Retailer_GetLatest").Scan(&retailers).Error
	if err != nil {
		return &[]Retailer{}, err
	}

	if len(retailers) > 0 {
		for i, _ := range retailers {
			err := db.Debug().Model(&City{}).Unscoped().Where("id = ?", retailers[i].CityID).Order("id desc").Take(&retailers[i].City).Error
			if err != nil {
				return &[]Retailer{}, err
			}

			/*if retailers[i].RetailerPaymentTermID != 0 {
				err := db.Debug().Model(&RetailerPaymentTerm{}).Unscoped().Where("id = ?", retailers[i].RetailerPaymentTermID).Order("id desc").Take(&retailers[i].RetailerPaymentTerm).Error
				if err != nil {
					return &[]Retailer{}, err
				}
			}

			if retailers[i].RetailerReimbursementCycleID != 0 {
				err := db.Debug().Model(&RetailerReimbursementCycle{}).Unscoped().Where("id = ?", retailers[i].RetailerReimbursementCycleID).Order("id desc").Take(&retailers[i].RetailerReimbursementCycle).Error
				if err != nil {
					return &[]Retailer{}, err
				}
			}*/
		}
	}

	return &retailers, nil
}

func (retailer *Retailer) FindRetailerByID(db *gorm.DB, retailerID uint64) (*Retailer, error) {
	var err error
	err = db.Debug().Model(&Retailer{}).Unscoped().Where("id = ?", retailerID).Order("created_at desc").Take(&retailer).Error
	if err != nil {
		return &Retailer{}, err
	}

	if retailer.ID != 0 {
		err := db.Debug().Model(&City{}).Unscoped().Where("id = ?", retailer.CityID).Order("id desc").Take(&retailer.City).Error
		if err != nil {
			return &Retailer{}, err
		}
	}

	return retailer, nil
}

func (retailer *Retailer) FindRetailerHistoryByID(db *gorm.DB, originalRetailerID uint64) (*[]Retailer, error) {
	var err error
	var retailers = []Retailer{}
	err = db.Debug().Model(&Retailer{}).Unscoped().Where("original_id = ?", originalRetailerID).Order("created_at desc").Find(&retailers).Error
	if err != nil {
		return &[]Retailer{}, err
	}

	if len(retailers) > 0 {
		for i, _ := range retailers {
			err := db.Debug().Model(&City{}).Unscoped().Where("id = ?", retailers[i].CityID).Order("id desc").Take(&retailers[i].City).Error
			if err != nil {
				return &[]Retailer{}, err
			}
		}
	}

	return &retailers, nil
}

func (retailer *Retailer) CreateRetailer(db *gorm.DB) (*Retailer, error) {
	var err error
	tx := db.Begin()
	err = db.Debug().Model(&Retailer{}).Create(&retailer).Error
	if err != nil {
		tx.Rollback()
		return &Retailer{}, err
	}

	err = db.Debug().Model(&Retailer{}).Where("id = ?", retailer.ID).Updates(
		Retailer{
			OriginalID: retailer.ID,
		}).Error
	if err != nil {
		tx.Rollback()
		return &Retailer{}, err
	}

	tx.Commit()
	return retailer, nil
}

func (retailer *Retailer) UpdateRetailer(db *gorm.DB) (*Retailer, error) {
	var err error
	dateTimeNow := time.Now()
	err = db.Debug().Model(&Retailer{}).Where("id = ?", retailer.ID).Updates(
		Retailer{
			SoldToName:                   retailer.SoldToName,
			Address_1:                    retailer.Address_1,
			Address_2:                    retailer.Address_2,
			Address_3:                    retailer.Address_3,
			CityID:                       retailer.CityID,
			ZipCode:                      retailer.ZipCode,
			Phone:                        retailer.Phone,
			Email:                        retailer.Email,
			RetailerPaymentTermID:        retailer.RetailerPaymentTermID,
			RetailerReimbursementCycleID: retailer.RetailerReimbursementCycleID,
			UpdatedAt:                    &dateTimeNow,
		}).Error

	if err != nil {
		return &Retailer{}, err
	}

	return retailer, nil
}

func (retailer *Retailer) DeactivateRetailer(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Retailer{}).Where("id = ?", retailer.ID).Delete(&Retailer{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (retailer *Retailer) ReactivateRetailer(db *gorm.DB) (*Retailer, error) {
	var err error
	tx := db.Begin()
	dateTimeNow := time.Now()
	err = db.Debug().Model(&Retailer{}).Unscoped().Where("id = ?", retailer.ID).Updates(
		Retailer{
			ReactivatedAt: &dateTimeNow,
		}).Error
	if err != nil {
		tx.Rollback()
		return &Retailer{}, err
	}

	retailer.ID = 0
	retailer.DeletedAt = nil
	retailer.ReactivatedAt = nil
	err = db.Debug().Model(&Retailer{}).Create(&retailer).Error
	if err != nil {
		tx.Rollback()
		return &Retailer{}, err
	}

	tx.Commit()
	return retailer, nil
}

/*func (retailer *Retailer) TerminateRetailerLater(db *gorm.DB) (int64, error) {
	var err error

	err = db.Debug().Model(&Retailer{}).Where("id = ?", retailer.ID).Updates(
		Retailer{
			DeletedAt: retailer.DeletedAt,
		}).Error

	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (retailer *Retailer) TerminateRetailerNow(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Retailer{}).Where("id = ?", retailer.ID).Delete(&Retailer{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}*/

func (Retailer) TableName() string {
	return "retailer"
}
