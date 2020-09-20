package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jojiejo/user-point/api/models"
)

//GetUserPointByUserID => Get User Point By User ID
func (server *Server) GetUserPointByUserID(c *gin.Context) {
	userID := c.Param("id")
	convertedUserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	userPoint := models.UserPoint{}
	userPoints, err := userPoint.FindPointHistoryByUserID(server.DB, convertedUserID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_user"] = "No user found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": userPoints,
	})
}

//UpdateUserPoint => Update user point by User ID
func (server *Server) UpdateUserPoint(c *gin.Context) {
	userID := c.Param("id")
	convertedUserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
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

	originalUser := models.User{}
	err = server.DB.Debug().
		Model(models.User{}).
		Where("id = ?", convertedUserID).
		Order("id desc").
		Take(&originalUser).
		Error
	if err != nil {
		errList["no_user"] = "No user found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	userPoint := models.UserPoint{}
	err = json.Unmarshal(body, &userPoint)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	userPoint.UserID = convertedUserID
	_, err = userPoint.CreateUserPoint(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_user"] = "No user found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	originalUser.CurrentPoint = originalUser.CurrentPoint + userPoint.Value
	updatedUser, err := originalUser.UpdateUserPoint(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_user"] = "No user found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": updatedUser,
	})
}
