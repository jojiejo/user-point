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

func (server *Server) GetBusinessTypes(c *gin.Context) {
	log.Printf("Begin => Get Business Types")

	bt := models.BusinessType{}
	bts, err := bt.FindBusinessTypes(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_business_type"] = "No business type found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedBusinessTypes, _ := json.Marshal(bts)
	log.Println("Get Business Types : ", string(stringifiedBusinessTypes))
	c.JSON(http.StatusOK, gin.H{
		"response": bts,
	})

	log.Printf("End => Get Business Types")
}

func (server *Server) GetBusinessType(c *gin.Context) {
	log.Printf("Begin => Get Business Type")
	btID := c.Param("id")
	convertedBTID, err := strconv.ParseUint(btID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	bt := models.BusinessType{}
	receivedBusinessType, err := bt.FindBusinessTypeByID(server.DB, convertedBTID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_business_type"] = "No business type found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedBusinessType, _ := json.Marshal(receivedBusinessType)
	log.Println("Get Business Type : ", string(stringifiedBusinessType))
	c.JSON(http.StatusOK, gin.H{
		"response": receivedBusinessType,
	})

	log.Printf("End => Get Business Type")
}

func (server *Server) CreateBusinessType(c *gin.Context) {
	log.Printf("Begin => Create Business Type")
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

	bt := models.BusinessType{}
	err = json.Unmarshal(body, &bt)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	bt.Prepare()
	errorMessages := bt.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	createdBusinessType, err := bt.CreateBusinessType(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": createdBusinessType,
	})

	log.Printf("End => Create Business Type")
}

func (server *Server) UpdateBusinessType(c *gin.Context) {
	log.Printf("Begin => Update Business Type")

	errList = map[string]string{}
	btID := c.Param("id")

	convertedBTID, err := strconv.ParseUint(btID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalBusinessType := models.BusinessType{}
	err = server.DB.Debug().Model(models.BusinessType{}).Where("id = ?", convertedBTID).Order("id desc").Take(&originalBusinessType).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_business_type"] = "No business type found"
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

	bt := models.BusinessType{}
	err = json.Unmarshal(body, &bt)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	bt.ID = originalBusinessType.ID
	bt.Prepare()
	errorMessages := bt.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	updatesBusinessType, err := bt.UpdateBusinessType(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": updatesBusinessType,
	})

	log.Printf("End => Update Business Type")
}
