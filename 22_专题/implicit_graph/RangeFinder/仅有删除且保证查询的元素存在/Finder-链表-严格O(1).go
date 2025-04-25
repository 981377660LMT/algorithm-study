package main

type Finder struct {
	n          int
	exist      []bool
	prev, next []int
}

// 建立一个包含0到n-1的集合.
func NewFinder(n int) *Finder {
	res := &Finder{
		n:     n,
		exist: make([]bool, n),
		prev:  make([]int, n),
		next:  make([]int, n),
	}
	for i := 0; i < n; i++ {
		res.exist[i] = true
		res.prev[i] = i - 1
		res.next[i] = i + 1
	}
	return res
}

// 0<=i<n.
func (fs *Finder) Has(i int) bool {
	return i >= 0 && i < fs.n && fs.exist[i]
}

// 0<=i<n.
func (fs *Finder) Erase(i int) bool {
	if !fs.Has(i) {
		return false
	}

	l, r := fs.prev[i], fs.next[i]
	if l >= 0 {
		fs.next[l] = r
	}
	if r < fs.n {
		fs.prev[r] = l
	}
	fs.exist[i] = false
	return true
}

// 返回`严格`小于i的最大元素,如果不存在,返回-1.
// !调用时需保证 Has(i)==true.
// 0<=i<n.
func (fs *Finder) Prev(i int) int {
	return fs.prev[i]
}

// 返回`严格`大于i的最小元素.如果不存在,返回n.
// !调用时需保证 Has(i)==true.
// 0<=i<n.
func (fs *Finder) Next(i int) int {
	return fs.next[i]
}
