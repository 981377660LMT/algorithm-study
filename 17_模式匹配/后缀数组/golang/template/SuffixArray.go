// https://github.com/EndlessCheng/codeforces-go/blob/646deb927bbe089f60fc0e9f43d1729a97399e5f/copypasta/strings.go#L556
// https://visualgo.net/zh/suffixarray
// !常用分隔符 #(35) $(36) _(95) |(124)
// SA-IS 与 DC3 的效率对比 https://riteme.site/blog/2016-6-19/sais.html#5
// 注：Go1.13 开始使用 SA-IS 算法
//
// - 支持sa/rank/lcp
// - 比较任意两个子串的字典序
// - 求出任意两个子串的最长公共前缀(lcp)

//  sa : 排第几的后缀是谁.
//  rank : 每个后缀排第几.
//  lcp : 排名相邻的两个后缀的最长公共前缀.
// 	lcp[0] = 0
// 	lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
//
//  "banana" -> sa: [5 3 1 0 4 2], rank: [3 2 5 1 4 0], lcp: [0 1 3 0 0 2]
//
//  !lcp(sa[i],sa[j]) = min(height[i+1..j])
//
// !api:
//  func NewSuffixArray(ords []int) *SuffixArray
//  func NewSuffixArrayWithString(s string) *SuffixArray
//  func (suf *SuffixArray) Lcp(a, b int, c, d int) int
//  func (suf *SuffixArray) CompareSubstr(a, b int, c, d int) int
//  func (suf *SuffixArray) LcpRange(left int, k int) (start, end int)
//	func (suf *SuffixArray) Count(start, end int) int
//  func GetSA(ords []int) (sa []int)
//  func UseSA(ords []int) (sa, rank, lcp []int)

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"math/bits"
)

func main() {
	// CF19C()
	// CF126()
	// CF432D()
	CF822E()

	// abc213f()
	// abc272f()

	// p2178()
	// P2870()
	// P3763()
	// P3804()
	// P4051()
	// P4248()

	// 重复次数最多的连续重复子串()
	// fmt.Println(所有后缀LCP之和("abaab"))

	// test()
}

// Deletion of Repeats (删除连续重复子串)
// https://www.luogu.com.cn/problem/CF19C
// 给定一个长度为n的数组(字符串)。每个字母最多出现10次.
// 要进行若干次删除操作，每次把字符串中的连续重复子串(repeat)形如 XX删除前一个 X 及其前面的所有字符。
// 优先删除最短的重复子串，如果有多个就选择最左边的。
// 输出最后的字符串。
//
// 1. 离散化
// !2. 连续重复子串 <=> lcp(sa[i],sa[j]) >= j-i
func CF19C() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	ords := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &ords[i])
	}
	D := NewDictionary[int]()
	for i, v := range ords {
		ords[i] = D.Id(v)
	}

	mp := make(map[int][]int)
	for i, v := range ords {
		mp[v] = append(mp[v], i)
	}
	S := NewSuffixArray(ords)
	type pair = struct{ length, start int32 }
	pq := NewHeap[pair](
		func(a, b pair) bool {
			if a.length != b.length {
				return a.length < b.length
			}
			return a.start < b.start
		},
		nil,
	)

	for _, indexes := range mp {
		for i, v1 := range indexes {
			for j := i + 1; j < len(indexes); j++ {
				v2 := indexes[j]
				if lcp := S.Lcp(v1, n, v2, n); lcp >= v2-v1 {
					pq.Push(pair{length: int32(v2 - v1), start: int32(v1)})
				}
			}
		}
	}

	start := int32(0)
	for pq.Len() > 0 {
		p := pq.Pop()
		if start <= p.start {
			start = p.start + p.length
		}
	}

	res := ords[start:]
	for i := 0; i < len(res); i++ {
		res[i] = D.Value(res[i])
	}
	fmt.Fprintln(out, len(res))
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

// Password
// https://www.luogu.com.cn/problem/CF126B
// 给你一个字符串S（|S|<=1e6）.
// !找到既是S前缀又是S后缀又在S中间出现过（既不是S前缀又不是S后缀）的子串.
// 如果不存在输出“Just a legend”.
//
// !是前缀：rank[0] 对应后缀的前缀即为整个串的前缀.
// !是后缀：说明这个后缀的字典序需要小于整个串的字典序，即排名小于rank[0].
// !在中间出现过：
func CF126() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const INF int = 1e18

	var s string
	fmt.Fscan(in, &s)

	S := NewSuffixArrayFromString(s)
	sa, rank, height := S.Sa, S.Rank, S.Height
	k := rank[0] // 整个串的排名

	n := len(s)
	resEnd := 0
	lcp := INF
	for i := k - 1; i >= 0; i-- { // 前缀的字典序一定小于整个串的字典序, 所以向上找
		lcp = min(lcp, height[i+1])
		length := n - sa[i]
		if length == lcp { // is border
			ok1 := k-i >= 2                      // 上方存在第三个相同的子串
			ok2 := k+1 < n && height[k+1] >= lcp // 下方存在第三个相同的子串
			if ok1 || ok2 {
				resEnd = max(resEnd, length)
			}
		}
	}

	if resEnd == 0 {
		fmt.Fprintln(out, "Just a legend")
	} else {
		fmt.Fprintln(out, s[:resEnd])
	}
}

