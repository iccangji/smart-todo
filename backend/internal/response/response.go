package response

import (
	"github.com/gin-gonic/gin"
)

func SuccessWithData(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{
		"success": true,
		"data":    data,
	})
}

func SuccessWithMessage(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": true,
		"message": message,
	})
}

func SucessWithPagination(c *gin.Context, status int, response PaginatedResponse) {
	c.JSON(status, gin.H{
		"success": true,
		"data":    response.Data,
		"meta":    response.Meta,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}
