package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetFees(c *gin.Context) {
	log.Printf("Begin => Get Fees")

	fee := models.Fee{}
	fees, err := fee.FindAllFees(server.DB)
	if err != nil {
		errString := "No fee found"
		log.Printf(errString)
		errList["no_fee"] = errString
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedFee, _ := json.Marshal(fees)
	log.Printf("Get Fees : ", string(stringifiedReceivedFee))
	c.JSON(http.StatusOK, gin.H{
		"response": fees,
	})

	log.Printf("End => Get Fees")
}

func (server *Server) GetInitialFees(c *gin.Context) {
	log.Printf("Begin => Get Initial Fees")

	fee := models.Fee{}
	fees, err := fee.FindIntialFees(server.DB)
	if err != nil {
		errString := "No fee found"
		log.Printf(errString)
		errList["no_fee"] = errString
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

func (server *Server) GetFee(c *gin.Context) {
	log.Printf("Begin => Get Fee by ID")
	feeID := c.Param("id")
	convertedFeeID, err := strconv.ParseUint(feeID, 10, 64)
	if err != nil {
		errString := "Invalid request"
		log.Printf(errString)
		errList["invalid_request"] = errString
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	fee := models.Fee{}
	feeReceived, err := fee.FindFeeByID(server.DB, convertedFeeID)
	if err != nil {
		errString := "Invalid request"
		log.Printf(errString)
		errList["no_fee"] = errString
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

	log.Printf("Begin => Get Fee by ID")
}
