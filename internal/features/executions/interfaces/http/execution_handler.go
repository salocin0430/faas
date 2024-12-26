package http

import (
	"faas/internal/features/executions/application/dto"
	"faas/internal/features/executions/application/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExecutionHandler struct {
	executionService *service.ExecutionService
}

func NewExecutionHandler(service *service.ExecutionService) *ExecutionHandler {
	return &ExecutionHandler{executionService: service}
}

func (h *ExecutionHandler) CreateExecution(c *gin.Context) {
	var req dto.CreateExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	execution, err := h.executionService.CreateExecution(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, execution)
}

func (h *ExecutionHandler) GetExecution(c *gin.Context) {
	executionID := c.Param("id")
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	execution, err := h.executionService.GetExecution(c.Request.Context(), executionID, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to access this execution"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "execution not found"})
		return
	}

	c.JSON(http.StatusOK, execution)
}

func (h *ExecutionHandler) ListExecutions(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	executions, err := h.executionService.ListUserExecutions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, executions)
}
