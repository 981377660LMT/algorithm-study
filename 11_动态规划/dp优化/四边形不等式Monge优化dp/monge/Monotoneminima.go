// 决策单调性(Monotone)函数f(i,j)的最小值
// Monge(四边形不等式) ⇒ Totally Monotone(TM) ⇒ Monotone なので、Monotone は弱い条件である。
// https://ei1333.github.io/luzhiled/snippets/dp/monotone-minima.html
// https://beet-aizu.github.io/library/algorithm/monotoneminima.cpp
// https://noshi91.github.io/algorithm-encyclopedia/monotone-minima

// 对于一个二元函数f(i,j) (0<=i<H, 0<=j<W),
// !如果对任意 p<q 满足 argmin(f(p,*))<=argmin(f(q,*)),
// !即f(i,j)取到最小值时,i如果变大,决策点j也随之变大,则称f(i,j)是关于i Monotone 的.
// !例如 f(i,j)=nums[j]+(j-i)^2 是关于i的Monotone函数(一次函数)
// 绝大多数的斜率优化也满足决策单调性（维护的下凸包，斜率也满足单调不减）

// 如果矩阵满足totally monotone,则可以用SMAWK算法O(H+W)求出每一行的最小值.

package main

import (
	"bufio"
	"fmt"
	"os"
)

// 对每个查询 0<=qi<q 求出 f(i,j) 取得最小值时的 (j, f(i,j)) (0<=j<W)
//
//	!f(i,j): 0<=i<=Q-1, 0<=j<=W-1
func Monotoneminima(Q, W int, f func(i, j int) int, isMin bool) [][2]int {
	dp := make([][2]int, Q) // dp[i] 表示第i行取到`最小值`的(索引,值)

	var dfs func(top, bottom, left, right int)
	dfs = func(top, bottom, left, right int) {
		if top > bottom {
			return
		}

		mid := (top + bottom) / 2
		index := -1
		res := 0
		for i := left; i <= right; i++ {
			tmp := f(mid, i)
			if isMin {
				if index == -1 || tmp < res {
					index = i
					res = tmp
				}
			} else {
				if index == -1 || tmp > res {
					index = i
					res = tmp
				}
			}
		}
		dp[mid] = [2]int{index, res}
		dfs(top, mid-1, left, index)
		dfs(mid+1, bottom, index, right)
	}

	dfs(0, Q-1, 0, W-1)
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
	res := Monotoneminima(n, n, dist, true)
	for i := 0; i < n; i++ {
		fmt.Fprintln(out, res[i][1])
	}
}
