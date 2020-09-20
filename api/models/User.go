package models

import (
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jojiejo/user-point/api/security"
)

//User => User of this system
type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
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

//Validate => Validate the input
func (u *User) ValidateInsertion() {
	var errorMessages = make(map[string]string)
	var err error
}

//ValidateUpdate =>
func (u *User) ValidateUpdate() {
	var errorMessages = make(map[string]string)
	var err error
}

//ValidateUpdatePoint
func (u *User) ValidateInsertion() {
	var errorMessages = make(map[string]string)
	var err error
}

func (unit *Unit) FindAllUnits(db *gorm.DB) (*[]Unit, error) {
	var err error
	units := []Unit{}
	err = db.Debug().Model(&Unit{}).Limit(100).Order("id asc").Find(&units).Error
	if err != nil {
		return &[]Unit{}, err
	}

	return &units, nil
}

func (Unit) TableName() string {
	return "unit"
}
