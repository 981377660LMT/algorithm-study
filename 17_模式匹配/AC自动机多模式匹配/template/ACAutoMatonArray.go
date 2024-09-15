// https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
// https://www.cnblogs.com/alex-wei/p/Common_String_Theory_Theory_automaton_related.html
// https://zhuanlan.zhihu.com/p/533603249
// 类似 KMP 的失配数组，失配指针的含义为：
// 当前结点所表示字符串的最长真后缀 (border)，使得该后缀作为某个单词的前缀出现。
// !https://www.cnblogs.com/alex-wei/p/Common_String_Theory_Theory_automaton_related.html
// !- 子串 = 前缀的后缀
//
//	  Trie树（AC自动机）的祖先节点 = 前缀
//	  Fail树的祖先节点 = 后缀
//	  字符串x在字符串y中出现的次数 = `Fail树中x的子树`与`Trie树中y到根的路径`的交点个数
//	  一个单词在匹配串S中出现次数之和，等于它在S的所有前缀中作为后缀出现 的次数之和。
//
//	 - 单词her、say、she、shr、he构成的AC自动机。
//	   (i)表示的是link指针指向的节点。
//
//	                      ""
//	                    /    \
//	                   /      \
//	                  h(1)     s
//	                 /       /   \
//	               e(2)    a     h(1)
//								 /      /     /   \
//	              r      y    e(2)    r
//
// - notes:
//  !0.节点也被称为 状态。每个节点对应某个模式串的某个前缀。
//  1. 如果我们既知道前缀信息（trie），又知道后缀信息（fail），就可以做字符串匹配：
//     前缀的后缀就是子串，只要遍历到所有前缀，对每个前缀做「后缀匹配」，就完成了字符串匹配（统计子串出现次数）
//  2. 动态增删模式串，同时询问「查询所有模式串在文本串中的出现次数」，可以改为离线（先把所有模式串加入 AC 自动机）
//     当对一个节点增删 end 标记时，如果只对这一个节点修改，那么询问就需要遍历整个 fail 链，太慢了
//     换个思路：改为 end 标记会对它的 fail 子树全部 +1/-1，这样就可以在遍历文本串时单点询问了
//  3. 简单的AC自动机只能进行不重复统计，即统计字典中多少个字符串出现在目标串中。
//     如果要重复计数的话，即同一个字符串出现在不同位置要重复计数，则依靠fail指针可能很耗时，因为需要遍历fail指针。
//     可以使用树（fail树）上的前缀和来累计求解。如果涉及离线修改，则可以使用dfs序+树状数组来维护这个前缀和。
//  4. 字符串之间的包含关系:
//     !"子串是某个前缀的后缀", 遍历文本串的每个前缀，在fail树上找到对应最长后缀，然后在trie树上定位到对应串.
//     !- 前缀在trie树上,一个字符串的所有前缀对应根节点到该节点的路径.
//     !- 后缀在fail树上,结点p的子树对应字符的后缀均为结点p对应的字符串.
//  5. 当 Children[i][c] 不能匹配文本串 text 中的某个字符时，Children[Fail[i]][c] 即为下一个待匹配节点.
//  6. 多串问题 => 字典树/广义SAM
//  7. 对于fail转移p1->p2，p2既是p1的最长真后缀，也是p2[i]的一个前缀.
//  8. ACAM 接受且仅接受以给定词典中的某一个单词"结尾"的字符串.
//  !9. 失配指针fail[p]的含义位p所表示字符串s[p]的最长真后缀，使得该后缀作为某个单词的前缀出现。
//  !11.Move(pos,char)表示往状态pos后添加字符char，所得字符串的"最长的"与某个单词的"前缀"匹配的"后缀"所表示的状态。
//  !12.终止节点代表一个单词，或"以一个单词结尾"。所有终止节点 组成的集合对应着 DFA 的 接受状态集合。
//  !13.fail树性质
//     - 对于节点p及其对应字符串pattern[p]，对于其子树内部所有节点q，都有pattern[q]是pattern[p]的后缀。
//     反之也成立，即q是p的后缀<=>p在q的子树内。(这里与后缀自动机fail树的性质相同)
//     - 若p是终止节点，则p的子树全部都是终止节点(border关系.)
//     - 一个单词p在匹配串s中出现次数之和，等于它在s的"所有前缀"中作为"后缀"出现的次数之和，等于 fail 树从p到根节点上单词节点的数量
//  !14.通常带修链求和要用到树剖，但查询具有特殊性质：一个端点是根。
//     !链加+单点查询 可以变为 单点加+子树求和。
//     !单点修改+链求和，可以转变为 子树修改+单点查询。(TODO:封装BIT方法.)
//     只要包含一个单点操作，一个链操作，均可以将链操作转化为子树操作，即可将时间复杂度更大的树剖 BIT 换成普通 BIT。
//  !15. 两种处理方式:
//     - 模式串信息dp预处理 + 文本串单点查询：得到模式串在trie树上的位置，然后将结点信息沿着fail树向下dp传递，文本串查询时只需要单点查询;
//     - 文本串跳fail查询：不使用dp预处理模式串信息，文本串查询时每次跳fail查询(注意如果求的是树链并，需要visited去重.)
//  !16. 模式串作为后缀，文本串作为前缀.
//  !17. AC自动机擅长多串模式匹配，还可以处理字符串权值带修的问题，而广义SAM就不擅长。
//       !AC自动机擅长处理"多个串在某个串中"的出现次数，而广义SAM擅长处理"某个串在多个串中"的出现次数.
//       eg: CF587F-Duff is Mad: 给定n个字符串，q次询问[start:end)中的字符串在s[index]中出现次数之和.
//           CF547E-Mike and Friends: 给定n个字符串，q次询问s[index]在[start:end)中的字符串出现次数之和.

