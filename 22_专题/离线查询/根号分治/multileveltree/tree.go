package main

func main() {

}

// 多级分块结构的树.
type MultiLevelTree[E any] struct {
	e  func() E
	op func(a, b E) E
}

func NewMultiLevelTree[E any](
	e func() E, op func(a, b E) E) *MultiLevelTree[E] {
	return &MultiLevelTree[E]{e: e, op: op}
}
func (m *MultiLevelTree[E]) Build(n int, f func(i int) E) {}

func (m *MultiLevelTree[E]) QueryRange(l, r int) E {
	panic("todo")
}
func (m *MultiLevelTree[E]) QueryAll() E {
	panic("todo")
}
func (m *MultiLevelTree[E]) Get(i int) E {
	panic("todo")
}

func (m *MultiLevelTree[E]) Update(i int, v E) {
	panic("todo")
}
func (m *MultiLevelTree[E]) Set(i int, v E) {
	panic("todo")
}
