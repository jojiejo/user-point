package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

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

func (server *Server) GetLatestRetailers(c *gin.Context) {
	retailer := models.Retailer{}

	retailers, err := retailer.FindAllLatestRetailers(server.DB)
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

func (server *Server) GetActiveRetailers(c *gin.Context) {
	retailer := models.Retailer{}

	retailers, err := retailer.FindAllActiveRetailers(server.DB)
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

func (server *Server) GetRetailerHistory(c *gin.Context) {
	originalRetailerID := c.Param("id")
	convertedOriginalRetailerID, err := strconv.ParseUint(originalRetailerID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}
	retailer := models.Retailer{}

	retailerReceived, err := retailer.FindRetailerHistoryByID(server.DB, convertedOriginalRetailerID)
	if err != nil {
		errList["no_retailer"] = "No retailer found"
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

	var count int
	err = server.DB.Debug().Model(models.Retailer{}).Where("sold_to_number = ?", retailer.SoldToNumber).Count(&count).Error
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if count > 0 {
		errList["sold_to_number"] = "Entered sold to number already exists"
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

	retailerCreated, err := retailer.CreateRetailer(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"response": retailerCreated,
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

	retailerUpdated, err := retailer.UpdateRetailer(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": retailerUpdated,
	})
}

func (server *Server) DeactivateRetailerNow(c *gin.Context) {
	errList = map[string]string{}
	retailerID := c.Param("id")

	retailerid, err := strconv.ParseUint(retailerID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid Request"
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

	_, err = originalRetailer.DeactivateRetailerNow(server.DB)
	if err != nil {
		errList["other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Selected retailer has been deactivated successfully.",
	})
}

func (server *Server) DeactivateRetailerLater(c *gin.Context) {
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
	err = server.DB.Debug().Model(models.Retailer{}).Unscoped().Where("id = ?", retailerid).Order("id desc").Take(&originalRetailer).Error
	if err != nil {
		errList["no_retailer"] = "No retailer found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	//Count whether there is still a relation active
	type ActiveRelation struct {
		RetailerCount int `json:"retailer_count"`
	}

	activeRelation := ActiveRelation{}
	err = server.DB.Debug().Raw("EXEC spAPI_RetailerSiteRelation_CountActiveRetailer ?", originalRetailer.OriginalID).Scan(&activeRelation).Error
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if activeRelation.RetailerCount > 0 {
		errList["linked_retailer"] = "Selected retailer is still linked to a site"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	//Check if the new deleted_at input is greater than the previous deleted_at
	if originalRetailer.DeletedAt != nil {
		dateTimeNow := time.Now()
		if dateTimeNow.After(*originalRetailer.DeletedAt) {
			errList["time_exceeded"] = "Ended at time field can not be updated"
			c.JSON(http.StatusNotFound, gin.H{
				"error": errList,
			})
			return
		}
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
	_, err = retailer.DeactivateRetailerLater(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": "Selected retailer has been deactivated successfully.",
	})
}

func (server *Server) ReactivateRetailer(c *gin.Context) {
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
	err = server.DB.Debug().Unscoped().Model(models.Retailer{}).Where("id = ?", retailerid).Order("id desc").Take(&originalRetailer).Error
	if err != nil {
		errList["no_retailer"] = "No retailer found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	originalRetailer.Prepare()
	errorMessages := originalRetailer.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if originalRetailer.DeletedAt == nil {
		errList["status_unprocessed"] = "The retailer has not been deactivated"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if originalRetailer.ReactivatedAt != nil {
		errList["status_unprocessed"] = "The retailer has been reactivated in prior"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	retailerReactivated, err := originalRetailer.ReactivateRetailer(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": retailerReactivated,
	})
}

/*func (server *Server) TerminateRetailerNow(c *gin.Context) {
	errList = map[string]string{}
	retailerID := c.Param("id")

	retailerid, err := strconv.ParseUint(retailerID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalRetailer := models.Retailer{}
	err = server.DB.Debug().Model(models.Retailer{}).Where("id = ?", retailerid).Order("id desc").Take(&originalRetailer).Error
	if err != nil {
		errList["no_site"] = "No retailer found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	_, err = originalRetailer.TerminateRetailerNow(server.DB)
	if err != nil {
		errList["other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Selected retailer has been terminated successfully.",
	})
}

func (server *Server) TerminateRetailerLater(c *gin.Context) {
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

	_, err = retailer.TerminateRetailerLater(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Selected retailer will be terminated at given time.",
	})
}*/

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
