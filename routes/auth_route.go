package routes

import (
	"warehouse-backend/controllers"

	"github.com/gin-gonic/gin"
)

// AuthRoutes mengatur rute autentikasi
func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controllers.RegisterUser)
		auth.POST("/login", controllers.LoginUser)
	}
}
