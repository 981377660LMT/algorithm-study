// 路径修改
// 路径查询
// 子树查询
// MaxPath树上二分

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	demo()
	// yosupoVertexAddPathSum()
	// rangeQueryOnTree()
	// yuki1197()
}

func demo() {
	{
		//    0
		//   / \
		//  1   2
		//     / \
		//    3   4

		tree := NewTree32(5)
		tree.AddEdge(0, 1, 0)
		tree.AddEdge(0, 2, 0)
		tree.AddEdge(2, 3, 0)
		tree.AddEdge(2, 4, 0)
		tree.Build(0)

		S := NewTreeMonoid32Lazy(tree, true)
		S.Build(func(vidOrEid int32) E { return int(vidOrEid) })
		fmt.Println(S.QuerySubtree(0)) // 10
		fmt.Println(S.QuerySubtree(1)) // 1
		fmt.Println(S.QuerySubtree(2)) // 9
		fmt.Println(S.QuerySubtree(3)) // 3
		fmt.Println(S.QuerySubtree(4)) // 4
		fmt.Println(S.QueryPath(1, 3)) // 6
		S.Update(3, 10)
		fmt.Println(S.QuerySubtree(0)) // 20

		fmt.Println(S.MaxPath(1, 3, func(x E) bool { return x < 4 }))  // 0
		fmt.Println(S.MaxPath(1, 3, func(x E) bool { return x < -1 })) // 0

		fmt.Println(S.GetAll())
		S.UpdateOuttree(2, 10)
		fmt.Println(S.GetAll())
		fmt.Println(S.Get(0))
	}
}

// https://judge.yosupo.jp/problem/vertex_add_path_sum
func yosupoVertexAddPathSum() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	weights := make([]int, n)
	for i := 0; i < int(n); i++ {
		fmt.Fscan(in, &weights[i])
	}
	tree := NewTree32(n)
	for i := 1; i < int(n); i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		tree.AddEdge(u, v, 0)
	}
	tree.Build(0)

	S := NewTreeMonoid32Lazy(tree, false)
	S.Build(func(vidOrEid int32) E { return weights[vidOrEid] })
	for i := 0; i < int(q); i++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 0 {
			var v, x int32
			fmt.Fscan(in, &v, &x)
			S.Update(v, int(x))
		} else {
			var u, v int32
			fmt.Fscan(in, &u, &v)
			fmt.Fprintln(out, S.QueryPath(u, v))
		}
	}
}

// 路径修改，路径查询，边权.
// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=GRL_5_E
func rangeQueryOnTree() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	tree := NewTree32(n)
	for i := int32(0); i < n; i++ {
		var k int32
		fmt.Fscan(in, &k)
		for j := int32(0); j < k; j++ {
			var c int32
			fmt.Fscan(in, &c)
			tree.AddEdge(i, c, 0)
		}
	}
	tree.Build(0)

	S := NewTreeMonoid32Lazy(tree, true)
	S.Build(func(vidOrEid int32) E { return 0 })

	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 0 {
			var v, w int32
			fmt.Fscan(in, &v, &w)
			S.UpdatePath(0, v, int(w))
		} else {
			var v int32
			fmt.Fscan(in, &v)
			fmt.Fprintln(out, S.QueryPath(0, v))
		}
	}
}

