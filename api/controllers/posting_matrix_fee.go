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

func (server *Server) GetPostingMatrixFees(c *gin.Context) {
	log.Printf("Begin => Get Posting Matrix Fees")

	pmf := models.PostingMatrixFee{}
	pmfs, err := pmf.FindPostingMatrixFees(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_posting_matrix"] = "No posting matrix product found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedPmfs, _ := json.Marshal(pmfs)
	log.Printf("Get Posting Matrix Fees : ", string(stringifiedPmfs))
	c.JSON(http.StatusOK, gin.H{
		"response": pmfs,
	})

	log.Printf("End => Get Posting Matrix Fees")
}

func (server *Server) GetPostingMatrixFee(c *gin.Context) {
	log.Printf("Begin => Get Posting Matrix by Fee")
	pmfID := c.Param("id")
	convertedPmfID, err := strconv.ParseUint(pmfID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	pmf := models.PostingMatrixFee{}
	pmfReceived, err := pmf.FindPostingMatrixFee(server.DB, convertedPmfID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_posting_matrix"] = "No posting matrix found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedPmpReceived, _ := json.Marshal(pmfReceived)
	log.Printf("Get Posting Matrix by Tax : ", string(stringifiedPmpReceived))
	c.JSON(http.StatusOK, gin.H{
		"response": pmfReceived,
	})

	log.Printf("End => Get Posting Matrix by Tax")
}

func (server *Server) CreatePostingMatrixFee(c *gin.Context) {
	log.Printf("Begin => Create Posting Matrix by Fee")
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

	pmf := models.PostingMatrixFee{}
	err = json.Unmarshal(body, &pmf)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	var count int
	err = server.DB.Debug().Model(models.PostingMatrixFee{}).Where("fee_id = ?", pmf.FeeID).Count(&count).Error
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if count > 0 {
		errString := "Entered fee already exists"
		log.Printf(errString)
		errList["name"] = errString
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	pmf.Prepare()
	errorMessages := pmf.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	pmfCreated, err := pmf.CreatePostingMatrixFee(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": pmfCreated,
	})

	log.Printf("End => Create Posting Matrix by Fee")
}

func (server *Server) UpdatePostingMatrixFee(c *gin.Context) {
	log.Printf("Begin => Update Posting Matrix Fee")

	errList = map[string]string{}
	pmfID := c.Param("id")
	pmfid, err := strconv.ParseUint(pmfID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalPmf := models.PostingMatrixFee{}
	err = server.DB.Debug().Model(models.Fee{}).Where("id = ?", pmfid).Order("id desc").Take(&originalPmf).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_posting_matrix"] = "No posting matrix found"
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

	pmf := models.PostingMatrixFee{}
	err = json.Unmarshal(body, &pmf)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	pmf.ID = originalPmf.ID
	pmf.FeeID = originalPmf.FeeID
	pmf.Prepare()
	errorMessages := pmf.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	pmfUpdated, err := pmf.UpdatePostingMatrixFee(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": pmfUpdated,
	})

	log.Printf("End => Update Posting Matrix Fee")
}
