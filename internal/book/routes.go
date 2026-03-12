package book

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	books := rg.Group("/books")
	{
		books.GET("", handler.List)
		books.GET("/:id", handler.GetByID)
		books.POST("", handler.Create)
		books.PUT("/:id", handler.Update)
		books.DELETE("/:id", handler.Delete)
	}
}
