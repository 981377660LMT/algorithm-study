// 后缀树是一种维护一个字符串所有后缀的数据结构, 同时也是后缀数组矩形区域构成的树.
// 后缀树的本质是后缀trie树的虚树.
//
// 优点：适合处理字典序问题(按照dfs序遍历后缀树就是按照字典序遍历子串)
// 缺点：不如SAM灵活
//
// https://maspypy.github.io/library/string/suffix_tree.hpp
// https://www.luogu.com.cn/blog/EternalAlexander/xuan-ku-hou-zhui-shu-mo-shu
// https://oi-wiki.org/string/suffix-tree/
// https://www.bilibili.com/video/BV1c741137GD
//
// !   - s = abbbab
//
// !   - suffix array
//
//	ab
//	abbbab
//	b
//	bab
//	bbab
//	bbbab
//
// !   - suffix tree (每个结点是一个startPos等价类)
//
//		 ab(1) --- bbab(2)
//		/
//	""(0)
//		\
//		 b(3) --- ab(4)
//		  \
//		   b(5) --- ab(6)
//		    \
//		     ---- bab(7)
//
// !   - [行1，行2，列1，列2] 区域范围
//
//	0 : [0 6 0 0]
//	1 : [0 2 0 2]
//	2 : [1 2 2 6]
//	3 : [2 6 0 1]
//	4 : [3 4 1 3]
//	5 : [4 6 1 2]
//	6 : [4 5 2 4]
//	7 : [5 6 2 5]
//
// !   - 表格展示
//
//	_________________________
//	|_1_|_1_|___|___|___|___|
//	|_1_|_1_|_2_|_2_|_2_|_2_|
//	|_3_|___|___|___|___|___|
//	|_3_|_4_|_4_|___|___|___|
//	|_3_|_5_|_6_|_6_|___|___|
//	|_3_|_5_|_7_|_7_|_7_|___|
//
// !note:
//
//   - 1. 定义字符串 S 的 后缀 trie 为将 S 的所有后缀插入至 trie 树中得到的字典树。
//     在后缀字典树上，每个叶节点都代表了原串的一个后缀.
//     !每个节点到根的路径都是原串的后缀的前缀，即为原串的一个子串.
//     !而这个节点的子树中叶结点的个数代表了`它是多少个后缀的前缀`，即为它在原串中的出现次数。
//     现在的问题是，这颗字典树的节点数是 O(n^2) 的。
//     观察这棵字典树，发现它有很多节点只有一个子节点，形成了若干条单链。
//     我们可以考虑将这些只有一个子节点的节点压缩起来.
//     记后缀 trie 中所有对应 S 的某个后缀的节点为后缀节点。
//     如果令后缀 trie 中所有拥有多于一个儿子的节点和后缀节点为关键点，定义只保留关键点，
//     !将非关键点形成的链压缩成一条边形成的压缩 trie 树为 后缀树 (Suffix Tree)。
//     与后缀字典树不同的是，后缀树的一条边可能有若干个字符。
//
//   - 2. 后缀树上每一个节点到根的路径都是 S 的一个非空子串。
//
// !  - 3. 后缀树的 DFS 序就是后缀数组(后缀树上dfs是按照字典序遍历子串).
//
//   - 4. 后缀树的一个子树对应到后缀数组上的一个区间。
//
//   - 5. 后缀树上每一个叶节点都唯一地对应着原串的一个后缀，两个叶节点的 LCA 对应最长串是它们对应的后缀的 LCP (最长公共前缀) 。
//
//   - 6. 后缀数组的 height 的结论可以理解为树上若干个节点的 LCA 等于 DFS 序最小的和最大的节点的 LCA.
//
// !  - 7. 算上空后缀的话，一个串的后缀字典树上的叶节点个数为 n+1，
//      而后缀树可以认为是这 n+1 个叶节点的虚树，因此节点数上界为 2n+1.
//      个上界可以由串 aaa⋯aaa 达到。
//
// !  - 8. 后缀树是反串的 SAM 的 parent 树，因为后缀树的一个节点的实质是一个 startPos 等价类.
//     而 SAM 的节点代表的是一个 endPos 等价类.
//
//   - 9. 本质不同子串数等于结点数(带权)之和.

package main

import (
	"bufio"
	"fmt"
	"index/suffixarray"
	"os"
	"reflect"
	"sort"
	"strings"
	"unsafe"
)

func main() {
	demo()

	// cf123d()
	// cf427d()
	// cf802I()
	// cf873F()

	// P3181()
	// p3804()
	// p4341()
	// p5341()

	// yukicoder2361()
	abc280()
}

