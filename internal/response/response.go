package response

import "github.com/gin-gonic/gin"

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
