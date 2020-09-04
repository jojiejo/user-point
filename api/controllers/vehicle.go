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

func (server *Server) GetVehicles(c *gin.Context) {
	log.Printf("Begin => Get Vehicles")

	v := models.Vehicle{}
	vs, err := v.FindVehicles(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_vehicle"] = "No vehicle found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedVehicles, _ := json.Marshal(vs)
	log.Printf("Get Vehicles : ", string(stringifiedVehicles))
	c.JSON(http.StatusOK, gin.H{
		"response": vs,
	})

	log.Printf("End => Get Vehicles")
}

func (server *Server) GetVehicle(c *gin.Context) {
	log.Printf("Begin => Get Vehicle")
	vID := c.Param("id")
	convertedVID, err := strconv.ParseUint(vID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	v := models.Vehicle{}
	receivedVehicle, err := v.FindVehicleByID(server.DB, convertedVID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_vehicle"] = "No vehicle found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedVehicle, _ := json.Marshal(receivedVehicle)
	log.Printf("Get Driver : ", string(stringifiedVehicle))
	c.JSON(http.StatusOK, gin.H{
		"response": receivedVehicle,
	})

	log.Printf("End => Get Vehicle")
}

func (server *Server) GetVehicleByPayer(c *gin.Context) {
	log.Printf("Begin => Get Vehicle By Payer")
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

	v := models.Vehicle{}
	receivedVehicle, err := v.FindVehicleByCCID(server.DB, convertedCCID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_vehicle"] = "No vehicle found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedVehicle, _ := json.Marshal(receivedVehicle)
	log.Printf("Get Driver : ", string(stringifiedVehicle))
	c.JSON(http.StatusOK, gin.H{
		"response": receivedVehicle,
	})

	log.Printf("End => Get Vehicle By Payer")
}

func (server *Server) CreateVehicle(c *gin.Context) {
	log.Printf("Begin => Create Vehicle")
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

	v := models.Vehicle{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	v.Prepare()
	errorMessages := v.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	createdVehicle, err := v.CreateVehicle(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": createdVehicle,
	})

	log.Printf("End => Create Vehicle")
}

func (server *Server) UpdateVehicle(c *gin.Context) {
	log.Printf("Begin => Update Vehicle")

	errList = map[string]string{}
	vID := c.Param("id")

	convertedVID, err := strconv.ParseUint(vID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalVehicle := models.Vehicle{}
	err = server.DB.Debug().Model(models.Vehicle{}).Where("v_id = ?", convertedVID).Order("id desc").Take(&originalVehicle).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_vehicle"] = "No vehicle found"
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

	v := models.Vehicle{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	v.ID = originalVehicle.ID
	v.Prepare()
	errorMessages := v.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	updatedVehicle, err := v.UpdateVehicle(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": updatedVehicle,
	})

	log.Printf("End => Update Vehicle")
}
