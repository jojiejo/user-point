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

func (server *Server) GetSiteTerminalRelationBySiteID(c *gin.Context) {
	siteID := c.Param("id")
	convertedsiteID, err := strconv.ParseUint(siteID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}
	siteTerminalRelation := models.SiteTerminalRelation{}

	siteTerminalRelationReceived, err := siteTerminalRelation.FindAllSiteTerminalRelationBySiteID(server.DB, convertedsiteID)
	if err != nil {
		errList["no_relation"] = "No relation found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": siteTerminalRelationReceived,
	})
}

func (server *Server) CreateSiteTerminalRelation(c *gin.Context) {
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}
	siteTerminalRelation := models.SiteTerminalRelation{}

	err = json.Unmarshal(body, &siteTerminalRelation)
	if err != nil {
		fmt.Println(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	errorMessages := siteTerminalRelation.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	relationCreated, err := siteTerminalRelation.CreateSiteTerminalRelation(server.DB)
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

func (server *Server) UpdateSiteTerminalRelation(c *gin.Context) {
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

	originalRelation := models.SiteTerminalRelation{}
	err = server.DB.Debug().Model(models.SiteTerminalRelation{}).Where("id = ?", relationid).Order("id desc").Take(&originalRelation).Error
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

	relation := models.SiteTerminalRelation{}
	err = json.Unmarshal(body, &relation)
	if err != nil {
		errList["unmarshal_error"] = "Can not unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	relation.ID = originalRelation.ID
	errorMessages := relation.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	relationUpdated, err := relation.UpdateSiteTerminalRelation(server.DB)
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

func (server *Server) UnlinkSiteTerminalRelation(c *gin.Context) {
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

	originalRelation := models.SiteTerminalRelation{}
	err = server.DB.Debug().Model(models.SiteTerminalRelation{}).Where("id = ?", relationid).Order("id desc").Take(&originalRelation).Error
	if err != nil {
		errList["no_relation"] = "No relation found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	_, err = originalRelation.UnlinkSiteTerminalRelation(server.DB)
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
