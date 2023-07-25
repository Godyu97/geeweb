package gee

import (
	"net/http"
	"path"
)

const ParamFilePath = "filepath"

func (g *RouterGroup) createStaticHandler(pattern string, fs http.FileSystem) HandlerFunc {
	reqPath := path.Join(g.prefix, pattern)
	//用于提供文件服务
	//StripPrefix 在fs处理请求时去掉前缀reqPath
	fileServer := http.StripPrefix(reqPath, http.FileServer(fs))
	return func(ctx *Context) {
		file := ctx.Param(ParamFilePath)
		//Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		//req.url.path作为文件路径
		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

// serve static files
// pattern 为请求该文件夹req路径
// root 为磁盘上的某个文件夹
func (g *RouterGroup) Static(pattern string, root string) {
	handler := g.createStaticHandler(pattern, http.Dir(root))
	urlPattern := path.Join(pattern, "/*"+ParamFilePath)
	g.GET(urlPattern, handler)
}
