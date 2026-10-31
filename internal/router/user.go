package router

import (
	"github.com/gin-gonic/gin"
	v2 "github.com/zkep/my-geektime/internal/api/v2"
)

func user(_, private *gin.RouterGroup) {
	api := v2.NewUser()
	mcpToken := v2.NewMCPToken()
	{
		private.GET("/user/list", api.List)
		private.POST("/user/status", api.Status)
		private.POST("/user/create", api.Create)
		private.POST("/user/delete", api.Delete)
		private.GET("/user/mcp/tokens", mcpToken.List)
		private.POST("/user/mcp/tokens", mcpToken.Create)
		private.POST("/user/mcp/tokens/delete", mcpToken.Delete)
	}
}
