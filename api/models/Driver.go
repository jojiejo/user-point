package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

//Driver => Driver struct
type Driver struct {
	ID        uint64     `gorm:"primary_key;auto_increment;column:card_holder_id" json:"id"`
	FleetID   string     `gorm:"not null;" json:"fleet_id"`
	Name      string     `gorm:"not null;column:card_holder_name" json:"name"`
	Address   string     `gorm:"column:card_holder_address" json:"address"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	CCID      uint64     `gorm:"not null;" json:"cc_id"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

//Prepare => Prepare driver string
func (d *Driver) Prepare() {
	d.FleetID = html.EscapeString(strings.TrimSpace(d.FleetID))
	d.Name = html.EscapeString(strings.TrimSpace(d.Name))
	d.Address = html.EscapeString(strings.TrimSpace(d.Address))
	d.Email = html.EscapeString(strings.TrimSpace(d.Email))
	d.Phone = html.EscapeString(strings.TrimSpace(d.Phone))
	d.CreatedAt = time.Now()
}

//Validate => Validate driver input
func (d *Driver) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if d.FleetID == "" {
		err = errors.New("Fleet ID field is required")
		errorMessages["fleet_id"] = err.Error()
	}

	if len(d.FleetID) > 10 {
		err = errors.New("Fleet ID field must contain maximum of 10 characters")
		errorMessages["fleet_id"] = err.Error()
	}

	if d.Name == "" {
		err = errors.New("Name field is required")
		errorMessages["name"] = err.Error()
	}

	return errorMessages
}

//FindDrivers => Find All Drivers
func (d *Driver) FindDrivers(db *gorm.DB) (*[]Driver, error) {
	var err error
	ds := []Driver{}
	err = db.Debug().Model(&Driver{}).
		Order("id, created_at desc").
		Find(&ds).Error

	if err != nil {
		return &[]Driver{}, err
	}

	return &ds, nil
}

//FindDriverByID => Find Driver By ID
func (d *Driver) FindDriverByID(db *gorm.DB, dID uint64) (*Driver, error) {
	var err error
	err = db.Debug().Model(&Vehicle{}).Unscoped().
		Where("card_holder_id = ?", dID).
		Order("created_at desc").
		Take(&d).Error

	if err != nil {
		return &Driver{}, err
	}

	return d, nil
}

//FindDriverByCCID => Find Driver by CC ID
func (d *Driver) FindDriverByCCID(db *gorm.DB, ccID uint64) (*[]Driver, error) {
	var err error
	ds := []Driver{}
	err = db.Debug().Model(&Driver{}).Unscoped().
		Where("card_holder_id = ?", ccID).
		Order("created_at desc").
		Take(&ds).Error

	if err != nil {
		return &ds, err
	}

	return &ds, nil
}

//CreateDriver => Create Driver
func (d *Driver) CreateDriver(db *gorm.DB) (*Driver, error) {
	var err error
	err = db.Debug().Model(&AccountClass{}).Create(&d).Error
	if err != nil {
		return &Driver{}, err
	}

	//Select created fee
	_, err = d.FindDriverByID(db, d.ID)
	if err != nil {
		return &Driver{}, err
	}

	return d, nil
}

//UpdateDriver => Update Driver
func (d *Driver) UpdateDriver(db *gorm.DB) (*Driver, error) {
	var err error
	dateTimeNow := time.Now()

	//Update the data
	err = db.Debug().Model(&d).Updates(
		map[string]interface{}{
			"fleet_id":         d.FleetID,
			"card_holder_name": d.Name,
			"updated_at":       dateTimeNow,
		}).Error

	if err != nil {
		return &Driver{}, err
	}

	_, err = d.FindDriverByID(db, d.ID)
	if err != nil {
		return &Driver{}, err
	}

	return d, nil
}

//TableName => Define Table
func (Driver) TableName() string {
	return "mstCardHolder"
}