// 最长的只出现一次的子串长度.
func maxSubstringLength(s string) int {
	_, ranges := SuffixTree(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
	res := -1
	for i := 1; i < len(ranges); i++ {
		rowStart, rowEnd, _, colEnd := ranges[i][0], ranges[i][1], ranges[i][2], ranges[i][3]
		freq := int(rowEnd - rowStart)
		if freq == 1 {
			res = max(res, int(colEnd))
		}
	}
	return res
}

// https://oi-wiki.org/string/suffix-tree/
func demo() {
	// s := "cabab"
	s := "abbbab"
	sa, _, lcp := SuffixArray32(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
	suffixTree, ranges := SuffixTreeFrom(sa, lcp)
	fmt.Println(suffixTree, ranges)
	start, end := RecoverSubstring(sa, 3, 1, 3)
	fmt.Println(s[start:end])
	fmt.Println(LocateSuffixes(int32(len(s)), suffixTree, ranges))
}

// CF123D String
// https://www.luogu.com.cn/problem/CF123D
// !枚举字符串 s 的每一个本质不同的子串 ss ，令 cnt(ss) 为子串 ss 在字符串 s 中出现的个数，求 ∑ cnt(ss)*(cnt(ss)+1)/2
// 建立后缀树，可以得到每个节点对应后缀数组上的 [行1，行2，列1，列2] 矩形区域.
// !(行2-行1) 表示此startPos出现次数, (列2-列1) 表示结点包含的压缩的字符串长度(个数).
func cf123d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	_, ranges := SuffixTree(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
	res := 0
	for i := 1; i < len(ranges); i++ {
		rowStart, rowEnd, colStart, colEnd := ranges[i][0], ranges[i][1], ranges[i][2], ranges[i][3]
		freq, nodeCount := int(rowEnd-rowStart), int(colEnd-colStart)
		res += (freq * (freq + 1) / 2) * nodeCount
	}
	fmt.Fprintln(out, res)
}

// Match & Catch
// https://www.luogu.com.cn/problem/CF427D
// 求两个串的最短公共唯一子串
// 令s12 := s1 + "#" + s2，遍历后缀树，对每个结点统计belong即可.
func cf427d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const INF int32 = 1e9
	var s1, s2 string
	fmt.Fscan(in, &s1, &s2)
	n1 := int32(len(s1))
	s12 := s1 + "#" + s2
	sa, _, height := SuffixArray32(int32(len(s12)), func(i int32) int32 { return int32(s12[i]) })
	tree, ranges := SuffixTreeFrom(sa, height)

	res := INF
	var dfs func(cur int32)
	dfs = func(cur int32) {
		rowStart, rowEnd, colStart, colEnd := ranges[cur][0], ranges[cur][1], ranges[cur][2], ranges[cur][3]
		freq, nodeCount := rowEnd-rowStart, colEnd-colStart
		if nodeCount > 0 && freq == 2 {
			belong1, belong2 := 0, 0
			for i := rowStart; i < rowEnd; i++ {
				if sa[i] < n1 {
					belong1++
				} else if sa[i] > n1 {
					belong2++
				}
			}
			if belong1 == 1 && belong2 == 1 {
				minLength := colStart + 1
				res = min32(res, minLength)
			}
		}
		for _, v := range tree[cur] {
			dfs(v)
		}
	}
	dfs(0)

	if res == INF {
		res = -1
	}
	fmt.Fprintln(out, res)
}

// Fake News (hard)
// https://www.luogu.com.cn/problem/CF802I
// 给出 s，求所有 s 的本质不同子串 ss 在 s 中的出现次数平方和，重复的子串只算一次。
// 同cf123d
func cf802I() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func() {
		var s string
		fmt.Fscan(in, &s)
		_, ranges := SuffixTree(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
		res := 0
		for i := 1; i < len(ranges); i++ {
			rowStart, rowEnd, colStart, colEnd := ranges[i][0], ranges[i][1], ranges[i][2], ranges[i][3]
			freq, nodeCount := int(rowEnd-rowStart), int(colEnd-colStart)
			res += (freq * freq) * nodeCount
		}
		fmt.Fprintln(out, res)
	}

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		solve()
	}
}

// Forbidden Indices
// https://codeforces.com/problemset/problem/873/F
// 给出一个字符串 s，一个 01 串，长度均为 n（n≤2e5）.
// !设 ss 为 s 的一个子串，求 `ss长度*不在被禁止位置结束的子串ss出现次数` 的最大值。
//
// 取反串，限制条件就变成了`不在被禁止位置开始的子串ss出现次数`, 转换成`禁止一些后缀`.
// 建立后缀树即可.
func cf873F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)
	var forbidden string
	fmt.Fscan(in, &forbidden)

	s, forbidden = reverseString(s), reverseString(forbidden)
	sa, _, height := SuffixArray32(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
	ok := make([]bool, n) // 按照sa数组顺序的ok的后缀起点.
	for i := 0; i < n; i++ {
		j := sa[i]
		ok[i] = forbidden[j] == '0'
	}
	okPreSum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		okPreSum[i] = okPreSum[i-1]
		if ok[i-1] {
			okPreSum[i]++
		}
	}

	_, ranges := SuffixTreeFrom(sa, height)
	res := 0
	for i := 1; i < len(ranges); i++ {
		rowStart, rowEnd := ranges[i][0], ranges[i][1]
		freq := okPreSum[rowEnd] - okPreSum[rowStart]
		length := int(ranges[i][3])
		res = max(res, freq*length)
	}
	fmt.Fprintln(out, res)
}

