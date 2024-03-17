// 循环节及border理论
// 循环节:
// 一个字符串的循环节是指一个非空字符串, 使得原字符串是由循环节重复若干次(>=2)得到的.
// 例如, "ababab"的循环节是"ab".
// https://www.cnblogs.com/alex-wei/p/Common_String_Theory_Theory.html

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

func main() {
	// P2375()
	// P3426()
	// P3435()
	// P3449()
	P3449V2()
	// P3538()
	// P4391()
	// P5829()

	// CF526D()

	// acwing143()
}

// 459. 重复的子字符串(是否存在循环节)
// https://leetcode.cn/problems/repeated-substring-pattern/description/
// 给定一个非空的字符串 s ，检查是否可以通过由它的一个子串重复多次构成。
func repeatedSubstringPattern(s string) bool {
	n := int32(len(s))
	next := GetNext(n, func(i int32) int32 { return int32(s[i]) })
	return Period(next, n-1) > 0
}

// P2375 [NOI2014] 动物园 (halfLink)
// https://www.luogu.com.cn/problem/P2375
// 对每个前缀求其最长的长度不超过串长一半的 border 在失配树上的深度.
func P2375() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e9 + 7

	// O(nlogn) 对每个前缀求其最长的长度不超过串长一半的 border 在失配树上的深度.
	query := func(s string) []int32 {
		n := int32(len(s))
		nexts := GetNext(n, func(i int32) int32 { return int32(s[i]) })
		db := NewDoubling(n+1, int(n+1))
		depths := make([]int32, n+1)
		for i := int32(1); i <= n; i++ {
			parent := nexts[i-1]
			db.Add(i, parent)
			depths[i] = depths[parent] + 1
		}
		db.Build()

		res := make([]int32, n) // res[i]表示前缀s[:i+1]的最长border在失配树上的深度
		for i := int32(1); i <= n; i++ {
			halfLen := i >> 1
			_, tmp := db.MaxStep(i, func(next int32) bool { return next > halfLen })
			halfLink := nexts[tmp-1]
			res[i-1] = depths[halfLink]
		}
		return res
	}

	// O(n) 对每个前缀求其最长的长度不超过串长一半的 border 在失配树上的深度.
	query2 := func(s string) []int32 {
		n := int32(len(s))
		f := func(i int32) int32 { return int32(s[i]) }
		nexts := GetNext(n, f)
		return GetHalfLinkLength(n, f, nexts)
	}

	var T int32
	fmt.Fscan(in, &T)
	for i := int32(0); i < T; i++ {
		var s string
		fmt.Fscan(in, &s)
		// lens := query(s)
		_ = query
		lens := query2(s)
		res := 1
		for _, v := range lens {
			res = res * int((v + 1)) % MOD
		}
		fmt.Fprintln(out, res)
	}
}

// P3426 [POI2005] SZA-Template
// https://www.luogu.com.cn/problem/P3426
// 你打算在纸上印一串字母。
// 为了完成这项工作，你决定刻一个印章。印章每使用一次，就会将印章上的所有字母印到纸上。
// 同一个位置的相同字符可以印多次。例如：用 aba 这个印章可以完成印制 ababa 的工作（中间的 a 被印了两次）。
// 但是，因为印上去的东西不能被抹掉，在同一位置上印不同字符是不允许的。
// 例如：用 aba 这个印章不可以完成印制 abcba 的工作。
// 因为刻印章是一个不太容易的工作，你希望印章的字符串长度尽可能小。
// !输出一个整数，代表印章上字符串长度的最小值。
//
// 印章为 ababbaba。
// 印制过程如下：
//
// ababbababbabababbabababbababbaba
// ababbaba(8)
//
//			ababbaba(13)
//	    	   ababbaba(20)
//	              ababbaba(27)
//	                   ababbaba(32)
//
// 印章是原串的一个border。
// 容易发现该border出现在了多个位置，且间距需要小于等于border的长度。
// !建出失配树，容易发现印章上的字符只有可能是n到根上的所有点u，且点u合法的条件是`子树内出现的点的最大间距不超过u`
//
// 维护maxGap的方法:
// 1. 只删除，不添加：直接用set维护前驱后继, maxGap只会不断增加.(本题)
// 2. 删除和添加都有：set维护前驱后继的同时，可删除堆/平衡树维护gap长度.
func P3426() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	n := int32(len(s))
	next := GetNext(n, func(i int32) int32 { return int32(s[i]) })
	tree := BuildFailTree(next)
	fs := NewFastSetFrom(int(n+1), func(i int) bool { return true })

	chain := []int32{} // n到根上的所有点
	for cur := int32(n); cur > 0; cur = next[cur-1] {
		chain = append(chain, cur)
	}
	chain = append(chain, 0)
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}

	maxGap := int32(1)
	erase := func(x int) {
		fs.Erase(x)
		prev := fs.Prev(x)
		if prev == -1 {
			return
		}
		next := fs.Next(x)
		if next == fs.n {
			return
		}
		maxGap = max32(maxGap, int32(next)-int32(prev))
	}

	var dfs func(int32, int32)
	dfs = func(cur int32, stop int32) {
		erase(int(cur))
		for _, next := range tree[cur] {
			if next != stop { // 保留这颗子树内的点，删除其他分支上的点
				dfs(next, cur)
			}
		}
	}

	res := n
	for i := 1; i < len(chain); i++ {
		dfs(chain[i-1], chain[i])
		if maxGap <= chain[i] {
			res = chain[i]
			break
		}
	}

	fmt.Fprintln(out, res)
}

