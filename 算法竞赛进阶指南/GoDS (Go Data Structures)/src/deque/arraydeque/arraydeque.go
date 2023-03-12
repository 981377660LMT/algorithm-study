// !用两个 slice 头对头拼在一起实现(每个slice头部只删除元素，不添加元素，互相弥补劣势)
// 在知道数据量的情况下，也可以直接创建一个两倍数据量大小的 slice，然后用两个下标表示头尾，初始化在 slice 正中
// https://github.dev/EndlessCheng/codeforces-go/blob/master/misc/atcoder/abc274/e

package arraydeque

import (
	"fmt"
	"strings"
)

func main() {
	queue := NewDeque(10)
	queue.Append(1)
	queue.Append(2)
	queue.Pop()
	queue.Append(4)

	fmt.Println(queue)
	queue.ForEach(func(value E, index int) {
		fmt.Println(value, index)
	})

	fmt.Println(queue.At(0))
	fmt.Println(queue.At(-1))

	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
}

type E = int

func NewDeque(cap int) *ArrayDeque {
	half := cap / 2
	return &ArrayDeque{
		left:  make([]E, 0, half+1),
		right: make([]E, 0, half+1),
	}
}

type ArrayDeque struct {
	left, right []E
}

func (queue *ArrayDeque) Append(value E) {
	queue.right = append(queue.right, value)
}

func (queue *ArrayDeque) AppendLeft(value E) {
	queue.left = append(queue.left, value)
}

func (queue *ArrayDeque) Pop() (value E) {
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

func (queue *ArrayDeque) PopLeft() (value E) {
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

func (queue *ArrayDeque) At(index int) (value E) {
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

func (queue *ArrayDeque) ForEach(f func(value E, index int)) {
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
	queue.ForEach(func(value E, _ int) {
		values = append(values, fmt.Sprintf("%v", value))
	})

	res = append(res, strings.Join(values, ", "), "}")
	return strings.Join(res, "")
}
