package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

//UserPoint => User of this system
type UserPoint struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
	Value     float32   `gorm:"not null" json:"value"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

//Prepare => Prepare the string & time
func (up *UserPoint) Prepare() {
	up.CreatedAt = time.Now()
	up.UpdatedAt = time.Now()
}

//ValidateInsertion => Validate the input when insert data
func (up *UserPoint) ValidateInsertion() map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	if up.UserID == 0 {
		err = errors.New("User ID is required")
		errorMessages["user_id"] = err.Error()
	}

	if up.Value == 0 {
		err = errors.New("Value is required")
		errorMessages["value"] = err.Error()
	}

	return errorMessages
}

//FindPointHistoryByUserID => Find point history by user id
func (up *UserPoint) FindPointHistoryByUserID(db *gorm.DB, userID uint64) (*[]UserPoint, error) {
	var err error
	ups := []UserPoint{}
	err = db.Debug().
		Model(&UserPoint{}).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&ups).
		Error

	if err != nil {
		return &[]UserPoint{}, err
	}

	return &ups, nil
}

//CreateUserPoint => Create user point
func (up *UserPoint) CreateUserPoint(db *gorm.DB) (*UserPoint, error) {
	var err error
	err = db.Debug().
		Model(&UserPoint{}).
		Create(&up).
		Error

	if err != nil {
		return &UserPoint{}, err
	}

	return up, nil
}

//TableName => Define tablename
func (UserPoint) TableName() string {
	return "user_point"
}
