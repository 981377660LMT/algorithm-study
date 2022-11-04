// !用两个 slice 头对头拼在一起实现(每个slice头部只删除元素，不添加元素，互相弥补劣势)
// https://github.dev/EndlessCheng/codeforces-go/blob/master/misc/atcoder/abc274/e

package arraydeque

import (
	"cmnx/src/deque"
	"fmt"
	"strings"
)

// Assert Deque implementation
var _ deque.Deque = (*ArrayDeque)(nil)

func main() {
	queue := NewArrayDeque(10)
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

	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
}

func NewArrayDeque(numElements int) *ArrayDeque {
	half := numElements / 2
	return &ArrayDeque{
		left:  make([]interface{}, 0, half+1),
		right: make([]interface{}, 0, half+1),
	}
}

type ArrayDeque struct {
	left, right []interface{}
}

func (queue *ArrayDeque) Append(value interface{}) {
	queue.right = append(queue.right, value)
}

func (queue *ArrayDeque) AppendLeft(value interface{}) {
	queue.left = append(queue.left, value)
}

func (queue *ArrayDeque) Pop() (value interface{}) {
	if queue.Len() == 0 {
		return
	}

	if len(queue.right) > 0 {
		queue.right, value = queue.right[:len(queue.right)-1], queue.right[len(queue.right)-1]
	} else {
		value, queue.left = queue.left[0], queue.left[1:]
	}

	return
}

func (queue *ArrayDeque) PopLeft() (value interface{}) {
	if queue.Len() == 0 {
		return
	}

	if len(queue.left) > 0 {
		queue.left, value = queue.left[:len(queue.left)-1], queue.left[len(queue.left)-1]
	} else {
		value, queue.right = queue.right[0], queue.right[1:]
	}

	return
}

func (queue *ArrayDeque) At(index int) (value interface{}) {
	n := queue.Len()
	if index < 0 {
		index += n
	}

	if index < 0 || index >= n {
		return
	}

	if index < len(queue.left) {
		value = queue.left[len(queue.left)-1-index]
	} else {
		value = queue.right[index-len(queue.left)]
	}

	return
}

func (queue *ArrayDeque) ForEach(f func(value interface{}, index int)) {
	leftLen := len(queue.left)
	for i := 0; i < leftLen; i++ {
		f(queue.left[i], i)
	}

	for i := 0; i < len(queue.right); i++ {
		f(queue.right[i], leftLen+i)
	}
}

func (queue *ArrayDeque) Len() int {
	return len(queue.left) + len(queue.right)
}

func (queue *ArrayDeque) String() string {
	res := []string{"ArrayDeque{"}
	values := []string{}
	queue.ForEach(func(value interface{}, index int) {
		values = append(values, fmt.Sprintf("%v", value))
	})

	res = append(res, strings.Join(values, ", "), "}")
	return strings.Join(res, "")
}
