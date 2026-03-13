package borrowing

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	borrowings := rg.Group("/borrowings")
	{
		borrowings.GET("", handler.List)
		borrowings.GET("/:id", handler.GetByID)
		borrowings.POST("", handler.Create)
		borrowings.POST("/:id/return", handler.Return)
	}
}