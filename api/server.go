package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jojiejo/user-point/api/controllers"
)

var server = controllers.Server{}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

//Run => Run the server
func Run() {
	var err error

	//Init Log File
	logFile, err := os.OpenFile("user-point.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	//to disable logging => log.SetOutput(ioutil.Discard)

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error populating ENV, %v", err)
	}

	fmt.Println("ENV values has been populated successfully")

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	apiPort := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	fmt.Printf("\nListening to port %s", apiPort)

	server.Run(apiPort)
}
