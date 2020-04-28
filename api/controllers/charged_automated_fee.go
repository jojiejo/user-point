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

func (server *Server) GetChargedAutomatedFeesOnSelectedAccount(c *gin.Context) {
	log.Printf("Begin => Get Charged Automated Fees on Selected Account")

	ccID := c.Param("id")
	convertedCCID, err := strconv.ParseUint(ccID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	chargedAutomatedFee := models.ChargedAutomatedFee{}
	receivedFees, err := chargedAutomatedFee.FindChargedAutomatedFeeByCCID(server.DB, convertedCCID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_fee"] = "No fee found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedFees, _ := json.Marshal(receivedFees)
	log.Printf("Get Charged Automated Fee on Selected Account : ", string(stringifiedReceivedFees))
	c.JSON(http.StatusOK, gin.H{
		"response": receivedFees,
	})

	log.Printf("End => Get Charged Automated Fees on Selected Account")
}

func (server *Server) GetChargedAutomatedFee(c *gin.Context) {
	log.Printf("Begin => Get Charged Automated Fee by ID")
	relationID := c.Param("id")
	convertedRelationID, err := strconv.ParseUint(relationID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	chargedAutomatedFee := models.ChargedAutomatedFee{}
	relationReceived, err := chargedAutomatedFee.FindChargedAutomatedFeeByID(server.DB, convertedRelationID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_fee"] = "No fee found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedRelationReceived, _ := json.Marshal(relationReceived)
	log.Printf("Get Charged Ad Hoc Fee : ", string(stringifiedRelationReceived))
	c.JSON(http.StatusOK, gin.H{
		"response": chargedAutomatedFee,
	})

	log.Printf("Begin => Get Charged Automated Fee by ID")
}

func (server *Server) UpdateAutomatedFee(c *gin.Context) {
	log.Printf("Begin => Update Charged Automated Fee by ID")

	errList = map[string]string{}
	relationID := c.Param("id")
	convertedRelationID, err := strconv.ParseUint(relationID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalRelation := models.ChargedAutomatedFee{}
	err = server.DB.Debug().Model(models.Payer{}).Where("id = ?", relationID).Order("id desc").Take(&originalRelation).Error
	if err != nil {
		errList["no_fee"] = "No fee found"
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

	relation := models.ChargedAutomatedFee{}
	err = json.Unmarshal(body, &relation)
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	relation.ID = originalRelation.ID
	relation.Prepare()
	errorMessages := relation.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	updatedRelation, err := relation.UpdateAutomatedFee(server.DB, convertedRelationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": updatedRelation,
	})

	log.Printf("End => Update Charged Automated Fee by ID")
}