// P3435 [POI2006] OKR-Periods of Words
// https://www.luogu.com.cn/problem/P3435
// !求给定字符串所有前缀的最大周期长度之和。
// 最大周期，最小border => 等价于求出每个点在失配树上的最浅非 0 祖先。
func P3435() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	queryPrefixMaxPeriod := func(s string) []int32 {
		n := int32(len(s))
		next := GetNext(n, func(i int32) int32 { return int32(s[i]) })
		tree := BuildFailTree(next)
		lca := NewLCA32(tree, []int32{0})
		res := make([]int32, n)
		for i := int32(1); i <= n; i++ {
			depth := lca.Depth[i]
			minBorder := lca.KthAncestor(i, depth-1)
			res[i-1] = i - minBorder
		}
		return res
	}

	queryPrefixMaxPeriod2 := func(s string) []int32 {
		n := int32(len(s))
		next := GetNext(n, func(i int32) int32 { return int32(s[i]) })
		tree := BuildFailTree(next)
		minBorder := make([]int32, n+1)
		var dfs func(int32, int32)
		dfs = func(u int32, border int32) {
			minBorder[u] = border
			for _, v := range tree[u] {
				dfs(v, border)
			}
		}
		for _, next := range tree[0] {
			dfs(next, next)
		}
		res := make([]int32, n)
		for i := int32(1); i <= n; i++ {
			res[i-1] = i - minBorder[i]
		}
		return res
	}

	var n int32
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)
	_ = queryPrefixMaxPeriod
	borders := queryPrefixMaxPeriod2(s)
	res := 0
	for _, v := range borders {
		res += int(v)
	}
	fmt.Fprintln(out, res)
}

// P3449 [POI2006] PAL-Palindromes (拼接回文串1)
// https://www.luogu.com.cn/problem/P3449
// 给定一些回文串，问这些串两两拼接(自己可以和自己拼接)，最终有多少个回文串.
//
// !结论：给定两个回文串a、b，当且仅当这`两个回文串的最短回文整周期串相同`时，a+b是回文串.“
// !具体而言，如果 n%(n-next[n-1]) == 0，则最短回文整周期串为s[:n-next[n-1]]，否则为s本身.
func P3449() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := range words {
		var m int32
		fmt.Fscan(in, &m)
		fmt.Fscan(in, &words[i])
	}

	group := make(map[string]int32) // 按照最短回文整周期串分组
	for _, word := range words {
		n := int32(len(word))
		next := GetNext(n, func(i int32) int32 { return int32(word[i]) })
		maybeMinPeriod := n - next[n-1]
		if n%maybeMinPeriod == 0 {
			group[word[:maybeMinPeriod]]++
		} else {
			group[word]++
		}
	}

	res := 0
	for _, v := range group {
		res += int(v) * int(v)
	}
	fmt.Fprintln(out, res)
}