// P3181 [HAOI2016] 找相同字符
// 求两个字符的相同子串对数.两个方案不同当且仅当这两个子串中有一个位置不同。
// https://www.luogu.com.cn/problem/P3181
//
// 令s12 := s1 + "#" + s2，遍历后缀树，对每个结点统计belong即可.
func P3181() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s1, s2 string
	fmt.Fscan(in, &s1, &s2)
	n1 := int32(len(s1))
	s12 := s1 + "#" + s2
	sa, _, height := SuffixArray32(int32(len(s12)), func(i int32) int32 { return int32(s12[i]) })
	tree, ranges := SuffixTreeFrom(sa, height)

	res := 0
	var dfs func(cur int32)
	dfs = func(cur int32) {
		rowStart, rowEnd, colStart, colEnd := ranges[cur][0], ranges[cur][1], ranges[cur][2], ranges[cur][3]
		freq, nodeCount := rowEnd-rowStart, colEnd-colStart
		if nodeCount > 0 && freq >= 2 {
			belong1, belong2 := 0, 0
			for i := rowStart; i < rowEnd; i++ {
				if sa[i] < n1 {
					belong1++
				} else if sa[i] > n1 {
					belong2++
				}
			}
			res += int(belong1) * int(belong2) * int(nodeCount)
		}
		for _, v := range tree[cur] {
			dfs(v)
		}
	}
	dfs(0)
	fmt.Fprintln(out, res)
}

// P3804 【模板】后缀自动机（SAM）
// https://www.luogu.com.cn/problem/P3804
// 请你求出 S 的所有出现次数不为 1 的子串的出现次数乘上该子串长度的最大值
func p3804() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	tree, ranges := SuffixTree(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
	res := 0

	var dfs func(cur int32)
	dfs = func(cur int32) {
		freq, nodeCount := ranges[cur][1]-ranges[cur][0], ranges[cur][3]-ranges[cur][2]
		if nodeCount > 0 && freq > 1 {
			maxLength := int(ranges[cur][3])
			res = max(res, maxLength*int(freq))
		}
		for _, v := range tree[cur] {
			dfs(v)
		}
	}
	dfs(0)

	fmt.Fprintln(out, res)
}

// P4341 [BJWC2010] 外星联络
// https://www.luogu.com.cn/problem/P4341
// 给一个字符串求所以出现次数大于 1 的子串所出现的次数。输出的顺序按对应的子串的字典序排列。
func p4341() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)

	tree, ranges := SuffixTree(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
	var res []int32
	var dfs func(cur int32)
	dfs = func(cur int32) {
		freq, nodeCount := ranges[cur][1]-ranges[cur][0], ranges[cur][3]-ranges[cur][2]
		if freq > 1 {
			for i := int32(0); i < nodeCount; i++ {
				res = append(res, freq)
			}
		}
		for _, v := range tree[cur] {
			dfs(v)
		}
	}
	dfs(0)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// P5341 [TJOI2019] 甲苯先生和大中锋的字符串
// https://www.luogu.com.cn/problem/P5341
// 给定字符串s, 求出现 k 次的子串中出现次数的最多的长度。如果不存在子串出现 k 次，则输出 −1 。
func p5341() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(s string, k int32) int {
		tree, ranges := SuffixTree(int32(len(s)), func(i int32) int32 { return int32(s[i]) })

		lengthCounter := make([]int, len(s)+2)
		add := func(min int32, max int32, val int) {
			lengthCounter[min] += val
			lengthCounter[max+1] -= val
		}
		build := func() {
			for i := 1; i < len(lengthCounter); i++ {
				lengthCounter[i] += lengthCounter[i-1]
			}
		}

		var dfs func(cur int32)
		dfs = func(cur int32) {
			freq, nodeCount := ranges[cur][1]-ranges[cur][0], ranges[cur][3]-ranges[cur][2]
			if nodeCount > 0 && freq == k {
				minLength, maxLength := ranges[cur][2]+1, ranges[cur][3]
				add(minLength, maxLength, 1)
			}
			for _, v := range tree[cur] {
				dfs(v)
			}
		}
		dfs(0)

		build()
		res, maxCount := -1, -1
		for i, v := range lengthCounter {
			if v > 0 && v >= maxCount {
				maxCount = v
				res = i
			}
		}
		return res
	}

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var s string
		var k int32
		fmt.Fscan(in, &s, &k)
		fmt.Fprintln(out, solve(s, k))
	}
}

