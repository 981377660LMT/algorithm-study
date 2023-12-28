// InterpolatePeriodicSequence
// 周期序列插值
// !用于发现周期性序列的循环节
// 例如 123[456][456][456]...

package main

import "fmt"

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 4, 5}
	seq := NewInterpolatePeriodicSequence(arr)
	for i := 0; i < 20; i++ {
		fmt.Println(seq.Get(i))
	}
}

type T = int

type InterpolatePeriodicSequence struct {
	rawSeq []T
	offset int
}

func NewInterpolatePeriodicSequence(seq []T) *InterpolatePeriodicSequence {
	revSeq := make([]T, len(seq))
	for i := range seq {
		revSeq[i] = seq[len(seq)-1-i]
	}
	z := ZAlgoNums(revSeq)
	z[0] = 0
	max_, maxIndex := -1, -1
	for i, v := range z {
		if v >= max_ {
			max_ = v
			maxIndex = i
		}
	}
	return &InterpolatePeriodicSequence{rawSeq: seq, offset: maxIndex}
}

func (ips *InterpolatePeriodicSequence) Get(index int) T {
	if index < len(ips.rawSeq) {
		return ips.rawSeq[index]
	}
	k := (index - (len(ips.rawSeq) - 1) + ips.offset - 1) / ips.offset
	index -= k * ips.offset
	return ips.rawSeq[index]
}

func ZAlgoNums(nums []T) []int {
	n := len(nums)
	z := make([]int, n)
	left, right := 0, 0
	for i := 1; i < n; i++ {
		z[i] = max(min(z[i-left], right-i+1), 0)
		for i+z[i] < n && nums[z[i]] == nums[i+z[i]] {
			left, right = i, i+z[i]
			z[i]++
		}
	}
	return z
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
