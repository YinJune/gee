package gee

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Engine struct {
	//存储路由和路由执行函数
	router *router
}

// 构造方法
func New() *Engine {
	return &Engine{router: newRouter()}
}

//注册路由
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.router.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.router.addRoute("POST", pattern, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

//所有请求都经由此方法分发
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.router.handle(NewContext(w, r))
}
