// 后缀数组矩形区域构成的树，严格来说不能算suffixTree
// https://maspypy.github.io/library/string/suffix_tree.hpp
//
//!   - s = abbbab
//
//!   - suffix array
//     ab
//     abbbab
//     b
//     bab
//     bbab
//     bbbab
//
//!   - suffix tree

//	  	 ab(1) --- bbab(2)
//	  	/
//	  ""(0)
//	  	\
//	  	 b(3) --- ab(4)
//	  	  \
//	  	   b(5) --- ab(6)
//	  	    \
//	  	     ---- bab(7)
//
//!   - [行1，行2，列1，列2] 区域范围
//     0 : [0 6 0 0]
//     1 : [0 2 0 2]
//     2 : [1 2 2 6]
//     3 : [2 6 0 1]
//     4 : [3 4 1 3]
//     5 : [4 6 1 2]
//     6 : [4 5 2 4]
//     7 : [5 6 2 5]
//
//!   - 表格展示
//    _________________________
//    |_1_|_1_|___|___|___|___|
//    |_1_|_1_|_2_|_2_|_2_|_2_|
//    |_3_|___|___|___|___|___|
//    |_3_|_4_|_4_|___|___|___|
//    |_3_|_5_|_6_|_6_|___|___|
//    |_3_|_5_|_7_|_7_|_7_|___|
//

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
	// demo()

	// cf123d()

	test()
}

func demo() {
	s := "abbbab"
	sa, _, lcp := SuffixArray32(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
	tree, ranges := SuffixTreeFrom(sa, lcp)

	// [[1 3] [2] [] [4 5] [] [6 7] [] []]
	// [[0 6 0 0] [0 2 0 2] [1 2 2 6] [2 6 0 1] [3 4 1 3] [4 6 1 2] [4 5 2 4] [5 6 2 5]]
	fmt.Println(tree, ranges)
}

func test() {
	CA := NewCartesianTreeSimple32([]int32{1, 2, 1, 4, 2, 1, 3})
	fmt.Println(CA.Root, CA.parent, "r")
	s := "aabbabbaa"
	_, ranges := SuffixTree(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
	checkRange := func(t [4]int32, xl, xr, yl, yr int32) {
		a, b, c, d := t[0], t[1], t[2], t[3]
		if a != xl || b != yl || c != xr || d != yr {
			fmt.Println(a, b, c, d, xl, xr, yl, yr)
			panic("error")
		}
	}
	checkRange(ranges[0], 0, 0, 9, 0)
	checkRange(ranges[1], 0, 0, 5, 1)
	checkRange(ranges[2], 1, 1, 3, 2)
	checkRange(ranges[3], 2, 2, 3, 9)
	checkRange(ranges[4], 3, 1, 5, 4)
	checkRange(ranges[5], 3, 4, 4, 5)
	checkRange(ranges[6], 4, 4, 5, 8)
	checkRange(ranges[7], 5, 0, 9, 1)
	checkRange(ranges[8], 5, 1, 7, 2)
	checkRange(ranges[9], 5, 2, 6, 3)
	checkRange(ranges[10], 6, 2, 7, 6)
	checkRange(ranges[11], 7, 1, 9, 3)
	checkRange(ranges[12], 7, 3, 8, 4)
	checkRange(ranges[13], 8, 3, 9, 7)
}

//	void solve() {
//	  STR(S);
//	  Suffix_Array X(S);
//	  auto [G, dat] = suffix_tree(X);
//	  ll ANS = 0;
//	  for (auto&& [L, R, lo, hi]: dat) {
//	    ll n = R - L;
//	    ANS += (hi - lo) * n * (n + 1) / 2;
//	  }
//	  print(ANS);
//	}
//

// CF123D String
// https://www.luogu.com.cn/problem/CF123D
// !给定一个字符串 s，定义 cnt(a) 为子串 a 在 s 中出现的次数,求∑cnt(a)*(cnt(a)+1)/2
func cf123d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	_, ranges := SuffixTree(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
	res := 0
	for i := 1; i < len(ranges); i++ {
		rowStart, rowEnd, colStart, colEnd := ranges[i][0], ranges[i][1], ranges[i][2], ranges[i][3]
		rowLen, colLen := int(rowEnd-rowStart), int(colEnd-colStart)
		res += colLen * rowLen * (rowLen + 1) / 2
	}
	fmt.Fprintln(out, res)
}

// TODO: https://yukicoder.me/problems/no/2361
// https://maspypy.github.io/library/test/yukicoder/2361.test.cpp
func yukicoder2361() {

}

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
	fmt.Println(ct.Range, height)
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
	fmt.Println(root, "rot")
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
