package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
)

// from https://atcoder.jp/users/ccppjsrb
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

// 特殊な入力形式に注意してください。 また、メモリ制限が通常より小さいことに注意してください。

// 頂点
// 1,2,…,N の
// N 頂点からなる無向グラフがあり、最初辺はありません。
// このグラフに対して、以下の
// Q 個のクエリを処理してください。

// 1
// u
// v
// タイプ
// 1 : 頂点
// u と頂点
// v との間に辺を追加する。
// 辺を追加する前の時点で、
// u と
// v は異なる連結成分に属する。(すなわち、グラフは常に森である。)

// 2
// u
// v
// タイプ
// 2 : 頂点
// u と頂点
// v の双方に隣接する頂点があるならその番号を答え、無ければ
// 0 と答える。
// グラフが常に森であることから、このクエリに対する解答は一意に定まることが示せる。

// 但し、上記のクエリは暗号化して与えられます。
// 本来のクエリは
// 3 つの整数
// A,B,C として定義され、これをもとに暗号化されたクエリが
// 3 つの整数
// a,b,c として与えられます。
// タイプ
// 2 のクエリのうち、先頭から
// k 個目のものに対する解答を
// X
// k
// ​
//   とします。 さらに、
// k=0 に対して
// X
// k
// ​
//  =0 と定義します。
// 与えられた
// a,b,c から以下の通りに
// A,B,C を復号してください。

// そのクエリより前に与えられたタイプ
// 2 のクエリの個数を
// l とする(そのクエリ自身は数えない)。このとき、以下の通りに復号せよ。
// A=1+(((a×(1+X
// l
// ​
//
//	))mod998244353)mod2)
//
// B=1+(((b×(1+X
// l
// ​
//
//	))mod998244353)modN)
//
// C=1+(((c×(1+X
// l
// ​
//
//	))mod998244353)modN)

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()
	const MOD int = 998244353

	n, q := io.NextInt(), io.NextInt()
	lct := NewLinkCutTree32(false)
	nodes := lct.Build(int32(n), func(i int32) E { return 1 })
	preRes := 0

	for i := 0; i < q; i++ {
		A, B, C := io.NextInt(), io.NextInt(), io.NextInt()
		A = 1 + (((A * (1 + preRes)) % MOD) % 2)
		B = 1 + (((B * (1 + preRes)) % MOD) % n)
		C = 1 + (((C * (1 + preRes)) % MOD) % n)
		B--
		C--
		if A == 1 {
			lct.LinkEdge(nodes[B], nodes[C])
		} else {
			if !lct.IsConnected(nodes[B], nodes[C]) {
				io.Println(0)
				preRes = 0
			} else {
				d := lct.QueryPath(nodes[B], nodes[C])
				if d != 3 {
					io.Println(0)
					preRes = 0
				} else {
					lct.Evert(nodes[B])
					preRes = int(lct.GetParent(nodes[C]).id + 1)
					io.Println(preRes)
				}
			}
		}
	}

}

type E = int32

func (*LinkCutTree32) rev(e E) E   { return e } // 区间反转
func (*LinkCutTree32) op(a, b E) E { return a + b }

type LinkCutTree32 struct {
	nodeId int32
	edges  map[[2]int32]struct{}
	check  bool
}

// check: AddEdge/RemoveEdge で辺の存在チェックを行うかどうか.
func NewLinkCutTree32(check bool) *LinkCutTree32 {
	return &LinkCutTree32{edges: make(map[[2]int32]struct{}), check: check}
}

func (lct *LinkCutTree32) Build(n int32, f func(i int32) E) []*treeNode {
	nodes := make([]*treeNode, n)
	for i := int32(0); i < n; i++ {
		nodes[i] = lct.Alloc(f(i))
	}
	return nodes
}

func (lct *LinkCutTree32) Alloc(e E) *treeNode {
	res := newTreeNode(e, lct.nodeId)
	lct.nodeId++
	return res
}

// t を根に変更する.
func (lct *LinkCutTree32) Evert(t *treeNode) {
	lct.expose(t)
	lct.toggle(t)
	lct.push(t)
}

// 存在していない辺 uv を新たに張る.
//
//	すでに存在している辺 uv に対しては何もしない.
func (lct *LinkCutTree32) LinkEdge(child, parent *treeNode) (ok bool) {
	if lct.check {
		if lct.IsConnected(child, parent) {
			return
		}
		id1, id2 := child.id, parent.id
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		lct.edges[[2]int32{id1, id2}] = struct{}{}
	}

	lct.Evert(child)
	lct.expose(parent)
	child.p = parent
	parent.r = child
	lct.update(parent)
	return true
}

// 存在している辺を切り離す.
//
//	存在していない辺に対しては何もしない.
func (lct *LinkCutTree32) CutEdge(u, v *treeNode) (ok bool) {
	if lct.check {
		id1, id2 := u.id, v.id
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		tuple := [2]int32{id1, id2}
		if _, has := lct.edges[tuple]; !has {
			return
		}
		delete(lct.edges, tuple)
	}

	lct.Evert(u)
	lct.expose(v)
	parent := v.l
	v.l = nil
	lct.update(v)
	parent.p = nil
	return true
}

// u と v の lca を返す.
//
//	u と v が異なる連結成分なら nullptr を返す.
//	!上記の操作は根を勝手に変えるため, 事前に Evert する必要があるかも.
func (lct *LinkCutTree32) LCA(u, v *treeNode) *treeNode {
	if !lct.IsConnected(u, v) {
		return nil
	}
	lct.expose(u)
	return lct.expose(v)
}

