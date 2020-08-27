package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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
	log.Printf("Get Products : ", string(stringifiedProduct))
	c.JSON(http.StatusOK, gin.H{
		"response": products,
	})

	log.Printf("End => Get Products")
}

func (server *Server) GetActiveProducts(c *gin.Context) {
	log.Printf("Begin => Get Active Products")

	product := models.Product{}
	products, err := product.FindActiveProducts(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_product"] = "No product found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedProduct, _ := json.Marshal(products)
	log.Println("Get Active Products : ", string(stringifiedProduct))
	c.JSON(http.StatusOK, gin.H{
		"response": products,
	})

	log.Printf("End => Get Active Products")
}

func (server *Server) GetProduct(c *gin.Context) {
	log.Printf("Begin => Get Product")
	productID := c.Param("id")
	convertedProductID, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	product := models.Product{}
	productReceived, err := product.FindProductByID(server.DB, convertedProductID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_product"] = "No product found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedProduct, _ := json.Marshal(productReceived)
	log.Printf("Get Product : ", string(stringifiedReceivedProduct))
	c.JSON(http.StatusOK, gin.H{
		"response": productReceived,
	})

	log.Printf("End => Get Product")
}

func (server *Server) CreateProduct(c *gin.Context) {
	log.Printf("Begin => Create Product")
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	product := models.Product{}
	err = json.Unmarshal(body, &product)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	product.Prepare()
	errorMessages := product.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	productCreated, err := product.CreateProduct(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": productCreated,
	})

	log.Printf("End => Create Product")
}

func (server *Server) UpdateProduct(c *gin.Context) {
	log.Printf("Begin => Update Product")

	errList = map[string]string{}
	productID := c.Param("id")

	convertedProductID, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalProduct := models.Product{}
	err = server.DB.Debug().Model(models.Product{}).Where("product_id = ?", convertedProductID).Order("product_id desc").Take(&originalProduct).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_product"] = "No product found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	product := models.Product{}
	err = json.Unmarshal(body, &product)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	product.ID = originalProduct.ID
	product.Prepare()
	errorMessages := product.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	productUpdated, err := product.UpdateProduct(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": productUpdated,
	})

	log.Printf("End => Update Product")
}
