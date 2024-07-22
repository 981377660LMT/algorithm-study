// 版本树/操作树

package main

type VersionTree struct {
	nodes   []*treeNode
	version int
}

// 创建一个新的版本树，maxOperation 表示最大操作数.
// 版本号从 0 开始.
func NewVersionTree(maxOperation int) *VersionTree {
	nodes := make([]*treeNode, 0, maxOperation+1)
	nodes = append(nodes, newTreeNode2(emptyOperation))
	return &VersionTree{nodes: nodes}
}

// 在当前版本上添加一个操作，返回新版本号.
func (t *VersionTree) Apply(apply, undo func()) int {
	newNode := newTreeNode(apply, undo)
	t.nodes = append(t.nodes, newNode)
	t.nodes[t.version].children = append(t.nodes[t.version].children, newNode)
	t.version = len(t.nodes) - 1
	return t.version
}

// 切换到指定版本.
func (t *VersionTree) SwitchVersion(version int) { t.version = version }

// 应用所有操作.
func (t *VersionTree) Run() { t.dfs(t.nodes[0]) }

// 获取当前版本号.
func (t *VersionTree) Version() int { return t.version }

func (t *VersionTree) dfs(root *treeNode) {
	root.operation.apply()
	for _, child := range root.children {
		t.dfs(child)
	}
	root.operation.undo()
}

type treeNode struct {
	children  []*treeNode
	operation *operation
}

func newTreeNode(apply, undo func()) *treeNode { return newTreeNode2(newOperation(apply, undo)) }
func newTreeNode2(op *operation) *treeNode     { return &treeNode{operation: op} }

type operation struct{ apply, undo func() }

var emptyOperation = newOperation(func() {}, func() {})

func newOperation(apply, undo func()) *operation { return &operation{apply: apply, undo: undo} }
