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

func (server *Server) GetSalesReps(c *gin.Context) {
	log.Printf("Begin => Get Sales Reps")

	sr := models.SalesRep{}
	srs, err := sr.FindSalesReps(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_sales_rep"] = "No sales rep found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedSalesReps, _ := json.Marshal(srs)
	log.Printf("Get Sales Reps : ", string(stringifiedSalesReps))
	c.JSON(http.StatusOK, gin.H{
		"response": srs,
	})

	log.Printf("End => Get Sales Reps")
}

func (server *Server) GetSalesRep(c *gin.Context) {
	log.Printf("Begin => Get Sales Rep")
	srID := c.Param("id")
	convertedSrID, err := strconv.ParseUint(srID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	sr := models.SalesRep{}
	receivedSR, err := sr.FindSalesRepByID(server.DB, convertedSrID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_sales_rep"] = "No sales rep found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedSR, _ := json.Marshal(receivedSR)
	log.Printf("Get Sales Rep : ", string(stringifiedReceivedSR))
	c.JSON(http.StatusOK, gin.H{
		"response": receivedSR,
	})

	log.Printf("End => Get Sales Rep")
}

func (server *Server) CreateSalesRep(c *gin.Context) {
	log.Printf("Begin => Create Sales Rep")
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

	sr := models.SalesRep{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	sr.Prepare()
	errorMessages := sr.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	createdSr, err := sr.CreateSalesRep(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": createdSr,
	})

	log.Printf("End => Create Sales Rep")
}

func (server *Server) UpdateSalesRep(c *gin.Context) {
	log.Printf("Begin => Update Sales Rep")

	errList = map[string]string{}
	srID := c.Param("id")

	convertedSrID, err := strconv.ParseUint(srID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalSalesRep := models.SalesRep{}
	err = server.DB.Debug().Model(models.Product{}).Where("id = ?", convertedSrID).Order("product_id desc").Take(&originalSalesRep).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_sales_rep"] = "No sales rep found"
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

	sr := models.SalesRep{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	sr.ID = originalSalesRep.ID
	sr.Prepare()
	errorMessages := sr.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	updatedSalesRep, err := sr.UpdateSalesRep(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": updatedSalesRep,
	})

	log.Printf("End => Update Sales Rep")
}
