package main

type Node struct {
	Bag      []int
	Children []int
}

// 宽度为2的树分解.
// O(V+E).
// https://ei1333.hateblo.jp/entry/2020/02/12/150319
type TreeDecompositionWidth2 struct {
	isTreeWidth2 bool
	nodes        []*Node
}

func NewTreeDecompositionWidth2(tree [][]int) *TreeDecompositionWidth2 {
	nodes := make([]*Node, 1)
	n := len(tree)
	uf := _NewUf(n)
	deg := make([]int, n)
	nexts := make([]map[int]struct{}, n)
	for v := 0; v < n; v++ {
		curNexts := make(map[int]struct{})
		nexts[v] = curNexts
		for _, u := range tree[v] {
			uf.Union(u, v)
			curNexts[u] = struct{}{}
		}
		deg[v] = len(nexts[v])
	}

	// 加虚拟结点使得图联通.
	leaders := make([]int, 0)
	for v := 0; v < n; v++ {
		if uf.Find(v) == v {
			leaders = append(leaders, v)
		}
	}
	for i := 0; i < len(leaders)-1; i++ {
		u := leaders[i]
		v := leaders[i+1]
		nexts[u][v] = struct{}{}
		nexts[v][u] = struct{}{}
		deg[u]++
		deg[v]++
	}

	// -2: removed and added to the tree
	// -1: not removed
	// >= 0: removed and not yet added to the tree
	states := make([]int, n)
	for i := range states {
		states[i] = -1
	}
	stack := make([]int, 0)
	for v := 0; v < n; v++ {
		if deg[v] <= 2 {
			stack = append(stack, v)
		}
	}

	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if states[v] != -1 {
			continue
		}
		node := &Node{Bag: []int{v}}
		x, y := -1, -1
		for u := range nexts[v] {
			if states[u] == -1 {
				if x == -1 {
					x = u
				} else {
					y = u
				}
				node.Bag = append(node.Bag, u)
			} else if states[u] > 0 {
				node.Children = append(node.Children, states[u])
				states[u] = -2
			}
		}
		if x != -1 {
			if y == -1 {
				deg[x]--
			} else {
				if _, ok := nexts[x][y]; !ok {
					nexts[x][y] = struct{}{}
					nexts[y][x] = struct{}{}
				} else {
					deg[x]--
					deg[y]--
				}
			}
		}
		for u := range nexts[v] {
			if states[u] == -1 && deg[u] <= 2 {
				stack = append(stack, u)
			}
		}
		deg[v] = 0
		states[v] = len(nodes)
		nodes = append(nodes, node)
	}

	treewidthIs2 := true
	for _, d := range deg {
		if d > 0 {
			treewidthIs2 = false
			break
		}
	}
	if treewidthIs2 {
		nodes[0].Children = append(nodes[0].Children, len(nodes)-1)
	}

	return &TreeDecompositionWidth2{
		isTreeWidth2: treewidthIs2,
		nodes:        nodes,
	}
}

func (jtw2 *TreeDecompositionWidth2) IsTreeWidth2() bool {
	return jtw2.isTreeWidth2
}

func (jtw2 *TreeDecompositionWidth2) GetNodes() []*Node {
	return jtw2.nodes
}

func (jtw2 *TreeDecompositionWidth2) Size() int {
	return len(jtw2.nodes)
}

func _NewUf(n int) *_UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_UnionFindArray{
		Part:   n,
		size:   n,
		rank:   rank,
		parent: parent,
	}
}

type _UnionFindArray struct {
	size   int
	Part   int
	rank   []int
	parent []int
}

func (ufa *_UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *_UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *_UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}