//  !18. 注意trie树和fail树共用结点，类似sam的fail树和dag共用结点.
//       在trie树/sam的dag往下走，相当于往后加字符，而在fail树往下走，相当于往前加字符.
//  !19. !AC自动机接受所有模式串的前缀.
//
// - applications:
//  1. "不能出现若干单词" 的字符串 计数或最优化 问题 => ac自动机上dp: 一般都是dfs(index,pos):长度为index的字符串，当前trie状态为pos.
//     枚举26种字符(字符集)转移.
//  2. 自定义结点信息题：用WordPos初始化状态，用SuffixLink转移.
//  !3. 查询x在y中的出现次数 => y到Trie树的根的每个结点上打标记(前缀)，查询fail树上x的子树内有标记的节点个数
//                           => 离线查询：trie树上dfs标记前缀, 树状数组查询fail树子树和.
//  !4. 每个模式串查询在文本串中出现的次数 => 文本串在每个结点上打标记(前缀)，查询fail树上模式串的子树内有标记的节点个数(子树和)。
//  !5. ac自动机+树状数组应用总结：
//    - P3966 单词: 对模式串的每个串，求他在所有模式串中的出现次数.
//             文本串在 AC 自动机上每经过一个节点就将节点权值增加1(标记所有前缀)，
//             每个单词的出现次数为它在 fail 树上的子树节点权值和(作为多少个前缀的后缀出现.)
//    - P2414 [NOI2011] 阿狸的打字机：离线"多次"查询一个串在另一个串的出现次数.
//             等价于:fail树中x的子树(对应一些更长的后缀)与trie树中y到根节点的路径(对应一些更短的前缀)的公共结点数.
//						 离线查询，将所有询问保存到y上，在 Trie树 上 dfs+回溯 即可.
//             !"链加、单点查询" 转化为 "单点加、子树和查询".
//    - P5840 [COCI2015] Divljak：对每个新加入的串的前缀求树链并.
//    - CF163E e-Government：动态求所有模式串在文本串中出现次数之和
//             !"单点修改+链上求和" 转化为 "子树修改+单点求和";
//    - CF547E Mike and Friends ：查询s[index]在[start:end)中的字符串出现次数之和.
//						 将区间查询转换为两个前缀的差，扫描线加入模式串，并对模式串每个前缀，标记每个后缀.
//             查询时单点查询即可.
//    - CF587F-Duff is Mad: 查询[start:end)中的字符串在s[index]中出现次数之和.
//             按照查询的文本串的长度进行根号分治.
//             !文本串为短串时，将所有模式串中的短串按照扫描线的顺序加入.扫描线+前缀和+树状数组处理.
//             对每个(短)模式串，加入时fail树子树和+1, 文本串查询时累加其前缀结点之和.类似mike and friends中的技巧.
//             !文本串为长串时，这样的串不超过sqrt(n)个.
//             转化为求出每个模式串在文本串中出现的次数(P3966 单词).
//             文本串每个前缀+1，dfs求出子树和即可.
//  !6. ac自动机线性dp linkWord

// !注意不能直接用link跳fail，会被卡到O(n^2).
// 某个串很长时，一直指向的是这个长串的局部.
// 需要使用linkWord.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {

	// SeparateString()

	// P3041()
	// P3121()
	// P3193()
	// P3796()
	// P3808()
	// P3966()
	// P4052()

	// P5357()
	// P5357_link()
	P5357_dp()

	// P7456()

	// CF808G()

}

func demo() {
	acm := NewACAutoMatonArray(26, 97)
	acm.AddString("abc")
	acm.AddString("bc")
	acm.BuildSuffixLink(false)
	fmt.Println(acm.GetCounter())
}

const INF32 int32 = 1e9 + 10

// 一个文本串，可以将一个字符改成*，
// 保证都是小写字母组成。问最少改多少次，可以完全不包含若干模式串。
// !AC自动机(多模式串匹配)
// https://atcoder.jp/contests/abc268/submissions/34752897
func SeparateString() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	var n int
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	acm := NewACAutoMatonArray(26, 97)
	for _, word := range words {
		acm.AddString(word)
	}
	acm.BuildSuffixLink(false)

	accept := acm.GetCounter()
	pos := int32(0)
	res := 0
	for _, char := range s {
		pos = acm.Move(pos, char)
		if accept[pos] > 0 {
			res++
			pos = 0
		}
	}

	fmt.Fprintln(out, res)
}

// P3041 [USACO12JAN] Video Game G
// https://www.luogu.com.cn/problem/P3041
// 给出n个字典单词，问长度为k的字符串最多可以包含多少个字典单词。
// n<=20,k<=1000,len(word[i])<=15
func P3041() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int32
	fmt.Fscan(in, &n, &k)
	words := make([]string, n)
	acm := NewACAutoMatonArray(26, 'A')
	for i := range words {
		var s string
		fmt.Fscan(in, &s)
		words[i] = s
		acm.AddString(s)
	}
	acm.BuildSuffixLink(true)

	size := acm.Size()
	counter := acm.GetCounter()
	memo := make([][]int32, k)
	for i := range memo {
		row := make([]int32, size)
		for j := range row {
			row[j] = -1
			memo[i] = row
		}
	}
	var dfs func(index, pos int32) int32
	dfs = func(index, pos int32) int32 {
		if index == k {
			return 0
		}
		if tmp := memo[index][pos]; tmp != -1 {
			return tmp
		}
		res := int32(0)
		for v := int32(65); v < 65+26; v++ {
			nextPos := acm.Move(pos, v)
			res = max32(res, dfs(index+1, nextPos)+counter[nextPos])
		}
		memo[index][pos] = res
		return res
	}

	res := dfs(0, 0)
	fmt.Println(res)
}

