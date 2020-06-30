package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetTelematicDeviceByCardID(c *gin.Context) {
	log.Printf("Begin => Get Telematic Device By Card ID")

	cardID := c.Param("id")
	td := models.CardTelematicDevice{}
	tdReceived, err := td.FindTelematicDeviceByCardID(server.DB, cardID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_transaction"] = "No telematic device found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedTelematicDevice, _ := json.Marshal(tdReceived)
	log.Printf("Get Telematic Device By Card ID : ", string(stringifiedTelematicDevice))
	c.JSON(http.StatusOK, gin.H{
		"response": tdReceived,
	})

	log.Printf("End => Get Telematic Device By Card ID")
}

func (server *Server) UpdateTelematicDevice(c *gin.Context) {
	log.Printf("Begin => Update Telematic Device")

	errList = map[string]string{}
	cardID := c.Param("id")
	originalCtd := models.CardTelematicDevice{}
	err := server.DB.Debug().Model(models.CardTelematicDevice{}).Where("card_id = ?", cardID).Order("card_id desc").Take(&originalCtd).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_card"] = "No card found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	ctd := models.CardTelematicDevice{}
	err = json.Unmarshal(body, &ctd)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	ctd.CardID = originalCtd.CardID
	ctd.Prepare()
	errorMessages := ctd.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	ctdUpdated, err := ctd.UpdateTelematicDevice(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": ctdUpdated,
	})

	log.Printf("End => Update Telematic Device")
}
