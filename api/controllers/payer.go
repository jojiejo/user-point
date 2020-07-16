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

func (server *Server) GetPayers(c *gin.Context) {
	payer := models.ShortenedPayer{}
	payers, err := payer.FindAllPayers(server.DB)
	if err != nil {
		errList["no_payer"] = "No payer found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": payers,
	})
}

func (server *Server) GetPayer(c *gin.Context) {
	payerID := c.Param("id")
	convertedPayerID, err := strconv.ParseUint(payerID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	payer := models.Payer{}
	payerReceived, err := payer.FindPayerByCCID(server.DB, convertedPayerID)
	if err != nil {
		errList["no_payer"] = "No payer found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": payerReceived,
	})
}

func (server *Server) UpdateConfiguration(c *gin.Context) {
	errList = map[string]string{}
	CCID := c.Param("id")
	ccid, err := strconv.ParseUint(CCID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalPayer := models.Payer{}
	err = server.DB.Debug().Model(models.Payer{}).Where("cc_id = ?", ccid).Order("cc_id desc").Take(&originalPayer).Error
	if err != nil {
		errList["no_payer"] = "No payer found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	payer := models.ShortenedPayer{}
	err = json.Unmarshal(body, &payer)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	payer.CCID = originalPayer.CCID
	payer.Prepare()
	errorMessages := payer.ValidateConfiguration()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	payerUpdated, err := payer.UpdatePayerConfiguration(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": payerUpdated,
	})
}

func (server *Server) GetTransactionInvoiceByPayer(c *gin.Context) {
	log.Printf("Begin => Get Transaction Invoice By Payer")
	CCID := c.Param("id")
	month := c.Param("month")
	year := c.Param("year")

	convertedYear, _ := strconv.Atoi(year)
	convertedMonth, _ := strconv.Atoi(month)
	firstDay := time.Date(convertedYear, time.Month(convertedMonth), 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, 0).Add(time.Nanosecond * -1)

	parsedFirstDay := firstDay.Format("20060102")
	parsedLastDay := lastDay.Format("20060102")

	convertedCCID, err := strconv.ParseUint(CCID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	tibp := models.TransactionInvoiceByPayer{}
	tibpReceived, err := tibp.FindInvoiceByCCIDAndDate(server.DB, convertedCCID, parsedFirstDay, parsedLastDay)
	if err != nil {
		log.Printf(err.Error())
		errList["no_posting_matrix"] = "No invoice found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedTIBP, _ := json.Marshal(tibpReceived)
	log.Printf("Get Transaction Invoice by Payer : ", string(stringifiedTIBP))
	c.JSON(http.StatusOK, gin.H{
		"response": tibpReceived,
	})

	log.Printf("End => Get Transaction Invoice By Payer")
}

func (server *Server) GetFeeInvoiceByPayer(c *gin.Context) {
	log.Printf("Begin => Get Fee Invoice By Payer")
	CCID := c.Param("id")
	month := c.Param("month")
	year := c.Param("year")

	convertedYear, _ := strconv.Atoi(year)
	convertedMonth, _ := strconv.Atoi(month)
	firstDay := time.Date(convertedYear, time.Month(convertedMonth), 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, 0).Add(time.Nanosecond * -1)

	parsedFirstDay := firstDay.Format("20060102")
	parsedLastDay := lastDay.Format("20060102")

	convertedCCID, err := strconv.ParseUint(CCID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	fibp := models.FeeInvoiceByPayer{}
	fibpReceived, err := fibp.FindFeeInvoiceByCCIDAndDate(server.DB, convertedCCID, parsedFirstDay, parsedLastDay)
	if err != nil {
		log.Printf(err.Error())
		errList["no_posting_matrix"] = "No invoice found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedFIBP, _ := json.Marshal(fibpReceived)
	log.Printf("Get Transaction Invoice by Payer : ", string(stringifiedFIBP))
	c.JSON(http.StatusOK, gin.H{
		"response": fibpReceived,
	})

	log.Printf("End => Get Fee Invoice By Payer")
}
