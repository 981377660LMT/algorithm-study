package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// K匿名序列()
	// 仓库建设()
	// 青蛙跳()
	// 特别行动队()
	打印文章()
}

// K匿名序列(带滑窗大小限制的CHT)
// https://www.acwing.com/problem/content/description/336/
// !给定一个长度为 n 的非严格递增整数序列，
// 每次操作可以将其中的一个数减少一，
// 问最少多少次操作后能够使得序列中的任何一个数在序列中都至少有 k−1 个数与之相同。
//
// 就相当于是把一个数组切割成若干组，操作之后每组都有至少k−1个数字相同，
// 花费就是这组的数字和sum，
// 再减去最小值min乘以这个组的cnt，
// 也就是sum−(min∗cnt)。
// 分成的组一定是连续的一段，那么就可以设计转移方程了
// dp[j] = min(dp[i] + (p[j] - p[i]) - nums[i]*(j-i)) 其中 j-i>=k
// 分离ij得到
// !dp[j] = min(-nums[i]*j + nums[i]*i + dp[i] - p[i] + p[j]) 其中 j-i>=k
// !直线为(-nums[i], nums[i]*i + dp[i] - p[i])
//
// !j-i>=k条件的限制：当j入队时，i=j-k+1，i>=0才能入队.
func K匿名序列() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const INF int = 1e18

	solve := func(nums []int, k int) int {
		n := len(nums)
		preSum := make([]int, n+1)
		for i := 1; i <= n; i++ {
			preSum[i] = preSum[i-1] + nums[i-1]
		}

		dp := make([]int, n+1)
		for i := range dp {
			dp[i] = INF
		}
		dp[0] = 0
		cht := NewConvexHullTrickDeque(true)
		for j := 0; j <= n; j++ {
			if i := j - k + 1; i >= 0 {
				best, _ := cht.QueryMonotoneInc(j)
				dp[j] = best + preSum[j]
				cht.AddLineMonotone(-nums[i], nums[i]*i+dp[i]-preSum[i], -1)
			}
		}
		return dp[n]
	}

	var T int
	fmt.Fscan(in, &T)
	for i := 0; i < T; i++ {
		var n, k int
		fmt.Fscan(in, &n, &k)
		nums := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &nums[j])
		}
		fmt.Fprintln(out, solve(nums, k))
	}
}

// 仓库建设(搬运货物)
// https://www.acwing.com/file_system/file/content/whole/index/content/4184038/
// 有 n 个工厂，由高到低分布在一座山上，工厂 1 在山顶，工厂 n 在山脚。
// 第 i 个工厂目前有成品 pi 件，在第 i 个工厂位置建立仓库的费用是 ci.
// 对于没有建立仓库的工厂，其产品被运往其他的仓库，
// 产品只能往山下运（只能运往编号更大的工厂的仓库），一件产品运送一个单位距离的费用是 1.
// 假设建立的仓库容量都足够大。工厂 i 与 1 的距离是 xi，问总费用最小值。
//
// 假设在i建工厂那么
// dp[i]=min(dp[j]+ xi*∑pi -∑(xk*pk)) + ci  (xi表示工厂i到1的距离，pk表示工厂k的产品数量)
// 对数量求前缀和得
// dp[i]=min(dp[j]+ dist[i]*(countPreSum[i]-countPreSum[j]) - (feePreSum[i]-feePreSum[j])) + costi
// 分离ij得到
// !dp[i]= -countPreSum[j]*dist[i] + dp[j] + feePreSum[j] - feePreSum[i] + costi + dist[i]*countPreSum[i]
// 直线为(-countPreSum[j], dp[j]+feePreSum[j]), x为dist[i]
func 仓库建设() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	factories := make([][3]int, 0, n) // dist, count, cost
	for i := 0; i < n; i++ {
		var dist, count, cost int
		fmt.Fscan(in, &dist, &count, &cost)
		factories = append(factories, [3]int{dist, count, cost})
	}

	countPreSum, feePreSum := make([]int, n+1), make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist, count := factories[i-1][0], factories[i-1][1]
		countPreSum[i] = countPreSum[i-1] + count
		feePreSum[i] = feePreSum[i-1] + count*dist
	}

	dp := make([]int, n+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	cht := NewConvexHullTrickDeque(true)
	for i := 0; i <= n; i++ {
		if i > 0 {
			dist, cost := factories[i-1][0], factories[i-1][2]
			best, _ := cht.QueryMonotoneInc(dist)
			dp[i] = best - feePreSum[i] + cost + dist*countPreSum[i]
		}
		if i < n {
			cht.AddLineMonotone(-countPreSum[i], dp[i]+feePreSum[i], -1)
		}
	}

	// !注意如果没有货物，可以不用建立仓库
	res := dp[n]
	for i := n - 1; i >= 0; i-- {
		count := factories[i][1]
		if count != 0 {
			break
		}
		if dp[i] < res {
			res = dp[i]
		}
	}
	fmt.Fprintln(out, res)
}

