package mcp

import (
	"net/http"

	"github.com/mark3labs/mcp-go/server"
)

func NewHTTPHandler() http.Handler {
	return server.NewStreamableHTTPServer(
		NewServer(),
		server.WithStateLess(true),
	)
}
