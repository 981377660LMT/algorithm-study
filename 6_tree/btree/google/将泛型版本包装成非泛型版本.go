// 将泛型版本包装成非泛型版本的技巧
// 类型定义（type definition）
// 通过这个 “包装”（或者称为 “薄封装”）技巧，可以在 保持旧的接口和调用方式（不显式声明泛型）的前提下，把实现细节委托给内部的泛型版本。这样即使 Go 1.18+ 新增了泛型，也无需改动原先依赖 BTree 的代码，实现了平滑升级和向后兼容。
package main

type ObjectG[E any] struct {
	value E
}

func NewObjectG[E any](value E) *ObjectG[E] {
	return &ObjectG[E]{value: value}
}

var Object ObjectG[any]