// !青蛙跳(青蛙过河3，代价为距离的平方和)
// !h1<h2<...<hn
// https://atcoder.jp/contests/dp/tasks/dp_z
// dp[j] = min(dp[j], dp[i] + c + (h[j] - h[i])^2) (0<=i<j<=n-1)
// !分离ij变形: dp[j]= min(-2*h[i]*h[j] + (h[i]^2+dp[i]) + h[j]^2 + c)
// !直线为(-2*h[i], h[i]^2+dp[i]), x为h[j]
func 青蛙跳() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, c int
	fmt.Fscan(in, &n, &c)
	heights := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &heights[i])
	}

	dp := make([]int, n)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	// cht := NewConvexHullTrickDeque(true)

	// !height不满足递增时，需要用lichao树
	hMin, hMax := heights[0], heights[0]
	for i := 1; i < n; i++ {
		h := heights[i]
		if h < hMin {
			hMin = h
		}
		if h > hMax {
			hMax = h
		}
	}
	cht := NewConvexHullTrickLichao(true, hMin, hMax+1)

	for i := 0; i < n; i++ {
		h := heights[i]
		if i > 0 {
			best, _ := cht.Query(h)
			dp[i] = best + h*h + c
		}
		cht.AddLine(-2*h, h*h+dp[i], -1)
	}

	fmt.Fprintln(out, dp[n-1])
}

// 特别行动队
// https://www.acwing.com/solution/content/107979/
// 记p为前缀和
// dp[j]=max(dp[j],dp[i]+a*(p[j]-p[i])^2+b*(p[j]-p[i])+c) (1<=i<j<=n)
// !变形: dp[j]=max(-2*a*pi*pj + (a*pi^2+dpi-b*pi) + a*pj^2+b*pj+c)
// !直线为(-2*a*pi, a*pi^2+dpi-b*pi), x为pj
func 特别行动队() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var a, b, c int
	fmt.Fscan(in, &a, &b, &c)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	preSum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		preSum[i] = preSum[i-1] + nums[i-1]
	}

	cht := NewConvexHullTrickDeque(false) // 注意这里是最大值
	dp := make([]int, n+1)
	dp[0] = 0
	for i := 0; i <= n; i++ {
		if i > 0 {
			p := preSum[i]
			best, _ := cht.QueryMonotoneInc(p)
			dp[i] = best + a*p*p + b*p + c
		}
		if i < n {
			p := preSum[i]
			cht.AddLineMonotone(-2*a*p, a*p*p-b*p+dp[i], i)
		}
	}

	fmt.Fprintln(out, dp[n])
}

// 打印文章
// https://www.acwing.com/problem/content/1096/
// 记p为前缀和
// dp[0]=0,dp[j]=min(dp[j],dp[i]+(p[j]-p[i])^2+c) (1<=i<j<=n)
// !变形: dp[j]=min(-2*pi*pj + (pi^2+dpi) + pj^2+c)
// !直线为(-2*pi, pi^2+dpi), x为pj
func 打印文章() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	preSum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		preSum[i] = preSum[i-1] + nums[i-1]
	}

	cht := NewConvexHullTrickDeque(true)
	dp := make([]int, n+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	for i := 0; i <= n; i++ {
		if i > 0 {
			p := preSum[i]
			best, _ := cht.QueryMonotoneInc(p)
			dp[i] = best + p*p + m
		}
		if i < n {
			p := preSum[i]
			cht.AddLineMonotone(-2*p, p*p+dp[i], i)
		}
	}

	fmt.Fprintln(out, dp[n])
}

//

const INF int = 1e18

type Line struct{ k, b, id int }
type LichaoNode struct {
	line Line
	l, r *LichaoNode
}
type ConvexHullTrickLichao struct {
	isMin        bool
	lower, upper int
	root         *LichaoNode
}

// 根据待查询的自变量x的上下界[lower,upper]建立CHTLichao.
func NewConvexHullTrickLichao(isMin bool, lower, upper int) *ConvexHullTrickLichao {
	return &ConvexHullTrickLichao{isMin: isMin, lower: lower, upper: upper}
}

