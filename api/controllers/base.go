package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jojiejo/user-point/api/middlewares"
	"github.com/jojiejo/user-point/api/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //Ms. SQL driver
)

//Server => Struct of Server Attributes
type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

var errList = make(map[string]string)

//Initialize => Init server
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error

	// Initialize DB Connection
	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		log.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		log.Printf("The service has been connected to the %s database", Dbdriver)
	}

	// Database Migration
	server.DB.Debug().AutoMigrate(
		&models.User{},
	)

	gin.SetMode(gin.ReleaseMode)
	server.Router = gin.Default()
	server.Router.Use(middlewares.CORSMiddleware())

	server.initializeRoutes()
}

//Run => Run server
func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
