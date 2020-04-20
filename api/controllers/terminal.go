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

// Terminal Overview
func (server *Server) GetTerminalOverview(c *gin.Context) {
	log.Printf("Begin => Get Terminal Overview")
	retailerID := c.Param("id")
	siteID := c.Param("site_id")
	convertedRetailerID, err := strconv.ParseUint(retailerID, 10, 64)
	convertedSiteID, err := strconv.ParseUint(siteID, 10, 64)
	log.Printf("retailer %d site %d", convertedRetailerID, convertedSiteID)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	terminal := models.Terminal{}
	terminalReceived, err := terminal.FindTerminalOverview(server.DB, convertedRetailerID, convertedSiteID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_terminal"] = "No terminal found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedTerminalReceived, _ := json.Marshal(terminalReceived)
	log.Printf("Get Terminal Overview : ", string(stringifiedTerminalReceived))
	c.JSON(http.StatusOK, gin.H{
		"response": terminalReceived,
	})

	log.Printf("End => Get Terminal Overview")
}

func (server *Server) GetTerminals(c *gin.Context) {
	log.Printf("Begin => Get Terminals")

	terminal := models.Terminal{}
	terminals, err := terminal.FindAllTerminals(server.DB)
	if err != nil {
		log.Printf("No terminal found")
		errList["no_terminal"] = "No terminal found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedTerminalReceived, _ := json.Marshal(terminals)
	log.Printf("Get Terminals : ", string(stringifiedTerminalReceived))
	c.JSON(http.StatusOK, gin.H{
		"response": terminals,
	})

	log.Printf("End => Get Terminals")
}

func (server *Server) GetLatestTerminals(c *gin.Context) {
	terminal := models.Terminal{}

	terminals, err := terminal.FindAllLatestTerminals(server.DB)
	if err != nil {
		errList["no_terminal"] = "No terminal found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": terminals,
	})
}

func (server *Server) GetTerminal(c *gin.Context) {
	terminalID := c.Param("id")
	convertedTerminalID, err := strconv.ParseUint(terminalID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}
	terminal := models.Terminal{}

	terminalReceived, err := terminal.FindTerminalByID(server.DB, convertedTerminalID)
	if err != nil {
		errList["no_terminal"] = "No terminal found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": terminalReceived,
	})
}

func (server *Server) GetTerminalHistory(c *gin.Context) {
	originalTerminalID := c.Param("id")
	convertedOriginalTerminalID, err := strconv.ParseUint(originalTerminalID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	terminal := models.Terminal{}
	terminalReceived, err := terminal.FindTerminalHistoryByID(server.DB, convertedOriginalTerminalID)
	if err != nil {
		errList["no_terminal"] = "No terminal found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": terminalReceived,
	})
}

func (server *Server) CreateTerminal(c *gin.Context) {
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	terminal := models.Terminal{}
	err = json.Unmarshal(body, &terminal)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	var count int
	err = server.DB.Debug().Model(models.Terminal{}).Where("terminal_sn = ? AND terminal_id = ? AND merchant_id = ?", terminal.TerminalSN, terminal.TerminalID, terminal.MerchantID).Count(&count).Error
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if count > 0 {
		errList["combination"] = "The combination of entered terminal serial number, terminal id, and merchant id already exist"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	terminal.Prepare()
	errorMessages := terminal.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	terminalCreated, err := terminal.CreateTerminal(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": terminalCreated,
	})
}

