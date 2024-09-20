package main

import (
	"cmp"
	"fmt"
	"slices"
	"time"
)

func main() {
	test()
}

// 78. 子集
// https://leetcode.cn/problems/subsets/description/
func subsets(nums []int) [][]int {
	res := make([][]int, 0, 1<<len(nums))
	Powerset(len(nums), func(subset []int) bool {
		tmp := make([]int, 0, len(subset))
		for _, index := range subset {
			tmp = append(tmp, nums[index])
		}
		res = append(res, tmp)
		return false
	})
	return res
}

// 131. 分割回文串
// https://leetcode.cn/problems/palindrome-partitioning/description/
func partition(s string) (res [][]string) {
	isPalindrome := func(start, end int) bool {
		for i, j := start, end-1; i < j; i, j = i+1, j-1 {
			if s[i] != s[j] {
				return false
			}
		}
		return true
	}

	Partitions(len(s), func(splits []int) bool {
		ok := true
		for i := 0; i < len(splits)-1; i++ {
			if !isPalindrome(splits[i], splits[i+1]) {
				ok = false
				break
			}
		}
		if ok {
			tmp := make([]string, 0, len(splits)-1)
			for i := 0; i < len(splits)-1; i++ {
				tmp = append(tmp, s[splits[i]:splits[i+1]])
			}
			res = append(res, tmp)
		}
		return false
	})
	return
}

// 46. 全排列
// https://leetcode.cn/problems/permutations/
func permute(nums []int) (res [][]int) {
	DistinctPermutations(nums, len(nums), func(perm []int) bool {
		res = append(res, append(perm[:0:0], perm...))
		return false
	})
	return
}

// 遍历子集.
func Powerset(n int, f func(subset []int) bool) {
	path := make([]int, 0, n)
	var dfs func(int) bool
	dfs = func(index int) bool {
		if index == n {
			return f(path)
		}
		if dfs(index + 1) {
			return true
		}
		path = append(path, index)
		if dfs(index + 1) {
			return true
		}
		path = path[:len(path)-1]
		return false
	}
	dfs(0)
}

// 遍历数组所有的分割方案，按照分割点将数组分割成若干段.
//
//	Partitions(3, func(splits []int) bool {
//	    for i := 0; i < len(splits)-1 ; i++ {
//	        fmt.Println(arr[splits[i]:splits[i+1]])
//	    }
//	    return false
//	})
func Partitions(n int, f func(splits []int) bool) {
	if n == 0 {
		return
	}
	path := make([]int, 0, n)
	path = append(path, 0)
	var dfs func(int) bool
	dfs = func(index int) bool {
		if index == n-1 {
			path = append(path, n)
			stop := f(path)
			path = path[:len(path)-1]
			return stop
		}
		if dfs(index + 1) {
			return true
		}
		path = append(path, index+1)
		if dfs(index + 1) {
			return true
		}
		path = path[:len(path)-1]
		return false
	}
	dfs(0)
}

// 将 n 个元素的集合分成 k 个部分，不允许为空.
func SetPartitions(n, k int, f func(parts [][]int) bool) {
	if k < 1 {
		panic("Can't partition in a negative or zero number of groups")
	}
	if k > n {
		return
	}
	parts := make([][]int, k)
	for i := range parts {
		parts[i] = make([]int, 0, n)
	}
	var dfs func(int, int) bool
	dfs = func(index, count int) bool {
		if index == n {
			if count == k {
				return f(parts)
			}
			return false
		}
		for i := 0; i < count; i++ {
			parts[i] = append(parts[i], index)
			if dfs(index+1, count) {
				return true
			}
			parts[i] = parts[i][:len(parts[i])-1]
		}
		if count < k {
			parts[count] = append(parts[count], index)
			if dfs(index+1, count+1) {
				return true
			}
			parts[count] = parts[count][:len(parts[count])-1]
		}
		return false
	}
	dfs(0, 0)
}

// 将 n 个元素的集合分成任意个部分.
func SetPartitionsAll(n int, f func(parts [][]int) bool) {
	parts := make([][]int, n)
	for i := range parts {
		parts[i] = make([]int, 0, n)
	}
	var dfs func(int, int) bool
	dfs = func(index, count int) bool {
		if index == n {
			return f(parts[:count])
		}
		for i := 0; i < count; i++ {
			parts[i] = append(parts[i], index)
			if dfs(index+1, count) {
				return true
			}
			parts[i] = parts[i][:len(parts[i])-1]
		}
		parts[count] = append(parts[count], index)
		if dfs(index+1, count+1) {
			return true
		}
		parts[count] = parts[count][:len(parts[count])-1]
		return false
	}
	dfs(0, 0)
}

