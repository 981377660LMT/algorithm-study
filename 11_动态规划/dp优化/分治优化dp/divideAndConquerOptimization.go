// 决策单调性优化dp.
// 単調最小値DP (aka. 分割統治DP) 优化 offlineDp
// https://ei1333.github.io/library/dp/divide-and-conquer-optimization.hpp
// !用于高速化 dp[k][j]=min(dp[k-1][i]+f(i,j)) (0<=i<j<=n) => !将区间[0,n)分成k组的最小代价
//
//	如果f满足决策单调性 那么对转移的每一行，可以采用 monotoneminima 寻找最值点
//	O(kn^2)优化到O(knlogn)
//
// https://www.cnblogs.com/purplevine/p/16990286.html
// https://www.cnblogs.com/alex-wei/p/DP_optimization_method_II.html

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func main() {
	CF833B()
	// CF868F()

	// P4360()
	// P5574()
}

// Ciel and Gondolas
// https://www.luogu.com.cn/problem/CF321E
func CF321E() {}

// CF833B-The Bakery (决策单调性+莫队维护区间颜色个数)
// https://www.luogu.com.cn/problem/CF833B
// 将一个数组分为k段，使得总价值最大。
// 一段区间的价值表示为区间内不同数字的个数。
// n=3e4,k<=50
//
// dp[i][j]=max{dp[i-1][k]+cost(k+1,j) 1<=k<j
// dp[i][j]意为前j个数被分成i段时的最大总价值.
func CF833B() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	D := NewDictionary()
	for i, v := range nums {
		nums[i] = D.Id(v)
	}

	counter := make([]int, D.Size())
	left, right, kind := 0, 0, 0
	add := func(i int) {
		counter[nums[i]]++
		if counter[nums[i]] == 1 {
			kind++
		}
	}
	remove := func(i int) {
		counter[nums[i]]--
		if counter[nums[i]] == 0 {
			kind--
		}
	}
	f := func(l, r int, _ int) int {
		for left > l {
			left--
			add(left)
		}
		for right < r {
			add(right)
			right++
		}
		for left < l {
			remove(left)
			left++
		}
		for right > r {
			right--
			remove(right)
		}
		return -kind // 要求最大值，因此取负
	}

	dp := DivideAndConquerOptimization(k, n, f)
	fmt.Fprintln(out, -dp[k][n])
}

// Yet Another Minimization Problem (决策单调性+莫队维护相同元素的对数)
// https://www.luogu.com.cn/problem/CF868F
// 有一个长度为 n 的序列，要求将其分成 k 个子段，每个子段的花费是子段内相同元素的对数，求最小花费。
// dp[k][i] 表示前 i 个元素分成 k 个子段的最小花费。
func CF868F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	D := NewDictionary()
	for i, v := range nums {
		nums[i] = D.Id(v)
	}

	counter := make([]int, D.Size())
	left, right, cost := 0, 0, 0
	add := func(i int) {
		cost += counter[nums[i]]
		counter[nums[i]]++
	}
	remove := func(i int) {
		counter[nums[i]]--
		cost -= counter[nums[i]]
	}
	f := func(l, r int, _ int) int {
		for left > l {
			left--
			add(left)
		}
		for right < r {
			add(right)
			right++
		}
		for left < l {
			remove(left)
			left++
		}
		for right > r {
			right--
			remove(right)
		}
		return cost
	}

	dp := DivideAndConquerOptimization(k, n, f)
	fmt.Fprintln(out, dp[k][n])
}

// P4360 [CEOI2004] 锯木厂选址
// https://www.luogu.com.cn/problem/P4360
func P4360() {}

// P5574 [CmdOI2019] 任务分配问题 (决策单调性+莫队+树状数组维护逆序对)
// https://www.luogu.com.cn/problem/P5574
func P5574() {}

// !f(i,j,step): 左闭右开区间[i,j)的代价(0<=i<j<=n)
//
//	可选:step表示当前在第几组(1<=step<=k)
func DivideAndConquerOptimization(k, n int, f func(i, j, step int) int) [][]int {
	dp := make([][]int, k+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	dp[0][0] = 0

	for k_ := 1; k_ <= k; k_++ {
		getCost := func(y, x int) int {
			if x >= y {
				return INF
			}
			return dp[k_-1][x] + f(x, y, k_)
		}
		res := monotoneminima(n+1, n+1, getCost)
		for j := 0; j <= n; j++ {
			dp[k_][j] = res[j][1]
		}
	}

	return dp
}

// 对每个 0<=i<H 求出 f(i,j) 取得最小值的 (j, f(i,j)) (0<=j<W)
func monotoneminima(H, W int, f func(i, j int) int) [][2]int {
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
			tmp := f(mid, i)
			if index == -1 || tmp < res { // !less if get min
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

type V = int
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}
func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary) Size() int {
	return len(d._idToValue)
}
