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
//   等于回文树的状态数 - 2（排除奇根和偶根两个状态）。
// !- 以i为结尾的回文子串个数
//   等于以i为结尾的回文树节点的深度。
// - 所有回文子串的数量
//   由于回文树的构造过程中，节点本身就是按照拓扑序插入，
//   因此只需要逆序枚举所有状态，将当前状态的出现次数加到其 fail 指针对应状态的出现次数上即可。
// - 两个串的公共回文子串数量
//   同时dfs遍历其相同的节点，累加两边节点的cnt的乘积即可
// - 最小回文划分
// - 回文分割dp/回文树dp 『Border Series在回文自动机上的运用』=> Palindrome series
//   https://www.cnblogs.com/Parsnip/p/12426971.html
//
// Usage:
// PalindromicTree(S): 构造回文树.
// Add(x): 末尾添加文字x，返回以这个字符为后缀的最长回文串的位置.
// AddString(S): 末尾添加字符串S.
// GetFrequency(): 每个顶点对应的回文串出现的次数.
// RecoverPalindrome(pos): 返回pos位置的回文串.
// UpdateDp(init, apply): 每次增加一个字符,用以当前字符为结尾的`所有本质不同回文串`更新dp值.
// BuildFailTree(): 构造fail树.
// GetNode(pos): 返回pos位置的回文串顶点.
// Size(): 回文树中的顶点个数.
//
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
//
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
//
// https://github.com/EndlessCheng/codeforces-go/blob/646deb927bbe089f60fc0e9f43d1729a97399e5f/copypasta/pam.go#L3
// 回文自动机（回文树）  Palindrome Automaton (PAM) / EerTree
// 如果我们既知道前缀信息（trie），又知道后缀信息（fail），就可以做字符串匹配
// 对于回文计数类问题（比如求一个字符串有多少个本质不同的回文子串，每个又分别出现多少次），Manacher 就心有余而力不足了
// 刚好，PAM 的每个节点就表示一个本质不同回文串
// !类似AC自动机，将回文子串建立联系.
// 回文树的每个节点对应一个回文子串.
// !类似在trie上的插入，不过现在走过一条边的含义变为"在当前回文串的`两边`加上字符c".
// 回文树有两个根：奇根、偶根.
// 奇根的子节点表示回文子串长度为奇数，偶根的子节点表示回文子串长度为偶数.
// 奇根长度为-1，失配指针指向偶根；偶根长度为0，失配指针指向奇根.
//
// 每次增加一个字符，本质不同的回文子串个数最多增加 1 个.一个字符串s的本质不同回文子串个数不超过|s|.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	demo()

	// CF17E()
	// CF835D()
	// CF906E()
	// CF932G()

	// P1659()
	// P3649()
	// P4287()
	// P4555()
	// P4762()
	// P5496()
	// P5555()

	// yuki263()
	// yuki273()
	// yuki465()
	// yuki2606()
}

