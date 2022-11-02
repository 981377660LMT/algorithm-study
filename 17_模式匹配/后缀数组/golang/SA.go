// Go1.13 开始使用 SA-IS 算法

package main

import (
	"fmt"
	"index/suffixarray"
	"reflect"
	"sort"
	"unsafe"
)

// https://github.dev/EndlessCheng/codeforces-go/copypasta/strings.go
func suffixArray(s string) ([]int32, []int, []int) {
	n := len(s)

	// !后缀数组 sa 排第几的是谁
	// sa[i] 表示n个后缀字典序中的第 i 个字符串在 s 中的位置
	// 也就是将s的n个后缀从小到大进行排序之后把排好序的后缀的开头位置顺次放入SA中。
	sa := *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New([]byte(s))).Elem().FieldByName("sa").Field(0).UnsafeAddr()))

	// !后缀名次数组 rank 你排第几
	// 后缀 s[i:] 位于后缀字典序中的第 rank[i] 个
	// 特别地，rank[0] 即 s 在后缀字典序中的排名，rank[n-1] 即 s[n-1:] 在字典序中的排名
	rank := make([]int, n)
	for i := range rank {
		rank[sa[i]] = i
	}

	// !高度数组 height 也就是排名相邻的两个后缀的最长公共前缀。
	// height[0] = 0
	// height[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	height := make([]int, n)
	h := 0
	for i, rk := range rank {
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := int(sa[rk-1]); i+h < n && j+h < n && s[i+h] == s[j+h]; h++ {
			}
		}
		height[rk] = h
	}

	return sa, rank, height
}

// !查询所有索引位置
func indexOfAll(rawString, searchString string) []int {
	sa := suffixarray.New([]byte(rawString))
	indexes := sa.Lookup([]byte(searchString), -1)
	sort.Ints(indexes)
	return indexes
}

func main() {
	fmt.Println(indexOfAll("bananaaaaa", "a")) // [1 3 5 6 7 8 9]

	sa, rank, height := suffixArray("abcdea")
	fmt.Println("sa:", sa)         // "abcdea" => [5 0 1 2 3 4]
	fmt.Println("rank:", rank)     // "abcdea" => [1 2 3 4 5 0]
	fmt.Println("height:", height) // "abcdea" => [0 1 0 0 0 0]

	sa, rank, height = suffixArray("你好啊")
	fmt.Println("sa:", sa)         // "abcdea" => [8 7 2 4 1 5 0 6 3]
	fmt.Println("rank:", rank)     // "abcdea" => [6 4 2 8 3 5 7 1 0]
	fmt.Println("height:", height) // "abcdea" => [0 0 0 0 0 1 0 0 1]

	fmt.Println([]byte("1000"))

}
