// gin 框架中的 HandlerChain

package main

import "fmt"

func main() {
	var finalHandler Handler = func(ctx *Context, req map[string]any) (map[string]any, error) {
		req["final_handler"] = true
		return req, nil
	}

	var middleware1 Handler = func(ctx *Context, req map[string]any) (map[string]any, error) {
		req["middleware1_preprocess"] = true
		ctx.Next(req)
		req["middleware1_postprocess"] = true
		return req, nil
	}

	var middleware2 Handler = func(ctx *Context, req map[string]any) (map[string]any, error) {
		req["middleware2_preprocess"] = true
		ctx.Next(req)
		req["middleware2_postprocess"] = true
		return req, nil
	}

	var middleware1WithAbort Handler = func(ctx *Context, req map[string]any) (map[string]any, error) {
		req["middleware1_preprocess"] = true
		ctx.Abort()
		return req, nil
	}

	{
		ctx := NewContext(middleware1, middleware2, finalHandler)
		params := map[string]any{}
		ctx.Next(params)
		fmt.Println(params)
	}

	{
		ctx := NewContext(middleware1WithAbort, middleware2, finalHandler)
		params := map[string]any{}
		ctx.Next(params)
		fmt.Println(params)
	}
}

// handler 链最大长度为 100
const MaxHandlersCnt = 100

type Handler func(ctx *Context, req map[string]any) (map[string]any, error)

type Context struct {
	index    int
	handlers []Handler
}

func NewContext(handlers ...Handler) *Context {
	return &Context{
		index:    -1,
		handlers: handlers,
	}
}

func (c *Context) Next(req map[string]any) {
	c.index++
	for c.index < len(c.handlers) && c.index <= MaxHandlersCnt {
		c.handlers[c.index](c, req)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = MaxHandlersCnt
}
