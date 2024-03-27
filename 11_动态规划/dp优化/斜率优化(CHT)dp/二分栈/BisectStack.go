// 决策单调性之"二分栈"
// 二分栈算法常用于有如下决策单调性的 DP 问题中：
// 每个决策点j只会被比它更前的决策点i(i<j)反超。
// 通常情况下的限制为：
// 贡献函数二阶导恒为非负，求最大值 或 二阶导恒为非正，求最小值。
// 还有一个二分队列，单调性与二分栈相反，但是没有本质区别.
//
// https://wenku.baidu.com/view/681d161ca300a6c30c229f70 1D1D动态规划优化初步：
//
// 使用一个栈来维护数据，栈中的每一个元素保存一个决策的起始位置与终了位置，
// 显然这些位置相互连接且依次递增，形如"00001111112223333..."。
// 当插入一个新的决策时，从后到前扫描栈，对于每一个老决策来说，做这样两件事:
// 1、如果在老决策的起点处还是新决策更好，则退栈，全额抛弃老决策，
// 将其区间合并至新决策中，继续扫描下一个决策。
// 2、如果在老决策的起点处是老决策好，则转折点必然在这个老决策的区间中;
// 二分查找之，然后新决策进栈，结束。
// 由于一个决策出栈之后再也不会进入，所以均摊时间为O(1)，
// 但是由于二分查找的存在，所以整个算法的时间复杂度为O(nlogn)。
//
// https://www.cnblogs.com/alex-wei/p/DP_optimization_method_II.html
// https://www.cnblogs.com/flashhu/p/9480669.html
// https://blog.csdn.net/Cherrt/article/details/120916488

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
)

var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {
	P5504()
}

func BisectStack() {}

// P3195 [HNOI2008] 玩具装箱
// https://www.luogu.com.cn/problem/P3195
// !dp[j] = min(dp[i]+(preSum[j]-preSum[i]+j-i-1-C)^2)
func P3195() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, C := io.NextInt(), io.NextInt()
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = io.NextInt()
	}

	preSum := make([]int, n+1)
	for i := 0; i < n; i++ {
		preSum[i+1] = preSum[i] + nums[i]
	}

	_ = C
}

// P5504 [JSOI2011] 柠檬(二分栈, TODO：斜率优化)
// https://www.luogu.com.cn/problem/P5504
// 将一个数组分成若干段，每一段中选一个数作为这一段的代表数x.
// 假设这个数有t个，那么这一段的得分为x*t*t.
// 求划分数组的最大得分.
// 数组长度<=1e5.
//
// dp[i]表示前i个数的最大得分.
// 一个关键的结论是：只有左右端点元素相同的区间才有可能转移, 成为最佳决策.
// !dp[j] = max(dp[i] + nums[i] * (idInGroup[j]-idInGroup[i]+1)^2), 其中i<j且nums[i]==nums[j].
// 注意，需要给每个v值维护一个单调栈，栈中存储的是区间的左端点.
func P5504() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()
	nums := make([]int, n)
	for i := range nums {
		nums[i] = io.NextInt()
	}
	D := NewDictionary()
	for i, v := range nums {
		nums[i] = D.Id(v)
	}

	counter := make([]int, D.Size())
	idInGroup := make([]int, n)
	for i, v := range nums {
		idInGroup[i] = counter[v]
		counter[v]++
	}

	dp := make([]int, n+1)
	cal := func(i, len int) int {
		return dp[i-1] + len*len*nums[i]
	}

	// TODO
	check := func(i, j int) int {
		return 0
	}

	stacks := make([][]int, D.Size())
	top := func(v int) int { return stacks[v][len(stacks[v])-1] }
	second := func(v int) int { return stacks[v][len(stacks[v])-2] }
	for i, v := range nums {
		for len(stacks[v]) > 1 && check(second(v), top(v)) <= check(top(v), i) {
			stacks[v] = stacks[v][:len(stacks[v])-1]
		}
		stacks[v] = append(stacks[v], i)
		for len(stacks[v]) > 1 && check(second(v), top(v)) <= idInGroup[i] {
			stacks[v] = stacks[v][:len(stacks[v])-1]
		}
		dp[i] = cal(top(v), idInGroup[i]-idInGroup[top(v)]+1)
	}

	io.Println(dp[n])
}

// P5665 [CSP-S2019] 划分
// https://www.luogu.com.cn/problem/P5665
// 将一段长度为n的序列划分成若干个区间，使得区间sum递增，且所有区间的平方和最小
func P5665() {}

// 2945. 找到最大非递减数组的长度
// https://leetcode.cn/problems/find-maximum-non-decreasing-array-length/description/
func findMaximumLength(nums []int) int {

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
