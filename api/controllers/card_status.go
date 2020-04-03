package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetAllCardStatus(c *gin.Context) {
	log.Printf("Begin => Get All Card Status")

	cs := models.CardStatus{}
	css, err := cs.FindAllCardStatus(server.DB)
	if err != nil {
		errString := "No card status found"
		log.Printf(errString)
		errList["no_card_status"] = errString
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedCS, _ := json.Marshal(css)
	log.Printf("Get All Card Status : ", string(stringifiedReceivedCS))
	c.JSON(http.StatusOK, gin.H{
		"response": css,
	})

	log.Printf("End => Get All Card Status")
}