// https://yukicoder.me/problems/no/1197
// No.1197 モンスターショー
// 树上移动距离之和(带修改版)
// 给定一个有 n 个节点的树，每个节点与一条边相连，边长为 1。
// 一开始，有 k 只史莱姆分别位于不同的节点上。现在，你需要支持以下操作：
// 1 i d: 将第 i 只史莱姆从其所在节点移到节点 d。
// !2 e :查询将所有史莱姆移到节点 e 时，史莱姆移动距离之和。
// 其中，史莱姆从一个节点到另一个节点的移动距离为其走过的边长之和。
// !固定0号根节点,统计每个史莱姆的深度之和,然后统计从根节点移动到各个点的距离之和.
func yuki1197() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k, q int32
	fmt.Fscan(in, &n, &k, &q)
	pos := make([]int32, k) // 史莱姆的初始位置
	for i := range pos {
		fmt.Fscan(in, &pos[i])
		pos[i]--
	}

	tree := NewTree32(n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u--
		v--
		tree.AddEdge(u, v, 1)
	}
	tree.Build(0)

	M := NewTreeMonoid32Lazy(tree, false)
	M.Build(func(vidOrEid int32) E { return 0 })
	for _, v := range pos {
		M.UpdatePath(0, v, 1)
	}

	depth := tree.Depth
	depSum := 0
	for _, v := range pos {
		depSum += int(depth[v])
	}

	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 { // 维护史莱姆的深度之和、根节点到各个结点的移动距离之和
			var slime, moveTo int32
			fmt.Fscan(in, &slime, &moveTo)
			slime, moveTo = slime-1, moveTo-1
			depSum -= int(depth[pos[slime]])
			M.UpdatePath(0, pos[slime], -1)
			pos[slime] = moveTo
			depSum += int(depth[pos[slime]])
			M.UpdatePath(0, pos[slime], 1)
		} else {
			var to int32
			fmt.Fscan(in, &to)
			to--
			res := int(depth[to])*int(k) + depSum // 所有史莱姆先走到根节点,再走到to的距离之和
			res -= 2 * M.QueryPath(0, to)         // 多算的距离
			res += int(2 * k)                     // 线段树里把0-0的移动距离统计为1,所以要减去2k
			fmt.Fprintln(out, res)
		}
	}
}

const commutative bool = true // !E是否可交换，即 op(a,b) == op(b,a)

type E = int
type Id = int

func e() E   { return 0 }
func id() Id { return 0 }
func op(left, right E) E {
	return left + right
}
func mapping(f Id, g E, size int) E {
	return f*size + g
}
func composition(f, g Id) Id {
	return f + g
}

type TreeMonoid32Lazy struct {
	edge      int32
	n         int32
	tree      *Tree32
	seg, segR *LazySegTree32
}

func NewTreeMonoid32Lazy(tree *Tree32, edge bool) *TreeMonoid32Lazy {
	var edgeValue int32
	if edge {
		edgeValue = 1
	}
	return &TreeMonoid32Lazy{edge: edgeValue, n: tree.n, tree: tree}
}

func (tag *TreeMonoid32Lazy) Build(f func(vidOrEid int32) E) {
	idToNode := tag.tree.IdToNode
	vToE := tag.tree.vToE
	if tag.edge == 0 {
		fv := func(i int32) E { return f(idToNode[i]) }
		tag.seg = NewLazySegTree32(tag.n, fv, e, id, op, mapping, composition)
		if !commutative {
			tag.segR = NewLazySegTree32(tag.n, fv, e, id, func(a, b E) E { return op(b, a) }, mapping, composition)
		}
	} else {
		fe := func(i int32) E {
			if i == 0 {
				return e()
			}
			return f(vToE[idToNode[i]])
		}
		tag.seg = NewLazySegTree32(tag.n, fe, e, id, op, mapping, composition)
		if !commutative {
			tag.segR = NewLazySegTree32(tag.n, fe, e, id, func(a, b E) E { return op(b, a) }, mapping, composition)
		}
	}
}

func (tag *TreeMonoid32Lazy) Set(i int32, x E) {
	if tag.edge != 0 {
		i = tag.tree.EToV(i)
	}
	i = tag.tree.Lid[i]
	tag.seg.Set(i, x)
	if !commutative {
		tag.segR.Set(i, x)
	}
}

func (tag *TreeMonoid32Lazy) Update(i int32, x E) {
	if tag.edge != 0 {
		i = tag.tree.EToV(i)
	}
	i = tag.tree.Lid[i]
	tag.seg.Multiply(i, x)
	if !commutative {
		tag.segR.Multiply(i, x)
	}
}

func (tag *TreeMonoid32Lazy) Get(i int32) E {
	if tag.edge != 0 {
		i = tag.tree.EToV(i)
	}
	return tag.seg.Get(tag.tree.Lid[i])
}

func (tag *TreeMonoid32Lazy) GetAll() []E {
	data := tag.seg.GetAll()
	tree := tag.tree
	if tag.edge == 0 {
		res := make([]E, tag.n)
		for v := int32(0); v < tag.n; v++ {
			res[v] = data[tree.Lid[v]]
		}
		return res
	} else {
		res := make([]E, tag.n-1)
		for i := int32(0); i < tag.n-1; i++ {
			res[i] = data[tree.Lid[tree.EToV(i)]]
		}
		return res
	}
}