// https://www.luogu.com.cn/problem/P3449
// !给定一些字符串，问这些串两两拼接(自己可以和自己拼接)，最终有多少个回文串(回文对).
// 不是回文串也能做.
// 三种情况:
// 1. 长度相等，abc+cba
// 2. len(s1)>len(s2)，abc(xxx)+cba，回文中心在s1上，说明s2的反串是A的前缀,且s1剩余后缀是回文串.
// 3. len(s1)<len(s2)，cba+(xxx)abc, 回文中心在s2上，说明s1的反串是B的后缀,且s2剩余前缀是回文串.
// 枚举前缀用trie，判断回文用manacher.
func P3449V2() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := range words {
		var m int32
		fmt.Fscan(in, &m)
		fmt.Fscan(in, &words[i])
	}

	manachers := make([]*Manacher, len(words))
	for i, word := range words {
		manachers[i] = NewManacher(word)
	}
	// 询问一个串的后缀s[start:)是否是回文串(这里空串不为回文串).
	isSuffixPalindrome := func(wordIndex int32, start int32) bool {
		n := int32(len(words[wordIndex]))
		if start >= n {
			return false
		}
		return manachers[wordIndex].IsPalindrome(start, n)
	}
	// 询问一个串的前缀s[:end)是否是回文串(这里空串不为回文串).
	isPrefixPalindrome := func(wordIndex int32, end int32) bool {
		if end <= 0 {
			return false
		}
		return manachers[wordIndex].IsPalindrome(0, end)
	}

	prefixTrie := NewTrie()
	suffixTrie := NewTrie()
	for _, word := range words {
		n := int32(len(word))
		prefixTrie.Insert(n, func(i int32) byte { return word[i] })
		suffixTrie.Insert(n, func(i int32) byte { return word[n-1-i] })
	}

	res := 0
	for wi, word := range words {
		m := int32(len(word))

		// 拼接的两个字符串长度相同
		res += int(suffixTrie.Count(m, func(i int32) byte { return word[i] }))

		root1 := suffixTrie.root
		// AB+A，枚举前缀A
		for j := 0; j < len(word); j++ {
			char := word[j]
			next, ok := root1.Children[char]
			if !ok {
				break
			}
			root1 = next
			if root1.EndCount > 0 && isSuffixPalindrome(int32(wi), int32(j+1)) {
				res += int(root1.EndCount)
			}
		}

		root2 := prefixTrie.root
		// B+AB，枚举后缀B
		for j := len(word) - 1; j >= 0; j-- {
			char := word[j]
			next, ok := root2.Children[char]
			if !ok {
				break
			}
			root2 = next
			if root2.EndCount > 0 && isPrefixPalindrome(int32(wi), int32(j)) {
				res += int(root2.EndCount)
			}
		}
	}

	fmt.Fprintln(out, res)

}

// P3538 [POI2012] OKR-A Horrible Poem (子串最小循环节)
// https://www.luogu.com.cn/problem/P3538
// !给定一个小写英文字母组成的字符串s和q个查询，每个查询包含两个整数start和end，表示询问s[start:end)的最小循环节长度.
// 如果不存在循环节(重复次数>=2)，返回end-start.
//
// !结论: 若s[start+k,end) == s[start:end-k)，则k为s[start:end)的一个循环节.
// 由于循环节长度为串长的约数，枚举验证即可.
func P3538() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)

	R := NewRollingHashUnsafe(131)
	table := R.Build(s)

	// 求s[start:end)的最小循环节长度
	query := func(start, end int32) (res int32) {
		m := end - start
		EnumerateFactors(m, func(k int32) bool {
			if R.Query(table, start+k, end) == R.Query(table, start, end-k) {
				res = k
				return true
			}
			return false
		})
		return
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var start, end int32
		fmt.Fscan(in, &start, &end)
		start--
		fmt.Fprintln(out, query(start, end))
	}
}

// P4391 [BOI2009] Radio Transmission 无线传输
// https://www.luogu.com.cn/problem/P4391
// !给定字符串s，需要找到一个尽可能短的前缀p满足s是pp..pp的子串.
// 答案为 n-next[n-1]，其中next为s的next数组.
func P4391() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)
	next := GetNext(n, func(i int32) int32 { return int32(s[i]) })
	fmt.Fprintln(out, n-next[n-1])
}

// P5829 【模板】失配树 (求解两个前缀的最长公共 border)
// https://www.luogu.com.cn/problem/P5829
// https://blog.csdn.net/Fighting_Peter/article/details/112724803
// 对s的每个前缀，将其向maxBorder连边，得到一棵树。树的每个结点表示前缀s[:i](0<=i<=n).
// 在 KMP 算法中注意到 next[i]<i，因此若从 next[i] 向 i 连边，我们最终会得到一棵有根树。这就是失配树(Border Tree/Fail Tree)。
// 失配树有很好的性质：
// !对于树上两个点a,b,a是b的祖先<=>s[:a+1]是s[:b+1]的border。这一点由 next 的性质可以得到.
// 因此，若需要查询u前缀和v前缀的最长公共 border, 只需要查询u和v在失配树上的 LCA 即可.
//
//   - aaaabbabbaa
//   - [[1 5 6 8 9] [2 7 10] [3] [4] [] [] [] [] [] [] []]
//   - 0:
//     1:a
//     2:aa
//     3:aaa
//     4:aaaa
//     5:aaaab
//     6:aaaabb
//     7:aaaabba
//     8:aaaabbab
//     9:aaaabbabb
//     10:aaaabbabba
//     11:aaaabbabbaa
//
// - 失配树(failTree)，每个点指向当前前缀的最长真后缀，深度表示当前border长度，父亲节点的编号一定比其儿子节点的编号要小.
//
//	         0
//	     / / | \ \
//	    1  5 6 8  9
//	  / | \
//	 2  7  10
//	/ \
//	3  11
//	|
//	4
func P5829() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	var q int
	fmt.Fscan(in, &q)

	n := int32(len(s))
	next := GetNext(n, func(i int32) int32 { return int32(s[i]) })
	tree := make([][]int32, n+1) // 结点i表示前缀s[:i]
	for i := int32(1); i <= n; i++ {
		p := next[i-1]
		tree[p] = append(tree[p], i)
	}

	L := NewLCA32(tree, []int32{0})
	for i := 0; i < q; i++ {
		var a, b int32
		fmt.Fscan(in, &a, &b)

		// !s[:a]和s[:b]的最长公共 border
		if a == 0 || b == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		lca := L.LCA(a, b)
		if lca == a || lca == b {
			lca = next[lca-1] // !Border 不能是自己，所以当 LCA 是两数之一时，要再往上跳一格
		}
		fmt.Fprintln(out, lca)
	}
}

