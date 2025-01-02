package http

import (
	"faas/internal/shared/infrastructure/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupSecretRoutes(router *gin.Engine, handler *SecretHandler, jwtSecret string) {
	secrets := router.Group("/api/secrets")
	secrets.Use(middleware.ExtractUserID(jwtSecret))
	{
		secrets.POST("", handler.CreateSecret)
		secrets.GET("", handler.ListSecrets)
		secrets.GET("/:id", handler.GetSecret)
		secrets.PUT("/:id", handler.UpdateSecret)
		secrets.DELETE("/:id", handler.DeleteSecret)
	}
}
