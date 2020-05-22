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

func (server *Server) GetPostingMatrixTaxes(c *gin.Context) {
	log.Printf("Begin => Get Posting Matrix Taxes")

	pmt := models.PostingMatrixTax{}
	pmts, err := pmt.FindPostingMatrixTaxes(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_posting_matrix"] = "No posting matrix product found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedPmts, _ := json.Marshal(pmts)
	log.Printf("Get Posting Matrix Fees : ", string(stringifiedPmts))
	c.JSON(http.StatusOK, gin.H{
		"response": pmts,
	})

	log.Printf("End => Get Posting Matrix Taxes")
}

func (server *Server) GetPostingMatrixTax(c *gin.Context) {
	log.Printf("Begin => Get Posting Matrix by Tax")
	pmtID := c.Param("id")
	convertedPmtID, err := strconv.ParseUint(pmtID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	pmt := models.PostingMatrixTax{}
	pmtReceived, err := pmt.FindPostingMatrixTax(server.DB, convertedPmtID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_posting_matrix"] = "No posting matrix found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedPmpReceived, _ := json.Marshal(pmtReceived)
	log.Printf("Get Posting Matrix by Tax : ", string(stringifiedPmpReceived))
	c.JSON(http.StatusOK, gin.H{
		"response": pmtReceived,
	})

	log.Printf("End => Get Posting Matrix by Tax")
}

func (server *Server) CreatePostingMatrixTax(c *gin.Context) {
	log.Printf("Begin => Create Posting Matrix by Tax")
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

	pmt := models.PostingMatrixTax{}
	err = json.Unmarshal(body, &pmt)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	var count int
	err = server.DB.Debug().Model(models.PostingMatrixTax{}).Where("gl_name = ?", pmt.GLName).Count(&count).Error
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if count > 0 {
		errString := "Entered tax already exists"
		log.Printf(errString)
		errList["name"] = errString
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	pmt.Prepare()
	errorMessages := pmt.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	pmtCreated, err := pmt.CreatePostingMatrixTax(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": pmtCreated,
	})

	log.Printf("End => Create Posting Matrix by Tax")
}

func (server *Server) UpdatePostingMatrixTax(c *gin.Context) {
	log.Printf("Begin => Update Posting Matrix Tax")

	errList = map[string]string{}
	pmtID := c.Param("id")
	pmtid, err := strconv.ParseUint(pmtID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalPmt := models.PostingMatrixTax{}
	err = server.DB.Debug().Model(models.PostingMatrixTax{}).Where("id = ?", pmtid).Order("id desc").Take(&originalPmt).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_posting_matrix"] = "No posting matrix found"
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

	pmt := models.PostingMatrixTax{}
	err = json.Unmarshal(body, &pmt)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	pmt.ID = originalPmt.ID
	pmt.Prepare()
	errorMessages := pmt.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	pmtUpdated, err := pmt.UpdatePostingMatrixTax(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": pmtUpdated,
	})

	log.Printf("End => Update Posting Matrix Tax")
}
