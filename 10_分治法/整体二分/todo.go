package main

func main() {

}

// 整体二分 Parallel Binary Search
// https://oi-wiki.org/misc/parallel-binsearch/
// https://codeforces.com/blog/entry/45578
// todo 整体二分解决静态区间第 k 小的优化 https://www.luogu.com.cn/blog/2-6-5-3-5/zheng-ti-er-fen-xie-jue-jing-tai-ou-jian-di-k-xiao-di-you-hua
// 模板题 https://www.luogu.com.cn/problem/P3527
// todo https://atcoder.jp/contests/agc002/tasks/agc002_d
//  https://www.hackerrank.com/contests/hourrank-23/challenges/selective-additions/problem
//  https://www.codechef.com/problems/MCO16504
func parallelBinarySearch(n int, qs []struct{ l, r, v int }) []int {
	// 读入询问时可以处理成左闭右开的形式

	ans := make([]int, n)
	tar := make([]int, n)
	for i := range tar {
		tar[i] = i
	}
	var f func([]int, int, int)
	f = func(tar []int, ql, qr int) {
		if len(tar) == 0 {
			return
		}
		if ql+1 == qr {
			for _, c := range tar {
				ans[c] = ql // qr
			}
			return
		}
		qm := (ql + qr) / 2
		for _, q := range qs[ql:qm] {
			_ = q
			// apply(q)

		}

		// 根据此刻查询的结果将 tar 分成左右两部分
		var left, right []int
		for _, who := range tar {
			_ = who

		}

		for _, q := range qs[ql:qm] {
			_ = q
			// rollback(q)

		}
		f(left, ql, qm)
		f(right, qm, qr)
	}
	f(tar, 0, len(qs)+1) // 这样可以将无法满足要求的 ans[i] 赋值为 len(qs)
	return ans
}
