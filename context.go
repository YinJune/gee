package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Method         string
	Path           string
	StatusCode     int
}

//构造函数
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		ResponseWriter: w,
		Request:        r,
		Method:         r.Method,
		Path:           r.URL.Path,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// 设置响应码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.ResponseWriter.WriteHeader(code)
}

//设置请求头
func (c *Context) SetHeader(key, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

//返回字符串
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain")
	c.ResponseWriter.Write([]byte(fmt.Sprintf(format, values...)))
}

//返回JSON
func (c *Context) JSON(code int, obj interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.ResponseWriter)
	err := encoder.Encode(obj)
	if err != nil {
		http.Error(c.ResponseWriter, err.Error(), 500)
	}
}

//返回二进制数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.ResponseWriter.Write(data)
}

//返回html
func (c *Context) HTML(code int, html string) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/html")
	c.ResponseWriter.Write([]byte(html))
}
