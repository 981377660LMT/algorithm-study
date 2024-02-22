// ChminStack/ChmaxStack
// 设计一个高效的数据结构，支持两种操作：
// 1.加入一个数x，然后所有元素与x取min；2.求容器内所有元素之和

package main

import "fmt"

func main() {
	S := NewClampableStack(true)
	S.AddAndClamp(1)
	S.AddAndClamp(2)
	S.AddAndClamp(1)
	fmt.Println(S.Sum()) // 3
	S = NewClampableStack(false)
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

// clampMin 为true时，支持容器内所有数与x取min；为false时，支持容器内所有数与x取max.
func NewClampableStack(clampMin bool) *ClampableStack {
	return &ClampableStack{clampMin: clampMin}
}

func (h *ClampableStack) AddAndClamp(x int) {
	newCount := 1
	if h.clampMin {
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
	} else {
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
