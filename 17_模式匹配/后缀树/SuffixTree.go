// 后缀数组举行区域构成的树，严格来说不能算suffixTree
// https://maspypy.github.io/library/string/suffix_tree.hpp
//
//   - s = abbbab
//
//   - suffix array
//     ab
//     abbbab
//     b
//     bab
//     bbab
//     bbbab
//
//   - suffix tree

//		 ab --- bbab
//		/
//
//	""
//
//		\
//		 b --- ab
//		  \
//		   b --- ab
//		    \
//		      --- bab
//
// ```
package main

import "fmt"

func main() {
	s := "abbbab"
	ords := make([]int, len(s))
	for i, b := range s {
		ords[i] = int(b)
	}
	sa, _, lcp := UseSA(ords)
	tree, ranges := SuffixTree(sa, lcp)
	fmt.Println(tree, ranges)
}

// 每个节点为后缀数组上的一个长方形区间.
func SuffixTree(sa, height []int) (directedTree [][]int, ranges [][4]int) {
	height = height[1:]
	n := len(sa)
	if n == 1 {
		directedTree = make([][]int, 2)
		directedTree[0] = append(directedTree[0], 1)
		ranges = append(ranges, [4]int{0, 1, 0, 0})
		ranges = append(ranges, [4]int{0, 1, 0, 1})
		return
	}

	var edges [][2]int
	ranges = append(ranges, [4]int{0, n, 0, 0})
	ct := NewCartesianTreeSimple(height, true)

	var dfs func(parent, index int, h int)
	dfs = func(parent, index int, h int) {
		left, right := ct.Range[index][0], ct.Range[index][1]+1
		hh := height[index]
		if h < hh {
			edges = append(edges, [2]int{parent, len(ranges)})
			parent = len(ranges)
			ranges = append(ranges, [4]int{left, right, h, hh})
		}
		if ct.leftChild[index] == -1 {
			if hh < n-sa[index] {
				edges = append(edges, [2]int{parent, len(ranges)})
				ranges = append(ranges, [4]int{index, index + 1, hh, n - sa[index]})
			}
		} else {
			dfs(parent, ct.leftChild[index], hh)
		}
		if ct.rigthChild[index] == -1 {
			if hh < n-sa[index+1] {
				edges = append(edges, [2]int{parent, len(ranges)})
				ranges = append(ranges, [4]int{index + 1, index + 2, hh, n - sa[index+1]})
			}
		} else {
			dfs(parent, ct.rigthChild[index], hh)
		}
	}

	root := ct.Root
	if height[root] > 0 {
		edges = append(edges, [2]int{0, 1})
		ranges = append(ranges, [4]int{0, n, 0, height[root]})
		dfs(1, root, height[root])
	} else {
		dfs(0, root, 0)
	}

	directedTree = make([][]int, len(ranges))
	for _, e := range edges {
		u, v := e[0], e[1]
		directedTree[u] = append(directedTree[u], v)
	}
	return
}

//	 sa : 排第几的后缀是谁.
//	 rank : 每个后缀排第几.
//	 lcp : 排名相邻的两个后缀的最长公共前缀.
//		lcp[0] = 0
//		lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
func UseSA(ords []int) (sa, rank, lcp []int) {
	n := len(ords)
	sa = GetSA(ords)

	rank = make([]int, n)
	for i := range rank {
		rank[sa[i]] = i
	}

	// !高度数组 lcp 也就是排名相邻的两个后缀的最长公共前缀。
	// lcp[0] = 0
	// lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	lcp = make([]int, n)
	h := 0
	for i, rk := range rank {
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := sa[rk-1]; i+h < n && j+h < n && ords[i+h] == ords[j+h]; h++ {
			}
		}
		lcp[rk] = h
	}

	return
}

func GetSA(ords []int) (sa []int) {
	if len(ords) == 0 {
		return []int{}
	}

	mn := mins(ords)
	for i, x := range ords {
		ords[i] = x - mn + 1
	}
	ords = append(ords, 0)
	n := len(ords)
	m := maxs(ords) + 1
	isS := make([]bool, n)
	isLms := make([]bool, n)
	lms := make([]int, 0, n)
	for i := 0; i < n; i++ {
		isS[i] = true
	}
	for i := n - 2; i > -1; i-- {
		if ords[i] == ords[i+1] {
			isS[i] = isS[i+1]
		} else {
			isS[i] = ords[i] < ords[i+1]
		}
	}
	for i := 1; i < n; i++ {
		isLms[i] = !isS[i-1] && isS[i]
	}
	for i := 0; i < n; i++ {
		if isLms[i] {
			lms = append(lms, i)
		}
	}
	bin := make([]int, m)
	for _, x := range ords {
		bin[x]++
	}

	induce := func() []int {
		sa := make([]int, n)
		for i := 0; i < n; i++ {
			sa[i] = -1
		}

		saIdx := make([]int, m)
		copy(saIdx, bin)
		for i := 0; i < m-1; i++ {
			saIdx[i+1] += saIdx[i]
		}
		for j := len(lms) - 1; j > -1; j-- {
			i := lms[j]
			x := ords[i]
			saIdx[x]--
			sa[saIdx[x]] = i
		}

		copy(saIdx, bin)
		s := 0
		for i := 0; i < m; i++ {
			s, saIdx[i] = s+saIdx[i], s
		}
		for j := 0; j < n; j++ {
			i := sa[j] - 1
			if i < 0 || isS[i] {
				continue
			}
			x := ords[i]
			sa[saIdx[x]] = i
			saIdx[x]++
		}

		copy(saIdx, bin)
		for i := 0; i < m-1; i++ {
			saIdx[i+1] += saIdx[i]
		}
		for j := n - 1; j > -1; j-- {
			i := sa[j] - 1
			if i < 0 || !isS[i] {
				continue
			}
			x := ords[i]
			saIdx[x]--
			sa[saIdx[x]] = i
		}

		return sa
	}

	sa = induce()

	lmsIdx := make([]int, 0, len(sa))
	for _, i := range sa {
		if isLms[i] {
			lmsIdx = append(lmsIdx, i)
		}
	}
	l := len(lmsIdx)
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = -1
	}
	ord := 0
	order[n-1] = ord
	for i := 0; i < l-1; i++ {
		j, k := lmsIdx[i], lmsIdx[i+1]
		for d := 0; d < n; d++ {
			jIsLms, kIsLms := isLms[j+d], isLms[k+d]
			if ords[j+d] != ords[k+d] || jIsLms != kIsLms {
				ord++
				break
			}
			if d > 0 && (jIsLms || kIsLms) {
				break
			}
		}
		order[k] = ord
	}
	b := make([]int, 0, l)
	for _, i := range order {
		if i >= 0 {
			b = append(b, i)
		}
	}
	var lmsOrder []int
	if ord == l-1 {
		lmsOrder = make([]int, l)
		for i, ord := range b {
			lmsOrder[ord] = i
		}
	} else {
		lmsOrder = GetSA(b)
	}
	buf := make([]int, len(lms))
	for i, j := range lmsOrder {
		buf[i] = lms[j]
	}
	lms = buf
	return induce()[1:]
}

func mins(a []int) int {
	mn := a[0]
	for _, x := range a {
		if x < mn {
			mn = x
		}
	}
	return mn
}

func maxs(a []int) int {
	mx := a[0]
	for _, x := range a {
		if x > mx {
			mx = x
		}
	}
	return mx
}

type CartesianTreeSimple struct {
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

func NewCartesianTreeSimple(nums []int, isMin bool) *CartesianTreeSimple {
	res := &CartesianTreeSimple{}
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