// P3121 [USACO15FEB] Censoring G
// https://www.luogu.com.cn/problem/P3121
// 在longer中不断删除toRemove中的字符(按顺序遍历toRemove查找到要删除的)，求剩下的字符串.
// 注意删除一个单词后可能会导致 s 中出现另一个toRemove中的单词.
// !保证toRemove中没有包含关系(子串关系)，即删除时保证只会删除一个单词.
func P3121() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var longer string
	fmt.Fscan(in, &longer)

	acm := NewACAutoMatonArray(26, 97)
	var n int
	fmt.Fscan(in, &n)
	toRemove := make([]string, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		toRemove[i] = s
		acm.AddString(s)
	}
	acm.BuildSuffixLink(true)

	posLen := make([]int, acm.Size()) // 每个状态对应的字符长度
	for i, p := range acm.WordPos {
		posLen[p] = len(toRemove[i])
	}

	pos := int32(0)
	stack := make([]int32, 0, len(longer))
	posRecord := make([]int32, len(longer))
	for i := range longer {
		v := int32(longer[i])
		pos = acm.Move(pos, v)
		posRecord[i] = pos
		stack = append(stack, int32(i))
		if wordLen := posLen[pos]; wordLen > 0 { // 由于模式串不存在包含关系，所以不会出现多个模式串同时匹配的情况
			stack = stack[:len(stack)-wordLen]
			if len(stack) > 0 {
				pos = posRecord[stack[len(stack)-1]]
			} else {
				pos = 0
			}
		}
	}

	res := make([]byte, 0, len(stack))
	for _, v := range stack {
		res = append(res, longer[v])
	}
	fmt.Fprintln(out, string(res))
}

// P3193 [HNOI2008] GT考试 (KMP+矩阵快速幂dp)
// 给定一些长度之和为m的字符串evil.
// 求有多少种长度为n的数字串不包含evil中的任意一个字符串.
// n<=1e9,m<=100.
//
// dp[i][j] 表示长度为i的准考证和A匹配到了第j位的方案数.
// O(m^3logn)
// 预处理出转移的邻接矩阵，等价于有向图转移n次的路径数.
func P3193() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, MOD int
	fmt.Fscan(in, &n, &m, &MOD)
	wordCount := 1
	evil := make([]string, wordCount)
	for i := 0; i < wordCount; i++ {
		fmt.Fscan(in, &evil[i])
	}
	acm := NewACAutoMatonArray(10, 48) // '0'-'9'
	for _, word := range evil {
		acm.AddString(word)
	}
	acm.BuildSuffixLink(true)
	counter := acm.GetCounter()

	size := acm.Size()
	T := make([][]int, size)
	for i := range T {
		T[i] = make([]int, size)
	}

	for pos := int32(0); pos < size; pos++ {
		if counter[pos] == 0 {
			for char := '0'; char <= '9'; char++ {
				nextPos := acm.Move(pos, char)
				if counter[nextPos] == 0 {
					T[pos][nextPos]++
				}
			}
		}
	}

	T = MatPow(T, n, MOD)
	res := 0
	for _, v := range T[0] {
		res = (res + v) % MOD
	}
	fmt.Fprintln(out, res)
}

// P3796 AC 自动机（简单版 II）
// https://www.luogu.com.cn/problem/P3796
// 有 N 个由小写字母组成的模式串以及一个文本串 T。
// 每个模式串可能会在文本串中出现多次。
// 你需要找出哪些模式串在文本串 T 中出现的次数最多。
func P3796() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(patterns []string, text string) (maxCount int32, ids []int32) {
		acm := NewACAutoMatonArray(26, 97)
		for _, pattern := range patterns {
			acm.AddString(pattern)
		}
		acm.BuildSuffixLink(false)

		indexes := acm.GetIndexes()
		hit := make([]int32, len(patterns))
		pos := int32(0)
		for _, char := range text {
			pos = acm.Move(pos, char)
			for _, id := range indexes[pos] {
				hit[id]++
				if hit[id] > maxCount {
					maxCount = hit[id]
					ids = []int32{id}
				} else if hit[id] == maxCount {
					ids = append(ids, id)
				}
			}
		}

		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		return
	}

	for {
		var n int32
		fmt.Fscan(in, &n)
		if n == 0 {
			break
		}
		patterns := make([]string, n)
		for i := int32(0); i < n; i++ {
			fmt.Fscan(in, &patterns[i])
		}
		var text string
		fmt.Fscan(in, &text)

		maxCount, ids := solve(patterns, text)
		fmt.Fprintln(out, maxCount)
		for _, id := range ids {
			fmt.Fprintln(out, patterns[id])
		}
	}
}