func demo() {
	// s := "eertree"
	s := "aaaaa"

	// 返回回文树、每个前缀结尾的回文串个数.
	NewPalindromicTreeAndEnds := func(n int, f func(int) int32) (tree *PalindromicTree, ends []int) {
		tree = NewPalindromicTree()
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

// Palisection (不相交回文子串对数)
// https://www.luogu.com.cn/problem/CF17E
// 给出一个串 S，找到所有回文子串，统计有多少个pair他们的区间是相交的。
// 要求对于 S 的每个前缀在线输出答案。
// !可以转化为求不相交的回文串对子数。
// 枚举左边回文的结束位置i(分割点), 右边的回文只要在i以后开始即可.
func CF17E() {
	const MOD int = 51123987

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)
	// n := 97453
	// s := strings.Repeat("a", n)

	// !以i为结尾的回文串的个数(等于在回文树中的深度)
	getEnds := func(n int, f func(int) int32) []int {
		P := NewPalindromicTree()
		ends := make([]int, n)
		counter := make([]int, n+2)
		for i := 0; i < n; i++ {
			pos := P.Add(f(i))
			link := P.Nodes[pos].Link
			counter[pos] = counter[link] + 1
			ends[i] = counter[pos]
		}
		return ends
	}

	ends1 := getEnds(n, func(i int) int32 { return int32(s[i]) })     // 原串以i为结尾的回文串的个数
	ends2 := getEnds(n, func(i int) int32 { return int32(s[n-1-i]) }) // 反串以i为结尾的回文串的个数
	sufSum := make([]int, n+1)
	for i := n - 1; i >= 0; i-- {
		sufSum[i] = sufSum[i+1] + int(ends2[n-1-i])
	}

	allCount := sufSum[0] % MOD
	allPair := allCount * (allCount - 1) / 2 % MOD
	noIntersectPair := 0
	for i := 0; i < n-1; i++ {
		noIntersectPair = (noIntersectPair + ends1[i]*sufSum[i+1]%MOD) % MOD
	}

	res := (allPair - noIntersectPair + MOD) % MOD
	fmt.Println(res)
}

// Palindromic characteristics (k阶回文子串个数)
// https://www.luogu.com.cn/problem/CF835D
// 给定一个长度为n (1≤ n ≤ 1e6)的字符串，对于1<=k<=n，求出k阶回文子串有多少个。
// k阶回文串的定义是：
// k = 1时，就是回文串的定义
// k > 1时，不仅其本身为回文串，并且它的左半边和右半边都是k - 1阶回文串（奇数串左右半边不包括中心）
//
// !通过节点u的 fail 指针找到长度不超过len(u)/2 的最长回文后缀
func CF835D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	T := NewPalindromicTree()
	T.AddString(s)

	getHalfLink := func(pos int32) int32 {
		// targetLen := T.Nodes[pos].Length / 2
		// res := pos
		// for T.Nodes[res].Length > targetLen {
		// 	res = T.Nodes[res].Link
		// }
		// return res
		return T.Nodes[pos].HalfLink
	}

	dp := make([]int, T.Size()) // !每一个节点对应的回文串的阶数
	for i := 2; i < T.Size(); i++ {
		half := getHalfLink(int32(i))
		halfLen := int(T.Nodes[half].Length)
		curLen := int(T.Nodes[i].Length)
		if halfLen == curLen/2 {
			dp[i] = dp[half] + 1
		} else {
			dp[i] = 1
		}
	}

	freq := T.GetFrequency() // !每一个节点对应的回文串的个数

	res := make([]int, len(s)+1)
	for i := 2; i < T.Size(); i++ {
		level := dp[i]
		res[level] += freq[i]
	}
	for i := len(s) - 1; i >= 0; i-- {
		res[i] += res[i+1]
	}
	for _, v := range res[1:] {
		fmt.Fprint(out, v, " ")
	}
}

// Reverses
// https://www.luogu.com.cn/problem/CF906E
// 给定两个长度为n的字符串s和t。可以翻转t的若干个不相交的区间，问最少翻转几个区间，使得s和t相同.
// TODO
func CF906E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s, t string
	fmt.Fscan(in, &s, &t)
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
	t := make([]int32, 0, n)
	for i := 0; i < n/2; i++ {
		t = append(t, int32(s[i]))
		t = append(t, int32(s[n-1-i]))
	}

	counter := make([]int, n+2) // 回文树上每个位置对应的dp值
	dp := make([]int, n+1)      // dp[i]表示前i个字符的分割次数
	dp[0] = 1
	tree := NewPalindromicTree()
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

		// !偶数长度时更新dp
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

// P1659 [国家集训队] 拉拉队排练
// https://www.luogu.com.cn/problem/P1659
// 询问一个串中所有奇回文按照长度降序排列，前k个奇回文的长度乘模MOD.
// 如果少于k个奇回文，输出-1.
func P1659() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 19930726
	qpow := func(a, b int) int {
		a %= MOD
		res := 1
		for ; b > 0; b >>= 1 {
			if b&1 > 0 {
				res = res * a % MOD
			}
			a = a * a % MOD
		}
		return res
	}

	var n, k int
	fmt.Fscan(in, &n, &k)
	var s string
	fmt.Fscan(in, &s)

	T := NewPalindromicTree()
	T.AddString(s)
	freq := T.GetFrequency()

	var odds [][2]int // (长度，出现次数)
	oddFreqSum := 0
	for i := 2; i < len(freq); i++ {
		len_ := int(T.Nodes[i].Length)
		if len_&1 == 1 {
			odds = append(odds, [2]int{len_, freq[i]})
			oddFreqSum += freq[i]
		}
	}

	if oddFreqSum < k {
		fmt.Fprintln(out, -1)
		return
	}

	sort.Slice(odds, func(i, j int) bool { return odds[i][0] > odds[j][0] })
	remain := k
	res := 1
	for _, item := range odds {
		len_, freq := item[0], item[1]
		if remain > freq {
			res = res * qpow(len_, freq) % MOD
			remain -= freq
		} else {
			res = res * qpow(len_, remain) % MOD
			break
		}
	}
	fmt.Fprintln(out, res)
}

