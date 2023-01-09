package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request method get"})
			return
		}

		c.Next()
	}
}

func PostMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request method post"})
			return
		}
		
		c.Next()
	}
}

func PutMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPut {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request method put"})
			return
		}
		
		c.Next()
	}
}

func PatchMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPatch {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request method patch"})
			return
		}
		
		c.Next()
	}
}

func DeleteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodDelete {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request method delete"})
			return
		}
		
		c.Next()
	}
}