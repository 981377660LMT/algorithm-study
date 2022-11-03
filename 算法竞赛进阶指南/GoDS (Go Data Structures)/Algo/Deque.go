// !用两个 slice 头对头拼在一起实现(每个slice头部只删除元素，不添加元素，互相弥补劣势)
// https://github.dev/EndlessCheng/codeforces-go/blob/master/misc/atcoder/abc274/e

package main

import (
	"fmt"
	"strings"
)

func main() {
	queue := NewSliceDeque(10)
	queue.Append(1)
	queue.Append(2)
	queue.Pop()
	queue.Append(4)
	queue.AppendLeft("a")

	fmt.Println(queue)
	queue.ForEach(func(value interface{}, index int) {
		fmt.Println(value, index)
	})

	fmt.Println(queue.At(0))
	fmt.Println(queue.At(-1))
}

func NewSliceDeque(numElements int) Deque {
	half := numElements / 2
	return &SliceDeque{
		left:  make([]interface{}, 0, half+1),
		right: make([]interface{}, 0, half+1),
	}
}

type Deque interface {
	Append(value interface{})
	AppendLeft(value interface{})
	Pop() (value interface{}, ok bool)
	PopLeft() (value interface{}, ok bool)
	At(index int) (value interface{}, ok bool)
	ForEach(func(value interface{}, index int))
	Len() int
}

type SliceDeque struct {
	left, right []interface{}
}

func (queue *SliceDeque) Append(value interface{}) {
	queue.right = append(queue.right, value)
}

func (queue *SliceDeque) AppendLeft(value interface{}) {
	queue.left = append(queue.left, value)
}

func (queue *SliceDeque) Pop() (value interface{}, ok bool) {
	if queue.Len() == 0 {
		return nil, false
	}

	if len(queue.right) > 0 {
		queue.right, value = queue.right[:len(queue.right)-1], queue.right[len(queue.right)-1]
	} else {
		value, queue.left = queue.left[0], queue.left[1:]
	}

	ok = true
	return
}

func (queue *SliceDeque) PopLeft() (value interface{}, ok bool) {
	if queue.Len() == 0 {
		return nil, false
	}

	if len(queue.left) > 0 {
		queue.left, value = queue.left[:len(queue.left)-1], queue.left[len(queue.left)-1]
	} else {
		value, queue.right = queue.right[0], queue.right[1:]
	}

	ok = true
	return
}

func (queue *SliceDeque) At(index int) (value interface{}, ok bool) {
	n := queue.Len()
	if index < 0 {
		index += n
	}

	if index < 0 || index >= n {
		return nil, false
	}

	if index < len(queue.left) {
		value = queue.left[len(queue.left)-1-index]
	} else {
		value = queue.right[index-len(queue.left)]
	}

	ok = true
	return
}

func (queue *SliceDeque) ForEach(f func(value interface{}, index int)) {
	leftLen := len(queue.left)
	for i := 0; i < leftLen; i++ {
		f(queue.left[i], i)
	}

	for i := 0; i < len(queue.right); i++ {
		f(queue.right[i], leftLen+i)
	}
}

func (queue *SliceDeque) Len() int {
	return len(queue.left) + len(queue.right)
}

func (queue *SliceDeque) String() string {
	res := []string{"SliceDeque{"}
	values := []string{}
	queue.ForEach(func(value interface{}, index int) {
		values = append(values, fmt.Sprintf("%v", value))
	})

	res = append(res, strings.Join(values, ", "), "}")
	return strings.Join(res, "")
}
