package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	yosupo()
}

// https://judge.yosupo.jp/problem/range_parallel_unionfind
// 并行合并，查询联通的二元组(i,j)的数量模998244353.
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var n, q int32
	fmt.Fscan(in, &n, &q)
	uf := NewRangeParallelUnionFindOnline(n)
	nums := make([]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	res := 0
	f := func(big, small int32) {
		res += nums[big] * nums[small]
		res %= MOD
		nums[big] += nums[small]
		nums[big] %= MOD
	}

	for i := int32(0); i < q; i++ {
		var len, a, b int32
		fmt.Fscan(in, &len, &a, &b)
		uf.UnionRange(a, a+len, b, b+len, f)
		fmt.Fprintln(out, res)
	}
}

type RangeParallelUnionFindOnline struct {
	n   int32
	log int32
	// ufs[i][a]==ufs[i][b] <=> [a,...,a+2^i) == [b,...,b+2^i)
	ufs []*unionFindArraySimple32
}

func NewRangeParallelUnionFindOnline(n int32) *RangeParallelUnionFindOnline {
	log := int32(1)
	for 1<<log < n {
		log++
	}
	ufs := make([]*unionFindArraySimple32, log)
	for i := int32(0); i < log; i++ {
		ufs[i] = newUnionFindArraySimple32(n - (1 << i) + 1)
	}
	return &RangeParallelUnionFindOnline{n: n, log: log, ufs: ufs}
}

func (uf *RangeParallelUnionFindOnline) Union(i, j int32, f func(big, small int32)) {
	uf.unionInner(0, i, j, f)
}

func (uf *RangeParallelUnionFindOnline) UnionRange(s1, e1 int32, s2, e2 int32, f func(big, small int32)) {
	if e1-s1 != e2-s2 {
		panic("invalid")
	}
	n := e1 - s1
	if n == 0 {
		return
	}
	if n == 1 {
		uf.unionInner(0, s1, s2, f)
		return
	}
	k := topbit(n - 1)
	uf.unionInner(k, s1, s2, f)
	uf.unionInner(k, e1-(1<<k), e2-(1<<k), f)
}

func (uf *RangeParallelUnionFindOnline) unionInner(k, l1, l2 int32, f func(big, small int32)) {
	if k == 0 {
		a, b := uf.ufs[0].Find(l1), uf.ufs[0].Find(l2)
		if a == b {
			return
		}
		uf.ufs[0].Union(a, b)
		c := uf.ufs[0].Find(a)
		f(c, a^b^c)
		return
	}
	if !uf.ufs[k].Union(l1, l2) {
		return
	}
	uf.unionInner(k-1, l1, l2, f)
	uf.unionInner(k-1, l1+(1<<(k-1)), l2+(1<<(k-1)), f)
}

type unionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func newUnionFindArraySimple32(n int32) *unionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &unionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *unionFindArraySimple32) Union(key1, key2 int32) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = int32(root1)
	u.Part--
	return true
}

func (u *unionFindArraySimple32) Find(key int32) int32 {
	if u.data[key] < 0 {
		return key
	}
	u.data[key] = u.Find(u.data[key])
	return u.data[key]
}

func (u *unionFindArraySimple32) GetSize(key int32) int32 {
	return -u.data[u.Find(key)]
}

func topbit(x int32) int32 {
	if x == 0 {
		return -1
	}
	return int32(31 - bits.LeadingZeros32(uint32(x)))
}