// P3808 AC 自动机（简单版）
// https://www.luogu.com.cn/problem/P3808
// 给定 n 个模式串 和一个文本串，求有多少个不同的模式串在文本串里出现过。
// 两个模式串不同当且仅当他们编号不同。
//
// !相同编号的串多次出现仅算一次，因此题目相当于求：
// 文本串在模式串建出的ACAM上匹配时经过的所有节点到根的"路径的并"上单词节点的个数。
// !上跳fail并用数组标记是否访问过(queryChain).
func P3808() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &words[i])
	}
	var s string
	fmt.Fscan(in, &s)

	acm := NewACAutoMatonArray(26, 97)
	for _, word := range words {
		acm.AddString(word)
	}
	acm.BuildSuffixLink(false)

	size := acm.Size()
	endCount := make([]int32, size)
	for _, pos := range acm.WordPos {
		endCount[pos]++
	}

	visited := make([]bool, size)
	queryChain := func(pos int32) int32 {
		res := int32(0)
		for pos > 0 {
			if visited[pos] {
				break
			}
			visited[pos] = true
			res += endCount[pos]
			pos = acm.link[pos]
		}
		return res
	}

	res := 0
	pos := int32(0)
	for _, char := range s {
		pos = acm.Move(pos, char)
		res += int(queryChain(pos))
	}
	fmt.Fprintln(out, res)
}

// https://www.luogu.com.cn/problem/P3966
// 一篇论文是由许多单词组成但小张发现一个单词会在论文中出现很多次。
// 他想知道每个单词分别在论文中出现了多少次。
//
// 同P5357
// 文本串(这里为所有模式串)在 AC 自动机上每经过一个节点就将其权值增加1，
// !则每个单词在中的出现次数即 在 fail 树上的子树节点权值和(子树和)。
func P3966() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	acm := NewACAutoMatonArray(26, 97)
	for _, word := range words {
		acm.AddString(word)
	}
	acm.BuildSuffixLink(false)

	tree := acm.BuildFailTree()
	values := make([]int, acm.Size())
	for _, word := range words {
		pos := int32(0)
		for _, char := range word {
			pos = acm.Move(pos, char)
			values[pos]++
		}
	}
	subValues := make([]int, acm.Size())
	var dfs func(int32)
	dfs = func(pos int32) {
		subValues[pos] = values[pos]
		for _, child := range tree[pos] {
			dfs(child)
			subValues[pos] += subValues[child]
		}
	}
	dfs(0)

	for _, pos := range acm.WordPos {
		fmt.Fprintln(out, subValues[pos])
	}

}

// P4052 [JSOI2007] 文本生成器
// https://www.luogu.com.cn/problem/P4052
// !给定一些模式串，求长度为m的所有文本串的个数，且该文本串至少包括一个模式串，答案对10007取模
// n<=60,targetLen<=100,len(word[i])<=100.
// 合法的情况总数并不好求，考虑求出不合法的情况（即不存在一个子串等于模式串），最后用总数26^m减去即可。
func P4052() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e4 + 7

	acm := NewACAutoMatonArray(26, 'A')

	var n, targetLen int32
	fmt.Fscan(in, &n, &targetLen)
	words := make([]string, n)
	lengthSum := 0
	for i := int32(0); i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		words[i] = s
		acm.AddString(s)
		lengthSum += len(s)
	}
	acm.BuildSuffixLink(true)

	size := acm.Size()
	counter := acm.GetCounter()

	memo := make([][]int, targetLen)
	for i := range memo {
		row := make([]int, size)
		for j := range row {
			row[j] = -1
			memo[i] = row
		}
	}
	var dfs func(index, pos int32) int
	dfs = func(index, pos int32) int {
		if index == targetLen {
			return 1
		}
		if tmp := memo[index][pos]; tmp != -1 {
			return tmp
		}
		res := 0
		for v := int32(65); v < 65+26; v++ {
			nextPos := acm.Move(pos, v)
			if counter[nextPos] > 0 {
				continue
			}
			res += dfs(index+1, nextPos)
			res %= MOD
		}

		memo[index][pos] = res
		return res
	}
	bad := dfs(0, 0)

	qpow := func(a, b int) int {
		res := 1
		for b > 0 {
			if b&1 == 1 {
				res = res * a % MOD
			}
			a = a * a % MOD
			b >>= 1
		}
		return res
	}

	res := (qpow(26, int(targetLen)) - bad) % MOD
	if res < 0 {
		res += MOD
	}

	fmt.Println(res)
}

// P5357 【模板】AC 自动机（二次加强版）
// https://www.luogu.com.cn/problem/P5357
// 分别求出每个模式串在文本串中出现的次数。
//
// 文本串在 AC 自动机上每经过一个节点就将其权值增加1，
// !则每个单词在中的出现次数即 在 fail 树上的子树节点权值和(子树和)。
func P5357() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &words[i])
	}
	var s string
	fmt.Fscan(in, &s)

	acm := NewACAutoMatonArray(26, 'a')
	for _, word := range words {
		acm.AddString(word)
	}
	acm.BuildSuffixLink(false)

	values := make([]int, acm.Size())
	pos := int32(0)
	for _, char := range s {
		pos = acm.Move(pos, char)
		values[pos]++
	}
	failTree := acm.BuildFailTree()
	subValues := make([]int, acm.Size())
	var dfs func(int32)
	dfs = func(pos int32) {
		subValues[pos] = values[pos]
		for _, child := range failTree[pos] {
			dfs(child)
			subValues[pos] += subValues[child]
		}
	}
	dfs(0)

	for _, pos := range acm.WordPos {
		fmt.Fprintln(out, subValues[pos])
	}
}

