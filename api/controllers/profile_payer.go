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

//GetProfilePayer (s) => Get Account / Payer Profile
func (server *Server) GetProfilePayer(c *gin.Context) {
	log.Printf("Begin => Get Account / Payer Profile")

	ccID := c.Param("id")
	convertedCcID, err := strconv.ParseUint(ccID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	pp := models.ProfilePayer{}
	pps, err := pp.FindProfilePayerByCCID(server.DB, convertedCcID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_account_profile"] = "No account profile found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedProfilePayer, _ := json.Marshal(pps)
	log.Println("Get Account / Payer Profile : ", string(stringifiedProfilePayer))
	c.JSON(http.StatusOK, gin.H{
		"response": pps,
	})

	log.Printf("End => Get Account / Payer Profile")
}

//CreateProfilePayer => Create Profile Payer
func (server *Server) CreateProfilePayer(c *gin.Context) {
	log.Printf("Begin => Create Profile Payer")
	errList = map[string]string{}

	ccID := c.Param("id")
	convertedCcID, err := strconv.ParseUint(ccID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalPayer := models.Payer{}
	_, err = originalPayer.FindPayerByCCID(server.DB, convertedCcID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_payer"] = "No payer found"
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

	profilePayerTrx := server.DB.Begin()

	pp := models.ProfilePayer{}
	err = json.Unmarshal(body, &pp)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	pp.ResProfile.Prepare()
	errorMessages := pp.ResProfile.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	createdProfileMaster, err := pp.ResProfile.CreateProfileMaster(server.DB)
	if err != nil {
		profilePayerTrx.Rollback()
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	pp.ResProfile.Velocity.ResProfileID = createdProfileMaster.ResProfileID
	pp.ResProfileID = createdProfileMaster.ResProfileID

	pp.CCID = originalPayer.CCID
	createdPayerProfile, err := pp.CreateProfilePayer(server.DB)
	if err != nil {
		profilePayerTrx.Rollback()
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	profilePayerTrx.Commit()
	c.JSON(http.StatusCreated, gin.H{
		"response": createdPayerProfile,
	})

	log.Printf("End => Create Profile Payer")
}

//UpdateProfilePayer => Update Profile Payer
func (server *Server) UpdateProfilePayer(c *gin.Context) {
	log.Printf("Begin => Create Profile Payer")
	errList = map[string]string{}

	ccID := c.Param("id")
	convertedCcID, err := strconv.ParseUint(ccID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalPayerProfile := models.ProfilePayer{}
	_, err = originalPayerProfile.FindProfilePayerByCCID(server.DB, convertedCcID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_account_profile"] = "No account profile found"
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

	pm := models.ProfileMaster{}
	err = json.Unmarshal(body, &pm)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	pm.ResProfileID = originalPayerProfile.ResProfileID
	pm.Velocity.ResProfileID = originalPayerProfile.ResProfileID
	pm.Prepare()
	errorMessages := pm.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	updatedProfileMaster, err := pm.UpdateProfileMaster(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"response": updatedProfileMaster,
	})

	log.Printf("End => Update Profile Payer")
}
