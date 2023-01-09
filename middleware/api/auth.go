package middleware

import (
	"fmt"
	"net/http"
	"pendekin/helper"

	"github.com/gin-gonic/gin"
)

func verifyCookie(c *gin.Context) error {
	cookie, err := c.Request.Cookie("auth_cookie")
	if err != nil {
		return fmt.Errorf("Cookie not found!")
	}

	if cookie.Value == "" {
		return fmt.Errorf("User unauthorized!")
	}

	_, err = helper.ValidateJWT(cookie.Value)
	if err != nil {
		return err
	}

	return nil
}

func AlreadyLoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := verifyCookie(c); err == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "You are already logged in!"})		
			return
		}
		
		c.Next()
	}
}

func IsLoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := verifyCookie(c); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Next()
	}
}