// P3649 [APIO2014] 回文串
// https://www.luogu.com.cn/problem/P3649
// !对所有不同的回文子串,求出(回文长度*回文出现次数)的最大值
func P3649() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	T := NewPalindromicTree()
	T.AddString(s)

	freq := T.GetFrequency()
	res := 0
	for pos := 2; pos < len(freq); pos++ {
		node := T.GetNode(pos)
		res = max(res, int(node.Length)*freq[pos])
	}
	fmt.Fprintln(out, res)
}

// P4287 [SHOI2011] 双倍回文
// https://www.luogu.com.cn/problem/P4287
// 如果x能够写成 xx'xx' (x'为x的倒置)，则x为双倍回文.
// 对于给定的字符串，计算它的最长双倍回文子串的长度。
// 对原串构建回文自动机,抽离fail树，从根开始dfs
// 设len[x]表示节点x表示的最长回文子串长度
// 在fail树上，x到根节点的路径上的点表示的字符串包含了x代表的回文子串的所有回文后缀
// 所以若dfs到了x，若len[x]为偶数，标记len[x]*2，如果在x的子树中能找到len为len[x]*2的点，那么len[x]*2就可以用来更新答案
func P4287() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)
	T := NewPalindromicTree()
	T.AddString(s)

	failTree := T.BuildFailTree()
	visitedLen := make([]bool, n+1)
	res := 0

	var dfs func(cur int)
	dfs = func(cur int) {
		len_ := int(T.Nodes[cur].Length)
		if len_ >= 4 && len_%4 == 0 && visitedLen[len_/2] {
			res = max(res, len_)
		}
		if len_ > 0 {
			visitedLen[len_] = true
		}
		for _, next := range failTree[cur] {
			dfs(next)
		}
		if len_ > 0 {
			visitedLen[len_] = false
		}
	}
	dfs(0)

	fmt.Fprintln(out, res)
}

// P4555 [国家集训队] 最长双回文串
// https://www.luogu.com.cn/problem/P4555
// 分别建两个回文自动机，一个是正序，一个是倒序，
// 然后通过自动机得出最长回文长度，接着枚举切割点就好了。
// 给定字符串s,求s的最长双回文子串T，即可将T分割成两部分X,Y，使得X和Y都是回文串。
// len(S)<=1e5
func P4555() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	tree1, tree2 := NewPalindromicTree(), NewPalindromicTree()
	preLen, sufLen := make([]int32, len(s)), make([]int32, len(s)) // !记录每个位置结尾的最长回文串长度
	for i := 0; i < len(s); i++ {
		p1 := tree1.Add(int32(s[i]))
		preLen[i] = tree1.Nodes[p1].Length
		p2 := tree2.Add(int32(s[len(s)-1-i]))
		sufLen[len(s)-1-i] = tree2.Nodes[p2].Length
	}

	res := int32(0)
	// 枚举分割点
	for i := 0; i < len(s)-1; i++ {
		// 以i为分割点的最长双回文串长度
		cur := preLen[i] + sufLen[i+1]
		if cur > res {
			res = cur
		}
	}

	fmt.Fprintln(out, res)
}

// P4762 [CERC2014] Virus synthesis
// https://www.luogu.com.cn/problem/P4762
// 初始有一个空串，有以下操作：
// 1.在串开头或者末尾添加一个字母。
// 2.将当前串复制一份并反接在后面。（AAAB → AAABBAAA）
// 询问要操作成给定的串需要多少次操作。
//
// !可以发现操作 2 越多越好，但是发现操作 2 只能得到长度为偶数的回文串，
// 所以答案其实就是 在一个长度为偶数的回文串的基础上暴力添加字符，即 min(dp[i]+n−len[i])
//
// 转移为：
// !dp[u] = min(dp[u] , dp[halfLink] + len(u)/2 - len(halfLink) + 1) (补全到一半后复制)
// TODO
func P4762() {
	// in := bufio.NewReader(os.Stdin)
	// out := bufio.NewWriter(os.Stdout)
	// defer out.Flush()

	// solve := func(s string) int {
	// 	n := len(s)
	// 	T := NewPalindromicTree()
	// 	T.AddString(s)
	// 	dp := make([]int, T.Size())

	// 	res := n
	// 	for i := 2; i < T.Size(); i++ {
	// 		len_ := int(T.Nodes[i].Length)
	// 		if len_&1 == 0 {
	// 			res = min(res, dp[i]+n-len_/2)
	// 		}
	// 	}
	// 	return res
	// }

	// var T int
	// fmt.Fscan(in, &T)
	// for i := 0; i < T; i++ {
	// 	var s string
	// 	fmt.Fscan(in, &s)
	// 	fmt.Println(solve(s))
	// }
}

