package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetAllFakturPajakRange(c *gin.Context) {
	log.Printf("Begin => Get All Faktur Pajak Range")

	fpr := models.FakturPajakRange{}
	fprs, err := fpr.FindFakturPajakRanges(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_faktur_pajak"] = "No faktur pajak range found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedFprs, _ := json.Marshal(fprs)
	log.Printf("Get All Faktur Pajak Range : ", string(stringifiedFprs))
	c.JSON(http.StatusOK, gin.H{
		"response": fprs,
	})

	log.Printf("End => Get All Faktur Pajak Range")
}

func (server *Server) GetNextAvailableFakturPajakNumber(c *gin.Context) {
	log.Printf("Begin => Get Next Available Faktur Pajak Number")

	fpr := models.FakturPajakRange{}
	_, err := fpr.FindFakturPajakNextAvailableNumber(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_faktur_pajak"] = "No faktur pajak range found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedFpr, _ := json.Marshal(fpr)
	log.Printf("Get Get Next Available Faktur Pajak Number : ", string(stringifiedFpr))
	c.JSON(http.StatusOK, gin.H{
		"response": fpr,
	})

	log.Printf("End => Get Next Available Faktur Pajak Number")
}

func (server *Server) GetAvailableFakturPajakRange(c *gin.Context) {
	log.Printf("Begin => Get Available Faktur Pajak Range")

	afpr := models.AvailableFakturPajakRange{}
	_, err := afpr.FindAvailableFakturPajakRange(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_faktur_pajak_range"] = "No faktur pajak range found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedAfpr, _ := json.Marshal(afpr)
	log.Printf("Get Available Faktur Pajak Range : ", string(stringifiedAfpr))
	c.JSON(http.StatusOK, gin.H{
		"response": afpr,
	})

	log.Printf("End => Get Next Available Faktur Pajak Range")
}

func (server *Server) CreateFakturPajakRange(c *gin.Context) {
	log.Printf("Begin => Create Faktur Pajak Range")
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

	fpr := models.FakturPajakRange{}
	err = json.Unmarshal(body, &fpr)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	var count int
	err = server.DB.Debug().Model(models.FakturPajakRange{}).Where("prefix = ? AND start_range = ? AND end_range = ?", fpr.Prefix, fpr.StartRange, fpr.EndRange).Count(&count).Error
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if count > 0 {
		errString := "Entered faktur pajak range already exists"
		log.Printf(errString)
		errList["faktur_pajak_range"] = errString
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	fpr.Prepare()
	errorMessages := fpr.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	fprCreated, err := fpr.CreateFakturPajakRange(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": fprCreated,
	})

	log.Printf("End => Create Faktur Pajak Range")
}
