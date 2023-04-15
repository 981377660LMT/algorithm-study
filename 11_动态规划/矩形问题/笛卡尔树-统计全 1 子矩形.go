package main

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
//  baseLine: 是否只统计高度为1的矩形.
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

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
