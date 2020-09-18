package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

//GetCardBlockReasons => Get Card Block Reasons
func (server *Server) GetCardBlockReasons(c *gin.Context) {
	log.Printf("Begin => Get Card Block Reasons")

	cbr := models.CardBlockReason{}
	cbrs, err := cbr.FindCardBlockReasons(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_card_block_reason"] = "No card block reason class found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedCardBlockReasons, _ := json.Marshal(cbrs)
	log.Println("Get Card Block Reasons : ", string(stringifiedCardBlockReasons))
	c.JSON(http.StatusOK, gin.H{
		"response": cbrs,
	})

	log.Printf("End => Get Card Block Reasons")
}

//GetCardBlockReason => Get Card Block Reason
func (server *Server) GetCardBlockReason(c *gin.Context) {
	log.Printf("Begin => Get Card Block Reason")
	cbrID := c.Param("id")
	convertedCbrID, err := strconv.ParseUint(cbrID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	cbr := models.CardBlockReason{}
	receivedCardBlockReason, err := cbr.FindCardBlockReasonByID(server.DB, convertedCbrID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_card_block_reason"] = "No card block reason found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedCardBlockReason, _ := json.Marshal(receivedCardBlockReason)
	log.Println("Get Sales Rep : ", string(stringifiedCardBlockReason))
	c.JSON(http.StatusOK, gin.H{
		"response": receivedCardBlockReason,
	})

	log.Printf("End => Get Card Block Reason")
}
