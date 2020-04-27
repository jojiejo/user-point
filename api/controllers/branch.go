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

func (server *Server) GetBranchByCCID(c *gin.Context) {
	CCID := c.Param("id")
	convertedCCID, err := strconv.ParseUint(CCID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	branch := models.ShortenedBranch{}
	branchReceived, err := branch.FindBranchByCCID(server.DB, convertedCCID)
	if err != nil {
		errList["no_branch"] = "No branch found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": branchReceived,
	})
}

func (server *Server) GetBranch(c *gin.Context) {
	log.Printf("Begin => Get Specific Branch")
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

	branch := models.Branch{}
	branchReceived, err := branch.FindBranchByID(server.DB, convertedBranchID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_payer"] = "No payer found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	log.Printf("Successfully Get Specific Branch")
	c.JSON(http.StatusOK, gin.H{
		"response": branchReceived,
	})
}

func (server *Server) UpdateCardGroupFlagInSelectedBranch(c *gin.Context) {
	log.Printf("Begin => Update Card Group Flag In Selected Branch")
	errList = map[string]string{}
	branchID := c.Param("id")

	branchid, err := strconv.ParseUint(branchID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalBranch := models.Branch{}
	_, err = originalBranch.FindBranchByID(server.DB, branchid)
	if err != nil {
		log.Printf(err.Error())
		errList["no_fee"] = "No branch found"
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

	branch := models.Branch{}
	err = json.Unmarshal(body, &branch)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	branch.SubCorporateID = originalBranch.SubCorporateID
	errorMessages := branch.ValidateConfiguration()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	branchUpdated, err := branch.UpdateConfiguration(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	var numberOfCardGroup int
	err = server.DB.Debug().Model(models.CardGroup{}).Where("sub_corporate_id = ?", originalBranch.SubCorporateID).Count(&numberOfCardGroup).Error
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if numberOfCardGroup == 0 {
		cg := models.CardGroup{}
		cg.SubCorporateID = originalBranch.SubCorporateID
		cg.ResProfileID = branch.ResProfileID
		cg.CardGroupName = "CG " + originalBranch.GSAPCustomerMasterData.ContactName_1
		_, err := cg.CreateCardGroup(server.DB)
		if err != nil {
			log.Printf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"response": branchUpdated,
	})

	log.Printf("End => Update Card Group Flag In Selected Branch")
}
