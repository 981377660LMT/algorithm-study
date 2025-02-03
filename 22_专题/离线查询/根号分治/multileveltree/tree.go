package main

func main() {

}

// 多级分块结构的树.
type MultiLevelTree[E any] struct {
	e  func() E
	op func(a, b E) E
}

func NewMultiLevelTree[E any](e func() E, op func(a, b E) E) *MultiLevelTree[E] {
	return &MultiLevelTree[E]{e: e, op: op}
}

func (m *MultiLevelTree[E]) Build(n int, f func(i int) E) {}
