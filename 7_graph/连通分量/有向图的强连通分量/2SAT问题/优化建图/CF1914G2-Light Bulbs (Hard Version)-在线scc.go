package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

func main() {
	CF1914G2()
}

// Light Bulbs (Hard Version)
// https://www.luogu.com.cn/problem/CF1914G2
// 有一个长度为 2n的带颜色序列，颜色 1∼n 各出现两次.
// 初始你可以对点集 S 染色。
// 然后，你可以不断进行两个操作。
// 1.如果一个数被染色了，则另一个相同的数也会被染色。
// 2.如果某两个相同的数被染色了，那么它们之间所有数都会被染色。
// 你要求出最小可能的 S 集合大小使得所有数都被染色，并且求出有多少 S 满足其大小到达最小。输出对 998244353 取模。
//
// 如果数 u 被染色，可以使得数 v 也被染色，那么我们就连一条 u 到 v 的有向边表示这种关系。
// !问题变成了，在一个有向图中，你可以标记若干个点，要求每个点都可以由被标记的点达到。
// !最后因为要考虑选哪些点作为起始点，使得图中所有节点都可以被遍历到，因此考虑对这个图求强连通分量缩点成DAG。
// 此时所有入度为 0 的点（缩点后）必然是要选择的点。
// 当把所有入度为 0 的点选择后，根据拓扑排序知道一定能遍历图中的所有点。
// !因此要选择的点数的最小值就是缩点后入度为0的点的数量，只需在所有该节点所表示的强连通分量中任意选择一个节点即可（因为强连通分量内的点可以互相到达）。
// !因此选择的方案数也知道了，就是每个入度为0的点所表示的强连通分量的大小的乘积。
// https://codeforces.com/contest/1914/submission/245626432
func CF1914G2() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	solve := func() {
		var n int
		fmt.Fscan(in, &n)
		nums := make([]int, 2*n)
		for i := 0; i < 2*n; i++ {
			fmt.Fscan(in, &nums[i])
			nums[i]--
		}

		colorPos := make([][]int, n)
		for i, v := range nums {
			colorPos[v] = append(colorPos[v], i)
		}
		set := NewFastSetFrom(2*n, func(i int) bool { return true })     // 维护原图未访问的位置.
		seg := NewSegmentTree(2*n, func(i int) E { return E{-INF, -1} }) // 维护反图未访问的位置.

		setUsed := func(curColor int, rev bool) {
			left, right := colorPos[curColor][0], colorPos[curColor][1]
			if !rev {
				// !原图上标记已访问的位置
				set.Erase(left)
				set.Erase(right)
			} else {
				// !反图上标记已访问的位置
				seg.Set(left, seg.e())
			}
		}
		findUnused := func(curColor int, rev bool) int {
			left, right := colorPos[curColor][0], colorPos[curColor][1]
			if !rev {
				// !原图上寻找未访问过的位置: [left+1,right) 之间的一个未访问的位置
				cand := set.Next(left + 1)
				if cand >= right {
					return -1
				}
				return nums[cand]
			} else {
				// !反图上寻找未访问过的位置：原图中哪种颜色可以到达当前颜色curColor？
				tmp1 := seg.Query(0, left) // 可以到达left
				right1, color1 := tmp1.right, tmp1.color
				if right1 > left {
					return color1
				}
				tmp2 := seg.Query(0, right) // 可以到达right
				right2, color2 := tmp2.right, tmp2.color
				if right2 > right {
					return color2
				}
				return -1
			}
		}
		resetSeg := func() { // 将反图上的所有点标记为未访问.
			for c := 0; c < n; c++ {
				left, right := colorPos[c][0], colorPos[c][1]
				seg.Set(left, E{right: right, color: c})
			}
		}

		resetSeg()
		count, belong := OnlineSCC(n, setUsed, findUnused) // n种颜色
		groups := make([][]int, count)                     // 拓扑序逆序分组
		for i, g := range belong {
			groups[g] = append(groups[g], i)
		}

		resetSeg()
		res := [2]int{0, 1}
		for i := len(groups) - 1; i >= 0; i-- { // 按照拓扑序遍历DAG,寻找入度为0的scc.
			group := groups[i]
			for _, v := range group {
				setUsed(v, true)
			}

			isIndegZero := true
			for _, v := range group {
				x := findUnused(v, true)
				if x != -1 {
					isIndegZero = false
					break
				}
			}

			if isIndegZero {
				res[0]++
				res[1] = res[1] * 2 * len(group) % MOD
			}
		}

		fmt.Fprintln(out, res[0], res[1])
	}

	var T int
	fmt.Fscan(in, &T)
	for i := 0; i < T; i++ {
		solve()
	}
}

const INF int = 1e18

// PointSetRangeMaxIndex

type E = struct{ right, color int }

func (*SegmentTree) e() E { return E{-INF, -1} }
func (*SegmentTree) op(a, b E) E {
	if a.right > b.right {
		return a
	}
	if a.right < b.right {
		return b
	}
	if a.color < b.color {
		return a
	}
	return b
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

type FastSet struct {
	n, lg int
	seg   [][]int
	size  int
}

func NewFastSet(n int) *FastSet {
	res := &FastSet{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func NewFastSetFrom(n int, f func(i int) bool) *FastSet {
	res := NewFastSet(n)
	for i := 0; i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
			res.size++
		}
	}
	for h := 0; h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *FastSet) Has(i int) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet) Insert(i int) bool {
	if fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	fs.size++
	return true
}

func (fs *FastSet) Erase(i int) bool {
	if !fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	fs.size--
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == len(cache) {
			break
		}
		d := cache[i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *FastSet) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *FastSet) Enumerate(start, end int, f func(i int)) {
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *FastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *FastSet) Size() int {
	return fs.size
}

func (*FastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*FastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}

// kosaraju 在线求有向图的强连通分量.常用于2-sat优化建图问题.
//
//	setUsed(v, rev)：将 v 设置为已使用, rev 表示是否是反图
//	findUnused(v, rev)：返回未使用过的点中与 v 相邻的点, rev 表示是否是反图.不存在时返回 -1.
//
// 返回强连通分量的个数和每个点所属的分量编号.
// !注意按照0到count-1的遍历结果为拓扑序逆序.
//
// 步骤:
// https://www.cnblogs.com/RioTian/p/14026585.html
// 1.通过后序遍历的方式遍历整个有向图，并且维护每个点的出栈顺序
// 2.根据出栈顺序从大到小遍历反图
// 3.对点u来说，在遍历反图时所有能够到达的v都和u在一个强连通分量当中
func OnlineSCC(n int, setUsed func(cur int, rev bool), findUnused func(cur int, rev bool) int) (count int, belong []int) {
	belong = make([]int, n)

	stack := make([]int, n)
	ptr := n

	visited := make([]bool, n)
	var dfs1 func(v int)
	dfs1 = func(v int) { // 原图
		visited[v] = true
		setUsed(v, false)
		for {
			to := findUnused(v, false)
			if to == -1 {
				break
			}
			dfs1(to)
		}
		ptr--
		stack[ptr] = v
	}
	for v := 0; v < n; v++ {
		if !visited[v] {
			dfs1(v)
		}
	}

	visited = make([]bool, n)
	var dfs2 func(v int)
	dfs2 = func(v int) { // 反图
		visited[v] = true
		belong[v] = count
		setUsed(v, true)
		for {
			to := findUnused(v, true)
			if to == -1 {
				break
			}
			dfs2(to)
		}
	}
	for _, v := range stack {
		if !visited[v] {
			dfs2(v)
			count++
		}
	}

	return
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
