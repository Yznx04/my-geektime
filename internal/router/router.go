package router

import (
	"embed"
	"net/http"
	"path"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/zkep/my-geektime/internal/config"
	"github.com/zkep/my-geektime/internal/global"
	mcpapi "github.com/zkep/my-geektime/internal/mcp"
	mw "github.com/zkep/my-geektime/internal/middleware"
)

func NewRouter(assets embed.FS) (*gin.Engine, error) {
	e := gin.Default()

	e.Use(mw.Cors())

	ef, err := static.EmbedFolder(assets, "web")
	if err != nil {
		return e, err
	}

	e.Use(static.Serve("/", ef))

	e.NoRoute(func(c *gin.Context) {
		reqPath := c.Request.URL.Path
		lastSlash := strings.LastIndex(reqPath, "/")
		fileName := reqPath[lastSlash+1:]
		hasExt := strings.Contains(fileName, ".")
		if !hasExt {
			c.FileFromFS("/", ef)
			return
		}
		c.Status(http.StatusNotFound)
	})

	if global.CONF.Storage.Driver == config.StorageLocal {
		e.StaticFS(path.Join("/", global.CONF.Storage.Bucket),
			gin.Dir(global.CONF.Storage.Directory, true))
	}

	e.Any("/mcp", mw.MCPAuth(), gin.WrapH(mcpapi.NewHTTPHandler()))

	public := e.Group("v2", mw.Timeout())
	private := e.Group("v2", mw.JWTMiddleware(), mw.Timeout())

	base(public, private)

	dict(public, private)

	product(public, private)

	task(public, private)

	user(public, private)

	file(public, private)

	setting(public, private)

	collect(public, private)

	return e, nil
}
