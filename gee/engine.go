package gee

import (
	"github.com/Godyu97/geeweb/common"
	"html/template"
	"net/http"
	"strings"
)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
	//html render
	htmlTemplates *template.Template
	htmlFuncMap   template.FuncMap
}

func New() *Engine {
	mux := &Engine{
		router: newRouter(),
	}
	mux.RouterGroup = &RouterGroup{engine: mux}
	mux.groups = []*RouterGroup{mux.RouterGroup}
	return mux
}

func Default() *Engine {
	mux := New()
	mux.Use(Logger(), Recovery())
	return mux
}

func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	e.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute(common.GET, pattern, handler)
}

// POST defines the method to add POST request
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute(common.POST, pattern, handler)
}

// Run defines the method to start a http server
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middlewares := make([]HandlerFunc, 0)
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, r)
	c.handlers = middlewares
	c.engine = e
	e.router.handle(c)
}
