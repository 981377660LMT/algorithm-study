// https://maspypy.github.io/library/ds/offline_query/parallel_binary_search.hpp
// https://oi-wiki.org/misc/parallel-binsearch/
// https://betrue12.hateblo.jp/entry/2019/08/14/152227
// 整体二分的主体思路就是把`多个查询`一起解决。
// !`单个查询可以二分答案解决，但是多个查询就会TLE`的这种场合，就可以考虑整体二分。
// 整体二分解决这样一类问题:
//  - 给定一个长度为n的操作序列,按顺序执行这些操作。
//  - 给定q个查询,每个查询形如:"条件qi为真(满足条件)是在第几次操作之后?"
//  - 对条件的判定具有单调性，即某个操作后qi为真,后续操作都会满足qi为真。
//
// 一些时候整体二分可以被持久化数据结构取代.

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	// 静态区间第 k 小
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	queries := make([][3]int, q) // [l, r) 中第 k+1 小 (0-indexed)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i][0], &queries[i][1], &queries[i][2])
		queries[i][0]--
	}

	// argsort
	I := make([]int, len(nums))
	for i := 0; i < len(nums); i++ {
		I[i] = i
	}
	sort.Slice(I, func(i, j int) bool {
		return nums[I[i]] < nums[I[j]]
	})

	bit := NewBitArray(n)
	init := func() { bit.Build(make([]int, n)) }
	update := func(t int) { bit.Apply(I[t], 1) }
	check := func(q int) bool {
		l, r, k := queries[q][0], queries[q][1], queries[q][2]
		return bit.ProdRange(l, r) > k-1
	}
	ok := ParallelBinarySearch(q, n, 0, init, update, check)
	for i := 0; i < q; i++ {
		t := ok[i]
		fmt.Fprintln(out, nums[I[t-1]])
	}
}

// 整体二分 (q, ok, ng, init, update, check).
//
//	给定一个`操作序列`和 q 个`查询`.
//	对每个查询，返回使得 check 第一次成立时的 update 的次数.
//	- q: 查询的数量.
//	- ok, ng: 二分的上下界.
//	- init(): 每次二分前的初始化.
//	- update(t): 执行操作序列中的第 t 次操作(0-indexed).
//	- check(q): 对查询 q 的判定.
func ParallelBinarySearch(
	q, ok, ng int,
	init func(), update func(t int), check func(q int) bool,
) []int {
	T := max(ok, ng)
	right, left := make([]int, q), make([]int, q)
	for i := 0; i < q; i++ {
		right[i], left[i] = ok, ng
	}

	for {
		checkT := make([]int, q)
		for i := 0; i < q; i++ {
			checkT[i] = -1
		}

		for i := 0; i < q; i++ {
			if abs(right[i]-left[i]) > 1 {
				checkT[i] = (right[i] + left[i]) / 2
			}
		}

		indeg := make([]int, T+1)
		for i := 0; i < q; i++ {
			t := checkT[i]
			if t != -1 {
				indeg[t+1]++
			}
		}
		for i := 0; i < T; i++ {
			indeg[i+1] += indeg[i]
		}
		total := indeg[T]
		if total == 0 {
			break
		}
		counter := make([]int, T+1)
		copy(counter, indeg)
		csr := make([]int, total)
		for i := 0; i < q; i++ {
			t := checkT[i]
			if t != -1 {
				csr[counter[t]] = i
				counter[t]++
			}
		}

		init()
		t := 0
		for _, q := range csr {
			for t < checkT[q] {
				update(t)
				t++
			}

			if check(q) {
				right[q] = t
			} else {
				left[q] = t
			}
		}
	}

	return right
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type BitArray struct {
	n    int
	log  int
	data []int
}

// 長さ n の 0で初期化された配列で構築する.
func NewBitArray(n int) *BitArray {
	return &BitArray{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

// 配列で構築する.
func NewBitArrayFrom(arr []int) *BitArray {
	res := NewBitArray(len(arr))
	res.Build(arr)
	return res
}

func (b *BitArray) Build(arr []int) {
	if b.n != len(arr) {
		panic("len of arr is not equal to n")
	}
	for i := 1; i <= b.n; i++ {
		b.data[i] = arr[i-1]
	}
	for i := 1; i <= b.n; i++ {
		j := i + (i & -i)
		if j <= b.n {
			b.data[j] += b.data[i]
		}
	}
}

// 要素 i に値 v を加える.
func (b *BitArray) Apply(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

// [0, r) の要素の総和を求める.
func (b *BitArray) Prod(r int) int {
	res := int(0)
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

// [l, r) の要素の総和を求める.
func (b *BitArray) ProdRange(l, r int) int {
	return b.Prod(r) - b.Prod(l)
}

// 区間[0,k]の総和がx以上となる最小のkを求める.数列が単調増加であることを要求する.
func (b *BitArray) LowerBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] < x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

// 区間[0,k]の総和がxを上回る最小のkを求める.数列が単調増加であることを要求する.
func (b *BitArray) UpperBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] <= x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}