// Prefixes and Suffixes
// https://www.luogu.com.cn/problem/CF432D
// “完美子串”定义为既是前缀又是后缀的子串.
// 给定一个字符串，求这些串的(长度，个数).
// 是前缀的后缀 <=> lcp(后缀x, s) == n - sa[x]
func CF432D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const INF int = 1e18

	var s string
	fmt.Fscan(in, &s)

	S := NewSuffixArrayFromString(s)
	sa, rank, height := S.Sa, S.Rank, S.Height
	k := rank[0] // 整个串的排名

	n := len(s)

	lcp := INF
	res := [][2]int{{n, 1}} // 整个串
	for i := k - 1; i >= 0; i-- {
		lcp = min(lcp, height[i+1])
		length := n - sa[i]
		if length == lcp { // is border (是前缀的后缀)
			start, end := sa[i], sa[i]+length
			count := S.Count(start, end)
			res = append(res, [2]int{length, count})
		}
	}

	fmt.Fprintln(out, len(res))
	for i := len(res) - 1; i >= 0; i-- {
		fmt.Fprintln(out, res[i][0], res[i][1])
	}
}

// CF822E Liar (拼接不相交子串的最少代价)
// https://www.luogu.com.cn/problem/solution/CF822E
// 给定两个字符串s，t，长度分别为n，m。
// 你需要选择s的若干个两两不相交的子串，然后将它们按照原先在s中出现的顺序合并起来，希望得到t。
// !拼接的子串个数能否不超过limit.
//
// n,m<=1e5,limit<=30.
//
// !dp[i][j] 表示 s 的前 i 个字符选出j个子串时，t 最多能匹配到的位置。
// 每次匹配尽量匹配完，即两个字符串两个后缀的lcp长度最优.
// dp[i][j] -> dp[i+1][j] 跳过匹配位置
func CF822E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	var s, t string
	var limit int
	fmt.Fscan(in, &n)
	fmt.Fscan(in, &s)
	fmt.Fscan(in, &m)
	fmt.Fscan(in, &t)
	fmt.Fscan(in, &limit)

	S := NewSuffixArray2FromString(s, t)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, limit+1)
	}
	for i := 0; i < n; i++ {
		for k := 0; k < limit+1; k++ {
			dp[i+1][k] = max(dp[i+1][k], dp[i][k]) // 不选s的这个后缀
			if k == limit {
				break
			}
			// 后缀s[i:]和后缀t[dp[i][k]:]的lcp
			lcp := S.Lcp(i, n, dp[i][k], m)
			dp[i+lcp][k+1] = max(dp[i+lcp][k+1], dp[i][k]+lcp)
		}
	}

	if dp[n][limit] == m {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}

// https://www.luogu.com.cn/problem/P2178
// P2178 [NOI2015] 品酒大会
// https://www.luogu.com.cn/problem/P2178
// 对于每个 0<=i<n，求有多少对后缀满足 len(lcp) ≥ i，以及满足条件的两个后缀的权值乘积的最大值。
//
// !求出height数组后，按 height 从大到小把所有sa数组中相邻后缀合并.用并查集维护.
func p2178() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const INF int = 1e18

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)
	scores := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &scores[i])
	}

	S := NewSuffixArrayFromString(s)
	sa, height := S.Sa, S.Height
	indexes := make([][]int, n)
	for i, h := range height {
		indexes[h] = append(indexes[h], i)
	}
	uf := NewUnionFindArray(n)
	groupInfo := make([][2]int, n) // (min,max)
	for i := 0; i < n; i++ {
		groupInfo[i] = [2]int{scores[sa[i]], scores[sa[i]]}
	}

	res1, res2 := make([]int, n), make([]int, n) // res1: 满足条件的对数, res2: 满足条件的两个后缀的权值乘积的最大值
	curPair, curMax := 0, -INF
	for i := n - 1; i >= 0; i-- {
		for _, j := range indexes[i] {
			if j == 0 {
				continue
			}

			uf.Union(
				j-1, j,
				func(big, small int) {
					bigSize, smallSize := uf.Size(big), uf.Size(small)
					curPair += bigSize * smallSize
					bigMin, smallMin := groupInfo[big][0], groupInfo[small][0]
					bigMax, smallMax := groupInfo[big][1], groupInfo[small][1]
					groupInfo[big][0] = min(bigMin, smallMin)
					groupInfo[big][1] = max(bigMax, smallMax)
					curMax = max(curMax, bigMax*smallMax)
					curMax = max(curMax, bigMin*smallMin)
				},
			)
		}

		res1[i], res2[i] = curPair, curMax
	}

	for i := 0; i < n; i++ {
		if res2[i] == -INF {
			res2[i] = 0
		}
	}

	for i := 0; i < n; i++ {
		fmt.Fprintln(out, res1[i], res2[i])
	}
}

// P2870 [USACO07DEC] Best Cow Line G
// https://www.luogu.com.cn/problem/P2870
// 从字符串首尾取字符最小化字典序
// 每输出 80 个字母需要一个换行。
func P2870() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	ords1, ords2 := make([]int, n), make([]int, n) // 原串,反串
	for i := 0; i < n; i++ {
		var c string
		fmt.Fscan(in, &c)
		v := int(c[0])
		ords1[i], ords2[n-1-i] = v, v
	}
	S := NewSuffixArray2(ords1, ords2)

	ptr1, ptr2 := 0, 0
	count := 0
	for count < n {
		if S.CompareSubstr(ptr1, n, ptr2, n) == -1 {
			fmt.Fprint(out, string(rune(ords1[ptr1])))
			ptr1++
		} else {
			fmt.Fprint(out, string(rune(ords2[ptr2])))
			ptr2++
		}
		count++

		if count%80 == 0 {
			fmt.Fprintln(out)
		}
	}
}

