// offline dp: 区间分成k等分 dp[k][j]=min(dp[k-1][i]+f(i,j)) (0<=i<j) (dfsIndexRemain)
// online dp: 区间分成任意分 dp[j]=min(dp[i]+f(i,j)) (0<=i<j) (dfsIndex)

// !https://qiita.com/tmaehara/items/0687af2cfb807cde7860 (深度好文)
// https://beet-aizu.github.io/library/algorithm/offlineonline.cpp
// https://ei1333.github.io/library/dp/online-offline-dp.hpp

// オフライン・オンライン変換：
// !如果offline问题存在复杂度O(M(n))的解,那么online问题存在复杂度O(M(n)logn)的解
// dp[j]=min(dp[i]+f(i,j)) (0<=i<j)
// O(n^2)优化到O(nlogn^2)
// 例子:
// !dp[j]=min(dp[i]+(x[j]-x[i]-a)^2)

package main

import (
	"bufio"
	"fmt"
	"os"
)

// !dist(i,j): 左闭右开区间[i,j)的代价(0<=i<j<=n)
func offlineOnlineDp(n int, dist func(i, j int) int) int {
	dp := make([]int, n+1)
	used := make([]bool, n+1)
	used[n] = true

	update := func(k, val int) {
		if !used[k] {
			dp[k] = val
		}
		dp[k] = min(dp[k], val) // min if get min
		used[k] = true
	}

	var dfs func(top, bottom, left, right int) // induce
	dfs = func(top, bottom, left, right int) {
		if top == bottom {
			return
		}
		mid := (top + bottom) / 2
		index := left
		res := dist(mid, index) + dp[index]
		for i := left; i <= right; i++ {
			tmp := dist(mid, i) + dp[i]
			if tmp < res { // !less if get min
				res = tmp
				index = i
			}
		}

		update(mid, res)
		dfs(top, mid, left, index)
		dfs(mid+1, bottom, index, right)
	}

	var solve func(left, right int)
	solve = func(left, right int) {
		if left+1 == right {
			update(left, dist(left, right)+dp[right])
			return
		}
		mid := (left + right) / 2
		solve(mid, right)
		dfs(left, mid, mid, right)
		solve(left, mid)
	}

	solve(0, n)
	return dp[0]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// 垃圾回收
	// https://yukicoder.me/problems/no/705
	// 公园里有n个垃圾
	// n个人准备去回收垃圾,
	// 每个人的起点在(startsi,0),沿x轴顺序递增排列
	// 每个垃圾的位置在(xi,yi),沿x轴顺序递增排列
	// !每个人只能回收他左边的垃圾
	// 每个人i回收垃圾j,j+1,...,i(0<=j<=i)的时间花费为
	// 1.曼哈顿距离的平方2.曼哈顿距离3.曼哈顿距离的三次方
	// !求回收所有垃圾的最短时间之和
	// dp[i]表示前i个人回收完前i个垃圾时的最短花费
	// dp[i]=min(dp[j]+abs((xi-xj))^2+yj^2) (0<=j<=i)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	starts, xs, ys := make([]int, n), make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &starts[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &xs[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &ys[i])
	}

	dist := func(i, j int) int {
		// 0<=i<j<=n
		a := abs(starts[j-1] - xs[i])
		b := abs(ys[i])
		return a + b
	}
	res := offlineOnlineDp(n, dist)
	fmt.Fprintln(out, res)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
