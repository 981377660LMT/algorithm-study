// 所有子树 mex (mex从1开始)

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
	LuoguP3605()
}

func smallestMissingValueSubtree(parents []int, nums []int) []int {
	n := len(parents)
	tree := make([][]int, n)
	for i := 1; i < n; i++ {
		tree[parents[i]] = append(tree[parents[i]], i)
	}

	res := make([]int, n)
	mex, counter := 1, [1e5 + 10]int{}
	update := func(root int) {
		counter[nums[root]]++
		for counter[mex] > 0 {
			mex++
		}
	}
	query := func(root int) {
		res[root] = mex
	}
	clear := func(root int) {
		counter[nums[root]]--
		if counter[mex] == 0 && nums[root] < mex {
			mex = nums[root]
		}
	}
	reset := func() {}

	dfu := NewDSUonTree(tree, 0)
	dfu.Run(update, query, clear, reset)
	return res
}

// https://www.luogu.com.cn/problem/P3605
// P3605 [USACO17JAN] Promotion Counting P
// 对每个结点，求子树内有多少个结点的权值比它大
func LuoguP3605() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}

	tree := make([][]int, n)
	for i := 1; i < n; i++ {
		var p int
		fmt.Fscan(in, &p)
		p--
		tree[p] = append(tree[p], i)
	}

	getRank, size := DiscretizeSparse(values, 0)
	newValues := make([]int, n)
	for i := 0; i < n; i++ {
		newValues[i] = getRank(values[i])
	}

	res := make([]int, n)
	bit := NewBitArray(size)
	update := func(root int) {
		bit.Add(newValues[root], 1)
	}
	query := func(root int) {
		res[root] = bit.QueryRange(newValues[root]+1, size)
	}
	clear := func(root int) {
		bit.Add(newValues[root], -1)
	}
	dsu := NewDSUonTree(tree, 0)
	dsu.Run(update, query, clear, nil)

	for i := 0; i < n; i++ {
		fmt.Fprintln(out, res[i])
	}
}

type DSUonTree struct {
	g                        [][]int
	n                        int
	subSize, euler, down, up []int
	idx                      int
	root                     int
}

func NewDSUonTree(tree [][]int, root int) *DSUonTree {
	res := &DSUonTree{
		g:       tree,
		n:       len(tree),
		subSize: make([]int, len(tree)),
		euler:   make([]int, len(tree)),
		down:    make([]int, len(tree)),
		up:      make([]int, len(tree)),
		root:    root,
	}

	res.dfs1(root, -1)
	res.dfs2(root, -1)
	return res
}

// update:添加root处的贡献
// query:查询root的子树的贡献并更新答案
// clear:退出轻儿子时清除root处的贡献
// reset:退出轻儿子时重置所有值(如果需要的话)
func (d *DSUonTree) Run(
	update func(root int),
	query func(root int),
	clear func(root int),
	reset func(),
) {
	var dsu func(cur, par int, keep bool)
	dsu = func(cur, par int, keep bool) {
		nexts := d.g[cur]
		for i := 1; i < len(nexts); i++ {
			if to := nexts[i]; to != par {
				dsu(to, cur, false)
			}
		}

		if d.subSize[cur] != 1 {
			dsu(nexts[0], cur, true)
		}

		if d.subSize[cur] != 1 {
			for i := d.up[nexts[0]]; i < d.up[cur]; i++ {
				update(d.euler[i])
			}
		}

		update(cur)
		query(cur)
		if !keep {
			for i := d.down[cur]; i < d.up[cur]; i++ {
				clear(d.euler[i])
			}
			if reset != nil {
				reset()
			}
		}
	}

	dsu(d.root, -1, false)
}

// 每个结点的欧拉序起点.
func (d *DSUonTree) Id(root int) int {
	return d.down[root]
}

func (d *DSUonTree) dfs1(cur, par int) int {
	d.subSize[cur] = 1
	nexts := d.g[cur]
	if len(nexts) >= 2 && nexts[0] == par {
		nexts[0], nexts[1] = nexts[1], nexts[0]
	}
	for i, next := range nexts {
		if next == par {
			continue
		}
		d.subSize[cur] += d.dfs1(next, cur)
		if d.subSize[next] > d.subSize[nexts[0]] {
			nexts[0], nexts[i] = nexts[i], nexts[0]
		}
	}
	return d.subSize[cur]
}

func (d *DSUonTree) dfs2(cur, par int) {
	d.euler[d.idx] = cur
	d.down[cur] = d.idx
	d.idx++
	for _, next := range d.g[cur] {
		if next == par {
			continue
		}
		d.dfs2(next, cur)
	}
	d.up[cur] = d.idx
}

type BitArray struct {
	n    int
	log  int
	data []int
}

func NewBitArray(n int) *BitArray {
	return &BitArray{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

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

func (b *BitArray) Add(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

// [0, r)
func (b *BitArray) Query(r int) int {
	res := 0
	for ; r > 0; r &= r - 1 {
		res += b.data[r]
	}
	return res
}

func (b *BitArray) QueryRange(l, r int) int {
	return b.Query(r) - b.Query(l)
}

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

func (b *BitArray) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
}

func DiscretizeSparse(nums []int, offset int) (getRank func(int) int, count int) {
	set := make(map[int]struct{})
	for _, v := range nums {
		set[v] = struct{}{}
	}
	count = len(set)
	allNums := make([]int, 0, count)
	for k := range set {
		allNums = append(allNums, k)
	}
	sort.Ints(allNums)
	getRank = func(x int) int { return sort.SearchInts(allNums, x) + offset }
	return
}