// Om Nom and Necklace
// https://www.luogu.com.cn/problem/CF526D
// 给定一个长度为 n 的字符串和一个正整数k，判断其每个前缀是否形如 ABABA...BA
// A、B 可以为空，也可以是一个字符串，A 有 k+1 个，B 有 k 个.
//
// 令C=A+B,那么字符串可以变成 CC……CA，其中 C 有 k 个，且 A 是 C 的前缀（由定义而知）。
// 枚举C...C的长度，如果满足条件，则当前字符串一定具有周期 i-next[i-1](可能不完整).
func CF526D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int32
	fmt.Fscan(in, &n, &k)
	var s string
	fmt.Fscan(in, &s)

	Z := ZAlgo(s)
	Z = append(Z, 0) // 哨兵
	next := GetNext(n, func(i int32) int32 { return int32(s[i]) })
	D := NewDiffArray(n + 1)

	for i := k; i <= n; i += k { // C...C可能的长度
		maybePeriod := i - next[i-1] // border性质
		if (i/k)%maybePeriod == 0 {
			start, end := i, i+min32(i/k, Z[i])+1
			D.Add(start, end, 1)
		}
	}

	res := strings.Builder{}
	for i := int32(1); i <= n; i++ {
		if D.Get(i) > 0 {
			res.WriteByte('1')
		} else {
			res.WriteByte('0')
		}
	}
	fmt.Fprintln(out, res.String())
}

// 对每个前缀s[:i+1]，求具有循环节的前缀的长度和对应的循环次数
// 要求循环次数>=2.
// https://www.acwing.com/problem/content/discussion/index/143/1/
func acwing143() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(s string) [][2]int32 {
		n := int32(len(s))
		res := make([][2]int32, n)
		next := GetNext(n, func(i int32) int32 { return int32(s[i]) })
		for i := int32(0); i < n; i++ {
			period := Period(next, i)
			if period > 0 {
				res[i][0] = i + 1
				res[i][1] = (i + 1) / period
			}
		}
		return res
	}
	ptr := 1
	for {
		var n int
		fmt.Fscan(in, &n)
		if n == 0 {
			break
		}
		var s string
		fmt.Fscan(in, &s)
		res := solve(s)
		fmt.Fprintf(out, "Test case #%d\n", ptr)
		ptr++
		for _, v := range res {
			if v[0] == 0 {
				continue
			}
			fmt.Fprintf(out, "%d %d\n", v[0], v[1])
		}
		fmt.Fprintln(out)
	}
}

// `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度.
func GetNext(n int32, f func(i int32) int32) []int32 {
	next := make([]int32, n)
	j := int32(0)
	for i := int32(1); i < n; i++ {
		for j > 0 && f(i) != f(j) {
			j = next[j-1]
		}
		if f(i) == f(j) {
			j++
		}
		next[i] = j
	}
	return next
}

// `halfLinkLength[i]`表示`[:i+1]`这一段子串长度不超过串长一半的最长的border长度.
func GetHalfLinkLength(n int32, f func(i int32) int32, nexts []int32) (halfLinkLength []int32) {
	depth := make([]int32, n+1)
	for i := int32(1); i <= n; i++ {
		parent := nexts[i-1]
		depth[i] = depth[parent] + 1
	}
	halfLinkLength = make([]int32, n)
	pos := int32(0)
	for i := int32(1); i < n; i++ {
		for pos > 0 && f(i) != f(pos) {
			pos = nexts[pos-1]
		}
		if f(i) == f(pos) {
			pos++
		}
		for pos > (i+1)>>1 {
			pos = nexts[pos-1]
		}
		halfLinkLength[i] = depth[pos]
	}
	return
}

// 求s的前缀[0:i+1)的最小周期.如果不存在,则返回0.
//
// !要求循环节次数>1 且 循环节完整.
//
//	0<=i<len(s).
func Period(next []int32, i int32) int32 {
	res := i + 1 - next[i]
	// !循环节次数>1 且 循环节完整
	if i+1 > res && (i+1)%res == 0 {
		return res
	}
	return 0
}

