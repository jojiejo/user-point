package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetPostingMatrixVAT(c *gin.Context) {
	log.Printf("Begin => Get Posting Matrix VAT")

	pmv := models.PostingMatrixVAT{}
	pmvs, err := pmv.FindPostingMatrixVATs(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_posting_matrix"] = "No posting matrix product found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedPmvs, _ := json.Marshal(pmvs)
	log.Printf("Get Posting Matrix VAT : ", string(stringifiedPmvs))
	c.JSON(http.StatusOK, gin.H{
		"response": pmvs,
	})

	log.Printf("End => Get Posting Matrix VAT")
}
