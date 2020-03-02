package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetRetailers(c *gin.Context) {
	retailer := models.Retailer{}

	retailers, err := retailer.FindAllRetailers(server.DB)
	if err != nil {
		errList["no_retailer"] = "No retailer found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": retailers,
	})
}

func (server *Server) GetRetailer(c *gin.Context) {
	retailerID := c.Param("id")
	convertedRetailerID, err := strconv.ParseUint(retailerID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}
	retailer := models.Retailer{}

	retailerReceived, err := retailer.FindRetailerByID(server.DB, convertedRetailerID)
	if err != nil {
		errList["no_retailer"] = "No Retailer Found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": retailerReceived,
	})
}

func (server *Server) CreateRetailer(c *gin.Context) {
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}
	retailer := models.Retailer{}

	err = json.Unmarshal(body, &retailer)
	if err != nil {
		fmt.Println(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	retailer.Prepare()
	errorMessages := retailer.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	siteCreated, err := retailer.CreateRetailer(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": siteCreated,
	})
}

func (server *Server) UpdateRetailer(c *gin.Context) {
	errList = map[string]string{}
	retailerID := c.Param("id")

	retailerid, err := strconv.ParseUint(retailerID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalRetailer := models.Retailer{}
	err = server.DB.Debug().Model(models.Retailer{}).Where("id = ?", retailerid).Order("id desc").Take(&originalRetailer).Error
	if err != nil {
		errList["no_retailer"] = "No retailer found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	retailer := models.Retailer{}
	err = json.Unmarshal(body, &retailer)
	if err != nil {
		errList["unmarshal_error"] = "Can not unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}
	retailer.ID = originalRetailer.ID

	retailer.Prepare()
	errorMessages := retailer.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	siteUpdated, err := retailer.UpdateRetailer(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": siteUpdated,
	})
}

func (server *Server) DeleteRetailer(c *gin.Context) {
	errList = map[string]string{}
	retailerID := c.Param("id")

	retailerid, err := strconv.ParseUint(retailerID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalRetailer := models.Retailer{}
	err = server.DB.Debug().Model(models.Retailer{}).Where("id = ?", retailerid).Order("id desc").Take(&originalRetailer).Error
	if err != nil {
		errList["no_site"] = "No site found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	_, err = originalRetailer.DeleteRetailer(server.DB)
	if err != nil {
		errList["other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Selected retailer has been deleted successfully.",
	})
}

func (server *Server) GetPaymentTerms(c *gin.Context) {
	paymentTerm := models.RetailerPaymentTerm{}

	paymentTerms, err := paymentTerm.FindAllRetailerPaymentTerms(server.DB)
	if err != nil {
		errList["no_payment_term"] = "No payment term found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": paymentTerms,
	})
}

func (server *Server) GetReimbursementCycles(c *gin.Context) {
	reimbursementCycle := models.RetailerReimbursementCycle{}

	reimbursementCycles, err := reimbursementCycle.FindAllRetailerReimbursementCycles(server.DB)
	if err != nil {
		errList["no_reimbursement_cycle"] = "No reimbursement cycle found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": reimbursementCycles,
	})
}
