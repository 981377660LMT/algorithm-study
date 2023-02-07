// 决策单调性(Monotone)
// Monge(四边形不等式) ⇒ Totally Monotone(TM) ⇒ Monotone なので、Monotone は弱い条件である。
// https://ei1333.github.io/luzhiled/snippets/dp/monotone-minima.html
// https://beet-aizu.github.io/library/algorithm/monotoneminima.cpp

// 对于一个二元函数f(i,j) (0<=i<H, 0<=j<W),
// !如果对任意 p<q 满足 argmin(f(p,*))<=argmin(f(q,*)),
// !即f(i,j)取到最小值时,j如果变大,i也变大,则称f(i,j)是关于i Monotone 的.
// !例如 f(i,j)=nums[j]+(j-i)^2 是关于i的Monotone函数(一次函数)

package main

import (
	"bufio"
	"fmt"
	"os"
)

// !dist(i,j): 闭区间[i,j]的代价(0<=i<=j<=W-1)
func monotoneminima(H, W int, dist func(i, j int) int) [][2]int {
	dp := make([][2]int, H) // dp[i] 表示第i行取到`最小值`的(索引,值)

	var dfs func(top, bottom, left, right int)
	dfs = func(top, bottom, left, right int) {
		if top > bottom {
			return
		}

		mid := (top + bottom) / 2
		index := -1
		res := 0
		for i := left; i <= right; i++ {
			tmp := dist(mid, i)
			if index == -1 || tmp < res { // less
				index = i
				res = tmp
			}
		}
		dp[mid] = [2]int{index, res}
		dfs(top, mid-1, left, index)
		dfs(mid+1, bottom, index, right)
	}

	dfs(0, H-1, 0, W-1)
	return dp
}

func main() {
	// AtCoder COLOCON -Colopl programming contest 2018- Final C - スペースエクスプローラー高橋君
	// https://blog.hamayanhamayan.com/entry/2018/01/21/161336
	// 给定长为n的数组a (n<=2e5)
	// !对每个i 求 a[j]+(j-i)^2 的最小值
	// 也可以用Convex Hull Trick 求出

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	dist := func(i, j int) int {
		return nums[j] + (j-i)*(j-i)
	}
	res := monotoneminima(n, n, dist)
	for i := 0; i < n; i++ {
		fmt.Fprintln(out, res[i][1])
	}
}
