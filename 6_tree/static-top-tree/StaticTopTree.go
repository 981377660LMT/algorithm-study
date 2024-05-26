// !StaticTopTree: 描述子树合并过程的深度为O(logN)的二叉树.
// !StaticTopTreeDp适用于：单点修改+全局dp值查询.
//

// !- 树形dp与树分解的联系：
//
//	type T = int
//	vertex := func (u int) T { return nums[u] }
//	addEdge := func (res T) T { return res }
//	addVertex := func (res T, u int) T { return d + nums[u] }
//	merge := func (l, r T) T { return l * r }
//	dp := func (u int) T {
//	  if len(tree[u]) == 0 {
//	!    return vertex(u)
//	  }
//	  children := make([]T, 0)
//	  for _, v := range tree[u] {
//	    childRes := dp(v)
//	!    children = append(children, addEdge(childRes))
//	  }
//	  res := children[0]
//	  for i := 1; i < len(children); i++ {
//	!    res = merge(res, children[i])
//	  }
//	!  return addVertex(res, u)
//	}
//
// !- 树形dp的另一种视角：StaticTopTree的五个操作
// 1. Vertex: 添加叶子节点.
// 2. AddEdge: 用一条轻边+虚拟根节点连接簇.
// 3. AddVertex: 用一个实际顶点代替虚拟根节点.
//     	     x       ◯
//       ->  ｜  ->  ｜
//     ◯     ◯ 		   ◯

// 4. Rake: 合并两个具有公共虚拟根节点的簇.
//     x    x          x
//     ｜ + ｜  ->    /  \
//     ◯    ◯  	    ◯    ◯

//  5. Compress: 用一条重边合并两个簇
//        ◯             ◯
//       /   		       /
//      ◯  +    ->    /
//     /             /
//    ◯             ◯
//

// !- StaticTopTree 构建(分解)方法:
//  1. 删除与根节点相连的一条重边(compress).
//  2. 删除根节点(保留与根节点相连的边，用虚拟根节点代替)（addVertex）.
//  3. 将虚拟根节点连接的子树分割 (rake).
//  4. 删除虚拟根节点连接每个子树的轻边，产生出了一个新的子树 (addEdge).
//
// !- 两种Cluster:
//  1. PointCluster: 根为虚拟根节点的子树形态.合并两个pointCluster的操作为Rake.
//  2. PathCluster: 根为真实顶点，重边连接多个子树.合并两个pathCluster的操作为Compress.

// ! - pointCluster は木 DP における部分木と対応している
//   - pathCluster 同士のマージは「変な形をした何か」同士のマージになる，一番難しい.

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	// abc351g()
}

// [ABC351G] Hash on Tree (DynamicTreeHash，动态树哈希)
// https://www.luogu.com.cn/problem/AT_abc351_g
// !dp[parent] = value[parent] + dp[child1] * dp[child2] * ... * dp[childn]
//
// cluster の外部との接点を、Top tree の用語を借りて boundary vertex と呼ぶことにします。
// path cluster は基本的に「根に近い方の boundary vertex」「根から遠い方の boundary vertex」という 2 つの boundary vertex を持ちます。
// (path cluster が根付き木の場合は例外で、2 つの boundary vertex が同一の頂点になります。)
// 問題によっては、この 2 つの boundary vertex に注目することで DP の遷移を構成しやすくなります。
// 適切な観察により、path cluster は「根から遠い方の boundary vertex にハッシュ値がx である根付き木が結合した時に、
// !根に近い方の boundary vertex を根とする根付き木のハッシュ値はax+b になる」となる値(a,b) を持てば良いことがわかります。
// このように path cluster に載せる情報を定義すると compress はアフィン関数の合成として定義することが出来ます。
// その他の関数も適切な考察により次のような関数を実装すればよいことがわかります。
func abc351g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	const MOD int = 998244353

	var n, q int32
	fmt.Fscan(in, &n, &q)
	tree := make([][]int32, n)
	for i := int32(1); i < n; i++ {
		var p int32
		fmt.Fscan(in, &p)
		p--
		tree[p] = append(tree[p], i)
	}
	nums := make([]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	stt := NewStaticTopTree(tree, 0)

	type Point = int // 子树哈希值

	// !pathCluster中离根远的端点哈希值为x时，离根近的端点哈希值为mul*x+add
	// 最后整棵树的哈希值即为add.
	type Path = struct{ mul, add int } // mul*x+add

	vertex := func(v int32) Path {
		return Path{mul: 1, add: nums[v]}
	}
	addEdge := func(p Path) Point {
		return p.add
	}
	addVertex := func(x Point, v int32) Path {
		return Path{mul: x, add: nums[v]}
	}
	rake := func(x, y Point) Point {
		return x * y % MOD
	}
	// p簇离根节点更近.
	// compress就是仿射变换.
	compress := func(p, c Path) Path {
		return Path{mul: p.mul * c.mul % MOD, add: (p.mul*c.add + p.add) % MOD}
	}

	dp := NewStaticTopTreeDP(stt, vertex, addEdge, addVertex, rake, compress)
	for i := int32(0); i < q; i++ {
		var v, x int
		fmt.Fscan(in, &v, &x)
		v--
		nums[v] = x
		newRes := dp.Update(int32(v))
		fmt.Fprintln(out, newRes.add)
	}
}

