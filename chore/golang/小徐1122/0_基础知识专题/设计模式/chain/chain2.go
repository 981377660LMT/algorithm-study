// grpc-go 框架中的 InterceptorChain
// 每当有一笔 grpc 请求到达服务端后，服务端会根据请求的 path 匹配到对应的 handler，
// 并且在执行 handler 之前，会先通过执行一个拦截器链 InterceptorChain，完成一些前后置附属逻辑的执行.
package main

import (
	"context"
	"fmt"
)

func main() {
	// 声明好 finnal handler
	var handler Handler = func(ctx context.Context, req []string) ([]string, error) {
		fmt.Printf("final handler is running, req: %+v", req)
		req = append(req, "finnal_handler")
		return req, nil
	}

	// 声明拦截器 1
	var interceptor1 Interceptor = func(ctx context.Context, req []string, handler Handler) ([]string, error) {
		fmt.Println("interceptor1 preprocess...")
		req = append(req, "interceptor1_preprocess")
		resp, err := handler(ctx, req)
		fmt.Println("interceptor1 postprocess")
		resp = append(resp, "interceptor1_postprocess")
		return resp, err
	}

	// 声明拦截器 2
	var interceptor2 Interceptor = func(ctx context.Context, req []string, handler Handler) ([]string, error) {
		fmt.Println("interceptor2 preprocess...")
		req = append(req, "interceptor2_preprocess")
		resp, _ := handler(ctx, req)
		fmt.Println("interceptor2 postprocess")
		resp = append(resp, "interceptor2_postprocess")
		return resp, nil
	}

	chainedInterceptor := chainInterceptors([]Interceptor{
		interceptor1, interceptor2,
	})

	resp, err := chainedInterceptor(context.Background(), nil, handler)
	if err != nil {
		fmt.Printf("error: %+v", err)
		return
	}
	fmt.Printf("resp: %+v", resp)

}

type Handler func(ctx context.Context, req []string) ([]string, error)

type Interceptor func(ctx context.Context, req []string, handler Handler) ([]string, error)

func chainInterceptors(interceptors []Interceptor) Interceptor {
	if len(interceptors) == 0 {
		return nil
	}
	// 返回一个拦截器 interceptor 类型的闭包执行函数
	return func(ctx context.Context, req []string, handler Handler) ([]string, error) {
		return interceptors[0](ctx, req, getChainHandler(interceptors, 0, handler))
	}
}

// 从 inteceptor list 的 index + 1 位置开始，结合 finnal handler，形成一个增强版的 handler
func getChainHandler(interceptors []Interceptor, index int, finalHandler Handler) Handler {
	if index == len(interceptors)-1 {
		return finalHandler
	}
	return func(ctx context.Context, req []string) ([]string, error) {
		return interceptors[index+1](ctx, req, getChainHandler(interceptors, index+1, finalHandler))
	}
}
