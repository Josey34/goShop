package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"message": message,
	})
}

func ErrorWithFields(c *gin.Context, statusCode int, message string, fields map[string]string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"message": message,
		"fields":  fields,
	})
}

func Paginated(c *gin.Context, data interface{}, page, limit, total int) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"meta": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
