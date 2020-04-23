package gee

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	//存储路由和路由执行函数
	router map[string]HandlerFunc
}

// 构造方法
func NewEngine() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}
func (e *Engine) addRoute(method, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	e.router[key] = handler
}

//注册路由
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

//所有请求都经由此方法分发
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(w, r)
	} else {
		fmt.Fprintf(w, "404 not found path %s", key)
	}
}
