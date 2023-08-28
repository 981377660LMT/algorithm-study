package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	words := make([]string, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &words[i])
	}

	var q int
	fmt.Fscan(in, &q)
	opt := make([][5]int, q)
	for i := 0; i < q; i++ {
		var kind, w1, w2, start, end int
		fmt.Fscan(in, &kind, &w1, &w2, &start, &end)
		w1--
		w2--
		start--
		opt[i] = [5]int{kind, w1, w2, start, end}
	}

	res := HashSwapping(words, opt)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// 3000ms
// https://atcoder.jp/contests/soundhound2018-summer-final/tasks/soundhound2018_summer_final_e
// opt[i] = [kind, w1, w2, start, end]
// kind == 1: 交换w1[start:end]和w2[start:end]
// kind == 2: 输出w1[start:end]的哈希值
// c1c2...cn的哈希值为 ord(c1)*BASE^(n-1) + ord(c2)*BASE^(n-2) + ... + ord(cn)*BASE^0 模MOD
func HashSwapping(words []string, opt [][5]int) []int {
	roots := make([]*Node, len(words))
	for i, word := range words {
		leaves := make([]E, len(word))
		for j := 0; j < len(word); j++ {
			leaves[j] = [2]int{BASE, myOrd(word[j])}
		}
		roots[i] = Build(leaves)
	}

	res := []int{}
	for _, o := range opt {
		kind, w1, w2, start, end := o[0], o[1], o[2], o[3], o[4]
		if kind == 1 {
			xl, xm, xr := Split3ByRank(roots[w1], start, end)
			yl, ym, yr := Split3ByRank(roots[w2], start, end)
			roots[w1] = Merge(xl, Merge(ym, xr))
			roots[w2] = Merge(yl, Merge(xm, yr))
		} else {
			cur := Query(&roots[w1], start, end)
			res = append(res, cur[1])
		}
	}
	return res
}

const INF int = 1e18
const MOD int = 1e9 + 7
const BASE int = 1e6

// myOrd(a) = 1, myOrd(b) = 2, ..., myOrd(z) = 26
var myOrd = func(c byte) int {
	return int(c - 'a' + 1)
}

type E = [2]int // [pow of base, hash]
func e() E      { return [2]int{1, 0} }
func op(e1, e2 E) E {
	return [2]int{(e1[0] * e2[0]) % MOD, (e2[1]*e1[0] + e1[1]) % MOD}
}

//
//
//

// 每个结点代表一段区间
type Node struct {
	left, right *Node
	value       E
	data        E
	size        int
	priority    uint64
}

func NewRoot() *Node {
	return nil
}

func NewNode(v E) *Node {
	res := &Node{value: v, data: v, size: 1, priority: _nextRand()}
	return res
}

// 合并`左右`两棵树，保证Rank有序.
func Merge(left, right *Node) *Node {
	if left == nil || right == nil {
		if left == nil {
			return right
		}
		return left
	}

	if left.priority > right.priority {
		left.right = Merge(left.right, right)
		return _pushUp(left)
	} else {
		right.left = Merge(left, right.left)
		return _pushUp(right)
	}
}

// split root to [0,k) and [k,n)
func SplitByRank(root *Node, k int) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}

	leftSize := Size(root.left)
	if k <= leftSize {
		first, second := SplitByRank(root.left, k)
		root.left = second
		return first, _pushUp(root)
	} else {
		first, second := SplitByRank(root.right, k-leftSize-1)
		root.right = first
		return _pushUp(root), second
	}
}

func Split3ByRank(root *Node, l, r int) (*Node, *Node, *Node) {
	if root == nil {
		return nil, nil, nil
	}
	root, right := SplitByRank(root, r)
	left, mid := SplitByRank(root, l)
	return left, mid, right
}

// Fold.
func Query(node **Node, start, end int) (res E) {
	if start >= end || node == nil {
		return e()
	}
	left1, right1 := SplitByRank(*node, start)
	left2, right2 := SplitByRank(right1, end-start)
	if left2 != nil {
		res = left2.data
	} else {
		res = e()
	}
	*node = Merge(left1, Merge(left2, right2))
	return
}

func Size(node *Node) int {
	if node == nil {
		return 0
	}
	return node.size
}

func _pushUp(node *Node) *Node {
	node.size = 1
	node.data = node.value
	if left := node.left; left != nil {
		node.size += left.size
		node.data = op(left.data, node.data)
	}
	if right := node.right; right != nil {
		node.size += right.size
		node.data = op(node.data, right.data)
	}
	return node
}

var _seed = uint64(time.Now().UnixNano()/2 + 1)

func _nextRand() uint64 {
	_seed ^= _seed << 7
	_seed ^= _seed >> 9
	return _seed
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------

// 按照leaves的顺序(不是值的顺序)构建一棵树.
func Build(leaves []E) *Node {
	if len(leaves) == 0 {
		return nil
	}
	var dfs func(l, r int) *Node
	dfs = func(l, r int) *Node {
		if r-l == 1 {
			node := NewNode(leaves[l])
			return _pushUp(node)
		}
		mid := (l + r) >> 1
		return Merge(dfs(l, mid), dfs(mid, r))
	}
	return dfs(0, len(leaves))
}