func P5357_link() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &words[i])
	}
	var s string
	fmt.Fscan(in, &s)

	acm := NewACAutoMatonArray(26, 'a')
	for _, word := range words {
		acm.AddString(word)
	}
	acm.BuildSuffixLink(true)

	wordIndexes := make([][]int32, acm.Size())
	for i, pos := range acm.WordPos {
		wordIndexes[pos] = append(wordIndexes[pos], int32(i))
	}

	res := make([]int32, n)
	pos := int32(0)
	for _, char := range s {
		pos = acm.Move(pos, char)
		for cur := pos; cur > 0; cur = acm.LinkWord(cur) {
			for _, idx := range wordIndexes[cur] {
				res[idx]++
			}
		}
	}

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// !推荐做法.
func P5357_dp() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &words[i])
	}
	var s string
	fmt.Fscan(in, &s)

	acm := NewACAutoMatonArray(26, 'a')
	for _, word := range words {
		acm.AddString(word)
	}
	acm.BuildSuffixLink(true)

	counter := make([]int32, acm.Size())
	pos := int32(0)
	for _, char := range s {
		pos = acm.Move(pos, char)
		counter[pos]++
	}
	acm.DpReverse(func(link, cur int32) { counter[link] += counter[cur] })

	for _, v := range acm.WordPos {
		fmt.Fprintln(out, counter[v])
	}
}

// P7456 [CERC2018] The ABCD Murderer (结合线段树优化DP)
// n 个模式串和一个文本串，问要拼成文本串最少需要几个模式串。其中模式串可重叠。
//
// 求出对于每个位置 i，以s[i]结尾的最长单词的长度maxLen 。
// dp[i]=min(dp[j]), i-maxLen[i]<=j<i
func P7456() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var count int32
	fmt.Fscan(in, &count)
	var text string
	fmt.Fscan(in, &text)
	patterns := make([]string, count)
	for i := range patterns {
		fmt.Fscan(in, &patterns[i])
	}

	acm := NewACAutoMatonArray(26, 97)
	for _, s := range patterns {
		acm.AddString(s)
	}
	acm.BuildSuffixLink(false)

	maxBorder := make([]int32, acm.Size()) // 每个状态(前缀)的最长border
	for i, pos := range acm.WordPos {
		maxBorder[pos] = max32(maxBorder[pos], int32(len(patterns[i])))
	}
	acm.Dp(func(from, to int32) { maxBorder[to] = max32(maxBorder[to], maxBorder[from]) })

	n := int32(len(text))
	dp := NewPointSetRangeMin(n+1, func(i int32) E { return INF32 })
	dp.Set(0, 0)
	pos := int32(0)
	for i := int32(1); i <= n; i++ {
		char := int32(text[i-1])
		pos = acm.Move(pos, char)
		preMin := dp.Query(i-maxBorder[pos], i)
		dp.Set(i, preMin+1)
	}

	res := dp.Get(n)
	if res <= n {
		fmt.Fprintln(out, res)
	} else {
		fmt.Fprintln(out, -1)
	}
}

// Anthem of Berland (AC自动机+dp)
// https://www.luogu.com.cn/problem/CF808G
// 给定s串和t串，其中s串包含小写字母和问号，t串只包含小写字母。
// 假设有k个问号.你需要把每个问号替换成一个小写字母，一共有26^k种方案。
// 求替换后t在s中出现的次数的最大值.
//
// len(s)*len(t)<=1e7
//
// dp[i][pos]表示s串的前i个字符，AC自动机上的状态为pos时，匹配次数的最大值.
// t串还可以换成多个模式串.
func CF808G() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s, t string
	fmt.Fscan(in, &s, &t)

	words := []string{t}
	acm := NewACAutoMatonArray(26, 97)
	for _, w := range words {
		acm.AddString(w)
	}
	acm.BuildSuffixLink(true)
	counter := acm.GetCounter()

	size := acm.Size()
	initDp := func() []int {
		res := make([]int, size)
		for i := range res {
			res[i] = -1
		}
		return res
	}

	dp := initDp()
	dp[0] = 0
	for _, char := range s {
		ndp := initDp()
		for input := 'a'; input <= 'z'; input++ {
			// 当前字符不合法，不能作为转移
			if char != '?' && char != input {
				continue
			}
			for pos := int32(0); pos < size; pos++ {
				if dp[pos] == -1 {
					continue
				}
				nextPos := acm.Move(pos, input)
				ndp[nextPos] = max(ndp[nextPos], dp[pos]+int(counter[nextPos]))
			}
		}
		dp = ndp
	}

	res := 0
	for _, v := range dp {
		res = max(res, v)
	}
	fmt.Fprintln(out, res)
}

// 1032. 字符流
// https://leetcode.cn/problems/stream-of-characters/description/
type StreamChecker struct {
	*ACAutoMatonArray
	pos   int32 // state in AhoCorasick
	count []int32
}

// words = ["abc", "xyz"] 且字符流中逐个依次加入 4 个字符 'a'、'x'、'y' 和 'z' ，
// 你所设计的算法应当可以检测到 "axyz" 的后缀 "xyz" 与 words 中的字符串 "xyz" 匹配。
func Constructor(words []string) StreamChecker {
	acm := NewACAutoMatonArray(26, 'a')
	for _, word := range words {
		acm.AddString(word)
	}
	acm.BuildSuffixLink(false)
	return StreamChecker{ACAutoMatonArray: acm, count: acm.GetCounter()}
}

// 从字符流中接收一个新字符，如果字符流中的任一非空后缀能匹配 words 中的某一字符串，返回 true ；否则，返回 false。
func (this *StreamChecker) Query(letter byte) bool {
	this.pos = this.Move(this.pos, int32(letter))
	return this.count[this.pos] > 0
}