// z算法求字符串每个后缀与原串的最长公共前缀长度
//
// z[0]=0
// z[i]是后缀s[i:]与s的最长公共前缀(LCP)的长度 (i>=1)
func ZAlgo(s string) []int32 {
	n := int32(len(s))
	z := make([]int32, n)
	left, right := int32(0), int32(0)
	for i := int32(0); i < n; i++ {
		z[i] = max32(min32(z[i-left], right-i+1), 0)
		for i+z[i] < n && s[z[i]] == s[i+z[i]] {
			left, right = i, i+z[i]
			z[i]++
		}
	}
	return z
}

func ZAlgoNums(nums []int32) []int32 {
	n := int32(len(nums))
	z := make([]int32, n)
	left, right := int32(0), int32(0)
	for i := int32(0); i < n; i++ {
		z[i] = max32(min32(z[i-left], right-i+1), 0)
		for i+z[i] < n && nums[z[i]] == nums[i+z[i]] {
			left, right = i, i+z[i]
			z[i]++
		}
	}
	return z
}

type S = string

type RollingHashUnsafe struct {
	base  uint
	power []uint
}

// 131/13331/1713302033171(回文素数)
func NewRollingHashUnsafe(base uint) *RollingHashUnsafe {
	return &RollingHashUnsafe{
		base:  base,
		power: []uint{1},
	}
}

func (r *RollingHashUnsafe) Build(s S) (hashTable []uint) {
	sz := len(s)
	hashTable = make([]uint, sz+1)
	for i := 0; i < sz; i++ {
		hashTable[i+1] = hashTable[i]*r.base + uint(s[i])
	}
	return hashTable
}

func (r *RollingHashUnsafe) Query(sTable []uint, start, end int32) uint {
	r.expand(end - start)
	return sTable[end] - sTable[start]*r.power[end-start]
}

func (r *RollingHashUnsafe) Combine(h1, h2 uint, h2len int32) uint {
	r.expand(h2len)
	return h1*r.power[h2len] + h2
}

func (r *RollingHashUnsafe) AddChar(hash uint, c byte) uint {
	return hash*r.base + uint(c)
}

// 两个字符串的最长公共前缀长度.
func (r *RollingHashUnsafe) LCP(sTable []uint, start1, end1 int32, tTable []uint, start2, end2 int32) int32 {
	len1 := end1 - start1
	len2 := end2 - start2
	len := min32(len1, len2)
	low := int32(0)
	high := len + 1
	for high-low > 1 {
		mid := (low + high) / 2
		if r.Query(sTable, start1, start1+mid) == r.Query(tTable, start2, start2+mid) {
			low = mid
		} else {
			high = mid
		}
	}
	return low
}

func (r *RollingHashUnsafe) expand(sz int32) {
	if int32(len(r.power)) < sz+1 {
		preSz := int32(len(r.power))
		r.power = append(r.power, make([]uint, sz+1-preSz)...)
		for i := preSz - 1; i < sz; i++ {
			r.power[i+1] = r.power[i] * r.base
		}
	}
}

type DiffArray struct {
	diff  []int
	dirty bool
}

func NewDiffArray(n int32) *DiffArray {
	return &DiffArray{
		diff: make([]int, n+1),
	}
}

func (d *DiffArray) Add(start, end int32, delta int) {
	if start < 0 {
		start = 0
	}
	if end >= int32(len(d.diff)) {
		end = int32(len(d.diff)) - 1
	}
	if start >= end {
		return
	}
	d.dirty = true
	d.diff[start] += delta
	d.diff[end] -= delta
}

func (d *DiffArray) Build() {
	if d.dirty {
		preSum := make([]int, len(d.diff))
		for i := 1; i < len(d.diff); i++ {
			preSum[i] = preSum[i-1] + d.diff[i]
		}
		d.diff = preSum
		d.dirty = false
	}
}

func (d *DiffArray) Get(pos int32) int {
	d.Build()
	return d.diff[pos]
}

func (d *DiffArray) GetAll() []int {
	d.Build()
	return d.diff[:len(d.diff)-1]
}

// 空间复杂度为O(1)的枚举因子.枚举顺序为从小到大.
func EnumerateFactors(n int32, f func(factor int32) (shouldBreak bool)) {
	if n <= 0 {
		return
	}
	i := int32(1)
	upper := int32(math.Sqrt(float64(n)))
	for ; i <= upper; i++ {
		if n%i == 0 {
			if f(i) {
				return
			}
		}
	}
	i--
	if i*i == n {
		i--
	}
	for ; i > 0; i-- {
		if n%i == 0 {
			if f(n / i) {
				return
			}
		}
	}
}

// 根据nexts数据构建fail树.
// 返回的fail树长度为n+1，结点i表示前缀s[:i].
func BuildFailTree(nexts []int32) [][]int32 {
	n := int32(len(nexts))
	tree := make([][]int32, n+1)
	for i := int32(1); i <= n; i++ {
		p := nexts[i-1]
		tree[p] = append(tree[p], i)
	}
	return tree
}