// https://atcoder.jp/contests/abc351/editorial/9868
// https://zenn.dev/blue_jam/articles/0526c70b74f6bb
type StaticTopTreeDP[PointCluster, PathCluster any] struct {
	stt *StaticTopTree

	// 节点初始值.
	vertex func(v int32) PathCluster
	// 子树的贡献值.
	addEdge func(t PathCluster) PointCluster
	// 将点v加入子树，向父节点返回答案.
	addVertex func(t PointCluster, v int32) PathCluster
	// 合并两个子树的贡献.
	rake func(x, y PointCluster) PointCluster
	// 合并两个pathCluster，p簇离根节点更近.
	compress func(p, c PathCluster) PathCluster

	dpPoint []PointCluster
	dpPath  []PathCluster
}

func NewStaticTopTreeDP[PointCluster, PathCluster any](
	tree *StaticTopTree,
	vertex func(v int32) PathCluster,
	addEdge func(t PathCluster) PointCluster,
	addVertex func(t PointCluster, v int32) PathCluster,
	rake func(x, y PointCluster) PointCluster,
	compress func(p, c PathCluster) PathCluster,
) *StaticTopTreeDP[PointCluster, PathCluster] {
	n := tree.Size()
	dpPoint := make([]PointCluster, n)
	dpPath := make([]PathCluster, n)
	res := &StaticTopTreeDP[PointCluster, PathCluster]{
		stt:    tree,
		vertex: vertex, addVertex: addVertex, addEdge: addEdge, rake: rake, compress: compress,
		dpPoint: dpPoint, dpPath: dpPath,
	}
	res._dfs(tree.Root)
	return res
}

func (sttdp *StaticTopTreeDP[PointCluster, PathCluster]) Update(u int32) PathCluster {
	nodes := sttdp.stt.Nodes
	// push up
	for u != -1 {
		sttdp._modify(u)
		u = nodes[u].p
	}
	return sttdp.dpPath[sttdp.stt.Root]
}

func (sttdp *StaticTopTreeDP[PointCluster, PathCluster]) _modify(k int32) {
	nodes, dpPath, dpPoint := sttdp.stt.Nodes, sttdp.dpPath, sttdp.dpPoint
	node := nodes[k]
	switch node.op {
	case vertex:
		dpPath[k] = sttdp.vertex(k)
	case compress:
		dpPath[k] = sttdp.compress(dpPath[node.l], dpPath[node.r])
	case rake:
		dpPoint[k] = sttdp.rake(dpPoint[node.l], dpPoint[node.r])
	case addEdge:
		dpPoint[k] = sttdp.addEdge(dpPath[node.l])
	case addVertex:
		dpPath[k] = sttdp.addVertex(dpPoint[node.l], k)
	}
}

func (sttdp *StaticTopTreeDP[PointCluster, PathCluster]) _dfs(u int32) {
	if u == -1 {
		return
	}
	node := sttdp.stt.Nodes
	sttdp._dfs(node[u].l)
	sttdp._dfs(node[u].r)
	sttdp._modify(u)
}

type opType int8

const (
	vertex    opType = iota // 叶子节点
	addVertex               // 非叶子节点
	addEdge
	rake
	compress
)

type sttNode struct {
	op      opType
	l, r, p int32
}

