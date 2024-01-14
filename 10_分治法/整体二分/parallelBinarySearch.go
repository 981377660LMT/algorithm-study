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
	"strings"
)

func main() {
	// StaticRangeKthSmallest()
	// 矩阵乘法()
	// 天天爱射击()
	// UnionSets()
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
	order := make([]int, len(nums))
	for i := 0; i < len(nums); i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		return nums[order[i]] < nums[order[j]]
	})

	bit := NewBitArray(n)
	reset := func() {
		bit.Build(make([]int, n))
	}
	mutate := func(id int) {
		bit.Update(order[id], 1)
	}
	predicate := func(qid int) bool {
		l, r, k := queries[qid][0], queries[qid][1], queries[qid][2]
		return bit.QueryRange(l, r) >= k
	}

	left := ParallelBinarySearch(n, q, reset, mutate, predicate)
	for i := 0; i < q; i++ {
		v := left[i]
		fmt.Fprintln(out, nums[order[v]])
	}
}

// 静态二维矩阵第 k 小
// P1527 [国家集训队] 矩阵乘法
// https://www.luogu.com.cn/problem/P1527
func 矩阵乘法() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n*n)
	for i := 0; i < n*n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	queries := make([][5]int, q) // [x1,x2) x [y1,y2) 中第 k 小 (1-indexed)
	for i := 0; i < q; i++ {
		var x1, y1, x2, y2, k int
		fmt.Fscan(in, &x1, &y1, &x2, &y2, &k)
		x1--
		y1--
		queries[i] = [5]int{x1, x2, y1, y2, k}
	}

	// argsort
	order := make([]int, n*n)
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		return nums[order[i]] < nums[order[j]]
	})

	bit := NewBIT2DDense(n, n)
	reset := func() {
		bit.Build(func(x, y int) int { return 0 })
	}
	mutate := func(id int) {
		bit.Add(order[id]/n, order[id]%n, 1)
	}
	predicate := func(qid int) bool {
		item := queries[qid]
		x1, x2, y1, y2, k := item[0], item[1], item[2], item[3], item[4]
		return bit.QueryRange(x1, x2, y1, y2) >= k
	}

	res := ParallelBinarySearch(n*n, q, reset, mutate, predicate)
	for i := 0; i < q; i++ {
		v := res[i]
		fmt.Fprintln(out, nums[order[v]])
	}
}

// P7424 [THUPC2017] 天天爱射击
// https://www.luogu.com.cn/problem/P7424
// 将【每个子弹射出之后有多少个木板会碎掉】转化为【每个木板会在第几次子弹射击之后碎掉】
// 子弹按照发射顺序给出
func 天天爱射击() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var board, bullet int
	fmt.Fscan(in, &board, &bullet)

	boardInfo := make([][3]int, board) // [left,right,hp]
	for i := 0; i < board; i++ {
		var left, right, hp int
		fmt.Fscan(in, &left, &right, &hp)
		boardInfo[i] = [3]int{left, right, hp}
	}

	bulletInfo := make([]int, bullet) // pos
	for i := 0; i < bullet; i++ {
		fmt.Fscan(in, &bulletInfo[i])
	}

	const N = 2e5 + 10
	bit := NewBitArray(N) // 每个子弹击中的位置计数
	reset := func() {
		bit = NewBitArray(N)
	}

	mutate := func(id int) {
		pos := bulletInfo[id]
		bit.Update(pos, 1)
	}

	predicate := func(qid int) bool {
		item := boardInfo[qid]
		left, right, hp := item[0], item[1], item[2]
		return bit.QueryRange(left, right+1) >= hp
	}

	// bullet 次修改，board 次查询
	left := ParallelBinarySearch(bullet, board, reset, mutate, predicate)
	counter := make([]int, bullet)
	for _, v := range left {
		if v < 0 || v >= bullet {
			continue
		}
		counter[v]++
	}
	for _, v := range counter {
		fmt.Fprintln(out, v)
	}
}

