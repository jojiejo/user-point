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

func (server *Server) GetProductGroups(c *gin.Context) {
	log.Printf("Begin => Get Product Groups")

	productGroup := models.ProductGroup{}
	productGroups, err := productGroup.FindProductGroups(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_product_group"] = "No product group found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedProductGroups, _ := json.Marshal(productGroups)
	log.Printf("Get Product Groups : ", string(stringifiedProductGroups))
	c.JSON(http.StatusOK, gin.H{
		"response": productGroups,
	})

	log.Printf("End => Get Product Groups")
}

func (server *Server) GetActiveProductGroup(c *gin.Context) {
	log.Printf("Begin => Get Active Product Group")

	productGroup := models.ProductGroup{}
	productGroups, err := productGroup.FindActiveProductGroups(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_product_group"] = "No product group found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedProductGroups, _ := json.Marshal(productGroups)
	log.Printf("Get Product Groups : ", string(stringifiedProductGroups))
	c.JSON(http.StatusOK, gin.H{
		"response": productGroups,
	})

	log.Printf("End => Get Active Product Group")
}

func (server *Server) GetProductGroup(c *gin.Context) {
	log.Printf("Begin => Get Product Group")
	productGroupID := c.Param("id")
	convertedProductGroupID, err := strconv.ParseUint(productGroupID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	productGroup := models.ProductGroup{}
	receivedProductGroup, err := productGroup.FindProductGroupByID(server.DB, convertedProductGroupID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_product_group"] = "No product group found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedProductGroup, _ := json.Marshal(receivedProductGroup)
	log.Printf("Get Product : ", string(stringifiedReceivedProductGroup))
	c.JSON(http.StatusOK, gin.H{
		"response": receivedProductGroup,
	})

	log.Printf("End => Get Product Group")
}

func (server *Server) CreateProductGroup(c *gin.Context) {
	log.Printf("Begin => Create Product Group")
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

	productGroup := models.ProductGroup{}
	err = json.Unmarshal(body, &productGroup)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	productGroup.Prepare()
	errorMessages := productGroup.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	createdProductGroup, err := productGroup.CreateProductGroup(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": createdProductGroup,
	})

	log.Printf("End => Create Product Group")
}

func (server *Server) UpdateProductGroup(c *gin.Context) {
	log.Printf("Begin => Update Product Group")

	errList = map[string]string{}
	productGroupID := c.Param("id")

	convertedProductGroupID, err := strconv.ParseUint(productGroupID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalProductGroup := models.ProductGroup{}
	err = server.DB.Debug().Model(models.Product{}).Where("product_group_id = ?", convertedProductGroupID).Order("product_group_id desc").Take(&originalProductGroup).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_product_group"] = "No product group found"
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

	productGroup := models.ProductGroup{}
	err = json.Unmarshal(body, &productGroup)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	productGroup.ID = originalProductGroup.ID
	productGroup.Prepare()
	errorMessages := productGroup.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	updatedProductGroup, err := productGroup.UpdateProductGroup(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": updatedProductGroup,
	})

	log.Printf("End => Update Product Group")
}
