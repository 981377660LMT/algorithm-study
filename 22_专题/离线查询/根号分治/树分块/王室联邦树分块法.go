// https://github.com/EndlessCheng/codeforces-go/blob/29153e2c702970aaccc69db6c4739c3103f95429/copypasta/graph_tree.go#LL1585C2-L1585C2
// https://github.com/hos-lyric/libra/blob/60b8b56ecae5860f81d75de28510d94336f5dad9/data_structure/tree_sqrt_decomposition.cpp#L15

// https://www.luogu.com.cn/blog/gxy001/shu-fen-kuai-xue-xi-bi-ji
// https://www.cnblogs.com/hua-dong/p/8275227.html
// https://ouuan.github.io/post/%E8%8E%AB%E9%98%9F%E5%B8%A6%E4%BF%AE%E8%8E%AB%E9%98%9F%E6%A0%91%E4%B8%8A%E8%8E%AB%E9%98%9F%E8%AF%A6%E8%A7%A3/#%E5%88%86%E5%9D%97%E6%96%B9%E5%BC%8F
// https://oi-wiki.org/ds/tree-decompose/
//
// ---
//
// https://www.cnblogs.com/IzayoiMiku/p/14691521.html
// 树分块的几种方法:
//
// - topCluster 分块法
// !- 王室联邦分块法
// !分块方式:满足每块大小在 [B,3B]，块内每个点到核心点路径上的所有点都在块内.
// !但是不保证每个块都是连通的.
//
// dfs，并创建一个栈，dfs一个点时先记录初始栈顶高度，
// 每dfs完当前节点的一棵子树就判断栈内（相对于刚开始dfs时）新增节点的数量是否≥B，
// 是则将栈内所有新增点分为同一块，核心点为当前dfs的点，
// 当前节点结束dfs时将当前节点入栈，
// 整个dfs结束后将栈内所有剩余节点归入已经分好的最后一个块。

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	n := 5
	tree := [][]int{{1}, {0, 2, 3}, {1, 4}, {1}, {2}}
	blockSize := 2
	fmt.Println(LimitSizeDecompose(n, tree, blockSize))
}

func P2325() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, b int
	fmt.Fscan(in, &n, &b)
	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	blockId, roots := LimitSizeDecompose(n, tree, b)
	fmt.Fprintln(out, len(roots))
	for _, v := range blockId {
		fmt.Fprint(out, v+1, " ")
	}
	fmt.Fprintln(out)
	for _, v := range roots {
		fmt.Fprint(out, v+1, " ")
	}
}

// 王室联邦分块法
//
//	https://github.com/EndlessCheng/codeforces-go/blob/29153e2c702970aaccc69db6c4739c3103f95429/copypasta/graph_tree.go#LL1585C2-L1585C2
//	https://www.luogu.com.cn/problem/P2325
func LimitSizeDecompose(n int, tree [][]int, blockSize int) (belong []int, blockRoot []int) {
	blockRoot = []int{}     // 每个块的根节点(关键点)
	belong = make([]int, n) // belong[i]表示节点i所属的块的编号, blockRoot[belong[i]]即为i所属块的根节点(关键点)
	stack := []int{}
	var dfs func(cur, pre int)
	dfs = func(cur, pre int) {
		size := len(stack)
		for _, next := range tree[cur] {
			if next != pre {
				dfs(next, cur)
				if len(stack)-size >= blockSize {
					blockRoot = append(blockRoot, cur)
					for len(stack) > size {
						belong[stack[len(stack)-1]] = len(blockRoot) - 1
						stack = stack[:len(stack)-1]
					}
				}
			}
		}
		stack = append(stack, cur)
	}

	dfs(0, -1)
	if len(blockRoot) == 0 {
		blockRoot = []int{0}
	}
	for _, v := range stack {
		belong[v] = len(blockRoot) - 1
	}
	return
}
