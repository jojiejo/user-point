package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jojiejo/user-point/api/security"
)

//User => User of this system
type User struct {
	ID        uint32     `gorm:"primary_key;auto_increment" json:"id"`
	Email     string     `gorm:"size:100;not null;unique" json:"email"`
	Password  string     `gorm:"size:100;not null;" json:"-"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

//HashPassword => Hash Password of User
func (u *User) HashPassword() error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

//Prepare => Prepare the string & time
func (u *User) Prepare() {
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

//ValidateInsertion => Validate the input when insert data
func (u *User) ValidateInsertion() map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	if u.Email == "" {
		err = errors.New("Email is required")
		errorMessages["email"] = err.Error()
	}

	if u.Password == "" {
		err = errors.New("Password is required")
		errorMessages["password"] = err.Error()
	}

	return errorMessages
}

//ValidateUpdate => Validate the input when update user data
func (u *User) ValidateUpdate() map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	if u.ID == 0 {
		err = errors.New("ID is required")
		errorMessages["email"] = err.Error()
	}

	return errorMessages
}

//ValidateUpdatePoint => Validate the input when update point
func (u *User) ValidateUpdatePoint() map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	if u.ID == 0 {
		err = errors.New("ID is required")
		errorMessages["email"] = err.Error()
	}

	return errorMessages
}

//TableName => Define tablename
func (User) TableName() string {
	return "user"
}
