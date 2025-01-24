package main

import (
	"log"
	"math/rand"
)

func main() {
	for i := 0; i < 1000; i++ {
		left := rand.Intn(10)
		right := rand.Intn(10)
		k := rand.Intn(10)
		res1 := countInRange(left, right, k)
		res2 := bruteForce(left, right, k)
		if res1 != res2 {
			log.Fatalf("left = %v, right = %v, k = %v, res1 = %v, res2 = %v", left, right, k, res1, res2)
		}
	}
	log.Println("pass")
}

// 左侧可选长度为L，右侧可选长度为R，长度不超过k的非空子数组个数.
// 左侧和右侧都包含当前元素.
func countInRange(left, right, k int) int {
	upper := right
	if upper > k {
		upper = k
	}
	if upper <= 0 {
		return 0
	}

	if left > k {
		return (k + k - upper + 1) * upper / 2
	}

	pos := k - left + 1
	if pos > upper {
		return upper * left
	}

	c1 := pos - 1
	res1 := c1 * left
	c2 := upper - pos + 1
	min_ := k - (upper - 1)
	max_ := k - (pos - 1)
	res2 := (min_ + max_) * c2 / 2
	return res1 + res2
}

func bruteForce(left int, right int, k int) int {
	var res int
	for L := 1; L <= left; L++ {
		for R := 1; R <= right; R++ {
			if L+R-1 <= k { // -1 是去除重复计算的当前元素.
				res++
			}
		}
	}
	return res
}
