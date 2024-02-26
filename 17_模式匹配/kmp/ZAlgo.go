package main

import "fmt"

func main() {
	fmt.Println(ZAlgo("ababab"))
}

func sumScores(s string) int64 {
	n := len(s)
	res := int64(0)
	z := ZAlgo(s)
	for i := 0; i < n; i++ {
		res += int64(z[i])
	}
	return res + int64(n)
}

// z算法求字符串每个后缀与原串的最长公共前缀长度
//
// z[0]=0
// z[i]是后缀s[i:]与s的最长公共前缀(LCP)的长度 (i>=1)
func ZAlgo(s string) []int {
	n := len(s)
	z := make([]int, n)
	left, right := 0, 0
	for i := 1; i < n; i++ {
		z[i] = max(min(z[i-left], right-i+1), 0)
		for i+z[i] < n && s[z[i]] == s[i+z[i]] {
			left, right = i, i+z[i]
			z[i]++
		}
	}
	return z
}

func ZAlgoNums(nums []int) []int {
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

type DiffArray struct {
	diff  []int
	dirty bool
}

func NewDiffArray(n int) *DiffArray {
	return &DiffArray{
		diff: make([]int, n+1),
	}
}

func (d *DiffArray) Add(start, end, delta int) {
	if start < 0 {
		start = 0
	}
	if end >= len(d.diff) {
		end = len(d.diff) - 1
	}
	if start >= end {
		return
	}
	d.dirty = true
	d.diff[start] += delta
	d.diff[end] -= delta
}

func (d *DiffArray) Build() {
	if d.dirty {
		preSum := make([]int, len(d.diff))
		for i := 1; i < len(d.diff); i++ {
			preSum[i] = preSum[i-1] + d.diff[i]
		}
		d.diff = preSum
		d.dirty = false
	}
}

func (d *DiffArray) Get(pos int) int {
	d.Build()
	return d.diff[pos]
}

func (d *DiffArray) GetAll() []int {
	d.Build()
	return d.diff[:len(d.diff)-1]
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
