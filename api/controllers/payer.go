package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

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

func (server *Server) UpdateInvoiceProduction(c *gin.Context) {
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
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	payer.CCID = originalPayer.CCID
	payer.Prepare()
	errorMessages := payer.ValidateInvoiceProduction()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	payerUpdated, err := payer.UpdateInvoiceProduction(server.DB)
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

func (server *Server) UpdateCredit(c *gin.Context) {
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
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	payer.CCID = originalPayer.CCID
	payer.Prepare()
	errorMessages := payer.ValidateCredit()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	payerUpdated, err := payer.UpdateCredit(server.DB)
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
