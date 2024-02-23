package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	CF1691D()
}

// 1504. 统计全 1 子矩形 O(ROW*COL)
// https://leetcode.cn/problems/count-submatrices-with-all-ones/
func numSubmat(mat [][]int) int {
	curMat := make([]int, len(mat[0]))
	res := 0
	for i := 0; i < len(mat); i++ {
		for j := 0; j < len(mat[0]); j++ {
			if mat[i][j] == 0 {
				curMat[j] = 0
			} else {
				curMat[j]++
			}
		}
		C := NewCartesianTree(curMat, true)
		res += C.CountSubrectangle(true)
	}
	return res
}

// Max GEQ Sum
// https://www.luogu.com.cn/problem/CF1691D
// 给定一个数组，判断每个子数组最大值是否不小于其和.
//
// 1. 单调栈找到每个元素nums[i]作为最大值时的左右边界区间.
// !2. 如果这个区间内包含nums[i]的最大子数组和>最大值, 则返回NO.
// 3. 否则返回YES.
func CF1691D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(nums []int) bool {
		n := len(nums)
		preSum := make([]int, n+1)
		for i := 0; i < n; i++ {
			preSum[i+1] = preSum[i] + nums[i]
		}

		Range := NewCartesianTree(nums, false).Range // 每个元素作为最大值时的左右边界.
		seg := NewSeg(n+1, func(i int) E { return [2]int{preSum[i], preSum[i]} })

		for i := 0; i < n; i++ {
			start, end := Range[i][0], Range[i][1]
			min_ := seg.Query(start, i+1)[0]
			max_ := seg.Query(i+1, end+1)[1]
			maxSum := max_ - min_ // 包含nums[i]的最大子数组和
			if maxSum > nums[i] {
				return false
			}
		}
		return true
	}

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &nums[i])
		}
		if solve(nums) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

type CartesianTree struct {
	// ![left, right) 每个元素作为最大/最小值时的左右边界.
	//  左侧为严格扩展, 右侧为非严格扩展.
	//  例如: [2, 1, 1, 5] => [[0 1] [0 4] [2 4] [3 4]]
	Range                         [][2]int
	Root                          int
	n                             int
	nums                          []int
	leftChild, rigthChild, parent []int
	isMin                         bool
}

func NewCartesianTree(nums []int, isMin bool) *CartesianTree {
	res := &CartesianTree{}
	n := len(nums)
	Range := make([][2]int, n)
	lch := make([]int, n)
	rch := make([]int, n)
	par := make([]int, n)

	for i := 0; i < n; i++ {
		Range[i] = [2]int{-1, -1}
		lch[i] = -1
		rch[i] = -1
		par[i] = -1
	}

	res.n = n
	res.nums = nums
	res.Range = Range
	res.leftChild = lch
	res.rigthChild = rch
	res.parent = par
	res.isMin = isMin

	if n == 1 {
		res.Range[0] = [2]int{0, 1}
		return res
	}

	compare := func(i, j int) bool {
		if isMin {
			return nums[i] < nums[j] || (nums[i] == nums[j] && i < j)
		}
		return nums[i] > nums[j] || (nums[i] == nums[j] && i < j)
	}

	stack := make([]int, 0)
	for i := 0; i < n; i++ {
		for len(stack) > 0 && compare(i, stack[len(stack)-1]) {
			res.leftChild[i] = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		res.Range[i][0] = 0
		if len(stack) > 0 {
			res.Range[i][0] = stack[len(stack)-1] + 1
		}
		stack = append(stack, i)
	}

	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && compare(i, stack[len(stack)-1]) {
			res.rigthChild[i] = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		res.Range[i][1] = n
		if len(stack) > 0 {
			res.Range[i][1] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}

	for i := 0; i < n; i++ {
		if res.leftChild[i] != -1 {
			res.parent[res.leftChild[i]] = i
		}
		if res.rigthChild[i] != -1 {
			res.parent[res.rigthChild[i]] = i
		}
	}
	for i := 0; i < n; i++ {
		if res.parent[i] == -1 {
			res.Root = i
		}
	}

	return res
}

// 以i处元素为最低点的矩形范围.
func (c *CartesianTree) MaxRectangleAt(i int) (left, right, height int) {
	left, right = c.Range[i][0], c.Range[i][1]
	height = c.nums[i]
	return
}

// 直方图中的最大矩形面积.
func (c *CartesianTree) MaxRectangleArea() int {
	if !c.isMin {
		panic("need min")
	}
	res := 0
	for i := 0; i < c.n; i++ {
		left, right, height := c.MaxRectangleAt(i)
		res = max(res, (right-left)*height)
	}
	return res
}

// 直方图中的矩形数量.
//
//	baseLine: 是否只统计高度为1的矩形.
func (c *CartesianTree) CountSubrectangle(onlyBaseLine bool) int {
	if !c.isMin {
		panic("need min")
	}
	res := 0
	for i := 0; i < c.n; i++ {
		left, right, height := c.MaxRectangleAt(i)
		x := height
		if !onlyBaseLine {
			x = height * (height + 1) / 2
		}
		res += x * (i - left + 1) * (right - i)
	}
	return res
}

// 还原笛卡尔树,返回树的有向邻接表和根节点.
func (c *CartesianTree) GetTree() (tree [][]int, root int) {
	tree = make([][]int, c.n)
	for i := 0; i < c.n; i++ {
		p := c.parent[i]
		if p == -1 {
			root = i
			continue
		}
		tree[p] = append(tree[p], i)
	}
	return
}

const INF int = 1e18

// RangeMinMax

type E = [2]int

func (*Seg) e() E        { return [2]int{INF, -INF} }
func (*Seg) op(a, b E) E { return [2]int{min(a[0], b[0]), max(a[1], b[1])} }
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

type Seg struct {
	n, size int
	seg     []E
}

func NewSeg(n int, f func(int) E) *Seg {
	res := &Seg{}
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

// [start, end)
func (st *Seg) Query(start, end int) E {
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
func (st *Seg) QueryAll() E { return st.seg[1] }