func (server *Server) UpdateTerminal(c *gin.Context) {
	errList = map[string]string{}
	terminalID := c.Param("id")

	terminalid, err := strconv.ParseUint(terminalID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalTerminal := models.Terminal{}
	err = server.DB.Debug().Model(models.Terminal{}).Where("id = ?", terminalid).Order("id desc").Take(&originalTerminal).Error
	if err != nil {
		errList["no_post"] = "No terminal found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	terminal := models.Terminal{}
	err = json.Unmarshal(body, &terminal)
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	terminal.ID = originalTerminal.ID
	terminal.Prepare()
	errorMessages := terminal.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	terminalUpdated, err := terminal.UpdateTerminal(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": terminalUpdated,
	})
}

func (server *Server) DeactivateTerminalNow(c *gin.Context) {
	errList = map[string]string{}
	terminalID := c.Param("id")

	terminalid, err := strconv.ParseUint(terminalID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalTerminal := models.Terminal{}
	err = server.DB.Debug().Model(models.Terminal{}).Where("id = ?", terminalid).Order("id desc").Take(&originalTerminal).Error
	if err != nil {
		errList["no_terminal"] = "No terminal found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	_, err = originalTerminal.DeactivateTerminalNow(server.DB)
	if err != nil {
		errList["other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Selected terminal has been deactivated successfully.",
	})
}

func (server *Server) DeactivateTerminalLater(c *gin.Context) {
	errList = map[string]string{}
	terminalID := c.Param("id")

	terminalid, err := strconv.ParseUint(terminalID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalTerminal := models.Terminal{}
	err = server.DB.Debug().Model(models.Terminal{}).Unscoped().Where("id = ?", terminalid).Order("id desc").Take(&originalTerminal).Error
	if err != nil {
		errList["no_post"] = "No terminal found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	if originalTerminal.DeletedAt != nil {
		dateTimeNow := time.Now()
		if dateTimeNow.After(*originalTerminal.DeletedAt) {
			errList["time_exceeded"] = "Ended at time field can not be updated"
			c.JSON(http.StatusNotFound, gin.H{
				"error": errList,
			})
			return
		}
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	terminal := models.Terminal{}
	err = json.Unmarshal(body, &terminal)
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	terminal.ID = originalTerminal.ID
	terminal.Prepare()

	_, err = terminal.DeactivateTerminalLater(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Selected terminal has been deactivated successfully.",
	})
}

func (server *Server) ReactivateTerminal(c *gin.Context) {
	errList = map[string]string{}
	terminalID := c.Param("id")

	terminalid, err := strconv.ParseUint(terminalID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalTerminal := models.Terminal{}
	err = server.DB.Debug().Unscoped().Model(models.Terminal{}).Where("id = ?", terminalid).Order("id desc").Take(&originalTerminal).Error
	if err != nil {
		errList["no_terminal"] = "No terminal found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	if originalTerminal.DeletedAt == nil {
		errList["status_unprocessed"] = "The terminal has not been deactivated"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if originalTerminal.ReactivatedAt != nil {
		errList["status_unprocessed"] = "The terminal has been reactivated in prior"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	originalSite := models.Site{}
	err = server.DB.Debug().Raw("EXEC spAPI_Site_GetLatestByID ?", originalTerminal.SiteID).Scan(&originalSite).Error
	if err != nil {
		errList["no_related_site"] = "No related site found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	dateTimeNow := time.Now()
	//Check if the new deleted_at input is greater than the previous deleted_at
	if originalSite.DeletedAt != nil {
		if dateTimeNow.After(*originalSite.DeletedAt) {
			errList["status_unprocessed"] = "The site related to this terminal has been deactivated"
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": errList,
			})
			return
		}
	}

	originalTerminal.Prepare()
	errorMessages := originalTerminal.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	terminalReactivated, err := originalTerminal.ReactivateTerminal(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": terminalReactivated,
	})
}

/*func (server *Server) TerminateTerminalLater(c *gin.Context) {
	errList = map[string]string{}
	terminalID := c.Param("id")

	terminalid, err := strconv.ParseUint(terminalID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalTerminal := models.Terminal{}
	err = server.DB.Debug().Model(models.Terminal{}).Where("id = ?", terminalid).Order("id desc").Take(&originalTerminal).Error
	if err != nil {
		errList["no_site"] = "No terminal found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	_, err = originalTerminal.TerminateTerminalLater(server.DB)
	if err != nil {
		errList["other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Selected site will be terminated at given time.",
	})
}*/
