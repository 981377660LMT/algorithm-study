package main

func EnumerateTree[T any](
	root T, getChildren func(node T) []T,
	visit func(node T) (stop bool),
) {
	var dfs func(node T) (stop bool)
	dfs = func(node T) (stop bool) {
		if visit(node) {
			return true
		}
		for _, child := range getChildren(node) {
			if dfs(child) {
				return true
			}
		}
		return false
	}

	dfs(root)
}

func main() {
	// Example usage
	type Node struct {
		Value    int
		Children []Node
	}

	root := Node{
		Value: 1,
		Children: []Node{
			{Value: 2, Children: []Node{{Value: 4}, {Value: 5}}},
			{Value: 3, Children: []Node{{Value: 6}}},
		},
	}

	getChildren := func(node Node) []Node {
		return node.Children
	}

	visit := func(node Node) bool {
		println(node.Value)
		return false // continue visiting
	}

	EnumerateTree(root, getChildren, visit)
}
