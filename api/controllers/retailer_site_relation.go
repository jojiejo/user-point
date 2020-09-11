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

func (server *Server) GetRetailerSiteRelationByRetailerID(c *gin.Context) {
	retailerID := c.Param("id")
	convertedRetailerID, err := strconv.ParseUint(retailerID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}
	retailerSiteRelation := models.RetailerSiteRelation{}

	retailerSiteRelationReceived, err := retailerSiteRelation.FindAllRetailerSiteRelationByRetailerID(server.DB, convertedRetailerID)
	if err != nil {
		errList["no_relation"] = "No relation found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": retailerSiteRelationReceived,
	})
}

//GetLatestRetailerSiteRelationByRetailerID => Get Latest Retailer Site Relation By Retailer ID
func (server *Server) GetLatestRetailerSiteRelationByRetailerID(c *gin.Context) {
	retailerID := c.Param("id")
	convertedRetailerID, err := strconv.ParseUint(retailerID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	retailerSiteRelation := models.RetailerSiteRelation{}
	retailerSiteRelationReceived, err := retailerSiteRelation.FindAllLatestRetailerSiteRelationByRetailerID(server.DB, convertedRetailerID)
	if err != nil {
		errList["no_relation"] = "No relation found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": retailerSiteRelationReceived,
	})
}

func (server *Server) CreateRetailerSiteRelation(c *gin.Context) {
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}
	retailerSiteRelation := models.RetailerSiteRelation{}

	err = json.Unmarshal(body, &retailerSiteRelation)
	if err != nil {
		fmt.Println(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	retailerSiteRelation.Prepare()
	errorMessages := retailerSiteRelation.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	relationCreated, err := retailerSiteRelation.CreateRetailerSiteRelation(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": relationCreated,
	})
}

func (server *Server) UpdateRetailerSiteRelation(c *gin.Context) {
	errList = map[string]string{}
	relationID := c.Param("relation_id")

	relationid, err := strconv.ParseUint(relationID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalRelation := models.RetailerSiteRelation{}
	err = server.DB.Debug().Model(models.RetailerSiteRelation{}).Where("id = ?", relationid).Order("id desc").Take(&originalRelation).Error
	if err != nil {
		errList["no_relation"] = "No relation found"
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

	relation := models.RetailerSiteRelation{}
	err = json.Unmarshal(body, &relation)
	if err != nil {
		errList["unmarshal_error"] = "Can not unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}
	relation.ID = originalRelation.ID

	relation.Prepare()
	errorMessages := relation.UpdateValidate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	relationUpdated, err := relation.UpdateRetailerSiteRelation(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": relationUpdated,
	})
}

func (server *Server) UnlinkRetailerSiteRelation(c *gin.Context) {
	errList = map[string]string{}
	relationID := c.Param("relation_id")

	relationid, err := strconv.ParseUint(relationID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalRelation := models.RetailerSiteRelation{}
	err = server.DB.Debug().Model(models.RetailerSiteRelation{}).Where("id = ?", relationid).Order("id desc").Take(&originalRelation).Error
	if err != nil {
		errList["no_relation"] = "No relation found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	if originalRelation.EndedAt != nil {
		dateTimeNow := time.Now()
		if dateTimeNow.After(*originalRelation.EndedAt) {
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

	relation := models.RetailerSiteRelation{}
	err = json.Unmarshal(body, &relation)
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}
	relation.ID = originalRelation.ID

	_, err = relation.UnlinkRetailerSiteRelation(server.DB)
	if err != nil {
		errList["other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Selected relation has been deleted successfully.",
	})
}
