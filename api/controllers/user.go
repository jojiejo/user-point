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

//GetUsers => Get Users
func (server *Server) GetUsers(c *gin.Context) {
	user := models.User{}
	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_user"] = "No user found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": users,
	})
}

//GetUserByID => Get User By ID
func (server *Server) GetUserByID(c *gin.Context) {
	ID := c.Param("id")
	convertedID, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	user := models.User{}
	receivedUser, err := user.FindUserByID(server.DB, convertedID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_user"] = "No user found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": receivedUser,
	})
}

//CreateUser => Create User
func (server *Server) CreateUser(c *gin.Context) {
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

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	user.Prepare()
	errorMessages := user.ValidateInsertion()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	createdUser, err := user.CreateUser(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"response": createdUser,
	})
}

//DeleteUser => Delete User
func (server *Server) DeleteUser(c *gin.Context) {
	errList = map[string]string{}

	ID := c.Param("id")
	convertedID, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalUser := models.User{}
	err = server.DB.Debug().
		Model(models.User{}).
		Where("id = ?", convertedID).
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

	_, err = originalUser.DeleteUser(server.DB)
	if err != nil {
		errList["other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Selected user has been deleted successfully.",
	})
}
