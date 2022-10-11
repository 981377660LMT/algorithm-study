// // 动态开点线段树

// package main

// type segmentTreeNode struct {
// 	leftNode, rightNode *segmentTreeNode
// 	left, right         int32
// 	value               int64
// 	isLazy              bool
// 	lazyValue           int64
// }

// func createSegmentTreeNode(left, right int32) *segmentTreeNode {
// 	node := &segmentTreeNode{}
// 	node.left = left
// 	node.right = right
// 	node.value = e()
// 	return node
// }

// // !动态开点线段树
// type SegmentTree struct {
// 	root         *segmentTreeNode
// 	lower, upper int32
// }

// func CreateSegmentTree(lower, upper int32) *SegmentTree {
// 	tree := &SegmentTree{}
// 	tree.lower = lower
// 	tree.upper = upper
// 	tree.root = createSegmentTreeNode(lower, upper)
// 	return tree
// }

// func pushUp(root *segmentTreeNode) {
// 	root.value = op(root.leftNode.value, root.rightNode.value)
// }

// func pushDown(root *segmentTreeNode, left, mid, right int32) {
// 	if root.leftNode == nil {
// 		root.leftNode = createSegmentTreeNode(left, mid)
// 	}

// 	if root.rightNode == nil {
// 		root.rightNode = createSegmentTreeNode(mid+1, right)
// 	}

// 	if root.isLazy {
// 		root.leftNode.isLazy = true
// 		root.rightNode.isLazy = true
// 		root.leftNode.lazyValue = updateTree(root.leftNode.lazyValue, root.lazyValue)
// 		root.rightNode.lazyValue = updateTree(root.rightNode.lazyValue, root.lazyValue)
// 		root.leftNode.value = updateTree(root.leftNode.value, root.lazyValue)
// 		root.rightNode.value = updateTree(root.rightNode.value, root.lazyValue)
// 		root.isLazy = false
// 		root.lazyValue = e()
// 	}
// }

// func (tree *SegmentTree) Build(nums []int64) {
// 	tree.build(tree.lower, tree.upper, tree.root, nums)
// }

// func (tree *SegmentTree) Update(left, right int32, value int64) {
// 	if left < tree.lower {
// 		left = tree.lower
// 	}

// 	if right > tree.upper {
// 		right = tree.upper
// 	}

// 	if left > right {
// 		return
// 	}

// 	tree.update(left, right, tree.lower, tree.upper, tree.root, value)
// }

// func (tree *SegmentTree) Query(left, right int32) int64 {
// 	if left < tree.lower {
// 		left = tree.lower
// 	}

// 	if right > tree.upper {
// 		right = tree.upper
// 	}

// 	if left > right {
// 		return e()
// 	}

// 	return tree.query(left, right, tree.lower, tree.upper, tree.root)
// }

// func (tree *SegmentTree) QueryAll() int64 {
// 	return tree.root.value
// }

// func (tree *SegmentTree) update(L, R, l, r int32, root *segmentTreeNode, value int64) {
// 	if L <= l && r <= R {
// 		root.isLazy = true
// 		root.lazyValue = updateTree(root.lazyValue, value)
// 		root.value = updateTree(root.value, value)
// 		return
// 	}

// 	mid := (l + r) >> 1
// 	pushDown(root, l, mid, r)
// 	if L <= mid {
// 		tree.update(L, R, l, mid, root.leftNode, value)
// 	}

// 	if R > mid {
// 		tree.update(L, R, mid+1, r, root.rightNode, value)
// 	}

// 	pushUp(root)
// }

// func (tree *SegmentTree) query(L, R, l, r int32, root *segmentTreeNode) int64 {
// 	if L <= l && r <= R {
// 		return root.value
// 	}

// 	mid := (l + r) >> 1
// 	pushDown(root, l, mid, r)
// 	res := e()

// 	if L <= mid {
// 		leftValue := tree.query(L, R, l, mid, root.leftNode)
// 		res = op(res, leftValue)
// 	}

// 	if R > mid {
// 		rightValue := tree.query(L, R, mid+1, r, root.rightNode)
// 		res = op(res, rightValue)
// 	}

// 	return res
// }

// func (tree *SegmentTree) build(l, r int32, root *segmentTreeNode, nums []int64) {
// 	root.left, root.right = l, r
// 	if l == r {
// 		root.value = nums[l-1]
// 		return
// 	}

// 	mid := (l + r) >> 1
// 	pushDown(root, l, mid, r)
// 	tree.build(l, mid, root.leftNode, nums)
// 	tree.build(mid+1, r, root.rightNode, nums)
// 	pushUp(root)
// }

// func min(a, b int64) int64 {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

// func max(a, b int64) int64 {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

// // !tree operations
// // 结点的的初始值
// func e() int64 {
// 	return 1e18 // 幺元
// }

// // 合并左右节点的value op需要满足结合律
// func op(leftValue, rightValue int64) int64 {
// 	return min(leftValue, rightValue)
// }

// // 结点的更新方式
// func updateTree(oldValue, newValue int64) int64 {
// 	return min(oldValue, newValue)
// 	return newValue
// }

// func main() {
// 	tree := CreateSegmentTree(1, 10)
// 	println(tree.QueryAll())

// 	tree.Update(1, 10, 1)
// 	println(tree.QueryAll())

// 	tree.Update(1, 5, 1)
// 	println(tree.QueryAll())
// }

// !TODO 有问题