package main

import (
	"go-jwt/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	const port = "8080"
	//port = os.Getenv("PORT")

	router := gin.New()
	// a Logger middleware that will write the logs to gin.DefaultWriter.
	// By default, gin.DefaultWriter = os.Stdout.
	router.Use(gin.Logger())
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access Granted for api 1"}) //setting the headers
	})
	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"sucess": "Aces Granted for api-2"})
	})
	router.Run(port)
}
