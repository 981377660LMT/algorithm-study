package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	CF932G()
}

// Palindrome Partition (偶回文串划分)
// https://www.luogu.com.cn/problem/CF932G
// 将一个字符串分成偶数段,假设分为s1,s2,...,s2k,要求 s[i]=s[2k-i](1<=i<2k).求合法的分割方案数模1e9+7.
//
// !构造一个字符串 T= s[1], s[2k], s[2], s[2k-1], ..., s[k], s[k+1]。
// !等价于求T的偶回文划分数(每个回文串长度为偶数).
func CF932G() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e9 + 7

	var s string
	fmt.Fscan(in, &s)
	if len(s)&1 == 1 {
		fmt.Fprintln(out, 0)
		return
	}

	n := len(s)
	t := make([]byte, 0, n)
	for i := 0; i < n/2; i++ {
		t = append(t, (s[i]))
		t = append(t, (s[n-1-i]))
	}

	counter := make([]int, n+2) // 回文树上每个位置对应的dp值
	dp := make([]int, n+1)      // dp[i]表示前i个字符的最少分割次数
	dp[0] = 1
	tree := NewPalindromicTreeArray()
	for i, c := range t {
		tree.Add(c)

		indexes := tree.UpdateDp(
			// 基于 s[start:i+1] 这一段回文初始化 pos 处的值
			func(pos, start int) {
				counter[pos] = dp[start]
			},
			// 基于 pre 的信息更新 pos 处的值
			func(pos, pre int) {
				counter[pos] += counter[pre]
				if counter[pos] >= MOD {
					counter[pos] -= MOD
				}
			},
		)

		if i&1 == 1 {
			for _, p := range indexes {
				dp[i+1] += counter[p]
				if dp[i+1] >= MOD {
					dp[i+1] -= MOD
				}
			}
		}

	}

	fmt.Fprintln(out, dp[n])
}

const SIGMA byte = 26
const OFFSET byte = 97

type Node struct {
	Next      [SIGMA]int32
	Link      int32 // suffix link，指向当前回文串的最长真回文后缀的位置
	Length    int32 // 结点代表的回文串的长度
	deltaLink int32 // u一直沿着link向上跳到第一使得diff[v] ≠ diff[u]的节点v，即u所在等差数列中长度最小的那个节点。
}

type PalindromicTreeArray struct {
	Ords    []byte
	Nodes   []*Node
	lastPos int32 // 当前字符串(原串前缀)的最长回文后缀
}

func NewPalindromicTreeArray() *PalindromicTreeArray {
	res := &PalindromicTreeArray{}
	res.Nodes = append(res.Nodes, res.newNode(0, -1)) // 奇根，长为 -1
	res.Nodes = append(res.Nodes, res.newNode(0, 0))  // 偶根，长为 0
	return res
}

// !添加一个字符,返回以这个字符为后缀的最长回文串的位置pos.
// 每次增加一个字符，本质不同的回文子串个数最多增加 1 个.
func (pt *PalindromicTreeArray) Add(x byte) int {
	x -= OFFSET
	pt.Ords = append(pt.Ords, x)
	cur := pt.findPrevPalindrome(pt.lastPos)
	hasKey := pt.Nodes[cur].Next[x] != -1
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

		if pt.diff(pt.lastPos) == pt.diff(newNode.Link) {
			newNode.deltaLink = pt.Nodes[newNode.Link].deltaLink
		} else {
			newNode.deltaLink = newNode.Link
		}
	}

	return int(pt.lastPos)
}

func (pt *PalindromicTreeArray) AddString(s string) {
	if len(s) == 0 {
		return
	}
	for i := 0; i < len(s); i++ {
		pt.Add(s[i])
	}
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
func (pt *PalindromicTreeArray) UpdateDp(init func(pos, start int), apply func(pos, pre int)) (indexes []int) {
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
func (pt *PalindromicTreeArray) Dp(f func(from, to int)) {
	for i := pt.Size() - 1; i >= 2; i-- {
		f(int(pt.Nodes[i].Link), i)
	}
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
func (pt *PalindromicTreeArray) BuildFailTree() [][]int {
	n := pt.Size()
	res := make([][]int, n)
	for i := 1; i < n; i++ {
		link := int(pt.Nodes[i].Link)
		res[link] = append(res[link], i)
	}
	return res
}

// 输出每个顶点代表的回文串.
func (pt *PalindromicTreeArray) GetPalindrome(pos int) []int {
	if pos == 0 {
		return []int{-1}
	}
	if pos == 1 {
		return []int{0}
	}
	var res []int
	// 在偶树/奇树中找到当前节点的回文串
	pt.outputDfs(0, pos, &res)
	pt.outputDfs(1, pos, &res)
	start := len(res) - 1
	if pt.Nodes[pos].Length&1 == 1 {
		start--
	}
	for i := start; i >= 0; i-- {
		res = append(res, res[i])
	}
	return res
}

// 回文树中的顶点个数.(包含两个奇偶虚拟顶点)
// 一个串的本质不同回文子串个数等于 Size()-2.
func (pt *PalindromicTreeArray) Size() int {
	return len(pt.Nodes)
}

// 返回pos位置的回文串顶点.
func (pt *PalindromicTreeArray) GetNode(pos int) *Node {
	return pt.Nodes[pos]
}

func (pt *PalindromicTreeArray) newNode(link, length int32) *Node {
	res := &Node{
		Next:      [SIGMA]int32{},
		Link:      link,
		Length:    length,
		deltaLink: -1,
	}
	for i := range res.Next {
		res.Next[i] = -1
	}
	return res
}

// 沿着失配指针找到第一个满足 x+s+x 是原串回文后缀的位置.
func (pt *PalindromicTreeArray) findPrevPalindrome(cur int32) int32 {
	pos := int32(len(pt.Ords) - 1)
	for {
		rev := pos - 1 - pt.Nodes[cur].Length
		// !插入当前字符的条件str[i]==str[i-len-1]
		if rev >= 0 && pt.Ords[rev] == pt.Ords[len(pt.Ords)-1] {
			break
		}
		cur = pt.Nodes[cur].Link
	}
	return cur
}

// 当前位置的回文串长度减去当前回文串的最长后缀回文串的长度.
func (pt *PalindromicTreeArray) diff(pos int32) int32 {
	if pt.Nodes[pos].Link <= 0 {
		return -1
	}
	return pt.Nodes[pos].Length - pt.Nodes[pt.Nodes[pos].Link].Length
}

func (pt *PalindromicTreeArray) outputDfs(cur, id int, res *[]int) bool {
	if cur == id {
		return true
	}
	for key, next := range pt.Nodes[cur].Next {
		if pt.outputDfs(int(next), id, res) {
			*res = append(*res, int(key))
			return true
		}
	}
	return false
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