func (tag *TreeMonoid32Lazy) QueryPath(from, to int32) E {
	pd := tag.tree.GetPathDecomposition(from, to, tag.edge)
	res := e()
	for i := 0; i < len(pd); i++ {
		res = op(res, tag.getProd(pd[i][0], pd[i][1]))
	}
	return res
}

func (tag *TreeMonoid32Lazy) QueryAll() E {
	if !commutative {
		panic("not implemented")
	}
	return tag.QueryAll()
}

func (tag *TreeMonoid32Lazy) QuerySubtree(u int32) E {
	if !commutative {
		panic("not implemented")
	}
	l, r := tag.tree.Lid[u], tag.tree.Rid[u]
	return tag.seg.Query(l+tag.edge, r)
}

func (tag *TreeMonoid32Lazy) UpdateSubtree(u int32, f Id) {
	l, r := tag.tree.Lid[u], tag.tree.Rid[u]
	tag.seg.Update(l+tag.edge, r, f)
	if !commutative {
		tag.segR.Update(l+tag.edge, r, f)
	}
}

// outtree: 除开u的子树外的其他节点
func (tag *TreeMonoid32Lazy) UpdateOuttree(u int32, f Id) {
	l, r := tag.tree.Lid[u], tag.tree.Rid[u]
	tag.seg.Update(tag.edge, l+tag.edge, f)
	tag.seg.Update(r, tag.n, f)
	if !commutative {
		tag.segR.Update(tag.edge, l+tag.edge, f)
		tag.segR.Update(r, tag.n, f)
	}
}

func (tag *TreeMonoid32Lazy) UpdatePath(u, v int32, f Id) {
	pd := tag.tree.GetPathDecomposition(u, v, tag.edge)
	for i := 0; i < len(pd); i++ {
		l, r := pd[i][0], pd[i][1]
		if l > r {
			l, r = r, l
		}
		tag.seg.Update(l, r+1, f)
		if !commutative {
			tag.segR.Update(l, r+1, f)
		}
	}
}

// 满足 check 为true的最远的节点.
// 如果不存在返回 -1.
func (tag *TreeMonoid32Lazy) MaxPath(from, to int32, check func(E) bool) int32 {
	if tag.edge != 0 {
		return tag.maxPathEdge(from, to, check)
	}
	if !check(tag.QueryPath(from, from)) {
		return -1
	}
	pd := tag.tree.GetPathDecomposition(from, to, tag.edge)
	val := e()
	idToNode := tag.tree.IdToNode
	for _, e := range pd {
		x := tag.getProd(e[0], e[1])
		if tmp := op(val, x); check(tmp) {
			val = tmp
			from = idToNode[e[1]]
			continue
		}
		checkTmp := func(x E) bool { return check(op(val, x)) }
		if e[0] <= e[1] {
			i := tag.seg.MaxRight(e[0], checkTmp)
			if i == e[0] {
				return from
			}
			return idToNode[i-1]
		} else {
			i := int32(0)
			if commutative {
				i = tag.seg.MinLeft(e[0]+1, checkTmp)
			} else {
				i = tag.segR.MinLeft(e[0]+1, checkTmp)
			}
			if i == e[0]+1 {
				return from
			}
			return idToNode[i]
		}
	}
	return to
}

func (tag *TreeMonoid32Lazy) maxPathEdge(from, to int32, check func(E) bool) int32 {
	if !check(e()) {
		return -1
	}
	lca := tag.tree.Lca(from, to)
	pd := tag.tree.GetPathDecomposition(from, lca, tag.edge)
	val := e()
	parent, idToNode := tag.tree.Parent, tag.tree.IdToNode
	for _, e := range pd {
		x := tag.getProd(e[0], e[1])
		if tmp := op(val, x); check(tmp) {
			val = tmp
			from = parent[idToNode[e[1]]]
			continue
		}
		checkTmp := func(x E) bool { return check(op(val, x)) }
		i := int32(0)
		if commutative {
			i = tag.seg.MinLeft(e[0]+1, checkTmp)
		} else {
			i = tag.segR.MinLeft(e[0]+1, checkTmp)
		}
		if i == e[0]+1 {
			return from
		}
		return parent[idToNode[i]]
	}
	pd = tag.tree.GetPathDecomposition(lca, to, tag.edge)
	for _, e := range pd {
		x := tag.getProd(e[0], e[1])
		if tmp := op(val, x); check(tmp) {
			val = tmp
			from = idToNode[e[1]]
			continue
		}
		checkTmp := func(x E) bool { return check(op(val, x)) }
		i := tag.seg.MaxRight(e[0], checkTmp)
		if i == e[0] {
			return from
		}
		return idToNode[i-1]
	}
	return to
}

