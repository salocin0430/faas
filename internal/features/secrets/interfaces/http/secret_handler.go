package http

import (
	"faas/internal/features/secrets/application/dto"
	"faas/internal/features/secrets/application/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SecretHandler struct {
	secretService *service.SecretService
}

func NewSecretHandler(secretService *service.SecretService) *SecretHandler {
	return &SecretHandler{
		secretService: secretService,
	}
}

func (h *SecretHandler) CreateSecret(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")

	var req dto.CreateSecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.secretService.CreateSecret(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *SecretHandler) GetSecret(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	secretID := c.Param("id")

	response, err := h.secretService.GetSecret(c.Request.Context(), userID, secretID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SecretHandler) ListSecrets(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")

	response, err := h.secretService.ListSecrets(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"secrets": response})
}

func (h *SecretHandler) UpdateSecret(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	secretID := c.Param("id")

	var req dto.UpdateSecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.secretService.UpdateSecret(c.Request.Context(), userID, secretID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SecretHandler) DeleteSecret(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	secretID := c.Param("id")

	if err := h.secretService.DeleteSecret(c.Request.Context(), userID, secretID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
