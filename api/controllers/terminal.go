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

func (server *Server) GetTerminals(c *gin.Context) {
	terminal := models.Terminal{}

	terminals, err := terminal.FindAllTerminals(server.DB)
	if err != nil {
		errList["No_terminals"] = "No terminal found"
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
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}
	terminal := models.Terminal{}

	terminalReceived, err := terminal.FindTerminalByID(server.DB, convertedTerminalID)
	if err != nil {
		errList["no_terminal"] = "No terminal Found"
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
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}
	terminal := models.Terminal{}

	err = json.Unmarshal(body, &terminal)
	if err != nil {
		fmt.Println(err.Error())
		errList["Unmarshal_error"] = "Cannot unmarshal body"
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

	siteCreated, err := terminal.CreateTerminal(server.DB)
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

func (server *Server) UpdateTerminal(c *gin.Context) {
	errList = map[string]string{}
	terminalID := c.Param("id")

	terminalid, err := strconv.ParseUint(terminalID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalTerminal := models.Terminal{}
	err = server.DB.Debug().Model(models.Site{}).Where("id = ?", terminalid).Order("id desc").Take(&originalTerminal).Error
	if err != nil {
		errList["No_post"] = "No terminal found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	terminal := models.Terminal{}
	err = json.Unmarshal(body, &terminal)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
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

func (server *Server) DeleteTerminal(c *gin.Context) {
	errList = map[string]string{}
	terminalID := c.Param("id")

	terminalid, err := strconv.ParseUint(terminalID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalTerminal := models.Terminal{}
	err = server.DB.Debug().Model(models.Site{}).Where("id = ?", terminalid).Order("id desc").Take(&originalTerminal).Error
	if err != nil {
		errList["No_post"] = "No Site Found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	_, err = originalTerminal.DeleteTerminal(server.DB)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Selected terminal has been deleted successfully.",
	})

}
