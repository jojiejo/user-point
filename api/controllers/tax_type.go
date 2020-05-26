package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetAllTaxTypes(c *gin.Context) {
	log.Printf("Begin => Get All Tax Types")

	taxType := models.TaxType{}
	taxTypes, err := taxType.FindAllTaxTypes(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_tax_type"] = "No tax type found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedTaxTypes, _ := json.Marshal(taxTypes)
	log.Printf("Get All Tax Types : ", string(stringifiedTaxTypes))
	c.JSON(http.StatusOK, gin.H{
		"response": taxTypes,
	})

	log.Printf("End => Get All Tax Types")
}
