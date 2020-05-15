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

func (server *Server) GetRebatePrograms(c *gin.Context) {
	log.Printf("Begin => Get Rebate Programs")

	rebateProgram := models.RebateProgram{}
	rebatePrograms, err := rebateProgram.FindRebatePrograms(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_rebate_program"] = "No rebate program found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedRebatePrograms, _ := json.Marshal(rebatePrograms)
	log.Printf("Get Rebate Programs : ", string(stringifiedRebatePrograms))
	c.JSON(http.StatusOK, gin.H{
		"response": rebatePrograms,
	})

	log.Printf("End => Get Rebate Programs")
}

func (server *Server) GetRebateProgramsByTypeID(c *gin.Context) {
	log.Printf("Begin => Get Rebate Programs By Type ID")

	rtID := c.Param("id")
	convertedRtID, err := strconv.ParseUint(rtID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	rebateProgram := models.RebateProgram{}
	rebatePrograms, err := rebateProgram.FindRebateProgramsByTypeID(server.DB, convertedRtID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_rebate_program"] = "No rebate program found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedRebatePrograms, _ := json.Marshal(rebatePrograms)
	log.Printf("Get Get Rebate Programs By Type ID : ", string(stringifiedRebatePrograms))
	c.JSON(http.StatusOK, gin.H{
		"response": rebatePrograms,
	})

	log.Printf("End => Get Rebate Programs By Type ID")
}

func (server *Server) GetRebateProgram(c *gin.Context) {
	log.Printf("Begin => Get Rebate Program by ID")
	rpID := c.Param("id")
	convertedRpID, err := strconv.ParseUint(rpID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	rp := models.RebateProgram{}
	rpReceived, err := rp.FindRebateProgramByID(server.DB, convertedRpID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_rebate_program"] = "No rebate program found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedRebateProgram, _ := json.Marshal(rpReceived)
	log.Printf("Get Fees : ", string(stringifiedReceivedRebateProgram))
	c.JSON(http.StatusOK, gin.H{
		"response": rpReceived,
	})

	log.Printf("End => Get Rebate Program by ID")
}

func (server *Server) CreateRebateProgram(c *gin.Context) {
	log.Printf("Begin => Create Rebate Program")
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

	rp := models.RebateProgram{}
	err = json.Unmarshal(body, &rp)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	rp.Prepare()
	errorMessages := rp.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	rpCreated, err := rp.CreateRebateProgram(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": rpCreated,
	})

	log.Printf("End => Create Rebate Program")
}

func (server *Server) UpdateRebateProgram(c *gin.Context) {
	log.Printf("Begin => Update Rebate Program")

	errList = map[string]string{}
	rpID := c.Param("id")

	rpid, err := strconv.ParseUint(rpID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalRp := models.RebateProgram{}
	err = server.DB.Debug().Model(models.RebateProgram{}).Where("id = ?", rpid).Order("id desc").Take(&originalRp).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_rebate_program"] = "No rebate program found"
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

	rp := models.RebateProgram{}
	err = json.Unmarshal(body, &rp)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	rp.ID = originalRp.ID
	rp.Prepare()
	errorMessages := rp.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	rpUpdated, err := rp.UpdateRebateProgram(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": rpUpdated,
	})

	log.Printf("End => Update Fee")
}
