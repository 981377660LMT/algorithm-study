// 给定一个集合，对该集合的所有子集，计算该子集的所有子集之和
//（这个「和」不一定是加法，可以是其它的满足合并性质的统计量）

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
)

func main() {
	// CF165E()
	// CF383E()
	// CF449D()
	// CF1234F()
	// ARC100E()

	// 维护集合最大和次大的写法()
	// demo()

	bruteForce1 := func(nums []int) int {
		res := 0
		for i := 0; i < len(nums); i++ {
			for j := i + 1; j < len(nums); j++ {
				res = max(res, nums[i]|nums[j])
			}
		}
		return res
	}
	bruteForce2 := func(nums []int) (maxOnesCount int, index1, index2 int) {
		maxOnesCount = 0
		for i := 0; i < len(nums); i++ {
			for j := i + 1; j < len(nums); j++ {
				if c := bits.OnesCount32(uint32(nums[i] & nums[j])); c > maxOnesCount {
					maxOnesCount = c
					index1, index2 = i, j
				}
			}
		}
		return
	}
	bruteForce3 := func(nums []int) int {
		res := 0
		for i := 0; i < len(nums); i++ {
			for j := i + 1; j < len(nums); j++ {
				if nums[i]&nums[j] == 0 {
					res++
				}
			}
		}
		return res
	}
	bruteForce4 := func(nums []int) (maxOnesCount int, index1, index2 int) {
		maxOnesCount = 0
		for i := 0; i < len(nums); i++ {
			for j := i + 1; j < len(nums); j++ {
				if c := bits.OnesCount32(uint32(nums[i] | nums[j])); c > maxOnesCount {
					maxOnesCount = c
					index1, index2 = i, j
				}
			}
		}
		return
	}

	T := 5
	for i := 0; i < T; i++ {
		nums := make([]int, 100)
		for j := range nums {
			nums[j] = rand.Intn(1e6)
		}

		actual1 := MaxBitwiseOrPair(nums)
		expected1 := bruteForce1(nums)
		if expected1 != actual1 {
			fmt.Println(expected1, actual1, nums)
			panic("not equal1")
		}

		actual2, index1, index2 := BitwiseAndPairWithMaxOnesCount(nums)
		expected2, index1, index2 := bruteForce2(nums)
		if expected2 != actual2 {
			fmt.Println(expected2, actual2, nums, index1, index2)
			panic("not equal2")
		}

		actual3 := 按位与为0的二元组(nums)
		expected3 := bruteForce3(nums)
		if expected3 != actual3 {
			fmt.Println(expected3, actual3, nums)
			panic("not equal3")
		}

		actual4, index1, index2 := BitwiseOrPairWithMaxOnesCount(nums)
		expected4, index1, index2 := bruteForce4(nums)
		if expected4 != actual4 {
			fmt.Println(expected4, actual4, nums)
			panic("not equal4")
		}
	}

	fmt.Println(BitwiseOrPairWithMaxOnesCount([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))

	fmt.Println("pass")
}

// 按位或最大的二元组
// 要求找到两个不同的下标i≠j，使得ai∣aj最大。
// nums[i]<=1e6, n<=1e6
// Max 2-OR Problem in O(MlogM)
// !O(UlogU)  Boolean orthogonal detection
// https://www.cs.toronto.edu/~deepkush/ovs.pdf
// https://cs.stackexchange.com/questions/82743/for-given-set-of-integers-find-and-count-the-pairs-with-maximum-value-of-bitwis
//
// 1.由于i=j时结果一定不是最优，因此i≠j的约束可以无视.
// !2.先将信息下推到子集，再对每个数，从高位比特开始枚举到最低比特，试填法判断这一位能否填1.
func MaxBitwiseOrPair(nums []int) int {
	log := max(bits.Len(uint(maxs(nums...))), 1)
	dp := make([]bool, 1<<log)
	for _, v := range nums {
		dp[v] = true
	}
	SosDp2(log, func(cur, super int) { dp[cur] = dp[cur] || dp[super] })

	res := 0
	for _, num := range nums {
		curOr := 0
		for bit := log - 1; bit >= 0; bit-- {
			curOr |= 1 << bit    // 试填法
			if !dp[curOr&^num] { // 无法提供
				curOr ^= 1 << bit
			}
		}
		res = max(res, curOr)
	}
	return res
}

// 按位与1最多的二元组
// !要求找到两个不同的下标i≠j，使得ai&aj包含的1最多。
// nums[i]<=1e6, n<=1e6
//
// !等价于：找到某个特殊的数v，使得v的超集在数组中出现不少于两次
// !将信息下推到子集，找到这个数v后，扫描一遍数组，从中任意拎两个v的超集出来.
func BitwiseAndPairWithMaxOnesCount(nums []int) (maxOnesCount int, index1, index2 int) {
	log := max(bits.Len(uint(maxs(nums...))), 1)
	dp := make([]int, 1<<log)
	for _, v := range nums {
		dp[v]++
	}
	SosDp2(log, func(cur, super int) { dp[cur] += dp[super] })

	maxOnesCount, index1, index2 = 0, -1, -1
	bestAnd := 0
	for i, v := range dp {
		if v >= 2 {
			if c := bits.OnesCount32(uint32(i)); c > maxOnesCount {
				maxOnesCount = c
				bestAnd = i
			}
		}
	}

	for i, v := range nums {
		if v&bestAnd == bestAnd {
			if index1 == -1 {
				index1 = i
			} else {
				index2 = i
				break
			}
		}
	}

	return
}

// 按位或1最多的二元组
// !要求找到两个不同的下标i≠j，使得ai|aj包含的1最多。
// https://taodaling.github.io/blog/2019/08/23/%E4%BA%8C%E8%BF%9B%E5%88%B6%E4%BD%8D%E8%BF%90%E7%AE%97/
// nums[i]<=1e6, n<=1e6
//
// 这里的不同是没有意义的，因为要或运算后包含最多的1，那么必定会选择两个不同的值。
// 我们可以这样求解，维护一个特殊的数组b，bi记录序列a中与i进行与运算后包含最多1的数。
// 这个数组可以这样求，首先对于任意序列中的数ai，很显然与bai=ai。
// 之后我们可以进行下推操作，将bi中的数下推给所有bj，其中j∈i。之后我们再进行一次上推操作，将每个数bi上推给自己的超集。
// 上面这个算法看起来应该是O(3^k)，但是我们可以仅枚举恰好比当前集合小1和大1的子集或超集，就可以将时间复杂度降低为O(k2^k).
// 现在我们要找或运算后拥有最多1的家伙了，直接暴力遍历每个序列a中的值，记现在的ai，
// !那么要找与ai或运算后含最多1的数，实际上就是在找b[(2^k−1)⊕ai]，其中⊕表示异或运算。
func BitwiseOrPairWithMaxOnesCount(nums []int) (maxOnesCount int, index1, index2 int) {
	log := max(bits.Len(uint(maxs(nums...))), 1)
	dp := make([]int, 1<<log) // dp[i] 表示与 i 进行与运算后包含最多1的数的下标.
	for i := range dp {
		dp[i] = -1
	}
	for i, v := range nums {
		dp[v] = i
	}
	SosDp2(log, func(cur, super int) { dp[cur] = dp[super] })
	SosDp1(log, func(cur, sub int) { dp[cur] = dp[sub] })

	mask := 1<<log - 1
	maxOnesCount, index1, index2 = 0, -1, -1
	bestOr := 0
	for i := range dp {
		need := mask ^ i
		if c := bits.OnesCount32(uint32(need)); c > maxOnesCount {
			maxOnesCount = c
			bestOr = need
		}
	}

	fmt.Println(bestOr, 987)

	for i, v := range nums {
		if v|bestOr == bestOr {
			if index1 == -1 {
				index1 = i
			} else {
				index2 = i
				break
			}
		}
	}
	return
}

// https://www.luogu.com.cn/problem/CF165E
// 给定一个数组.对每个nums[i]，问是否存在nums[j]满足nums[i]&nums[j]=0.
// 存在则输出nums[j],否则输出-1
// !nums[i]<=4e6
//
// !其实就是问两个二进制集合有没有交。
func CF165E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	log := max(bits.Len(uint(maxs(nums...))), 1)
	dp := make([]int, 1<<log)
	for i := range dp {
		dp[i] = -1
	}
	for _, v := range nums {
		dp[v] = v
	}

	SosDp1(log, func(cur, sub int) {
		if dp[sub] != -1 {
			dp[cur] = dp[sub]
		}
	})

	mask := 1<<log - 1
	for _, v := range nums {
		fmt.Fprint(out, dp[mask^v], " ")
	}
}

