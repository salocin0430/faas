package http

import (
	"fmt"
	"net/http"

	"faas/internal/features/functions/application/dto"
	"faas/internal/features/functions/application/service"

	"github.com/gin-gonic/gin"
)

type FunctionHandler struct {
	functionService *service.FunctionService
}

func NewFunctionHandler(functionService *service.FunctionService) *FunctionHandler {
	return &FunctionHandler{
		functionService: functionService,
	}
}

func (h *FunctionHandler) CreateFunction(c *gin.Context) {
	var req dto.CreateFunctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener userID del header X-User-ID
	fmt.Println("Creating function for user:", c.GetHeader("X-User-ID"))
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id not found"})
		return
	}

	function, err := h.functionService.CreateFunction(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, function)
}

func (h *FunctionHandler) GetFunction(c *gin.Context) {
	functionID := c.Param("id")
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	function, err := h.functionService.GetFunction(c.Request.Context(), functionID, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to access this function"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "function not found"})
		return
	}

	c.JSON(http.StatusOK, function)
}

func (h *FunctionHandler) ListUserFunctions(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	functions, err := h.functionService.ListUserFunctions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, functions)
}

func (h *FunctionHandler) DeleteFunction(c *gin.Context) {
	functionID := c.Param("id")
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.functionService.DeleteFunction(c.Request.Context(), functionID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
