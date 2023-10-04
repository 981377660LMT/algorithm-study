// https://www.luogu.com.cn/problem/CF786C
package main

func main() {

}

// 如果 f(i) 的计算结果近似 n/i
// 可以对 [1,n] 值域分治，如果区间内的结果都相同，则不再分治
// 时间复杂度类似整除分块 O(f(n)√n)
// https://codeforces.com/problemset/problem/786/C
func floorDivide(n int, f func(int) int) []int {
	ans := make([]int, n+1)
	var solve func(int, int)
	solve = func(l, r int) {
		if l > r {
			return
		}
		resL, resR := f(l), f(r)
		if resL == resR {
			for i := l; i <= r; i++ {
				ans[i] = resL
			}
			return
		}
		ans[l] = resL
		ans[r] = resR
		mid := (l + r) / 2
		solve(l+1, mid)
		solve(mid+1, r-1)
	}
	solve(1, n)
	return ans[1:]
}