// 2781. 最长合法子字符串的长度
// https://leetcode.cn/problems/length-of-the-longest-valid-substring/
// 给你一个字符串 word 和一个字符串数组 forbidden 。
// 如果一个字符串不包含 forbidden 中的任何字符串，我们称这个字符串是 合法 的。
// 请你返回字符串 word 的一个 最长合法子字符串 的长度。
// 子字符串 指的是一个字符串中一段连续的字符，它可以为空。
//
// 1 <= word.length <= 1e5
// word 只包含小写英文字母。
// 1 <= forbidden.length <= 1e5
// !1 <= forbidden[i].length <= 1e5
// !sum(len(forbidden)) <= 1e7
// forbidden[i] 只包含小写英文字母。
//
// 思路:
// !子串是前缀的后缀.
// 类似字符流, 需要处理出每个位置为结束字符的包含至少一个模式串的`最短后缀`.
// !那么此时左端点就对应这个位置+1
func longestValidSubstring(word string, forbidden []string) int {
	acm := NewACAutoMatonArray(26, 97)
	for _, w := range forbidden {
		acm.AddString(w)
	}
	acm.BuildSuffixLink(false)

	minLen := make([]int32, acm.Size()) // 每个状态(前缀)匹配的最短单词
	for i := range minLen {
		minLen[i] = INF32
	}
	for i, pos := range acm.WordPos {
		minLen[pos] = min32(minLen[pos], int32(len(forbidden[i])))
	}
	acm.Dp(func(from, to int32) { minLen[to] = min32(minLen[to], minLen[from]) })

	res, left, pos := int32(0), int32(0), int32(0)
	for right, char := range word {
		pos = acm.Move(pos, char)
		left = max32(left, int32(right)-minLen[pos]+2)
		res = max32(res, int32(right)-left+1)
	}
	return int(res)
}

// https://leetcode.cn/problems/multi-search-lcci/
// 给定一个较长字符串big和一个包含较短字符串的数组smalls，
// 设计一个方法，根据smalls中的每一个较短字符串，对big进行搜索。
// !输出smalls中的字符串在big里出现的所有位置positions，
// 其中positions[i]为smalls[i]出现的所有位置。
func multiSearch(big string, smalls []string) [][]int {
	acm := NewACAutoMatonArray(26, 'a')
	for _, s := range smalls {
		acm.AddString(s)
	}
	acm.BuildSuffixLink(false)

	pos := int32(0)
	res := make([][]int, len(smalls))
	indexes := acm.GetIndexes()
	for i := 0; i < len(big); i++ {
		pos = acm.Move(pos, int32(big[i]))
		for _, j := range indexes[pos] {
			res[j] = append(res[j], i-len(smalls[j])+1) // !i-len(smalls[j])+1: 模式串在big中的起始位置
		}
	}
	return res
}

// 不调用 BuildSuffixLink 就是Trie，调用 BuildSuffixLink 就是AC自动机.
// 每个状态对应Trie中的一个结点，也对应一个前缀.
type ACAutoMatonArray struct {
	WordPos            []int32   // WordPos[i] 表示加入的第i个模式串对应的节点编号(单词结点).
	Parent             []int32   // parent[v] 表示节点v的父节点.
	Children           [][]int32 // children[v][c] 表示节点v通过字符c转移到的节点.
	BfsOrder           []int32   // 结点的拓扑序,0表示虚拟节点.
	Depth              []int32   // !每个节点的深度.也就是对应的模式串前缀的长度.
	link               []int32   // 又叫fail.指向当前trie节点(对应一个前缀)的最长真后缀对应结点，例如"bc"是"abc"的最长真后缀.
	linkWord           []int32
	sigma              int32 // 字符集大小.
	offset             int32 // 字符集的偏移量.
	needUpdateChildren bool  // 是否需要更新children数组.
}

func NewACAutoMatonArray(sigma, offset int32) *ACAutoMatonArray {
	res := &ACAutoMatonArray{sigma: sigma, offset: offset}
	res.Clear()
	return res
}

// 添加一个字符串，返回最后一个字符对应的节点编号.
func (trie *ACAutoMatonArray) AddString(str string) int32 {
	if len(str) == 0 {
		return 0
	}
	pos := int32(0)
	for _, s := range str {
		ord := s - trie.offset
		if trie.Children[pos][ord] == -1 {
			trie.Children[pos][ord] = trie.newNode2(pos, ord)
		}
		pos = trie.Children[pos][ord]
	}
	trie.WordPos = append(trie.WordPos, pos)
	return pos
}

// 在pos位置添加一个字符，返回新的节点编号.
func (trie *ACAutoMatonArray) AddChar(pos, ord int32) int32 {
	ord -= trie.offset
	if trie.Children[pos][ord] != -1 {
		return trie.Children[pos][ord]
	}
	trie.Children[pos][ord] = trie.newNode2(pos, ord)
	return trie.Children[pos][ord]
}

// pos: DFA的状态集, ord: DFA的字符集
func (trie *ACAutoMatonArray) Move(pos, ord int32) int32 {
	ord -= trie.offset
	if trie.needUpdateChildren {
		return trie.Children[pos][ord]
	}
	for {
		nexts := trie.Children[pos]
		if nexts[ord] != -1 {
			return nexts[ord]
		}
		if pos == 0 {
			return 0
		}
		pos = trie.link[pos]
	}
}

// 自动机中的节点(状态)数量，包括虚拟节点0.
func (trie *ACAutoMatonArray) Size() int32 {
	return int32(len(trie.Children))
}