// https://www.luogu.com.cn/problem/CF449D
// 求nums中按位与为0的非空集合个数，模1e9+7.
// (取反后等价于按位或为全集.)
// nums[i]<=1e6.
//
// !dp[i] 表示子集包含 i 的方案数.
func CF449D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e9 + 7

	pow := func(x, n int) int {
		x %= MOD
		res := 1 % MOD
		for ; n > 0; n >>= 1 {
			if n&1 == 1 {
				res = res * x % MOD
			}
			x = x * x % MOD
		}
		return res
	}

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	log := max(bits.Len(uint(maxs(nums...))), 1)
	counter := make([]int, 1<<log)
	for _, v := range nums {
		counter[v]++
	}
	SosDp2(log, func(cur, super int) {
		counter[cur] = (counter[cur] + counter[super]) % MOD
	})

	dp := make([]int, 1<<log) // dp[i] 表示子集包含 i 的方案数.
	for i := range dp {
		dp[i] = pow(2, counter[i])
	}
	SosDp2(log, func(cur, super int) {
		dp[cur] = (dp[cur] - dp[super] + MOD) % MOD
	})

	fmt.Fprintln(out, dp[0])
}

// https://www.luogu.com.cn/problem/CF383E
// 给出 n 个长度为 3 的单词（每个字母都为 a 到 x 中的任意一个），
// 我们说该单词是正确的，当且仅当该单词中含有至少一个元音。元音可以是范围内的任意字母。
// 对于所有2^24种元音集合，求出所有情况下正确单词个数的平方的异或和。
//
// !平方的异或和不太能维护，考虑暴力枚举子集计算答案。
// 正难则反，求非法的数量.
func CF383E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := range words {
		fmt.Fscan(in, &words[i])
	}

	log := 24 // 'a'-'x' 24 位
	MASK := 1<<log - 1
	dp := make([]int, 1<<log) // !非法个数
	for _, word := range words {
		state := 0
		for _, c := range word {
			state |= 1 << (c - 'a')
		}
		dp[state^MASK]++
	}
	SosDp2(log, func(cur, super int) {
		dp[cur] += dp[super]
	})
	res := 0
	for i := 0; i < 1<<log; i++ {
		x := n - dp[i]
		res ^= x * x
	}
	fmt.Fprintln(out, res)
}

