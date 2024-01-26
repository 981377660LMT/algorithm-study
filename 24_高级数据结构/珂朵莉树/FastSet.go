// 又叫做 64-ary tree
// !时间复杂度:O(log64n)
// https://zhuanlan.zhihu.com/p/107238627
// https://www.luogu.com.cn/blog/RuntimeErrror/ni-suo-fou-zhi-dao-di-shuo-ju-jie-gou-van-emde-boas-shu
// 使用场景:
// 1. 在存储IP地址的时候， 需要快速查找某个IP地址（2 ^32大小)是否在访问的列表中，
//    或者需要找到比这个IP地址大一点或者小一点的IP作为重新分配的IP。
// 2. 一条路上开了很多商店，用int来表示商店的位置（假设位置为1-256之间的数），
//    不断插入，删除商店，同时需要找到离某个商店最近的商店在哪里。

// !Insert/Erase/Prev/Next/Has/Enumerate
// https://maspypy.github.io/library/ds/fastset.hpp

// ! 注意频繁查找普通属性耗时(res.B), 把B写在代码里面或者定义成const,
// ! 会快很多(200ms->30ms), 编译器会优化

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fs := NewFastSetFrom(10, func(i int) bool { return i%2 == 0 })
	fmt.Println(fs)
}

func demo() {
	demo := func() {
		n := int(1e7)
		fs := NewFastSet(n)
		time1 := time.Now()
		for i := 0; i < n; i++ {
			fs.Insert(i)
			fs.Next(i)
			fs.Prev(i)
			fs.Has(i)
			fs.Erase(i)
			fs.Insert(i)
		}
		fmt.Println(time.Since(time1))
	}
	_ = demo

	// https://judge.yosupo.jp/problem/predecessor_problem
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	const N = 1e7 + 10
	set := NewFastSet(N)
	var s string
	fmt.Fscan(in, &s)
	for i, v := range s {
		if v == '1' {
			set.Insert(i)
		}
	}

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		switch op {
		case 0:
			var k int
			fmt.Fscan(in, &k)
			set.Insert(k)
		case 1:
			var k int
			fmt.Fscan(in, &k)
			set.Erase(k)
		case 2:
			var k int
			fmt.Fscan(in, &k)
			if set.Has(k) {
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintln(out, 0)
			}
		case 3:
			var k int
			fmt.Fscan(in, &k)
			ceiling := set.Next(k)
			if ceiling < N {
				fmt.Fprintln(out, ceiling)
			} else {
				fmt.Fprintln(out, -1)
			}
		case 4:
			var k int
			fmt.Fscan(in, &k)
			floor := set.Prev(k)
			fmt.Fprintln(out, floor)

		}
	}
}

type FastSet struct {
	n, lg int
	seg   [][]int
}

func NewFastSet(n int) *FastSet {
	res := &FastSet{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func NewFastSetFrom(n int, f func(i int) bool) *FastSet {
	res := NewFastSet(n)
	for i := 0; i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
		}
	}
	for h := 0; h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *FastSet) Has(i int) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet) Insert(i int) bool {
	if fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	return true
}

func (fs *FastSet) Erase(i int) bool {
	if !fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == len(cache) {
			break
		}
		d := cache[i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *FastSet) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *FastSet) Enumerate(start, end int, f func(i int)) {
	x := start - 1
	for {
		x = fs.Next(x + 1)
		if x >= end {
			break
		}
		f(x)
	}
}

func (fs *FastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *FastSet) Size() int {
	return fs.n
}

func (*FastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*FastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}
