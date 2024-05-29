// 超级洗衣机环上版本
// https://leetcode.cn/problems/super-washing-machines/

package main

import (
	"fmt"
	"sort"
)

func main() {
	// machines = [1,0,5]
	fmt.Println(CandyAssignProblemSimple([]int{6, 6, 3}, []int{3, 3, 9})) // [2,0,1]
}

// 环上邻位移动，使得数组从cur变成target，返回每个位置的移动次数.
// res[i] 表示有多少个糖果从i移动到i+1，可能为负数.
func CandyAssignProblemSimple(cur []int, target []int) []int {
	n := int32(len(cur))
	if n == 0 {
		return nil
	}
	res := make([]int, n)
	for i := int32(0); i < n-1; i++ {
		prev := i - 1
		if prev < 0 {
			prev += n
		}
		res[i] = res[prev] + cur[i] - target[i]
	}
	kthB := append(res[:0:0], res...)
	for i := int32(0); i < n; i++ {
		kthB[i] = -kthB[i]
	}
	half := (n + 1) / 2
	sort.Ints(kthB)
	x := kthB[half-1]
	for i := int32(0); i < n; i++ {
		res[i] += x
	}
	return res
}
