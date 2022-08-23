package routes

import "github.com/gin-gonic/gin"

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("user/signup", handlers.Signup())
	incomingRoutes.POST("user/login", handlers.Login())
}
