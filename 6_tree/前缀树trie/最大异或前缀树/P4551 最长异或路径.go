// 给定一棵 n 个点的带权树，结点下标从 0 开始到 n。
// 寻找树中找两个结点，求最长的异或路径。
// 异或路径指的是指两个结点之间唯一路径上的所有边权的异或。
//
// 树上差分
// 树上 x 到 y 的路径上所有边权的 xor 结果就等于 `D[x] xor D[y]`。
// !其中 D[x]表示根节点到 x 的异或值,重叠路径抵消了(前缀异或)
// 所以，`问题就变成了从 D[1]~D[N]这 N 个数中选出两个，xor 的结果最大`
// 时间复杂度O(n)

// https://www.luogu.com.cn/problem/P4551

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	edges := make([][3]int, n-1)
	for i := range edges {
		fmt.Fscan(in, &edges[i][0], &edges[i][1], &edges[i][2])
		edges[i][0]--
		edges[i][1]--
	}

	tree := make([][][2]int, n)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		tree[u] = append(tree[u], [2]int{v, w})
		tree[v] = append(tree[v], [2]int{u, w})
	}

	res := 0
	trie := NewXORTrie(1 << 32)
	xorToRoot := make([]int, 0, n-1)
	var dfs func(cur, parent, preXor int)
	dfs = func(cur, parent, preXor int) {
		xorToRoot = append(xorToRoot, preXor)
		for _, e := range tree[cur] {
			next, weight := e[0], e[1]
			if next == parent {
				continue
			}
			dfs(next, cur, preXor^weight)
		}
	}
	dfs(0, -1, 0)

	for _, v := range xorToRoot {
		trie.Insert(v)
		res = max(res, trie.Query(v))
	}

	fmt.Fprintln(out, res)
}

type Node struct {
	count    int
	children [2]*Node // 数组比 left,right 更快
}

type XORTrieSimple struct {
	bit  int
	root *Node
}

func NewXORTrie(upper int) *XORTrieSimple {
	return &XORTrieSimple{
		bit:  bits.Len(uint(upper)),
		root: &Node{},
	}
}

func (bt *XORTrieSimple) Insert(num int) {
	root := bt.root
	for i := bt.bit - 1; i >= 0; i-- {
		bit := (num >> i) & 1
		if root.children[bit] == nil {
			root.children[bit] = &Node{}
		}
		root = root.children[bit]
		root.count++
	}
	return
}

// 必须保证num存在于trie中.
func (bt *XORTrieSimple) Remove(num int) {
	root := bt.root
	for i := bt.bit - 1; i >= 0; i-- {
		bit := (num >> i) & 1
		root = root.children[bit]
		root.count--
	}
}

func (bt *XORTrieSimple) Query(num int) (maxXor int) {
	root := bt.root
	for i := bt.bit - 1; i >= 0; i-- {
		if root == nil {
			return
		}
		bit := (num >> i) & 1
		if root.children[bit^1] != nil && root.children[bit^1].count > 0 {
			maxXor |= 1 << i
			bit ^= 1
		}
		root = root.children[bit]
	}
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
