package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math/bits"
	"os"
	"runtime/debug"
	"strconv"
)

var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func init() {
	debug.SetGCPercent(-1)
}

// 区间最大异或和
// https://www.luogu.com.cn/problem/P4735
// https://www.acwing.com/problem/content/258/
// 给定一个非负整数序列nums，初始长度为n。
// 有q个操作，有以下两种操作类型:
//  - A x: 添加操作，表示在序列末尾添加一个数x，序列的长度+1。
//  - Q start end x: 询问操作，你需要找到一个位置p，
//    !满足start≤p<end，使得x与后缀的异或和 x^(nums[pos]^nums[pos+1]^...^nums[n-1]) 最大，输出最大是多少。
// !n,m<=3e5 0<=x<=1e7

// !解法
// 维护前缀异或 查询变为 preXor[pos]^preXor[n]^x
// 即求 前缀异或与(x^preXor[n])的最大值
// `用持久化trie01来维护前缀异或`，第i个版本为插入了nums[i]后的trie树
// 怎么查询区间[start,end)的最大值呢？
// !限制右端点-> 查询版本为end的trie树.
// !限制左端点-> 查询过程中节点的最后更新时间>=start.
// 这样就保证了查询的是区间[start,end)的最大值.
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, q := io.NextInt(), io.NextInt()
	nums := make([]int, n)
	for i := range nums {
		nums[i] = io.NextInt()
	}
	operations := make([][4]int, q)
	for i := range operations {
		op := io.Text()
		if op == "A" {
			x := io.NextInt()
			operations[i] = [4]int{0, x, 0, 0}
		} else {
			start, end, x := io.NextInt(), io.NextInt(), io.NextInt()
			start--
			operations[i] = [4]int{1, start, end, x}
		}
	}

	maxVersion := n + q + 1
	preXor := make([]int, 0, maxVersion)
	preXor = append(preXor, 0)

	trie := NewXorTrieSimplePersistent(int(1e7+10), true)
	roots := make([]*XorTrieNode, 0, maxVersion)
	initRoot := trie.NewRoot()
	roots = append(roots, trie.Insert(initRoot, 0, 0)) // 前缀和首部的0

	add := func(value int, index int) {
		preXor = append(preXor, preXor[len(preXor)-1]^value)
		roots = append(roots, trie.Insert(roots[len(roots)-1], preXor[len(preXor)-1], index))
	}

	for i, v := range nums {
		add(v, i+1)
	}

	curLen := n
	for _, operation := range operations {
		kind := operation[0]
		if kind == 0 {
			x := operation[1]
			add(x, curLen+1)
			curLen++
		} else {
			start, end, x := operation[1], operation[2], operation[3]
			// fmt.Println("start, end, x:", start, end, x, preXor, x^preXor[len(preXor)-1])
			// trie.Enumerate(roots[end-1], func(num int) { fmt.Println(num) })
			maxXor, _ := trie.Query(roots[end-1], x^preXor[len(preXor)-1], func(node *XorTrieNode) bool {
				return node.data >= start
			})
			fmt.Fprintln(out, maxXor)
		}
	}

}

func demo() {
	trie := NewXorTrieSimplePersistent(1e9, true)
	root := trie.NewRoot()
	root1 := trie.Insert(root, 1, 1)
	root2 := trie.Insert(root1, 2, 2)
	root3 := trie.Remove(root2, 1)
	fmt.Println(trie.Query(root1, 1, nil))
	fmt.Println(trie.Query(root2, 1, nil))
	trie.Enumerate(root1, func(num int) { fmt.Println(num) })
	trie.Enumerate(root2, func(num int) { fmt.Println(num) })
	trie.Enumerate(root3, func(num int) { fmt.Println(num) })
}

// https://leetcode.cn/problems/maximum-xor-of-two-numbers-in-an-array/
func findMaximumXOR(nums []int) int {
	max_ := 1
	for _, num := range nums {
		max_ = max(max_, num)
	}
	tree := NewXorTrieSimplePersistent(max_, false)
	root := tree.NewRoot()
	maxXor := 0
	for _, num := range nums {
		root = tree.Insert(root, num, 0)
		cand, _ := tree.Query(root, num, nil)
		maxXor = max(maxXor, cand)
	}
	return maxXor
}

