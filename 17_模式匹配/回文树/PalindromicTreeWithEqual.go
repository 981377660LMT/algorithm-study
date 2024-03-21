// https://blog.csdn.net/Clove_unique/article/details/53750322
// https://ei1333.github.io/library/string/palindromic-tree.hpp
// https://oi-wiki.org/string/pam/
// 文字列 Sが与えられる. Parindromic Tree は Sに含まれるすべての回文を頂点とした木である.
// 長さが -1 と 0 の超頂点を用意する.
// 各頂点からは, その頂点に対応する回文の両端に同じ文字を 1文字加えてできる回文の頂点に辺を張ることで木を構成する.
// 特に長さ 1の回文は-1の超頂点から, 長さ 2の回文は 0の超頂点から辺を張ることになる.

// !さらに Suffix Link として, 頂点の回文に含まれる`最大の回文接尾辞`に対応する頂点に辺を張る.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// demo()
	P3501()
}

func demo() {
	// s := "eertree"
	s := "aaaaa"

	// 返回回文树、每个前缀结尾的回文串个数.
	NewPalindromicTreeAndEnds := func(n int, f func(int) int32) (tree *PalindromicTree, ends []int) {
		tree = NewPalindromicTree(func(i, j int) bool {
			if i == j {
				return true
			}
			return f(i) != f(j)
		})
		ends = make([]int, n)
		counter := make([]int, n+2)
		for i := 0; i < n; i++ {
			pos := tree.Add(f(i))
			link := tree.Nodes[pos].Link
			counter[pos] = counter[link] + 1
			ends[i] = counter[pos]
		}
		return
	}

	pt, ends := NewPalindromicTreeAndEnds(len(s), func(i int) int32 { return int32(s[i]) })
	for i := 0; i < pt.Size(); i++ {
		start, end := pt.RecoverPalindrome(i)
		fmt.Println(start, end, s[start:end])
	}
	fmt.Println("ends", ends)                  // 每个前缀结尾的回文串个数
	fmt.Println("count", pt.Size()-2)          // 本质不同回文子串个数(不超过O(n))
	fmt.Println("size", pt.GetFrequency()[2:]) // 每个顶点对应的回文串出现的次数

	fmt.Println(pt.BuildFailTree())
	for i := 2; i < pt.Size(); i++ {
		fmt.Println(pt.GetNode(i).Indexes)
	}
}

// P3501 [POI2010] ANT-Antisymmetry (!对于一个0/1序列，求出其中异或意义下回文的子串数量。)
// https://www.luogu.com.cn/problem/P3501
// 对于一个01字符串，如果将这个字符串0和1取反后，再将整个串反过来和原串一样，就称作“反对称”字符串。
// 比如00001111和010101就是反对称的，1001就不是。
// 现在给出一个长度为N的01字符串，求它有多少个子串是反对称的。
func P3501() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s01 string
	fmt.Fscan(in, &s01)

	p := NewPalindromicTree(func(i, j int) bool { return s01[i] == s01[j] })
	p.AddString(s01)
	freq := p.GetFrequency()
	res := 0
	for i := 2; i < p.Size(); i++ {
		res += freq[i]
		start, end := p.RecoverPalindrome(i)
		fmt.Println(start, end, s01[start:end])
	}
	fmt.Fprintln(out, res)

}

type Node struct {
	Next      map[int32]int32 // 子节点.
	Link      int32           // suffix link，指向当前回文串的最长真回文后缀的位置
	Length    int32           // 结点代表的回文串的长度
	Indexes   []int32         // 哪些位置的最长回文后缀
	HalfLink  int32           // 长度不超过 len(u)//2 的最长回文后缀的位置
	deltaLink int32           // u一直沿着link向上跳到第一使得diff[v] ≠ diff[u]的节点v，即u所在等差数列中长度最小的那个节点。
}

type PalindromicTree struct {
	Ords    []int32
	Nodes   []*Node
	lastPos int32 // 当前字符串(原串前缀)的最长回文后缀
	equal   func(i, j int) bool
}

func NewPalindromicTree(equal func(i, j int) bool) *PalindromicTree {
	res := &PalindromicTree{equal: equal}
	res.Nodes = append(res.Nodes, res.newNode(0, -1)) // 奇根，长为 -1
	res.Nodes = append(res.Nodes, res.newNode(0, 0))  // 偶根，长为 0
	return res
}

// !添加一个字符,返回以这个字符为后缀的最长回文串的位置pos.
// 每次增加一个字符，本质不同的回文子串个数最多增加 1 个.
// !以i位置结尾的回文串个数等于此节点在树中的深度.
func (pt *PalindromicTree) Add(x int32) int {
	pos := int32(len(pt.Ords))
	pt.Ords = append(pt.Ords, x)
	cur := pt.findPrevPalindrome(pt.lastPos)
	_, hasKey := pt.Nodes[cur].Next[x]
	if !hasKey {
		pt.Nodes[cur].Next[x] = int32(len(pt.Nodes))
	}
	pt.lastPos = pt.Nodes[cur].Next[x]
	if !hasKey {
		newNode := pt.newNode(-1, pt.Nodes[cur].Length+2)
		pt.Nodes = append(pt.Nodes, newNode)
		if newNode.Length == 1 {
			newNode.Link = 1
		} else {
			newNode.Link = pt.Nodes[pt.findPrevPalindrome(pt.Nodes[cur].Link)].Next[x]
		}

		if newNode.Length <= 2 {
			newNode.HalfLink = newNode.Link
		} else {
			halfNode := pt.Nodes[pt.Nodes[cur].HalfLink]
			for !pt.equal(int(pos), int(pos-halfNode.Length-1)) || (halfNode.Length+2)*2 > newNode.Length {
				halfNode = pt.Nodes[halfNode.Link]
			}
			newNode.HalfLink = halfNode.Next[x]
		}

		if pt.diff(pt.lastPos) == pt.diff(newNode.Link) {
			newNode.deltaLink = pt.Nodes[newNode.Link].deltaLink
		} else {
			newNode.deltaLink = newNode.Link
		}
	}

	pt.Nodes[pt.lastPos].Indexes = append(pt.Nodes[pt.lastPos].Indexes, pos)
	return int(pt.lastPos)
}