// P3763 [TJOI2017] DNA (允许失配k次的字符串匹配)
// https://www.luogu.com.cn/problem/P3763
// 给定一个长n1的文本串和长n2的模式串，对文本串每个长为n2的子串，问允许失配3次的情况下，有多少个子串与模式串匹配。
//
// 使用lcp加速匹配.
func P3763() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(text, pattern string, k int) (res int) {
		n1, n2 := len(text), len(pattern)
		if n1 < n2 {
			return 0
		}

		S := NewSuffixArray2FromString(text, pattern)
		check := func(start, end int) bool {
			failCount := 0
			pos := start
			for pos < end {
				// 移动到下一个待匹配位置继续匹配
				if text[pos] != pattern[pos-start] {
					failCount++
					if failCount > k {
						return false
					}
					pos++
				} else {
					pos += S.Lcp(pos, end, n2-(end-pos), n2)
				}
			}
			return true
		}

		for i := 0; i+n2 <= n1; i++ {
			if check(i, i+n2) {
				res++
			}
		}
		return
	}

	var T int32
	fmt.Fscan(in, &T)
	for i := int32(0); i < T; i++ {
		var text, pattern string
		fmt.Fscan(in, &text, &pattern)
		fmt.Fprintln(out, solve(text, pattern, 3))
	}
}

// P3804 【模板】后缀自动机（SAM）
// https://www.luogu.com.cn/problem/P3804
// 给定一个长度为 n 的只包含小写字母的字符串 s。
// !对于所有 s 的出现次数不为 1 的子串，设其 value值为该 子串出现的次数 × 该子串的长度。
// 请计算，value 的最大值是多少。
// n <= 1e6
//
// !子串出现次数乘以次数的最大值-直方图最大矩形
// 直方图最大矩形
// 子串长度看成高,lcp范围看成宽
// https://www.acwing.com/solution/content/25201/
func P3804() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	if len(s) <= 1 {
		fmt.Fprintln(out, 0)
		return
	}
	S := NewSuffixArrayFromString(s)
	heights := S.Height
	L, R := GetRange(heights, false, false, false) // 求每个元素作为最小值的影响范围(区间)
	res := 0
	for i := 0; i < len(heights); i++ {
		res = max(res, heights[i]*(R[i]-L[i]+2))
	}
	fmt.Fprintln(out, res)
}

// P4051 [JSOI2007] 字符加密
// https://www.luogu.com.cn/problem/P4051
// 给定一个字符s，将所有轮转字符串按照字典序排序，输出排序后的每个串的最后一个字符.
//
// 变为s+s后排序，输出原串结尾字符即可
func P4051() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var s string
	fmt.Fscan(in, &s)
	n := len(s)
	s += s
	S := NewSuffixArrayFromString(s)
	sa := S.Sa
	for _, v := range sa {
		if v < n {
			fmt.Fprint(out, string(s[v+n-1]))
		}
	}
}

// P4248 [AHOI2013] 差异
// https://www.luogu.com.cn/problem/P4248
// 定义两个字符串 S 和 T 的差异 diff(S,T) 为这两个串的长度之和减去两倍的这两个串的最长公共前缀的长度。
// 求所有的 diff(S,T) 之和。
func P4248() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	n := len(s)
	res := n*(n-1)*(n+1)/2 - 2*所有后缀LCP之和(s)
	fmt.Fprintln(out, res)
}

// P6095 [JSOI2015] 串分割
// https://www.luogu.com.cn/problem/P6095
// 后缀数组一个很好的作用就是可以用它来按照字典序二分字符串
func P6095() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	var s string
	fmt.Fscan(in, &s)
}

// 对0<=i<j<n,求lcp(s[i:],s[j:])之和.
func 所有后缀LCP之和(s string) int {
	S := NewSuffixArrayFromString(s)
	height := S.Height

	n := len(s)
	res := 0
	// !lcp(sa[i],sa[j]) = min(height[i+1..j])
	clampMaxStack := NewClampableStack(false) // 截断最大值的单调栈
	for i := 0; i < n; i++ {
		clampMaxStack.AddAndClamp(height[i])
		res += clampMaxStack.Sum() //  sa[i]与左侧所有后缀的lcp和
	}
	clampMaxStack.Clear()
	for i := n - 1; i >= 0; i-- {
		res += clampMaxStack.Sum() // sa[i]与右侧所有后缀的lcp和
		clampMaxStack.AddAndClamp(height[i])
	}
	return res / 2
}

// Periodic Substring
// 重复次数最多的连续重复子串 (nlogn)
// https://codeforces.com/edu/course/2/lesson/2/5/practice/contest/269656/problem/F
// https://blog.nowcoder.net/n/47821f2464e146ea86d83b224a91d855
// https://blog.nowcoder.net/n/f9c3bcdf807546bd9c8d8cc43df84079
// !枚举连续串的长度 |s|，按照 |s| 对整个串进行分块，对相邻两块的块首计算 LCP(i,i+len)，然后看是否还能再重复一次
func 重复次数最多的连续重复子串() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	n := len(s)
	S := NewSuffixArrayFromString(s)
	res := 1
	for len_ := 1; len_ < n; len_++ {
		for start := 0; start+len_ < n; start += len_ {
			repeatLen := S.Lcp(start, n, start+len_, n)
			repeatCount := repeatLen/len_ + 1
			// 前面可能还有 (len_ - repeatLen%len_) 个字符在第一个重复子串中
			if p := start - (len_ - repeatLen%len_); p >= 0 && S.Lcp(p, n, p+len_, n) >= len_ {
				repeatCount++
			}
			if repeatCount > res {
				res = repeatCount
			}
		}
	}

	fmt.Fprintln(out, res)

}

