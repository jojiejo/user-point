package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

//User => User of this system
type User struct {
	ID           uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Email        string     `gorm:"size:100;not null;unique" json:"email"`
	CurrentPoint float32    `gorm:"not null;default:0" json:"current_point"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
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

	return errorMessages
}

//FindAllUsers => Find all users
func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().
		Model(&User{}).
		Order("id, created_at desc").
		Find(&users).
		Error

	if err != nil {
		return &[]User{}, err
	}

	return &users, nil
}

//FindUserByID => Find user by ID
func (u *User) FindUserByID(db *gorm.DB, ID uint64) (*User, error) {
	var err error
	err = db.Debug().
		Model(&User{}).
		Where("id = ?", ID).
		Order("created_at desc").
		Take(&u).
		Error

	if err != nil {
		return &User{}, err
	}

	return u, nil
}

//CreateUser => Create user
func (u *User) CreateUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().
		Model(&User{}).
		Create(&u).
		Error

	if err != nil {
		return &User{}, err
	}

	return u, nil
}

//UpdateUserPoint => Update user point
func (u *User) UpdateUserPoint(db *gorm.DB) (*User, error) {
	var err error
	dateTimeNow := time.Now()
	err = db.Debug().Model(&u).Updates(
		map[string]interface{}{
			"current_point": u.CurrentPoint,
			"updated_at":    dateTimeNow,
		}).Error

	if err != nil {
		return &User{}, err
	}

	_, err = u.FindUserByID(db, u.ID)
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

//DeleteUser => Delete user
func (u *User) DeleteUser(db *gorm.DB) (int64, error) {
	db = db.Debug().
		Model(&User{}).
		Where("id = ?", u.ID).
		Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

//TableName => Define tablename
func (User) TableName() string {
	return "user"
}
