package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetRebatePayerRelations(c *gin.Context) {
	log.Printf("Begin => Get Rebate Payer Relations")

	rebatePayerRelation := models.RebatePayer{}
	rebatePayerRelations, err := rebatePayerRelation.FindRebatePayerRelations(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_rebate_payer"] = "No rebate payer relation found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedRebatePayerRelations, _ := json.Marshal(rebatePayerRelations)
	log.Printf("Get Rebate Payer Relations : ", string(stringifiedRebatePayerRelations))
	c.JSON(http.StatusOK, gin.H{
		"response": rebatePayerRelations,
	})

	log.Printf("End => Get Rebate Payer Relations")
}