// F - Common Prefixes-每个后缀与所有后缀的LCP长度和
// https://atcoder.jp/contests/abc213/tasks/abc213_f
// 定义LCP(X,Y)为字符串X,Y的公共前缀长度(LCP)。
// 给定长度为N的字符串S，设S表示从第i个字符开始的S的后缀(就是后缀数组里的那些后缀)。
// !计算出:对于k=1,2,...,N,LCP(Sk,S1)+LCP(Sk,S2)+ +...+LCP(Sk,SN)的值。
// !即求每个后缀与所有后缀的公共前缀长度和。
// n<=1e6
//
// https://blog.hamayanhamayan.com/entry/2021/08/09/010405
func abc213f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)

	S := NewSuffixArrayFromString(s)
	sa, height := S.Sa, S.Height

	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[sa[i]] = n - sa[i]
	}

	// !lcp(sa[i],sa[j]) = min(height[i+1..j])
	clampMaxStack := NewClampableStack(false) // 截断最大值的单调栈
	for i := 0; i < n; i++ {
		clampMaxStack.AddAndClamp(height[i])
		res[sa[i]] += clampMaxStack.Sum() // sa[i]与左侧所有后缀的lcp和
	}
	clampMaxStack.Clear()
	for i := n - 1; i >= 0; i-- {
		res[sa[i]] += clampMaxStack.Sum() // sa[i]与右侧所有后缀的lcp和
		clampMaxStack.AddAndClamp(height[i])
	}

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// F - Two Strings(所有循环串的比较计数)
// https://atcoder.jp/contests/abc272/tasks/abc272_f
// 给定两个长为n的字符串s和t
// 问s和t的所有轮转的子串中 s的轮转子串有多少个字典序 <= t的轮转子串
//
// 技巧:
// 需要一起比较s和t的所有轮转字串的字典序
// !构造一个新的字符串 s+s+'#'+t+t+'|'
// (注意题目要的是小于等于, 这样保证两个字符串在比较完长度为n后S后面的#小于T中任意一个字符。)
// !后缀数组求出每个串的rank
// !然后在t的rank中 用s的每个子串rank二分出t中的pos
func abc272f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s, t string
	fmt.Fscan(in, &s, &t)

	SMALL, BIG := "#", "|"
	sstt := s + s + SMALL + t + t + BIG
	S := NewSuffixArrayFromString(sstt)
	rank := S.Rank
	sRank, tRank := rank[:n], rank[2*n+1:2*n+1+n]
	sort.Ints(tRank)
	res := 0
	for _, r := range sRank {
		res += n - sort.SearchInts(tRank, r)
	}
	fmt.Fprintln(out, res)
}

// !不同子串长度之和
// 枚举每个后缀，计算前缀总数，再减掉重复
func diffSum(s string) int {
	n := len(s)
	ords := make([]int, n)
	for i, c := range s {
		ords[i] = int(c)
	}
	_, _, height := UseSA(ords)
	res := n * (n + 1) * (n + 2) / 6 // 所有子串长度1到n的平方和
	for _, h := range height {
		res -= h * (h + 1) / 2
	}
	return res
}

// https://leetcode.cn/problems/largest-merge-of-two-strings/
// 1754. 构造字典序最大的合并字符串
// 从字符串首尾取字符最小化字典序
func largestMerge(word1 string, word2 string) string {
	ords1, ords2 := make([]int, len(word1)), make([]int, len(word2))
	for i, c := range word1 {
		ords1[i] = int(c)
	}
	for i, c := range word2 {
		ords2[i] = int(c)
	}
	S := NewSuffixArray2(ords1, ords2)

	n1, n2 := len(word1), len(word2)
	sb := strings.Builder{}

	i, j := 0, 0
	for i < len(word1) && j < len(word2) {
		if S.CompareSubstr(i, n1, j, n2) == 1 {
			sb.WriteByte(word1[i])
			i++
		} else {
			sb.WriteByte(word2[j])
			j++
		}
	}

	sb.WriteString(word1[i:])
	sb.WriteString(word2[j:])

	return sb.String()
}

// 2223. 构造字符串的总得分和
// https://leetcode.cn/problems/sum-of-scores-of-built-strings/
func sumScores(s string) int64 {
	sa := NewSuffixArrayFromString(s)
	n := len(s)
	res := 0
	for i := 0; i < n; i++ {
		res += sa.Lcp(0, n, i, n)
	}
	return int64(res)
}

func test() {
	n := int(1000)
	ords := make([]int, n)
	for i := 1; i < n; i++ {
		ords[i] = i * i
		ords[i] ^= ords[i-1]
	}

	S := NewSuffixArray(ords)
	S2 := NewSuffixArray(ords)
	LcpRange2 := func(left, k int) (start, end int) {
		curRank := S2.Rank[left]
		for i := curRank; i >= 0; i-- {
			sa := S2.Sa[i]
			if S2.Lcp(sa, n, left, n) >= k {
				start = i
			} else {
				break
			}
		}
		for i := curRank; i < n; i++ {
			sa := S2.Sa[i]
			if S2.Lcp(sa, n, left, n) >= k {
				end = i + 1
			} else {
				break
			}
		}
		if start == 0 && end == 0 {
			return -1, -1
		}
		return
	}

	count2 := func(ords []int, start, end int) int {

		target := ords[start:end]
		len_ := end - start
		if len_ == 0 {
			return 0
		}
		res := 0

		for i := 0; i+len_ <= n; i++ {
			curOrds := ords[i : i+len_]
			allOk := true
			for j := 0; j < len_; j++ {
				if curOrds[j] != target[j] {
					allOk = false
					break
				}
			}
			if allOk {
				res++
			}
		}
		return res
	}

	// lcpRange
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			start1, end1 := S.LcpRange(i, j)
			start2, end2 := LcpRange2(i, j)
			if start1 != start2 || end1 != end2 {
				fmt.Println(i, j, start1, end1, start2, end2)
				panic("")
			}
		}
	}

	// count
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			c1 := S.Count(i, j)
			c2 := count2(ords, i, j)
			if c1 != c2 {
				fmt.Println(i, j, c1, c2)
				panic("")
			}
		}
	}

	fmt.Println("pass")
}

func demo() {
	s := "abca"
	ords := make([]int, len(s))
	for i, c := range s {
		ords[i] = int(c)
	}
	sa, rank, height := UseSA(ords)
	fmt.Println(sa, rank, height)
}

