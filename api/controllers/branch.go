package controllers

import (
	"log"
	"net/http"
	"strconv"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetBranchByCCID(c *gin.Context) {
	CCID := c.Param("id")
	convertedCCID, err := strconv.ParseUint(CCID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	branch := models.ShortenedBranch{}
	branchReceived, err := branch.FindBranchByCCID(server.DB, convertedCCID)
	if err != nil {
		errList["no_branch"] = "No branch found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": branchReceived,
	})
}

func (server *Server) GetBranch(c *gin.Context) {
	log.Printf("Begin => Get Specific Branch")
	branchID := c.Param("id")
	convertedBranchID, err := strconv.ParseUint(branchID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	branch := models.Branch{}
	branchReceived, err := branch.FindBranchByID(server.DB, convertedBranchID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_payer"] = "No payer found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	log.Printf("Successfully Get Specific Branch")
	c.JSON(http.StatusOK, gin.H{
		"response": branchReceived,
	})
}
