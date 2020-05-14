package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type RebateProgram struct {
	ID                      uint64                `gorm:"primary_key;auto_increment" json:"id"`
	Name                    string                `gorm:"not null;size:100" json:"name"`
	RebateTypeID            uint64                `gorm:"not null;" json:"rebate_type_id"`
	RebateType              RebateType            `json:"rebate_type"`
	RebateCalculationTypeID uint64                `gorm:"not null;" json:"rebate_calculation_type_id"`
	RebateCalculationType   RebateCalculationType `json:"rebate_calculation_type"`
	RebatePeriodID          uint64                `gorm:"not null;" json:"rebate_period_id"`
	RebatePeriod            RebatePeriod          `json:"rebate_period"`
	Site                    []Site                `gorm:"many2many:rebate_program_site;association_autoupdate:false;association_jointable_foreignkey:rebate_program_id;association_jointable_foreignkey:site_id" json:"site"`
	Product                 []Product             `gorm:"many2many:rebate_program_product;association_autoupdate:false;association_jointable_foreignkey:rebate_program_id;association_jointable_foreignkey:product_id" json:"product"`
	Tier                    []RebateProgramTier   `gorm:"foreignkey:RebateProgramID;association_foreignkey:ID" json:"tier"`
	StartedAt               *time.Time            `json:"started_at"`
	EndedAt                 *time.Time            `json:"ended_at"`
	CreatedAt               time.Time             `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt               time.Time             `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt               *time.Time            `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

func (rp *RebateProgram) Prepare() {
	rp.Name = html.EscapeString(strings.TrimSpace(rp.Name))
	rp.CreatedAt = time.Now()
}

func (rp *RebateProgram) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if rp.Name == "" {
		err = errors.New("Name field is required")
		errorMessages["name"] = err.Error()
	}

	if rp.RebateTypeID < 1 {
		err = errors.New("Rebate type field is required")
		errorMessages["rebate_type"] = err.Error()
	}

	if rp.RebateCalculationTypeID < 1 {
		err = errors.New("Rebate calculation field is required")
		errorMessages["rebate_calculation"] = err.Error()
	}

	if rp.RebatePeriodID < 1 {
		err = errors.New("Rebate period field is required")
		errorMessages["rebate_period"] = err.Error()
	}

	return errorMessages
}

func (rp *RebateProgram) FindRebatePrograms(db *gorm.DB) (*[]RebateProgram, error) {
	var err error
	rps := []RebateProgram{}
	err = db.Debug().Model(&RebateProgram{}).Unscoped().
		Preload("Site").
		Preload("Product").
		Preload("RebateType").
		Preload("RebatePeriod").
		Preload("RebateCalculationType").
		Preload("Tier").
		Order("id, created_at desc").
		Find(&rps).Error

	if err != nil {
		return &[]RebateProgram{}, err
	}

	return &rps, nil
}

func (rp *RebateProgram) FindRebateProgramByID(db *gorm.DB, rpID uint64) (*RebateProgram, error) {
	var err error
	err = db.Debug().Model(&Fee{}).Unscoped().
		Preload("Site").
		Preload("Product").
		Preload("RebateType").
		Preload("RebatePeriod").
		Preload("RebateCalculationType").
		Preload("Tier").
		Where("id = ?", rpID).
		Order("created_at desc").
		Take(&rp).Error

	if err != nil {
		return &RebateProgram{}, err
	}

	return rp, nil
}

func (rp *RebateProgram) CreateRebateProgram(db *gorm.DB) (*RebateProgram, error) {
	var err error
	err = db.Debug().Model(&Fee{}).Create(&rp).Error
	if err != nil {
		return &RebateProgram{}, err
	}

	//Select created fee
	_, err = rp.FindRebateProgramByID(db, rp.ID)
	if err != nil {
		return &RebateProgram{}, err
	}

	return rp, nil
}

func (rp *RebateProgram) UpdateRebateProgram(db *gorm.DB) (*RebateProgram, error) {
	var err error
	dateTimeNow := time.Now()

	//Update tier
	err = db.Debug().Where("rebate_program_id = ?", rp.ID).Delete(RebateProgramTier{}).Error
	if err != nil {
		return &RebateProgram{}, err
	}

	//Update site
	err = db.Debug().Model(&rp).Where("id = ?", rp.ID).Association("Site").Replace(rp.Site).Error
	if err != nil {
		return &RebateProgram{}, err
	}

	//Update product
	err = db.Debug().Model(&rp).Where("id = ?", rp.ID).Association("Product").Replace(rp.Product).Error
	if err != nil {
		return &RebateProgram{}, err
	}

	err = db.Debug().Model(&rp).Updates(
		map[string]interface{}{
			"name":                       rp.Name,
			"rebate_type_id":             rp.RebateTypeID,
			"rebate_calculation_type_id": rp.RebateCalculationTypeID,
			"rebate_period_id":           rp.RebatePeriodID,
			"started_at":                 rp.StartedAt,
			"ended_at":                   rp.EndedAt,
			"updated_at":                 dateTimeNow,
		}).Error

	if err != nil {
		return &RebateProgram{}, err
	}

	//Select updated fee
	_, err = rp.FindRebateProgramByID(db, rp.ID)
	if err != nil {
		return &RebateProgram{}, err
	}

	return rp, nil
}

func (rp *RebateProgram) DeactivateRebateProgramLater(db *gorm.DB) (int64, error) {
	var err error
	err = db.Debug().Model(&rp).Unscoped().Updates(
		RebateProgram{
			EndedAt: rp.EndedAt,
		}).Error

	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (RebatePeriod) TableName() string {
	return "rebate_period"
}