type SuffixArray struct {
	Sa     []int // 排名第i的后缀是谁.
	Rank   []int // 后缀s[i:]的排名是多少.
	Height []int // 排名相邻的两个后缀的最长公共前缀.Height[0] = 0,Height[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	Ords   []int
	n      int
	minSt  *LinearRMQ // 维护lcp的最小值
}

// !ord值很大时,需要先离散化.
// !ords[i]>=0.
func NewSuffixArray(ords []int) *SuffixArray {
	ords = append(ords[:0:0], ords...)
	res := &SuffixArray{n: len(ords), Ords: ords}
	sa, rank, lcp := res._useSA(ords)
	res.Sa, res.Rank, res.Height = sa, rank, lcp
	return res
}

func NewSuffixArrayFromString(s string) *SuffixArray {
	ords := make([]int, len(s))
	for i, c := range s {
		ords[i] = int(c)
	}
	return NewSuffixArray(ords)
}

// 求任意两个子串s[a,b)和s[c,d)的最长公共前缀(lcp).
func (suf *SuffixArray) Lcp(a, b int, c, d int) int {
	if a >= b || c >= d {
		return 0
	}
	cand := suf._lcp(a, c)
	return min(cand, min(b-a, d-c))
}

// 比较任意两个子串s[a,b)和s[c,d)的字典序.
//
//	s[a,b) < s[c,d) 返回-1.
//	s[a,b) = s[c,d) 返回0.
//	s[a,b) > s[c,d) 返回1.
func (suf *SuffixArray) CompareSubstr(a, b int, c, d int) int {
	len1, len2 := b-a, d-c
	lcp := suf.Lcp(a, b, c, d)
	if len1 == len2 && lcp >= len1 {
		return 0
	}
	if lcp >= len1 || lcp >= len2 { // 一个是另一个的前缀
		if len1 < len2 {
			return -1
		}
		return 1
	}
	if suf.Rank[a] < suf.Rank[c] {
		return -1
	}
	return 1
}

// 与 s[left:] 的 lcp 大于等于 k 的后缀数组(sa)上的区间.
// 如果不存在,返回(-1,-1).
func (suf *SuffixArray) LcpRange(left int, k int) (start, end int) {
	if k > suf.n-left {
		return -1, -1
	}
	if k == 0 {
		return 0, suf.n
	}
	if suf.minSt == nil {
		suf.minSt = NewLinearRMQ(suf.Height)
	}
	i := suf.Rank[left] + 1
	start = suf.minSt.MinLeft(i, func(e int) bool { return e >= k }) - 1 // 向左找
	end = suf.minSt.MaxRight(i, func(e int) bool { return e >= k })      // 向右找
	return
}

// 查询s[start:end)在s中的出现次数.
func (suf *SuffixArray) Count(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > suf.n {
		end = suf.n
	}
	if start >= end {
		return 0
	}
	a, b := suf.LcpRange(start, end-start)
	return b - a
}

func (suf *SuffixArray) Print(n int, f func(i int) int, sa []int) {
	for _, v := range sa {
		s := make([]int, 0, n-v)
		for i := v; i < n; i++ {
			s = append(s, f(i))
		}
		fmt.Println(s)
	}
}

// 求任意两个后缀s[i:]和s[j:]的最长公共前缀(lcp).
func (suf *SuffixArray) _lcp(i, j int) int {
	if suf.minSt == nil {
		suf.minSt = NewLinearRMQ(suf.Height)
	}
	if i == j {
		return suf.n - i
	}
	r1, r2 := suf.Rank[i], suf.Rank[j]
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	return suf.minSt.Query(r1+1, r2+1)
}

func (suf *SuffixArray) _getSA(ords []int) (sa []int) {
	if len(ords) == 0 {
		return []int{}
	}
	mn := mins(ords)
	for i, x := range ords {
		ords[i] = x - mn + 1
	}
	ords = append(ords, 0)
	n := len(ords)
	m := maxs(ords) + 1
	isS := make([]bool, n)
	isLms := make([]bool, n)
	lms := make([]int, 0, n)
	for i := 0; i < n; i++ {
		isS[i] = true
	}
	for i := n - 2; i > -1; i-- {
		if ords[i] == ords[i+1] {
			isS[i] = isS[i+1]
		} else {
			isS[i] = ords[i] < ords[i+1]
		}
	}
	for i := 1; i < n; i++ {
		isLms[i] = !isS[i-1] && isS[i]
	}
	for i := 0; i < n; i++ {
		if isLms[i] {
			lms = append(lms, i)
		}
	}
	bin := make([]int, m)
	for _, x := range ords {
		bin[x]++
	}

	induce := func() []int {
		sa := make([]int, n)
		for i := 0; i < n; i++ {
			sa[i] = -1
		}

		saIdx := make([]int, m)
		copy(saIdx, bin)
		for i := 0; i < m-1; i++ {
			saIdx[i+1] += saIdx[i]
		}
		for j := len(lms) - 1; j > -1; j-- {
			i := lms[j]
			x := ords[i]
			saIdx[x]--
			sa[saIdx[x]] = i
		}

		copy(saIdx, bin)
		s := 0
		for i := 0; i < m; i++ {
			s, saIdx[i] = s+saIdx[i], s
		}
		for j := 0; j < n; j++ {
			i := sa[j] - 1
			if i < 0 || isS[i] {
				continue
			}
			x := ords[i]
			sa[saIdx[x]] = i
			saIdx[x]++
		}

		copy(saIdx, bin)
		for i := 0; i < m-1; i++ {
			saIdx[i+1] += saIdx[i]
		}
		for j := n - 1; j > -1; j-- {
			i := sa[j] - 1
			if i < 0 || !isS[i] {
				continue
			}
			x := ords[i]
			saIdx[x]--
			sa[saIdx[x]] = i
		}

		return sa
	}

	sa = induce()

	lmsIdx := make([]int, 0, len(sa))
	for _, i := range sa {
		if isLms[i] {
			lmsIdx = append(lmsIdx, i)
		}
	}
	l := len(lmsIdx)
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = -1
	}
	ord := 0
	order[n-1] = ord
	for i := 0; i < l-1; i++ {
		j, k := lmsIdx[i], lmsIdx[i+1]
		for d := 0; d < n; d++ {
			jIsLms, kIsLms := isLms[j+d], isLms[k+d]
			if ords[j+d] != ords[k+d] || jIsLms != kIsLms {
				ord++
				break
			}
			if d > 0 && (jIsLms || kIsLms) {
				break
			}
		}
		order[k] = ord
	}
	b := make([]int, 0, l)
	for _, i := range order {
		if i >= 0 {
			b = append(b, i)
		}
	}
	var lmsOrder []int
	if ord == l-1 {
		lmsOrder = make([]int, l)
		for i, ord := range b {
			lmsOrder[ord] = i
		}
	} else {
		lmsOrder = suf._getSA(b)
	}
	buf := make([]int, len(lms))
	for i, j := range lmsOrder {
		buf[i] = lms[j]
	}
	lms = buf
	return induce()[1:]
}

