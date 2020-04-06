package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetSites(c *gin.Context) {
	log.Printf("Begin => Get Sites")
	site := models.Site{}

	sites, err := site.FindAllSites(server.DB)
	if err != nil {
		log.Printf("No site found")
		errList["no_site"] = "No site found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	log.Printf("Successfully Get Sites")
	c.JSON(http.StatusOK, gin.H{
		"response": sites,
	})
	log.Printf("End => Get Sites")
}

func (server *Server) GetLatestSites(c *gin.Context) {
	log.Printf("Begin => Get Latest Sites")
	site := models.Site{}

	sites, err := site.FindAllLatestSites(server.DB)
	if err != nil {
		log.Printf("No latest site found")
		errList["no_site"] = "No latest site found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	log.Printf("Successfully Get Latest Sites")
	c.JSON(http.StatusOK, gin.H{
		"response": sites,
	})

	log.Printf("End => Get Latest Sites")
}

func (server *Server) GetActiveSites(c *gin.Context) {
	log.Printf("Begin => Get Latest Sites")
	site := models.Site{}

	sites, err := site.FindAllActiveSites(server.DB)
	if err != nil {
		log.Printf("No active site found")
		errList["no_site"] = "No latest site found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	log.Printf("Successfully Get Latest Sites")
	c.JSON(http.StatusOK, gin.H{
		"response": sites,
	})

	log.Printf("End => Get Active Sites")
}

func (server *Server) GetSite(c *gin.Context) {
	log.Printf("Begin => Get Site by ID")
	siteID := c.Param("id")
	convertedSiteID, err := strconv.ParseUint(siteID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	site := models.Site{}
	siteReceived, err := site.FindSiteByID(server.DB, convertedSiteID)
	if err != nil {
		errList["no_site"] = "No site found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": siteReceived,
	})
}

func (server *Server) GetSiteHistory(c *gin.Context) {
	originalSiteID := c.Param("id")
	convertedOriginalSiteID, err := strconv.ParseUint(originalSiteID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}
	site := models.Site{}

	siteReceived, err := site.FindSiteHistoryByID(server.DB, convertedOriginalSiteID)
	if err != nil {
		errList["no_site"] = "No site found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": siteReceived,
	})
}

func (server *Server) GetTerminalBySiteID(c *gin.Context) {
	siteID := c.Param("id")
	convertedsiteID, err := strconv.ParseUint(siteID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	terminal := models.Terminal{}
	terminalReceived, err := terminal.FindAllTerminalBySiteID(server.DB, convertedsiteID)
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

func (server *Server) CreateSite(c *gin.Context) {
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	site := models.Site{}
	err = json.Unmarshal(body, &site)
	if err != nil {
		fmt.Println(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	var count int
	err = server.DB.Debug().Model(models.Site{}).Where("ship_to_number = ?", site.ShipToNumber).Count(&count).Error
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if count > 0 {
		errList["ship_to_number"] = "Entered ship to number already exists"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	site.Prepare()
	errorMessages := site.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	siteCreated, err := site.CreateSite(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": siteCreated,
	})
}

func (server *Server) UpdateSite(c *gin.Context) {
	errList = map[string]string{}
	siteID := c.Param("id")

	siteid, err := strconv.ParseUint(siteID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalSite := models.Site{}
	err = server.DB.Debug().Model(models.Site{}).Where("id = ?", siteid).Order("id desc").Take(&originalSite).Error
	if err != nil {
		errList["no_site"] = "No site found"
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

	site := models.Site{}
	err = json.Unmarshal(body, &site)
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}
	site.ID = originalSite.ID

	site.Prepare()
	errorMessages := site.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	siteUpdated, err := site.UpdateSite(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": siteUpdated,
	})
}

func (server *Server) DeactivateSiteLater(c *gin.Context) {
	errList = map[string]string{}
	siteID := c.Param("id")

	siteid, err := strconv.ParseUint(siteID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalSite := models.Site{}
	err = server.DB.Debug().Unscoped().Model(models.Site{}).Where("id = ?", siteid).Order("id desc").Take(&originalSite).Error
	if err != nil {
		errList["no_site"] = "No site found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	//Count whether there is still a relation active
	dateTimeNow := time.Now()
	var activeRelationWithNullEndedCount int
	err = server.DB.Debug().Model(models.RetailerSiteRelation{}).Unscoped().Where("site_id = ? AND started_at <= ? AND ended_at IS NULL", originalSite.OriginalID, dateTimeNow).Count(&activeRelationWithNullEndedCount).Error
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	var activeRelationWithFilledEndedCount int
	err = server.DB.Debug().Model(models.RetailerSiteRelation{}).Unscoped().Where("site_id = ? AND started_at <= ? AND ended_at >= ?", originalSite.OriginalID, dateTimeNow, dateTimeNow).Count(&activeRelationWithFilledEndedCount).Error
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if activeRelationWithNullEndedCount > 0 || activeRelationWithFilledEndedCount > 0 {
		errList["linked_retailer"] = "Selected site is still linked to a retailer"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	//Count wheter there is still a terminal active
	var activeTerminalWithNullEndedCount int
	err = server.DB.Debug().Model(models.Terminal{}).Unscoped().Where("site_id = ? AND created_at <= ? AND deleted_at IS NULL", originalSite.OriginalID, dateTimeNow).Count(&activeTerminalWithNullEndedCount).Error
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if activeTerminalWithNullEndedCount > 0 {
		errList["linked_terminal"] = "Selected site is still linked to a terminal"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	var activeTerminalWithFilledEndedCount int
	err = server.DB.Debug().Model(models.Terminal{}).Unscoped().Where("site_id = ? AND created_at <= ? AND deleted_at >= ?", originalSite.OriginalID, dateTimeNow, dateTimeNow).Count(&activeTerminalWithFilledEndedCount).Error
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if activeTerminalWithNullEndedCount > 0 || activeTerminalWithFilledEndedCount > 0 {
		errList["linked_terminal"] = "Selected site is still linked to a terminal"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	//Check if the new deleted_at input is greater than the previous deleted_at
	if originalSite.DeletedAt != nil {
		dateTimeNow := time.Now()
		if dateTimeNow.After(*originalSite.DeletedAt) {
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

	site := models.Site{}
	err = json.Unmarshal(body, &site)
	if err != nil {
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}
	site.ID = originalSite.ID

	site.Prepare()
	_, err = site.DeactivateSiteLater(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": "Selected site has been deactivated successfully.",
	})
}

func (server *Server) DeactivateSiteNow(c *gin.Context) {
	errList = map[string]string{}
	siteID := c.Param("id")

	siteid, err := strconv.ParseUint(siteID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalSite := models.Site{}
	err = server.DB.Debug().Model(models.Site{}).Where("id = ?", siteid).Order("id desc").Take(&originalSite).Error
	if err != nil {
		errList["no_site"] = "No site found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	_, err = originalSite.DeactivateSiteNow(server.DB)
	if err != nil {
		errList["other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Selected site has been deactivated successfully.",
	})
}

func (server *Server) ReactivateSite(c *gin.Context) {
	errList = map[string]string{}
	siteID := c.Param("id")

	siteid, err := strconv.ParseUint(siteID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalSite := models.Site{}
	err = server.DB.Debug().Unscoped().Model(models.Site{}).Where("id = ?", siteid).Order("id desc").Take(&originalSite).Error
	if err != nil {
		errList["no_site"] = "No site found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	if originalSite.DeletedAt == nil {
		errList["status_unprocessed"] = "The site has not been deactivated"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	if originalSite.ReactivatedAt != nil {
		errList["status_unprocessed"] = "The site has been reactivated in prior"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	originalSite.Prepare()
	errorMessages := originalSite.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	siteReactivated, err := originalSite.ReactivateSite(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": siteReactivated,
	})
}

/*func (server *Server) TerminateSiteNow(c *gin.Context) {
	errList = map[string]string{}
	siteID := c.Param("id")

	siteid, err := strconv.ParseUint(siteID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalSite := models.Site{}
	err = server.DB.Debug().Model(models.Site{}).Where("id = ?", siteid).Order("id desc").Take(&originalSite).Error
	if err != nil {
		errList["no_post"] = "No site found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	_, err = originalSite.TerminateSiteNow(server.DB)
	if err != nil {
		errList["other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Selected site has been deleted successfully.",
	})
}

func (server *Server) TerminateSiteLater(c *gin.Context) {
	errList = map[string]string{}
	siteID := c.Param("id")

	siteid, err := strconv.ParseUint(siteID, 10, 64)
	if err != nil {
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalSite := models.Site{}
	err = server.DB.Debug().Model(models.Site{}).Where("id = ?", siteid).Order("id desc").Take(&originalSite).Error
	if err != nil {
		errList["no_retailer"] = "No retailer found"
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

	site := models.Site{}
	err = json.Unmarshal(body, &site)
	if err != nil {
		errList["unmarshal_error"] = "Can not unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}
	site.ID = originalSite.ID

	_, err = site.TerminateSiteLater(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Selected site will be terminated at given time.",
	})
}
*/

func (server *Server) GetSiteTypes(c *gin.Context) {
	siteType := models.SiteType{}

	siteTypes, err := siteType.FindAllSiteTypes(server.DB)
	if err != nil {
		errList["no_site_type"] = "No site type found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": siteTypes,
	})
}