// P5496 【模板】回文自动机（PAM）
// https://www.luogu.com.cn/problem/P5496
// 每个位置结尾的回文子串个数(这个节点所表示的回文串中，有多少个后缀也是回文)
// nodes[ch].num = nodes[nodes[ch].fail].num + 1;
func P5496() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	T := NewPalindromicTree()

	res := make([]int32, len(s))       // 每个位置结尾的回文子串个数
	counter := make([]int32, len(s)+2) // !统计每个回文树上每个位置结尾的回文子串个数
	for i := 0; i < len(s); i++ {
		curChar := s[i]
		if i >= 1 {
			curChar = (s[i]-97+byte(res[i-1]%26))%26 + 97
		}

		pos := T.Add(int32(curChar))
		link := T.Nodes[pos].Link
		counter[pos] = counter[link] + 1 // !转移
		res[i] = counter[pos]
		fmt.Fprint(out, res[i], " ")
	}
}

// P5555 秩序魔咒
// https://www.luogu.com.cn/problem/P5555
// 给定两个字符串，求最长公共回文子串的长度，以及有多少并列最长的子串。
func P5555() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	var s, t string
	fmt.Fscan(in, &s, &t)
	T := NewPalindromicTree()
	T.AddString(s + "><" + t)
	dps := make([]int, T.Size())
	dpt := make([]int, T.Size())
	lenS := int32(len(s))
	for i := 0; i < T.Size(); i++ {
		for _, j := range T.Nodes[i].Indexes { // 回文出现位置
			if j < lenS {
				dps[i]++
			} else if j >= lenS+2 {
				dpt[i]++
			}
		}
	}

	res := 0
	maxLen := 0
	for i := T.Size() - 1; i >= 2; i-- { // 按照拓扑序遍历本质不同回文
		if dps[i] > 0 && dpt[i] > 0 {
			if int(T.Nodes[i].Length) > maxLen {
				maxLen = int(T.Nodes[i].Length)
				res = 1
			} else if int(T.Nodes[i].Length) == maxLen {
				res += 1
			}
		}
		// 因为求最长，所以可以不向上转移
		// link := T.Nodes[i].Link
		// dps[link] += dps[i]
		// dpt[link] += dpt[i]
	}

	fmt.Fprintln(out, maxLen, res)
}

// No.263 Common Palindromes Extra
// 求两个字符串的公共回文子串的对数 n<=5e5
// https://yukicoder.me/problems/no/263
func yuki263() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s, t string
	fmt.Fscan(in, &s, &t)
	T := NewPalindromicTree()
	T.AddString(s + "><" + t)
	dps := make([]int, T.Size())
	dpt := make([]int, T.Size())
	lenS := int32(len(s))
	for i := 0; i < T.Size(); i++ {
		for _, j := range T.Nodes[i].Indexes { // 回文出现位置
			if j < lenS {
				dps[i]++
			} else if j >= lenS+2 {
				dpt[i]++
			}
		}
	}

	res := 0
	for i := T.Size() - 1; i >= 2; i-- { // 按照拓扑序遍历本质不同回文
		res += dps[i] * dpt[i]
		link := T.Nodes[i].Link
		dps[link] += dps[i]
		dpt[link] += dpt[i]
	}
	fmt.Fprintln(out, res)
}