func (suf *SuffixArray) _useSA(ords []int) (sa, rank, lcp []int) {
	n := len(ords)
	sa = suf._getSA(ords)

	rank = make([]int, n)
	for i := range rank {
		rank[sa[i]] = i
	}

	// !高度数组 lcp 也就是排名相邻的两个后缀的最长公共前缀。
	// lcp[0] = 0
	// lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	lcp = make([]int, n)
	h := 0
	for i, rk := range rank {
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := int(sa[rk-1]); i+h < n && j+h < n && ords[i+h] == ords[j+h]; h++ {
			}
		}
		lcp[rk] = h
	}

	return
}

type LinearRMQ struct {
	n     int
	nums  []int
	small []int
	large [][]int
}

// n: 序列长度.
// less: 入参为两个索引,返回值表示索引i处的值是否小于索引j处的值.
//
//	消除了泛型.
func NewLinearRMQ(nums []int) *LinearRMQ {
	n := len(nums)
	res := &LinearRMQ{n: n, nums: nums}
	stack := make([]int, 0, 64)
	small := make([]int, 0, n)
	var large [][]int
	large = append(large, make([]int, 0, n>>6))
	for i := 0; i < n; i++ {
		for len(stack) > 0 && nums[stack[len(stack)-1]] > nums[i] {
			stack = stack[:len(stack)-1]
		}
		tmp := 0
		if len(stack) > 0 {
			tmp = small[stack[len(stack)-1]]
		}
		small = append(small, tmp|(1<<(i&63)))
		stack = append(stack, i)
		if (i+1)&63 == 0 {
			large[0] = append(large[0], stack[0])
			stack = stack[:0]
		}
	}

	for i := 1; (i << 1) <= n>>6; i <<= 1 {
		csz := n>>6 + 1 - (i << 1)
		v := make([]int, csz)
		for k := 0; k < csz; k++ {
			back := large[len(large)-1]
			v[k] = res._getMin(back[k], back[k+i])
		}
		large = append(large, v)
	}

	res.small = small
	res.large = large
	return res
}

// 查询区间`[start, end)`中的最小值.
func (rmq *LinearRMQ) Query(start, end int) int {
	if start >= end {
		panic(fmt.Sprintf("start(%d) should be less than end(%d)", start, end))
	}
	end--
	left := start>>6 + 1
	right := end >> 6
	if left < right {
		msb := bits.Len64(uint64(right-left)) - 1
		cache := rmq.large[msb]
		i := (left-1)<<6 + bits.TrailingZeros64(uint64(rmq.small[left<<6-1]&(^0<<(start&63))))
		cand1 := rmq._getMin(i, cache[left])
		j := right<<6 + bits.TrailingZeros64(uint64(rmq.small[end]))
		cand2 := rmq._getMin(cache[right-(1<<msb)], j)
		return rmq.nums[rmq._getMin(cand1, cand2)]
	}
	if left == right {
		i := (left-1)<<6 + bits.TrailingZeros64(uint64(rmq.small[left<<6-1]&(^0<<(start&63))))
		j := left<<6 + bits.TrailingZeros64(uint64(rmq.small[end]))
		return rmq.nums[rmq._getMin(i, j)]
	}
	return rmq.nums[right<<6+bits.TrailingZeros64(uint64(rmq.small[end]&(^0<<(start&63))))]
}

