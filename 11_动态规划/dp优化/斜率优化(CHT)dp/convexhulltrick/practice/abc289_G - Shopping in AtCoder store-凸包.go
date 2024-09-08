// G - Shopping in AtCoder store
// https://atcoder.jp/contests/abc289/tasks/abc289_g
// 商品定价
// n个顾客，每个人有一个购买的欲望wanti,
// m件物品，每一件物品有一个价值pricei,
// !每一个顾客会买商品当且仅当 wanti + pricei >= 定价
// !现在要求对每一个商品定价，求出它的最大销售值（数量*定价）

// 1<=n,m<=2e5
// https://blog.csdn.net/sophilex/article/details/129014785

// 1. 每一个商品的定价一定是从 wanti+pricei 中选出 (否则可以抬高定价)
//    !对于商品i，我们的销售额就是max{(j+1)*(wj+pi)}(0<=j<n)
// 2. 我们令横坐标为pi，纵坐标为对应的价值，不同j的选择对应不同的总价值。
//    显然我们最后是要找一个凸包，最后的答案就是横坐标对应的凸包上的点的纵坐标了

// 实现:
// 对商品原价 x(横坐标x),直线 linej 表示有 j+1 个顾客购买了这个商品,商品定价确定为 wantj+x
// !那么直线j对应的总收益为 yj = (j+1)*(wantj+x)
// 维护这k条直线即可

// 建模总结:
// 用 CHT 直接调模板，这种题关键是构造 j 条直线
// 1.什么作为横坐标,什么作为纵坐标 => 要查询的自变量(每个商品原价为x时)作为横坐标,答案作为纵坐标(此时的答案就是总收益)
// 2.什么作为(j条)直线 => 每条直线对应不同购买人数j，购买人数j唯一确定了一条直线
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
	sort.Slice(wants, func(i, j int) bool { // 倒序
		return wants[i] > wants[j]
	})

	// yj = (j+1)*(wantj+x) =  (j+1)*x + (j+1)*wantj
	cht := NewConvexHullTrickLichao(false, 0, 1e9+10)
	for j := 0; j < len(wants); j++ {
		cht.AddLine(j+1, wants[j]*(j+1), -1)
	}

	res := make([]int, len(prices))
	for i := 0; i < len(prices); i++ {
		res[i], _ = cht.Query(prices[i])
	}
	return res
}

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
	upper++
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