// O(logN) 追加一条直线k*x+b, id为直线的编号.
func (cht *ConvexHullTrickLichao) AddLine(k, b, id int) {
	if !cht.isMin {
		k, b = -k, -b
	}
	line := Line{k, b, id}
	cht.root = cht.addLine(cht.root, line, cht.lower, cht.upper, cht.getY(line, cht.lower), cht.getY(line, cht.upper))
}

// O(logN^2) 追加一条左闭右开的线段[start,end)，所在直线k*x+b, id为线段的编号.
func (cht *ConvexHullTrickLichao) AddSegment(start, end, k, b, id int) {
	if !cht.isMin {
		k, b = -k, -b
	}
	line := Line{k, b, id}
	cht.root = cht.addSegment(cht.root, line, start, end-1, cht.lower, cht.upper, cht.getY(line, cht.lower), cht.getY(line, cht.upper))
}

// O(logN) 查询k*x+b的最小/大值.如果不存在直线,返回的id为-1.
func (cht *ConvexHullTrickLichao) Query(x int) (res, id int) {
	res, id = cht.query(cht.root, cht.lower, cht.upper, x)
	if !cht.isMin {
		res = -res
	}
	return
}

func (cht *ConvexHullTrickLichao) addLine(t *LichaoNode, x Line, l, r, xL, xR int) *LichaoNode {
	if t == nil {
		return &LichaoNode{line: x}
	}

	tL, tR := cht.getY(t.line, l), cht.getY(t.line, r)
	if tL <= xL && tR <= xR {
		return t
	} else if tL >= xL && tR >= xR {
		t.line = x
		return t
	} else {
		mid := (l + r) >> 1
		if mid == r {
			mid--
		}

		tM, xM := cht.getY(t.line, mid), cht.getY(x, mid)
		if tM > xM {
			t.line, x = x, t.line
			if xL >= tL {
				t.l = cht.addLine(t.l, x, l, mid, tL, tM)
			} else {
				t.r = cht.addLine(t.r, x, mid+1, r, tM+x.k, tR)
			}
		} else {
			if tL >= xL {
				t.l = cht.addLine(t.l, x, l, mid, xL, xM)
			} else {
				t.r = cht.addLine(t.r, x, mid+1, r, xM+x.k, xR)
			}
		}

		return t
	}
}

func (cht *ConvexHullTrickLichao) addSegment(t *LichaoNode, x Line, a, b, l, r, xL, xR int) *LichaoNode {
	if r < a || b < l {
		return t
	}
	if a <= l && r <= b {
		y := Line{x.k, x.b, x.id}
		return cht.addLine(t, y, l, r, xL, xR)
	}

	if t != nil {
		tL, tR := cht.getY(t.line, l), cht.getY(t.line, r)
		if tL <= xL && tR <= xR {
			return t
		}
	} else {
		t = &LichaoNode{line: Line{0, INF, -1}}
	}

	mid := (l + r) >> 1
	if mid == r {
		mid--
	}
	xM := cht.getY(x, mid)
	t.l = cht.addSegment(t.l, x, a, b, l, mid, xL, xM)
	t.r = cht.addSegment(t.r, x, a, b, mid+1, r, xM+x.k, xR)
	return t
}

func (cht *ConvexHullTrickLichao) query(t *LichaoNode, l, r, x int) (res, id int) {
	if t == nil {
		res, id = INF, -1
		return
	}
	if l == r {
		res, id = cht.getY(t.line, x), t.line.id
		return
	}

	mid := (l + r) >> 1
	if mid == r {
		mid--
	}

	res, id = cht.getY(t.line, x), t.line.id
	if x <= mid {
		cand, candId := cht.query(t.l, l, mid, x)
		if cand < res {
			res, id = cand, candId
		}
	} else {
		cand, candId := cht.query(t.r, mid+1, r, x)
		if cand < res {
			res, id = cand, candId
		}
	}
	return
}

func (cht *ConvexHullTrickLichao) getY(line Line, x int) int {
	return line.k*x + line.b
}

type ConvexHullTrickDeque struct {
	isMin bool
	dq    *Deque
}

func NewConvexHullTrickDeque(isMin bool) *ConvexHullTrickDeque {
	return &ConvexHullTrickDeque{
		isMin: isMin,
		dq:    &Deque{},
	}
}

