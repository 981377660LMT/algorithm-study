// 子序列自动机
// 如果值域很大可以用哈希表/数组记录 pos 然后二分查找 https://www.luogu.com.cn/problem/P5826
// LC727 https://leetcode.cn/problems/minimum-window-subsequence/
// LC792 https://leetcode.cn/problems/number-of-matching-subsequences/
// LC2014 https://leetcode.cn/problems/longest-subsequence-repeated-k-times/
// LC466 https://leetcode.cn/problems/count-the-repetitions/
// http://codeforces.com/problemset/problem/91/A
// - https://www.luogu.com.cn/problem/P9572?contestId=124047
// - 【子串】 LC686 https://leetcode.cn/problems/repeated-string-match/
// https://codeforces.com/contest/1845/problem/C
// - 相关 LC2350 https://leetcode.cn/problems/shortest-impossible-sequence-of-rolls/
subsequenceAutomaton := func(s string) {
	// build nxt
	// nxt[i][j] 表示在 i 右侧的字符 j 的最近位置
	pos := [26]int{}
	for i := range pos {
		pos[i] = len(s)
	}
	nxt := make([][26]int, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		nxt[i] = pos
		pos[s[i]-'a'] = i
		// 注意这行在下面。我个人更喜欢这种写法，不喜欢写 nxt... +1，让一个指针指向当前位置右侧太奇怪了，不如这种写法清晰
	}

	// 返回是 s 的子序列的最长的 t 的前缀的长度
	match := func(t string) int {
		if t == "" || s == "" {
			return 0
		}
		i, j := 0, 0
		if t[0] == s[0] {
			j = 1 // t[0] 匹配 ok
		}
		for ; j < len(t); j++ {
			i = nxt[i][t[j]-'a']
			if i == len(s) {
				break
			}
		}
		return j
	}
	_ = match
}