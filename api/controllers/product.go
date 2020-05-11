package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetProducts(c *gin.Context) {
	log.Printf("Begin => Get Products")

	product := models.Product{}
	products, err := product.FindProducts(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_product"] = "No product found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedProduct, _ := json.Marshal(products)
	log.Printf("Get Rebate Calculation Types : ", string(stringifiedProduct))
	c.JSON(http.StatusOK, gin.H{
		"response": products,
	})

	log.Printf("End => Get Products")
}