// https://www.luogu.com.cn/problem/CF1234F
// 给你一个字符串S，你可以翻转一次S的任意一个子串。
// 问翻转后S的子串中各个字符都不相同的最长子串长度。
// 字符集<=20
//
// !翻转操作本质上是在 S 中找两个不相交子串，且拼接后不含相同字符。
// !等价于:给定一个数组，求满足ai&aj==0的ai|aj二进制中1的个数的最大值.
// 我们可以枚举 ai，求出 ai的补集的子集中的二进制中 1 的个数的最大值，
// 其与 ai的二进制中 1 的个数之和的最大值就是答案。
func CF1234F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	set := make(map[int]struct{})
	n := len(s)
	for left := 0; left < n; left++ {
		state := 0
		for right := left; right < n; right++ {
			cur := 1 << (s[right] - 'a')
			if state&cur > 0 {
				break
			} else {
				state |= cur
			}
			set[state] = struct{}{}
		}
	}

	nums := make([]int, 0, len(set))
	for state := range set {
		nums = append(nums, state)
	}

	// 给定一个数组，求满足ai&aj==0的ai|aj二进制中1的个数的最大值.
	// nums[i]<=1e6
	solve := func(nums []int) int {
		log := max(bits.Len(uint(maxs(nums...))), 1)
		dp := make([]int, 1<<log)
		for _, v := range nums {
			dp[v] = bits.OnesCount(uint(v))
		}

		SosDp1(log, func(cur, sub int) {
			dp[cur] = max(dp[cur], dp[sub])
		})

		res := 0
		MASK := 1<<log - 1
		for s := range dp {
			res = max(res, dp[s]+dp[s^MASK])
		}
		return res
	}

	fmt.Fprintln(out, solve(nums))
}