func (trie *ACAutoMatonArray) Empty() bool {
	return len(trie.Children) == 1
}

// 构建后缀链接(失配指针).
// needUpdateChildren 表示是否需要更新children数组(连接trie图).
//
// !move调用较少时，设置为false更快.
func (trie *ACAutoMatonArray) BuildSuffixLink(needUpdateChildren bool) {
	trie.needUpdateChildren = needUpdateChildren
	trie.link = make([]int32, len(trie.Children))
	for i := range trie.link {
		trie.link[i] = -1
	}
	trie.BfsOrder = make([]int32, len(trie.Children))
	head, tail := 0, 0
	trie.BfsOrder[tail] = 0
	tail++
	for head < tail {
		v := trie.BfsOrder[head]
		head++
		for i, next := range trie.Children[v] {
			if next == -1 {
				continue
			}
			trie.BfsOrder[tail] = next
			tail++
			f := trie.link[v]
			for f != -1 && trie.Children[f][i] == -1 {
				f = trie.link[f]
			}
			trie.link[next] = f
			if f == -1 {
				trie.link[next] = 0
			} else {
				trie.link[next] = trie.Children[f][i]
			}
		}
	}
	if !needUpdateChildren {
		return
	}
	for _, v := range trie.BfsOrder {
		for i, next := range trie.Children[v] {
			if next == -1 {
				f := trie.link[v]
				if f == -1 {
					trie.Children[v][i] = 0
				} else {
					trie.Children[v][i] = trie.Children[f][i]
				}
			}
		}
	}
}

// !对当前文本串后缀，找到每个模式串单词匹配的最大前缀.
// `linkWord`指向当前节点的最长后缀对应的节点.
// 区别于`_link`,`linkWord`指向的节点对应的单词不会重复.
// 即不会出现`_link`指向某个长串局部的恶化情况.
func (trie *ACAutoMatonArray) LinkWord(pos int32) int32 {
	if len(trie.linkWord) == 0 {
		hasWord := make([]bool, len(trie.Children))
		for _, p := range trie.WordPos {
			hasWord[p] = true
		}
		trie.linkWord = make([]int32, len(trie.Children))
		link, linkWord := trie.link, trie.linkWord
		for _, v := range trie.BfsOrder {
			if v != 0 {
				p := link[v]
				if hasWord[p] {
					linkWord[v] = p
				} else {
					linkWord[v] = linkWord[p]
				}
			}
		}
	}
	return trie.linkWord[pos]
}

func (trie *ACAutoMatonArray) Clear() {
	trie.WordPos = trie.WordPos[:0]
	trie.Parent = trie.Parent[:0]
	trie.Depth = trie.Depth[:0]
	trie.Children = trie.Children[:0]
	trie.link = trie.link[:0]
	trie.linkWord = trie.linkWord[:0]
	trie.BfsOrder = trie.BfsOrder[:0]
	trie.newNode()
}

// 获取每个状态包含的模式串的个数.
func (trie *ACAutoMatonArray) GetCounter() []int32 {
	counter := make([]int32, len(trie.Children))
	for _, pos := range trie.WordPos {
		counter[pos]++
	}
	for _, v := range trie.BfsOrder {
		if v != 0 {
			counter[v] += counter[trie.link[v]]
		}
	}
	return counter
}

// 获取每个状态包含的模式串的索引.(模式串长度和较小时使用)
// fail指针每次命中，都至少有一个比指针深度更长的单词出现，因此每个位置最坏情况下不超过O(sqrt(n))次命中
// O(n*sqrt(n))
// TODO: roaring bitmaps 优化空间复杂度.
func (trie *ACAutoMatonArray) GetIndexes() [][]int32 {
	res := make([][]int32, len(trie.Children))
	for i, pos := range trie.WordPos {
		res[pos] = append(res[pos], int32(i))
	}
	for _, v := range trie.BfsOrder {
		if v != 0 {
			from, to := trie.link[v], v
			arr1, arr2 := res[from], res[to]
			arr3 := make([]int32, len(arr1)+len(arr2))
			i, j, k := 0, 0, 0
			for i < len(arr1) && j < len(arr2) {
				if arr1[i] < arr2[j] {
					arr3[k] = arr1[i]
					i++
				} else if arr1[i] > arr2[j] {
					arr3[k] = arr2[j]
					j++
				} else {
					arr3[k] = arr1[i]
					i++
					j++
				}
				k++
			}
			copy(arr3[k:], arr1[i:])
			k += len(arr1) - i
			copy(arr3[k:], arr2[j:])
			k += len(arr2) - j
			arr3 = arr3[:k:k]
			res[to] = arr3
		}
	}
	return res
}

// 按照拓扑序自顶向下进行转移(EnumerateFail).
func (trie *ACAutoMatonArray) Dp(f func(link, cur int32)) {
	for _, v := range trie.BfsOrder {
		if v != 0 {
			f(trie.link[v], v)
		}
	}
}

// 按照拓扑序逆序自底向上进行转移(EnumerateFailReverse).
func (trie *ACAutoMatonArray) DpReverse(f func(link, cur int32)) {
	for i := len(trie.BfsOrder) - 1; i >= 0; i-- {
		v := trie.BfsOrder[i]
		if v != 0 {
			f(trie.link[v], v)
		}
	}
}

func (trie *ACAutoMatonArray) BuildFailTree() [][]int32 {
	res := make([][]int32, trie.Size())
	trie.Dp(func(pre, cur int32) {
		res[pre] = append(res[pre], cur)
	})
	return res
}

