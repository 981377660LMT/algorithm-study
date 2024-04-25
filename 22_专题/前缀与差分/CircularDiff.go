// 环形差分

package main

import "fmt"

func main() {
	D := NewCircularDiff(10)
	D.AddCircle(8, 20, 2)
	D.Build()
	fmt.Println(D.GetAll())
}

// 环形差分
type CircularDiff struct {
	n    int
	diff []int
}

func NewCircularDiff(n int) *CircularDiff {
	return &CircularDiff{n: n, diff: make([]int, n+1)}
}

func (cd *CircularDiff) Add(start, end, x int) {
	if start < 0 {
		start = 0
	}
	if end > cd.n {
		end = cd.n
	}
	if start >= end {
		return
	}
	cd.diff[start] += x
	cd.diff[end] -= x
}

func (cd *CircularDiff) AddCircle(start, end, x int) {
	if start >= end {
		return
	}
	n := cd.n
	loop := (end - start) / n
	if loop > 0 {
		cd.Add(0, n, x*loop)
	}
	if (end-start)%n == 0 {
		return
	}
	start %= n
	end %= n
	if start < end {
		cd.Add(start, end, x)
	} else {
		cd.Add(start, n, x)
		if end > 0 {
			cd.Add(0, end, x)
		}
	}
}

func (cd *CircularDiff) Build() {
	for i := 1; i <= cd.n; i++ {
		cd.diff[i] += cd.diff[i-1]
	}
}

func (cd *CircularDiff) Get(i int) int {
	return cd.diff[i]
}

func (cd *CircularDiff) GetAll() []int {
	return cd.diff[:cd.n]
}
