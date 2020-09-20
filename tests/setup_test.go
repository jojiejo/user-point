package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql" //Ms. SQL driver
	"github.com/joho/godotenv"
	"github.com/jojiejo/user-point/api/controllers"
	"github.com/jojiejo/user-point/api/models"
)

var server = controllers.Server{}
var userInstance = models.User{}
var userPointInstance = models.UserPoint{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("./../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	Database()

	os.Exit(m.Run())
}

func Database() {
	var err error

	DbDriver := os.Getenv("DB_DRIVER")
	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
	server.DB, err = gorm.Open(DbDriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", DbDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", DbDriver)
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshing user table")
	return nil
}

func seedUser() (models.User, error) {
	user := models.User{
		Email: "djodi.ramadhan@example.com",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func seedUsers() ([]models.User, error) {
	var err error
	if err != nil {
		return nil, err
	}
	users := []models.User{
		models.User{
			Email: "djodi@example.com",
		},
		models.User{
			Email: "ramadhan@example.com",
		},
	}

	for i := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}

	return users, nil
}

func refreshUserAndUserPointTable() error {
	err := server.DB.DropTableIfExists(&models.UserPoint{}, &models.User{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.UserPoint{}, &models.User{}).Error
	if err != nil {
		return err
	}

	err = server.DB.Model(&models.UserPoint{}).AddForeignKey("user_id", "user(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshing user and user point tables")
	return nil
}

func seedUserAndUserPoint() (models.User, models.UserPoint, error) {

	user := models.User{
		Email: "djodi@example.com",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, models.UserPoint{}, err
	}

	userPoint := models.UserPoint{
		Value:  -10,
		UserID: user.ID,
	}

	err = server.DB.Model(&models.UserPoint{}).Create(&userPoint).Error
	if err != nil {
		return models.User{}, models.UserPoint{}, err
	}

	return user, userPoint, nil
}
