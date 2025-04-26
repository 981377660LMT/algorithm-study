// https://github.com/981377660LMT/ts/issues/806
//
// golang里的F-Bounded多态
// 两步走

package main

import "fmt"

// F-Bounded多态
type Constraint[T any] interface {
	Value() T
}

func DoSomething[T Constraint[T]](c T) T {
	fmt.Println("Doing something")
	return c
}

// 测试F-Bounded多态
type MyType struct{ name int }

func (m MyType) Value() MyType { return m }
func (m MyType) OtherMethod()  {}

func main() {
	obj := MyType{name: 1}
	res := DoSomething(obj)
	res.OtherMethod()
}