func newSttNode(op opType, l, r int32) *sttNode {
	return &sttNode{op: op, l: l, r: r, p: -1}
}

type StaticTopTree struct {
	Root  int32
	Tree  [][]int32
	Nodes []*sttNode // 子树合并(分解)过程.
}

func NewStaticTopTree(rootedTree [][]int32, root int32) *StaticTopTree {
	res := &StaticTopTree{Tree: rootedTree}
	res._dfs(root)
	n := len(rootedTree)
	vs := make([]*sttNode, n, 4*n)
	for i := 0; i < n; i++ {
		vs[i] = newSttNode(vertex, -1, -1)
	}
	res.Nodes = vs
	a, _ := res._compress(root)
	res.Root = a
	res.Nodes = res.Nodes[:len(res.Nodes):len(res.Nodes)]
	return res
}

func (stt *StaticTopTree) Size() int32 {
	return int32(len(stt.Nodes))
}

func (stt *StaticTopTree) _dfs(u int32) int32 {
	size, heavy := int32(1), int32(0)
	for i, v := range stt.Tree[u] {
		subsize := stt._dfs(v)
		size += subsize
		if heavy < subsize {
			heavy = subsize
			stt.Tree[u][0], stt.Tree[u][i] = stt.Tree[u][i], stt.Tree[u][0]
		}
	}
	return size
}

func (stt *StaticTopTree) _compress(u int32) (int32, int32) {
	a, b := stt._addVertex(u)
	chs := make([][2]int32, 0)
	chs = append(chs, [2]int32{a, b})
	tree := stt.Tree
	for len(tree[u]) > 0 {
		u = tree[u][0] // heavy child
		a, b := stt._addVertex(u)
		chs = append(chs, [2]int32{a, b})
	}
	return stt._merge(compress, chs)
}

func (stt *StaticTopTree) _makeNode(t opType, l, r, k int32) int32 {
	if k == -1 {
		k = int32(len(stt.Nodes))
		stt.Nodes = append(stt.Nodes, newSttNode(t, l, r))
	} else {
		stt.Nodes[k] = newSttNode(t, l, r)
	}
	if l != -1 {
		stt.Nodes[l].p = k
	}
	if r != -1 {
		stt.Nodes[r].p = k
	}
	return k
}

func (stt *StaticTopTree) _merge(t opType, a [][2]int32) (int32, int32) {
	if len(a) == 1 {
		a, b := a[0][0], a[0][1]
		return a, b
	}
	sizeSum := int32(0)
	for i := 0; i < len(a); i++ {
		sizeSum += a[i][1]
	}
	var b, c [][2]int32
	for i := 0; i < len(a); i++ {
		it, size := a[i][0], a[i][1]
		if sizeSum > size {
			b = append(b, [2]int32{it, size})
		} else {
			c = append(c, [2]int32{it, size})
		}
		sizeSum -= size * 2
	}
	l, lSize := stt._merge(t, b)
	r, rSize := stt._merge(t, c)
	return stt._makeNode(t, l, r, -1), lSize + rSize
}

func (stt *StaticTopTree) _addEdge(u int32) (int32, int32) {
	it, size := stt._compress(u)
	return stt._makeNode(addEdge, it, -1, -1), size
}

func (stt *StaticTopTree) _rake(u int32) (int32, int32) {
	var chs [][2]int32
	tree := stt.Tree
	for i := 1; i < len(tree[u]); i++ {
		a, b := stt._addEdge(tree[u][i])
		chs = append(chs, [2]int32{a, b})
	}
	return stt._merge(rake, chs)
}

func (stt *StaticTopTree) _addVertex(u int32) (int32, int32) {
	if len(stt.Tree[u]) < 2 {
		return stt._makeNode(vertex, -1, -1, u), 1
	} else {
		it, size := stt._rake(u)
		return stt._makeNode(addVertex, it, -1, u), size + 1
	}
}

// 无根树转有根树.
func ToRootedTree32(tree [][]int32, root int32) [][]int32 {
	n := len(tree)
	rootedTree := make([][]int32, n)
	visited := make([]bool, n)
	visited[root] = true
	queue := []int32{root}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, next_ := range tree[cur] {
			if !visited[next_] {
				visited[next_] = true
				queue = append(queue, next_)
				rootedTree[cur] = append(rootedTree[cur], next_)
			}
		}
	}
	return rootedTree
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