func (tag *TreeMonoid32Lazy) getProd(a, b int32) E {
	if commutative {
		if a <= b {
			return tag.seg.Query(a, b+1)
		}
		return tag.seg.Query(b, a+1)
	} else {
		if a <= b {
			return tag.seg.Query(a, b+1)
		}
		return tag.segR.Query(b, a+1)
	}
}

// !template
type LazySegTree32 struct {
	n           int32
	size        int32
	log         int32
	data        []E
	lazy        []Id
	e           func() E
	id          func() Id
	op          func(E, E) E
	mapping     func(Id, E, int) E
	composition func(Id, Id) Id
}

func NewLazySegTree32(
	n int32, f func(int32) E,
	e func() E, id func() Id,
	op func(E, E) E, mapping func(Id, E, int) E, composition func(Id, Id) Id,
) *LazySegTree32 {
	tree := &LazySegTree32{e: e, id: id, op: op, mapping: mapping, composition: composition}
	tree.n = n
	tree.log = int32(bits.Len32(uint32(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := int32(0); i < n; i++ {
		tree.data[tree.size+i] = f(i)
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

// 查询切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTree32) Query(left, right int32) E {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return tree.e()
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	sml, smr := tree.e(), tree.e()
	for left < right {
		if left&1 != 0 {
			sml = tree.op(sml, tree.data[left])
			left++
		}
		if right&1 != 0 {
			right--
			smr = tree.op(tree.data[right], smr)
		}
		left >>= 1
		right >>= 1
	}
	return tree.op(sml, smr)
}
func (tree *LazySegTree32) QueryAll() E {
	return tree.data[1]
}
func (tree *LazySegTree32) GetAll() []E {
	for i := int32(1); i < tree.size; i++ {
		tree.pushDown(i)
	}
	res := make([]E, tree.n)
	copy(res, tree.data[tree.size:tree.size+tree.n])
	return res
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTree32) Update(left, right int32, f Id) {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	l2, r2 := left, right
	for left < right {
		if left&1 != 0 {
			tree.propagate(left, f)
			left++
		}
		if right&1 != 0 {
			right--
			tree.propagate(right, f)
		}
		left >>= 1
		right >>= 1
	}
	left = l2
	right = r2
	for i := int32(1); i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTree32) MinLeft(right int32, predicate func(data E) bool) int32 {
	if right == 0 {
		return 0
	}
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown((right - 1) >> i)
	}
	res := tree.e()
	for {
		right--
		for right > 1 && right&1 != 0 {
			right >>= 1
		}
		if !predicate(tree.op(tree.data[right], res)) {
			for right < tree.size {
				tree.pushDown(right)
				right = right<<1 | 1
				if tmp := tree.op(tree.data[right], res); predicate(tmp) {
					res = tmp
					right--
				}
			}
			return right + 1 - tree.size
		}
		res = tree.op(tree.data[right], res)
		if (right & -right) == right {
			break
		}
	}
	return 0
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTree32) MaxRight(left int32, predicate func(data E) bool) int32 {
	if left == tree.n {
		return tree.n
	}
	left += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(left >> i)
	}
	res := tree.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(tree.op(res, tree.data[left])) {
			for left < tree.size {
				tree.pushDown(left)
				left <<= 1
				if tmp := tree.op(res, tree.data[left]); predicate(tmp) {
					res = tmp
					left++
				}
			}
			return left - tree.size
		}
		res = tree.op(res, tree.data[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return tree.n
}

// 单点查询(不需要 pushUp/op 操作时使用)
func (tree *LazySegTree32) Get(index int32) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *LazySegTree32) Set(index int32, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := int32(1); i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *LazySegTree32) Multiply(index int32, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = tree.op(tree.data[index], e)
	for i := int32(1); i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *LazySegTree32) pushUp(root int32) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *LazySegTree32) pushDown(root int32) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}
func (tree *LazySegTree32) propagate(root int32, f Id) {
	size := 1 << (tree.log - int32((bits.Len32(uint32(root)) - 1)) /**topbit**/)
	tree.data[root] = tree.mapping(f, tree.data[root], size)
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}

func (tree *LazySegTree32) String() string {
	var sb []string
	sb = append(sb, "[")
	for i := int32(0); i < tree.n; i++ {
		if i != 0 {
			sb = append(sb, ", ")
		}
		sb = append(sb, fmt.Sprintf("%v", tree.Get(i)))
	}
	sb = append(sb, "]")
	return strings.Join(sb, "")
}

type neighbor = struct {
	to   int32
	eid  int32
	cost int
}

type Tree32 struct {
	Lid, Rid      []int32
	IdToNode      []int32
	Depth         []int32
	DepthWeighted []int
	Parent        []int32
	Head          []int32 // 重链头
	Tree          [][]neighbor
	Edges         [][2]int32
	vToE          []int32 // 节点v的父边的id
	n             int32
}

func NewTree32(n int32) *Tree32 {
	res := &Tree32{Tree: make([][]neighbor, n), Edges: make([][2]int32, 0, n-1), n: n}
	return res
}

func (t *Tree32) AddEdge(u, v int32, w int) {
	eid := int32(len(t.Edges))
	t.Tree[u] = append(t.Tree[u], neighbor{to: v, eid: eid, cost: w})
	t.Tree[v] = append(t.Tree[v], neighbor{to: u, eid: eid, cost: w})
	t.Edges = append(t.Edges, [2]int32{u, v})
}

func (t *Tree32) AddDirectedEdge(from, to int32, cost int) {
	eid := int32(len(t.Edges))
	t.Tree[from] = append(t.Tree[from], neighbor{to: to, eid: eid, cost: cost})
	t.Edges = append(t.Edges, [2]int32{from, to})
}

func (t *Tree32) Build(root int32) {
	if root != -1 && int32(len(t.Edges)) != t.n-1 {
		panic("edges count != n-1")
	}
	n := t.n
	t.Lid = make([]int32, n)
	t.Rid = make([]int32, n)
	t.IdToNode = make([]int32, n)
	t.Depth = make([]int32, n)
	t.DepthWeighted = make([]int, n)
	t.Parent = make([]int32, n)
	t.Head = make([]int32, n)
	t.vToE = make([]int32, n)
	for i := int32(0); i < n; i++ {
		t.Depth[i] = -1
		t.Head[i] = root
		t.vToE[i] = -1
	}
	if root != -1 {
		t._dfsSize(root, -1)
		time := int32(0)
		t._dfsHld(root, &time)
	} else {
		time := int32(0)
		for i := int32(0); i < n; i++ {
			if t.Depth[i] == -1 {
				t._dfsSize(i, -1)
				t._dfsHld(i, &time)
			}
		}
	}
}

// 从v开始沿着重链向下收集节点.
func (t *Tree32) HeavyPathAt(v int32) []int32 {
	path := []int32{v}
	for {
		a := path[len(path)-1]
		for _, e := range t.Tree[a] {
			if e.to != t.Parent[a] && t.Head[e.to] == v {
				path = append(path, e.to)
				break
			}
		}
		if path[len(path)-1] == a {
			break
		}
	}
	return path
}

// 返回重儿子，如果没有返回 -1.
func (t *Tree32) HeavyChild(v int32) int32 {
	k := t.Lid[v] + 1
	if k == t.n {
		return -1
	}
	w := t.IdToNode[k]
	if t.Parent[w] == v {
		return w
	}
	return -1
}

// 从v开始向上走k步.
func (t *Tree32) KthAncestor(v, k int32) int32 {
	if k > t.Depth[v] {
		return -1
	}
	for {
		u := t.Head[v]
		if t.Lid[v]-k >= t.Lid[u] {
			return t.IdToNode[t.Lid[v]-k]
		}
		k -= t.Lid[v] - t.Lid[u] + 1
		v = t.Parent[u]
	}
}

func (t *Tree32) Lca(u, v int32) int32 {
	for {
		if t.Lid[u] > t.Lid[v] {
			u, v = v, u
		}
		if t.Head[u] == t.Head[v] {
			return u
		}
		v = t.Parent[t.Head[v]]
	}
}

func (t *Tree32) LcaRooted(u, v, root int32) int32 {
	return t.Lca(u, v) ^ t.Lca(u, root) ^ t.Lca(v, root)
}

func (t *Tree32) Dist(a, b int32) int32 {
	c := t.Lca(a, b)
	return t.Depth[a] + t.Depth[b] - 2*t.Depth[c]
}

func (t *Tree32) DistWeighted(a, b int32) int {
	c := t.Lca(a, b)
	return t.DepthWeighted[a] + t.DepthWeighted[b] - 2*t.DepthWeighted[c]
}

// c 是否在 p 的子树中.c和p不能相等.
func (t *Tree32) InSubtree(c, p int32) bool {
	return t.Lid[p] <= t.Lid[c] && t.Lid[c] < t.Rid[p]
}

// 从 a 开始走 k 步到 b.
func (t *Tree32) Jump(a, b, k int32) int32 {
	if k == 1 {
		if a == b {
			return -1
		}
		if t.InSubtree(b, a) {
			return t.KthAncestor(b, t.Depth[b]-t.Depth[a]-1)
		}
		return t.Parent[a]
	}
	c := t.Lca(a, b)
	dac := t.Depth[a] - t.Depth[c]
	dbc := t.Depth[b] - t.Depth[c]
	if k > dac+dbc {
		return -1
	}
	if k <= dac {
		return t.KthAncestor(a, k)
	}
	return t.KthAncestor(b, dac+dbc-k)
}

func (t *Tree32) SubtreeSize(v int32) int32 {
	return t.Rid[v] - t.Lid[v]
}

func (t *Tree32) SubtreeSizeRooted(v, root int32) int32 {
	if v == root {
		return t.n
	}
	x := t.Jump(v, root, 1)
	if t.InSubtree(v, x) {
		return t.Rid[v] - t.Lid[v]
	}
	return t.n - t.Rid[x] + t.Lid[x]
}

func (t *Tree32) CollectChild(v int32) []int32 {
	var res []int32
	for _, e := range t.Tree[v] {
		if e.to != t.Parent[v] {
			res = append(res, e.to)
		}
	}
	return res
}

// 收集与 v 相邻的轻边.
func (t *Tree32) CollectLight(v int32) []int32 {
	var res []int32
	skip := true
	for _, e := range t.Tree[v] {
		if e.to != t.Parent[v] {
			if !skip {
				res = append(res, e.to)
			}
			skip = false
		}
	}
	return res
}

func (tree *Tree32) RestorePath(from, to int32) []int32 {
	res := []int32{}
	composition := tree.GetPathDecomposition(from, to, 0)
	for _, e := range composition {
		a, b := e[0], e[1]
		if a <= b {
			for i := a; i <= b; i++ {
				res = append(res, tree.IdToNode[i])
			}
		} else {
			for i := a; i >= b; i-- {
				res = append(res, tree.IdToNode[i])
			}
		}
	}
	return res
}

// 返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
//
//	!eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
func (tree *Tree32) GetPathDecomposition(u, v int32, edge int32) [][2]int32 {
	up, down := [][2]int32{}, [][2]int32{}
	lid, head, parent := tree.Lid, tree.Head, tree.Parent
	for {
		if head[u] == head[v] {
			break
		}
		if lid[u] < lid[v] {
			down = append(down, [2]int32{lid[head[v]], lid[v]})
			v = parent[head[v]]
		} else {
			up = append(up, [2]int32{lid[u], lid[head[u]]})
			u = parent[head[u]]
		}
	}
	if lid[u] < lid[v] {
		down = append(down, [2]int32{lid[u] + edge, lid[v]})
	} else if lid[v]+edge <= lid[u] {
		up = append(up, [2]int32{lid[u], lid[v] + edge})
	}
	for i := 0; i < len(down)/2; i++ {
		down[i], down[len(down)-1-i] = down[len(down)-1-i], down[i]
	}
	return append(up, down...)
}

// 遍历路径上的 `[起点,终点)` 欧拉序 `左闭右开` 区间.
func (tree *Tree32) EnumeratePathDecomposition(u, v int32, edge int32, f func(start, end int32)) {
	head, lid, parent := tree.Head, tree.Lid, tree.Parent
	for {
		if head[u] == head[v] {
			break
		}
		if lid[u] < lid[v] {
			a, b := lid[head[v]], lid[v]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			v = parent[head[v]]
		} else {
			a, b := lid[u], lid[head[u]]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			u = parent[head[u]]
		}
	}
	if lid[u] < lid[v] {
		a, b := lid[u]+edge, lid[v]
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	} else if lid[v]+edge <= lid[u] {
		a, b := lid[u], lid[v]+edge
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	}
}

// 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
func (tree *Tree32) Id(root int32) (int32, int32) {
	return tree.Lid[root], tree.Rid[root]
}

// 返回返回边 u-v 对应的 欧拉序起点编号, 1 <= eid <= n-1., 0-indexed.
func (tree *Tree32) Eid(u, v int32) int32 {
	if tree.Lid[u] > tree.Lid[v] {
		return tree.Lid[u]
	}
	return tree.Lid[v]
}

// 点v对应的父边的边id.如果v是根节点则返回-1.
func (tre *Tree32) VToE(v int32) int32 {
	return tre.vToE[v]
}

// 第i条边对应的深度更深的那个节点.
func (tree *Tree32) EToV(i int32) int32 {
	u, v := tree.Edges[i][0], tree.Edges[i][1]
	if tree.Parent[u] == v {
		return u
	}
	return v
}

func (tree *Tree32) ELid(u int32) int32 {
	return 2*tree.Lid[u] - tree.Depth[u]
}

func (tree *Tree32) ERid(u int32) int32 {
	return 2*tree.Rid[u] - tree.Depth[u] - 1
}

func (t *Tree32) _dfsSize(cur, pre int32) {
	size := t.Rid
	t.Parent[cur] = pre
	if pre != -1 {
		t.Depth[cur] = t.Depth[pre] + 1
	} else {
		t.Depth[cur] = 0
	}
	size[cur] = 1
	nexts := t.Tree[cur]
	for i := int32(len(nexts)) - 2; i >= 0; i-- {
		e := nexts[i+1]
		if t.Depth[e.to] == -1 {
			nexts[i], nexts[i+1] = nexts[i+1], nexts[i]
		}
	}
	hldSize := int32(0)
	for i, e := range nexts {
		to := e.to
		if t.Depth[to] == -1 {
			t.DepthWeighted[to] = t.DepthWeighted[cur] + e.cost
			t.vToE[to] = e.eid
			t._dfsSize(to, cur)
			size[cur] += size[to]
			if size[to] > hldSize {
				hldSize = size[to]
				if i != 0 {
					nexts[0], nexts[i] = nexts[i], nexts[0]
				}
			}
		}
	}
}

func (t *Tree32) _dfsHld(cur int32, times *int32) {
	t.Lid[cur] = *times
	*times++
	t.Rid[cur] += t.Lid[cur]
	t.IdToNode[t.Lid[cur]] = cur
	heavy := true
	for _, e := range t.Tree[cur] {
		to := e.to
		if t.Depth[to] > t.Depth[cur] {
			if heavy {
				t.Head[to] = t.Head[cur]
			} else {
				t.Head[to] = to
			}
			heavy = false
			t._dfsHld(to, times)
		}
	}
}

// 路径 [a,b] 与 [c,d] 的交集.
// 如果为空则返回 {-1,-1}，如果只有一个交点则返回 {x,x}，如果有两个交点则返回 {x,y}.
func (t *Tree32) PathIntersection(a, b, c, d int32) (int32, int32) {
	ab := t.Lca(a, b)
	ac := t.Lca(a, c)
	ad := t.Lca(a, d)
	bc := t.Lca(b, c)
	bd := t.Lca(b, d)
	cd := t.Lca(c, d)
	x := ab ^ ac ^ bc // meet(a,b,c)
	y := ab ^ ad ^ bd // meet(a,b,d)
	if x != y {
		return x, y
	}
	z := ac ^ ad ^ cd
	if x != z {
		x = -1
	}
	return x, x
}

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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
