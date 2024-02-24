// TODO: 后缀自动机版本

package main

import (
	"bufio"
	"fmt"
	"index/suffixarray"
	"os"
	"reflect"
	"unsafe"
)

func main() {
	// P3975()

	cf128B()
}

// P3975 [TJOI2015] 弦论
// https://www.luogu.com.cn/problem/P3975
func P3975() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	var b, k int
	fmt.Fscan(in, &b, &k)

	unique := b == 0
	start, end, ok := KthSmallestSubstring(
		int32(len(s)), func(i int32) int32 { return int32(s[i]) }, k,
		unique,
	)
	if !ok {
		fmt.Fprintln(out, -1)
		return
	}
	fmt.Fprintln(out, s[start:end])
}

// https://www.luogu.com.cn/problem/CF128B
func cf128B() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	var k int
	fmt.Fscan(in, &k)

	start, end, ok := KthSmallestSubstring(
		int32(len(s)), func(i int32) int32 { return int32(s[i]) }, k,
		false,
	)
	if !ok {
		fmt.Fprintln(out, "No such line.")
		return
	}
	fmt.Fprintln(out, s[start:end])
}

// 字典序第k小的子串.k>=1.
// unique 表示是否对所有子串去重.
func KthSmallestSubstring(n int32, f func(i int32) int32, k int, unique bool) (start, end int, ok bool) {
	if unique {
		return solveUnique(n, f, k)
	} else {
		return solveWithoutUnique(n, f, k)
	}
}

// 直接遍历height数组, 每个后缀会增加 n-sa[i]-height[i] 个新串，并且是按照字典序加入的.
func solveUnique(n int32, f func(i int32) int32, k int) (start, end int, ok bool) {
	sa, _, height := SuffixArray32(n, f)
	remain := k
	for i := int32(0); i < n; i++ {
		newCount := int(n - sa[i] - height[i])
		if remain > newCount {
			remain -= newCount
		} else {
			start = int(sa[i])
			end = start + int(height[i]) + remain
			ok = true
			return
		}
	}
	return
}

// 遍历后缀树.
func solveWithoutUnique(n int32, f func(i int32) int32, k int) (start, end int, ok bool) {
	if k > int(n)*(int(n)+1)/2 {
		return
	}

	sa, _, height := SuffixArray32(n, f)
	tree, ranges := SuffixTreeFrom(sa, height)
	remain := k
	var dfs func(cur int32) bool
	dfs = func(cur int32) bool {
		freq, length := int(ranges[cur][1]-ranges[cur][0]), int(ranges[cur][3]-ranges[cur][2])
		count := freq * length
		if remain <= count {
			remain--
			div, mod := remain/freq, remain%freq
			// RecoverSubstring
			row := int(ranges[cur][0]) + mod
			colEnd := int(ranges[cur][2]) + div + 1
			start = int(sa[row])
			end = int(sa[row]) + colEnd
			ok = true
			return true
		}
		remain -= count
		for _, next := range tree[cur] {
			if dfs(next) {
				return true
			}
		}
		return false
	}
	dfs(0)
	return
}

// directTree: 后缀树, 从 0 开始编号, 0 结点为虚拟根节点.
// ranges: 每个结点对应后缀数组上的 [行1，行2，列1，列2] 矩形区域.
// !(行2-行1) 表示此startPos出现次数, (列2-列1) 表示结点包含的压缩的字符串长度(个数).
func SuffixTree(n int32, f func(i int32) int32) (directedTree [][]int32, ranges [][4]int32) {
	sa, _, lcp := SuffixArray32(n, f)
	return SuffixTreeFrom(sa, lcp)
}

