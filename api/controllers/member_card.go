package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

//GenerateBearerCard => Generate Bearer Card
func (server *Server) GenerateBearerCard(c *gin.Context) {
	log.Printf("Begin => Generate Bearer Card")
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	ac := models.AccountClass{}
	err = json.Unmarshal(body, &ac)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	ac.Prepare()
	errorMessages := ac.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	createdAccountClass, err := ac.CreateAccountClass(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": createdAccountClass,
	})

	log.Printf("End => Generate Bearer Card")
}

//GenerateVehicleCard => Generate Vehicle Card
func (server *Server) GenerateVehicleCard(c *gin.Context) {
	log.Printf("Begin => Generate Vehicle Card")
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	ac := models.AccountClass{}
	err = json.Unmarshal(body, &ac)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	ac.Prepare()
	errorMessages := ac.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	createdAccountClass, err := ac.CreateAccountClass(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": createdAccountClass,
	})

	log.Printf("End => Generate Vehicle Class")
}
