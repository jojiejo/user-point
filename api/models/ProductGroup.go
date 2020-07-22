package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type ProductGroup struct {
	ID        uint64     `gorm:"primary_key;auto_increment;column:product_group_id" json:"id"`
	Code      string     `gorm:"not null;column:product_group_code" json:"code"`
	Name      string     `gorm:"not null;size:100;column:product_group_name" json:"name"`
	StartedAt *time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (pg *ProductGroup) Prepare() {
	pg.Code = html.EscapeString(strings.TrimSpace(pg.Code))
	pg.Name = html.EscapeString(strings.TrimSpace(pg.Name))
	pg.CreatedAt = time.Now()
}

func (pg *ProductGroup) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if pg.Code == "" {
		err = errors.New("Code field is required")
		errorMessages["code"] = err.Error()
	}

	if pg.Name == "" {
		err = errors.New("Name field is required")
		errorMessages["name"] = err.Error()
	}

	if pg.StartedAt == nil {
		err = errors.New("Started at field is required")
		errorMessages["material"] = err.Error()
	}

	return errorMessages
}

func (pg *ProductGroup) FindProductGroups(db *gorm.DB) (*[]ProductGroup, error) {
	var err error
	pgs := []ProductGroup{}
	err = db.Debug().Model(&ProductGroup{}).
		Order("product_group_id, created_at desc").
		Find(&pgs).Error

	if err != nil {
		return &[]ProductGroup{}, err
	}

	return &pgs, nil
}

func (pg *ProductGroup) FindProductGroupByID(db *gorm.DB, productGroupID uint64) (*ProductGroup, error) {
	var err error
	err = db.Debug().Model(&ProductGroup{}).Unscoped().
		Where("product_group_id = ?", productGroupID).
		Order("created_at desc").
		Take(&pg).Error

	if err != nil {
		return &ProductGroup{}, err
	}

	return pg, nil
}

func (pg *ProductGroup) CreateProductGroup(db *gorm.DB) (*ProductGroup, error) {
	var err error
	err = db.Debug().Model(&ProductGroup{}).Create(&pg).Error
	if err != nil {
		return &ProductGroup{}, err
	}

	//Select created product group
	_, err = pg.FindProductGroupByID(db, pg.ID)
	if err != nil {
		return &ProductGroup{}, err
	}

	return pg, nil
}

func (pg *ProductGroup) UpdateProductGroup(db *gorm.DB) (*ProductGroup, error) {
	var err error
	dateTimeNow := time.Now()

	//Update the data
	err = db.Debug().Model(&pg).Updates(
		map[string]interface{}{
			"product_group_code": pg.Code,
			"product_group_name": pg.Name,
			"started_at":         pg.StartedAt,
			"ended_at":           pg.EndedAt,
			"updated_at":         dateTimeNow,
		}).Error

	if err != nil {
		return &ProductGroup{}, err
	}

	//Select updated product
	_, err = pg.FindProductGroupByID(db, pg.ID)
	if err != nil {
		return &ProductGroup{}, err
	}

	return pg, nil
}

func (ProductGroup) TableName() string {
	return "mstProductGroup"
}
