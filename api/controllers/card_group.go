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

// Sub corporate = Sub account = branch
func (server *Server) GetCardGroupsByBranchID(c *gin.Context) {
	log.Printf("Begin => Get Card Groups by Branch")

	branchID := c.Param("id")
	convertedBranchID, err := strconv.ParseUint(branchID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	cg := models.CardGroup{}
	cgs, err := cg.FindAllCardGroupsByBranchID(server.DB, convertedBranchID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_card_group"] = "No card group found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	log.Printf("Successfully Get Card Groups")
	c.JSON(http.StatusOK, gin.H{
		"response": cgs,
	})
	log.Printf("End => Get Card Groups by Branch")
}

func (server *Server) GetCardGroupByID(c *gin.Context) {
	log.Printf("Begin => Get Specific Card Group")

	cgID := c.Param("id")
	convertedCgID, err := strconv.ParseUint(cgID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	cg := models.CardGroup{}
	cgs, err := cg.FindCardGroupByID(server.DB, convertedCgID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_card_group"] = "No card group found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	log.Printf("Successfully Get Specific Card Group")
	c.JSON(http.StatusOK, gin.H{
		"response": cgs,
	})
	log.Printf("End => Get Specific Card Group")
}

func (server *Server) CreateCardGroup(c *gin.Context) {
	log.Printf("Begin => Create Card Group")
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

	cg := models.CardGroup{}
	err = json.Unmarshal(body, &cg)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	cg.Prepare()
	errorMessages := cg.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		log.Printf("Error: ", errList)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	cgCreated, err := cg.CreateCardGroup(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	stringifiedCgCreated, _ := json.Marshal(cgCreated)
	log.Printf("New Card Group : ", string(stringifiedCgCreated))
	c.JSON(http.StatusCreated, gin.H{
		"response": cgCreated,
	})
	log.Printf("End => Create Card Group")
}

func (server *Server) UpdateCardGroup(c *gin.Context) {
	log.Printf("Begin => Update Card Group")
	errList = map[string]string{}

	cgID := c.Param("id")
	convertedCgID, err := strconv.ParseUint(cgID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalCg := models.CardGroup{}
	err = server.DB.Debug().Model(models.CardGroup{}).Where("card_group_code = ?", convertedCgID).Order("card_group_code desc").Take(&originalCg).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_card_group"] = "No card group found"
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

	cg := models.CardGroup{}
	err = json.Unmarshal(body, &cg)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	cg.CardGroupCode = originalCg.CardGroupCode
	cg.Prepare()
	errorMessages := cg.Validate()
	if len(errorMessages) > 0 {
		log.Printf(err.Error())
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	cgUpdated, err := cg.UpdateCardGroup(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	stringifiedCgUpdated, _ := json.Marshal(cgUpdated)
	log.Printf("Updated Card Group : ", string(stringifiedCgUpdated))
	c.JSON(http.StatusOK, gin.H{
		"response": cgUpdated,
	})
	log.Printf("End => Update Card Group")
}

func (server *Server) DeactivateCardGroup(c *gin.Context) {
	log.Printf("Begin => Deactivate Card Group")
	errList = map[string]string{}

	cgID := c.Param("id")
	convertedCgID, err := strconv.ParseUint(cgID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalCg := models.CardGroup{}
	err = server.DB.Debug().Model(models.CardGroup{}).Where("card_group_code = ?", convertedCgID).Order("card_group_code desc").Take(&originalCg).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_card_group"] = "No card group found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	_, err = originalCg.DeleteCardGroup(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errList,
		})
		return
	}

	log.Printf("Selected card group has been deleted successfully.")
	c.JSON(http.StatusOK, gin.H{
		"response": "Selected card group has been deleted successfully.",
	})
	log.Printf("End => Deactivate Card Group")
}
