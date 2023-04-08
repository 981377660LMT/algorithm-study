// NewDisjointSparseTableXor
// XorDisjointSparseTable 异或st表

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://yukicoder.me/problems/no/1891
	// !给定一个长为2的幂次的数组
	// 区间仿射变换,区间查询时每个下标异或上给定的数
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	leaves := make([]S, n)
	for i := range leaves {
		fmt.Fscan(in, &leaves[i].mul, &leaves[i].add)
	}

	seg := NewXorDisjointSparseTable(leaves)

	for i := 0; i < q; i++ {
		var start, end, indexXor, x int
		fmt.Fscan(in, &start, &end, &indexXor, &x)
		res := seg.Query(start, end, indexXor)
		fmt.Fprintln(out, (res.mul*x+res.add)%MOD)
	}
}

// RangeAffineRangeComposite
const MOD int = 998244353

type S = struct{ mul, add int }

func (*XorDisjointSparseTable) e() S { return S{1, 0} }
func (*XorDisjointSparseTable) op(a, b S) S {
	return S{a.mul * b.mul % MOD, (a.add*b.mul + b.add) % MOD}
}

type XorDisjointSparseTable struct {
	log  int
	data [][]S
}

// DisjointSparseTableXor 支持半群的区间静态查询.
//  op:只需要满足结合律 op(op(a,b),c) = op(a,op(b,c)).
//  !区间查询时下标可以异或上 indexXor.
//  !nums 的长度必须要是2的幂.
func NewXorDisjointSparseTable(nums []S) *XorDisjointSparseTable {
	res := &XorDisjointSparseTable{}
	n := len(nums)
	log := 0
	for 1<<log < n {
		log++
	}
	if 1<<log != n {
		panic("len(nums) must be power of 2")
	}
	data := make([][]S, log+1)
	data[0] = make([]S, 1<<log)
	copy(data[0], nums)
	for k := 0; k < log; k++ {
		data[k+1] = make([]S, 1<<log)
		for i := 0; i < 1<<log; i++ {
			data[k+1][i] = res.op(data[k][i], data[k][i^(1<<k)])
		}
	}

	res.log = log
	res.data = data
	return res
}

// Calculate prod_{l<=i<r} A[x xor i], in O(log N) time.
func (st *XorDisjointSparseTable) Query(start, end, indexXor int) S {
	xl, xr := st.e(), st.e()
	for k := 0; k <= st.log; k++ {
		if start >= end {
			break
		}
		if start&(1<<k) != 0 {
			xl = st.op(xl, st.data[k][start^indexXor])
			start += 1 << k
		}
		if end&(1<<k) != 0 {
			end -= 1 << k
			xr = st.op(st.data[k][end^indexXor], xr)
		}
	}
	return st.op(xl, xr)
}
