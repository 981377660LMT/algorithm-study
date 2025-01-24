// • 单例模式宏观介绍
// • 饿汉式单例模式实现思路
// • 懒汉式单例模式实现推演
// • Golang 单例工具 sync.Once 源码解析

package main

import "sync"

func main() {

}

// #region 懒汉

var (
	s   *singleton
	mux sync.Mutex
)

type singleton struct {
}

func newSingleton() *singleton {
	return &singleton{}
}

func (s *singleton) Work() {
}

type Instance interface {
	Work()
}

func GetInstance() Instance {
	if s != nil {
		return s
	}
	mux.Lock()
	defer mux.Unlock()
	if s != nil {
		return s
	}
	s = newSingleton()
	return s
}

// #endregion

// #region sync.Once

var (
	s2   *singleton2
	once sync.Once
)

type singleton2 struct {
}

func newSingleton2() *singleton2 {
	return &singleton2{}
}

func (s *singleton2) Work() {
}

type Instance2 interface {
	Work()
}

func GetInstance2() Instance2 {
	once.Do(func() {
		s2 = newSingleton2()
	})
	return s2
}

// #endregion
