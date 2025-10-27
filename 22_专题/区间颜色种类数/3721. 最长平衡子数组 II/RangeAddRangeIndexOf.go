// RangeAddRangeIndexOf

package main

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

// 3721. 最长平衡子数组 II
// https://leetcode.cn/problems/longest-balanced-subarray-ii/
// 给你一个整数数组 nums。
// 如果子数组中 不同偶数 的数量等于 不同奇数 的数量，则称该 子数组 是 平衡的 。
// 返回 最长 平衡子数组的长度。
// 子数组 是数组中连续且 非空 的一段元素序列。
func longestBalanced(nums []int) int {
	n := len(nums)
	res := 0
	last := make(map[int]int)
	presum := NewRangeAddRangeIndexOf(n+1, func(i int) int { return 0 }, 50)
	for i := 1; i <= n; i++ {
		v := nums[i-1]
		b := v%2*2 - 1
		if j := last[v]; j == 0 {
			presum.Update(i, n+1, b)
		} else {
			presum.Update(j, i, -b)
		}

		last[v] = i
		sum := presum.Get(i)
		if tmp := presum.IndexOf(sum, 0, i-res); tmp != -1 {
			res = max(res, i-tmp)
		}
	}
	return res
}

type RangeAddRangeIndexOf struct {
	nums       []int
	belong     []int
	blockStart []int
	blockEnd   []int
	blockCount int
	blockLazy  []int
	blockPos   []map[int]int
}

func NewRangeAddRangeIndexOf(
	n int, f func(i int) int,
	blockSize int,
) *RangeAddRangeIndexOf {
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = f(i)
	}

	if blockSize <= 0 {
		blockSize = max(1, int(math.Sqrt(float64(n/4))))
	}
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	belong := make([]int, n)
	for i := 0; i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := 0; i < n; i++ {
		belong[i] = i / blockSize
	}

	res := &RangeAddRangeIndexOf{
		nums:       nums,
		belong:     belong,
		blockStart: blockStart,
		blockEnd:   blockEnd,
		blockCount: blockCount,
		blockLazy:  make([]int, blockCount),
		blockPos:   make([]map[int]int, blockCount),
	}

	for i := 0; i < blockCount; i++ {
		res.rebuild(i)
	}
	return res
}

func (riar *RangeAddRangeIndexOf) Update(start, end int, v int) {
	if start < 0 {
		start = 0
	}
	if end > len(riar.nums) {
		end = len(riar.nums)
	}
	if start >= end {
		return
	}
	bid1, bid2 := riar.belong[start], riar.belong[end-1]
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			riar.nums[i] += v
		}
		riar.rebuild(bid1)
		return
	}
	for i := start; i < riar.blockEnd[bid1]; i++ {
		riar.nums[i] += v
	}
	riar.rebuild(bid1)
	for i := bid1 + 1; i < bid2; i++ {
		riar.blockLazy[i] += v
	}
	for i := riar.blockStart[bid2]; i < end; i++ {
		riar.nums[i] += v
	}
	riar.rebuild(bid2)
}

// 返回区间[start,end)中第一个等于target的下标, 若不存在则返回-1.
func (riar *RangeAddRangeIndexOf) IndexOf(target int, start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > len(riar.nums) {
		end = len(riar.nums)
	}
	if start >= end {
		return -1
	}
	bid1, bid2 := riar.belong[start], riar.belong[end-1]
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			if riar.nums[i]+riar.blockLazy[bid1] == target {
				return i
			}
		}
		return -1
	}
	for i := start; i < riar.blockEnd[bid1]; i++ {
		if riar.nums[i]+riar.blockLazy[bid1] == target {
			return i
		}
	}
	for i := bid1 + 1; i < bid2; i++ {
		if j, ok := riar.blockPos[i][target-riar.blockLazy[i]]; ok {
			return j
		}
	}
	for i := riar.blockStart[bid2]; i < end; i++ {
		if riar.nums[i]+riar.blockLazy[bid2] == target {
			return i
		}
	}
	return -1
}

func (riar *RangeAddRangeIndexOf) Get(i int) int {
	if i < 0 || i >= len(riar.nums) {
		panic("index out of range")
	}
	return riar.nums[i] + riar.blockLazy[riar.belong[i]]
}

func (riar *RangeAddRangeIndexOf) String() string {
	res := make([]string, 0, len(riar.nums))
	for i := 0; i < len(riar.nums); i++ {
		res = append(res, fmt.Sprintf("%d", riar.Get(i)))
	}
	return "[" + strings.Join(res, ",") + "]"
}

