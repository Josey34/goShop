package middleware

import (
	"errors"
	"net/http"

	domainerrors "github.com/Josey34/goshop/domain/errors"
	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var domainErr *domainerrors.DomainError
		if errors.As(err, &domainErr) {
			switch domainErr.Code {
			case domainerrors.CodeNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": domainErr.Message})
			case domainerrors.CodeValidation:
				c.JSON(http.StatusBadRequest, gin.H{"error": domainErr.Message, "fields": domainErr.Fields})
			case domainerrors.CodeConflict:
				c.JSON(http.StatusConflict, gin.H{"error": domainErr.Message})
			case domainerrors.CodeUnauthorized:
				c.JSON(http.StatusUnauthorized, gin.H{"error": domainErr.Message})
			case domainerrors.CodeInsufficientStock:
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": domainErr.Message})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
	}
}
