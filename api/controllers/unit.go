package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetAllUnits(c *gin.Context) {
	log.Printf("Begin => Get All Units")

	unit := models.Unit{}
	units, err := unit.FindAllUnits(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_unit"] = "No unit found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedUnits, _ := json.Marshal(units)
	log.Printf("Get All Card Status : ", string(stringifiedReceivedUnits))
	c.JSON(http.StatusOK, gin.H{
		"response": units,
	})

	log.Printf("End => Get All Units")
}
