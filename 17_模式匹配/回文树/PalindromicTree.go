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
// - 本质不同回文子串个数
//   等于回文树的状态数（排除奇根和偶根两个状态）。
// - 所有回文子串的数量
//   由于回文树的构造过程中，节点本身就是按照拓扑序插入，
//   因此只需要逆序枚举所有状态，将当前状态的出现次数加到其 fail 指针对应状态的出现次数上即可。
// - 两个串的公共回文子串数量
//   同时dfs遍历其相同的节点，累加两边节点的cnt的乘积即可
// - 最小回文划分

// Usage:
// PalindromicTree(S): 文字列 Sに対応する Palindromic Tree を構築する.
// Add(x): 末尾に文字 x を追加する. 追加した文字を末尾とする最長回文接尾辞の頂点番号を返す.
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

// !类似AC自动机，将回文子串建立联系.
// 回文树的每个节点对应一个回文子串.
// !类似在trie上的插入，不过现在走过一条边的含义变为"在当前回文串的`两边`加上字符c".
// 回文树有两个根：奇根、偶根.
// 奇根的子节点表示回文子串长度为奇数，偶根的子节点表示回文子串长度为偶数.
// 奇根长度为-1，失配指针指向偶根；偶根长度为0，失配指针指向奇根.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yuki2606()
}

func demo() {
	pt := NewPalindromicTree()
	pt.AddString("eertree")

	fmt.Println(pt.Size() - 2)     // 本质不同回文子串个数(不超过O(n))
	fmt.Println(pt.GetFrequency()) // 每个顶点对应的回文串出现的次数
	for pos := 0; pos < pt.Size(); pos++ {
		res := pt.GetPalindrome(pos)
		fmt.Println(res)
	}
}

// https://yukicoder.me/problems/no/2606
// 给定一个字符串s.
// 向一个空字符x串插入字符，如果x为回文，则获得 x的长度*x在s中出现的次数 的分数.
// 求最终可能的最大分数.
// n<=2e5
func yuki2606() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	T := NewPalindromicTree()
	T.AddString(s)
	counter := T.GetFrequency()
	n := T.Size()
	dp := make([]int, n)
	for i := 2; i < n; i++ {
		node := T.GetNode(i)
		count := counter[i]
		length := int(node.Length)
		fail := node.Link
		dp[i] = max(dp[i], dp[fail]+count*length)
	}

	fmt.Fprintln(out, maxs(dp))
}

type Node struct {
	Next      map[int32]int32 // 子节点.
	Link      int32           // suffix link，指向当前回文串的最长真回文后缀的位置
	Length    int32           // 结点代表的回文串的长度
	Indexes   []int32         // 哪些最长真回文后缀为当前节点对应的回文
	deltaLink int32
}

type PalindromicTree struct {
	Ords    []int32
	Nodes   []*Node
	lastPos int32 // 当前的最长回文后缀
}

func NewPalindromicTree() *PalindromicTree {
	res := &PalindromicTree{}
	res.Nodes = append(res.Nodes, res.newNode(0, -1)) // 奇根，长为 -1
	res.Nodes = append(res.Nodes, res.newNode(0, 0))  // 偶根，长为 0
	return res
}

// !添加一个字符,返回以这个字符为后缀的最长回文串的位置pos.
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
		pt.Nodes = append(pt.Nodes, pt.newNode(-1, pt.Nodes[cur].Length+2))
		if pt.Nodes[len(pt.Nodes)-1].Length == 1 {
			pt.Nodes[len(pt.Nodes)-1].Link = 1
		} else {
			pt.Nodes[len(pt.Nodes)-1].Link = pt.Nodes[pt.findPrevPalindrome(pt.Nodes[cur].Link)].Next[x]
		}

		if pt.diff(pt.lastPos) == pt.diff(pt.Nodes[len(pt.Nodes)-1].Link) {
			pt.Nodes[len(pt.Nodes)-1].deltaLink = pt.Nodes[pt.Nodes[len(pt.Nodes)-1].Link].deltaLink
		} else {
			pt.Nodes[len(pt.Nodes)-1].deltaLink = pt.Nodes[len(pt.Nodes)-1].Link
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

// 在每次调用AddChar(x)之后使用,更新dp.
//   - init(pos, start): 初始化顶点pos的dp值,对应回文串s[start:].
//   - apply(pos, prePos): 用prePos(fail指针指向的位置)更新pos.
//     返回值: 本次更新的回文串的顶点.
func (pt *PalindromicTree) UpdateDp(init func(pos, start int), apply func(pos, pre int)) (updated []int) {
	i := int32(len(pt.Ords) - 1)
	id := pt.lastPos
	for pt.Nodes[id].Length > 0 {
		init(int(id), int(i+1-pt.Nodes[pt.Nodes[id].deltaLink].Length-pt.diff(id)))
		if pt.Nodes[id].deltaLink != pt.Nodes[id].Link {
			apply(int(id), int(pt.Nodes[id].Link))
		}
		updated = append(updated, int(id))
		id = pt.Nodes[id].deltaLink
	}
	return
}

// 求出每个顶点(每个后缀)代表的回文串出现的次数.
func (pt *PalindromicTree) GetFrequency() []int {
	res := make([]int, pt.Size())
	for i := pt.Size() - 1; i >= 1; i-- { // 除去根节点(奇根)
		res[i] += len(pt.Nodes[i].Indexes)
		res[pt.Nodes[i].Link] += res[i] // 长回文包含短回文
	}
	return res
}

// 当前字符的本质不同回文串个数.
func (pt *PalindromicTree) CountPalindromes() int {
	res := 0
	for i := 1; i < pt.Size(); i++ { // 除去根节点(奇根)
		res += (len(pt.Nodes[i].Indexes))
	}
	return res
}

// 输出每个顶点代表的回文串.
func (pt *PalindromicTree) GetPalindrome(pos int) []int {
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
		if rev >= 0 && pt.Ords[rev] == pt.Ords[len(pt.Ords)-1] {
			break
		}
		cur = pt.Nodes[cur].Link
	}
	return cur
}

// 当前位置的回文串长度减去当前回文串的最长后缀回文串的长度.
func (pt *PalindromicTree) diff(pos int32) int32 {
	if pt.Nodes[pos].Link <= 0 {
		return -1
	}
	return pt.Nodes[pos].Length - pt.Nodes[pt.Nodes[pos].Link].Length
}

func (pt *PalindromicTree) outputDfs(cur, id int, res *[]int) bool {
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

func maxs(nums []int) int {
	res := nums[0]
	for _, v := range nums {
		if v > res {
			res = v
		}
	}
	return res
}
