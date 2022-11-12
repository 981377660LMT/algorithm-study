package eulertour

import (
	"fmt"
	"math/bits"
)

func demo() {
	// https://maspypy.com/euler-tour-%e3%81%ae%e3%81%8a%e5%8b%89%e5%bc%b7 里的例子
	// 0-1 0-5 1-2 1-4 2-3
	tree := [][]edge{
		{{1, 1}, {5, 5}},
		{{0, 0}, {2, 2}, {4, 4}},
		{{1, 1}, {3, 3}},
		{{2, 2}},
		{{1, 1}},
		{{0, 0}},
	}
	nodeValues := []int{1, 2, 3, 4, 5, 6} // 每个节点的权值
	et := NewEulerTour(len(tree), tree, 0)

	nodePresum := make([]int, 1, len(et.Tour)) // 结点权值前缀和，欧拉序列一一对应 (如果要修改权值，需要树状数组维护)
	edgePresum := make([]int, 1, len(et.Tour)) // 边权值前缀和，欧拉序列一一对应 (如果要修改权值，需要树状数组维护)
	var dfs func(cur, pre int)
	dfs = func(cur, pre int) {
		nodePresum = append(nodePresum, nodeValues[cur]+nodePresum[len(nodePresum)-1])
		for _, edge := range tree[cur] {
			if next := edge.next; next != pre {
				edgePresum = append(edgePresum, edge.weight+edgePresum[len(edgePresum)-1])
				dfs(next, cur)
				edgePresum = append(edgePresum, -edge.weight+edgePresum[len(edgePresum)-1])      // !回溯添加逆元
				nodePresum = append(nodePresum, -nodeValues[next]+nodePresum[len(nodePresum)-1]) // !回溯添加逆元
			}
		}
	}
	dfs(0, -1)
	fmt.Println(nodePresum, edgePresum)

	// 求2-5路径上的点权和
	res1 := 0
	et.QueryPathNode(2, 5, func(l, r int) {
		res1 += nodePresum[r] - nodePresum[l]
	})
	fmt.Println(res1) // 12

	// 求2-5路径上的边权和
	res2 := 0
	et.QueryPathEdge(2, 5, func(l, r int) {
		fmt.Println(l, r)
		res2 += edgePresum[r-1] - edgePresum[l-1] // !注意边权这里要减1
	})
	fmt.Println(res2) // 8
}

// !具有逆元的monoid (群)可以由欧拉序列O(logn)查询
// !不具有逆元的monoid (幺半群)可以由重链剖分O(logn*logn)查询
// 由于预处理 ST 表是基于一个长度为 2n 的序列，所以常数上是比倍增算法要大的。内存占用也比倍增要大一倍左右（这点可忽略）
// 优点是查询的复杂度低，适用于查询量大的情形
// https://nyaannyaan.github.io/library/tree/euler-tour.hpp
// https://github.dev/EndlessCheng/codeforces-go/blob/cca30623b9ac0f3333348ca61b4894cd00b753cc/copypasta/graph_tree.go#L822
// https://maspypy.com/euler-tour-%e3%81%ae%e3%81%8a%e5%8b%89%e5%bc%b7
// https://blog.csdn.net/weixin_45539557/article/details/114483543
type EulerTour struct {
	n, root int
	tree    [][]edge

	// 欧拉序列
	Tour []int
	//进入节点的时间戳,从0开始编号
	down []int
	// 离开节点的时间戳,从0开始编号
	up []int
	// 深度序列,和欧拉序列一一对应,根节点深度为0
	depth []int
	// 到根节点的距离
	distToRoot []int

	// RMQ查询深度序列区间最小值
	rmq [][]pair
}

type edge struct{ next, weight int }
type pair struct{ i, v int }
type callback = func(left, right int)

func NewEulerTour(n int, tree [][]edge, root int) *EulerTour {
	eulerTour := &EulerTour{n: n, tree: tree, root: root}
	eulerTour.buildTour()
	eulerTour.buildRMQ()
	return eulerTour
}

func (euler *EulerTour) LCA(root1, root2 int) int {
	down1, down2 := euler.down[root1], euler.down[root2]
	if down1 > down2 {
		down1, down2 = down2, down1
	}
	minDepthIndex := euler.queryRMQ(down1, down2+1)
	return euler.Tour[minDepthIndex]
}

// 路径点权(欧拉序列+前缀和计算)
// https://maspypy.com/euler-tour-%e3%81%ae%e3%81%8a%e5%8b%89%e5%bc%b7
func (euler *EulerTour) QueryPathNode(root1, root2 int, cb callback) {
	lca := euler.LCA(root1, root2)
	cb(euler.down[lca], euler.down[root1]+1)   // 链1
	cb(euler.down[lca]+1, euler.down[root2]+1) // 链2
}

// 路径边权(欧拉序列+前缀和计算)
// https://maspypy.com/euler-tour-%e3%81%ae%e3%81%8a%e5%8b%89%e5%bc%b7
func (euler *EulerTour) QueryPathEdge(root1, root2 int, cb callback) {
	lca := euler.LCA(root1, root2)
	cb(euler.down[lca]+1, euler.down[root1]+1) // 链1
	cb(euler.down[lca]+1, euler.down[root2]+1) // 链2
}

func (euler *EulerTour) QuerySubTree(root int, cb callback) {
	cb(euler.down[root], euler.up[root])
}

func (euler *EulerTour) QueryDist(root1, root2 int) int {
	return euler.distToRoot[root1] + euler.distToRoot[root2] - 2*euler.distToRoot[euler.LCA(root1, root2)]
}

func (euler *EulerTour) buildTour() {
	euler.Tour = make([]int, 0, euler.n*2-1)
	euler.down = make([]int, euler.n)
	euler.up = make([]int, euler.n)
	euler.depth = make([]int, 0, euler.n*2-1)
	euler.distToRoot = make([]int, euler.n)
	euler.buildTourDfs(euler.root, -1, 0, 0)
}

func (euler *EulerTour) buildTourDfs(cur, pre, dep, dist int) {
	euler.down[cur] = len(euler.Tour)
	euler.Tour = append(euler.Tour, cur)
	euler.depth = append(euler.depth, dep)
	euler.distToRoot[cur] = dist
	for _, edge := range euler.tree[cur] {
		if edge.next != pre {
			euler.buildTourDfs(edge.next, cur, dep+1, dist+edge.weight)
			euler.Tour = append(euler.Tour, cur)
			euler.up[cur] = len(euler.Tour)
			euler.depth = append(euler.depth, dep)
		}
	}
}

func (euler *EulerTour) buildRMQ() {
	n := len(euler.Tour)
	size := bits.Len(uint(n))
	euler.rmq = make([][]pair, size)
	for i := range euler.rmq {
		euler.rmq[i] = make([]pair, n)
	}
	for i := 0; i < n; i++ {
		euler.rmq[0][i] = pair{i, euler.depth[i]}
	}
	for i := 1; i < size; i++ {
		for j := 0; j+(1<<i) <= n; j++ {
			a, b := euler.rmq[i-1][j], euler.rmq[i-1][j+(1<<(i-1))]
			if a.v < b.v {
				euler.rmq[i][j] = a
			} else {
				euler.rmq[i][j] = b
			}
		}
	}
}

// 0-based
func (euler *EulerTour) queryRMQ(left, right int) int {
	k := bits.Len(uint(right-left)) - 1
	a, b := euler.rmq[k][left], euler.rmq[k][right-1<<k]
	if a.v < b.v {
		return a.i
	}
	return b.i
}
