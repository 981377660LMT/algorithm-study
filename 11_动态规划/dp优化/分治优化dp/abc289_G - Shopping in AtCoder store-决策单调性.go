// G - Shopping in AtCoder store
// 商品定价
// n个顾客，每个人有一个购买的欲望wanti,
// m件物品，每一件物品有一个价值pricei,
// !每一个顾客会买商品当且仅当 wanti + pricei >= 定价
// !现在要求对每一个商品定价，求出它的最大销售值（数量*定价）

// 1<=n,m<=2e5
// https://www.cnblogs.com/linyihdfj/p/17114515.html

// 1. 每一个商品的定价一定是从 wanti+pricei 中选出 (否则可以抬高定价)
//    !对于商品i，我们的销售额就是 f(i,j) = max{(j+1)*(wj+pi)}(0<=j<n) (j为购买人数)
// 2. 将价格和欲望排序后，函数f(i,j)具有决策单调性
//    对于每个查询qi,当原价增大时,取到最大值时的购买人数j也会增大,因此f(i,j)是关于i的Monotone函数
//    所以可以分治解决(Monotone Minima)
// !不要忘记排序

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	wants := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &wants[i])
	}
	prices := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &prices[i])
	}

	res := shoppingInAtCoderStore(wants, prices)
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

func shoppingInAtCoderStore(wants []int, prices []int) []int {
	sort.Slice(wants, func(i, j int) bool {
		return wants[i] > wants[j]
	})
	qi := make([]int, len(prices))
	for i := range qi {
		qi[i] = i
	}
	sort.Slice(qi, func(i, j int) bool {
		return prices[qi[i]] < prices[qi[j]]
	})
	sort.Ints(prices)

	// f(i,j) = max{(j+1)*(wj+pi)}(0<=j<n)
	mm := monotoneminima(len(prices), len(wants), func(i, j int) int {
		return (j + 1) * (wants[j] + prices[i])
	}, false)
	res := make([]int, len(prices))
	for i := 0; i < len(prices); i++ {
		res[qi[i]] = mm[i][1]
	}
	return res
}

// 对每个查询 0<=qi<q 求出 f(i,j) 取得最小值时的 (j, f(i,j)) (0<=j<W)
//  !f(i,j): 0<=i<=Q-1, 0<=j<=W-1
func monotoneminima(Q, W int, f func(i, j int) int, isMin bool) [][2]int {
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
