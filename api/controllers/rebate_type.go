package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetRebateTypes(c *gin.Context) {
	log.Printf("Begin => Get Rebate Types")

	rebateType := models.RebateType{}
	rebateTypes, err := rebateType.FindRebateTypes(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_rebate_type"] = "No rebate type found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedTypes, _ := json.Marshal(rebateTypes)
	log.Printf("Get Rebate Types : ", string(stringifiedTypes))
	c.JSON(http.StatusOK, gin.H{
		"response": rebateTypes,
	})

	log.Printf("End => Get Rebate Types")
}
