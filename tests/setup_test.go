package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"fleethub.shell.co.id/api/controllers"
	"fleethub.shell.co.id/api/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql" //Ms. SQL driver
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var siteInstance = models.Site{}

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
	DBURL := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	server.DB, err = gorm.Open(DbDriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", DbDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", DbDriver)
	}
}