func (rmq *LinearRMQ) _getMin(i, j int) int {
	if rmq.nums[i] < rmq.nums[j] {
		return i
	}
	return j
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (st *LinearRMQ) MaxRight(left int, check func(e int) bool) int {
	if left == st.n {
		return st.n
	}
	ok, ng := left, st.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(st.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (st *LinearRMQ) MinLeft(right int, check func(e int) bool) int {
	if right == 0 {
		return 0
	}
	ok, ng := right, -1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(st.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 用于求解`两个字符串s和t`相关性质的后缀数组.
type SuffixArray2 struct {
	SA     *SuffixArray
	offset int
}

// !ord值很大时,需要先离散化.
// !ords[i]>=0.
func NewSuffixArray2(ords1, ords2 []int) *SuffixArray2 {
	newNums := append(ords1, ords2...)
	sa := NewSuffixArray(newNums)
	return &SuffixArray2{SA: sa, offset: len(ords1)}
}

func NewSuffixArray2FromString(s, t string) *SuffixArray2 {
	ords1 := make([]int, len(s))
	for i, c := range s {
		ords1[i] = int(c)
	}
	ords2 := make([]int, len(t))
	for i, c := range t {
		ords2[i] = int(c)
	}
	return NewSuffixArray2(ords1, ords2)
}

// 求任意两个子串s[a,b)和t[c,d)的最长公共前缀(lcp).
func (suf *SuffixArray2) Lcp(a, b int, c, d int) int {
	return suf.SA.Lcp(a, b, c+suf.offset, d+suf.offset)
}

// 比较任意两个子串s[a,b)和t[c,d)的字典序.
//
//	s[a,b) < t[c,d) 返回-1.
//	s[a,b) = t[c,d) 返回0.
//	s[a,b) > t[c,d) 返回1.
func (suf *SuffixArray2) CompareSubstr(a, b int, c, d int) int {
	return suf.SA.CompareSubstr(a, b, c+suf.offset, d+suf.offset)
}

// !注意内部会修改ords.
//
//	 sa : 排第几的后缀是谁.
//	 rank : 每个后缀排第几.
//	 lcp : 排名相邻的两个后缀的最长公共前缀.
//		lcp[0] = 0
//		lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
func UseSA(ords []int) (sa, rank, lcp []int) {
	n := len(ords)
	sa = GetSA(ords)

	rank = make([]int, n)
	for i := range rank {
		rank[sa[i]] = i
	}

	// !高度数组 lcp 也就是排名相邻的两个后缀的最长公共前缀。
	// lcp[0] = 0
	// lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	lcp = make([]int, n)
	h := 0
	for i, rk := range rank {
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := sa[rk-1]; i+h < n && j+h < n && ords[i+h] == ords[j+h]; h++ {
			}
		}
		lcp[rk] = h
	}

	return
}

// 注意内部会修改ords.
func GetSA(ords []int) (sa []int) {
	if len(ords) == 0 {
		return []int{}
	}

	mn := mins(ords)
	for i, x := range ords {
		ords[i] = x - mn + 1
	}
	ords = append(ords, 0)
	n := len(ords)
	m := maxs(ords) + 1
	isS := make([]bool, n)
	isLms := make([]bool, n)
	lms := make([]int, 0, n)
	for i := 0; i < n; i++ {
		isS[i] = true
	}
	for i := n - 2; i > -1; i-- {
		if ords[i] == ords[i+1] {
			isS[i] = isS[i+1]
		} else {
			isS[i] = ords[i] < ords[i+1]
		}
	}
	for i := 1; i < n; i++ {
		isLms[i] = !isS[i-1] && isS[i]
	}
	for i := 0; i < n; i++ {
		if isLms[i] {
			lms = append(lms, i)
		}
	}
	bin := make([]int, m)
	for _, x := range ords {
		bin[x]++
	}

	induce := func() []int {
		sa := make([]int, n)
		for i := 0; i < n; i++ {
			sa[i] = -1
		}

		saIdx := make([]int, m)
		copy(saIdx, bin)
		for i := 0; i < m-1; i++ {
			saIdx[i+1] += saIdx[i]
		}
		for j := len(lms) - 1; j > -1; j-- {
			i := lms[j]
			x := ords[i]
			saIdx[x]--
			sa[saIdx[x]] = i
		}

		copy(saIdx, bin)
		s := 0
		for i := 0; i < m; i++ {
			s, saIdx[i] = s+saIdx[i], s
		}
		for j := 0; j < n; j++ {
			i := sa[j] - 1
			if i < 0 || isS[i] {
				continue
			}
			x := ords[i]
			sa[saIdx[x]] = i
			saIdx[x]++
		}

		copy(saIdx, bin)
		for i := 0; i < m-1; i++ {
			saIdx[i+1] += saIdx[i]
		}
		for j := n - 1; j > -1; j-- {
			i := sa[j] - 1
			if i < 0 || !isS[i] {
				continue
			}
			x := ords[i]
			saIdx[x]--
			sa[saIdx[x]] = i
		}

		return sa
	}

	sa = induce()

	lmsIdx := make([]int, 0, len(sa))
	for _, i := range sa {
		if isLms[i] {
			lmsIdx = append(lmsIdx, i)
		}
	}
	l := len(lmsIdx)
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = -1
	}
	ord := 0
	order[n-1] = ord
	for i := 0; i < l-1; i++ {
		j, k := lmsIdx[i], lmsIdx[i+1]
		for d := 0; d < n; d++ {
			jIsLms, kIsLms := isLms[j+d], isLms[k+d]
			if ords[j+d] != ords[k+d] || jIsLms != kIsLms {
				ord++
				break
			}
			if d > 0 && (jIsLms || kIsLms) {
				break
			}
		}
		order[k] = ord
	}
	b := make([]int, 0, l)
	for _, i := range order {
		if i >= 0 {
			b = append(b, i)
		}
	}
	var lmsOrder []int
	if ord == l-1 {
		lmsOrder = make([]int, l)
		for i, ord := range b {
			lmsOrder[ord] = i
		}
	} else {
		lmsOrder = GetSA(b)
	}
	buf := make([]int, len(lms))
	for i, j := range lmsOrder {
		buf[i] = lms[j]
	}
	lms = buf
	return induce()[1:]
}

type ClampableStackItem = struct {
	value int
	count int32
}

type ClampableStack struct {
	clampMin bool
	total    int
	count    int
	stack    []ClampableStackItem
}

// clampMin：
//
//	为true时，调用AddAndClamp(x)后，容器内所有数最小值被截断(小于x的数变成x)；
//	为false时，调用AddAndClamp(x)后，容器内所有数最大值被截断(大于x的数变成x).
func NewClampableStack(clampMin bool) *ClampableStack {
	return &ClampableStack{clampMin: clampMin}
}

func (h *ClampableStack) AddAndClamp(x int) {
	newCount := 1
	if h.clampMin {
		for len(h.stack) > 0 {
			top := h.stack[len(h.stack)-1]
			if top.value > x {
				break
			}
			h.stack = h.stack[:len(h.stack)-1]
			v, c := top.value, int(top.count)
			h.total -= v * c
			newCount += c
		}
	} else {
		for len(h.stack) > 0 {
			top := h.stack[len(h.stack)-1]
			if top.value < x {
				break
			}
			h.stack = h.stack[:len(h.stack)-1]
			v, c := top.value, int(top.count)
			h.total -= v * c
			newCount += c
		}
	}
	h.total += x * newCount
	h.count++
	h.stack = append(h.stack, ClampableStackItem{value: x, count: int32(newCount)})
}

func (h *ClampableStack) Sum() int {
	return h.total
}

func (h *ClampableStack) Len() int {
	return h.count
}

func (h *ClampableStack) Clear() {
	h.stack = h.stack[:0]
	h.total = 0
	h.count = 0
}

// 求每个元素作为最值的影响范围(闭区间).
func GetRange(nums []int, isMax, isLeftStrict, isRightStrict bool) (leftMost, rightMost []int) {
	compareLeft := func(stackValue, curValue int) bool {
		if isLeftStrict && isMax {
			return stackValue <= curValue
		} else if isLeftStrict && !isMax {
			return stackValue >= curValue
		} else if !isLeftStrict && isMax {
			return stackValue < curValue
		} else {
			return stackValue > curValue
		}
	}

	compareRight := func(stackValue, curValue int) bool {
		if isRightStrict && isMax {
			return stackValue <= curValue
		} else if isRightStrict && !isMax {
			return stackValue >= curValue
		} else if !isRightStrict && isMax {
			return stackValue < curValue
		} else {
			return stackValue > curValue
		}
	}

	n := len(nums)
	leftMost, rightMost = make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		rightMost[i] = n - 1
	}

	stack := []int{}
	for i := 0; i < n; i++ {
		for len(stack) > 0 && compareRight(nums[stack[len(stack)-1]], nums[i]) {
			rightMost[stack[len(stack)-1]] = i - 1
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	stack = []int{}
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && compareLeft(nums[stack[len(stack)-1]], nums[i]) {
			leftMost[stack[len(stack)-1]] = i + 1
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	return
}

type MonoQueueValue = int
type MonoQueue struct {
	MinQueue       []MonoQueueValue
	_minQueueCount []int32
	_less          func(a, b MonoQueueValue) bool
	_len           int
}

func NewMonoQueue(less func(a, b MonoQueueValue) bool) *MonoQueue {
	return &MonoQueue{
		_less: less,
	}
}

func (q *MonoQueue) Append(value MonoQueueValue) *MonoQueue {
	count := int32(1)
	for len(q.MinQueue) > 0 && q._less(value, q.MinQueue[len(q.MinQueue)-1]) {
		q.MinQueue = q.MinQueue[:len(q.MinQueue)-1]
		count += q._minQueueCount[len(q._minQueueCount)-1]
		q._minQueueCount = q._minQueueCount[:len(q._minQueueCount)-1]
	}
	q.MinQueue = append(q.MinQueue, value)
	q._minQueueCount = append(q._minQueueCount, count)
	q._len++
	return q
}

func (q *MonoQueue) Popleft() {
	q._minQueueCount[0]--
	if q._minQueueCount[0] == 0 {
		q.MinQueue = q.MinQueue[1:]
		q._minQueueCount = q._minQueueCount[1:]
	}
	q._len--
}

func (q *MonoQueue) Head() MonoQueueValue {
	return q.MinQueue[0]
}

func (q *MonoQueue) Min() MonoQueueValue {
	return q.MinQueue[0]
}

func (q *MonoQueue) Len() int {
	return q._len
}

func (q *MonoQueue) String() string {
	sb := []string{}
	for i := 0; i < len(q.MinQueue); i++ {
		sb = append(sb, fmt.Sprintf("%v", pair{q.MinQueue[i], q._minQueueCount[i]}))
	}
	return fmt.Sprintf("MonoQueue{%v}", strings.Join(sb, ", "))
}

type pair struct {
	value MonoQueueValue
	count int32
}

func (p pair) String() string {
	return fmt.Sprintf("(value: %v, count: %v)", p.value, p.count)
}

type UnionFindArray struct {
	data []int
}

func NewUnionFindArray(n int) *UnionFindArray {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArray{data: data}
}

// f: 合并两个分组前的钩子函数.
func (ufa *UnionFindArray) Union(key1, key2 int, f func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1, root2 = root2, root1
	}
	if f != nil {
		f(root1, root2)
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

func (ufa *UnionFindArray) Size(key int) int {
	return -ufa.data[ufa.Find(key)]
}

func NewHeap[H any](less func(a, b H) bool, nums []H) *Heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &Heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap[H any] struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap[H]) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap[H]) Len() int { return len(h.data) }

func (h *Heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap[H]) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root
		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}
		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}
		if minIndex == root {
			return
		}
		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}

type Dictionary[V comparable] struct {
	_idToValue []V
	_valueToId map[V]int32
}

// A dictionary that maps values to unique ids.
func NewDictionary[V comparable]() *Dictionary[V] {
	return &Dictionary[V]{
		_valueToId: map[V]int32{},
	}
}
func (d *Dictionary[V]) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return int(res)
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = int32(id)
	return id
}
func (d *Dictionary[V]) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary[V]) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary[V]) Size() int {
	return len(d._idToValue)
}

func mins(a []int) int {
	mn := a[0]
	for _, x := range a {
		if x < mn {
			mn = x
		}
	}
	return mn
}

func maxs(a []int) int {
	mx := a[0]
	for _, x := range a {
		if x > mx {
			mx = x
		}
	}
	return mx
}

func min(a, b int) int {
	if a <= b {
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

func max(a, b int) int {
	if a >= b {
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
