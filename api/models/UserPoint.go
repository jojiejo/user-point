package models

import (
	"errors"
	"time"
)

//UserPoint => User of this system
type UserPoint struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
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

//TableName => Define tablename
func (UserPoint) TableName() string {
	return "user_point"
}
