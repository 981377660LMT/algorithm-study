// 返回nums的各个子集的元素和的排序后的结果

// func
package main

import "fmt"

func main() {
	fmt.Println(subsetSorted([]int{1, 2}))
}

// 返回nums的各个子集的元素和的排序后的结果.
//  O(2^n)
func subsetSorted(nums []int) []int {
	dp := []int{0}
	for _, v := range nums {
		ndp := make([]int, len(dp))
		for i, w := range dp {
			ndp[i] = w + v
		}
		dp = merge(dp, ndp)
	}
	return dp
}

func merge(a, b []int) []int {
	i, n := 0, len(a)
	j, m := 0, len(b)
	res := make([]int, 0, n+m)
	for {
		if i == n {
			return append(res, b[j:]...)
		}
		if j == m {
			return append(res, a[i:]...)
		}
		if a[i] < b[j] { // 改成 > 为降序
			res = append(res, a[i])
			i++
		} else {
			res = append(res, b[j])
			j++
		}
	}
}
