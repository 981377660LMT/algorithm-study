package main

func main() {

}

// 统计数组的所有子序列的 GCD 的不同个数，复杂度 O(Clog^2C)
// LC1819 https://leetcode-cn.com/problems/number-of-different-subsequences-gcds/
// 我的题解 https://leetcode.cn/problems/number-of-different-subsequences-gcds/solution/ji-bai-100mei-ju-gcdxun-huan-you-hua-pyt-get7/
func countDifferentSubsequenceGCDs(a []int, gcd func(int, int) int) (ans int) {
	const mx int = 4e5 //
	has := [mx + 1]bool{}
	for _, v := range a {
		has[v] = true
	}
	for i := 1; i <= mx; i++ {
		g := 0
		for j := i; j <= mx && g != i; j += i { // 枚举 i 的倍数 j
			if has[j] { // 如果 j 在 nums 中
				g = gcd(g, j) // 更新最大公约数
			}
		}
		if g == i { // 找到一个答案
			ans++
		}
	}
	return
}
