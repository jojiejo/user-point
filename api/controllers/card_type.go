package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetCardTypes(c *gin.Context) {
	log.Printf("Begin => Get Card Types")

	cardType := models.CardType{}
	cardTypes, err := cardType.FindCardTypes(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_card_types"] = "No card type found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedCardTypes, _ := json.Marshal(cardTypes)
	log.Printf("Get Card Types : ", string(stringifiedCardTypes))
	c.JSON(http.StatusOK, gin.H{
		"response": cardTypes,
	})

	log.Printf("End => Get Card Types")
}

func (server *Server) GetCardType(c *gin.Context) {
	log.Printf("Begin => Get Card Type")
	cardTypeID := c.Param("id")
	convertedCardTypeID, err := strconv.ParseUint(cardTypeID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	cardType := models.CardType{}
	cardTypeReceived, err := cardType.FindCardTypeByID(server.DB, convertedCardTypeID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_card_type"] = "No card type found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedCardType, _ := json.Marshal(cardTypeReceived)
	log.Printf("Get Card Type : ", string(stringifiedReceivedCardType))
	c.JSON(http.StatusOK, gin.H{
		"response": cardTypeReceived,
	})

	log.Printf("End => Get Card Type")
}

func (server *Server) CreateCardType(c *gin.Context) {
	log.Printf("Begin => Create Card Type")
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

	cardType := models.CardType{}
	err = json.Unmarshal(body, &cardType)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	cardType.Prepare()
	errorMessages := cardType.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	cardTypeCreated, err := cardType.CreateCardType(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": cardTypeCreated,
	})

	log.Printf("End => Create Card Type")
}

func (server *Server) UpdateCardType(c *gin.Context) {
	log.Printf("Begin => Update Card Type")

	errList = map[string]string{}
	cardTypeID := c.Param("id")

	convertedCardTypeID, err := strconv.ParseUint(cardTypeID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalCardType := models.CardType{}
	err = server.DB.Debug().Model(models.CardType{}).Where("card_type_id = ?", convertedCardTypeID).Order("card_type_id desc").Take(&originalCardType).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_card_type"] = "No card type found"
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

	cardType := models.CardType{}
	err = json.Unmarshal(body, &cardType)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	cardType.ID = originalCardType.ID
	cardType.Prepare()
	errorMessages := cardType.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	cardTypeUpdated, err := cardType.UpdateCardType(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": cardTypeUpdated,
	})

	log.Printf("End => Update Card Type")
}
