package main

type Nested[T any] interface {
	IsNested() bool       // 当前元素是否是嵌套的.
	Flatten() []Nested[T] // 展开嵌套元素.如果当前元素是嵌套的,则展开嵌套元素,否则返回空数组.
	Get() T               // 获取当前元素的值.
}

type Collapse[T any] struct {
	stack [][]Nested[T]
}

func NewCollapse[T any](ns ...Nested[T]) *Collapse[T] {
	return &Collapse[T]{stack: [][]Nested[T]{ns}}
}

func (c *Collapse[T]) Next() T {
	// 由于保证调用 Next 之前会调用 HasNext，直接返回栈顶列表的队首元素，将其弹出队首并返回
	queue := c.stack[len(c.stack)-1]
	res := queue[0].Get()
	c.stack[len(c.stack)-1] = queue[1:]
	return res
}

func (c *Collapse[T]) HasNext() bool {
	for len(c.stack) > 0 {
		queue := c.stack[len(c.stack)-1]
		if len(queue) == 0 { // 当前队列为空, 回退到上一层
			c.stack = c.stack[:len(c.stack)-1]
			continue
		}
		nest := queue[0]
		if !nest.IsNested() {
			return true
		}
		// 若队首元素为列表，则将其弹出队列并入栈
		c.stack[len(c.stack)-1] = queue[1:]
		c.stack = append(c.stack, nest.Flatten())
	}
	return false
}
