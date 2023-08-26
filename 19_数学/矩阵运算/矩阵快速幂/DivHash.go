package main

import (
	"fmt"
	"math"
)

func main() {
	hash := NewDivHash(100)
	fmt.Println(hash.Get(1))
	fmt.Println(hash.Get(2))
	fmt.Println(hash.Get(3))
	fmt.Println(hash.Get(4))
	fmt.Println(hash.Get(5))
	fmt.Println(hash.Get(99))
	fmt.Println(hash.Get(19))
	fmt.Println(hash.Get(100))
}

type DivHash struct {
	n, sn int
}

// 根据n//i的值来哈希.需要2*sqrt(n)的空间.
func NewDivHash(n int) *DivHash {
	res := &DivHash{n: n}
	res.sn = int(math.Sqrt(float64(n)))
	return res
}

func (h *DivHash) Get(x int) int {
	if x <= h.sn {
		return x
	}
	return h.sn*2 - h.n/x + 1
}
