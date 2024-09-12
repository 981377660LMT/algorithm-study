package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

// No.2290 UnUnion Find (查询不连通的点)
// https://yukicoder.me/problems/no/2290
// 1 u v: 连接u和v
// 2 v: 输出一个与v不连通的点.如果不存在这样的点,输出-1.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	uf := NewUnionFindArraySimple32(n)
	fs := NewFastSet32From(n, func(i int32) bool { return true }) // 维护哪些点的根是自己

	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var u, v int32
			fmt.Fscan(in, &u, &v)
			u, v = u-1, v-1
			ru, rv := uf.Find(u), uf.Find(v)
			uf.Union(ru, rv, nil)
			if ru != uf.Find(ru) {
				fs.Erase(ru)
			}
			if rv != uf.Find(rv) {
				fs.Erase(rv)
			}
		} else {
			var x int32
			fmt.Fscan(in, &x)
			x--
			min_, max_ := fs.Next(0), fs.Prev(n-1)
			if uf.Find(min_) != uf.Find(x) {
				fmt.Fprintln(out, min_+1)
			} else if uf.Find(max_) != uf.Find(x) {
				fmt.Fprintln(out, max_+1)
			} else {
				fmt.Fprintln(out, -1)
			}
		}
	}
}

type FastSet32 struct {
	n, lg int32
	seg   [][]uint64
	size  int32
}

func NewFastSet32(n int32) *FastSet32 {
	res := &FastSet32{n: n}
	seg := [][]uint64{}
	n_ := n
	for {
		seg = append(seg, make([]uint64, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = int32(len(seg))
	return res
}

func NewFastSet32From(n int32, f func(i int32) bool) *FastSet32 {
	res := NewFastSet32(n)
	for i := int32(0); i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
			res.size++
		}
	}
	for h := int32(0); h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *FastSet32) Has(i int32) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet32) Insert(i int32) bool {
	if fs.Has(i) {
		return false
	}
	for h := int32(0); h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	fs.size++
	return true
}

func (fs *FastSet32) Erase(i int32) bool {
	if !fs.Has(i) {
		return false
	}
	for h := int32(0); h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	fs.size--
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet32) Next(i int32) int32 {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := int32(0); h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == int32(len(cache)) {
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
func (fs *FastSet32) Prev(i int32) int32 {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := int32(0); h < fs.lg; h++ {
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
func (fs *FastSet32) Enumerate(start, end int32, f func(i int32)) {
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *FastSet32) String() string {
	res := []string{}
	for i := int32(0); i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(int(i)))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *FastSet32) Size() int32 {
	return fs.size
}

func (*FastSet32) bsr(x uint64) int32 {
	return 63 - int32(bits.LeadingZeros64(x))
}

func (*FastSet32) bsf(x uint64) int32 {
	return int32(bits.TrailingZeros64(x))
}

type UnionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func NewUnionFindArraySimple32(n int32) *UnionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple32) Union(key1, key2 int32, beforeMerge func(big, small int32)) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	if beforeMerge != nil {
		beforeMerge(root1, root2)
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = root1
	u.Part--
	return true
}

func (u *UnionFindArraySimple32) Find(key int32) int32 {
	root := key
	for u.data[root] >= 0 {
		root = u.data[root]
	}
	for key != root {
		key, u.data[key] = u.data[key], root
	}
	return root
}

func (u *UnionFindArraySimple32) Size(key int32) int32 {
	return -u.data[u.Find(key)]
}
