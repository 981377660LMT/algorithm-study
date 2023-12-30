package main

func main() {
	// 又叫高速ゼータ変換

	// 高维前缀和 SOS DP (Sum over Subsets)
	// 给一个集合，对该集合的所有子集，计算该子集的所有子集之和（这个「和」不一定是加法，可以是其它的满足合并性质的统计量）
	// https://codeforces.com/blog/entry/45223
	// Some SOS DP Insights https://codeforces.com/blog/entry/105247
	// 大量习题 https://blog.csdn.net/weixin_38686780/article/details/100109753
	//
	// https://codeforces.com/problemset/problem/1234/F
	//    求满足 ai&aj=0 的 ai|aj 的二进制 1 的个数的最大值
	//    思路是转换成求每个 ai 的补集的 SOS，维护子集二进制 1 的个数的最大值
	// https://www.hackerearth.com/zh/problem/algorithm/special-pairs-5-3ee6b3fe-3d8a1606/
	//    求 ai&aj=0 的 (i,j) 对数，0<=ai<=1e6
	//    思路和上面类似，转换成求每个 ai 的补集的 SOS
	//    注：另一种解法是求 FWT(cnt)[0]
	// 转换成求集合中最大次大 https://atcoder.jp/contests/arc100/tasks/arc100_c
	// 求下标最大次大，且不需要在乎 k 的上限的写法 https://codeforces.com/problemset/problem/1554/B
	// https://codeforces.com/problemset/problem/165/E
	// 容斥 https://codeforces.com/problemset/problem/449/D
	// todo https://codeforces.com/problemset/problem/1208/F
	//  https://codeforces.com/problemset/problem/800/D
	//  https://codeforces.com/problemset/problem/383/E
	//  https://www.luogu.com.cn/problem/P6442
	// https://codeforces.com/problemset/problem/1523/D
	// 十进制 https://atcoder.jp/contests/arc136/tasks/arc136_d
	sosDP := func(a []int) []int {
		// 从子集转移的写法
		const mx = 20 // bits.Len(uint(max(a))
		dp := make([]int, 1<<mx)
		for _, v := range a {
			dp[v]++
		}
		for i := 0; i < mx; i++ {
			for s := 0; s < 1<<mx; s++ {
				s |= 1 << i
				// 将 s 的子集 s^1<<i 的统计量合并到 s 中
				dp[s] += dp[s^1<<i]
			}
		}

		{
			// 从超集转移的写法
			for i := 0; i < mx; i++ {
				for s := 1<<mx - 1; s >= 0; s-- {
					if s>>i&1 == 0 {
						dp[s] += dp[s|1<<i]
					}
				}
			}
		}

		{
			// 维护集合最大和次大的写法
			type pair struct{ fi, se int }
			dp := make([]pair, 1<<mx)
			for i := 0; i < mx; i++ {
				for s := 0; s < 1<<mx; s++ {
					s |= 1 << i
					p, q := dp[s], dp[s^1<<i]
					if q.se > p.fi {
						dp[s] = q
					} else if q.fi > p.fi {
						dp[s] = pair{q.fi, p.fi}
					} else if q.fi > p.se {
						dp[s].se = q.fi
					}
				}
			}
		}

		return dp
	}
}
