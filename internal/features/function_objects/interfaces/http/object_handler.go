package http

import (
	"faas/internal/features/function_objects/application/dto"
	"faas/internal/features/function_objects/application/service"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ObjectHandler struct {
	objectService *service.ObjectService
}

func NewObjectHandler(service *service.ObjectService) *ObjectHandler {
	return &ObjectHandler{
		objectService: service,
	}
}

func (h *ObjectHandler) CreateObject(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	req := &dto.CreateObjectRequest{
		FunctionID: c.Param("function_id"),
		Name:       c.Param("name"),
	}

	// Abrir y leer el archivo completo
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Procesar con los bytes
	response, err := h.objectService.CreateObject(c.Request.Context(), req, data, file.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *ObjectHandler) GetObject(c *gin.Context) {
	functionID := c.Param("function_id")
	objectName := c.Param("name")

	obj, data, err := h.objectService.GetObject(c.Request.Context(), functionID, objectName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "object not found"})
		return
	}

	// Configurar headers
	c.Header("Content-Type", obj.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", obj.Name))

	// Enviar archivo
	c.Data(http.StatusOK, obj.ContentType, data)
}

func (h *ObjectHandler) ListObjects(c *gin.Context) {
	functionID := c.Param("function_id")

	objects, err := h.objectService.ListObjects(c.Request.Context(), functionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, objects)
}

func (h *ObjectHandler) DeleteObject(c *gin.Context) {
	functionID := c.Param("function_id")
	objectName := c.Param("name")

	if err := h.objectService.DeleteObject(c.Request.Context(), functionID, objectName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Implementar resto de handlers...
