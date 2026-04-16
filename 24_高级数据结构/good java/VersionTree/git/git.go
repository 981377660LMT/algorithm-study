// G - Copy Query
// https://atcoder.jp/contests/abc453/tasks/abc453_g
//
// 用 Git 的视角看：
// 1. 每个数列 A_i 是一个 branch.
// 2. 类型 1 操作是让 branch X checkout 到 branch Y 的 head.
// 3. 类型 2 操作是从当前 head 提交一个修改某个位置的新 commit.
// 4. 类型 3 操作是在当前 head 上挂一个查询，最后统一 DFS 执行.
//
// 1. Head(branch): 获取指定分支指向的 commit.
// 2. CopyBranch(dst, src): 让 dst 指向 src 当前的 head.
// 3. Reset(branch, commit): 把 branch 重置到某个历史 commit.
// 4. Commit(branch, apply, invert): 在指定分支上提交一个新 commit.
// 5. Query(branch, query): 在指定分支 head 上挂查询.
// 6. Execute(): 统一 DFS 执行所有 commit 和查询.
// 7. CommitFrom(branch, parent, apply, invert): 从指定父版本派生新 commit，
//    适合“基于第 k 个版本生成新版本”的题，不必先 Reset 再 Commit.
//
// - CommitGraph
// 	负责 commit 节点、父子关系、挂查询、DFS 执行 apply/invert。
// - Git
// 	负责 branch head、复制分支、重置、提交修改、挂查询。

package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://atcoder.jp/contests/abc453/tasks/abc453_g
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	fmt.Fscan(in, &n, &m, &q)

	git := NewGit(int32(n), int32(q))
	values := make([]int, m)
	bit := NewBitArray(m)
	res := make([]int, q)
	for i := range res {
		res[i] = -1
	}

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		switch op {
		case 1:
			var x, y int
			fmt.Fscan(in, &x, &y)
			git.CopyBranch(BranchID(x-1), BranchID(y-1))
		case 2:
			var x, y, z int
			fmt.Fscan(in, &x, &y, &z)
			branch := BranchID(x - 1)
			pos := y - 1
			newValue := z
			var oldValue int
			git.Commit(branch, func() bool {
				oldValue = values[pos]
				bit.Add(pos, newValue-oldValue)
				values[pos] = newValue
				return true
			}, func() {
				bit.Add(pos, oldValue-newValue)
				values[pos] = oldValue
			})
		case 3:
			var x, l, r int
			fmt.Fscan(in, &x, &l, &r)
			queryIndex := i
			left := l - 1
			right := r
			git.Query(BranchID(x-1), func() {
				res[queryIndex] = bit.QueryRange(left, right)
			})
		}
	}

	git.Execute()
	for _, v := range res {
		if v != -1 {
			fmt.Fprintln(out, v)
		}
	}
}

type CommitID int32
type BranchID int32

type Git struct {
	graph *CommitGraph
	heads []CommitID
}

func NewGit(maxBranchCount int32, maxCommitCount int32) *Git {
	graph := NewCommitGraph(maxCommitCount)
	heads := make([]CommitID, maxBranchCount)
	for i := range heads {
		heads[i] = graph.Root()
	}
	return &Git{graph: graph, heads: heads}
}

func (g *Git) BranchCount() int32 {
	return int32(len(g.heads))
}

func (g *Git) Head(branchID BranchID) CommitID {
	return g.heads[branchID]
}

// 把一个分支指针直接挪到另一个分支当前指向的位置.
func (g *Git) CopyBranch(dst, src BranchID) {
	g.heads[dst] = g.heads[src]
}

func (g *Git) Reset(branchID BranchID, commitID CommitID) {
	g.heads[branchID] = commitID
}

// 在指定分支上继续提交一个 commit.
func (g *Git) Commit(branchID BranchID, apply func() bool, invert func()) CommitID {
	head := g.heads[branchID]
	newHead := g.graph.CommitFrom(head, apply, invert)
	g.heads[branchID] = newHead
	return newHead
}

// 基于指定父 commit 创建一个新 commit，再把 branch 指向这个新 commit.
func (g *Git) CommitFrom(branchID BranchID, parent CommitID, apply func() bool, invert func()) CommitID {
	newHead := g.graph.CommitFrom(parent, apply, invert)
	g.heads[branchID] = newHead
	return newHead
}

func (g *Git) Query(branchID BranchID, query func()) {
	g.graph.AttachQuery(g.heads[branchID], query)
}

func (g *Git) Execute() {
	g.graph.Execute()
}

type commitNode struct {
	parent   CommitID
	children []CommitID
	queries  []func()
	apply    func() bool
	invert   func()
}

type CommitGraph struct {
	nodes []commitNode
}

func NewCommitGraph(maxCommitCount int32) *CommitGraph {
	nodes := make([]commitNode, 1, maxCommitCount+1)
	nodes[0] = commitNode{
		parent: -1,
		apply:  func() bool { return false },
		invert: func() {},
	}
	return &CommitGraph{nodes: nodes}
}

func (g *CommitGraph) Root() CommitID {
	return 0
}

// 从 parent 派生一个新 commit.
func (g *CommitGraph) CommitFrom(parent CommitID, apply func() bool, invert func()) CommitID {
	if apply == nil {
		apply = func() bool { return false }
	}
	if invert == nil {
		invert = func() {}
	}
	newID := CommitID(len(g.nodes))
	g.nodes = append(g.nodes, commitNode{
		parent: parent,
		apply:  apply,
		invert: invert,
	})
	g.nodes[parent].children = append(g.nodes[parent].children, newID)
	return newID
}

func (g *CommitGraph) AttachQuery(commitID CommitID, query func()) {
	if query == nil {
		return
	}
	g.nodes[commitID].queries = append(g.nodes[commitID].queries, query)
}

func (g *CommitGraph) Execute() {
	g.dfs(g.Root())
}

func (g *CommitGraph) dfs(root CommitID) {
	node := &g.nodes[root]
	ok := node.apply()
	for _, query := range node.queries {
		query()
	}
	for _, child := range node.children {
		g.dfs(child)
	}
	if ok {
		node.invert()
	}
}

// Point Add Range Sum, 0-based.
type BITArray struct {
	n     int
	total int
	data  []int
}

func NewBitArray(n int) *BITArray {
	return &BITArray{n: n, data: make([]int, n)}
}

func (b *BITArray) Add(index int, v int) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

func (b *BITArray) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

func (b *BITArray) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}