// https://atcoder.jp/contests/code-thanks-festival-2017-open/tasks/code_thanks_festival_2017_h
// 给定n个集合,初始时第i个集合只有一个元素i (i=1,2,...,n)
// 之后进行m次合并操作,每次合并ai和bi所在的集合
// 如果ai和bi在同一个集合,则无事发生
// 给定q个询问,问ai和bi是在第几次操作后第一次连通的,如果不连通则输出-1
func UnionSets() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	mutaions := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &mutaions[i][0], &mutaions[i][1])
		mutaions[i][0]--
		mutaions[i][1]--
	}
	var q int
	fmt.Fscan(in, &q)
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i][0], &queries[i][1])
		queries[i][0]--
		queries[i][1]--
	}

	uf := NewUnionFindArray(n)
	reset := func() {
		uf = NewUnionFindArray(n)
	}
	mutate := func(id int) {
		mutation := mutaions[id]
		u, v := mutation[0], mutation[1]
		uf.Union(u, v)
	}
	predicate := func(qid int) bool {
		query := queries[qid]
		u, v := query[0], query[1]
		return uf.Find(u) == uf.Find(v)
	}

	res := ParallelBinarySearch(m, q, reset, mutate, predicate)
	for i := 0; i < q; i++ {
		v := res[i]
		if v == m {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, v+1)
		}
	}
}

