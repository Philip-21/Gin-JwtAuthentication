package routes

import (
	"go-jwt/handlers"
	"go-jwt/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	///the user will have a token before using this route
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", handlers.GetUsers())
	incomingRoutes.GET("/users/:user_id", handlers.GetUser())
}
