package http

import (
	"log"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, handler *UserHandler) {
	// Middleware de debug para todas las rutas
	r.Use(func(c *gin.Context) {
		log.Printf("Request path: %s, X-User-ID: %s", c.Request.URL.Path, c.GetHeader("X-User-ID"))
		c.Next()
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
	}

	users := r.Group("/api/users")
	{
		users.GET("", handler.ListUsers)
		users.GET("/:id", handler.GetUser)
		users.DELETE("/:id", handler.DeleteUser)
	}
}
