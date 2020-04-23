package gee

import "net/http"

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) getRoute(method, pattern string) (HandlerFunc, bool) {
	key := method + "-" + pattern
	handler, ok := r.handlers[key]
	return handler, ok
}

func (r *router) handle(c *Context) {
	handler, ok := r.getRoute(c.Method, c.Path)
	if ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 not found %s \n", c.Path)
	}
}
