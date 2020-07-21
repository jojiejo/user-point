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

func (server *Server) GetIndustryClassifications(c *gin.Context) {
	log.Printf("Begin => Get Industry Classifications")

	industryClassification := models.IndustryClassification{}
	industryClassifications, err := industryClassification.FindIndustryClassifications(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_industry_classification"] = "No industry classification found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedIndustryClassifications, _ := json.Marshal(industryClassifications)
	log.Printf("Get Industry Classifications : ", string(stringifiedIndustryClassifications))
	c.JSON(http.StatusOK, gin.H{
		"response": industryClassifications,
	})

	log.Printf("End => Get Industry Classifications")
}

func (server *Server) GetIndustryClassification(c *gin.Context) {
	log.Printf("Begin => Get Industry Classification")
	industryClassificationID := c.Param("id")
	convertedIndustryClassificationID, err := strconv.ParseUint(industryClassificationID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	industryClassification := models.IndustryClassification{}
	receivedIndustryClassification, err := industryClassification.FindIndustryClassification(server.DB, convertedIndustryClassificationID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_industry_classification"] = "No industry classification found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedIndustryClassification, _ := json.Marshal(receivedIndustryClassification)
	log.Printf("Get Industry Classification : ", string(stringifiedReceivedIndustryClassification))
	c.JSON(http.StatusOK, gin.H{
		"response": receivedIndustryClassification,
	})

	log.Printf("End => Get Industry Classification")
}

func (server *Server) CreateIndustryClassification(c *gin.Context) {
	log.Printf("Begin => Create Industry Classification")
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

	industryClassification := models.IndustryClassification{}
	err = json.Unmarshal(body, &industryClassification)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	industryClassification.Prepare()
	errorMessages := industryClassification.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	createdIndustryClassification, err := industryClassification.CreateIndustryClassification(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": createdIndustryClassification,
	})

	log.Printf("End => Create Industry Classification")
}

func (server *Server) UpdateIndustryClassification(c *gin.Context) {
	log.Printf("Begin => Update Industry Classification")

	errList = map[string]string{}
	industryClassificationID := c.Param("id")

	convertedIndustryClassificationID, err := strconv.ParseUint(industryClassificationID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalIndustryClassification := models.IndustryClassification{}
	err = server.DB.Debug().Model(models.CardType{}).Where("id = ?", convertedIndustryClassificationID).Order("id desc").Take(&originalIndustryClassification).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_industry_classification"] = "No industry classification found"
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

	industryClassification := models.IndustryClassification{}
	err = json.Unmarshal(body, &industryClassification)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	industryClassification.ID = originalIndustryClassification.ID
	industryClassification.Prepare()
	errorMessages := industryClassification.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	updatedIndustryClassification, err := industryClassification.UpdateIndustryClassification(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": updatedIndustryClassification,
	})

	log.Printf("End => Update Industry Classification")
}
