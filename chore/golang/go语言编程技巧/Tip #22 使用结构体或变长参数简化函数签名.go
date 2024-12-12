// 为了保持代码整洁，可以考虑以下两种策略：
//
// 结构体作为参数 (struct options)
// !变长参数 (variadic options)
//
// 这里主要介绍变长参数的方式，即使用变长参数来简化函数签名。
// 变长参数设置参数默认值比使用结构体作为参数更方便，你不需要隐藏它，直接在 ConnectToDatabase 函数中设置默认值即可。

package main

import "fmt"

func main() {
	ConnectToDatabase(WithHost("localhost"), WithPort(3306), WithUsername("root"), WithPassword("root"))
}

type ServiceConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	SSL      bool
}

type ServiceOptionFunc func(*ServiceConfig)

func WithHost(host string) ServiceOptionFunc { return func(cfg *ServiceConfig) { cfg.Host = host } }
func WithPort(port int) ServiceOptionFunc    { return func(cfg *ServiceConfig) { cfg.Port = port } }
func WithUsername(username string) ServiceOptionFunc {
	return func(cfg *ServiceConfig) { cfg.Username = username }
}
func WithPassword(password string) ServiceOptionFunc {
	return func(cfg *ServiceConfig) { cfg.Password = password }
}
func WithSSL(ssl bool) ServiceOptionFunc { return func(cfg *ServiceConfig) { cfg.SSL = ssl } }

func ConnectToDatabase(optionFuncs ...ServiceOptionFunc) {
	cfg := ServiceConfig{} // Default configuration
	for _, option := range optionFuncs {
		option(&cfg)
	}

	// Connect to database using cfg
	fmt.Printf("Connecting to database with config: %+v\n", cfg)
}