// No.273 回文分解
// https://yukicoder.me/problems/no/273
// !将字符串s分成若干段,保证每段都是回文串(段数>=2)
// 求最长的回文串的最大值
// 2<=len(s)<=1e5
func yuki273() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const INF int = 1e18

	var s string
	fmt.Fscan(in, &s)
	n := len(s)
	maxLen := make([]int, n+2) // dp1[i]表示以i回文树位置为结尾的最长回文串的长度
	for i := range maxLen {
		maxLen[i] = -INF
	}
	dp := make([]int, n+1) // dp2[i]表示前i个字符的最长回文串的长度
	for i := range dp {
		dp[i] = -INF
	}
	dp[0] = 1

	T := NewPalindromicTree()
	for i, c := range s {
		T.Add(c)

		indexes := T.UpdateDp(
			// 基于 s[start:i+1] 这一段回文初始化 pos 处的值
			func(pos, start int) {
				if i+1-start == n {
					maxLen[pos] = dp[start] // 段数需要>=2
				} else {
					maxLen[pos] = max(dp[start], i+1-start)
				}
			},
			// 基于 pre 的信息更新 pos 处的值
			func(pos, pre int) {
				if int(T.Nodes[pos].Length) == n {
					maxLen[pos] = max(maxLen[pos], maxLen[pre]) // 段数需要>=2
				} else {
					maxLen[pos] = max(maxLen[pos], max(maxLen[pre], int(T.Nodes[pos].Length)))
				}
			},
		)

		for _, p := range indexes {
			dp[i+1] = max(dp[i+1], maxLen[p])
		}
	}

	fmt.Fprintln(out, dp[n])
}

// No.465 PPPPPPPPPPPPPPPPAPPPPPPPP
// https://yukicoder.me/problems/no/465
// 假定p1,p2,p3是长度大于等于1的回文,A是长度大于等于1的字符串.
// 令PPAP=p1+p2+A+p3 求有多少组(p1,p2,A,p3) 满足 PPAP 等于给出的字符串s.
func yuki465() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	n := len(s)
	// dp1: 前缀是否是回文串
	// dp2: 前缀能拆成p1+p2的方案数
	dp1, dp2 := make([]int, n+1), make([]int, n+1)
	buf := make([]int, n+2)
	tree := NewPalindromicTree()
	for i, c := range s {
		pos := tree.Add(c)
		if int(tree.Nodes[pos].Length) == i+1 {
			dp1[i+1] = 1
		}

		indexes := tree.UpdateDp(
			// 初始化顶点pos的dp值,对应回文串s[start:i+1]
			func(pos, start int) {
				buf[pos] = dp1[start]
			},
			// 用pre更新pos
			func(pos, pre int) {
				buf[pos] += buf[pre]
			},
		)

		for _, p := range indexes {
			dp2[i+1] += buf[p]
		}
	}

	tree2 := NewPalindromicTree()
	res, sum := 0, 0
	for i := n - 1; i >= 0; i-- {
		res += dp2[i] * sum
		pos := tree2.Add(int32(s[i]))
		if int(tree2.Nodes[pos].Length) == n-i {
			sum++
		}
	}

	fmt.Fprintln(out, res)
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

// 5. 最长回文子串
// https://leetcode.cn/problems/longest-palindromic-substring/
func longestPalindrome(s string) string {
	start, end, maxLen := 0, 0, 0
	tree := NewPalindromicTree()
	for i, c := range s {
		pos := tree.Add(c)
		if int(tree.Nodes[pos].Length) > maxLen {
			maxLen = int(tree.Nodes[pos].Length)
			start, end = i-maxLen+1, i+1
		}
	}
	return s[start:end]
}

// 132. 分割回文串 II
// https://leetcode.cn/problems/palindrome-partitioning-ii
// 给你一个字符串 s，请你将 s 分割成一些子串，使每个子串都是回文。
// 返回符合要求的 最少分割次数 。
func minCut(s string) int {
	const INF int = 1e18

	n := len(s)
	minCount := make([]int, n+2) // 回文树上每个位置对应的dp值
	for i := range minCount {
		minCount[i] = INF
	}
	dp := make([]int, n+1) // dp[i]表示前i个字符的最少分割次数
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0

	T := NewPalindromicTree()
	for i, c := range s {
		T.Add(c)

		indexes := T.UpdateDp(
			// 基于 s[start:i+1] 这一段回文初始化 pos 处的值
			func(pos, start int) {
				minCount[pos] = dp[start]
			},
			// 基于 pre 的信息更新 pos 处的值
			func(pos, pre int) {
				minCount[pos] = min(minCount[pos], minCount[pre])
			},
		)

		for _, p := range indexes {
			dp[i+1] = min(dp[i+1], minCount[p]+1)
		}
	}

	return dp[n] - 1
}

// !647. 回文子串-回文子串个数(以i为结尾的回文子串个数)
// https://leetcode.cn/problems/palindromic-substrings/description/
func countSubstrings(s string) int {
	// tree := NewPalindromicTree()
	// ends := make([]int, len(s))      // 每个位置结尾的回文子串个数
	// counter := make([]int, len(s)+2) // !统计每个回文树上每个位置结尾的回文子串个数
	// for i, c := range s {
	// 	pos := tree.Add(c)
	// 	counter[pos] = counter[tree.Nodes[pos].Link] + 1 // !转移
	// 	ends[i] = counter[pos]
	// }
	// res := 0
	// for i := 0; i < len(s); i++ {
	// 	res += ends[i]
	// }
	// return res
	tree := NewPalindromicTree()
	tree.AddString(s)
	counter := tree.GetFrequency()
	res := 0
	for _, v := range counter[2:] {
		res += v
	}
	return res
}

