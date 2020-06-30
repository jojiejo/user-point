package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetAllTransactionForManualSettlement(c *gin.Context) {
	log.Printf("Begin => Get All Transaction for Manual Settlement")

	transactionDateFrom := c.Param("dateFrom")
	transactionDateTo := c.Param("dateTo")
	trx := models.TransactionForManualSettlement{}
	trxs, err := trx.FindTransactionForManualSettlement(server.DB, transactionDateFrom, transactionDateTo)
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

func (server *Server) ManualSettle(c *gin.Context) {
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	trx := models.TransactionForManualSettlement{}
	err = json.Unmarshal(body, &trx)
	if err != nil {
		fmt.Println(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	trx.Prepare()
	errorMessages := trx.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	manualSettlement, err := trx.ManualSettle(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": manualSettlement,
	})
}

func (server *Server) ManualSettleAllTransaction(c *gin.Context) {
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	trx := models.TransactionForManualSettlement{}
	err = json.Unmarshal(body, &trx)
	if err != nil {
		fmt.Println(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	trx.Prepare()
	errorMessages := trx.ValidateForSettleAllTransaction()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	err = trx.ManualSettleAllTransaction(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": "All transactions have been settled ",
	})
}
