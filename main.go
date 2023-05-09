package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	routes "golang-execrise/routes"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.New()
	router.Use(CORS())
	router.Use(gin.Logger())

	routes.HealthRoutes(router)
	routes.DecksRoutes(router)

	// // Create a new deck of cards
	// router.POST("/deck", newDeckHandler)

	// Start server
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}