func (pt *PalindromicTree) AddString(s string) {
	if len(s) == 0 {
		return
	}
	for _, v := range s {
		pt.Add(v)
	}
}

// 给定节点位置,返回其代表的回文串.
func (pt *PalindromicTree) RecoverPalindrome(pos int) (start, end int) {
	if pos < 2 {
		return
	}
	node := pt.Nodes[pos]
	end = int(node.Indexes[0]) + 1
	start = end - int(node.Length)
	return
}

// Palindrome Series 优化DP
// https://zhuanlan.zhihu.com/p/537113907
// https://zhuanlan.zhihu.com/p/92874690
// https://www.cnblogs.com/Parsnip/p/12426971.html
// Palindrome Series 使用的情况为：枚举所有的回文后缀，这时就可以把dp转移的复杂度从O(n)变成O(logn),且常数极小。
// 在每次调用Add(x)之后使用,用以当前字符为结尾的`所有本质不同回文串`更新dp值.
//   - init(pos, start): 初始化顶点pos的dp值,对应回文串s[start:i+1].
//   - apply(pos, prePos): 用prePos(fail指针指向的位置)更新pos.
//     返回值: 以S[i]为结尾的回文的位置.
func (pt *PalindromicTree) UpdateDp(init func(pos, start int), apply func(pos, pre int)) (indexes []int) {
	i := int32(len(pt.Ords) - 1)
	id := pt.lastPos
	for pt.Nodes[id].Length > 0 {
		init(int(id), int(i+1-pt.Nodes[pt.Nodes[id].deltaLink].Length-pt.diff(id)))
		if pt.Nodes[id].deltaLink != pt.Nodes[id].Link {
			apply(int(id), int(pt.Nodes[id].Link))
		}
		indexes = append(indexes, int(id))
		id = pt.Nodes[id].deltaLink
	}
	return
}

// 按照拓扑序进行转移.
// from: 后缀连接, to: 当前节点
func (pt *PalindromicTree) Dp(f func(from, to int)) {
	for i := pt.Size() - 1; i >= 2; i-- {
		f(int(pt.Nodes[i].Link), i)
	}
}

// 求出每个顶点对应的回文串出现的次数.
func (pt *PalindromicTree) GetFrequency() []int {
	res := make([]int, pt.Size())
	// !节点编号从大到小，就是 fail 树的拓扑序
	for i := pt.Size() - 1; i >= 1; i-- { // 除去根节点(奇根)
		res[i] += len(pt.Nodes[i].Indexes)
		res[pt.Nodes[i].Link] += res[i] // 长回文包含短回文
	}
	return res
}

// eg: "eertree" -> [[1] [2 4 5] [3 7] [8] [6] [] [] [] []]
//
//											0(奇根)
//											|
//							-------	1(偶根)
//	           /        |  \
//		        /         |   \
//		       /          |		 \
//		      /           |     \
//		     /            |      \
//		    2(e)          4(r)   5(t)
//		   /    \         |
//			3(ee)  7(ertre) 6(rtr)
//		  |
//		  8(eertree)
//
// 0 为奇根, 1 为偶根.
// 每条链对应一个前缀，沿着链往下走，相当于遍历当前前缀的所有回文后缀.
func (pt *PalindromicTree) BuildFailTree() [][]int {
	n := pt.Size()
	res := make([][]int, n)
	for i := 1; i < n; i++ {
		link := int(pt.Nodes[i].Link)
		res[link] = append(res[link], i)
	}
	return res
}

// 回文树中的顶点个数.(包含两个奇偶虚拟顶点)
// 一个串的本质不同回文子串个数等于 Size()-2.
func (pt *PalindromicTree) Size() int {
	return len(pt.Nodes)
}

// 返回pos位置的回文串顶点.
func (pt *PalindromicTree) GetNode(pos int) *Node {
	return pt.Nodes[pos]
}

func (pt *PalindromicTree) newNode(link, length int32) *Node {
	return &Node{
		Next:      make(map[int32]int32),
		Link:      link,
		Length:    length,
		deltaLink: -1,
	}
}

// 沿着失配指针找到第一个满足 x+s+x 是原串回文后缀的位置.
func (pt *PalindromicTree) findPrevPalindrome(cur int32) int32 {
	pos := int32(len(pt.Ords) - 1)
	for {
		rev := pos - 1 - pt.Nodes[cur].Length
		// !插入当前字符的条件str[i]==str[i-len-1]
		if rev >= 0 && ((pt.Ords[rev] != pt.Ords[len(pt.Ords)-1]) || rev == int32(len(pt.Ords)-1)) {
			break
		}
		cur = pt.Nodes[cur].Link
	}
	return cur
}

// 当前位置的回文串长度减去当前回文串的最长后缀回文串的长度.
func (pt *PalindromicTree) diff(pos int32) int32 {
	curNode := pt.Nodes[pos]
	if curNode.Link <= 0 {
		return -1
	}
	return curNode.Length - pt.Nodes[curNode.Link].Length
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxs(nums []int) int {
	res := nums[0]
	for _, v := range nums {
		if v > res {
			res = v
		}
	}
	return res
}

func mins(nums []int) int {
	res := nums[0]
	for _, v := range nums {
		if v < res {
			res = v
		}
	}
	return res
}
