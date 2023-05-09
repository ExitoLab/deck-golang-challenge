package routes

import (
	"github.com/gin-gonic/gin"

	controller "golang-execrise/controllers"
)

// DecksRoutes function
func DecksRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("create/deck", controller.CreateNewDeckHandler())
	incomingRoutes.GET("open/deck/:deck_id", controller.FindDeckByDeckIDHandler())
	incomingRoutes.GET("draw/deck/:deck_id", controller.DrawCardDeckByDeckIDHandler())
}
