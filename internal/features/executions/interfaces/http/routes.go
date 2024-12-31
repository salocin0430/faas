package http

import (
	"faas/internal/shared/infrastructure/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupExecutionRoutes(r *gin.Engine, handler *ExecutionHandler, jwtSecret string) {
	executions := r.Group("/api/executions")
	executions.Use(middleware.ExtractUserID(jwtSecret))
	{
		executions.POST("", handler.CreateExecution)
		executions.GET("/:id", handler.GetExecution)
		executions.GET("", handler.ListExecutions)
	}
}
