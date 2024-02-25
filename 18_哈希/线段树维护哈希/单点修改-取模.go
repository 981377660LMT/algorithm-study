package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	cf1200e()
}

// Compress Words
// https://www.luogu.com.cn/problem/CF1200E
// 给定n个字符串，将其从左到右依次合并，并在合并时去重(相邻的字符串的公共前后缀要只留下一个)
// 比如 abcabc 和 bcd 接起来就是 abcabcd。
// 输出答案.
func cf1200e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	const LIMIT = 1 << 20
	seg := NewSegmentTree(LIMIT, func(i int) E { return e() })
	res := strings.Builder{}
	ptr := 0

	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		tmpSeg := NewSegmentTree(len(s), func(i int) E { return FromElement(uint(s[i])) })
		start := 0
		for len_ := 0; len_ <= len(s); len_++ {
			if ptr < len_ {
				break
			}
			x := seg.Query(ptr-len_, ptr)
			y := tmpSeg.Query(0, len_)
			if x == y {
				start = len_
			}
		}
		for j := start; j < len(s); j++ {
			seg.Set(ptr, FromElement(uint(s[j])))
			ptr++
			res.WriteByte(s[j])
		}
	}

	fmt.Fprintln(out, res.String())
}

// 线段树维护多项式哈希
const BASE uint = 131
const MOD uint = 999999937 // 999999937/999999929/999999893/999999797/999999761/999999757/999999751/999999739

type E = [2]uint

func FromElement(v uint) E { return E{BASE, v} } // pow of base, val
func e() E                 { return E{1, 0} }
func op(a, b E) E          { return E{a[0] * b[0] % MOD, (a[1]*b[0] + b[1]) % MOD} }

type SegmentTree struct {
	n, size int
	seg     []E
}

func NewSegmentTree(n int, f func(int) E) *SegmentTree {
	res := &SegmentTree{}
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *SegmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTree) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTree) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return e()
	}
	leftRes, rightRes := e(), e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return op(leftRes, rightRes)
}
func (st *SegmentTree) QueryAll() E { return st.seg[1] }
func (st *SegmentTree) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}
