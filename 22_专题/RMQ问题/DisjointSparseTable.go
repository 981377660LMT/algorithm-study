package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/staticrmq
	// 5e5 - 750ms
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]E, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}
	st := NewDisjointSparse(nums)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		fmt.Fprintln(out, st.Query(l, r))
	}
}

func demo() {
	nums := []E{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	st := NewDisjointSparse(nums)
	fmt.Println(st.Query(0, 10))
	fmt.Println(st.MinLeft(10, func(e E) bool { return e >= 5 }))
}

const INF int = 1e18

type E = int

func (*DisjointSparseTable) e() E { return 0 }
func (*DisjointSparseTable) op(e1, e2 E) E {
	return gcd(e1, e2)
}
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type DisjointSparseTable struct {
	n    int
	data [][]E
}

// DisjointSparseTable 支持幺半群的区间静态查询.
//  eg: 区间乘积取模/区间仿射变换...
func NewDisjointSparse(leaves []E) *DisjointSparseTable {
	res := &DisjointSparseTable{}
	n := len(leaves)
	log := 1
	for (1 << log) < n {
		log++
	}
	data := make([][]E, log)
	data[0] = append(data[0], leaves...)
	for i := 1; i < log; i++ {
		data[i] = append(data[i], data[0]...)
		v := data[i]
		b := 1 << i
		for m := b; m <= n; m += 2 * b {
			l, r := m-b, min(m+b, n)
			for j := m - 1; j >= l+1; j-- {
				v[j-1] = res.op(v[j-1], v[j])
			}
			for j := m; j < r-1; j++ {
				v[j+1] = res.op(v[j], v[j+1])
			}
		}
	}
	res.n = n
	res.data = data
	return res
}

func (ds *DisjointSparseTable) Query(start, end int) E {
	if start >= end {
		return ds.e()
	}
	end--
	if start == end {
		return ds.data[0][start]
	}
	k := 31 - bits.LeadingZeros32(uint32(start^end))
	return ds.op(ds.data[k][start], ds.data[k][end])
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (ds *DisjointSparseTable) MaxRight(left int, check func(e E) bool) int {
	if left == ds.n {
		return ds.n
	}
	ok, ng := left, ds.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(ds.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (ds *DisjointSparseTable) MinLeft(right int, check func(e E) bool) int {
	if right == 0 {
		return 0
	}
	ok, ng := right, -1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(ds.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
