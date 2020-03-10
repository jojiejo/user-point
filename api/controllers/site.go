package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetSites(c *gin.Context) {
	site := models.Site{}

	sites, err := site.FindAllSites(server.DB)
	if err != nil {
		errList["no_site"] = "No site found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": sites,
	})
}

func (server *Server) GetLatestSites(c *gin.Context) {
	site := models.Site{}

	sites, err := site.FindAllLatestSites(server.DB)
	if err != nil {
		errList["no_retailer"] = "No retailer found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": sites,
	})
}

func (server *Server) GetSite(c *gin.Context) {
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
		errList["no_post"] = "No site found"
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

func (server *Server) DeactivateSite(c *gin.Context) {
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

	_, err = originalSite.DeactivateSite(server.DB)
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