func demo() {
	// n个操作，第i个操作将curSum增加i+1.
	// q个查询，第i个查询形如：curSum是否大于等于i+1.
	// 对于每个查询，输出第一个满足条件的操作的编号.

	curSum := 0
	res := ParallelBinarySearch(
		10, 10,
		func() {
			curSum = 0
		},
		func(mutationId int) {
			curSum += mutationId + 1
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
	reset func(), // 重置操作序列，一共调用 logn 次.
	mutate func(mutationId int), // 执行第 mutationId 次操作，一共调用 nlogn 次.
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
func (b *BitArray) Update(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

// [0, r) の要素の総和を求める.
func (b *BitArray) QueryPrefix(r int) int {
	res := int(0)
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

// [l, r) の要素の総和を求める.
func (b *BitArray) QueryRange(l, r int) int {
	return b.QueryPrefix(r) - b.QueryPrefix(l)
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

type BIT2DPointAddRangeSum struct {
	H, W int
	data []int
}

// 指定二维树状数组的高度和宽度.
func NewBIT2DDense(h, w int) *BIT2DPointAddRangeSum {
	res := &BIT2DPointAddRangeSum{H: h, W: w, data: make([]int, h*w)}
	return res
}

func (ft *BIT2DPointAddRangeSum) Build(f func(x, y int) int) {
	H, W := ft.H, ft.W
	for x := 0; x < H; x++ {
		for y := 0; y < W; y++ {
			ft.data[W*x+y] = f(x, y)
		}
	}
	for x := 1; x <= H; x++ {
		for y := 1; y <= W; y++ {
			ny := y + (y & -y)
			if ny <= W {
				ft.data[ft.idx(x, ny)] += ft.data[ft.idx(x, y)]
			}
		}
	}
	for x := 1; x <= H; x++ {
		for y := 1; y <= W; y++ {
			nx := x + (x & -x)
			if nx <= H {
				ft.data[ft.idx(nx, y)] += ft.data[ft.idx(x, y)]
			}
		}
	}
}

// 点 (x,y) 的值加上 val.
func (ft *BIT2DPointAddRangeSum) Add(x, y int, val int) {
	x++
	for x <= ft.H {
		ft.addX(x, y, val)
		x += x & -x
	}
}

// [lx,rx) * [ly,ry)
func (ft *BIT2DPointAddRangeSum) QueryRange(lx, rx, ly, ry int) int {
	pos, neg := 0, 0
	for lx < rx {
		pos += ft.sumX(rx, ly, ry)
		rx -= rx & -rx
	}
	for rx < lx {
		neg += ft.sumX(lx, ly, ry)
		lx -= lx & -lx
	}
	return pos - neg
}

// [0,rx) * [0,ry)
func (ft *BIT2DPointAddRangeSum) QueryPrefix(rx, ry int) int {
	pos := 0
	for rx > 0 {
		pos += ft.sumXPrefix(rx, ry)
		rx -= rx & -rx
	}
	return pos
}

func (ft *BIT2DPointAddRangeSum) String() string {
	res := make([][]string, ft.H)
	for i := 0; i < ft.H; i++ {
		res[i] = make([]string, ft.W)
		for j := 0; j < ft.W; j++ {
			res[i][j] = fmt.Sprintf("%v", ft.QueryRange(i, i+1, j, j+1))
		}

	}
	return fmt.Sprintf("%v", res)
}

func (ft *BIT2DPointAddRangeSum) idx(x, y int) int {
	return ft.W*(x-1) + (y - 1)
}

func (ft *BIT2DPointAddRangeSum) addX(x, y int, val int) {
	y++
	for y <= ft.W {
		ft.data[ft.idx(x, y)] += val
		y += y & -y
	}
}

func (ft *BIT2DPointAddRangeSum) sumX(x, ly, ry int) int {
	pos, neg := 0, 0
	for ly < ry {
		pos += ft.data[ft.idx(x, ry)]
		ry -= ry & -ry
	}
	for ry < ly {
		neg += ft.data[ft.idx(x, ly)]
		ly -= ly & -ly
	}
	return pos - neg
}

func (ft *BIT2DPointAddRangeSum) sumXPrefix(x, ry int) int {
	pos := 0
	for ry > 0 {
		pos += ft.data[ft.idx(x, ry)]
		ry -= ry & -ry
	}
	return pos
}

func DiscretizeCompressed(nums []int, offset int) (getRank func(int) int, count int) {
	set := make(map[int]struct{}, len(nums))
	for _, v := range nums {
		set[v] = struct{}{}
	}
	count = len(set)
	rank := make([]int, 0, count)
	for v := range set {
		rank = append(rank, v)
	}
	sort.Ints(rank)
	mp := make(map[int]int, count)
	for i, v := range rank {
		mp[v] = i + offset
	}
	getRank = func(v int) int { return mp[v] }
	count = len(nums)
	return
}

type BITRangeAddRangeSum struct {
	n     int
	tree1 []int
	tree2 []int
}

func NewBITRangeAddRangeSum(n int) *BITRangeAddRangeSum {
	return &BITRangeAddRangeSum{
		n:     n,
		tree1: make([]int, n+1),
		tree2: make([]int, n+1),
	}
}

// 切片内[start, end)的每个元素加上delta.
//
//	0<=start<=end<=n
func (b *BITRangeAddRangeSum) AddRange(start, end, delta int) {
	end--
	b._add(start, delta)
	b._add(end+1, -delta)
}

func (b *BITRangeAddRangeSum) QueryPrefix(index int) (res int) {
	index++
	if index > b.n {
		index = b.n
	}
	for i := index; i > 0; i &= i - 1 {
		res += index*b.tree1[i] - b.tree2[i]
	}
	return res
}

// 求切片内[start, end)的和.
//
//	0<=start<=end<=n
func (b *BITRangeAddRangeSum) QueryRange(start, end int) int {
	end--
	return b.QueryPrefix(end) - b.QueryPrefix(start-1)
}

func (b *BITRangeAddRangeSum) String() string {
	res := []string{}
	for i := 0; i < b.n; i++ {
		res = append(res, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BITRangeAddRangeSum: [%v]", strings.Join(res, ", "))
}

func (b *BITRangeAddRangeSum) _add(index, delta int) {
	index++
	for i := index; i <= b.n; i += i & -i {
		b.tree1[i] += delta
		b.tree2[i] += (index - 1) * delta
	}
}

type UnionFindArray struct {
	data []int
}

func NewUnionFindArray(n int) *UnionFindArray {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArray{data: data}
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1, root2 = root2, root1
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

func (ufa *UnionFindArray) Clear() {
	for i := range ufa.data {
		ufa.data[i] = -1
	}
}
