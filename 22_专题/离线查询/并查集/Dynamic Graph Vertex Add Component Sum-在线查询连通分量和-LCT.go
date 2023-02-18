// https://nyaannyaan.github.io/library/graph/offline-dynamic-connectivity.hpp

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int = 1e18

func main() {
	// https://judge.yosupo.jp/problem/dynamic_graph_vertex_add_component_sum
	// 0 u v 连接u v (保证u v不连接)
	// 1 u v 断开u v  (保证u v连接)
	// 2 u x 将u的值加上x
	// 3 u 输出u所在连通块的值

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}
	queries := make([][]int, 0, q)
	for i := 0; i < q; i++ {
		var op, u, v, add int
		fmt.Fscan(in, &op)
		if op == 0 || op == 1 {
			fmt.Fscan(in, &u, &v)
			queries = append(queries, []int{op, u, v})
		} else if op == 2 {
			fmt.Fscan(in, &u, &add)
			queries = append(queries, []int{op, u, add})
		} else {
			fmt.Fscan(in, &u)
			queries = append(queries, []int{op, u})
		}
	}

	odc := NewOffLineDynamicConnectivity(n, q)
	for i, q := range queries {
		if q[0] == 0 {
			odc.AddEdge(i, q[1], q[2])
		} else if q[0] == 1 {
			odc.RemoveEdge(i, q[1], q[2])
		}
	}
	odc.Build()

	lct := NewLinkCutTreeSubTree(false) // without edge check
	nodes := lct.Build(values)
	res := make([]int, q)
	for i := range res {
		res[i] = -INF
	}
	add := func(u, v int) { lct.LinkEdge(nodes[u], nodes[v]) }
	remove := func(u, v int) { lct.CutEdge(nodes[u], nodes[v]) }
	query := func(qi int) {
		if queries[qi][0] == 2 { // 2 u x 将u的值加上x
			u, add := queries[qi][1], queries[qi][2]
			lct.Set(nodes[u], lct.Get(nodes[u])+add)
		} else if queries[qi][0] == 3 { // 3 u 输出u所在连通块的值
			u := queries[qi][1]
			lct.Evert(nodes[u])
			res[qi] = lct.QuerySubTree(nodes[u])
		}
	}

	odc.Run(add, remove, query)
	for _, v := range res {
		if v != -INF {
			fmt.Fprintln(out, v)
		}
	}

}

type pair = struct{ first, second int }
type OffLineDynamicConnectivity struct {
	n, q, segsz     int
	uf              *rollbackUnionFind
	seg, qadd, qdel [][]pair
	cnt             map[pair]*pair // map的value是不可寻址的,需要使用指针类型代替
}

// Init with n nodes and q queries.
func NewOffLineDynamicConnectivity(n, q int) *OffLineDynamicConnectivity {
	res := &OffLineDynamicConnectivity{
		n:     n,
		q:     q,
		segsz: 1,
		uf:    newRollBackUnionFind(n),
		qadd:  make([][]pair, q),
		qdel:  make([][]pair, q),
		cnt:   make(map[pair]*pair),
	}
	for res.segsz < q {
		res.segsz *= 2
	}
	res.seg = make([][]pair, res.segsz*2)
	return res
}

func (odc *OffLineDynamicConnectivity) Build() {
	for i := 0; i < odc.q; i++ {
		for _, e := range odc.qadd[i] {
			if _, ok := odc.cnt[e]; !ok {
				odc.cnt[e] = &pair{0, 0}
			}
			dat := odc.cnt[e]
			if dat.second++; dat.second == 1 {
				dat.first = i
			}
		}

		for _, e := range odc.qdel[i] {
			if _, ok := odc.cnt[e]; !ok {
				odc.cnt[e] = &pair{0, 0}
			}
			dat := odc.cnt[e]
			if dat.second--; dat.second == 0 {
				odc.segment(e, dat.first, i)
			}
		}
	}

	for e, dat := range odc.cnt {
		if dat.second != 0 {
			odc.segment(e, dat.first, odc.q)
		}
	}
}

// Add an edge u<->v at time t.
func (odc *OffLineDynamicConnectivity) AddEdge(t, u, v int) {
	if u > v {
		u, v = v, u
	}
	odc.qadd[t] = append(odc.qadd[t], pair{u, v})
}

// Remove an edge u<->v at time t.
func (odc *OffLineDynamicConnectivity) RemoveEdge(t, u, v int) {
	if u > v {
		u, v = v, u
	}
	odc.qdel[t] = append(odc.qdel[t], pair{u, v})
}

func (odc *OffLineDynamicConnectivity) Run(add, remove func(u, v int), query func(qi int)) {
	odc.dfs(add, remove, query, 1, 0, odc.segsz)
}

