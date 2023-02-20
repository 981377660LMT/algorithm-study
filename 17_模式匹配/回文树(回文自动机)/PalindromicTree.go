// https://blog.csdn.net/Clove_unique/article/details/53750322
// https://ei1333.github.io/library/string/palindromic-tree.hpp
// https://oi-wiki.org/string/pam/
// 文字列 Sが与えられる. Parindromic Tree は Sに含まれるすべての回文を頂点とした木である.
// 長さが -1 と 0 の超頂点を用意する.
// 各頂点からは, その頂点に対応する回文の両端に同じ文字を 1文字加えてできる回文の頂点に辺を張ることで木を構成する.
// 特に長さ 1の回文は-1の超頂点から, 長さ 2の回文は 0の超頂点から辺を張ることになる.

// !さらに Suffix Link として, 頂点の回文に含まれる`最大の回文接尾辞`に対応する頂点に辺を張る.
// 例えば eertree からは ee, reer からは r に Suffix Link が張られることになる.
// Suffix Link からなる木を Suffix Link Tree と呼ぶことにする.

// Parindromic Tree は, 長さ -1の超頂点, 長さ 0の超頂点を根とした木および
// Suffix Link Tree の 3つの木構造を同時に管理するデータ構造とみなせる.

// !文字列に含まれる全てのユニークな回文の個数は超頂点 -1,0を除いた頂点数,
// !i番目の文字が最後尾となるような回文の個数は対応する頂点の Suffix Link Tree の深さと一致する.

// 应用
// 本质不同回文子串个数
//   等于回文树的状态数（排除奇根和偶根两个状态）。
// 回文子串出现次数
//   由于回文树的构造过程中，节点本身就是按照拓扑序插入，
//   因此只需要逆序枚举所有状态，将当前状态的出现次数加到其 fail 指针对应状态的出现次数上即可。
// 最小回文划分

// Usage:
// PalindromicTree(S): 文字列 Sに対応する Palindromic Tree を構築する.
// AddChar(x): 末尾に文字 x を追加する. 追加した文字を末尾とする最長回文接尾辞の頂点番号を返す.
// AddString(S): 末尾に文字列 S を追加する.
// GetFrequency(): 各頂点に対応する回文の出現回数を返す.
// GetPalindrome(pos): 頂点 pos に対応する回文を返す.
// GetNode(pos): 頂点 pos に対応するノードを返す.
// Size(): 頂点数を返す.

// 回文树里有两棵树，分别记录长度为奇数和偶数的回文串
// 0号表示回文串长度为偶数的树，1号表示回文串长度为奇数的树
// !每个节点代表一个回文串，记录转移c，如果在这个回文串前后都加上字符c形成的回文串原串的后缀,那么就添加儿子,否则就沿着fail指针往上找
// !每一个节点记录一个fail指针，指向这个回文串的最长后缀回文串
// 特殊地，0号点的fail指针指向1，非0、1号点并且后缀中不存在回文串的节点不指向它本身，而是指向0.
// !回文自动机的实质就是按顺序添加字符，每添加一个字符都要找出`以这个字符为后缀`的最长的回文串
// 例如当前字符为eertr,原来的指针指在rtr
// 现在追加e,会先判断树中是否存在 e+rtr+e ,如果不存在就会沿着fail指针一直往上找
// 即从 rtr -> r,然后看 是否存在 e+r+e ,如果不存在就继续往上找,知道找到一个存在的回文后缀(e).
// 最后将 eertr 指向这个位置e,也就是它的suffixLink.

package main

import "fmt"

func main() {
	pt := NewPalindromicTree("")

	fmt.Println(pt.Size() - 2)     // 本质不同回文子串个数(不超过O(n))
	fmt.Println(pt.GetFrequency()) // 每个顶点对应的回文串出现的次数
	for pos := 0; pos < pt.Size(); pos++ {
		fmt.Println(pt.GetPalindrome(pos))
	}

	s := "avaa"
	for i := 0; i < len(s); i++ {
		pos := pt.AddChar(s[i])
		node := pt.GetNode(pos)
		fmt.Println(node.Fail)
	}
}

type PalindromicTree struct {
	Chars []byte
	Nodes []*Node
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
	res.Nodes = append(res.Nodes, res.newNode(0, -1)) // 長さ -1 (奇数)
	res.Nodes = append(res.Nodes, res.newNode(0, 0))  // 長さ 0 (偶数)
	res.AddString(s)
	return res
}