// 每个节点为后缀数组上的一个矩形区间.
func SuffixTreeFrom(sa, height []int32) (directedTree [][]int32, ranges [][4]int32) {
	height = height[1:]
	n := int32(len(sa))
	if n == 1 {
		directedTree = make([][]int32, 2)
		directedTree[0] = append(directedTree[0], 1)
		ranges = append(ranges, [4]int32{0, 1, 0, 0})
		ranges = append(ranges, [4]int32{0, 1, 0, 1})
		return
	}

	var edges [][2]int32
	ranges = append(ranges, [4]int32{0, n, 0, 0})
	ct := NewCartesianTreeSimple32(height)

	var dfs func(p, idx int32, h int32)
	dfs = func(p, idx int32, h int32) {
		left, right := ct.Range[idx][0], ct.Range[idx][1]+1
		hh := height[idx]
		if h < hh {
			m := int32(len(ranges))
			edges = append(edges, [2]int32{p, m})
			p = m
			ranges = append(ranges, [4]int32{left, right, h, hh})
		}

		if ct.leftChild[idx] == -1 {
			if hh < n-sa[idx] {
				edges = append(edges, [2]int32{p, int32(len(ranges))})
				ranges = append(ranges, [4]int32{idx, idx + 1, hh, n - sa[idx]})
			}
		} else {
			dfs(p, ct.leftChild[idx], hh)
		}

		if ct.rigthChild[idx] == -1 {
			if hh < n-sa[idx+1] {
				edges = append(edges, [2]int32{p, int32(len(ranges))})
				ranges = append(ranges, [4]int32{idx + 1, idx + 2, hh, n - sa[idx+1]})
			}
		} else {
			dfs(p, ct.rigthChild[idx], hh)
		}
	}

	root := ct.Root
	if height[root] > 0 {
		edges = append(edges, [2]int32{0, 1})
		ranges = append(ranges, [4]int32{0, n, 0, height[root]})
		dfs(1, root, height[root])
	} else {
		dfs(0, root, 0)
	}

	directedTree = make([][]int32, len(ranges))
	for _, e := range edges {
		u, v := e[0], e[1]
		directedTree[u] = append(directedTree[u], v)
	}
	return
}

func SuffixArray32(n int32, f func(i int32) int32) (sa, rank, height []int32) {
	s := make([]byte, 0, n*4)
	for i := int32(0); i < n; i++ {
		v := f(i)
		s = append(s, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
	_sa := *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New(s)).Elem().FieldByName("sa").Field(0).UnsafeAddr()))
	sa = make([]int32, 0, n)
	for _, v := range _sa {
		if v&3 == 0 {
			sa = append(sa, v>>2)
		}
	}
	rank = make([]int32, n)
	for i := int32(0); i < n; i++ {
		rank[sa[i]] = i
	}
	height = make([]int32, n)
	h := int32(0)
	for i := int32(0); i < n; i++ {
		rk := rank[i]
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := sa[rk-1]; i+h < n && j+h < n && f(i+h) == f(j+h); h++ {
			}
		}
		height[rk] = h
	}
	return
}

type CartesianTreeSimple32 struct {
	// ![left, right) 每个元素作为最大/最小值时的左右边界.
	//  左侧为严格扩展, 右侧为非严格扩展.
	//  例如: [2, 1, 1, 5] => [[0 1] [0 4] [2 4] [3 4]]
	Range                         [][2]int32
	Root                          int32
	n                             int32
	nums                          []int32
	leftChild, rigthChild, parent []int32
}

// min
func NewCartesianTreeSimple32(nums []int32) *CartesianTreeSimple32 {
	res := &CartesianTreeSimple32{}
	n := int32(len(nums))
	Range := make([][2]int32, n)
	lch := make([]int32, n)
	rch := make([]int32, n)
	par := make([]int32, n)

	for i := int32(0); i < n; i++ {
		Range[i] = [2]int32{-1, -1}
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

	if n == 1 {
		res.Range[0] = [2]int32{0, 1}
		return res
	}

	less := func(i, j int32) bool {
		return (nums[i] < nums[j]) || (nums[i] == nums[j] && i < j)
	}

	stack := make([]int32, 0)
	for i := int32(0); i < n; i++ {
		for len(stack) > 0 && less(i, stack[len(stack)-1]) {
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
		for len(stack) > 0 && less(i, stack[len(stack)-1]) {
			res.rigthChild[i] = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		res.Range[i][1] = n
		if len(stack) > 0 {
			res.Range[i][1] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}

	for i := int32(0); i < n; i++ {
		if res.leftChild[i] != -1 {
			res.parent[res.leftChild[i]] = i
		}
		if res.rigthChild[i] != -1 {
			res.parent[res.rigthChild[i]] = i
		}
	}
	for i := int32(0); i < n; i++ {
		if res.parent[i] == -1 {
			res.Root = i
		}
	}

	return res
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

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