type LCA32 struct {
	Depth, Parent           []int32
	Tree                    [][]int32
	lid, rid, top, heavySon []int32
	idToNode                []int32
	dfnId                   int32
}

func NewLCA32(tree [][]int32, roots []int32) *LCA32 {
	n := len(tree)
	lid := make([]int32, n) // vertex => dfn
	rid := make([]int32, n)
	top := make([]int32, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int32, n)    // 深度
	parent := make([]int32, n)   // 父结点
	heavySon := make([]int32, n) // 重儿子
	idToNode := make([]int32, n)
	res := &LCA32{
		Tree:     tree,
		lid:      lid,
		rid:      rid,
		top:      top,
		Depth:    depth,
		Parent:   parent,
		heavySon: heavySon,
		idToNode: idToNode,
	}
	for _, root := range roots {
		res._build(root, -1, 0)
		res._markTop(root, root)
	}
	return res
}

// 树上多个点的 LCA，就是 DFS 序最小的和 DFS 序最大的这两个点的 LCA.
func (hld *LCA32) LCAMultiPoint(nodes []int32) int32 {
	if len(nodes) == 1 {
		return nodes[0]
	}
	if len(nodes) == 2 {
		return hld.LCA(nodes[0], nodes[1])
	}
	minDfn, maxDfn := int32(1<<31-1), int32(-1)
	for _, root := range nodes {
		if hld.lid[root] < minDfn {
			minDfn = hld.lid[root]
		}
		if hld.lid[root] > maxDfn {
			maxDfn = hld.lid[root]
		}
	}
	u, v := hld.idToNode[minDfn], hld.idToNode[maxDfn]
	return hld.LCA(u, v)
}

func (hld *LCA32) LCA(u, v int32) int32 {
	for {
		if hld.lid[u] > hld.lid[v] {
			u, v = v, u
		}
		if hld.top[u] == hld.top[v] {
			return u
		}
		v = hld.Parent[hld.top[v]]
	}
}

func (hld *LCA32) Dist(u, v int32) int32 {
	return hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)]
}

func (hld *LCA32) EnumeratePathDecomposition(u, v int32, vertex bool, f func(start, end int32)) {
	for {
		if hld.top[u] == hld.top[v] {
			break
		}
		if hld.lid[u] < hld.lid[v] {
			a, b := hld.lid[hld.top[v]], hld.lid[v]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			v = hld.Parent[hld.top[v]]
		} else {
			a, b := hld.lid[u], hld.lid[hld.top[u]]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			u = hld.Parent[hld.top[u]]
		}
	}

	edgeInt := int32(1)
	if vertex {
		edgeInt = 0
	}

	if hld.lid[u] < hld.lid[v] {
		a, b := hld.lid[u]+edgeInt, hld.lid[v]
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	} else if hld.lid[v]+edgeInt <= hld.lid[u] {
		a, b := hld.lid[u], hld.lid[v]+edgeInt
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	}
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
func (hld *LCA32) KthAncestor(root, k int32) int32 {
	if k > hld.Depth[root] {
		return -1
	}
	for {
		u := hld.top[root]
		if hld.lid[root]-k >= hld.lid[u] {
			return hld.idToNode[hld.lid[root]-k]
		}
		k -= hld.lid[root] - hld.lid[u] + 1
		root = hld.Parent[u]
	}
}

// 从 from 节点跳向 to 节点,跳过 step 个节点(0-indexed)
//
//	返回跳到的节点,如果不存在这样的节点,返回-1
func (hld *LCA32) Jump(from, to, step int32) int32 {
	if step == 1 {
		if from == to {
			return -1
		}
		if hld.IsInSubtree(to, from) {
			return hld.KthAncestor(to, (hld.Depth[to] - hld.Depth[from] - 1))
		}
		return hld.Parent[from]
	}
	c := hld.LCA(from, to)
	dac := hld.Depth[from] - hld.Depth[c]
	dbc := hld.Depth[to] - hld.Depth[c]
	if step > dac+dbc {
		return -1
	}
	if step <= dac {
		return hld.KthAncestor(from, step)
	}
	return hld.KthAncestor(to, dac+dbc-step)
}

// child 是否在 root 的子树中 (child和root不能相等)
func (hld *LCA32) IsInSubtree(child, root int32) bool {
	return hld.lid[root] <= hld.lid[child] && hld.lid[child] < hld.rid[root]
}

func (hld *LCA32) _build(cur, pre, dep int32) int32 {
	subSize, heavySize, heavySon := int32(1), int32(0), int32(-1)
	for _, next := range hld.Tree[cur] {
		if next != pre {
			nextSize := hld._build(next, cur, dep+1)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
		}
	}
	hld.Depth[cur] = dep
	hld.heavySon[cur] = heavySon
	hld.Parent[cur] = pre
	return subSize
}

func (hld *LCA32) _markTop(cur, top int32) {
	hld.top[cur] = top
	hld.lid[cur] = hld.dfnId
	hld.idToNode[hld.dfnId] = cur
	hld.dfnId++
	if hld.heavySon[cur] != -1 {
		hld._markTop(hld.heavySon[cur], top)
		for _, next := range hld.Tree[cur] {
			if next != hld.heavySon[cur] && next != hld.Parent[cur] {
				hld._markTop(next, next)
			}
		}
	}
	hld.rid[cur] = hld.dfnId
}

type DoublingSimple struct {
	n   int32
	log int32
	to  []int32
}

func NewDoubling(n int32, maxStep int) *DoublingSimple {
	res := &DoublingSimple{n: n, log: int32(bits.Len(uint(maxStep)))}
	size := n * res.log
	res.to = make([]int32, size)
	for i := int32(0); i < size; i++ {
		res.to[i] = -1
	}
	return res
}

func (d *DoublingSimple) Add(from, to int32) {
	d.to[from] = to
}

func (d *DoublingSimple) Build() {
	n := d.n
	for k := int32(0); k < d.log-1; k++ {
		for v := int32(0); v < n; v++ {
			w := d.to[k*n+v]
			next := (k+1)*n + v
			if w == -1 {
				d.to[next] = -1
				continue
			}
			d.to[next] = d.to[k*n+w]
		}
	}
}

// 求从 `from` 状态开始转移 `step` 次的最终状态的编号。
// 不存在时返回 -1。
func (d *DoublingSimple) Jump(from int32, step int) (to int32) {
	to = from
	for k := int32(0); k < d.log; k++ {
		if to == -1 {
			return
		}
		if step&(1<<k) != 0 {
			to = d.to[k*d.n+to]
		}
	}
	return
}

// 求从 `from` 状态开始转移 `step` 次，满足 `check` 为 `true` 的最大的 `step` 以及最终状态的编号。
func (d *DoublingSimple) MaxStep(from int32, check func(next int32) bool) (step int, to int32) {
	for k := d.log - 1; k >= 0; k-- {
		tmp := d.to[k*d.n+from]
		if tmp == -1 {
			continue
		}
		if check(tmp) {
			step |= 1 << k
			from = tmp
		}
	}
	to = from
	return
}

type FastSet struct {
	n, lg int
	seg   [][]int
	size  int
}

func NewFastSet(n int) *FastSet {
	res := &FastSet{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func NewFastSetFrom(n int, f func(i int) bool) *FastSet {
	res := NewFastSet(n)
	for i := 0; i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
			res.size++
		}
	}
	for h := 0; h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *FastSet) Has(i int) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet) Insert(i int) bool {
	if fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	fs.size++
	return true
}

