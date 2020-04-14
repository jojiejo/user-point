package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetInitialFees(c *gin.Context) {
	log.Printf("Begin => Get Initial Fees")

	fee := models.Fee{}
	fees, err := fee.FindIntialFees(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_fee"] = "No fee found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedFee, _ := json.Marshal(fees)
	log.Printf("Get Initial Fees : ", string(stringifiedReceivedFee))
	c.JSON(http.StatusOK, gin.H{
		"response": fees,
	})

	log.Printf("End => Get Initial Fees")
}

func (server *Server) GetInitialFee(c *gin.Context) {
	log.Printf("Begin => Get Fee by ID")
	feeID := c.Param("id")
	convertedFeeID, err := strconv.ParseUint(feeID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	fee := models.Fee{}
	feeReceived, err := fee.FindFeeByID(server.DB, convertedFeeID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_fee"] = "No fee found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedFee, _ := json.Marshal(feeReceived)
	log.Printf("Get Fees : ", string(stringifiedReceivedFee))
	c.JSON(http.StatusOK, gin.H{
		"response": fee,
	})

	log.Printf("End => Get Fee by ID")
}

func (server *Server) CreateFee(c *gin.Context) {
	log.Printf("Begin => Create Fee")
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

	fee := models.Fee{}
	err = json.Unmarshal(body, &fee)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	var count int
	err = server.DB.Debug().Model(models.Fee{}).Where("name = ?", fee.Name).Count(&count).Error
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if count > 0 {
		errString := "Entered name already exists"
		log.Printf(errString)
		errList["name"] = errString
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	fee.Prepare()
	errorMessages := fee.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	feeCreated, err := fee.CreateFee(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": feeCreated,
	})

	log.Printf("End => Create Fee")
}

func (server *Server) UpdateFee(c *gin.Context) {
	log.Printf("Begin => Update Fee")

	errList = map[string]string{}
	feeID := c.Param("id")

	feeid, err := strconv.ParseUint(feeID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalFee := models.Fee{}
	err = server.DB.Debug().Model(models.Fee{}).Where("id = ?", feeid).Order("id desc").Take(&originalFee).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_fee"] = "No fee found"
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

	fee := models.Fee{}
	err = json.Unmarshal(body, &fee)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	fee.ID = originalFee.ID
	fee.Prepare()
	errorMessages := fee.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	feeUpdated, err := fee.UpdateFee(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": feeUpdated,
	})

	log.Printf("End => Update Fee")
}

func (server *Server) DeactivateFee(c *gin.Context) {
	log.Printf("Begin => Deactivate Fee")

	errList = map[string]string{}
	feeID := c.Param("id")
	feeid, err := strconv.ParseUint(feeID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalFee := models.Fee{}
	err = server.DB.Debug().Unscoped().Model(models.Site{}).Where("id = ?", feeid).Order("id desc").Take(&originalFee).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_fee"] = "No fee found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	//Check if the new deleted_at input is greater than the previous deleted_at
	if originalFee.DeletedAt != nil {
		dateTimeNow := time.Now()
		if dateTimeNow.After(*originalFee.DeletedAt) {
			errString := "Ended at time field can not be updated"
			log.Printf(errString)
			errList["time_exceeded"] = errString
			c.JSON(http.StatusNotFound, gin.H{
				"error": errList,
			})
			return
		}
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

	fee := models.Fee{}
	err = json.Unmarshal(body, &fee)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	fee.ID = originalFee.ID
	fee.Prepare()
	_, err = fee.DeactivateFeeLater(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Selected fee has been deactivated successfully.",
	})

	log.Printf("End => Deactivate Fee")
}
