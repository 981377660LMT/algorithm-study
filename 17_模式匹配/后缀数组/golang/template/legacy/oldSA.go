// Go1.13 开始使用 SA-IS 算法

package main

import (
	"fmt"
	"index/suffixarray"
	"reflect"
	"unsafe"
)

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

// https://github.dev/EndlessCheng/codeforces-go/copypasta/strings.go
func suffixArray(s string) (sa []int32, rank, height []int) {
	n := len(s)

	// !后缀数组 sa 排第几的是谁
	// sa[i] 表示n个后缀字典序中的第 i 个字符串在 s 中的位置
	// 也就是将s的n个后缀从小到大进行排序之后把排好序的后缀的开头位置顺次放入SA中。
	sa = *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New([]byte(s))).Elem().FieldByName("sa").Field(0).UnsafeAddr()))

	// !后缀名次数组 rank 你排第几
	// 后缀 s[i:] 位于后缀字典序中的第 rank[i] 个
	// 特别地，rank[0] 即 s 在后缀字典序中的排名，rank[n-1] 即 s[n-1:] 在字典序中的排名
	rank = make([]int, n)
	for i := range rank {
		rank[sa[i]] = i
	}

	// !高度数组 height 也就是排名相邻的两个后缀的最长公共前缀。
	// height[0] = 0
	// height[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	height = make([]int, n)
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

	return
}

func main() {

	sa, rank, height := suffixArray("abcdea")
	fmt.Println("sa:", sa)         // "abcdea" => [5 0 1 2 3 4]
	fmt.Println("rank:", rank)     // "abcdea" => [1 2 3 4 5 0]
	fmt.Println("height:", height) // "abcdea" => [0 1 0 0 0 0]

	sa, rank, height = suffixArray("你好啊")
	fmt.Println("sa:", sa)         // "abcdea" => [8 7 2 4 1 5 0 6 3]
	fmt.Println("rank:", rank)     // "abcdea" => [6 4 2 8 3 5 7 1 0]
	fmt.Println("height:", height) // "abcdea" => [0 0 0 0 0 1 0 0 1]

	fmt.Println([]byte("1000"))

	sa, rank, height = suffixArrayNums([]int{1, 0, 0, 0})
	fmt.Println("sa:", sa) // "abcdea" => [8 7 2 4 1 5 0 6 3]
	fmt.Println("rank:", rank)
	fmt.Println("height:", height)
}
