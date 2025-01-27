// 常规适配模式和 interface 适配模式两种类型
//
// !常规适配模式：
//   适配器实现了目标接口，然后在适配器中调用被适配对象的方法
// !interface 适配模式：
//   https://www.bilibili.com/video/BV14M411H7J9
//   http://github.com/golang/go/wiki/CodeReviewComments#interface
//   interface 的定义应该由模块的使用方而非实现方来进行定义.
//   只有这样，使用方才能根据自己的视角，对模块进行最合适或者说最贴合自己使用需求的抽象定义.//
//   golang的接口是隐式的，只要实现了接口的方法，就是实现了接口
//   不需要显式声明对 interface 的 implement 操作，只需要实现了 interface 的所有方法，就自动会被编译器识别为 interface 的一种实现.

package main

import "fmt"

func main() {
	{
		// 结构体适配模式（对象适配器）
		// 核心特征：通过组合已有结构体实现接口适配
		typeC := &TypeC{}
		usb := NewTypeCAdapter(typeC)
		usb.Connect()
	}

}

// #region 结构体适配模式（对象适配器）

// 目标接口（客户端期望的）
type USB interface {
	Connect()
}

// 被适配结构体（不兼容接口）
type TypeC struct{}

func (t *TypeC) TypeCConnect() {
	fmt.Println("Type-C connected")
}

// 结构体适配器
type TypeCAdapter struct {
	typeC *TypeC
}

func NewTypeCAdapter(t *TypeC) USB {
	return &TypeCAdapter{typeC: t}
}

func (a *TypeCAdapter) Connect() {
	a.typeC.TypeCConnect()
	fmt.Println("Converted to USB connection")
}

// #endregion

// #region 接口转换模式（Go 特色实现）
// #endregion