// https://yukicoder.me/problems/no/2361
// 给定一个长为n的字符串s和q个询问.
// 每个询问给出一个区间[start,end), 问有多少个子串的字典序严格小于s[start:end).
//
// !1.将查询的子串 s[start:end) 表示为"起点为start的后缀的长为(end-start)的一个前缀".
// !2.离线查询，按照后缀起点将查询分组，便于后续按照字典序遍历查询.
// !3.线段树维护sa数组上的查询长度最小值, 每次将查询长度最小的查询取出，保证遍历后缀树节点时可以按照字典序处理查询.
func yukicoder2361() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	var s string
	fmt.Fscan(in, &s)
	queries := make([][2]int32, q)
	for i := range queries {
		fmt.Fscan(in, &queries[i][0], &queries[i][1])
		queries[i][0]--
	}

	sa, rank, height := SuffixArray32(n, func(i int32) int32 { return int32(s[i]) })
	type pair struct{ length, qid int32 } // 长度, 询问id
	queryGroups := make([][]pair, n)      // 按照sa数组下标分组
	for i := range queries {
		start, end := queries[i][0], queries[i][1]
		saIndex := rank[start]
		queryGroups[saIndex] = append(queryGroups[saIndex], pair{length: end - start, qid: int32(i)})
	}
	for _, group := range queryGroups {
		// 长度短的查询排在数组末尾，先取出
		sort.Slice(group, func(i, j int) bool { return group[i].length > group[j].length })
	}
	seg := NewSegmentTree(int(n), func(i int) E { return E{value: INF32, index: -1} }) // !维护每个saIndex对应的查询长度最小值
	updateRMQ := func(saIndex int32) {
		group := queryGroups[saIndex]
		if len(group) == 0 {
			seg.Set(int(saIndex), E{value: INF32, index: -1})
		} else {
			minLength := group[len(group)-1].length
			seg.Set(int(saIndex), E{value: minLength, index: saIndex})
		}
	}
	for i := int32(0); i < n; i++ {
		updateRMQ(i)
	}

	res := make([]int, q)
	suffixTree, ranges := SuffixTreeFrom(sa, height)
	smaller := 0
	var dfs func(cur int32)
	dfs = func(cur int32) {
		rowStart, rowEnd, colStart, colEnd := int(ranges[cur][0]), int(ranges[cur][1]), int(ranges[cur][2]), int(ranges[cur][3])
		freq, nodeCount := rowEnd-rowStart, colEnd-colStart
		minLength, maxLength := colStart+1, colEnd

		// !按照字典序取出存在于当前结点(矩形)内部的所有查询
		for {
			item := seg.Query(rowStart, rowEnd)
			queryLength, saIndex := int(item.value), item.index
			if queryLength > maxLength {
				break
			}
			group := &queryGroups[saIndex]
			qid := (*group)[len(*group)-1].qid
			*group = (*group)[:len(*group)-1]
			updateRMQ(saIndex)

			res[qid] = smaller + freq*(queryLength-minLength) // 整个矩形区域的子串个数
		}

		smaller += freq * nodeCount
		for _, next := range suffixTree[cur] {
			dfs(next)
		}
	}
	dfs(0)

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// https://www.luogu.com.cn/problem/AT_abc280_h （多个字符串的第k小子串）
// !给定n个字符串和q个查询,每次询问第 k 小子串,并输出它属于哪个字符串和它在原字符串中的位置.
// 查询的k从1开始递增.
//
// 离线所有询问，按照字典序遍历，走到排名对应的节点求得答案即可
func abc280() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const BIG rune = 'z' + 1

	var m int32
	fmt.Fscan(in, &m)
	starts, ends := make([]int32, m), make([]int32, m)
	sb := strings.Builder{}
	for i := int32(0); i < m; i++ {
		var s string
		fmt.Fscan(in, &s)
		sb.WriteRune(BIG)
		starts[i] = int32(sb.Len())
		sb.WriteString(s)
		ends[i] = int32(sb.Len())
	}

	s := sb.String()
	n := int32(len(s))
	sa, rank, height := SuffixArray32(n, func(i int32) int32 { return int32(s[i]) })

	// belong: 每个后缀所属的字符串编号
	// size: 每个后缀在所属字符串中的长度
	// at: 每个后缀在所属字符串中的起始位置
	belong, size, at := make([]int32, n), make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		belong[i] = -1
		size[i] = -1
		at[i] = -1
	}
	for i := int32(0); i < m; i++ {
		for j := starts[i]; j < ends[i]; j++ {
			p := rank[j]
			belong[p] = i
			size[p] = ends[i] - j
			at[p] = j - starts[i]
		}
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i])
		queries[i]--
	}

	queryPtr := 0
	smaller := 0
	suffixTree, ranges := SuffixTreeFrom(sa, height)
	res := make([][3]int32, 0, q) // (belong,start,end)
	var dfs func(cur int32)
	dfs = func(cur int32) {
		if queryPtr == q {
			return
		}
		rowStart, rowEnd, colStart, colEnd := ranges[cur][0], ranges[cur][1], ranges[cur][2], ranges[cur][3]

		// !只有长度在[1,size[rowStart])的子串才是合法子串(不带分隔符)
		if colStart >= size[rowStart] {
			return
		}
		colEnd = min32(colEnd, size[rowStart]) // 保证不超过字符串长度

		freq, nodeCount := rowEnd-rowStart, colEnd-colStart
		curCount := int(freq) * int(nodeCount)
		for ; queryPtr < q; queryPtr++ {
			// 在这个区域中寻找字典序第k小的子串(0<=k<curCount)
			k := queries[queryPtr] - smaller
			if k >= curCount {
				break
			}
			div := k / int(freq)
			bid := belong[rowStart]
			start := at[rowStart]
			end := start + colStart + int32(div) + 1
			res = append(res, [3]int32{bid, start, end})
		}

		smaller += curCount
		for _, next := range suffixTree[cur] {
			dfs(next)
		}
	}
	dfs(0)

	for _, v := range res {
		fmt.Fprintln(out, v[0]+1, v[1]+1, v[2])
	}
}

