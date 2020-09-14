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

func (server *Server) GetDrivers(c *gin.Context) {
	log.Printf("Begin => Get Drivers")

	d := models.Driver{}
	ds, err := d.FindDrivers(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_driver"] = "No driver found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedDrivers, _ := json.Marshal(ds)
	log.Println("Get Drivers : ", string(stringifiedDrivers))
	c.JSON(http.StatusOK, gin.H{
		"response": ds,
	})

	log.Printf("End => Get Drivers")
}

func (server *Server) GetDriver(c *gin.Context) {
	log.Printf("Begin => Get Driver")
	dID := c.Param("id")
	convertedDID, err := strconv.ParseUint(dID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	d := models.Driver{}
	receivedDriver, err := d.FindDriverByID(server.DB, convertedDID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_driver"] = "No driver found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedDriver, _ := json.Marshal(receivedDriver)
	log.Println("Get Driver : ", string(stringifiedDriver))
	c.JSON(http.StatusOK, gin.H{
		"response": receivedDriver,
	})

	log.Printf("End => Get Driver")
}

func (server *Server) GetDriverByPayer(c *gin.Context) {
	log.Printf("Begin => Get Driver By Payer")
	ccID := c.Param("id")
	convertedCCID, err := strconv.ParseUint(ccID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	d := models.Driver{}
	receivedDriver, err := d.FindDriverByID(server.DB, convertedCCID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_driver"] = "No driver found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedDriver, _ := json.Marshal(receivedDriver)
	log.Println("Get Driver by Payer: ", string(stringifiedDriver))
	c.JSON(http.StatusOK, gin.H{
		"response": receivedDriver,
	})

	log.Printf("End => Get Driver By Payer")
}

func (server *Server) CreateDriver(c *gin.Context) {
	log.Printf("Begin => Create Driver")
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

	d := models.Driver{}
	err = json.Unmarshal(body, &d)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	d.Prepare()
	errorMessages := d.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	createdDriver, err := d.CreateDriver(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": createdDriver,
	})

	log.Printf("End => Create Driver")
}

func (server *Server) UpdateDriver(c *gin.Context) {
	log.Printf("Begin => Update Driver")

	errList = map[string]string{}
	dID := c.Param("id")

	convertedDID, err := strconv.ParseUint(dID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalDriver := models.Driver{}
	err = server.DB.Debug().Model(models.Driver{}).Where("card_holder_id = ?", convertedDID).Order("id desc").Take(&originalDriver).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_driver"] = "No driver found"
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

	d := models.Driver{}
	err = json.Unmarshal(body, &d)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	d.ID = originalDriver.ID
	d.Prepare()
	errorMessages := d.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	updatedDriver, err := d.UpdateDriver(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": updatedDriver,
	})

	log.Printf("End => Update Driver")
}
