package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetGlobalVariables(c *gin.Context) {
	log.Printf("Begin => Get Global Variables")
	gv := models.GlobalVariable{}
	gvs, err := gv.FindAllGlobalVariables(server.DB)
	if err != nil {
		log.Printf("No global variable found")
		errList["no_global_variable"] = "No global variable found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	log.Printf("Successfully Get Global Variables")
	c.JSON(http.StatusOK, gin.H{
		"response": gvs,
	})
	log.Printf("End => Get Global Variables")
}

func (server *Server) GetGlobalVariableDetail(c *gin.Context) {
	gvdID := c.Param("id")
	convertedGvdID, err := strconv.ParseUint(gvdID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	gvd := models.GlobalVariableDetail{}
	gvdReceived, err := gvd.FindGlobalVariableDetailByID(server.DB, convertedGvdID)
	if err != nil {
		errList["no_global_variable_detail"] = "No global variable detail found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": gvdReceived,
	})
}

func (server *Server) GetGlobalVariableDetailByGlobalVariableID(c *gin.Context) {
	gvID := c.Param("id")
	convertedGVID, err := strconv.ParseUint(gvID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	gvds := models.GlobalVariableDetail{}
	gvdsReceived, err := gvds.FindGlobalVariableDetailByGlobalVariableID(server.DB, convertedGVID)
	if err != nil {
		errList["no_global_variable_detail"] = "No global variable detail found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": gvdsReceived,
	})
}

func (server *Server) CreateGlobalVariableDetail(c *gin.Context) {
	errList = map[string]string{}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	gvd := models.GlobalVariableDetail{}
	err = json.Unmarshal(body, &gvd)
	if err != nil {
		fmt.Println(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	var count int
	err = server.DB.Debug().Model(models.GlobalVariableDetail{}).Where("name = ?", gvd.Name).Count(&count).Error
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if count > 0 {
		errList["name"] = "Entered name already exists"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	gvd.Prepare()
	errorMessages := gvd.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	gvdCreated, err := gvd.CreateGlobalVariableDetail(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": gvdCreated,
	})
}

func (server *Server) UpdateGlobalVariableDetail(c *gin.Context) {
	errList = map[string]string{}
	gvdID := c.Param("id")

	gvdid, err := strconv.ParseUint(gvdID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalGvd := models.GlobalVariableDetail{}
	err = server.DB.Debug().Model(models.GlobalVariableDetail{}).Where("id = ?", gvdid).Order("id desc").Take(&originalGvd).Error
	if err != nil {
		errList["no_global_variable_detail"] = "No global variable detail found"
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

	gvd := models.GlobalVariableDetail{}
	err = json.Unmarshal(body, &gvd)
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	gvd.ID = originalGvd.ID
	gvd.Prepare()
	errorMessages := gvd.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	gvdReceived, err := gvd.UpdateGlobalVariableDetail(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": gvdReceived,
	})
}
