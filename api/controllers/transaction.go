package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetAllTransactionForManualSettlement(c *gin.Context) {
	log.Printf("Begin => Get All Transaction for Manual Settlement")

	transactionDate := c.Param("date")
	trx := models.TransactionForManualSettlement{}
	trxs, err := trx.FindTransactionForManualSettlement(server.DB, transactionDate)
	if err != nil {
		log.Printf(err.Error())
		errList["no_transaction"] = "No transaction found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedTransactions, _ := json.Marshal(trxs)
	log.Printf("Get Get All Transaction for Manual Settlement : ", string(stringifiedTransactions))
	c.JSON(http.StatusOK, gin.H{
		"response": trxs,
	})

	log.Printf("End => Get All Transaction for Manual Settlement")
}
