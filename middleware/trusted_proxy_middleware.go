package middleware

import (
	"github.com/gin-gonic/gin"
)

func TrustedProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sem lógica de verificação de IP. Vai direto para a próxima função.
		c.Next()
	}
}
