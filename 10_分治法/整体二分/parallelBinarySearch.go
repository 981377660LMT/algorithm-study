// parallelBinarySearch/parallelSortSearch
// 並列二分探索
//
// https://betrue12.hateblo.jp/entry/2019/08/14/152227
// https://ouuan.github.io/post/parallel-binary-search/
// https://maspypy.github.io/library/ds/offline_query/parallel_binary_search.hpp
// https://oi-wiki.org/misc/parallel-binsearch/
//
// 整体二分的主体思路就是把`多个查询`一起解决。
// !`单个查询可以二分答案解决，但是多个查询就会TLE`的这种场合，就可以考虑整体二分。
// 整体二分解决这样一类问题:
//   - 给定一个长度为n的操作序列, 按顺序执行这些操作;
//   - 给定q个查询,每个查询形如:"条件qi为真(满足条件)是在第几次操作之后?".
//     !要求对条件为真的判定具有单调性，即某个操作后qi为真,后续操作都会满足qi为真.
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
	// demo()
	StaticRangeKthSmallest()
}

// 静态区间第 k 小
// https://www.luogu.com.cn/problem/P3834
func StaticRangeKthSmallest() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	queries := make([][3]int, q) // [l, r) 中第 k 小 (1-indexed)
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
		return bit.ProdRange(l, r) >= k
	}
	ok := ParallelBinarySearch(n, q, update, init, check)
	for i := 0; i < q; i++ {
		t := ok[i]
		fmt.Fprintln(out, nums[I[t]])
	}
}

func demo() {
	// n个操作，第i个操作将curSum增加i+1.
	// q个查询，第i个查询形如：curSum是否大于等于i+1.
	// 对于每个查询，输出第一个满足条件的操作的编号.

	curSum := 0
	res := ParallelBinarySearch(
		10, 10,
		func(mutationId int) {
			curSum += mutationId + 1
		},
		func() {
			curSum = 0
		},
		func(queryId int) bool {
			return curSum >= 560
		},
	)

	fmt.Println(res)
}

// 整体二分解决这样一类问题:
//   - 给定一个长度为n的操作序列, 按顺序执行这些操作;
//   - 给定q个查询,每个查询形如:"条件qi为真(满足条件)是在第几次操作之后?".
//     !要求对条件为真的判定具有单调性，即某个操作后qi为真,后续操作都会满足qi为真.
//
// 返回:
//   - -1 => 不需要操作就满足条件的查询.
//   - [0, n) => 满足条件的最早的操作的编号(0-based).
//   - n => 执行完所有操作后都不满足条件的查询.
//
// https://betrue12.hateblo.jp/entry/2019/08/14/152227
func ParallelBinarySearch(
	n, q int,
	mutate func(mutationId int), // 执行第 mutationId 次操作，一共调用 nlogn 次.
	reset func(), // 重置操作序列，一共调用 logn 次.
	predicate func(queryId int) bool, // 判断第 queryId 次查询是否满足条件，一共调用 qlogn 次.
) []int {
	left, right := make([]int, q), make([]int, q)
	for i := 0; i < q; i++ {
		left[i], right[i] = 0, n
	}

	// 不需要操作就满足条件的查询
	for i := 0; i < q; i++ {
		if predicate(i) {
			right[i] = -1
		}
	}

	for {
		mids := make([]int, q)
		for i := range mids {
			mids[i] = -1
		}
		for i := 0; i < q; i++ {
			if left[i] <= right[i] {
				mids[i] = (left[i] + right[i]) >> 1
			}
		}

		// csr 数组保存二元对 (qi,mid).
		indeg := make([]int, n+2)
		for i := 0; i < q; i++ {
			mid := mids[i]
			if mid != -1 {
				indeg[mid+1]++
			}
		}
		for i := 0; i < n+1; i++ {
			indeg[i+1] += indeg[i]
		}
		total := indeg[n+1]
		if total == 0 {
			break
		}
		counter := append(indeg[:0:0], indeg...)
		csr := make([]int, total)
		for i := 0; i < q; i++ {
			mid := mids[i]
			if mid != -1 {
				csr[counter[mid]] = i
				counter[mid]++
			}
		}

		reset()
		times := 0
		for _, pos := range csr {
			for times < mids[pos] {
				mutate(times)
				times++
			}
			if predicate(pos) {
				right[pos] = times - 1
			} else {
				left[pos] = times + 1
			}
		}
	}

	return right
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
