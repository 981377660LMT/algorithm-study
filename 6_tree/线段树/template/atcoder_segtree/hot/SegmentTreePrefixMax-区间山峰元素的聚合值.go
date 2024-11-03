package main

import "fmt"

func main() {
	heights := []int{2, 1, 3, 3, 1}
	sums := []int{2, 3, 6, 10, 11}
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	seg := NewPrefixMaxSegTree(e, op)
	seg.Build(int32(len(heights)), func(i int32) (int, int) {
		return heights[i], sums[i]
	})

	fmt.Println(seg.Query(0, 5)) // 32
	fmt.Println(seg.QueryAll())  // 11
	fmt.Println(seg.GetAll())    // [2 1 3 3 1] [2 3 6 10 11]
	seg.Set(1, 2, 4)
	fmt.Println(seg.Query(1, 4)) // 20
}

const INF int = 1e18

type Data[E any] struct {
	max       int
	sum, rsum E
}

// 山峰元素的定义为：作为前缀最大值的元素(从左向右看可以看到.)
// 求区间山峰元素的聚合值.
// O(qlog^2n).
type PrefixMaxSegTree[E any] struct {
	n, size, log int32
	e            func() E
	op           func(E, E) E
	data         []Data[E]
}

func NewPrefixMaxSegTree[E any](e func() E, op func(E, E) E) *PrefixMaxSegTree[E] {
	return &PrefixMaxSegTree[E]{e: e, op: op}
}

func (seg *PrefixMaxSegTree[E]) Build(n int32, f func(i int32) (int, E)) {
	log := int32(1)
	for 1<<log < n {
		log++
	}
	size := int32(1 << log)
	data := make([]Data[E], size<<1)
	for i := range data {
		data[i] = Data[E]{-INF, seg.e(), seg.e()}
	}
	for i := int32(0); i < n; i++ {
		k, x := f(i)
		data[size+i] = Data[E]{k, x, seg.e()}
	}
	seg.n, seg.size, seg.log, seg.data = n, size, log, data
	for i := size - 1; i > 0; i-- {
		seg.update(i)
	}
}

func (seg *PrefixMaxSegTree[E]) Set(i int32, key int, e E) {
	i += seg.size
	seg.data[i] = Data[E]{key, e, seg.e()}
	for i > 1 {
		i >>= 1
		seg.update(i)
	}
}

func (seg *PrefixMaxSegTree[E]) QueryAll() E {
	return seg.data[1].sum
}

func (seg *PrefixMaxSegTree[E]) Query(l, r int32) E {
	k := -INF
	var suf []int32
	l += seg.size
	r += seg.size
	res := seg.e()
	for l < r {
		if l&1 == 1 {
			res = seg.op(res, seg.dfs(l, k))
			k = max(k, seg.data[l].max)
			l++
		}
		if r&1 == 1 {
			r--
			suf = append(suf, r)
		}
		l >>= 1
		r >>= 1
	}
	for i := len(suf) - 1; i >= 0; i-- {
		res = seg.op(res, seg.dfs(suf[i], k))
		k = max(k, seg.data[suf[i]].max)
	}
	return res
}

func (seg *PrefixMaxSegTree[E]) Get(i int32) (int, E) {
	return seg.data[seg.size+i].max, seg.data[seg.size+i].sum
}

func (seg *PrefixMaxSegTree[E]) GetAll() ([]int, []E) {
	key := make([]int, seg.n)
	x := make([]E, seg.n)
	for i := int32(0); i < seg.n; i++ {
		key[i] = seg.data[seg.size+i].max
		x[i] = seg.data[seg.size+i].sum
	}
	return key, x
}

func (seg *PrefixMaxSegTree[E]) update(i int32) {
	seg.data[i].max = max(seg.data[i<<1].max, seg.data[i<<1|1].max)
	seg.data[i].rsum = seg.dfs(i<<1|1, seg.data[i<<1].max)
	seg.data[i].sum = seg.op(seg.data[i<<1].sum, seg.data[i].rsum)
}

func (seg *PrefixMaxSegTree[E]) dfs(v int32, k int) E {
	if seg.size <= v {
		if k <= seg.data[v].max {
			return seg.data[v].sum
		}
		return seg.e()
	}
	if k <= seg.data[v<<1].max {
		return seg.op(seg.dfs(v<<1, k), seg.data[v].rsum)
	}
	return seg.dfs(v<<1|1, k)
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
