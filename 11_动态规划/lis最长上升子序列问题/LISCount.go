package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// 求LIS方案数
// LC673 https://leetcode-cn.com/problems/number-of-longest-increasing-subsequence/
// 673. 最长递增子序列的个数
func CountLIS(nums []int, strict bool) int {
	lis := [][]int{}         // 长度为i的lis的结尾的值
	countPreSum := [][]int{} // sum(lis[i][:j])

	for _, v := range nums {
		target := v
		if !strict {
			target++
		}
		pos := sort.Search(len(lis), func(i int) bool { return lis[i][len(lis[i])-1] >= target })
		count := 1
		if pos > 0 {
			i := sort.Search(len(lis[pos-1]), func(i int) bool { return lis[pos-1][i] <= target })
			count = countPreSum[pos-1][len(countPreSum[pos-1])-1] - countPreSum[pos-1][i]
		}
		if pos == len(lis) {
			lis = append(lis, []int{v})
			countPreSum = append(countPreSum, []int{0, count})
		} else {
			lis[pos] = append(lis[pos], v)
			countPreSum[pos] = append(countPreSum[pos], countPreSum[pos][len(countPreSum[pos])-1]+count)
		}
	}
	count := countPreSum[len(countPreSum)-1]
	return count[len(count)-1]
}

const MOD int = 1e9 + 7

// 求长为length的LIS个数.
// 赤壁之战 https://www.acwing.com/problem/content/299/
// 定义 dp[i][j] 表示 a[:j+1] 的长度为 i 且以 a[j] 结尾的 LIS
// 则有 dp[i][j] = ∑dp[i-1][k]  (k<j && a[k]<a[j])
// 注意到当 j 增加 1 时，只多了 k=j 这一个新决策，这样可以用树状数组来维护
// 复杂度 O(n*length*logn)
func CountLISWithLength(nums []int, length int, strict bool) int {
	nums = append(nums[:0:0], nums...)
	_nums := append(nums[:0:0], nums...)
	sort.Ints(_nums)
	for i, v := range nums {
		nums[i] = sort.SearchInts(_nums, v) + 2 // 离散化成从 2 开始的序列
	}

	n := len(nums)
	var tree []int
	add := func(i, val int) {
		for ; i < n+2; i += i & -i {
			tree[i] = (tree[i] + val) % MOD
		}
	}
	sum := func(i int) (res int) {
		for ; i > 0; i &= i - 1 {
			res = (res + tree[i]) % MOD
		}
		return
	}

	dp := make([][]int, length+1)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	for i := 1; i <= length; i++ {
		tree = make([]int, n+2)
		tmp1, tmp2 := dp[i-1], dp[i]
		if i == 1 {
			add(1, 1)
		}
		if strict {
			for j, v := range nums {
				tmp2[j] = sum(v - 1)
				add(v, tmp1[j])
			}
		} else {
			for j, v := range nums {
				tmp2[j] = sum(v)
				add(v, tmp1[j])
			}
		}
	}

	res := 0
	for _, v := range dp[length] {
		res = (res + v) % MOD
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)

	for i := 1; i <= T; i++ {
		var n, m int
		fmt.Fscan(in, &n, &m)
		nums := make([]int, n)
		for j := range nums {
			fmt.Fscan(in, &nums[j])
		}
		fmt.Fprintf(out, "Case #%d: %d\n", i, CountLISWithLength(nums, m, true))
	}
}
