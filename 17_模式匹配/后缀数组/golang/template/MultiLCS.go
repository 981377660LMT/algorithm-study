// k个字符串的最长公共子串
// https://www.acwing.com/problem/content/2813/
// 二分找出最低的height，使某个高于height的区间内包含所有的path

package main

import (
	"bufio"
	"fmt"
	"index/suffixarray"
	"os"
	"reflect"
	"sort"
	"unsafe"
)

func main() {
	p5546()
}

// https://www.luogu.com.cn/problem/P2463
func p5546() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	words := make([][]int32, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		words[i] = make([]int32, len(s))
		for j, b := range s {
			words[i][j] = int32(b - 'a' + 1)
		}
	}
	fmt.Fprintln(out, MultiLCS(words))
}

// 1923. 最长公共子路径
// https://leetcode.cn/problems/longest-common-subpath/solution/hou-zhui-shu-zu-er-fen-da-an-by-endlessc-ocar/
func longestCommonSubpath(n int, paths [][]int) int {
	path32 := make([][]int32, len(paths))
	for i, p := range paths {
		for _, v := range p {
			path32[i] = append(path32[i], int32(v))
		}
	}
	return MultiLCS(path32)
}

const INF int32 = 1e9 + 10

func MultiLCS(ords [][]int32) (res int) {
	sb := []int32{}
	cand := INF // 二分右边界
	for _, p := range ords {
		cand = min32(cand, int32(len(p)))
		sb = append(sb, INF) // dummy
		sb = append(sb, p...)
	}
	n, k := int32(len(sb)), int32(len(ords))

	// 标记每个元素属于哪个数组, dummy 为 -1
	belong := make([]int32, n)
	id := int32(-1)
	for i, v := range sb {
		if v == INF {
			id++
			belong[i] = -1
		} else {
			belong[i] = id
		}
	}

	sa, _, height := SuffixArray32(int32(len(sb)), func(i int32) int32 { return int32(sb[i]) })

	// 二分求答案，找到连续的 height[i] >= mid 的区间，使得区间内包含所有的 path
	return sort.Search(int(cand), func(mid int) bool {
		mid32 := int32(mid)
		mid32++ // bisect_right
		visited := make([]int32, k)
		for i := int32(0); i < n; i++ {
			if height[i] < mid32 || belong[sa[i]] == -1 {
				continue
			}

			continuousCount := int32(0)
			for start := i; i < n && height[i] >= mid32; i++ {
				// 检查 sa[i] 和 sa[i-1]
				if j := belong[sa[i]]; j != -1 && visited[j] != start {
					visited[j] = start
					continuousCount++
				}
				if j := belong[sa[i-1]]; j != -1 && visited[j] != start {
					visited[j] = start
					continuousCount++
				}
			}

			if continuousCount == k { // 凑齐了来自所有数组的元素
				return false
			}
		}

		return true
	})
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
