package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetPayerAssociations(c *gin.Context) {
	log.Printf("Begin => Get Payer Associations")

	pas := models.PayerAssociation{}
	payers, err := pas.FindPayerAssociations(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_payer"] = "No payer association found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": payers,
	})

	log.Printf("End => Get Payer Associations")
}

func (server *Server) GetPayerByPayerAssociationID(c *gin.Context) {
	log.Printf("Begin => Get Payer by Payer Association ID")

	paID := c.Param("id")
	payer := models.ShortenedPayer{}
	payerReceived, err := payer.FindPayerByPayerAssociationNumber(server.DB, paID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_payer_association"] = "No payer association found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedPayer, _ := json.Marshal(payerReceived)
	log.Printf("Get Payers : ", string(stringifiedReceivedPayer))
	c.JSON(http.StatusOK, gin.H{
		"response": payerReceived,
	})

	log.Printf("End => Get Fee by ID")
}
