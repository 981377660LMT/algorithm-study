
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int = 1e18

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	xs := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &xs[i])
		xs[i]--
	}

	wm := NewSegmentTreeFractionalCascading(xs)
	res := n * (n - 1)
	for i := 0; i < n; i++ {
		res -= wm.Query(0, i, 0, xs[i])       // 左侧不相撞的
		res -= wm.Query(i+1, n, xs[i]+1, 2*n) // 右侧不相撞的
	}

	fmt.Fprintln(out, res/2)
}

type SegmentTreeFractionalCascading struct {
	seg, ll, rr [][]int
	sz          int
}

func NewSegmentTreeFractionalCascading(array []int) *SegmentTreeFractionalCascading {
	res := &SegmentTreeFractionalCascading{}
	n := len(array)
	sz := 1
	for sz < n {
		sz <<= 1
	}
	tmp := 2*sz - 1
	seg := make([][]int, tmp)
	ll := make([][]int, tmp)
	rr := make([][]int, tmp)
	for k := 0; k < n; k++ {
		seg[k+sz-1] = append(seg[k+sz-1], array[k])
	}
	for k := sz - 2; k >= 0; k-- {
		a, b := 2*k+1, 2*k+2
		seg[k] = make([]int, len(seg[a])+len(seg[b]))
		ll[k] = make([]int, len(seg[k])+1)
		rr[k] = make([]int, len(seg[k])+1)
		seg[k] = append(seg[a], seg[b]...)
		tail1, tail2 := 0, 0
		for i := 0; i < len(seg[k]); i++ {
			for tail1 < len(seg[a]) && seg[a][tail1] < seg[k][i] {
				tail1++
			}
			for tail2 < len(seg[b]) && seg[b][tail2] < seg[k][i] {
				tail2++
			}
			ll[k][i] = tail1
			rr[k][i] = tail2
		}
		ll[k][len(seg[k])] = len(seg[a])
		rr[k][len(seg[k])] = len(seg[b])
	}
	res.seg = seg
	res.ll = ll
	res.rr = rr
	res.sz = sz
	return res
}

// 查询区间 [start, end) 中，[floor, ceiling) 范围内的数的个数.
func (st *SegmentTreeFractionalCascading) Query(start, end, floor, ceiling int) int {
	floor = sort.SearchInts(st.seg[0], floor)
	ceiling = sort.SearchInts(st.seg[0], ceiling)
	return st._query(start, end, floor, ceiling, 0, 0, st.sz)
}

func (st *SegmentTreeFractionalCascading) _query(a, b, lower, upper, k, l, r int) int {
	if a >= r || b <= l {
		return 0
	}
	if a <= l && r <= b {
		return upper - lower
	}
	return st._query(a, b, st.ll[k][lower], st.ll[k][upper], 2*k+1, l, (l+r)>>1) + st._query(a, b, st.rr[k][lower], st.rr[k][upper], 2*k+2, (l+r)>>1, r)
}
