package main

func main() {

}

const MOD int = 998244353

func SolveInclusionExclusion(nums []int) (res int) {
	n := len(nums)
	for state := 0; state < 1<<n; state++ {
		cur := 0
		count := 0
		for i, v := range nums {
			if state>>i&1 > 0 {
				// 视情况而定，有时候包含元素 i 表示考虑这种情况，有时候表示不考虑这种情况
				_ = v
				count++
			}
		}
		if count&1 > 0 {
			cur = -cur // 某些题目是 == 0
		}

		res = (res + cur) % MOD
	}
	if res < 0 {
		res += MOD
	}
	return
}
