// ChminStack/ChmaxStack
// 设计一个高效的数据结构，支持两种操作：
// 1.加入一个数x，然后所有元素与x取min；2.求容器内所有元素之和
//
// 一般配合后缀数组的height数组使用.

package main

import "fmt"

func main() {
	S := NewClampableStack(false)
	S.AddAndClamp(1)
	S.AddAndClamp(2)
	S.AddAndClamp(1)
	fmt.Println(S.Sum()) // 3
	S = NewClampableStack(true)
	S.AddAndClamp(1)
	S.AddAndClamp(2)
	S.AddAndClamp(1)
	fmt.Println(S.Sum()) // 5
}

type H = struct {
	value int
	count int32
}

type ClampableStack struct {
	clampMin bool
	total    int
	count    int
	stack    []H
}

// clampMin：
//  为true时，调用AddAndClamp(x)后，容器内所有数最小值被截断(小于x的数变成x)；
//  为false时，调用AddAndClamp(x)后，容器内所有数最大值被截断(大于x的数变成x).
func NewClampableStack(clampMin bool) *ClampableStack {
	return &ClampableStack{clampMin: clampMin}
}

func (h *ClampableStack) AddAndClamp(x int) {
	newCount := 1
	if h.clampMin {
		for len(h.stack) > 0 {
			top := h.stack[len(h.stack)-1]
			if top.value > x {
				break
			}
			h.stack = h.stack[:len(h.stack)-1]
			v, c := top.value, int(top.count)
			h.total -= v * c
			newCount += c
		}
	} else {
		for len(h.stack) > 0 {
			top := h.stack[len(h.stack)-1]
			if top.value < x {
				break
			}
			h.stack = h.stack[:len(h.stack)-1]
			v, c := top.value, int(top.count)
			h.total -= v * c
			newCount += c
		}
	}
	h.total += x * newCount
	h.count++
	h.stack = append(h.stack, H{value: x, count: int32(newCount)})
}

func (h *ClampableStack) Sum() int {
	return h.total
}

func (h *ClampableStack) Len() int {
	return h.count
}

func (h *ClampableStack) Clear() {
	h.stack = h.stack[:0]
	h.total = 0
	h.count = 0
}
