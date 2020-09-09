package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Vehicle struct {
	ID        uint64     `gorm:"primary_key;auto_increment;column:v_id" json:"id"`
	VehicleID string     `gorm:"not null;" json:"vehicle_id"`
	CCID      uint64     `gorm:"not null;" json:"cc_id"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (v *Vehicle) Prepare() {
	v.VehicleID = html.EscapeString(strings.TrimSpace(v.VehicleID))
	v.CreatedAt = time.Now()
}

func (v *Vehicle) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if v.VehicleID == "" {
		err = errors.New("Vehicle ID field is required")
		errorMessages["vehicle_id"] = err.Error()
	}

	return errorMessages
}

func (v *Vehicle) FindVehicles(db *gorm.DB) (*[]Vehicle, error) {
	var err error
	vs := []Vehicle{}
	err = db.Debug().Model(&Vehicle{}).
		Order("id, created_at desc").
		Find(&vs).Error

	if err != nil {
		return &[]Vehicle{}, err
	}

	return &vs, nil
}

func (v *Vehicle) FindVehicleByID(db *gorm.DB, vID uint64) (*Vehicle, error) {
	var err error
	err = db.Debug().Model(&Vehicle{}).Unscoped().
		Where("v_id = ?", vID).
		Order("created_at desc").
		Take(&v).Error

	if err != nil {
		return &Vehicle{}, err
	}

	return v, nil
}

func (v *Vehicle) FindVehicleByCCID(db *gorm.DB, ccID uint64) (*Vehicle, error) {
	var err error
	err = db.Debug().Model(&Vehicle{}).Unscoped().
		Where("cc_id = ?", ccID).
		Order("created_at desc").
		Take(&v).Error

	if err != nil {
		return &Vehicle{}, err
	}

	return v, nil
}

func (v *Vehicle) CreateVehicle(db *gorm.DB) (*Vehicle, error) {
	var err error
	err = db.Debug().Model(&AccountClass{}).Create(&v).Error
	if err != nil {
		return &Vehicle{}, err
	}

	//Select created fee
	_, err = v.FindVehicleByID(db, v.ID)
	if err != nil {
		return &Vehicle{}, err
	}

	return v, nil
}

func (v *Vehicle) UpdateVehicle(db *gorm.DB) (*Vehicle, error) {
	var err error
	dateTimeNow := time.Now()

	//Update the data
	err = db.Debug().Model(&v).Updates(
		map[string]interface{}{
			"vehicle_id": v.VehicleID,
			"cc_id":      v.CCID,
			"updated_at": dateTimeNow,
		}).Error

	if err != nil {
		return &Vehicle{}, err
	}

	//Select updated sales rep
	_, err = v.FindVehicleByID(db, v.ID)
	if err != nil {
		return &Vehicle{}, err
	}

	return v, nil
}

func (Vehicle) TableName() string {
	return "mstVehicle"
}
