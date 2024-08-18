// 100409. 找出最大的 N 位 K 回文数
// https://leetcode.cn/problems/find-the-largest-palindrome-divisible-by-k/solutions/2884548/tong-yong-zuo-fa-jian-tu-dfsshu-chu-ju-t-m3pu/找出最大的N位回文数，使得该回文数可以被 K 整除。
// 求最大的被k整除的n位回文数。
// 1 <= n <= 1e5
// 1 <= k <= 9
// !dp求方案.
// 由于我们只关心回文数模 k 的值是否为 0，可以把：
// 当前从右到左填到第 i 位。
// 已填入的数字模 k 的值为 j。
// 作为状态 (i,j)。
// 填入数字后，回文模k的值变成了
// !j2 = (j+d*(10^i+10^(n-i-1)))%k
// !注意特判n为奇数且i=mid-1的情况(mid=ceil(n/2))，此时模值变成了 (j+d*10^mid)%k。
// 因为要求最大，所以这里用dfs简化dp.

package main

import "fmt"

func main() {
	// n = 5, k = 6
	fmt.Println(largestPalindrome(5, 6)) // 99999
}

func largestPalindrome(n int, k int) string {
	pow10 := make([]int, n)
	pow10[0] = 1
	for i := 1; i < n; i++ {
		pow10[i] = pow10[i-1] * 10 % k
	}

	mid := (n + 1) / 2
	dp := make([]int, (mid+1)*k)
	pre := make([]int, (mid+1)*k)
	preValue := make([]int, (mid+1)*k)
	for i := range dp {
		dp[i] = -1
		pre[i] = -1
		preValue[i] = -1
	}

	var dfs func(int, int) int
	dfs = func(pos, mod int) int {
		if pos == mid {
			if mod == 0 {
				return 1
			}
			return 0
		}
		hash := pos*k + mod
		if dp[hash] != -1 {
			return dp[hash]
		}
		for d := 9; d >= 0; d-- {
			var nMod int
			if n&1 == 1 && pos == mid-1 {
				nMod = (mod + d*pow10[pos]) % k
			} else {
				nMod = (mod + d*(pow10[pos]+pow10[n-pos-1])) % k
			}
			if dfs(pos+1, nMod) == 1 {
				dp[hash] = 1
				nHash := (pos+1)*k + nMod
				pre[nHash] = hash
				preValue[nHash] = d
				return 1
			}
		}
		dp[hash] = 0
		return 0
	}

	dfs(0, 0)

	res := make([]byte, n)
	cur := mid * k
	ptr := mid - 1
	for pre[cur] != -1 {
		res[ptr] = '0' + byte(preValue[cur])
		res[n-1-ptr] = res[ptr]
		ptr--
		cur = pre[cur]
	}
	return string(res)
}