// https://www.luogu.com.cn/problem/AT_arc100_c
// 给定一个长为2^n的数组，对每个1<=k<=2^n-1，
// 找出最大的nums[i]+nums[j] (i|j<=k)
// n<=18
// 维护集合最大和次大即可，最终输出的每一个答案都是一个前缀最大值。
func ARC100E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, 1<<n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	log := n
	type pair struct{ max1, max2 int }
	dp := make([]pair, 1<<log)
	for i, v := range nums {
		dp[i] = pair{max1: v, max2: -1}
	}
	SosDp1(log, func(cur, sub int) {
		curDp, subDp := dp[cur], dp[sub]
		if subDp.max2 > curDp.max1 {
			dp[cur] = subDp
		} else if subDp.max1 > curDp.max1 {
			dp[cur] = pair{subDp.max1, curDp.max1}
		} else if subDp.max1 > curDp.max2 {
			dp[cur].max2 = subDp.max1
		}
	})

	res := 0
	for i := 1; i < 1<<log; i++ {
		res = max(res, dp[i].max1+dp[i].max2)
		fmt.Fprintln(out, res)
	}
}

// https://atcoder.jp/contests/arc136/tasks/arc136_d
// 10进制的情形

// 有多少对按位与为0的二元组/与运算为0的二元组个数
// nums[i]<=1e6, n<=1e6
// !信息上推到超集，再统计答案，注意除以2.
func 按位与为0的二元组(nums []int) int {
	log := max(bits.Len(uint(maxs(nums...))), 1)
	dp := make([]int, 1<<log)
	for _, v := range nums {
		dp[v]++
	}
	SosDp1(log, func(cur, sub int) {
		dp[cur] += dp[sub]
	})
	mask := 1<<log - 1
	res := 0
	for _, v := range nums {
		res += dp[mask^v]
	}
	return res / 2
}

// 从子集转移的 SOS DP，将信息上推到超集
// SubSetZeta
func SosDp1(log int, f func(cur, sub int)) {
	for i := 0; i < log; i++ {
		mask := 1 << i
		for s := 0; s < 1<<log; s++ {
			s |= mask
			f(s, s^mask) // 将 s 的子集 s^1<<i 的统计量合并到 s 中
		}
	}
}

// 从超集转移的 SOS DP，将信息下推到子集
// SuperSetZeta
func SosDp2(log int, f func(cur, super int)) {
	for i := 0; i < log; i++ {
		mask := 1 << i
		for s := 0; s < 1<<log; s++ {
			t := s ^ mask
			if s < t {
				f(s, t) // 将 s 的超集 s|1<<i 的统计量合并到 s 中
			}
		}
	}
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

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
