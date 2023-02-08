package main

// !1.青蛙跳
// https://atcoder.jp/contests/dp/tasks/dp_z
// dp[j] = min(dp[j], dp[i] + c + (h[j] - h[i])^2) (0<=i<j<=n-1)
// !变形: dp[j]= min(-2*h[i]*h[j] + (h[i]^2+dp[i]) + h[j]^2 + c)
func frog3(h []int, c int) int {
	n := len(h)

	// 这里Query的自变量是h[j],题目里1<=hi<=1e6
	// lower, upper := 1, int(1e6)
	// cht := NewConvexHullTrickLichao(true, lower, upper)

	cht := NewConvexHullTrickDeque(true)
	dp := make([]int, n)
	dp[0] = 0
	cht.AddLine(-2*h[0], h[0]*h[0]+dp[0], 0)
	for j := 1; j < n; j++ {
		best, _ := cht.Query(h[j])
		dp[j] = best + h[j]*h[j] + c
		cht.AddLine(-2*h[j], h[j]*h[j]+dp[j], j)
	}

	return dp[n-1]
}

// 2.特别行动队
// https://www.acwing.com/problem/content/336/
// 记p为前缀和
// dp[j]=max(dp[j],dp[i]+a*(p[j]-p[i])^2+b*(p[j]-p[i])+c) (1<=i<j<=n)
// !变形: dp[j]=max(-2*a*pi*pj + (a*pi^2+dpi-b*pi) + a*pj^2+b*pj+c)
func 特别行动队(h []int, a, b, c int) int {
	n := len(h)
	preSum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		preSum[i] = preSum[i-1] + h[i-1]
	}

	cht := NewConvexHullTrickDeque(false) // 注意这里是最大值
	dp := make([]int, n+1)
	dp[0] = 0
	cht.AddLine(-2*a*preSum[0], a*preSum[0]*preSum[0]+dp[0]-b*preSum[0], 0)

	for j := 1; j < n+1; j++ {
		best, _ := cht.Query(preSum[j])
		dp[j] = best + a*preSum[j]*preSum[j] + b*preSum[j] + c
		cht.AddLine(-2*a*preSum[j], a*preSum[j]*preSum[j]+dp[j]-b*preSum[j], j)
	}

	return dp[n]
}

// 3.打印文章
// https://www.acwing.com/problem/content/1096/
// 记p为前缀和
// dp[0]=0,dp[j]=min(dp[j],dp[i]+(p[j]-p[i])^2+c) (1<=i<j<=n)
// !变形: dp[j]=min(-2*pi*pj + (pi^2+dpi) + pj^2+c)
func 打印文章(h []int, c int) int {
	n := len(h)
	preSum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		preSum[i] = preSum[i-1] + h[i-1]
	}

	cht := NewConvexHullTrickDeque(true)
	dp := make([]int, n+1)
	dp[0] = 0
	cht.AddLine(-2*preSum[0], preSum[0]*preSum[0]+dp[0], 0)

	for j := 1; j < n+1; j++ {
		best, _ := cht.Query(preSum[j])
		dp[j] = best + preSum[j]*preSum[j] + c
		cht.AddLine(-2*preSum[j], preSum[j]*preSum[j]+dp[j], j)
	}

	return dp[n]
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
func (cht *ConvexHullTrickDeque) AddLine(k, b, id int) {
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
//  需要保证x是单调递增的.
//  如果不存在直线,返回的id为-1.
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
//  需要保证x是单调递减的.
//  如果不存在直线,返回的id为-1.
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

//
//
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
