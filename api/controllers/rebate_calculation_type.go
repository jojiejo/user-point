package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetRebateCalculationTypes(c *gin.Context) {
	log.Printf("Begin => Get Rebate Calculation Types")

	rebateCalculationType := models.RebateCalculationType{}
	rebateCalculationTypes, err := rebateCalculationType.FindRebateCalculationTypes(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_rebate_calculation_type"] = "No rebate calculation type found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedCalculationTypes, _ := json.Marshal(rebateCalculationTypes)
	log.Printf("Get Rebate Calculation Types : ", string(stringifiedCalculationTypes))
	c.JSON(http.StatusOK, gin.H{
		"response": rebateCalculationTypes,
	})

	log.Printf("End => Get Rebate Calculation Types")
}