// directTree: 后缀树, 从 0 开始编号, 0 结点为虚拟根节点.
// ranges: 每个结点对应后缀数组上的 [行1，行2，列1，列2] 矩形区域.
// !(行2-行1) 表示此startPos出现次数, (列2-列1) 表示结点包含的压缩的字符串长度(个数).
// 对应后缀sa编号: [rowStart, rowEnd)
// 对应字符串长度: [colStart+1, colEnd+1)
func SuffixTree(n int32, f func(i int32) int32) (directedTree [][]int32, ranges [][4]int32) {
	sa, _, lcp := SuffixArray32(n, f)
	return SuffixTreeFrom(sa, lcp)
}

// 每个节点为后缀数组上的一个矩形区间.
func SuffixTreeFrom(sa, height []int32) (directedTree [][]int32, ranges [][4]int32) {
	height = height[1:]
	n := int32(len(sa))
	if n == 1 {
		directedTree = make([][]int32, 2)
		directedTree[0] = append(directedTree[0], 1)
		ranges = append(ranges, [4]int32{0, 1, 0, 0})
		ranges = append(ranges, [4]int32{0, 1, 0, 1})
		return
	}

	var edges [][2]int32
	ranges = append(ranges, [4]int32{0, n, 0, 0})
	ct := NewCartesianTreeSimple32(height)

	var dfs func(p, idx int32, h int32)
	dfs = func(p, idx int32, h int32) {
		left, right := ct.Range[idx][0], ct.Range[idx][1]+1
		hh := height[idx]
		if h < hh {
			m := int32(len(ranges))
			edges = append(edges, [2]int32{p, m})
			p = m
			ranges = append(ranges, [4]int32{left, right, h, hh})
		}

		if ct.leftChild[idx] == -1 {
			if hh < n-sa[idx] {
				edges = append(edges, [2]int32{p, int32(len(ranges))})
				ranges = append(ranges, [4]int32{idx, idx + 1, hh, n - sa[idx]})
			}
		} else {
			dfs(p, ct.leftChild[idx], hh)
		}

		if ct.rigthChild[idx] == -1 {
			if hh < n-sa[idx+1] {
				edges = append(edges, [2]int32{p, int32(len(ranges))})
				ranges = append(ranges, [4]int32{idx + 1, idx + 2, hh, n - sa[idx+1]})
			}
		} else {
			dfs(p, ct.rigthChild[idx], hh)
		}
	}

	root := ct.Root
	if height[root] > 0 {
		edges = append(edges, [2]int32{0, 1})
		ranges = append(ranges, [4]int32{0, n, 0, height[root]})
		dfs(1, root, height[root])
	} else {
		dfs(0, root, 0)
	}

	directedTree = make([][]int32, len(ranges))
	for _, e := range edges {
		u, v := e[0], e[1]
		directedTree[u] = append(directedTree[u], v)
	}
	return
}

// 给定后缀数组上的范围 [row, colStart, colEnd]，求出这个区间对应的字符串s[start:end)。
func RecoverSubstring(sa []int32, row int32, colStart, colEnd int32) (start, end int32) {
	start = sa[row] + colStart
	end = sa[row] + colEnd
	return
}

func LocateSuffixes(n int32, suffixTree [][]int32, ranges [][4]int32) (pos []int32) {
	pos = make([]int32, n)
	var dfs func(int32)
	dfs = func(cur int32) {
		isLeaf := len(suffixTree[cur]) == 0
		if isLeaf {
			rowStart, rowEnd := ranges[cur][0], ranges[cur][1]
			fmt.Println(cur, rowStart, rowEnd, "*")
			for i := rowStart; i < rowEnd; i++ {
				pos[i] = cur
			}
		} else {
			for _, v := range suffixTree[cur] {
				dfs(v)
			}
		}
	}
	dfs(0)
	return
}

func SuffixArray32(n int32, f func(i int32) int32) (sa, rank, height []int32) {
	s := make([]byte, 0, n*4)
	for i := int32(0); i < n; i++ {
		v := f(i)
		s = append(s, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
	_sa := *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New(s)).Elem().FieldByName("sa").Field(0).UnsafeAddr()))
	sa = make([]int32, 0, n)
	for _, v := range _sa {
		if v&3 == 0 {
			sa = append(sa, v>>2)
		}
	}
	rank = make([]int32, n)
	for i := int32(0); i < n; i++ {
		rank[sa[i]] = i
	}
	height = make([]int32, n)
	h := int32(0)
	for i := int32(0); i < n; i++ {
		rk := rank[i]
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := sa[rk-1]; i+h < n && j+h < n && f(i+h) == f(j+h); h++ {
			}
		}
		height[rk] = h
	}
	return
}

type CartesianTreeSimple32 struct {
	// ![left, right) 每个元素作为最大/最小值时的左右边界.
	//  左侧为严格扩展, 右侧为非严格扩展.
	//  例如: [2, 1, 1, 5] => [[0 1] [0 4] [2 4] [3 4]]
	Range                         [][2]int32
	Root                          int32
	n                             int32
	nums                          []int32
	leftChild, rigthChild, parent []int32
}

