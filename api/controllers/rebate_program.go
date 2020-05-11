package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetRebatePrograms(c *gin.Context) {
	log.Printf("Begin => Get Rebate Programs")

	rebateProgram := models.RebateProgram{}
	rebatePrograms, err := rebateProgram.FindRebatePrograms(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_rebate_program"] = "No rebate program found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedRebatePrograms, _ := json.Marshal(rebatePrograms)
	log.Printf("Get Rebate Programs : ", string(stringifiedRebatePrograms))
	c.JSON(http.StatusOK, gin.H{
		"response": rebatePrograms,
	})

	log.Printf("End => Get Rebate Programs")
}

func (server *Server) GetRebateProgram(c *gin.Context) {
	log.Printf("Begin => Get Rebate Program by ID")
	rpID := c.Param("id")
	convertedRpID, err := strconv.ParseUint(rpID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	rp := models.RebateProgram{}
	rpReceived, err := rp.FindRebateProgramByID(server.DB, convertedRpID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_rebate_program"] = "No rebate program found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedRebateProgram, _ := json.Marshal(rpReceived)
	log.Printf("Get Fees : ", string(stringifiedReceivedRebateProgram))
	c.JSON(http.StatusOK, gin.H{
		"response": rpReceived,
	})

	log.Printf("End => Get Rebate Program by ID")
}
