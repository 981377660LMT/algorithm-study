package main

import (
	"fmt"
	"reflect"
)

func main() {
	Test_implement()
}

type MyInterface interface {
	MethodA()
	MethodB()
}

type MyClass struct{}

func NewMyClass() *MyClass {
	return &MyClass{}
}

// MyClass 实现了 MyInterface 声明的所有方法
func (m *MyClass) MethodA() {}

func (m *MyClass) MethodB() {}

func Test_implement() {
	// 获取 myClass 的类型
	myClassTyp := reflect.TypeOf(NewMyClass())
	// 获取 myInterface 的类型
	myInterTyp := reflect.TypeOf((*MyInterface)(nil)).Elem()
	// 判断是否具有实现关系
	fmt.Println(myClassTyp.Implements(myInterTyp))
}
