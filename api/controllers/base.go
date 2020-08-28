package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

    "github.com/jlaffaye/ftp"

	"fleethub.shell.co.id/api/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql" //Ms. SQL driver
)

type Server struct {
	DB     	*gorm.DB
	Router 	*gin.Engine
	FtpConn *ftp.ServerConn
}

var errList = make(map[string]string)

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
    var ftpHost = os.Getenv("FTP_HOST")

    // INITIALIZE DB CONNECTION
	DBURL := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", DbUser, DbPassword, DbHost, DbPort, DbName)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		log.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	}
	log.Printf("The service has been connected to the %s database", Dbdriver)

    // INITIALIZE FTP CONNECTION
    server.FtpConn, err = ftp.Dial(ftpHost)
    if err != nil {
        log.Println(err.Error())
    }
    log.Println("FTP Connection Success")

    // FTP LOGIN
    err = server.FtpConn.Login(os.Getenv("FTP_USER"), os.Getenv("FTP_PASSWORD"))
    if err != nil {
        log.Println(err.Error())
    }
    log.Println("FTP Login Success")

	gin.SetMode(gin.ReleaseMode)
	server.Router = gin.Default()
	server.Router.Use(middlewares.CORSMiddleware())

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
