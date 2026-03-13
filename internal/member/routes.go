package member

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	members := rg.Group("/members")
	{
		members.GET("", handler.List)
		members.GET("/:id", handler.GetByID)
		members.POST("", handler.Create)
		members.PUT("/:id", handler.Update)
		members.DELETE("/:id", handler.Delete)
	}
}