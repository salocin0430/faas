package http

import (
	"github.com/gin-gonic/gin"
)

func SetupExecutionRoutes(r *gin.Engine, handler *ExecutionHandler) {
	executions := r.Group("/api/executions")
	{
		executions.POST("", handler.CreateExecution)
		executions.GET("/:id", handler.GetExecution)
		executions.GET("", handler.ListExecutions)
	}
}
