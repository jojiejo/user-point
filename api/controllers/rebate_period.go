package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetRebatePeriods(c *gin.Context) {
	log.Printf("Begin => Get Rebate Periods")

	rebatePeriod := models.RebatePeriod{}
	rebatePeriods, err := rebatePeriod.FindRebatePeriods(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_rebate_type"] = "No rebate type found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedPeriods, _ := json.Marshal(rebatePeriods)
	log.Printf("Get Rebate Periods : ", string(stringifiedPeriods))
	c.JSON(http.StatusOK, gin.H{
		"response": rebatePeriods,
	})

	log.Printf("End => Get Rebate Periods")
}
