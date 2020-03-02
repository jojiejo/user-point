package controllers

import (
	"net/http"
	"strconv"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetCities(c *gin.Context) {
	city := models.City{}

	cities, err := city.FindAllCities(server.DB)
	if err != nil {
		errList["no_city"] = "No city found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": cities,
	})
}

func (server *Server) GetCity(c *gin.Context) {
	cityID := c.Param("id")
	convertedCityID, err := strconv.ParseUint(cityID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}
	city := models.City{}

	cityReceived, err := city.FindCityByID(server.DB, convertedCityID)
	if err != nil {
		errList["no_province"] = "No province found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": cityReceived,
	})
}