// !添加一个字符,返回以这个字符为后缀的最长回文串的位置pos.
func (pt *PalindromicTree) AddChar(x byte) int {
	pos := len(pt.Chars)
	pt.Chars = append(pt.Chars, x)
	// 如果在这个回文串前后都加上字符c形成的回文串原串的后缀,那么就添加儿子,否则就沿着fail指针往上找
	cur := pt.findPrevPalindrome(pt.ptr)
	_, hasKey := pt.Nodes[cur].next[x]
	if !hasKey {
		pt.Nodes[cur].next[x] = len(pt.Nodes)
	}
	pt.ptr = pt.Nodes[cur].next[x]
	if !hasKey {
		pt.Nodes = append(pt.Nodes, pt.newNode(-1, pt.Nodes[cur].Len+2))
		if pt.Nodes[len(pt.Nodes)-1].Len == 1 {
			pt.Nodes[len(pt.Nodes)-1].Fail = 1
		} else {
			pt.Nodes[len(pt.Nodes)-1].Fail = pt.Nodes[pt.findPrevPalindrome(pt.Nodes[cur].Fail)].next[x]
		}

		if pt.diff(pt.ptr) == pt.diff(pt.Nodes[len(pt.Nodes)-1].Fail) {
			pt.Nodes[len(pt.Nodes)-1].deltaLink = pt.Nodes[pt.Nodes[len(pt.Nodes)-1].Fail].deltaLink
		} else {
			pt.Nodes[len(pt.Nodes)-1].deltaLink = pt.Nodes[len(pt.Nodes)-1].Fail
		}
	}

	pt.Nodes[pt.ptr].Indexes = append(pt.Nodes[pt.ptr].Indexes, pos)
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
func (pt *PalindromicTree) UpdateDp(init func(pos, start int), apply func(pos, pre int)) (update []int) {
	i := len(pt.Chars) - 1
	id := pt.ptr
	for pt.Nodes[id].Len > 0 {
		init(id, i+1-pt.Nodes[pt.Nodes[id].deltaLink].Len-pt.diff(id))
		if pt.Nodes[id].deltaLink != pt.Nodes[id].Fail {
			apply(id, pt.Nodes[id].Fail)
		}
		update = append(update, id)
		id = pt.Nodes[id].deltaLink
	}
	return
}

// 求出每个顶点代表的回文串出现的次数.
func (pt *PalindromicTree) GetFrequency() []int {
	res := make([]int, pt.Size())
	for i := pt.Size() - 1; i > 0; i-- {
		res[i] += len(pt.Nodes[i].Indexes)
		res[pt.Nodes[i].Fail] += res[i] // 长回文包含短回文
	}
	return res
}

// 以当前字符结尾的回文串个数.
func (pt *PalindromicTree) CountPalindromes() int {
	res := 0
	for i := 1; i < pt.Size(); i++ {
		res += len(pt.Nodes[i].Indexes)
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
	if pt.Nodes[pos].Len&1 == 1 {
		start--
	}
	for i := start; i >= 0; i-- {
		res = append(res, res[i])
	}
	return []string{string(res)}
}

// 回文树中的顶点个数.(包含两个奇偶虚拟顶点)
func (pt *PalindromicTree) Size() int {
	return len(pt.Nodes)
}

// 返回pos位置的回文串顶点.
func (pt *PalindromicTree) GetNode(pos int) *Node {
	return pt.Nodes[pos]
}

// 当前位置的回文串长度减去当前回文串的最长后缀回文串的长度.
func (pt *PalindromicTree) diff(pos int) int {
	if pt.Nodes[pos].Fail <= 0 {
		return -1
	}
	return pt.Nodes[pos].Len - pt.Nodes[pt.Nodes[pos].Fail].Len
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
		rev := pos - 1 - pt.Nodes[cur].Len
		if rev >= 0 && pt.Chars[rev] == pt.Chars[len(pt.Chars)-1] {
			break
		}
		cur = pt.Nodes[cur].Fail
	}
	return cur
}

func (pt *PalindromicTree) outputDfs(v, id int, res *[]byte) bool {
	if v == id {
		return true
	}
	for key, next := range pt.Nodes[v].next {
		if pt.outputDfs(next, id, res) {
			*res = append(*res, key)
			return true
		}
	}
	return false
}