// 1745. 分割回文串 IV
// 能否划分成三段非空回文
// https://leetcode.cn/problems/palindrome-partitioning-iv/description/
func checkPartitioning(s string) bool {
	n := len(s)
	// dp1: 前缀是否是回文串
	// dp2: 前缀是否能拆成p1+p2
	dp1, dp2 := make([]int, n+1), make([]int, n+1)
	buf := make([]int, n+2)
	tree := NewPalindromicTree()
	for i, c := range s {
		pos := tree.Add(c)
		if int(tree.Nodes[pos].Length) == i+1 {
			dp1[i+1] = 1
		}

		indexes := tree.UpdateDp(
			// 初始化顶点pos的dp值,对应回文串s[start:i+1]
			func(pos, start int) {
				buf[pos] = dp1[start]
			},
			// 用pre更新pos
			func(pos, pre int) {
				buf[pos] += buf[pre]
			},
		)

		for _, p := range indexes {
			dp2[i+1] += buf[p]
		}
	}

	tree2 := NewPalindromicTree()
	res := 0
	for i := n - 1; i >= 0; i-- {
		pos := tree2.Add(int32(s[i]))
		if int(tree2.Nodes[pos].Length) == n-i {
			res += dp2[i]
		}
	}

	return res > 0
}

// 1960. 两个回文子字符串长度的最大乘积
// https://leetcode.cn/problems/maximum-product-of-the-length-of-two-palindromic-substrings/description/
// 你需要找到两个 `不重叠` 的回文 子字符串，它们的长度都必须为 奇数 ，使得它们长度的乘积最大。
// 2 <= s.length <= 1e5
func maxProduct(s string) int64 {
	// s[i]结尾的最长奇回文串长度
	getDp := func(n int, f func(i int) int32) []int {
		res := make([]int, n+1)
		T := NewPalindromicTree()

		memo := make([]int32, n+2)
		for i := range memo {
			memo[i] = -1
		}
		// 沿着后缀链接向上跳，直到找到一个长度为奇数的回文串
		var find func(int32) int32
		find = func(pos int32) int32 {
			if pos == 0 {
				return 0
			}
			if memo[pos] != -1 {
				return memo[pos]
			}
			if T.Nodes[pos].Length&1 == 1 {
				return pos
			}
			memo[pos] = find(T.Nodes[pos].Link)
			return memo[pos]
		}

		for i := 0; i < n; i++ {
			pos := int32(T.Add(f(i)))
			pos = find(pos)
			res[i+1] = max(res[i], int(T.Nodes[pos].Length))
		}
		return res
	}

	n := len(s)
	preMax := getDp(n, func(i int) int32 { return int32(s[i]) })
	sufMax := getDp(n, func(i int) int32 { return int32(s[n-i-1]) })
	for i, j := 0, len(sufMax)-1; i < j; i, j = i+1, j-1 {
		sufMax[i], sufMax[j] = sufMax[j], sufMax[i]
	}
	res := 0
	for i := 0; i < n; i++ {
		res = max(res, preMax[i]*sufMax[i])
	}
	return int64(res)
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
}

func NewPalindromicTree() *PalindromicTree {
	res := &PalindromicTree{}
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
			for pt.Ords[pos] != pt.Ords[pos-halfNode.Length-1] || (halfNode.Length+2)*2 > newNode.Length {
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

// 求出每个顶点对应的回文串出现的次数(注意不是每个以每个前缀结尾的回文串的次数,而是每个顶点代表的回文串的次数)
func (pt *PalindromicTree) GetFrequency() []int {
	res := make([]int, pt.Size())
	// !节点编号从大到小，就是 fail 树的拓扑序
	for i := pt.Size() - 1; i >= 1; i-- { // 除去根节点(奇根)
		res[i] += len(pt.Nodes[i].Indexes)
		res[pt.Nodes[i].Link] += res[i] // 长回文包含短回文
	}
	return res
}

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
		if rev >= 0 && pt.Ords[rev] == pt.Ords[len(pt.Ords)-1] {
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
