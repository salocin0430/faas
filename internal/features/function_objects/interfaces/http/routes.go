package http

import "github.com/gin-gonic/gin"

func SetupObjectRoutes(r *gin.Engine, handler *ObjectHandler) {
	objects := r.Group("/api/function-objects")
	{
		objects.POST("/:function_id/:name", handler.CreateObject)
		objects.GET("/:function_id/:name", handler.GetObject)
		objects.GET("/:function_id", handler.ListObjects)
		objects.DELETE("/:function_id/:name", handler.DeleteObject)
	}
}
