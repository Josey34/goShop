package middleware

import "github.com/gin-gonic/gin"

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}
