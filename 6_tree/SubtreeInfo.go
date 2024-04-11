package main

import "fmt"

func main() {
	//   0
	//  / \
	// 1   2
	//    / \
	//   3   4
	//      /
	//     5

	tree := [][]int32{{1, 2}, {0}, {0, 3, 4}, {2}, {2, 5}, {4}}
	parent, depth, subsize, height := GetSubtreeInfo(tree, 0)
	fmt.Println(parent, depth, subsize, height)
}

// 获取子树信息.
// height[i] 表示以 i 为根的子树的高度(距离最远的叶子节点的距离).
func GetSubtreeInfo(tree [][]int32, root int32) (parent, depth, subsize, height []int32) {
	n := int32(len(tree))
	parent, depth, subsize, height = make([]int32, n), make([]int32, n), make([]int32, n), make([]int32, n)
	topological := make([]int32, n)
	topological[0] = root
	parent[root], depth[root] = root, 0
	r := int32(1)
	for l := int32(0); l < r; l++ {
		i := topological[l]
		for _, j := range tree[i] {
			if j != parent[i] {
				topological[r] = j
				r++
				parent[j], depth[j] = i, depth[i]+1
			}
		}
	}

	for r--; r >= 0; r-- {
		i := topological[r]
		subsize[i], height[i] = 1, 0
		for _, j := range tree[i] {
			if j != parent[i] {
				subsize[i] += subsize[j]
				height[i] = max32(height[i], height[j]+1)
			}
		}
	}
	parent[root] = -1
	return
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
