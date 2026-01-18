package middleware

import (
	"log"
	"net/http"

	"OpsGo/internal/application/dto"

	"github.com/gin-gonic/gin"
)

// Recovery 恢复中间件，捕获panic
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, dto.Error(500, "Internal server error"))
				c.Abort()
			}
		}()
		c.Next()
	}
}
