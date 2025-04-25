package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fs := NewFinder2(5)
	fmt.Println(fs) // Finder{0, 1, 2, 3, 4}

	fs.Erase(2)
	fmt.Println(fs) // Finder{0, 1, 3, 4}

	fmt.Println("Has(2)?", fs.Has(2))    // false
	fmt.Println("Next(2) =", fs.Next(2)) // 3
	fmt.Println("Prev(2) =", fs.Prev(2)) // 1

	debug := func() {
		for i := 0; i < fs.n; i++ {
			fmt.Printf("Prev(%d) = %d\n", i, fs.Prev(i))
		}
		for i := 0; i < fs.n; i++ {
			fmt.Printf("Next(%d) = %d\n", i, fs.Next(i))
		}
	}
	_ = debug

	// debug()
}

type Finder2 struct {
	n          int
	prev, next []int
}

// 建立一个包含0到n-1的集合.
func NewFinder2(n int) *Finder2 {
	res := &Finder2{
		n:    n,
		prev: make([]int, n),
		next: make([]int, n),
	}
	for i := 0; i < n; i++ {
		res.prev[i] = i - 1
		res.next[i] = i + 1
	}
	return res
}

// 0<=i<n.
func (fs *Finder2) Has(i int) bool {
	if i < 0 || i >= fs.n {
		return false
	}
	return !(fs.prev[i] == i && fs.next[i] == i)
}

// 0<=i<n.
func (fs *Finder2) Erase(i int) bool {
	if !fs.Has(i) {
		return false
	}
	p, q := fs.prev[i], fs.next[i]
	if p >= 0 {
		fs.next[p] = q
	}
	if q < fs.n {
		fs.prev[q] = p
	}
	fs.prev[i] = i
	fs.next[i] = i
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
// 0<=i<n.
func (fs *Finder2) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}
	if fs.Has(i) {
		return i
	}
	fmt.Println(11)
	return fs.next[i]
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
// 0<=i<n.
func (fs *Finder2) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}
	if fs.Has(i) {
		return i
	}
	return fs.prev[i]
}

func (fs *Finder2) String() string {
	var res []string
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("Finder{%v}", strings.Join(res, ", "))
}
