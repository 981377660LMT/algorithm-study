// 单个chan做十分简洁
// 不仅保证一开始Foo先做
// 而且最后Bar结束时发送，正好Foo也在等待接受，所以都不会卡住

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
	n  int
	ch chan struct{}
}

func NewFooBar(n int) *FooBar {
	return &FooBar{n: n, ch: make(chan struct{})}
}

func (fb *FooBar) Foo(printFoo func()) {
	for i := 0; i < fb.n; i++ {
		// printFoo() outputs "foo". Do not change or remove this line.
		printFoo()
		fb.ch <- struct{}{}
		<-fb.ch
	}
}

func (fb *FooBar) Bar(printBar func()) {
	for i := 0; i < fb.n; i++ {
		<-fb.ch
		// printBar() outputs "bar". Do not change or remove this line.
		printBar()
		fb.ch <- struct{}{}
	}
}
