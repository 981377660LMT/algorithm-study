package main

func treeDpRecursively(n int, tree [][]int) (parent, order []int) {
	parent = make([]int, n)
	for i := range parent {
		parent[i] = -2
	}
	parent[0] = -1
	stack := []int{0}
	order = make([]int, 0, n)
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range tree[u] {
			if parent[v] == -2 {
				parent[v] = u
				stack = append(stack, v)
			}
		}
	}
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	return
}

func main() {
	n := 5
	edges := [][]int{
		{0, 1},
		{0, 2},
		{1, 3},
		{1, 4},
	}
	tree := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	parent, order := treeDpRecursively(n, tree)
	// Output the parent and order arrays
	for i := 0; i < n; i++ {
		println("Parent of", i, ":", parent[i])
	}
	for i := 0; i < n; i++ {
		println("Order", i, ":", order[i])
	}
}
