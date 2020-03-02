package middlewares

import (
	"net/http"

	"fleethub.shell.co.id/api/auth"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	errList := make(map[string]string)
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			errList["unauthorized"] = "Unauthorized"
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errList,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// CORSMiddleware enables using React end-point
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
