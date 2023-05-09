package routes

import (
	controller "golang-execrise/controllers"

	"github.com/gin-gonic/gin"
)

// HealthRoutes function
func HealthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/health", controller.HealthCheck())
}
