package response

import "github.com/gin-gonic/gin"

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	payload := gin.H{
		"message": message,
	}

	if data != nil {
		payload["data"] = data
	}

	c.JSON(statusCode, payload)
}

func Error(c *gin.Context, statusCode int, message string, errors interface{}) {
	payload := gin.H{
		"message": message,
	}

	if errors != nil {
		payload["errors"] = errors
	}

	c.JSON(statusCode, payload)
}

func PaginatedSuccess(c *gin.Context, statusCode int, message string, data interface{}, meta PaginationMeta) {
	c.JSON(statusCode, gin.H{
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}
