package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

//GetBranchByCCID => Get Branch By CCID
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

//GetBranch => Get Branch by Branch ID
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

//UpdateCardGroupFlagInSelectedBranch => Update Card Group Flag In Selected Branch
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

//GetTransactionInvoiceByBranch => Get Transaction Invoice By Branch
func (server *Server) GetTransactionInvoiceByBranch(c *gin.Context) {
	log.Printf("Begin => Get Transaction Invoice By Branch")
	subCorporateID := c.Param("id")
	month := c.Param("month")
	year := c.Param("year")

	convertedYear, _ := strconv.Atoi(year)
	convertedMonth, _ := strconv.Atoi(month)
	firstDay := time.Date(convertedYear, time.Month(convertedMonth), 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, 0).Add(time.Nanosecond * -1)

	parsedFirstDay := firstDay.Format("20060102")
	parsedLastDay := lastDay.Format("20060102")

	convertedSubCorporateID, err := strconv.ParseUint(subCorporateID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	tibb := models.TransactionInvoiceByBranch{}
	tibbReceived, err := tibb.FindInvoiceBySubCorporateIDAndDate(server.DB, convertedSubCorporateID, parsedFirstDay, parsedLastDay)
	if err != nil {
		log.Printf(err.Error())
		errList["no_posting_matrix"] = "No invoice found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedTIBB, _ := json.Marshal(tibbReceived)
	log.Println("Get Transaction Invoice by Payer : ", string(stringifiedTIBB))
	c.JSON(http.StatusOK, gin.H{
		"response": tibbReceived,
	})

	log.Printf("End => Get Transaction Invoice By Branch")
}

//GetFeeInvoiceByBranch => Get Fee Invoice By Branch
func (server *Server) GetFeeInvoiceByBranch(c *gin.Context) {
	log.Printf("Begin => Get Fee Invoice By Branch")
	subCorporateID := c.Param("id")
	month := c.Param("month")
	year := c.Param("year")

	convertedYear, _ := strconv.Atoi(year)
	convertedMonth, _ := strconv.Atoi(month)
	firstDay := time.Date(convertedYear, time.Month(convertedMonth), 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, 0).Add(time.Nanosecond * -1)

	parsedFirstDay := firstDay.Format("20060102")
	parsedLastDay := lastDay.Format("20060102")

	convertedSubCorporateID, err := strconv.ParseUint(subCorporateID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	fibb := models.FeeInvoiceByBranch{}
	tibbReceived, err := fibb.FindFeeInvoiceBySubCorporateIDAndDate(server.DB, convertedSubCorporateID, parsedFirstDay, parsedLastDay)
	if err != nil {
		log.Printf(err.Error())
		errList["no_posting_matrix"] = "No invoice found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedFIBB, _ := json.Marshal(tibbReceived)
	log.Println("Get Transaction Invoice by Payer : ", string(stringifiedFIBB))
	c.JSON(http.StatusOK, gin.H{
		"response": tibbReceived,
	})

	log.Printf("End => Get Fee Invoice By Branch")
}
