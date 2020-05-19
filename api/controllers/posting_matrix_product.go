package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetPostingMatrixProducts(c *gin.Context) {
	log.Printf("Begin => Get Posting Matrix Products")

	postingMatrixProduct := models.PostingMatrixProduct{}
	postingMatrixProducts, err := postingMatrixProduct.FindPostingMatrixProducts(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_rebate_payer"] = "No posting matrix product found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedPostingMatrixProducts, _ := json.Marshal(postingMatrixProducts)
	log.Printf("Get Posting Matrix Products : ", string(stringifiedPostingMatrixProducts))
	c.JSON(http.StatusOK, gin.H{
		"response": postingMatrixProducts,
	})

	log.Printf("End => Get Posting Matrix Products")
}
