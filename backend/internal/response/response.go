package response

import "github.com/gin-gonic/gin"

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

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}