func (odc *OffLineDynamicConnectivity) dfs(add, remove func(u, v int), query func(qi int), id, l, r int) {
	if odc.q <= l {
		return
	}
	state := odc.uf.GetState()
	var es []pair
	for _, e := range odc.seg[id] {
		u, v := e.first, e.second
		if !odc.uf.IsConnected(u, v) {
			odc.uf.Union(u, v)
			add(u, v)
			es = append(es, e)
		}
	}

	if l+1 == r {
		query(l)
	} else {
		odc.dfs(add, remove, query, id<<1, l, (l+r)>>1)
		odc.dfs(add, remove, query, id<<1|1, (l+r)>>1, r)
	}

	for _, e := range es {
		remove(e.first, e.second)
	}
	odc.uf.Rollback(state)
}

func (odc *OffLineDynamicConnectivity) segment(e pair, l, r int) {
	left, right := l+odc.segsz, r+odc.segsz
	for left < right {
		if left&1 == 1 {
			odc.seg[left] = append(odc.seg[left], e)
			left++
		}
		if right&1 == 1 {
			right--
			odc.seg[right] = append(odc.seg[right], e)
		}
		left >>= 1
		right >>= 1
	}
}

func newRollBackUnionFind(n int) *rollbackUnionFind {
	data := make([]int, n)
	for i := range data {
		data[i] = -1
	}
	return &rollbackUnionFind{data: data}
}

type rollbackUnionFind struct {
	innerSnap int
	data      []int
	history   []struct{ a, b int }
}

// 撤销上一次合并操作.
func (uf *rollbackUnionFind) Undo() bool {
	if len(uf.history) == 0 {
		return false
	}
	uf.data[uf.history[len(uf.history)-1].a] = uf.history[len(uf.history)-1].b
	uf.history = uf.history[:len(uf.history)-1]
	uf.data[uf.history[len(uf.history)-1].a] = uf.history[len(uf.history)-1].b
	uf.history = uf.history[:len(uf.history)-1]
	return true
}

// 回滚到指定的状态.
//  state 为 -1 表示回滚到上一次 `SnapShot` 时保存的状态.
//  其他值表示回滚到合并(Union) `state` 次后的状态.
func (uf *rollbackUnionFind) Rollback(state int) bool {
	if state == -1 {
		state = uf.innerSnap
	}
	state <<= 1
	if state < 0 || state > len(uf.history) {
		return false
	}
	for state < len(uf.history) {
		uf.Undo()
	}
	return true
}

// 获取当前合并(Union)被调用的次数.
func (uf *rollbackUnionFind) GetState() int {
	return len(uf.history) >> 1
}

// 保存并查集当前的状态.
func (uf *rollbackUnionFind) Snapshot() {
	uf.innerSnap = len(uf.history) >> 1
}

func (uf *rollbackUnionFind) Union(x, y int) bool {
	x, y = uf.Find(x), uf.Find(y)
	uf.history = append(uf.history, struct{ a, b int }{x, uf.data[x]})
	uf.history = append(uf.history, struct{ a, b int }{y, uf.data[y]})
	if x == y {
		return false
	}
	if uf.data[x] > uf.data[y] {
		x, y = y, x
	}
	uf.data[x] += uf.data[y]
	uf.data[y] = x
	return true
}

func (uf *rollbackUnionFind) Find(x int) int {
	cur := x
	for uf.data[cur] >= 0 {
		cur = uf.data[cur]
	}
	return cur
}

func (uf *rollbackUnionFind) IsConnected(x, y int) bool { return uf.Find(x) == uf.Find(y) }

func (uf *rollbackUnionFind) GetSize(x int) int { return -uf.data[uf.Find(x)] }

func (uf *rollbackUnionFind) GetGroups() [][]int {
	mp := make(map[int][]int)
	for i := range uf.data {
		mp[uf.Find(i)] = append(mp[uf.Find(i)], i)
	}
	var res [][]int
	for _, g := range mp {
		res = append(res, g)
	}
	return res
}

func (uf *rollbackUnionFind) String() string {
	groups := uf.GetGroups()
	sort.Slice(groups, func(i, j int) bool { return groups[i][0] < groups[j][0] })
	res := []string{}
	for _, g := range groups {
		res = append(res, fmt.Sprintf("%v", g))
	}
	return fmt.Sprintf("state = %d, groups = %v", uf.GetState(), res)
}

type E = int // 子树和

func (*TreeNode) e() E                { return 0 }
func (*TreeNode) op(this, other E) E  { return this + other }
func (*TreeNode) inv(this, other E) E { return this - other }

type LinkCutTreeSubTree struct {
	nodeId int
	edges  map[struct{ u, v int }]struct{}
	check  bool
}

// check: AddEdge/RemoveEdge で辺の存在チェックを行うかどうか.
func NewLinkCutTreeSubTree(check bool) *LinkCutTreeSubTree {
	return &LinkCutTreeSubTree{edges: make(map[struct{ u, v int }]struct{}), check: check}
}

// 各要素の値を vs[i] としたノードを生成し, その配列を返す.
func (lct *LinkCutTreeSubTree) Build(vs []E) []*TreeNode {
	nodes := make([]*TreeNode, len(vs))
	for i, v := range vs {
		nodes[i] = lct.Alloc(v)
	}
	return nodes
}

// 要素の値を v としたノードを生成する.
func (lct *LinkCutTreeSubTree) Alloc(key E) *TreeNode {
	res := newTreeNode(key, lct.nodeId)
	lct.nodeId++
	lct.update(res)
	return res
}

