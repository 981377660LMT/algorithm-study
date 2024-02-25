package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	son [26]*Node
	s   string
}

// Spy Syndrome 2
// https://www.luogu.com.cn/problem/CF633C
// 给定长为(n <= 10000)的主串，给(m <= 100000)个不超过 1000 的子串，
// 总长不超过 1000000。问主串是否能由子串拼出，每个子串可以多次使用
//
// trie 树 + 记忆化搜索(trie树上dp)，多模式串匹配
// 首先把所有单词的小写形式插入一颗 trie 树，末尾节点记录原始单词。
// dp[i] 表示前缀 s[0..i] 能否拆分成若干个单词，从 s 的末尾往前记忆化搜索即可得到答案
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)
	var m int
	fmt.Fscan(in, &m)

	trie := &Node{}
	for i := 0; i < m; i++ {
		var w string
		fmt.Fscan(in, &w)

		// 将单词的小写形式插入 trie 树
		root := trie
		for _, c := range w {
			if c < 'a' {
				c += 32
			}
			c -= 'a'
			if root.son[c] == nil {
				root.son[c] = &Node{}
			}
			root = root.son[c]
		}
		root.s = w // 末尾节点记录原单词
	}

	memo := make([]int8, n)
	var dfs func(int32) int8
	dfs = func(p int32) (res int8) {
		if p < 0 {
			return 1
		}
		m := &memo[p]
		if *m != 0 {
			return *m
		}
		defer func() { *m = res }()

		cur := trie
		for i := p; i >= 0; i-- {
			cur = cur.son[s[i]-'a']
			if cur == nil {
				return -1
			}
			if cur.s != "" && dfs(i-1) == 1 {
				fmt.Fprint(out, cur.s, " ")
				return 1
			}
		}
		return -1
	}

	dfs(int32(n - 1))
}
