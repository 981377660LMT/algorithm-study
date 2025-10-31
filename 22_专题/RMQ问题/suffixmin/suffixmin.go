package main

import (
	"fmt"
	"math"
)

func main() {
	// 按右端点在线回答 RMQ [l,r] 的最小值
	// 思路：r 从小到大扫；每到 r，把 a[r] 用 Set 放入结构；
	// 当有查询的右端点等于 r 时，直接用 Query(l) 得到 min(a[l..r])。
	arr := []int{5, 3, 7, 2, 6}
	type Q struct{ L, R, ID int }
	qs := []Q{
		{L: 0, R: 2, ID: 0}, // min(arr[0..2]) = min(5,3,7) = 3
		{L: 1, R: 3, ID: 1}, // min(arr[1..3]) = min(3,7,2) = 2
		{L: 3, R: 4, ID: 2}, // min(arr[3..4]) = min(2,6)   = 2
	}

	buckets := make([][]Q, len(arr))
	for _, q := range qs {
		buckets[q.R] = append(buckets[q.R], q)
	}

	res := make([]int, len(qs))
	sm := NewSuffixMin(len(arr), func(a, b int) bool { return a < b }, math.MaxInt/2)
	for r := 0; r < len(arr); r++ {
		sm.Set(r, arr[r])
		for _, q := range buckets[r] {
			res[q.ID] = sm.Query(q.L)
		}
	}
	for i, q := range qs {
		fmt.Printf("min(arr[%d..%d]) = %d\n", q.L, q.R, res[i])
	}
}

// SuffixMin 维护“末尾代入”的后缀最小值结构：
// - 按下标递增调用 Set(i, x)；
// - Query(i) 返回 min(a[i..p-1])，其中 p 为已设置的元素个数。
type SuffixMin[T any] struct {
	n       int
	p       int
	uf      *U
	res     []T
	st      []pair[T] // (value, uf root)
	less    func(a, b T) bool
	neutral T
}

type pair[T any] struct {
	val  T
	root int
}

func NewSuffixMin[T any](n int, less func(a, b T) bool, neutral T) *SuffixMin[T] {
	sm := &SuffixMin[T]{
		n:       n,
		uf:      NewUnionFindArray(n + 1),
		res:     make([]T, n+1),
		st:      make([]pair[T], 0, n),
		less:    less,
		neutral: neutral,
	}
	for i := 0; i <= n; i++ {
		sm.res[i] = neutral
	}
	return sm
}

// Set 设定 a[i]=x（必须按 i 从 0..N-1 递增调用）。
func (sm *SuffixMin[T]) Set(i int, x T) {
	if sm.p != i {
		panic("SuffixMin.Set: i must equal current p (append-only)")
	}
	for len(sm.st) > 0 && !sm.less(sm.st[len(sm.st)-1].val, x) {
		sm.uf.Union(sm.p, sm.st[len(sm.st)-1].root)
		sm.st = sm.st[:len(sm.st)-1]
	}
	r := sm.uf.Find(sm.p)
	sm.p++
	sm.res[r] = x
	sm.st = append(sm.st, pair[T]{val: x, root: r})
}

func (sm *SuffixMin[T]) Query(start int) T {
	if start > sm.p {
		panic("SuffixMin.Query: i must be <= p")
	}
	return sm.res[sm.uf.Find(start)]
}

func (sm *SuffixMin[T]) Reset() {
	sm.uf.Reset()
	for i := 0; i <= sm.n; i++ {
		sm.res[i] = sm.neutral
	}
	sm.st = sm.st[:0]
	sm.p = 0
}

type U struct {
	n    int
	data []int
}

func NewUnionFindArray(n int) *U {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &U{
		n:    n,
		data: data,
	}
}

// 按秩合并.
func (ufa *U) Union(a, b int) bool {
	ra, rb := ufa.Find(a), ufa.Find(b)
	if ra == rb {
		return false
	}
	if ufa.data[ra] > ufa.data[rb] {
		ra, rb = rb, ra
	}
	ufa.data[ra] += ufa.data[rb]
	ufa.data[rb] = ra
	return true
}

func (ufa *U) Find(x int) int {
	r := x
	for ufa.data[r] >= 0 {
		r = ufa.data[r]
	}
	for x != r {
		x, ufa.data[x] = ufa.data[x], r
	}
	return r
}

func (ufa *U) Reset() {
	for i := 0; i < ufa.n; i++ {
		ufa.data[i] = -1
	}
}