// t を根に変更する.
func (lct *LinkCutTreeSubTree) Evert(t *TreeNode) {
	lct.expose(t)
	lct.toggle(t)
	lct.push(t)
}

func (lct *LinkCutTreeSubTree) LinkEdge(child, parent *TreeNode) (ok bool) {
	if lct.check {
		if lct.IsConnected(child, parent) {
			return
		}
		id1, id2 := child.id, parent.id
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		tuple := struct{ u, v int }{id1, id2}
		lct.edges[tuple] = struct{}{}
	}

	lct.Evert(child)
	lct.expose(parent)
	child.p = parent
	parent.r = child
	lct.update(parent)
	return true
}

func (lct *LinkCutTreeSubTree) CutEdge(u, v *TreeNode) (ok bool) {
	if lct.check {
		id1, id2 := u.id, v.id
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		tuple := struct{ u, v int }{id1, id2}
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
//  u と v が異なる連結成分なら nullptr を返す.
//  !上記の操作は根を勝手に変えるため, 事前に Evert する必要があるかも.
func (lct *LinkCutTreeSubTree) QueryLCA(u, v *TreeNode) *TreeNode {
	if !lct.IsConnected(u, v) {
		return nil
	}
	lct.expose(u)
	return lct.expose(v)
}

func (lct *LinkCutTreeSubTree) QueryKthAncestor(x *TreeNode, k int) *TreeNode {
	lct.expose(x)
	for x != nil {
		lct.push(x)
		if x.r != nil && x.r.cnt > k {
			x = x.r
		} else {
			if x.r != nil {
				k -= x.r.cnt
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

// t を根とする部分木の要素の値の和を返す.
//  !Evert を忘れない！
func (lct *LinkCutTreeSubTree) QuerySubTree(t *TreeNode) E {
	lct.expose(t)
	return t.op(t.key, t.sub)
}

// t の値を v に変更する.
func (lct *LinkCutTreeSubTree) Set(t *TreeNode, key E) *TreeNode {
	lct.expose(t)
	t.key = key
	lct.update(t)
	return t
}

// t の値を返す.
func (lct *LinkCutTreeSubTree) Get(t *TreeNode) E {
	return t.key
}

// u と v が同じ連結成分に属する場合は true, そうでなければ false を返す.
func (lct *LinkCutTreeSubTree) IsConnected(u, v *TreeNode) bool {
	return u == v || lct.GetRoot(u) == lct.GetRoot(v)
}

func (lct *LinkCutTreeSubTree) expose(t *TreeNode) *TreeNode {
	rp := (*TreeNode)(nil)
	for cur := t; cur != nil; cur = cur.p {
		lct.splay(cur)
		if cur.r != nil {
			cur.Add(cur.r)
		}
		cur.r = rp
		if cur.r != nil {
			cur.Erase(cur.r)
		}
		lct.update(cur)
		rp = cur
	}
	lct.splay(t)
	return rp
}

func (lct *LinkCutTreeSubTree) update(t *TreeNode) {
	t.cnt = 1
	if t.l != nil {
		t.cnt += t.l.cnt
	}
	if t.r != nil {
		t.cnt += t.r.cnt
	}

	t.Merge(t.l, t.r)
}

func (lct *LinkCutTreeSubTree) rotr(t *TreeNode) {
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

func (lct *LinkCutTreeSubTree) rotl(t *TreeNode) {
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

func (lct *LinkCutTreeSubTree) toggle(t *TreeNode) {
	t.l, t.r = t.r, t.l
	t.rev = !t.rev
}

func (lct *LinkCutTreeSubTree) push(t *TreeNode) {
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

func (lct *LinkCutTreeSubTree) splay(t *TreeNode) {
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

func (lct *LinkCutTreeSubTree) GetRoot(t *TreeNode) *TreeNode {
	lct.expose(t)
	for t.l != nil {
		lct.push(t)
		t = t.l
	}
	return t
}

type TreeNode struct {
	key, sum, sub E
	rev           bool
	cnt           int
	id            int
	l, r, p       *TreeNode
}

func newTreeNode(key E, id int) *TreeNode {
	res := &TreeNode{key: key, sum: key, cnt: 1, id: id}
	res.sub = res.e()
	return res
}

func (n *TreeNode) IsRoot() bool {
	return n.p == nil || (n.p.l != n && n.p.r != n)
}

func (n *TreeNode) Add(other *TreeNode)   { n.sub = n.op(n.sub, other.sum) }
func (n *TreeNode) Erase(other *TreeNode) { n.sub = n.inv(n.sub, other.sum) }
func (n *TreeNode) Merge(n1, n2 *TreeNode) {
	var tmp1, tmp2 E
	if n1 != nil {
		tmp1 = n1.sum
	} else {
		tmp1 = n.e()
	}

	if n2 != nil {
		tmp2 = n2.sum
	} else {
		tmp2 = n.e()
	}

	n.sum = n.op(n.op(tmp1, n.key), n.op(n.sub, tmp2))
}
