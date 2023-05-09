package controllers

import (
	"github.com/gin-gonic/gin"
)

func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"Success": "The application is running"})
	}
}
