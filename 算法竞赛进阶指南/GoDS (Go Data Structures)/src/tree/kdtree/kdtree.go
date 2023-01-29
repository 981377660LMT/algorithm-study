// kdtree 二维矩形计数 (如果可以离线的话,可以用二维树状数组RectangleBIT)
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	// !你有一个N×N(N<=5e5)的棋盘，每个格子内有一个整数，初始时的时候全部为 0，现在需要维护两种操作：
	// 1 x y v：将(x,y)位置的数加上v
	// 2 x1 y1 x2 y2：询问(x1,y1)到(x2,y2)的矩形区域内所有数的和
	// 3 终止程序
	// 输入文件第一行一个正整数 N。
	// 接下来每行一个操作。每条命令除第一个数字之外，均要异或上一次输出的答案last_ans，初始时 last_ans =0。
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	kdTree := NewKdTree()

	var n, last, kind, x1, y1, x2, y2, val int
	fmt.Fscan(in, &n)
	for {
		fmt.Fscan(in, &kind)
		if kind == 3 {
			break
		}
		if kind == 1 {
			fmt.Fscan(in, &x1, &y1, &val)
			kdTree.Add([2]int{x1 ^ last, y1 ^ last}, val^last)
		} else {
			fmt.Fscan(in, &x1, &y1, &x2, &y2)
			last = kdTree.Query(x1^last, y1^last, x2^last, y2^last)
			fmt.Fprintln(out, last)
		}
	}
}

// https://github.dev/EndlessCheng/codeforces-go/blob/a0733fa7a046673ff42a058b0dca7852646fbf3b/copypasta/kd_tree.go#L7s
type kdTree struct {
	root *kdNode
}

func NewKdTree() *kdTree {
	rand.Seed(time.Now().UnixNano())
	return &kdTree{}
}

func (t *kdTree) Add(p [2]int, val int) { t.root = t.root.add(p, val, 0) }

func (t *kdTree) Query(x1, y1, x2, y2 int) int { return t.root.query(x1, y1, x2, y2) }

type kdNode struct {
	lr          [2]*kdNode
	p, mi, mx   [2]int // 0 为 x，1 为 y
	sz, val, sm int
}

func (kdNode) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (kdNode) max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (o *kdNode) size() int {
	if o != nil {
		return o.sz
	}
	return 0
}

func (o *kdNode) sum() int {
	if o != nil {
		return o.sm
	}
	return 0
}

func (o *kdNode) maintain() {
	o.sz = o.lr[0].size() + o.lr[1].size() + 1
	o.sm = o.lr[0].sum() + o.lr[1].sum() + o.val
	for i := 0; i < 2; i++ {
		o.mi[i] = o.p[i]
		o.mx[i] = o.p[i]
		for _, ch := range o.lr {
			if ch != nil {
				o.mi[i] = o.min(o.mi[i], ch.mi[i])
				o.mx[i] = o.max(o.mx[i], ch.mx[i])
			}
		}
	}
}

func (o *kdNode) nodes() []*kdNode {
	nodes := make([]*kdNode, 0, o.size())
	var f func(*kdNode)
	f = func(o *kdNode) {
		if o != nil {
			nodes = append(nodes, o)
			f(o.lr[0])
			f(o.lr[1])
		}
	}
	f(o)
	rand.Shuffle(len(nodes), func(i, j int) { nodes[i], nodes[j] = nodes[j], nodes[i] })
	return nodes
}

func divideKDT(a []*kdNode, k, dim int) {
	for l, r := 0, len(a)-1; l < r; {
		v := a[l].p[dim]
		i, j := l, r+1
		for {
			for i++; i < r && a[i].p[dim] < v; i++ {
			}
			for j--; j > l && a[j].p[dim] > v; j-- {
			}
			if i >= j {
				break
			}
			a[i], a[j] = a[j], a[i]
		}
		a[l], a[j] = a[j], a[l]
		if j == k {
			break
		} else if j < k {
			l = j + 1
		} else {
			r = j - 1
		}
	}
}

// 另一种实现是选择的维度要满足其内部点的分布的差异度最大，见 https://oi-wiki.org/ds/kdt/
func buildKDT(nodes []*kdNode, dim int) *kdNode {
	if len(nodes) == 0 {
		return nil
	}
	m := len(nodes) / 2
	divideKDT(nodes, m, dim)
	o := nodes[m]
	o.lr[0] = buildKDT(nodes[:m], dim^1)
	o.lr[1] = buildKDT(nodes[m+1:], dim^1)
	o.maintain()
	return o
}

func (o *kdNode) rebuild(dim int) *kdNode { return buildKDT(o.nodes(), dim) }

func (o *kdNode) add(p [2]int, val, dim int) *kdNode {
	if o == nil {
		o = &kdNode{p: p, val: val}
		o.maintain()
		return o
	}
	if p[dim] < o.p[dim] {
		o.lr[0] = o.lr[0].add(p, val, dim^1)
	} else {
		o.lr[1] = o.lr[1].add(p, val, dim^1)
	}
	o.maintain()
	if sz := o.size() * 3; o.lr[0].size()*4 > sz || o.lr[1].size()*4 > sz { // alpha=3/4
		return o.rebuild(dim)
	}
	return o
}

func (o *kdNode) query(x1, y1, x2, y2 int) (res int) {
	if o == nil || outRect(x1, y1, x2, y2, o.mi[0], o.mi[1], o.mx[0], o.mx[1]) {
		return
	}
	if inRect(x1, y1, x2, y2, o.mi[0], o.mi[1], o.mx[0], o.mx[1]) {
		return o.sm
	}
	if inRect(x1, y1, x2, y2, o.p[0], o.p[1], o.p[0], o.p[1]) { // 根在询问矩形内
		res = o.val
	}
	res += o.lr[0].query(x1, y1, x2, y2) + o.lr[1].query(x1, y1, x2, y2)
	return
}

// 矩形 X-Y 在矩形 x-y 内
func inRect(x1, y1, x2, y2, X1, Y1, X2, Y2 int) bool {
	return x1 <= X1 && X2 <= x2 && y1 <= Y1 && Y2 <= y2
}

// 矩形 X-Y 在矩形 x-y 外
func outRect(x1, y1, x2, y2, X1, Y1, X2, Y2 int) bool {
	return X2 < x1 || X1 > x2 || Y2 < y1 || Y1 > y2
}
