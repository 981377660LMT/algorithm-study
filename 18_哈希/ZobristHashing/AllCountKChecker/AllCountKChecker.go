// 判断数据结构中每个数出现的次数是否均为k.
// 等价于:
//  1. 数据结构中每个数出现的次数均为k的倍数：异或哈希.
//  2. 数据结构中每个数出现的次数均不超过k：双指针.
//     在右指针扫到 i 的时候，不停将左指针向右移动并减去这个桶的出现次数，
//     直到 nums[i] 的出现次数小于等于 k 为止。此时再统计答案，两个限制都可以满足。

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
)

func main() {
	CF1418G()
}

// https://www.luogu.com.cn/problem/solution/CF1418G
func CF1418G() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	fmt.Fprintln(out, CountSubarrayWithFrequencyEqualToK(nums, 3))
}

// https://leetcode.cn/problems/count-complete-substrings
func countCompleteSubstrings(word string, k int) int {
	n := len(word)
	ords := make([]int, n)
	for i := 0; i < n; i++ {
		ords[i] = int(word[i] - 'a')
	}

	groups := [][]int{}
	ptr := 0
	for ptr < n {
		leader := ords[ptr]
		group := []int{leader}
		ptr++
		for ptr < n && abs(ords[ptr]-ords[ptr-1]) <= 2 {
			group = append(group, ords[ptr])
			ptr++
		}
		groups = append(groups, group)
	}

	res := 0
	for _, group := range groups {
		res += CountSubarrayWithFrequencyEqualToK(group, k)
	}
	return res
}

type Value = int

func RandomHash(min, max uint64) func(value Value) uint64 {
	pool := make(map[Value]uint64)
	f := func(value Value) uint64 {
		if hash, ok := pool[value]; ok {
			return hash
		}
		rand := rand.Uint64()%(max-min+1) + min
		pool[value] = rand
		return rand
	}
	return f
}

// 统计满足`每个元素出现的次数均为k`条件的子数组的个数.
func CountSubarrayWithFrequencyEqualToK(arr []Value, k int) int {
	n := len(arr)
	if n == 0 || k <= 0 || k > n {
		return 0
	}

	R := RandomHash(1, math.MaxUint64/uint64(n))

	pool := make(map[Value]int)
	getId := func(value Value) int {
		if id, ok := pool[value]; ok {
			return id
		}
		id := len(pool)
		pool[value] = id
		return id
	}
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = getId(arr[i])
	}
	counter := make([]int, len(pool))
	random := make([]uint64, n)
	for i := 0; i < n; i++ {
		random[i] = R(nums[i])
	}
	hashPreSum := make([]uint64, n+1)
	for i := 0; i < n; i++ {
		hashPreSum[i+1] = hashPreSum[i]
		hashPreSum[i+1] -= uint64(counter[nums[i]]) * random[i]
		counter[nums[i]] = (counter[nums[i]] + 1) % k
		hashPreSum[i+1] += uint64(counter[nums[i]]) * random[i]
	}

	countPreSum := make(map[uint64]int)
	countPreSum[0] = 1
	counter = make([]int, len(pool))
	res, left := 0, 0
	for right, num := range nums {
		counter[num]++
		for counter[num] > k {
			counter[nums[left]]--
			countPreSum[hashPreSum[left]]--
			left++
		}
		res += countPreSum[hashPreSum[right+1]]
		countPreSum[hashPreSum[right+1]]++
	}

	return res
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