// min
func NewCartesianTreeSimple32(nums []int32) *CartesianTreeSimple32 {
	res := &CartesianTreeSimple32{}
	n := int32(len(nums))
	Range := make([][2]int32, n)
	lch := make([]int32, n)
	rch := make([]int32, n)
	par := make([]int32, n)

	for i := int32(0); i < n; i++ {
		Range[i] = [2]int32{-1, -1}
		lch[i] = -1
		rch[i] = -1
		par[i] = -1
	}

	res.n = n
	res.nums = nums
	res.Range = Range
	res.leftChild = lch
	res.rigthChild = rch
	res.parent = par

	if n == 1 {
		res.Range[0] = [2]int32{0, 1}
		return res
	}

	less := func(i, j int32) bool {
		return (nums[i] < nums[j]) || (nums[i] == nums[j] && i < j)
	}

	stack := make([]int32, 0)
	for i := int32(0); i < n; i++ {
		for len(stack) > 0 && less(i, stack[len(stack)-1]) {
			res.leftChild[i] = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		res.Range[i][0] = 0
		if len(stack) > 0 {
			res.Range[i][0] = stack[len(stack)-1] + 1
		}
		stack = append(stack, i)
	}

	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && less(i, stack[len(stack)-1]) {
			res.rigthChild[i] = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		res.Range[i][1] = n
		if len(stack) > 0 {
			res.Range[i][1] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}

	for i := int32(0); i < n; i++ {
		if res.leftChild[i] != -1 {
			res.parent[res.leftChild[i]] = i
		}
		if res.rigthChild[i] != -1 {
			res.parent[res.rigthChild[i]] = i
		}
	}
	for i := int32(0); i < n; i++ {
		if res.parent[i] == -1 {
			res.Root = i
		}
	}

	return res
}

// PointSetRangeMinIndex
const INF32 int32 = 1e9 + 10

type E = struct{ value, index int32 }

func (*SegmentTree) e() E {
	return E{value: INF32, index: -1}
}
func (*SegmentTree) op(a, b E) E {
	if a.value < b.value {
		return a
	}
	if a.value > b.value {
		return b
	}
	if a.index < b.index {
		return a
	}
	return b
}

type SegmentTree struct {
	n, size int
	seg     []E
}

func NewSegmentTree(n int, f func(int) E) *SegmentTree {
	res := &SegmentTree{}
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
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
func NewSegmentTreeFrom(leaves []E) *SegmentTree {
	res := &SegmentTree{}
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
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
func (st *SegmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTree) Update(index int, value E) {
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
func (st *SegmentTree) Query(start, end int) E {
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
func (st *SegmentTree) QueryAll() E { return st.seg[1] }
func (st *SegmentTree) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree) MaxRight(left int, predicate func(E) bool) int {
	if left == st.n {
		return st.n
	}
	left += st.size
	res := st.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(st.op(res, st.seg[left])) {
			for left < st.size {
				left <<= 1
				if tmp := st.op(res, st.seg[left]); predicate(tmp) {
					res = tmp
					left++
				}
			}
			return left - st.size
		}
		res = st.op(res, st.seg[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return st.n
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree) MinLeft(right int, predicate func(E) bool) int {
	if right == 0 {
		return 0
	}
	right += st.size
	res := st.e()
	for {
		right--
		for right > 1 && right&1 == 1 {
			right >>= 1
		}
		if !predicate(st.op(st.seg[right], res)) {
			for right < st.size {
				right = right<<1 | 1
				if tmp := st.op(st.seg[right], res); predicate(tmp) {
					res = tmp
					right--
				}
			}
			return right + 1 - st.size
		}
		res = st.op(st.seg[right], res)
		if right&-right == right {
			break
		}
	}
	return 0
}

// 返回树压缩后保留的节点编号和新的树.
// !新的树保留了原树的边权.
func CompressTree32(rawTree *Tree32, nodes []int32, directed bool) (rawId []int32, newTree *Tree32) {
	rawId = append(nodes[:0:0], nodes...)
	sort.Slice(rawId, func(i, j int) bool { return rawTree.LID[rawId[i]] < rawTree.LID[rawId[j]] })
	n := int32(len(rawId))
	for i := int32(0); i < n; i++ {
		j := i + 1
		if j == n {
			j = 0
		}
		rawId = append(rawId, rawTree.LCA(rawId[i], rawId[j]))
	}
	// remainNodes = append(remainNodes, rawTree.IdToNode[0])
	sort.Slice(rawId, func(i, j int) bool { return rawTree.LID[rawId[i]] < rawTree.LID[rawId[j]] })

	unique := func(a []int32) []int32 {
		visited := make(map[int32]struct{})
		newNums := []int32{}
		for _, v := range a {
			if _, ok := visited[v]; !ok {
				visited[v] = struct{}{}
				newNums = append(newNums, v)
			}
		}
		return newNums
	}

	rawId = unique(rawId)
	n = int32(len(rawId))
	newTree = NewTree32(n)

	stack := []int32{0}
	for i := int32(1); i < n; i++ {
		for {
			p := rawId[stack[len(stack)-1]]
			v := rawId[i]
			if rawTree.IsInSubtree(v, p) {
				break
			}
			stack = stack[:len(stack)-1]
		}
		newTree.AddDirectedEdge(stack[len(stack)-1], i)
		if !directed {
			newTree.AddDirectedEdge(i, stack[len(stack)-1])
		}
		stack = append(stack, i)
	}
	newTree.Build(0)
	return
}

type Tree32 struct {
	Tree          [][]int32 // (next)
	Depth         []int32
	Parent        []int32
	LID, RID      []int32 // 欧拉序[in,out)
	IdToNode      []int32
	top, heavySon []int32
	timer         int32
}

func NewTree32From(tree [][]int32) *Tree32 {
	n := int32(len(tree))
	lid := make([]int32, n)
	rid := make([]int32, n)
	IdToNode := make([]int32, n)
	top := make([]int32, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int32, n)    // 深度
	parent := make([]int32, n)   // 父结点
	heavySon := make([]int32, n) // 重儿子
	for i := range parent {
		parent[i] = -1
	}

	return &Tree32{
		Tree:     tree,
		Depth:    depth,
		Parent:   parent,
		LID:      lid,
		RID:      rid,
		IdToNode: IdToNode,
		top:      top,
		heavySon: heavySon,
	}
}

func NewTree32(n int32) *Tree32 {
	return NewTree32From(make([][]int32, n))
}

// 添加无向边 u-v.
func (tree *Tree32) AddEdge(u, v int32) {
	tree.Tree[u] = append(tree.Tree[u], v)
	tree.Tree[v] = append(tree.Tree[v], u)
}

// 添加有向边 u->v.
func (tree *Tree32) AddDirectedEdge(u, v int32) {
	tree.Tree[u] = append(tree.Tree[u], v)
}

// root:0-based
//
//	当root设为-1时，会从0开始遍历未访问过的连通分量
func (tree *Tree32) Build(root int32) {
	if root != -1 {
		tree.build(root, -1, 0)
		tree.markTop(root, root)
	} else {
		for i := int32(0); i < int32(len(tree.Tree)); i++ {
			if tree.Parent[i] == -1 {
				tree.build(i, -1, 0)
				tree.markTop(i, i)
			}
		}
	}
}

// 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
func (tree *Tree32) Id(root int32) (int32, int32) {
	return tree.LID[root], tree.RID[root]
}

// 返回返回边 u-v 对应的 欧拉序起点编号, 1 <= eid <= n-1., 0-indexed.
func (tree *Tree32) Eid(u, v int32) int32 {
	if tree.LID[u] > tree.LID[v] {
		return tree.LID[u]
	}
	return tree.LID[v]
}

func (tree *Tree32) LCA(u, v int32) int32 {
	for {
		if tree.LID[u] > tree.LID[v] {
			u, v = v, u
		}
		if tree.top[u] == tree.top[v] {
			return u
		}
		v = tree.Parent[tree.top[v]]
	}
}

func (tree *Tree32) RootedLCA(u, v int32, root int32) int32 {
	return tree.LCA(u, v) ^ tree.LCA(u, root) ^ tree.LCA(v, root)
}

func (tree *Tree32) RootedParent(u int32, root int32) int32 {
	return tree.Jump(u, root, 1)
}

func (tree *Tree32) Dist(u, v int32) int32 {
	return tree.Depth[u] + tree.Depth[v] - 2*tree.Depth[tree.LCA(u, v)]
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
//	kthAncestor(root,0) == root
func (tree *Tree32) KthAncestor(root, k int32) int32 {
	if k > tree.Depth[root] {
		return -1
	}
	for {
		u := tree.top[root]
		if tree.LID[root]-k >= tree.LID[u] {
			return tree.IdToNode[tree.LID[root]-k]
		}
		k -= tree.LID[root] - tree.LID[u] + 1
		root = tree.Parent[u]
	}
}

// 从 from 节点跳向 to 节点,跳过 step 个节点(0-indexed)
//
//	返回跳到的节点,如果不存在这样的节点,返回-1
func (tree *Tree32) Jump(from, to, step int32) int32 {
	if step == 1 {
		if from == to {
			return -1
		}
		if tree.IsInSubtree(to, from) {
			return tree.KthAncestor(to, tree.Depth[to]-tree.Depth[from]-1)
		}
		return tree.Parent[from]
	}
	c := tree.LCA(from, to)
	dac := tree.Depth[from] - tree.Depth[c]
	dbc := tree.Depth[to] - tree.Depth[c]
	if step > dac+dbc {
		return -1
	}
	if step <= dac {
		return tree.KthAncestor(from, step)
	}
	return tree.KthAncestor(to, dac+dbc-step)
}

func (tree *Tree32) CollectChild(root int32) []int32 {
	res := []int32{}
	for _, next := range tree.Tree[root] {
		if next != tree.Parent[root] {
			res = append(res, next)
		}
	}
	return res
}

// 返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
//
//	!eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
func (tree *Tree32) GetPathDecomposition(u, v int32, vertex bool) [][2]int32 {
	up, down := [][2]int32{}, [][2]int32{}
	for {
		if tree.top[u] == tree.top[v] {
			break
		}
		if tree.LID[u] < tree.LID[v] {
			down = append(down, [2]int32{tree.LID[tree.top[v]], tree.LID[v]})
			v = tree.Parent[tree.top[v]]
		} else {
			up = append(up, [2]int32{tree.LID[u], tree.LID[tree.top[u]]})
			u = tree.Parent[tree.top[u]]
		}
	}
	edgeInt := int32(1)
	if vertex {
		edgeInt = 0
	}
	if tree.LID[u] < tree.LID[v] {
		down = append(down, [2]int32{tree.LID[u] + edgeInt, tree.LID[v]})
	} else if tree.LID[v]+edgeInt <= tree.LID[u] {
		up = append(up, [2]int32{tree.LID[u], tree.LID[v] + edgeInt})
	}
	for i := 0; i < len(down)/2; i++ {
		down[i], down[len(down)-1-i] = down[len(down)-1-i], down[i]
	}
	return append(up, down...)
}

// 遍历路径上的 `[起点,终点)` 欧拉序 `左闭右开` 区间.
func (tree *Tree32) EnumeratePathDecomposition(u, v int32, vertex bool, f func(start, end int32)) {
	for {
		if tree.top[u] == tree.top[v] {
			break
		}
		if tree.LID[u] < tree.LID[v] {
			a, b := tree.LID[tree.top[v]], tree.LID[v]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			v = tree.Parent[tree.top[v]]
		} else {
			a, b := tree.LID[u], tree.LID[tree.top[u]]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			u = tree.Parent[tree.top[u]]
		}
	}

	edgeInt := int32(1)
	if vertex {
		edgeInt = 0
	}

	if tree.LID[u] < tree.LID[v] {
		a, b := tree.LID[u]+edgeInt, tree.LID[v]
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	} else if tree.LID[v]+edgeInt <= tree.LID[u] {
		a, b := tree.LID[u], tree.LID[v]+edgeInt
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	}
}

func (tree *Tree32) GetPath(u, v int32) []int32 {
	res := []int32{}
	composition := tree.GetPathDecomposition(u, v, true)
	for _, e := range composition {
		a, b := e[0], e[1]
		if a <= b {
			for i := a; i <= b; i++ {
				res = append(res, tree.IdToNode[i])
			}
		} else {
			for i := a; i >= b; i-- {
				res = append(res, tree.IdToNode[i])
			}
		}
	}
	return res
}

// 以root为根时,结点v的子树大小.
func (tree *Tree32) SubSize(v, root int32) int32 {
	if root == -1 {
		return tree.RID[v] - tree.LID[v]
	}
	if v == root {
		return int32(len(tree.Tree))
	}
	x := tree.Jump(v, root, 1)
	if tree.IsInSubtree(v, x) {
		return tree.RID[v] - tree.LID[v]
	}
	return int32(len(tree.Tree)) - tree.RID[x] + tree.LID[x]
}

// child 是否在 root 的子树中 (child和root不能相等)
func (tree *Tree32) IsInSubtree(child, root int32) bool {
	return tree.LID[root] <= tree.LID[child] && tree.LID[child] < tree.RID[root]
}

// 寻找以 start 为 top 的重链 ,heavyPath[-1] 即为重链底端节点.
func (tree *Tree32) GetHeavyPath(start int32) []int32 {
	heavyPath := []int32{start}
	cur := start
	for tree.heavySon[cur] != -1 {
		cur = tree.heavySon[cur]
		heavyPath = append(heavyPath, cur)
	}
	return heavyPath
}

// 结点v的重儿子.如果没有重儿子,返回-1.
func (tree *Tree32) GetHeavyChild(v int32) int32 {
	k := tree.LID[v] + 1
	if k == int32(len(tree.Tree)) {
		return -1
	}
	w := tree.IdToNode[k]
	if tree.Parent[w] == v {
		return w
	}
	return -1
}

func (tree *Tree32) ELID(u int32) int32 {
	return 2*tree.LID[u] - tree.Depth[u]
}

func (tree *Tree32) ERID(u int32) int32 {
	return 2*tree.RID[u] - tree.Depth[u] - 1
}

func (tree *Tree32) build(cur, pre, dep int32) int32 {
	subSize, heavySize, heavySon := int32(1), int32(0), int32(-1)
	for _, next := range tree.Tree[cur] {
		if next != pre {
			nextSize := tree.build(next, cur, dep+1)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
		}
	}
	tree.Depth[cur] = dep
	tree.heavySon[cur] = heavySon
	tree.Parent[cur] = pre
	return subSize
}

func (tree *Tree32) markTop(cur, top int32) {
	tree.top[cur] = top
	tree.LID[cur] = tree.timer
	tree.IdToNode[tree.timer] = cur
	tree.timer++
	heavySon := tree.heavySon[cur]
	if heavySon != -1 {
		tree.markTop(heavySon, top)
		for _, next := range tree.Tree[cur] {
			if next != heavySon && next != tree.Parent[cur] {
				tree.markTop(next, next)
			}
		}
	}
	tree.RID[cur] = tree.timer
}

func reverseString(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, r := range s {
		n--
		runes[n] = r
	}
	return string(runes)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
