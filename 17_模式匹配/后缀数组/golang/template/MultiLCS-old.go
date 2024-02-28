package main

import (
	"fmt"
	"index/suffixarray"
	"reflect"
	"unsafe"
)

func main() {
	// p5546()
	// 	5
	// [[0,1,0,1,0,1,0,1,0],[0,1,3,0,1,4,0,1,0]]
	n := 5
	paths := [][]int{{0, 1, 0, 1, 0, 1, 0, 1, 0}, {0, 1, 3, 0, 1, 4, 0, 1, 0}}
	fmt.Println(longestCommonSubpath(n, paths))
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
	row := int32(len(ords))
	sb := []int32{}
	lid, rid := make([]int32, row), make([]int32, row)
	minLen := INF
	for i, p := range ords {
		minLen = min32(minLen, int32(len(p)))
		sb = append(sb, INF) // splitter
		lid[i] = int32(len(sb))
		sb = append(sb, p...)
		rid[i] = int32(len(sb))
	}

	n := int32(len(sb))
	_, rank, height := SuffixArray32(n, func(i int32) int32 { return int32(sb[i]) })
	belong := make([]int32, n)
	for i := range belong {
		belong[i] = -1
	}
	for sid := int32(0); sid < int32(row); sid++ {
		for j := lid[sid]; j < rid[sid]; j++ {
			p := rank[j]
			belong[p] = sid
		}
	}

	// 二分答案，遍历高度数组找到连续的 height[i] >= mid 的区间，使得区间内包含所有的 path
	check := func(mid int32) bool {
		visited := make([]int32, row)
		for i := range visited {
			visited[i] = -1
		}
		for i := int32(0); i < n; i++ {
			if height[i] < mid || belong[i] == -1 {
				continue
			}

			first := i // group leader
			continuousCount := int32(0)
			for ; i < n && height[i] >= mid; i++ {
				// 检查 i
				if j := belong[i]; j != -1 && visited[j] != first {
					visited[j] = first
					continuousCount++
				}
				// 检查 sa[i-1]
				if j := belong[i-1]; j != -1 && visited[j] != first {
					visited[j] = first
					continuousCount++
				}
			}

			if continuousCount == row { // 凑齐了来自所有数组的元素
				return true
			}
		}
		return false
	}

	left, right := int32(1), minLen
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return int(right)

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
