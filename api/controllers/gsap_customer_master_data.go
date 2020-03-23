package controllers

import (
	"net/http"
	"strconv"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetCustomerMasterData(c *gin.Context) {
	MCMSID := c.Param("id")
	convertedMCMSID, err := strconv.ParseUint(MCMSID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	gsapCustomerMasterDatum := models.GSAPCustomerMasterData{}
	dataReceived, err := gsapCustomerMasterDatum.FindDataByMCMSID(server.DB, convertedMCMSID)
	if err != nil {
		errList["no_customer"] = "No customer found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": dataReceived,
	})
}
