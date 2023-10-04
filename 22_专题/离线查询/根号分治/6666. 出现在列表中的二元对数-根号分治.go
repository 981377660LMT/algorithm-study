// 6666. 出现在列表中的二元对数
// https://leetcode.com/discuss/interview-question/3517350/

package main

import (
	"fmt"
	"math"
	"math/rand"
)

func main() {
	// check
	for i := 0; i < 100; i++ {
		nums := make([]int, 100)
		for j := 0; j < 100; j++ {
			nums[j] = rand.Intn(5)
		}
		pairs := make([][]int, 100)
		for j := 0; j < 100; j++ {
			pairs[j] = []int{rand.Intn(5), rand.Intn(5)}
		}

		if CountPairs(nums, pairs) != CountPairsNaive(nums, pairs) {
			fmt.Println("error")
		}
	}

	fmt.Println("ok")
}

// 给定一个数组nums和一个二元对列表pairs，找出数组中有多少个二元对(A[i], A[j])，
// 使得i < j且(A[i] + A[j])出现在二元对列表中。

// 按照`集合大小`根号分治:
// !遍历小集合，更新大集合
func CountPairs(nums []int, pairs [][]int) int {
	uniquePairs := [][]int{}
	visited := map[[2]int]struct{}{}
	for _, pair := range pairs {
		hash := [2]int{pair[0], pair[1]}
		if _, ok := visited[hash]; !ok {
			uniquePairs = append(uniquePairs, pair)
			visited[hash] = struct{}{}
		}
	}

	rightToLefts := map[int][]int{}
	for _, pair := range uniquePairs {
		left, right := pair[0], pair[1]
		rightToLefts[right] = append(rightToLefts[right], left)
	}

	// 如果右端点对应的lefts不超过sqrt则是小集合，否则是大集合
	// 大集合的个数不超过sqrt个
	threshold := int(math.Sqrt(float64(len(uniquePairs))) + 1)
	leftToRights := map[int][]int{}
	for _, pair := range uniquePairs {
		left, right := pair[0], pair[1]
		if len(rightToLefts[right]) > threshold {
			leftToRights[left] = append(leftToRights[left], right)
		}
	}

	res := 0
	leftCounter, rightCounter := map[int]int{}, map[int]int{}
	for _, num := range nums {
		if lefts := rightToLefts[num]; len(lefts) <= threshold {
			for _, left := range lefts {
				res += leftCounter[left]
			}
		} else {
			res += rightCounter[num]
		}

		leftCounter[num]++
		rights := leftToRights[num]
		for _, right := range rights {
			rightCounter[right]++
		}
	}

	return res
}

func CountPairsNaive(nums []int, pairs [][]int) int {
	visited := map[[2]int]struct{}{}
	for _, pair := range pairs {
		hash := [2]int{pair[0], pair[1]}
		visited[hash] = struct{}{}
	}

	res := 0
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			hash := [2]int{nums[i], nums[j]}
			if _, ok := visited[hash]; ok {
				res++
			}
		}
	}
	return res
}
