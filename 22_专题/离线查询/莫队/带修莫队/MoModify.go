// 带有时间序列的莫队,时间复杂度O(n^5/3)
// https://maspypy.github.io/library/ds/offline_query/mo_3d.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/mo.go

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
	"sort"
)

func main() {
	// https://www.luogu.com.cn/problem/P1903
	// Q L R 查询第L支画笔到第R支画笔中共有几种不同颜色的画笔。
	// R P Col 把第P支画笔替换为颜色 Col

	// n,q<=1e5
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	mo := NewMoModify(nums, q)
	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op == "Q" {
			var l, r int
			fmt.Fscan(in, &l, &r)
			l--
			mo.AddQuery(l, r)
		} else {
			var p, col int
			fmt.Fscan(in, &p, &col)
			p--
			mo.AddModify(p, col)
		}
	}

	const N int = 1e6 + 10
	counter, kind := [N]int{}, 0
	res := make([]int, q)
	for i := range res {
		res[i] = -1
	}
	add := func(value T) {
		if counter[value] == 0 {
			kind++
		}
		counter[value]++
	}
	remove := func(value T) {
		counter[value]--
		if counter[value] == 0 {
			kind--
		}
	}
	query := func(qid int) {
		res[qid] = kind
	}

	mo.Run(add, remove, query)
	for _, v := range res {
		if v != -1 {
			fmt.Fprintln(out, v)
		}
	}
}

type T = int
type MoModify struct {
	chunkSize int
	nums      []T
	queries   []query
	modifies  []modify
}

type query struct{ leftBlock, rightBlock, left, right, time, qid int }
type modify struct{ pos, val int }

func NewMoModify(nums []T, q int) *MoModify {
	n := len(nums)
	nums_ := make([]T, n+1)
	copy(nums_[1:], nums)
	chunkSize := int(math.Round(math.Pow(float64(n), 2.0/3)))
	// chunkSize := max(1, n/max(1, int(math.Sqrt(float64(q*2/3)))))
	return &MoModify{chunkSize: chunkSize, nums: nums_}
}

// 添加一个查询，查询范围为`左闭右开区间` [left, right).
//  0 <= left <= right <= n.
func (mo *MoModify) AddQuery(left, right int) {
	left++
	mo.queries = append(
		mo.queries,
		query{left / mo.chunkSize, (right + 1) / mo.chunkSize, left, (right + 1), len(mo.modifies), len(mo.queries)},
	)
}

// 添加一个修改，修改位置为 pos, 修改值为 val.
//  0 <= pos < n.
func (mo *MoModify) AddModify(pos int, val T) {
	pos++
	mo.modifies = append(mo.modifies, modify{pos, val})
}

// 返回每个查询的结果.
//  add: 将数据添加到窗口. delta: 1 表示向右移动，-1 表示向左移动.
//  remove: 将数据从窗口移除. delta: 1 表示向右移动，-1 表示向左移动.
//  query: 查询窗口内的数据.
func (mo *MoModify) Run(
	add func(value T),
	remove func(value T),
	query func(qid int),
) {
	sort.Slice(mo.queries, func(i, j int) bool {
		a, b := mo.queries[i], mo.queries[j]
		if a.leftBlock != b.leftBlock {
			return a.leftBlock < b.leftBlock
		}
		if a.rightBlock != b.rightBlock {
			if a.leftBlock&1 == 0 {
				return a.rightBlock < b.rightBlock
			}
			return a.rightBlock > b.rightBlock
		}
		if a.rightBlock&1 == 0 {
			return a.time < b.time
		}
		return a.time > b.time
	})

	left, right, now := 1, 1, 0

	for _, q := range mo.queries {
		for ; right < q.right; right++ {
			add(mo.nums[right])
		}
		for ; left < q.left; left++ {
			remove(mo.nums[left])
		}
		for left > q.left {
			left--
			add(mo.nums[left])
		}
		for right > q.right {
			right--
			remove(mo.nums[right])
		}
		for ; now < q.time; now++ {
			m := mo.modifies[now]
			p, v := m.pos, m.val
			if q.left <= p && p < q.right {
				remove(mo.nums[p])
				add(v)
			}
			mo.nums[p], mo.modifies[now].val = v, mo.nums[p]
		}
		for now > q.time {
			now--
			m := mo.modifies[now]
			p, v := m.pos, m.val
			if q.left <= p && p < q.right {
				remove(mo.nums[p])
				add(v)
			}
			mo.nums[p], mo.modifies[now].val = v, mo.nums[p]
		}

		query(q.qid)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