// 遍历无重复排列.
func DistinctPermutations[T cmp.Ordered](arr []T, r int, f func(perm []T) bool) {
	if len(arr) == 0 || r < 1 || r > len(arr) {
		return
	}
	arr = append(arr[:0:0], arr...)
	slices.Sort(arr)

	if r == len(arr) {
		for {
			if f(arr) {
				break
			}
			if !nextPermutation(arr) {
				break
			}
		}
		return
	}

	head, tail := arr[:r:r], arr[r:len(arr):len(arr)]
	for len(tail) < len(arr) {
		tail = append(tail, arr[len(arr)-1]) // buf
	}
	tailLen := len(arr) - r
	for {
		if f(head) {
			return
		}
		pivot := tail[tailLen-1]
		i := r - 1
		found := false
		for ; i >= 0; i-- {
			if head[i] < pivot {
				found = true
				break
			}
			pivot = head[i]
		}
		if !found {
			return
		}

		found = false
		for j := 0; j < tailLen; j++ {
			if tail[j] > head[i] {
				head[i], tail[j] = tail[j], head[i]
				found = true
				break
			}
		}
		if !found {
			for j := r - 1; j >= 0; j-- {
				if head[j] > head[i] {
					head[i], head[j] = head[j], head[i]
					break
				}
			}
		}

		for j, k := tailLen, r-1; k > i; {
			tail[j] = head[k]
			j++
			k--
		}
		for j, k := i+1, 0; j < r; {
			head[j] = tail[k]
			j++
			k++
		}
		for j, k := 0, r-i-1; j < tailLen; {
			tail[j] = tail[k]
			j++
			k++
		}
	}
}

// 原地返回下一个字典序的排列.
// 不包含重复排列.
func nextPermutation[T cmp.Ordered](nums []T) bool {
	i := len(nums) - 1
	for i > 0 && nums[i-1] >= nums[i] {
		i--
	}
	if i == 0 {
		return false
	}
	last := i - 1
	j := len(nums) - 1
	for nums[j] <= nums[last] {
		j--
	}
	nums[last], nums[j] = nums[j], nums[last]
	for i, j := last+1, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
	return true
}

// 遍历无重复组合.
func DistinctCombinations[T cmp.Ordered](arr []T, r int, f func(comb []T) bool) {
	if len(arr) == 0 || r < 1 || r > len(arr) {
		return
	}

	uniqueEverSeen := make([][]int, len(arr))
	for i := range arr {
		var indexes []int
		visited := make(map[T]struct{})
		for j := i; j < len(arr); j++ {
			if _, has := visited[arr[j]]; has {
				continue
			}
			visited[arr[j]] = struct{}{}
			indexes = append(indexes, j)
		}
		uniqueEverSeen[i] = indexes
	}

	path := make([]T, 0, r)
	var dfs func(int, int) bool
	dfs = func(index, count int) bool {
		if count == r {
			return f(path)
		}
		if index == len(arr) {
			return false
		}
		for _, curIndex := range uniqueEverSeen[index] {
			path = append(path, arr[curIndex])
			if dfs(curIndex+1, count+1) {
				return true
			}
			path = path[:len(path)-1]
		}
		return false
	}
	dfs(0, 0)
}

func test() {
	{
		time1 := time.Now()
		Powerset(27, func(subset []int) bool {
			return false
		})
		fmt.Println(time.Since(time1)) // 900ms
	}

	{
		time1 := time.Now()
		Partitions(27, func(splits []int) bool {
			return false
		})
		fmt.Println(time.Since(time1)) // 550ms
	}

	{
		time1 := time.Now()
		SetPartitions(13, 5, func(parts [][]int) bool {
			return false
		})
		fmt.Println(time.Since(time1)) // 65ms
	}

	{
		time1 := time.Now()
		SetPartitionsAll(13, func(parts [][]int) bool {
			return false
		})
		fmt.Println(time.Since(time1)) // 150ms
	}

	{
		time1 := time.Now()
		count := 0
		DistinctPermutations([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 10, func(perm []int) bool {
			count++
			return false
		})
		fmt.Println(count, "hello")
		fmt.Println(time.Since(time1)) // 400ms
	}

	{
		time1 := time.Now()
		DistinctCombinations([]int{0, 0, 1, 0}, 2, func(comb []int) bool {
			fmt.Println(comb)
			return false
		})
		fmt.Println(time.Since(time1))
	}
}