func (fs *FastSet) Erase(i int) bool {
	if !fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	fs.size--
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == len(cache) {
			break
		}
		d := cache[i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *FastSet) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *FastSet) Enumerate(start, end int, f func(i int)) {
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *FastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *FastSet) Size() int {
	return fs.size
}

func (*FastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*FastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}

type TrieNode struct {
	Children map[byte]*TrieNode
	EndCount int32
}

func NewTrieNode() *TrieNode { return &TrieNode{Children: map[byte]*TrieNode{}} }

type Trie struct{ root *TrieNode }

func NewTrie() *Trie { return &Trie{root: NewTrieNode()} }

func (t *Trie) Insert(n int32, f func(int32) byte) {
	node := t.root
	for i := int32(0); i < n; i++ {
		char := f(i)
		if v, ok := node.Children[char]; ok {
			node = v
		} else {
			newNode := NewTrieNode()
			node.Children[char] = newNode
			node = newNode
		}
	}
	node.EndCount++
}

func (t *Trie) Count(n int32, f func(int32) byte) int32 {
	node := t.root
	for i := int32(0); i < n; i++ {
		char := f(i)
		if v, ok := node.Children[char]; ok {
			node = v
		} else {
			return 0
		}
	}
	return node.EndCount
}

type Sequence = string

type Manacher struct {
	n          int32
	seq        Sequence
	oddRadius  []int32
	evenRadius []int32
	maxOdd1    []int32
	maxOdd2    []int32
	maxEven1   []int32
	maxEven2   []int32
}

func NewManacher(seq Sequence) *Manacher {
	m := &Manacher{
		n:   int32(len(seq)),
		seq: seq,
	}
	return m
}

// 查询切片s[start:end]是否为回文串.
// 空串不为回文串.
func (ma *Manacher) IsPalindrome(start, end int32) bool {
	n := ma.n
	if start < 0 {
		start += n
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += n
	}
	if end > n {
		end = n
	}
	if start >= end {
		return false
	}

	len_ := end - start
	mid := (start + end) >> 1
	if len_&1 == 1 {
		return ma.GetOddRadius()[mid] >= len_>>1+1
	} else {
		return ma.GetEvenRadius()[mid] >= len_>>1
	}
}

