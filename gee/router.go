package gee

import (
	"fmt"
	"github.com/Godyu97/geeweb/common"
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(ctx *Context)

// 动态路由的支持
// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	parts := parsePattern(pattern)
	key := MakeRouterKey(method, pattern)
	_, ok := r.roots[method]
	if ok == false {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	reqParts := parsePattern(pattern)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if ok == false {
		return nil, nil
	}
	n := root.search(reqParts, 0)
	if n != nil {
		ps := parsePattern(n.pattern)
		for i, p := range ps {
			//:uri
			if p[0] == common.TrieFlag1 {
				params[p[1:]] = reqParts[i]
			}
			//*uri
			if p[0] == common.TrieFlag2 && len(p) > 1 {
				params[p[1:]] = strings.Join(reqParts[i:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := MakeRouterKey(c.Method, n.pattern)
		if handler, ok := r.handlers[key]; ok {
			handler(c)
		} else {
			c.String(http.StatusInternalServerError, "500 aeDPRbPi: %s\n", c.Path)
		}
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

func MakeRouterKey(method string, pattern string) string {
	return fmt.Sprintf("%s_%s", method, pattern)
}
