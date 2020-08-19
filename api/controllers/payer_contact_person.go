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

func (server *Server) GetPayerContactPersons(c *gin.Context) {
	log.Printf("Begin => Payer Contact Persons")

	ccID := c.Param("id")
	convertedCCID, err := strconv.ParseUint(ccID, 10, 64)

	contactPerson := models.PayerContactPerson{}
	contactPersons, err := contactPerson.FindPayerContactPersons(server.DB, convertedCCID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_contact_person"] = "No contact person found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedContactPersons, _ := json.Marshal(contactPersons)
	log.Printf("Get Contact Persons : ", string(stringifiedContactPersons))
	c.JSON(http.StatusOK, gin.H{
		"response": contactPersons,
	})

	log.Printf("End => Payer Contact Persons")
}

func (server *Server) GetPayerContactPerson(c *gin.Context) {
	log.Printf("Begin => Get Payer Contact Person")

	pcpID := c.Param("id")
	convertedContactPersonID, err := strconv.ParseUint(pcpID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	contactPerson := models.PayerContactPerson{}
	contactPersonReceived, err := contactPerson.FindPayerContactPersonByID(server.DB, convertedContactPersonID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_product"] = "No contact person found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedContactPerson, _ := json.Marshal(contactPersonReceived)
	log.Printf("Get Product : ", string(stringifiedReceivedContactPerson))
	c.JSON(http.StatusOK, gin.H{
		"response": contactPersonReceived,
	})

	log.Printf("End => Get Payer Contact Person")
}

func (server *Server) CreatePayerContactPerson(c *gin.Context) {
	log.Printf("Begin => Create Payer Contact Person")
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

	contactPerson := models.PayerContactPerson{}
	err = json.Unmarshal(body, &contactPerson)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	contactPerson.Prepare()
	errorMessages := contactPerson.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	contactPersonCreated, err := contactPerson.CreatePayerContactPerson(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": contactPersonCreated,
	})

	log.Printf("End => Create Payer Contact Person")
}

func (server *Server) UpdatePayerContactPerson(c *gin.Context) {
	log.Printf("Begin => Update Payer Contact Person")

	errList = map[string]string{}
	contactPersonID := c.Param("id")
	convertedContactPersonID, err := strconv.ParseUint(contactPersonID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalContactPerson := models.PayerContactPerson{}
	err = server.DB.Debug().Model(models.PayerContactPerson{}).Where("id = ?", convertedContactPersonID).Order("id desc").Take(&originalContactPerson).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_product"] = "No contact person found"
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

	contactPerson := models.PayerContactPerson{}
	err = json.Unmarshal(body, &contactPerson)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	contactPerson.ID = originalContactPerson.ID
	contactPerson.Prepare()
	errorMessages := contactPerson.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	contactPersonUpdated, err := contactPerson.UpdatePayerContactPerson(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": contactPersonUpdated,
	})

	log.Printf("End => Update Payer Contact Person")
}