// 获取每个中心点的奇回文半径`radius`.
// 回文为`[pos-radius+1:pos+radius]`.
func (ma *Manacher) GetOddRadius() []int32 {
	if ma.oddRadius != nil {
		return ma.oddRadius
	}
	n := ma.n
	ma.oddRadius = make([]int32, n)
	left, right := int32(0), int32(-1)
	for i := int32(0); i < n; i++ {
		var k int32
		if i > right {
			k = 1
		} else {
			k = min32(ma.oddRadius[left+right-i], right-i+1)
		}
		for i-k >= 0 && i+k < n && ma.seq[i-k] == ma.seq[i+k] {
			k++
		}
		ma.oddRadius[i] = k
		k--
		if i+k > right {
			left = i - k
			right = i + k
		}
	}
	return ma.oddRadius
}

// 获取每个中心点的偶回文半径`radius`.
// 回文为`[pos-radius:pos+radius]`.
func (ma *Manacher) GetEvenRadius() []int32 {
	if ma.evenRadius != nil {
		return ma.evenRadius
	}
	n := ma.n
	ma.evenRadius = make([]int32, n)
	left, right := int32(0), int32(-1)
	for i := int32(0); i < n; i++ {
		var k int32
		if i > right {
			k = 0
		} else {
			k = min32(ma.evenRadius[left+right-i+1], right-i+1)
		}
		for i-k-1 >= 0 && i+k < n && ma.seq[i-k-1] == ma.seq[i+k] {
			k++
		}
		ma.evenRadius[i] = k
		k--
		if i+k > right {
			left = i - k - 1
			right = i + k
		}
	}
	return ma.evenRadius
}

// 以s[index]开头的最长奇回文子串的长度.
func (ma *Manacher) GetLongestOddStartsAt(index int32) int32 {
	ma._initOdds()
	return ma.maxOdd1[index]
}

// 以s[index]结尾的最长奇回文子串的长度.
func (ma *Manacher) GetLongestOddEndsAt(index int32) int32 {
	ma._initOdds()
	return ma.maxOdd2[index]
}

// 以s[index]开头的最长偶回文子串的长度.
func (ma *Manacher) GetLongestEvenStartsAt(index int32) int32 {
	ma._initEvens()
	return ma.maxEven1[index]
}

// 以s[index]结尾的最长偶回文子串的长度.
func (ma *Manacher) GetLongestEvenEndsAt(index int32) int32 {
	ma._initEvens()
	return ma.maxEven2[index]
}

func (ma *Manacher) Len() int32 {
	return ma.n
}

func (ma *Manacher) _initOdds() {
	if ma.maxOdd1 != nil {
		return
	}
	n := ma.n
	ma.maxOdd1 = make([]int32, n)
	ma.maxOdd2 = make([]int32, n)
	for i := int32(0); i < n; i++ {
		ma.maxOdd1[i] = 1
		ma.maxOdd2[i] = 1
	}
	for i := int32(0); i < n; i++ {
		radius := ma.GetOddRadius()[i]
		start, end := i-radius+1, i+radius-1
		length := 2*radius - 1
		ma.maxOdd1[start] = max32(ma.maxOdd1[start], length)
		ma.maxOdd2[end] = max32(ma.maxOdd2[end], length)
	}
	for i := int32(0); i < n; i++ {
		if i-1 >= 0 {
			ma.maxOdd1[i] = max32(ma.maxOdd1[i], ma.maxOdd1[i-1]-2)
		}
		if i+1 < n {
			ma.maxOdd2[i] = max32(ma.maxOdd2[i], ma.maxOdd2[i+1]-2)
		}
	}
}

func (ma *Manacher) _initEvens() {
	if ma.maxEven1 != nil {
		return
	}
	n := ma.n
	ma.maxEven1 = make([]int32, n)
	ma.maxEven2 = make([]int32, n)
	for i := int32(0); i < n; i++ {
		radius := ma.GetEvenRadius()[i]
		if radius == 0 {
			continue
		}
		start := i - radius
		end := start + 2*radius - 1
		length := 2 * radius
		ma.maxEven1[start] = max32(ma.maxEven1[start], length)
		ma.maxEven2[end] = max32(ma.maxEven2[end], length)
	}
	for i := int32(0); i < n; i++ {
		if i-1 >= 0 {
			ma.maxEven1[i] = max32(ma.maxEven1[i], ma.maxEven1[i-1]-2)
		}
		if i+1 < n {
			ma.maxEven2[i] = max32(ma.maxEven2[i], ma.maxEven2[i+1]-2)
		}
	}
}

func ReverseString(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, r := range s {
		n--
		runes[n] = r
	}
	return string(runes)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a > b {
		return b
	}
	return a
}
