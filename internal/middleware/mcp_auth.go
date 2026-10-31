package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/zkep/my-geektime/internal/global"
	"github.com/zkep/my-geektime/internal/mcptoken"
)

func MCPAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.JWT == nil {
			abortMCPAuth(c)
			return
		}
		cfg := *global.JWT
		cfg.TokenLookup = "header:Authorization,cookie:analogjwt,query:token"
		raw, err := TokenFromRequest(&cfg, c)
		if err != nil {
			abortMCPAuth(c)
			return
		}
		if token, err1 := cfg.ParseToken(raw); err1 == nil {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if exp, ok1 := claims["exp"].(float64); ok1 && int64(exp) >= time.Now().UTC().Unix() {
					c.Set(global.Identity, claims[global.Identity])
					c.Set(global.Role, claims[global.Role])
					c.Next()
					return
				}
			}
		}
		uid, roleId, err1 := mcptoken.Authenticate(c.Request.Context(), raw)
		if err1 != nil {
			abortMCPAuth(c)
			return
		}
		c.Set(global.Identity, uid)
		c.Set(global.Role, float64(roleId))
		c.Next()
	}
}

func abortMCPAuth(c *gin.Context) {
	c.Header("WWW-Authenticate", `Bearer realm="my-geektime"`)
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"jsonrpc": "2.0",
		"id":      nil,
		"error": gin.H{
			"code":    -32001,
			"message": "unauthorized: use Authorization: Bearer <token>, the analogjwt cookie, or ?token=<token>; token can be a JWT or a long-lived MCP token",
		},
	})
}
