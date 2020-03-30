package controllers

import (
	"encoding/json"
	"log"
	"net/http"

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
