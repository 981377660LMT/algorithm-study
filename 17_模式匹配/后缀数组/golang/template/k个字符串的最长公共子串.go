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

const INF int = 1e18

// 1923. 最长公共子路径
// https://leetcode.cn/problems/longest-common-subpath/solution/hou-zhui-shu-zu-er-fen-da-an-by-endlessc-ocar/
func MultiLCS(ords [][]int) (res int) {
	sb := []int{}
	cand := INF // 二分右边界
	for _, p := range ords {
		cand = min(cand, len(p))
		sb = append(sb, INF) // dummy
		sb = append(sb, p...)
	}
	n, k := len(sb), len(ords)

	// 标记每个元素属于哪个数组
	belong := make([]int, n)
	id := -1
	for i, v := range sb {
		if v == INF {
			id++
			belong[i] = k
		} else {
			belong[i] = id
		}
	}

	sa, _, height := suffixArrayNums(sb)

	// 二分求答案
	return sort.Search(cand, func(limit int) bool {
		limit++ // bisect_right
		visited := make([]int, k)
		for i := 1; i < n; i++ {
			if height[i] < limit {
				continue
			}
			count := 0
			for start := i; i < n && height[i] >= limit; i++ {
				// 检查 sa[i] 和 sa[i-1]
				if j := belong[sa[i]]; j < k && visited[j] != start {
					visited[j] = start
					count++
				}
				if j := belong[sa[i-1]]; j < k && visited[j] != start {
					visited[j] = start
					count++
				}
			}
			if count == k { // 凑齐了来自所有数组的元素
				return false
			}
		}

		return true
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func suffixArrayNums(nums []int) (sa []int32, rank, height []int) {
	n := len(nums)
	s := make([]byte, 0, n*4)
	for _, v := range nums {
		s = append(s, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
	_sa := *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New(s)).Elem().FieldByName("sa").Field(0).UnsafeAddr()))
	sa = make([]int32, 0, n)
	for _, v := range _sa {
		if v&3 == 0 {
			sa = append(sa, v>>2)
		}
	}
	rank = make([]int, n)
	for i := range rank {
		rank[sa[i]] = i
	}
	height = make([]int, n)
	h := 0
	for i, rk := range rank {
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := int(sa[rk-1]); i+h < n && j+h < n && nums[i+h] == nums[j+h]; h++ {
			}
		}
		height[rk] = h
	}

	return
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	words := make([][]int, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		words[i] = make([]int, len(s))
		for j, b := range s {
			words[i][j] = int(b - 'a' + 1)
		}
	}
	fmt.Fprintln(out, MultiLCS(words))

}