func (lct *LinkCutTree32) KthAncestor(x *treeNode, k int32) *treeNode {
	lct.expose(x)
	for x != nil {
		lct.push(x)
		if x.r != nil && x.r.sz > k {
			x = x.r
		} else {
			if x.r != nil {
				k -= x.r.sz
			}
			if k == 0 {
				return x
			}
			k--
			x = x.l
		}
	}
	return nil
}

func (lct *LinkCutTree32) GetParent(t *treeNode) *treeNode {
	lct.expose(t)
	p := t.l
	if p == nil {
		return nil
	}
	for {
		lct.push(p)
		if p.r == nil {
			return p
		}
		p = p.r
	}
}

func (lct *LinkCutTree32) Jump(from, to *treeNode, k int32) *treeNode {
	lct.Evert(to)
	lct.expose(from)
	for {
		lct.push(from)
		rs := int32(0)
		if from.r != nil {
			rs = from.r.sz
		}
		if k < rs {
			from = from.r
			continue
		}
		if k == rs {
			break
		}
		k -= rs + 1
		from = from.l
	}
	lct.splay(from)
	return from
}

// u から根までのパス上の頂点の値を二項演算でまとめた結果を返す.
func (lct *LinkCutTree32) QueryToRoot(u *treeNode) E {
	lct.expose(u)
	return u.sum
}

// u から v までのパス上の頂点の値を二項演算でまとめた結果を返す.
func (lct *LinkCutTree32) QueryPath(u, v *treeNode) E {
	lct.Evert(u)
	return lct.QueryToRoot(v)
}

// t の値を v に変更する.
func (lct *LinkCutTree32) Set(t *treeNode, v E) {
	lct.expose(t)
	t.key = v
	lct.update(t)
}

// t の値を返す.
func (lct *LinkCutTree32) Get(t *treeNode) E {
	return t.key
}

// u と v が同じ連結成分に属する場合は true, そうでなければ false を返す.
func (lct *LinkCutTree32) IsConnected(u, v *treeNode) bool {
	return u == v || lct.GetRoot(u) == lct.GetRoot(v)
}

// t の根を返す.
func (lct *LinkCutTree32) GetRoot(t *treeNode) *treeNode {
	lct.expose(t)
	for t.l != nil {
		lct.push(t)
		t = t.l
	}
	return t
}

func (lct *LinkCutTree32) expose(t *treeNode) *treeNode {
	var rp *treeNode
	for cur := t; cur != nil; cur = cur.p {
		lct.splay(cur)
		cur.r = rp
		lct.update(cur)
		rp = cur
	}
	lct.splay(t)
	return rp
}

func (lct *LinkCutTree32) update(t *treeNode) *treeNode {
	t.sz = 1
	t.sum = t.key
	if t.l != nil {
		t.sz += t.l.sz
		t.sum = lct.op(t.l.sum, t.sum)
	}
	if t.r != nil {
		t.sz += t.r.sz
		t.sum = lct.op(t.sum, t.r.sum)
	}
	return t
}

func (lct *LinkCutTree32) rotr(t *treeNode) {
	x := t.p
	y := x.p
	x.l = t.r
	if t.r != nil {
		t.r.p = x
	}
	t.r = x
	x.p = t
	lct.update(x)
	lct.update(t)
	t.p = y
	if y != nil {
		if y.l == x {
			y.l = t
		}
		if y.r == x {
			y.r = t
		}
		lct.update(y)
	}
}

func (lct *LinkCutTree32) rotl(t *treeNode) {
	x := t.p
	y := x.p
	x.r = t.l
	if t.l != nil {
		t.l.p = x
	}
	t.l = x
	x.p = t
	lct.update(x)
	lct.update(t)
	t.p = y
	if y != nil {
		if y.l == x {
			y.l = t
		}
		if y.r == x {
			y.r = t
		}
		lct.update(y)
	}
}

func (lct *LinkCutTree32) toggle(t *treeNode) {
	t.l, t.r = t.r, t.l
	t.sum = lct.rev(t.sum)
	t.rev = !t.rev
}

func (lct *LinkCutTree32) push(t *treeNode) {
	if t.rev {
		if t.l != nil {
			lct.toggle(t.l)
		}
		if t.r != nil {
			lct.toggle(t.r)
		}
		t.rev = false
	}
}

func (lct *LinkCutTree32) splay(t *treeNode) {
	lct.push(t)
	for !t.IsRoot() {
		q := t.p
		if q.IsRoot() {
			lct.push(q)
			lct.push(t)
			if q.l == t {
				lct.rotr(t)
			} else {
				lct.rotl(t)
			}
		} else {
			r := q.p
			lct.push(r)
			lct.push(q)
			lct.push(t)
			if r.l == q {
				if q.l == t {
					lct.rotr(q)
					lct.rotr(t)
				} else {
					lct.rotl(t)
					lct.rotr(t)
				}
			} else {
				if q.r == t {
					lct.rotl(q)
					lct.rotl(t)
				} else {
					lct.rotr(t)
					lct.rotl(t)
				}
			}
		}
	}
}

type treeNode struct {
	l, r, p  *treeNode
	key, sum E
	rev      bool
	sz       int32
	id       int32
}

func newTreeNode(v E, id int32) *treeNode {
	return &treeNode{key: v, sum: v, sz: 1, id: id}
}

func (n *treeNode) IsRoot() bool {
	return n.p == nil || (n.p.l != n && n.p.r != n)
}

func (n *treeNode) String() string {
	return fmt.Sprintf("key: %v, sum: %v, sz: %v, rev: %v", n.key, n.sum, n.sz, n.rev)
}
