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

func (server *Server) GetAccountClasses(c *gin.Context) {
	log.Printf("Begin => Get Account Classes")

	ac := models.AccountClass{}
	acs, err := ac.FindAccountClasses(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_account_class"] = "No account class found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedAccountClasses, _ := json.Marshal(acs)
	log.Printf("Get Sales Reps : ", string(stringifiedAccountClasses))
	c.JSON(http.StatusOK, gin.H{
		"response": acs,
	})

	log.Printf("End => Get Account Classes")
}

func (server *Server) GetAccountClass(c *gin.Context) {
	log.Printf("Begin => Get Account Class")
	acID := c.Param("id")
	convertedACID, err := strconv.ParseUint(acID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	ac := models.AccountClass{}
	receivedAccountClass, err := ac.FindAccountClassByID(server.DB, convertedACID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_account_class"] = "No account class found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedAccountClass, _ := json.Marshal(receivedAccountClass)
	log.Printf("Get Sales Rep : ", string(stringifiedAccountClass))
	c.JSON(http.StatusOK, gin.H{
		"response": receivedAccountClass,
	})

	log.Printf("End => Get Account Class")
}

func (server *Server) CreateAccountClass(c *gin.Context) {
	log.Printf("Begin => Create Account Class")
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

	ac := models.AccountClass{}
	err = json.Unmarshal(body, &ac)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	ac.Prepare()
	errorMessages := ac.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	createdAccountClass, err := ac.CreateAccountClass(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": createdAccountClass,
	})

	log.Printf("End => Create Account Class")
}

func (server *Server) UpdateAccountClass(c *gin.Context) {
	log.Printf("Begin => Update Account Class")

	errList = map[string]string{}
	acID := c.Param("id")

	convertedACID, err := strconv.ParseUint(acID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalAccountClass := models.AccountClass{}
	err = server.DB.Debug().Model(models.AccountClass{}).Where("id = ?", convertedACID).Order("id desc").Take(&originalAccountClass).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_account_class"] = "No account class found"
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

	ac := models.AccountClass{}
	err = json.Unmarshal(body, &ac)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	ac.ID = originalAccountClass.ID
	ac.Prepare()
	errorMessages := ac.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	updatedAccountClass, err := ac.UpdateAccountClass(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": updatedAccountClass,
	})

	log.Printf("End => Update Account Class")
}
