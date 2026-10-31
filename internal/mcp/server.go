package mcp

import (
	"github.com/mark3labs/mcp-go/server"
)

const (
	ServerName    = "my-geektime"
	ServerVersion = "1.0.0"
)

func NewServer() *server.MCPServer {
	s := server.NewMCPServer(
		ServerName,
		ServerVersion,
		server.WithToolCapabilities(false),
		server.WithRecovery(),
		server.WithInstructions("提供本地缓存的极客时间课程访问能力："+
			"list_courses 浏览我的课程；search_courses 全文检索课程/正文/评论；"+
			"get_course 查看课程目录；get_article 阅读文章正文；list_comments 阅读评论。"),
	)
	registerTools(s)
	return s
}