func (riar *RangeAddRangeIndexOf) rebuild(bid int) {
	pos := make(map[int]int)
	start, end := riar.blockStart[bid], riar.blockEnd[bid]
	for i := end - 1; i >= start; i-- {
		pos[riar.nums[i]] = i
	}
	riar.blockPos[bid] = pos
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

func testRangeAddRangeIndex() {
	fmt.Println("--- Running Tests for RangeAddRangeIndexOf ---")

	// 测试 1: 基本初始化和单块内更新/查询
	fmt.Println("Test 1: Basic Init, Single Block Update/Query")
	{
		// 初始化数组为 [0, 1, 2, ..., 19]
		riar := NewRangeAddRangeIndexOf(20, func(i int) int { return i }, 5)
		// 检查初始值
		check(riar.IndexOf(5, 0, 20), 5, "Test 1.1: Initial IndexOf")

		// 在单个块内更新 [2, 5) -> [2, 3, 4] 的值加 10
		riar.Update(2, 5, 10)
		// 检查更新后的值
		check(riar.IndexOf(12, 0, 20), 2, "Test 1.2: IndexOf after update") // nums[2]+10 = 2+10=12
		check(riar.IndexOf(13, 0, 20), 3, "Test 1.3: IndexOf after update") // nums[3]+10 = 3+10=13
		check(riar.IndexOf(14, 0, 20), 4, "Test 1.4: IndexOf after update") // nums[4]+10 = 4+10=14
		// 检查未更新的值
		check(riar.IndexOf(1, 0, 20), 1, "Test 1.5: Unchanged value")
		check(riar.IndexOf(5, 0, 20), 5, "Test 1.6: Unchanged value")
	}

	// 测试 2: 跨多块更新/查询
	fmt.Println("\nTest 2: Multi-Block Update/Query")
	{
		riar := NewRangeAddRangeIndexOf(20, func(i int) int { return i }, 5)
		// 更新 [3, 18) 的值加 100
		riar.Update(3, 18, 100)

		// 检查第一个部分块
		check(riar.IndexOf(103, 0, 20), 3, "Test 2.1: First partial block") // 3+100
		check(riar.IndexOf(104, 0, 20), 4, "Test 2.2: First partial block") // 4+100

		// 检查中间的完整块
		check(riar.IndexOf(105, 0, 20), 5, "Test 2.3: Full block")  // 5+100
		check(riar.IndexOf(114, 0, 20), 14, "Test 2.4: Full block") // 14+100

		// 检查最后一个部分块
		check(riar.IndexOf(115, 0, 20), 15, "Test 2.5: Last partial block") // 15+100
		check(riar.IndexOf(117, 0, 20), 17, "Test 2.6: Last partial block") // 17+100

		// 检查未更新的边界
		check(riar.IndexOf(2, 0, 20), 2, "Test 2.7: Unchanged boundary")
		check(riar.IndexOf(18, 0, 20), 18, "Test 2.8: Unchanged boundary")
	}

	// 测试 3: 查询不到目标值
	fmt.Println("\nTest 3: IndexOf Not Found")
	{
		riar := NewRangeAddRangeIndexOf(20, func(i int) int { return i }, 5)
		riar.Update(0, 20, 10) // 所有值加 10
		check(riar.IndexOf(5, 0, 20), -1, "Test 3.1: Target not found")
		check(riar.IndexOf(100, 0, 20), -1, "Test 3.2: Target not found")
	}

	// 测试 4: 多个匹配项，返回第一个
	fmt.Println("\nTest 4: Multiple Occurrences")
	{
		riar := NewRangeAddRangeIndexOf(20, func(i int) int { return i % 5 }, 5) // [0,1,2,3,4,0,1,2,3,4,...]
		check(riar.IndexOf(3, 0, 20), 3, "Test 4.1: First occurrence")
		check(riar.IndexOf(3, 4, 20), 8, "Test 4.2: First occurrence in sub-range")
		check(riar.IndexOf(3, 9, 20), 13, "Test 4.3: First occurrence in sub-range")
	}

	// 测试 5: 边缘情况
	fmt.Println("\nTest 5: Edge Cases")
	{
		riar := NewRangeAddRangeIndexOf(20, func(i int) int { return i }, 5)
		// 空范围更新
		riar.Update(5, 5, 1000)
		check(riar.IndexOf(1005, 0, 20), -1, "Test 5.1: Empty range update")

		// 空范围查询
		check(riar.IndexOf(5, 5, 5), -1, "Test 5.2: Empty range query")

		// 负数增量
		riar.Update(0, 10, -5)
		check(riar.IndexOf(0, 0, 20), 5, "Test 5.3: Negative delta") // 5 + (-5) = 0
	}

	fmt.Println("\n--- All Tests for RangeAddRangeIndexOf Passed ---")
}

// 辅助函数，用于检查结果
func check(got, want interface{}, name string) {
	if !reflect.DeepEqual(got, want) {
		panic(fmt.Sprintf("FAIL: %s, got %v, want %v", name, got, want))
	}
	fmt.Printf("PASS: %s\n", name)
}