type Data = int

type XorTrieNode struct {
	count    int
	chidlren [2]*XorTrieNode
	data     Data // 结点上的自定义信息，例如这个结点最后一次被更新的时间(lastIndex)
}

// 简单的可持久化01字典树.
type XorTrieSimplePersistent struct {
	bit        int
	persistent bool
}

func NewXorTrieSimplePersistent(upper int, persistent bool) *XorTrieSimplePersistent {
	return &XorTrieSimplePersistent{bit: bits.Len(uint(upper)), persistent: persistent}
}

func (trie *XorTrieSimplePersistent) NewRoot() *XorTrieNode {
	return nil
}

func (trie *XorTrieSimplePersistent) Copy(node *XorTrieNode) *XorTrieNode {
	if node == nil || !trie.persistent {
		return node
	}
	return &XorTrieNode{
		count:    node.count,
		chidlren: node.chidlren,
		data:     node.data, // copy
	}
}

func (trie *XorTrieSimplePersistent) Insert(root *XorTrieNode, num int, data Data) *XorTrieNode {
	if root == nil {
		root = &XorTrieNode{}
	}
	return trie._insert(root, num, trie.bit-1, data)
}

// 必须保证num存在于trie中.
func (trie *XorTrieSimplePersistent) Remove(root *XorTrieNode, num int) *XorTrieNode {
	if root == nil {
		panic("can not remove from nil root")
	}
	return trie._remove(root, num, trie.bit-1)
}

// 查询num与root中的数异或的最大值.
// !如果root为nil,返回 0,nil.
// !canSwap: 判断是否能切换到另一位的子树.传nil表示不需要自定义判断.
func (trie *XorTrieSimplePersistent) Query(root *XorTrieNode, num int, canSwap func(node *XorTrieNode) bool) (maxXor int, node *XorTrieNode) {
	if root == nil {
		return
	}

	if canSwap == nil {
		for k := trie.bit - 1; k >= 0; k-- {
			bit := (num >> k) & 1
			if root.chidlren[bit^1] != nil {
				bit ^= 1
				maxXor |= 1 << k
			}
			root = root.chidlren[bit]
		}
		return maxXor, root
	} else {
		for k := trie.bit - 1; k >= 0; k-- {
			bit := (num >> k) & 1
			if root.chidlren[bit^1] != nil && canSwap(root.chidlren[bit^1]) {
				bit ^= 1
				maxXor |= 1 << k
			}
			root = root.chidlren[bit]
		}
		return maxXor, root
	}
}

func (trie *XorTrieSimplePersistent) Enumerate(root *XorTrieNode, f func(num int)) {
	trie._dfs(root, trie.bit-1, 0, f)
}

func (trie *XorTrieSimplePersistent) _insert(root *XorTrieNode, num int, depth int, data Data) *XorTrieNode {
	root = trie.Copy(root)
	root.count++
	root.data = data
	if depth < 0 {
		return root
	}
	bit := (num >> depth) & 1
	if root.chidlren[bit] == nil {
		root.chidlren[bit] = &XorTrieNode{}
	}
	root.chidlren[bit] = trie._insert(root.chidlren[bit], num, depth-1, data)
	return root
}

func (trie *XorTrieSimplePersistent) _remove(root *XorTrieNode, num int, depth int) *XorTrieNode {
	root = trie.Copy(root)
	root.count--
	if depth < 0 {
		return root
	}
	bit := (num >> depth) & 1
	root.chidlren[bit] = trie._remove(root.chidlren[bit], num, depth-1)
	return root
}

func (trie *XorTrieSimplePersistent) _dfs(root *XorTrieNode, depth int, curXor int, f func(num int)) {
	if root == nil {
		return
	}
	if depth < 0 {
		if root.count > 0 {
			f(curXor)
		}
		return
	}
	trie._dfs(root.chidlren[0], depth-1, curXor, f)
	trie._dfs(root.chidlren[1], depth-1, curXor|(1<<depth), f)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
