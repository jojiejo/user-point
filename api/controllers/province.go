package controllers

import (
	"net/http"
	"strconv"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetProvinces(c *gin.Context) {
	province := models.Province{}

	provinces, err := province.FindAllProvinces(server.DB)
	if err != nil {
		errList["no_province"] = "no province found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": provinces,
	})
}

func (server *Server) GetProvince(c *gin.Context) {
	provinceID := c.Param("id")
	convertedProvinceID, err := strconv.ParseUint(provinceID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}
	province := models.Province{}

	provinceReceived, err := province.FindProvinceByID(server.DB, convertedProvinceID)
	if err != nil {
		errList["no_province"] = "No Province Found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": provinceReceived,
	})
}

func (server *Server) GetCitiesByProvinceID(c *gin.Context) {
	provinceID := c.Param("id")
	convertedProvinceID, err := strconv.ParseUint(provinceID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}
	city := models.City{}

	cityReceived, err := city.FindCityByProvinceID(server.DB, convertedProvinceID)
	if err != nil {
		errList["no_city"] = "No City Found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": cityReceived,
	})
}
