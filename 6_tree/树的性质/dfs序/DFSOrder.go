// DFS 序
//            0 [0,5)
//           /       \
//          /         \
//        1 [0,3)      2 [3,4)
//       /    \
//      /      \
//    3 [0,1)   4[1,2)

package main

import "fmt"

func main() {
	n := 5
	edges := [][]int{{0, 1}, {0, 2}, {1, 3}, {1, 4}}
	tree := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}
	dfsOrder := NewDfsOrder(n, tree, 0)
	fmt.Println(dfsOrder.QueryRange(1, 1))
	fmt.Println(dfsOrder.QueryRange(2, 2))
	fmt.Println(dfsOrder.QueryRange(3, 3))
	fmt.Println(dfsOrder.QueryId(1))
	fmt.Println(dfsOrder.QueryId(2))
	fmt.Println(dfsOrder.QueryId(3))
}

type DfsOrder struct {
	Starts, Ends []int
	n            int
	dfsId        int // 从0开始
	tree         [][]int
}

func NewDfsOrder(n int, tree [][]int, root int) *DfsOrder {
	res := &DfsOrder{
		Starts: make([]int, n),
		Ends:   make([]int, n),
		n:      n,
		tree:   tree,
	}

	res.dfs(root, -1)
	return res
}

// 求子树映射到的区间[start, end).
//  0 <= start < end <= n.
func (d *DfsOrder) QueryRange(u, v int) (start, end int) {
	start, end = d.Starts[v], d.Ends[v]+1
	return
}

// 求root自身的dfs序.
//  0 <= id < n.
func (d *DfsOrder) QueryId(root int) (id int) {
	return d.Ends[root]
}

// 判断root是否是child的祖先.
func (d *DfsOrder) IsAncestor(root, child int) bool {
	left1, right1 := d.Starts[root], d.Ends[root]
	left2, right2 := d.Starts[child], d.Ends[child]
	return left1 <= left2 && right2 <= right1
}

func (d *DfsOrder) dfs(u, fa int) {
	d.Starts[u] = d.dfsId
	for _, v := range d.tree[u] {
		if v != fa {
			d.dfs(v, u)
		}
	}
	d.Ends[u] = d.dfsId
	d.dfsId++
}
