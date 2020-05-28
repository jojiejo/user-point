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

func (server *Server) GetAllTaxes(c *gin.Context) {
	log.Printf("Begin => Get All Taxes")

	tax := models.Tax{}
	taxes, err := tax.FindAllTaxes(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_tax"] = "No tax found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedTaxes, _ := json.Marshal(taxes)
	log.Printf("Get All Taxes : ", string(stringifiedTaxes))
	c.JSON(http.StatusOK, gin.H{
		"response": taxes,
	})

	log.Printf("End => Get All Taxes")
}

func (server *Server) GetTax(c *gin.Context) {
	log.Printf("Begin => Get Tax")
	taxID := c.Param("id")
	convertedTaxID, err := strconv.ParseUint(taxID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	tax := models.Tax{}
	taxReceived, err := tax.FindTax(server.DB, convertedTaxID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_tax"] = "No tax found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedTax, _ := json.Marshal(taxReceived)
	log.Printf("Get Tax : ", string(stringifiedTax))
	c.JSON(http.StatusOK, gin.H{
		"response": taxReceived,
	})

	log.Printf("End => Get Tax")
}

func (server *Server) CreateTax(c *gin.Context) {
	log.Printf("Begin => Create Tax")
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

	tax := models.Tax{}
	err = json.Unmarshal(body, &tax)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	var count int
	err = server.DB.Debug().Model(models.Tax{}).Where("tax_type_id = ? AND province_id = ? AND (ended_at IS NULL OR ended_at > ?)", tax.TaxTypeID, tax.ProvinceID, tax.StartedAt).Count(&count).Error
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if count > 0 {
		err = tax.UpdateRelatedTax(server.DB)
		if err != nil {
			log.Printf(err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err,
			})
			return
		}
	}

	tax.Prepare()
	errorMessages := tax.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	taxCreated, err := tax.CreateTax(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": taxCreated,
	})

	log.Printf("End => Create Tax")
}

func (server *Server) UpdateTax(c *gin.Context) {
	log.Printf("Begin => Update Tax")

	errList = map[string]string{}
	taxID := c.Param("id")
	taxid, err := strconv.ParseUint(taxID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalTax := models.Tax{}
	err = server.DB.Debug().Model(models.Fee{}).Where("id = ?", taxid).Order("id desc").Take(&originalTax).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_tax"] = "No tax found"
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

	tax := models.Tax{}
	err = json.Unmarshal(body, &tax)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	tax.ID = originalTax.ID
	tax.Prepare()
	taxUpdated, err := tax.UpdateTax(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": taxUpdated,
	})

	log.Printf("End => Update Tax")
}
