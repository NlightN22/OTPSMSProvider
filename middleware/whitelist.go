package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type WhitelistMiddleware struct {
	allowed []string
}

func NewWhitelistMiddleware(allowed []string) *WhitelistMiddleware {
	return &WhitelistMiddleware{allowed: allowed}
}

func (m *WhitelistMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if parsed := net.ParseIP(ip); parsed != nil && parsed.IsLoopback() {
			c.Next()
			return
		}
		if len(m.allowed) == 0 {
			c.Next()
			return
		}
		for _, w := range m.allowed {
			if strings.TrimSpace(w) == ip {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "IP not allowed"})
	}
}
