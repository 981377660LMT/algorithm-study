package main

type VersionTree struct {
	nodes   []*TreeNode
	version int
}

// 创建一个新的版本树，maxOperation 表示最大操作数.
// 版本号从 0 开始.
func NewVersionTree(maxOperation int) *VersionTree {
	nodes := make([]*TreeNode, 0, maxOperation+1)
	nodes = append(nodes, NewTreeNode(EMPTY_OPERATION))
	return &VersionTree{nodes: nodes}
}

// 在当前版本上添加一个操作，返回新版本号.
func (t *VersionTree) Apply(operation *Operation) int {
	newNode := NewTreeNode(operation)
	t.nodes = append(t.nodes, newNode)
	t.nodes[t.version].children = append(t.nodes[t.version].children, newNode)
	t.version = len(t.nodes) - 1
	return t.version
}

// 切换到指定版本.
func (t *VersionTree) SwitchVersion(version int) {
	t.version = version
}

// 应用所有操作.
func (t *VersionTree) Run() {
	t.dfs(t.nodes[0])
}

// 获取当前版本号.
func (t *VersionTree) GetVersion() int {
	return t.version
}

func (t *VersionTree) dfs(root *TreeNode) {
	root.operation.apply()
	for _, child := range root.children {
		t.dfs(child)
	}
	root.operation.undo()
}

type TreeNode struct {
	children  []*TreeNode
	operation *Operation
}

func NewTreeNode(operation *Operation) *TreeNode {
	return &TreeNode{operation: operation}
}

var EMPTY_OPERATION = NewOperation(func() {}, func() {})

type Operation struct {
	apply func()
	undo  func()
}

func NewOperation(apply, undo func()) *Operation {
	return &Operation{apply: apply, undo: undo}
}
