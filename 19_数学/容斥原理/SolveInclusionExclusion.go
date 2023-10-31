package main

func main() {

}

const MOD int = 998244353

// 容斥原理 (PIE, the principle of inclusion and exclusion)
// 参考《挑战程序设计竞赛》P296
// https://codeforces.com/blog/entry/64625
// https://ac.nowcoder.com/acm/contest/6219/C
//
// 多重集组合数 https://codeforces.com/problemset/problem/451/E
// https://codeforces.com/problemset/problem/1342/E
// 如何将问题转化成可以容斥的结构 https://codeforces.com/problemset/problem/1228/E
// 不重不漏 https://codeforces.com/problemset/problem/1007/B
// 与 SOS DP 结合 https://codeforces.com/problemset/problem/449/D
// 用因子容斥 https://codeforces.com/problemset/problem/900/D
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
