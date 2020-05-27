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

func (server *Server) GetPostingMatrixProducts(c *gin.Context) {
	log.Printf("Begin => Get Posting Matrix Products")

	postingMatrixProduct := models.PostingMatrixProduct{}
	postingMatrixProducts, err := postingMatrixProduct.FindPostingMatrixProducts(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_posting_matrix"] = "No posting matrix product found"
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

func (server *Server) GetPostingMatrixProduct(c *gin.Context) {
	log.Printf("Begin => Get Posting Matrix by Product")
	pmpID := c.Param("id")
	convertedPmpID, err := strconv.ParseUint(pmpID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	pmp := models.PostingMatrixProduct{}
	pmpReceived, err := pmp.FindPostingMatrixProduct(server.DB, convertedPmpID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_posting_matrix"] = "No posting matrix found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedPmpReceived, _ := json.Marshal(pmpReceived)
	log.Printf("Get Posting Matrix by Product : ", string(stringifiedPmpReceived))
	c.JSON(http.StatusOK, gin.H{
		"response": pmpReceived,
	})

	log.Printf("End => Get Posting Matrix by Product")
}

func (server *Server) CreatePostingMatrixProduct(c *gin.Context) {
	log.Printf("Begin => Create Posting Matrix by Product")
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

	pmp := models.PostingMatrixProduct{}
	err = json.Unmarshal(body, &pmp)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	var count int
	err = server.DB.Debug().Model(models.PostingMatrixProduct{}).Where("product_id = ?", pmp.ProductID).Count(&count).Error
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if count > 0 {
		errString := "Entered product already exists"
		log.Printf(errString)
		errList["name"] = errString
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	pmp.Prepare()
	errorMessages := pmp.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	pmpCreated, err := pmp.CreatePostingMatrixProduct(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": pmpCreated,
	})

	log.Printf("End => Create Posting Matrix by Product")
}

func (server *Server) UpdatePostingMatrixProduct(c *gin.Context) {
	log.Printf("Begin => Update Posting Matrix Product")

	errList = map[string]string{}
	pmpID := c.Param("id")
	pmpid, err := strconv.ParseUint(pmpID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalPmp := models.PostingMatrixProduct{}
	err = server.DB.Debug().Model(models.Fee{}).Where("id = ?", pmpid).Order("id desc").Take(&originalPmp).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_posting_matrix"] = "No posting matrix found"
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

	pmp := models.PostingMatrixProduct{}
	err = json.Unmarshal(body, &pmp)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	pmp.ID = originalPmp.ID
	pmp.ProductID = originalPmp.ProductID
	pmp.Prepare()
	errorMessages := pmp.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	pmpUpdated, err := pmp.UpdatePostingMatrixProduct(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": pmpUpdated,
	})

	log.Printf("End => Update Posting Matrix Product")
}