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
	// StampRally()
	CF1706E()
	// StaticRangeKthSmallest()
	// 矩阵乘法()
	// 天天爱射击()
	// UnionSets()
}

// https://atcoder.jp/contests/agc002/tasks/agc002_d
// 一张连通图，q 次询问从两个点 x 和 y 出发，
// 希望经过的点数量等于 z（每个点可以重复经过，但是重复经过只计算一次）
// 求经过的边最大编号最小是多少。
//
// 整体二分，这里的操作序列为：将边按照编号从小到大加入到图中。
func StampRally() []int {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		edges[i] = [2]int{u, v}
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([][3]int, q)
	for i := 0; i < q; i++ {
		var x, y, z int
		fmt.Fscan(in, &x, &y, &z)
		x--
		y--
		queries[i] = [3]int{x, y, z}
	}

	uf := NewUnionFindArray(n)
	reset := func() {
		uf.Clear()
	}
	mutate := func(mutationId int) {
		u, v := edges[mutationId][0], edges[mutationId][1]
		uf.Union(u, v)
	}
	// 当前点数量是否大于等于 z
	predicate := func(queryId int) bool {
		u, v, z := queries[queryId][0], queries[queryId][1], queries[queryId][2]
		if uf.Find(u) == uf.Find(v) {
			size := uf.GetSize(u)
			return size >= z
		} else {
			size1, size2 := uf.GetSize(u), uf.GetSize(v)
			return size1+size2 >= z
		}
	}

	res := ParallelBinarySearch(m, q, reset, mutate, predicate)
	for i := range res {
		res[i]++
	}
	return res
}

// https://www.luogu.com.cn/problem/CF1706E
// Qpwoeirut and Vertices
// 给出 n 个点， m 条边的不带权连通无向图， q 次询问至少要加完编号前多少的边，
// 才能使得 [start,end) 中的所有点两两连通。
//
// !把[start,end]区间内联通条件拆成(start,start+1),(start+1,start+2),...,(end-2,end-1)
// 分别联通至少需要有几条边，答案取max即可.
func CF1706E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)

	solve := func() {
		var n, m, q int
		fmt.Fscan(in, &n, &m, &q)
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u, v = u-1, v-1
			edges[i] = [2]int{u, v}
		}

		uf := NewUnionFindArray(n)
		reset := func() {
			uf.Clear()
		}
		mutate := func(mutationId int) {
			u, v := edges[mutationId][0], edges[mutationId][1]
			uf.Union(u, v)
		}
		predicate := func(queryId int) bool {
			return uf.Find(queryId) == uf.Find(queryId+1)
		}

		res := ParallelBinarySearch(m, n-1, reset, mutate, predicate)

		seg := NewSegmentTreeFrom(res)
		for i := 0; i < q; i++ {
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			end--
			if start == end {
				fmt.Fprint(out, 0, " ")
			} else {
				fmt.Fprint(out, seg.Query(start, end)+1, " ")
			}
		}
		fmt.Fprintln(out)
	}

	for t := 0; t < T; t++ {
		solve()
	}
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

func DiscretizeCompressed(nums []int, offset int) (getRank func(value int) int, getValue func(rank int) int, count int) {
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
	getValue = func(r int) int { return rank[r-offset] }
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

func (ufa *UnionFindArray) GetSize(key int) int {
	return -ufa.data[ufa.Find(key)]
}

const INF int = 1e18

// PointSetRangeMax

type E = int

func (*SegmentTree) e() E        { return -INF }
func (*SegmentTree) op(a, b E) E { return max(a, b) }
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

type SegmentTree struct {
	n, size int
	seg     []E
}

func NewSegmentTree(n int, f func(int) E) *SegmentTree {
	res := &SegmentTree{}
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func NewSegmentTreeFrom(leaves []E) *SegmentTree {
	res := &SegmentTree{}
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *SegmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTree) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTree) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}
func (st *SegmentTree) QueryAll() E { return st.seg[1] }
func (st *SegmentTree) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree) MaxRight(left int, predicate func(E) bool) int {
	if left == st.n {
		return st.n
	}
	left += st.size
	res := st.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(st.op(res, st.seg[left])) {
			for left < st.size {
				left <<= 1
				if tmp := st.op(res, st.seg[left]); predicate(tmp) {
					res = tmp
					left++
				}
			}
			return left - st.size
		}
		res = st.op(res, st.seg[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return st.n
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree) MinLeft(right int, predicate func(E) bool) int {
	if right == 0 {
		return 0
	}
	right += st.size
	res := st.e()
	for {
		right--
		for right > 1 && right&1 == 1 {
			right >>= 1
		}
		if !predicate(st.op(st.seg[right], res)) {
			for right < st.size {
				right = right<<1 | 1
				if tmp := st.op(st.seg[right], res); predicate(tmp) {
					res = tmp
					right--
				}
			}
			return right + 1 - st.size
		}
		res = st.op(st.seg[right], res)
		if right&-right == right {
			break
		}
	}
	return 0
}
