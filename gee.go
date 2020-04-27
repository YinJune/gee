package gee

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Engine struct {
	//一个engine是一个顶层group
	*RouterGroup
	//存储路由和路由执行函数
	router *router
	//保存所有分组
	groups []*RouterGroup
}
type RouterGroup struct{
	prefix string
	middlewares []HandlerFunc
	parent *RouterGroup
	engine *Engine
}
// 构造方法
func New() *Engine {
	engine:=&Engine{}
	engine.router= newRouter()
	engine.RouterGroup=&RouterGroup{engine:engine}
	engine.groups=[]*RouterGroup{engine.RouterGroup}
	return engine 
}
//添加路由分组
func (group *RouterGroup) Group(prefix string) *RouterGroup{
	g:=&RouterGroup{
		prefix:group.prefix+prefix,
		parent:group,
		engine:group.engine,
	}
	group.engine.groups=append(group.engine.groups,g)
	return g
}

func (g *RouterGroup) addRoute(method ,pattern string,handler HandlerFunc){
	pattern=g.prefix+pattern
	g.engine.router.addRoute(method,pattern,handler)
}

//注册路由
func (g *RouterGroup) GET(pattern string, handler HandlerFunc) {
	g.addRoute("GET", pattern, handler)
}

func (g *RouterGroup) POST(pattern string, handler HandlerFunc) {
	g.addRoute("POST", pattern, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

//所有请求都经由此方法分发
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.router.handle(NewContext(w, r))
}
