package http

import (
	"faas/internal/shared/infrastructure/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupFunctionRoutes(r *gin.Engine, handler *FunctionHandler, jwtSecret string) {
	api := r.Group("/api/functions")
	api.Use(middleware.ExtractUserID(jwtSecret))
	{
		api.POST("", handler.CreateFunction)
		api.GET("", handler.ListUserFunctions)
		api.GET("/:id", handler.GetFunction)
		api.DELETE("/:id", handler.DeleteFunction)
	}
}
