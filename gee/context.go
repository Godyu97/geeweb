package gee

import (
	"encoding/json"
	"fmt"
	"github.com/Godyu97/geeweb/common"
	"net/http"
)

type H map[string]any

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
	// engine
	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(k, v string) {
	c.Writer.Header().Set(k, v)
}

// response
func (c *Context) String(code int, format string, values ...any) {
	c.SetHeader(common.ContentType, common.Text)
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}
func (c *Context) Json(code int, obj any) {
	c.SetHeader(common.ContentType, common.Json)
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		//500
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}
func (c *Context) Fail(code int, err error) {
	c.index = len(c.handlers)
	c.Json(code, H{"err": err.Error()})
}
