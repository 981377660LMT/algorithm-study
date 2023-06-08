// 异或最小生成树
// https://www.luogu.com.cn/problem/CF888G
// 给定 n 个结点的无向完全图。每个点有一个点权ai
// 点i和点j之间的边权为ai⊕aj
// 求最小生成树的权值和
// n<=2e5 ai<=2**30

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	trie := NewBinaryTrie(30)
	counter := make(map[int]map[int]struct{}) // val => {id1, id2, ...} 维护id
	for i := 0; i < n; i++ {
		trie.Add(nums[i], 1)
		if counter[nums[i]] == nil {
			counter[nums[i]] = make(map[int]struct{})
		}
		counter[nums[i]][i] = struct{}{}
	}

	edges, _ := OnlineMST(
		n,
		func(u int) {
			trie.Add(nums[u], -1)
			delete(counter[nums[u]], u)
		},
		func(u int) {
			trie.Add(nums[u], 1)
			counter[nums[u]][u] = struct{}{}
		},
		func(u int) (id, cost int) {
			trie.XorAll(nums[u])
			cost, ok := trie.Min()
			trie.XorAll(nums[u])
			if !ok {
				return -1, -1
			}

			ids := counter[nums[u]^cost]
			if len(ids) == 0 {
				return -1, -1
			}
			for k := range ids {
				return k, cost
			}

			return -1, -1
		},
	)

	res := 0
	for _, e := range edges {
		res += e[2]
	}
	fmt.Fprintln(out, res)
}

// Brouvka
//  不预先给出图，而是指定一个函数 findUnused 来找到未使用过的点中与u权值最小的点。
//  findUnused(u)：返回 unused 中与 u 权值最小的点 v 和边权 cost
//                如果不存在，返回 (-1,*)
func OnlineMST(
	n int,
	setUsed func(u int), setUnused func(u int), findUnused func(u int) (v int, cost int),
) (res [][3]int, ok bool) {
	uf := NewUnionFindArray(n)
	for {
		updated := false
		groups := make([][]int, n)
		cand := make([][3]int, n) // [u, v, cost]
		for v := 0; v < n; v++ {
			cand[v] = [3]int{-1, -1, -1}
		}

		for v := 0; v < n; v++ {
			groups[uf.Find(v)] = append(groups[uf.Find(v)], v)
		}

		for v := 0; v < n; v++ {
			if uf.Find(v) != v {
				continue
			}
			for _, x := range groups[v] {
				setUsed(x)
			}
			for _, x := range groups[v] {
				y, cost := findUnused(x)
				if y == -1 {
					continue
				}
				a, c := cand[v][0], cand[v][2]
				if a == -1 || cost < c {
					cand[v] = [3]int{x, y, cost}
				}
			}
			for _, x := range groups[v] {
				setUnused(x)
			}
		}

		for v := 0; v < n; v++ {
			if uf.Find(v) != v {
				continue
			}
			a, b, c := cand[v][0], cand[v][1], cand[v][2]
			if a == -1 {
				continue
			}
			updated = true
			if uf.Union(a, b) {
				res = append(res, [3]int{a, b, c})
			}
		}

		if !updated {
			break
		}
	}

	if len(res) != n-1 {
		return nil, false
	}
	return res, true
}

func NewUnionFindArray(n int) *_UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_UnionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type _UnionFindArray struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *_UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}

	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *_UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *_UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *_UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *_UnionFindArray) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
}

func (ufa *_UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}

type BinaryTrie struct {
	xorLazy int
	root    *node
	maxLog  int
}

type node struct {
	count int // 以该节点为结尾的数的个数
	next  [2]*node
}

// maxLog: max log of x
// useId: whether to record the id of the number
func NewBinaryTrie(maxLog int) *BinaryTrie {
	return &BinaryTrie{root: newNode(), maxLog: maxLog}
}

func (bt *BinaryTrie) Add(num, count int) {
	bt.root = bt.add(bt.root, num, bt.maxLog, count, true)
}

func (bt *BinaryTrie) Count(num int) int {
	node := bt.find(bt.root, num, bt.maxLog)
	if node == nil {
		return 0
	}
	return node.count
}

// 0<=k<exist.
func (bt *BinaryTrie) Kth(k int) (res int, ok bool) {
	return bt.kthElement(bt.root, k, bt.maxLog)
}

func (bt *BinaryTrie) Max() (res int, ok bool) {
	return bt.Kth(bt.root.count - 1)
}

func (bt *BinaryTrie) Min() (res int, ok bool) {
	return bt.Kth(0)
}

func (bt *BinaryTrie) CountLess(num int) int {
	return bt.countLess(bt.root, num, bt.maxLog)
}

func (bt *BinaryTrie) BisectLeft(num int) int {
	return bt.CountLess(num)
}

func (bt *BinaryTrie) CountLessOrEqual(num int) int {
	return bt.CountLess(num + 1)
}

func (bt *BinaryTrie) BisectRight(num int) int {
	return bt.CountLessOrEqual(num)
}

func (bt *BinaryTrie) XorAll(x int) {
	bt.xorLazy ^= x
}

func (bt *BinaryTrie) add(t *node, bit, depth, x int, need bool) *node {
	if depth == -1 {
		t.count += x
	} else {
		f := (bt.xorLazy >> depth) & 1
		to := &t.next[f^((bit>>depth)&1)]
		if *to == nil {
			*to = newNode()
			need = false
		}
		*to = bt.add(*to, bit, depth-1, x, false)
		t.count += x
	}
	return t
}

func (bt *BinaryTrie) find(t *node, bit, depth int) *node {
	if depth == -1 {
		return t
	}
	f := (bt.xorLazy >> depth) & 1
	to := t.next[f^((bit>>depth)&1)]
	if to == nil {
		return nil
	}
	return bt.find(to, bit, depth-1)
}

func (bt *BinaryTrie) kthElement(t *node, k, bitIndex int) (int, bool) {
	if t == nil {
		return 0, false
	}
	if bitIndex == -1 || t == nil {
		return 0, true
	}
	f := (bt.xorLazy >> bitIndex) & 1
	count := 0
	if t.next[f] != nil {
		count = t.next[f].count
	}
	if count <= k {
		res, ok := bt.kthElement(t.next[f^1], k-count, bitIndex-1)
		res |= 1 << bitIndex
		return res, ok
	}
	return bt.kthElement(t.next[f], k, bitIndex-1)
}

func (bt *BinaryTrie) countLess(t *node, bit, bitIndex int) int {
	if bitIndex == -1 {
		return 0
	}
	res := 0
	f := (bt.xorLazy >> bitIndex) & 1
	if (bit>>bitIndex)&1 == 1 && t.next[f] != nil {
		res += t.next[f].count
	}
	if t.next[f^(bit>>bitIndex&1)] != nil {
		res += bt.countLess(t.next[f^(bit>>bitIndex&1)], bit, bitIndex-1)
	}
	return res
}

func newNode() *node {
	return &node{}
}
