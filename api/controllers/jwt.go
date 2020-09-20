package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jojiejo/user-point/api/auth"
)

//GetJWT => GenerateJWT
func (server *Server) GetJWT(c *gin.Context) {
	token, err := auth.CreateToken()
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": token,
	})
}
