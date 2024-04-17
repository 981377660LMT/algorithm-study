// 按位异或1最多的二元组
// 你看到“异或”的时候，应该意识到它不具有包容性，因此想继续使用SOSDp只会无功而返
// 但是可以bfs求解

package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"time"
)

func main() {
	bruteForce := func(nums []int) (maxOnesCount int, index1, index2 int) {
		for i, v1 := range nums {
			for j, v2 := range nums {
				if i == j {
					continue
				}
				cand := bits.OnesCount32(uint32(v1 ^ v2))
				if cand > maxOnesCount {
					maxOnesCount = cand
					index1, index2 = i, j
				}
			}
		}
		return
	}

	for i := 0; i < 10; i++ {
		nums := make([]int, 1e4)
		for j := range nums {
			nums[j] = rand.Intn(1e5)
		}
		res1, i1, i2 := bruteForce(nums)
		res2, i3, i4 := BitwiseXorPairWithMaxOnesCount(nums)
		if res1 != res2 || bits.OnesCount32(uint32(nums[i1]^nums[i2])) != res1 || bits.OnesCount32(uint32(nums[i3]^nums[i4])) != res1 {
			println("failed")
		}
	}

	nums := make([]int, 1e6)
	for j := range nums {
		nums[j] = rand.Intn(1e6)
	}
	time1 := time.Now()
	_, _, _ = BitwiseXorPairWithMaxOnesCount(nums)
	fmt.Println(time.Since(time1).Seconds())

	fmt.Println("pass")
}

// !要求找到两个不同的下标i≠j，使得ai^aj包含的1最多。
// nums[i]<=1e6, n<=1e6
func BitwiseXorPairWithMaxOnesCount(nums []int) (maxOnesCount int, index1, index2 int) {
	log := max(bits.Len(uint(maxs(nums))), 1)
	dist := make([]int32, 1<<log) // 凑出值为v的数缺少的bit位数
	for i := range dist {
		dist[i] = 1 << log
	}
	queue := []int32{}
	for _, v := range nums {
		if dist[v] > 0 {
			dist[v] = 0
			queue = append(queue, int32(v))
		}
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for i := 0; i < log; i++ {
			if next := cur ^ (1 << i); dist[next] > dist[cur]+1 {
				dist[next] = dist[cur] + 1
				queue = append(queue, next)
			}
		}
	}

	mask := 1<<log - 1
	best1 := 0
	for i, v := range nums {
		cand := log - int(dist[v^mask])
		if cand > maxOnesCount {
			maxOnesCount = cand
			best1 = v
			index1 = i
		}
	}

	for i, v := range nums {
		if i == index1 {
			continue
		}
		if bits.OnesCount32(uint32(v^best1)) == maxOnesCount {
			index2 = i
			break
		}
	}
	if index1 > index2 {
		index1, index2 = index2, index1
	}
	return
}

func maxs(nums []int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
