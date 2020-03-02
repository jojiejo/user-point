package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type GlobalVariable struct {
	ID        int       `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:50" json:"name"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (gv *GlobalVariable) Validate() map[string]string {
	var err error

	var errorMessages = make(map[string]string)

	if gv.Name == "" {
		err = errors.New("Required Name")
		errorMessages["required_name"] = err.Error()
	}

	return errorMessages
}

func (gv *GlobalVariable) FindAllGlobalVariables(db *gorm.DB) (*[]GlobalVariable, error) {
	var err error
	globalVariables := []GlobalVariable{}
	err = db.Debug().Model(&GlobalVariable{}).Limit(100).Order("created_at desc").Find(&globalVariables).Error
	if err != nil {
		return &[]GlobalVariable{}, err
	}
	if len(globalVariables) > 0 {
		return &[]GlobalVariable{}, err
	}
	return &globalVariables, nil
}
