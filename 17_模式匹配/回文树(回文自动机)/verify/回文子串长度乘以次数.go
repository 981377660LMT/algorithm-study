// https://www.luogu.com.cn/problem/P5496
// !对所有不同的回文子串,求出(回文长度*回文出现次数)的最大值

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	tree := NewPalindromicTree(s)
	freq := tree.GetFrequency()
	res := 0
	for pos := 2; pos < len(freq); pos++ {
		node := tree.GetNode(pos)
		res = max(res, node.Len*freq[pos])
	}
	fmt.Fprintln(out, res)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type PalindromicTree struct {
	Chars []byte
	nodes []*Node
	ptr   int // 新添加一个字母后所形成的最长回文后缀表示的节点
}

type Node struct {
	Len       int          // 结点代表的回文串的长度
	Indexes   []int        // 以哪些索引结尾的最长回文后缀是当前回文串
	Fail      int          // 指向当前回文串的最长回文后缀的位置
	next      map[byte]int // 当前回文串前后都加上字符c形成的回文串
	deltaLink int          // 以当前回文串结尾的长度不同的回文串的位置(ertre -> e)?
}

func NewPalindromicTree(s string) *PalindromicTree {
	res := &PalindromicTree{}
	res.nodes = append(res.nodes, res.newNode(0, -1)) // 長さ -1 (奇数)
	res.nodes = append(res.nodes, res.newNode(0, 0))  // 長さ 0 (偶数)
	res.AddString(s)
	return res
}

// !添加一个字符,返回以这个字符为后缀的最长回文串的位置pos.
func (pt *PalindromicTree) AddChar(x byte) int {
	pos := len(pt.Chars)
	pt.Chars = append(pt.Chars, x)
	// 如果在这个回文串前后都加上字符c形成的回文串原串的后缀,那么就添加儿子,否则就沿着fail指针往上找
	cur := pt.findPrevPalindrome(pt.ptr)
	_, hasKey := pt.nodes[cur].next[x]
	if !hasKey {
		pt.nodes[cur].next[x] = len(pt.nodes)
	}
	pt.ptr = pt.nodes[cur].next[x]
	if !hasKey {
		pt.nodes = append(pt.nodes, pt.newNode(-1, pt.nodes[cur].Len+2))
		if pt.nodes[len(pt.nodes)-1].Len == 1 {
			pt.nodes[len(pt.nodes)-1].Fail = 1
		} else {
			pt.nodes[len(pt.nodes)-1].Fail = pt.nodes[pt.findPrevPalindrome(pt.nodes[cur].Fail)].next[x]
		}

		if pt.diff(pt.ptr) == pt.diff(pt.nodes[len(pt.nodes)-1].Fail) {
			pt.nodes[len(pt.nodes)-1].deltaLink = pt.nodes[pt.nodes[len(pt.nodes)-1].Fail].deltaLink
		} else {
			pt.nodes[len(pt.nodes)-1].deltaLink = pt.nodes[len(pt.nodes)-1].Fail
		}
	}

	pt.nodes[pt.ptr].Indexes = append(pt.nodes[pt.ptr].Indexes, pos)
	return pt.ptr
}

func (pt *PalindromicTree) AddString(s string) {
	if len(s) == 0 {
		return
	}
	for i := 0; i < len(s); i++ {
		pt.AddChar(s[i])
	}
}

// 在每次调用AddChar(x)之后使用,更新dp.
//  - init(pos, start): 初始化顶点pos的dp值,对应回文串s[start:].
//  - apply(pos, prePos): 用prePos(fail指针指向的位置)更新pos.
//  返回值: 本次更新的回文串的顶点.
func (pt *PalindromicTree) UpdateDp(init func(int, int), apply func(int, int)) (update []int) {
	i := len(pt.Chars) - 1
	id := pt.ptr
	for pt.nodes[id].Len > 0 {
		init(id, i+1-pt.nodes[pt.nodes[id].deltaLink].Len-pt.diff(id))
		if pt.nodes[id].deltaLink != pt.nodes[id].Fail {
			apply(id, pt.nodes[id].Fail)
		}
		update = append(update, id)
		id = pt.nodes[id].deltaLink
	}
	return
}

// 求出每个顶点代表的回文串出现的次数.
func (pt *PalindromicTree) GetFrequency() []int {
	res := make([]int, pt.Size())
	for i := pt.Size() - 1; i > 0; i-- {
		res[i] += len(pt.nodes[i].Indexes)
		res[pt.nodes[i].Fail] += res[i] // 长回文包含短回文
	}
	return res
}

// 以当前字符结尾的回文串个数.
func (pt *PalindromicTree) CountPalindromes() int {
	res := 0
	for i := 1; i < pt.Size(); i++ {
		res += len(pt.nodes[i].Indexes)
	}
	return res
}

// 输出每个顶点代表的回文串.
func (pt *PalindromicTree) GetPalindrome(pos int) []string {
	if pos == 0 {
		return []string{"-1"}
	}
	if pos == 1 {
		return []string{"0"}
	}
	res := []byte{}
	pt.outputDfs(0, pos, &res)
	pt.outputDfs(1, pos, &res)
	start := len(res) - 1
	if pt.nodes[pos].Len&1 == 1 {
		start--
	}
	for i := start; i >= 0; i-- {
		res = append(res, res[i])
	}
	return []string{string(res)}
}

// 回文树中的顶点个数.(包含两个奇偶虚拟顶点)
func (pt *PalindromicTree) Size() int {
	return len(pt.nodes)
}

// 返回pos位置的回文串顶点.
func (pt *PalindromicTree) GetNode(pos int) *Node {
	return pt.nodes[pos]
}

// 当前位置的回文串长度减去当前回文串的最长后缀回文串的长度.
func (pt *PalindromicTree) diff(pos int) int {
	if pt.nodes[pos].Fail <= 0 {
		return -1
	}
	return pt.nodes[pos].Len - pt.nodes[pt.nodes[pos].Fail].Len
}

func (pt *PalindromicTree) newNode(suf, pLen int) *Node {
	return &Node{
		next:      make(map[byte]int),
		Fail:      suf,
		Len:       pLen,
		deltaLink: -1,
	}
}

// 沿着失配指针找到第一个满足 x+s+x 是原串回文后缀的位置.
func (pt *PalindromicTree) findPrevPalindrome(cur int) int {
	pos := len(pt.Chars) - 1
	for {
		rev := pos - 1 - pt.nodes[cur].Len
		if rev >= 0 && pt.Chars[rev] == pt.Chars[len(pt.Chars)-1] {
			break
		}
		cur = pt.nodes[cur].Fail
	}
	return cur
}

func (pt *PalindromicTree) outputDfs(v, id int, res *[]byte) bool {
	if v == id {
		return true
	}
	for key, next := range pt.nodes[v].next {
		if pt.outputDfs(next, id, res) {
			*res = append(*res, key)
			return true
		}
	}
	return false
}
