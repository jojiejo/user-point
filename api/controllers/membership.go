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

func (server *Server) GetMemberships(c *gin.Context) {
	log.Printf("Begin => Get Memberships")

	membership := models.Membership{}
	memberships, err := membership.FindMemberships(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_card_types"] = "No card type found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedMemberships, _ := json.Marshal(memberships)
	log.Printf("Get Card Types : ", string(stringifiedMemberships))
	c.JSON(http.StatusOK, gin.H{
		"response": memberships,
	})

	log.Printf("End => Get Memberships")
}

func (server *Server) GetMembership(c *gin.Context) {
	log.Printf("Begin => Get Card Type")
	membershipID := c.Param("id")
	convertedMembershipID, err := strconv.ParseUint(membershipID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	membership := models.Membership{}
	receivedMembership, err := membership.FindMembershipByID(server.DB, convertedMembershipID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_membership"] = "No membership found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedMembership, _ := json.Marshal(receivedMembership)
	log.Printf("Get Card Type : ", string(stringifiedReceivedMembership))
	c.JSON(http.StatusOK, gin.H{
		"response": receivedMembership,
	})

	log.Printf("End => Get Card Type")
}

func (server *Server) CreateMembership(c *gin.Context) {
	log.Printf("Begin => Create Membership")
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

	membership := models.Membership{}
	err = json.Unmarshal(body, &membership)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	membership.Prepare()
	errorMessages := membership.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	createdMembership, err := membership.CreateMembership(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": createdMembership,
	})

	log.Printf("End => Create Membership")
}

func (server *Server) UpdateMembership(c *gin.Context) {
	log.Printf("Begin => Update Membership")

	errList = map[string]string{}
	membershipID := c.Param("id")

	convertedMembershipID, err := strconv.ParseUint(membershipID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalMembership := models.Membership{}
	err = server.DB.Debug().Model(models.CardType{}).Where("membership_code = ?", convertedMembershipID).Order("membership_code desc").Take(&originalMembership).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_membership"] = "No membership found"
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

	membership := models.Membership{}
	err = json.Unmarshal(body, &membership)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	membership.ID = originalMembership.ID
	membership.Prepare()
	errorMessages := membership.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	updatedMembership, err := membership.UpdateMembership(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": updatedMembership,
	})

	log.Printf("End => Update Membership")
}
