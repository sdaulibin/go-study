package gee

import "net/http"

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (router *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	router.handlers[key] = handler
}

func (router *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler,ok := router.handlers[key];ok{
		handler(c)
	} else {
		c.String(http.StatusNotFound,"404 not found %s\n",c.Path)
	}
}