// 追加一条直线.需要保证斜率k是单调递增或者是单调递减的.
func (cht *ConvexHullTrickDeque) AddLineMonotone(k, b, id int) {
	if !cht.isMin {
		k, b = -k, -b
	}

	line := Line{k, b, id}
	if cht.dq.Empty() {
		cht.dq.AppendLeft(line)
		return
	}

	if cht.dq.Front().k <= k {
		if cht.dq.Front().k == k {
			if cht.dq.Front().b <= b {
				return
			}
			cht.dq.PopLeft()
		}
		for cht.dq.Len() >= 2 && cht.check(line, cht.dq.Front(), cht.dq.At(1)) {
			cht.dq.PopLeft()
		}
		cht.dq.AppendLeft(line)
	} else {
		if cht.dq.Back().k == k {
			if cht.dq.Back().b <= b {
				return
			}
			cht.dq.Pop()
		}
		for cht.dq.Len() >= 2 && cht.check(cht.dq.At(cht.dq.Len()-2), cht.dq.Back(), line) {
			cht.dq.Pop()
		}
		cht.dq.Append(line)
	}
}

// O(logn) 查询 k*x + b 的最小(大)值以及对应的直线id.
// 如果不存在直线,返回的id为-1.
func (cht *ConvexHullTrickDeque) Query(x int) (res, lineId int) {
	if cht.dq.Empty() {
		res, lineId = INF, -1
		if !cht.isMin {
			res = -INF
		}
		return
	}

	left, right := -1, cht.dq.Len()-1
	for left+1 < right {
		mid := (left + right) >> 1
		a, _ := cht.getY(cht.dq.At(mid), x)
		b, _ := cht.getY(cht.dq.At(mid+1), x)
		if a >= b {
			left = mid
		} else {
			right = mid
		}
	}

	res, lineId = cht.getY(cht.dq.At(right), x)
	if !cht.isMin {
		res = -res
	}
	return
}

// O(1) 查询 k*x + b 的最小(大)值以及对应的直线id.
//
//	需要保证x是单调递增的.
//	如果不存在直线,返回的id为-1.
func (cht *ConvexHullTrickDeque) QueryMonotoneInc(x int) (res, lineId int) {
	if cht.dq.Empty() {
		res, lineId = INF, -1
		if !cht.isMin {
			res = -INF
		}
		return
	}

	for cht.dq.Len() >= 2 {
		a, _ := cht.getY(cht.dq.Front(), x)
		b, _ := cht.getY(cht.dq.At(1), x)
		if a < b {
			break
		}
		cht.dq.PopLeft()
	}

	res, lineId = cht.getY(cht.dq.Front(), x)
	if !cht.isMin {
		res = -res
	}
	return
}

// O(1) 查询 k*x + b 的最小(大)值以及对应的直线id.
//
//	需要保证x是单调递减的.
//	如果不存在直线,返回的id为-1.
func (cht *ConvexHullTrickDeque) QueryMonotoneDec(x int) (res, lineId int) {
	if cht.dq.Empty() {
		res, lineId = INF, -1
		if !cht.isMin {
			res = -INF
		}
		return
	}

	for cht.dq.Len() >= 2 {
		a, _ := cht.getY(cht.dq.Back(), x)
		b, _ := cht.getY(cht.dq.At(cht.dq.Len()-2), x)
		if a < b {
			break
		}
		cht.dq.Pop()
	}

	res, lineId = cht.getY(cht.dq.Back(), x)
	if !cht.isMin {
		res = -res
	}
	return
}

func (cht *ConvexHullTrickDeque) check(a, b, c Line) bool {
	if b.b == a.b || c.b == b.b {
		return sign(b.k-a.k)*sign(c.b-b.b) >= sign(c.k-b.k)*sign(b.b-a.b)
	}
	return (b.k-a.k)*sign(c.b-b.b)*abs(c.b-b.b) >= (c.k-b.k)*sign(b.b-a.b)*abs(b.b-a.b)
}

func (cht *ConvexHullTrickDeque) getY(line Line, x int) (int, int) {
	return line.k*x + line.b, line.id
}

func sign(x int) int {
	if x == 0 {
		return 0
	} else if x > 0 {
		return 1
	} else {
		return -1
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type E = Line
type Deque struct{ l, r []E }

func (q Deque) Empty() bool     { return len(q.l) == 0 && len(q.r) == 0 }
func (q Deque) Len() int        { return len(q.l) + len(q.r) }
func (q *Deque) AppendLeft(v E) { q.l = append(q.l, v) }
func (q *Deque) Append(v E)     { q.r = append(q.r, v) }
func (q *Deque) PopLeft() (v E) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v E) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque) Front() E {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque) Back() E {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque) At(i int) E {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}
