package middleware

import (
	"go-service/internal/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 從 cookie/session 拿 wallet
		session, _ := c.Cookie("session_id") // 看你實際名稱

		if session != "" {
			// ⚠️ 這裡你要改成你實際解析 session 的方式
			// 先假設 session 就是 wallet（先測）
			c.Set(auth.ContextWalletAddress, session)
		}

		c.Next()
	}
}