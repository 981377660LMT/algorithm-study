package main

import (
	"fmt"
	"sync"
)

// 用于本地测试
func main() {
	n := 5
	foobar := NewFooBar(n)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		foobar.Bar(func() { fmt.Print("bar") })
	}()

	go func() {
		defer wg.Done()
		foobar.Foo(func() { fmt.Print("foo") })
	}()

	wg.Wait()
}

type FooBar struct {
	n     int
	fooCh chan struct{}
	barCh chan struct{}
}

func NewFooBar(n int) *FooBar {
	fb := &FooBar{
		n:     n,
		fooCh: make(chan struct{}, 1),
		barCh: make(chan struct{}, 1),
	}
	fb.fooCh <- struct{}{}
	return fb
}

func (fb *FooBar) Foo(printFoo func()) {
	for i := 0; i < fb.n; i++ {
		<-fb.fooCh
		printFoo()
		fb.barCh <- struct{}{}
	}
}

func (fb *FooBar) Bar(printBar func()) {
	for i := 0; i < fb.n; i++ {
		<-fb.barCh
		printBar()
		fb.fooCh <- struct{}{}
	}
}
