package gee

import (
	"net/http"
	"strings"
)

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

//解析url 且只能允许存在一个*匹配符
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
	root, ok := r.roots[method]
	if !ok {
		root = &node{}
		r.roots[method] = root
	}
	parts := parsePattern(pattern)
	root.insert(pattern, parts, 0)
}

//根据请求url找到对应的node（精确匹配/模糊匹配）如果是模糊匹配解析参数
func (r *router) getRoute(method, pattern string) (*node, map[string]string) {
	//请求访问的url
	searchParts := parsePattern(pattern)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	node := root.search(searchParts, 0)
	if node == nil {
		return nil, nil
	}
	//注册路由时的url
	parts := parsePattern(node.pattern)
	//根据请求时访问的url和匹配的路由可以对应出url参数
	for index, part := range parts {
		if part[0] == ':' {
			key := part[1:]
			params[key] = searchParts[index]
		}
		//只能有一个*
		if part[0] == '*' && len(part) > 1 {
			key := part[1:]
			params[key] = strings.Join(searchParts[index:], "/")
			break
		}
	}

	return node, params
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		//将请求路径参数赋给上下文
		c.Params = params
		//找到请求路径对应的执行函数
		key := c.Method + "-" + n.pattern //这里是node的pattern
		//执行
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 not found %s \n", c.Path)
	}
}