func (trie *ACAutoMatonArray) BuildTrieTree() [][]int32 {
	res := make([][]int32, trie.Size())
	for i := int32(1); i < trie.Size(); i++ {
		res[trie.Parent[i]] = append(res[trie.Parent[i]], i)
	}
	return res
}

// 返回str在trie树上的节点位置.如果不存在，返回0.
func (trie *ACAutoMatonArray) Search(str string) int32 {
	if len(str) == 0 {
		return 0
	}
	pos := int32(0)
	for _, char := range str {
		if pos >= int32(len(trie.Children)) || pos < 0 {
			return 0
		}
		ord := char - trie.offset
		if next := trie.Children[pos][ord]; next == -1 {
			return 0
		} else {
			pos = next
		}
	}
	return pos
}

func (trie *ACAutoMatonArray) newNode() int32 {
	nexts := make([]int32, trie.sigma)
	for i := range nexts {
		nexts[i] = -1
	}
	trie.Children = append(trie.Children, nexts)
	trie.Parent = append(trie.Parent, -1)
	trie.Depth = append(trie.Depth, 0)
	return int32(len(trie.Children) - 1)
}

func (trie *ACAutoMatonArray) newNode2(parent int32, char int32) int32 {
	node := trie.newNode()
	trie.Parent[node] = parent
	trie.Depth[node] = trie.Depth[parent] + 1
	return node
}

// !Point Add Range Sum, 0-based.
type BITArray struct {
	n     int32
	total int
	data  []int
}

func NewBitArray(n int32) *BITArray {
	res := &BITArray{n: n, data: make([]int, n)}
	return res
}

func NewBitArrayFrom(n int32, f func(i int32) int) *BITArray {
	total := 0
	data := make([]int, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
		total += data[i]
	}
	for i := int32(1); i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] += data[i-1]
		}
	}
	return &BITArray{n: n, total: total, data: data}
}

func (b *BITArray) Add(index int32, v int) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *BITArray) QueryPrefix(end int32) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *BITArray) QueryRange(start, end int32) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}

func (b *BITArray) QueryAll() int {
	return b.total
}

func (b *BITArray) String() string {
	sb := []string{}
	for i := int32(0); i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
}

type BITRangeAddPointGetArray struct {
	bit *BITArray
}

func NewBITRangeAddPointGet(n int32) *BITRangeAddPointGetArray {
	return &BITRangeAddPointGetArray{bit: NewBitArray(n)}
}

func NewBITRangeAddPointGetFrom(n int32, f func(i int32) int) *BITRangeAddPointGetArray {
	return &BITRangeAddPointGetArray{bit: NewBitArrayFrom(n, f)}
}

func (b *BITRangeAddPointGetArray) AddRange(start, end int32, delta int) {
	if start < 0 {
		start = 0
	}
	if end > b.bit.n {
		end = b.bit.n
	}
	if start >= end {
		return
	}
	b.bit.Add(start, delta)
	b.bit.Add(end, -delta)
}

func (b *BITRangeAddPointGetArray) Get(index int32) int {
	return b.bit.QueryPrefix(index + 1)
}

func (b *BITRangeAddPointGetArray) String() string {
	res := []string{}
	for i := int32(0); i < b.bit.n; i++ {
		res = append(res, fmt.Sprintf("%d", b.Get(i)))
	}
	return fmt.Sprintf("BITRangeAddPointGetArray: [%v]", strings.Join(res, ", "))
}

// PointSetRangeMin

type E = int32

func (*PointSetRangeMin) e() E        { return INF32 }
func (*PointSetRangeMin) op(a, b E) E { return min32(a, b) }

type PointSetRangeMin struct {
	n, size int32
	seg     []E
}

func NewPointSetRangeMin(n int32, f func(int32) E) *PointSetRangeMin {
	res := &PointSetRangeMin{}
	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := int32(0); i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func NewSegmentTreeFrom(leaves []E) *PointSetRangeMin {
	res := &PointSetRangeMin{}
	n := int32(len(leaves))
	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := int32(0); i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *PointSetRangeMin) Get(index int32) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *PointSetRangeMin) Set(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *PointSetRangeMin) Update(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *PointSetRangeMin) Query(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}
func (st *PointSetRangeMin) QueryAll() E { return st.seg[1] }
func (st *PointSetRangeMin) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}

func MatMul(m1, m2 [][]int, mod int) [][]int {
	res := make([][]int, len(m1))
	for i := range res {
		res[i] = make([]int, len(m2[0]))
	}
	for i := 0; i < len(m1); i++ {
		for k := 0; k < len(m2); k++ {
			for j := 0; j < len(m2[0]); j++ {
				res[i][j] = (res[i][j] + m1[i][k]*m2[k][j]) % mod
				if res[i][j] < 0 {
					res[i][j] += mod
				}
			}
		}
	}
	return res
}

func MatPow(mat [][]int, exp int, mod int) [][]int {
	n := len(mat)
	matCopy := make([][]int, n)
	for i := 0; i < n; i++ {
		matCopy[i] = make([]int, n)
		copy(matCopy[i], mat[i])
	}
	res := make([][]int, n)
	for i := 0; i < n; i++ {
		res[i] = make([]int, n)
		res[i][i] = 1
	}
	for exp > 0 {
		if exp&1 == 1 {
			res = MatMul(res, matCopy, mod)
		}
		matCopy = MatMul(matCopy, matCopy, mod)
		exp >>= 1
	}
	return res
}

func Pow(a, b, mod int) int {
	a %= mod
	res := 1
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a >= b {
		return a
	}
	return b
}
