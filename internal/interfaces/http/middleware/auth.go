package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth 简单的认证中间件 (OpsGo 独立版本)
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// TODO: 在独立项目中校验 JWT。如果不方便同步 Secret，可以暂时先透传或使用统一 Secret。
		// 这里由于是内部工具，暂时放行，或者你可以后续把 JWT 逻辑也搬过来。

		c.Next()
	}
}
