// No.1826 Fruits Collecting-打地鼠(一维)
// https://yukicoder.me/problems/no/1826
// https://atcoder.jp/contests/abc266/editorial/4664 (打地鼠二维)

// 一个人在原点,起始坐标时刻都为0,每次可以向左或向右移动一格或者不动,每次移动需要1秒
// 有n个地鼠,分别在时刻t[i]于坐标x[i]出现,打中这个地鼠可以获得分数s[i]
// 求出最大的分数
// n<=2e5 ti<=1e9 -1e9<=xi<=1e9 1<=si<=1e9

// 记dp[t][x]为t时刻在x位置的最大分数
// !则 dp[ti][xi] = max(dp[tj][xj])+si
// 其中第j个地鼠后可以打到第i个地鼠的条件为
// xi-(tj-ti)<=xj<=xi+(tj-ti)
// !令 Ai = ti+xi, Bi = ti-xi (45度旋转)
// 将式子变形得到 Aj>=Ai 且 Bj>=Bi
// !则 dp[Ai][Bi] = max(dp[Aj][Bj])+si  (Aj>=Ai 且 Bj>=Bi)
// - 解法1：将Ai和Bi排序遍历,用一个数据结构维护Bi上的前缀最大值
// - 解法2: 使用二维树状数组求二维前缀最大值
// !45°回転すると平面操作っぽくなる

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

	var n int
	fmt.Fscan(in, &n)
	goods := make([][3]int, n) // time, pos, score
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &goods[i][0], &goods[i][1], &goods[i][2])
	}

	txs := make([][3]int, 0, n)
	set := map[int]struct{}{0: {}} // start pos
	for i := 0; i < n; i++ {
		t, x, s := goods[i][0], goods[i][1], goods[i][2]
		if abs(x) <= t { // reachable
			txs = append(txs, [3]int{t + x, t - x, s})
			set[t-x] = struct{}{}
		}
	}

	sort.Slice(txs, func(i, j int) bool {
		if txs[i][0] == txs[j][0] {
			return txs[i][1] < txs[j][1]
		}
		return txs[i][0] < txs[j][0]
	})
	sorted := make([]int, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)
	mp := make(map[int]int, len(sorted))
	for i, v := range sorted {
		mp[v] = i
	}

	seg := NewFenwickTreePrefix(len(sorted))
	seg.Set(0, 0)
	for i := 0; i < len(txs); i++ {
		x, s := txs[i][1], txs[i][2]
		pos := mp[x]
		curMax := seg.Query(pos+1) + s
		seg.Set(pos, curMax)
	}
	fmt.Fprintln(out, seg.Query(len(sorted)))

}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type S = int

const INF int = 1e18

func (*FenwickTreePrefix) e() S        { return -INF }
func (*FenwickTreePrefix) op(a, b S) S { return max(a, b) }
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type FenwickTreePrefix struct {
	n    int
	data []S
}

func NewFenwickTreePrefix(n int) *FenwickTreePrefix {
	res := &FenwickTreePrefix{n, make([]S, n+1)}
	for i := 0; i < n+1; i++ {
		res.data[i] = res.e()
	}
	return res
}

func NewFenwickTreePrefixWithSlice(nums []S) *FenwickTreePrefix {
	n := len(nums)
	res := &FenwickTreePrefix{n, make([]S, n+1)}
	for i := 1; i < n+1; i++ {
		res.data[i] = nums[i-1]
	}
	for i := 1; i < n+1; i++ {
		if j := i + (i & -i); j <= n {
			res.data[j] = res.op(res.data[j], res.data[i])
		}
	}
	return res
}

// 单点更新index处的元素.
// 0 <= index < n
func (f *FenwickTreePrefix) Set(index int, value S) {
	for index++; index <= f.n; index += index & -index {
		f.data[index] = f.op(f.data[index], value)
	}
}

// 查询前缀区间 [0,right) 的值.
// 0 <= right <= n
func (f *FenwickTreePrefix) Query(right int) S {
	res := f.e()
	if right > f.n {
		right = f.n
	}
	for ; right > 0; right -= right & -right {
		res = f.op(res, f.data[right])
	}
	return res
}
