// https://colobu.com/gotips/060.html

package main

import "sync"

var (
	once     sync.Once
	instance *Config
)

type Config struct{}

func loadConfig() *Config { return &Config{} }

func GetConfig() *Config {
	once.Do(func() { instance = loadConfig() })
	return instance
}

// 一个sync.Once会跟踪两件事:
// 1. 一个原子计数器(或标志),有0和1两个值。
// 2. 一个用于保护慢路径的互斥锁。
//
// - 快速路径：原子计数器判断函数是否执行过时，执行过则后续的调用可以快速跳过并无需等待
// - 慢速路径：若计数器不为0，则会触发慢速路径。sync.Once 会进入慢速模式并调用 doSlow(f)
//   !doSlow：为什么已经获取了互斥锁，还需要 判断 o.done.Load() == 0 ?
