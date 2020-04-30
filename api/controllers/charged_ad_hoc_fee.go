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

func (server *Server) GetChargedAdHocFees(c *gin.Context) {
	log.Printf("Begin => Get Charged Ad Hoc Fees")

	adHocFee := models.ChargedAdHocFee{}
	adHocFees, err := adHocFee.FindAdHocFees(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_ad_hoc_fee"] = "No ad hoc fee found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedAdHocFees, _ := json.Marshal(adHocFees)
	log.Printf("Get Charged Ad Hoc Fees : ", string(stringifiedAdHocFees))
	c.JSON(http.StatusOK, gin.H{
		"response": adHocFees,
	})

	log.Printf("End => Get Charged Ad Hoc Fees")
}

func (server *Server) GetChargedAdHocFee(c *gin.Context) {
	log.Printf("Begin => Get Charged Ad Hoc Fee by ID")
	relationID := c.Param("id")
	convertedRelationID, err := strconv.ParseUint(relationID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	chargedAdHocFee := models.ChargedAdHocFee{}
	relationReceived, err := chargedAdHocFee.FindChargedAdHocFeeByID(server.DB, convertedRelationID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_fee"] = "No fee found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedRelationReceived, _ := json.Marshal(relationReceived)
	log.Printf("Get Charged Ad Hoc Fee : ", string(stringifiedRelationReceived))
	c.JSON(http.StatusOK, gin.H{
		"response": chargedAdHocFee,
	})

	log.Printf("Begin => Get Charged Ad Hoc Fee by ID")
}

func (server *Server) ChargeAdHocFee(c *gin.Context) {
	log.Printf("Begin => Charge Ad Hoc Fee")
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

	chargedAdHocFee := models.ChargedAdHocFee{}
	err = json.Unmarshal(body, &chargedAdHocFee)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	chargedAdHocFee.Prepare()
	errorMessages := chargedAdHocFee.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	feeCharged, err := chargedAdHocFee.ChargeAdHocFee(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"response": feeCharged,
	})

	log.Printf("End => Create Fee")
}

func (server *Server) CheckBulkChargeAdHocFee(c *gin.Context) {
	log.Printf("Begin => Check Bulk Charge Ad Hoc Fee")
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

	bulkChargeAdHocFee := models.BulkChargeAdHocFee{}
	err = json.Unmarshal(body, &bulkChargeAdHocFee)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	bulkChargeAdHocFee.Prepare()
	errorMessages := bulkChargeAdHocFee.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	checkedFieldBulkUpload, errorMessages := bulkChargeAdHocFee.BulkCheckAdHocFee(server.DB)
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		c.JSON(http.StatusNotFound, gin.H{
			"error": errorMessages,
		})
		return
	} else {
		stringifiedCheckedFieldBulkUpload, _ := json.Marshal(checkedFieldBulkUpload)
		log.Printf("Get Bulk Uploaded Field : ", string(stringifiedCheckedFieldBulkUpload))
		c.JSON(http.StatusOK, gin.H{
			"response": checkedFieldBulkUpload,
		})
	}

	log.Printf("End => Check Bulk Charge Ad Hoc Fee")
}

func (server *Server) BulkChargeAdHocFee(c *gin.Context) {
	log.Printf("Begin => Charge Ad Hoc Fee")
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

	bulkChargeAdHocFee := models.BulkChargeAdHocFee{}
	err = json.Unmarshal(body, &bulkChargeAdHocFee)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	bulkChargeAdHocFee.Prepare()
	errorMessages := bulkChargeAdHocFee.ChargeValidate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	fieldBulkCharge, errorMessages := bulkChargeAdHocFee.BulkChargeAdHocFee(server.DB)
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		c.JSON(http.StatusNotFound, gin.H{
			"error": errorMessages,
		})
		return
	} else {
		stringifiedFieldBulkCharge, _ := json.Marshal(fieldBulkCharge)
		log.Printf("Get Bulk Charge Ad Hoc Fee : ", string(stringifiedFieldBulkCharge))
		c.JSON(http.StatusOK, gin.H{
			"response": fieldBulkCharge,
		})
	}

	log.Printf("End => Bulk Charge Ad Hoc Fee")
}
