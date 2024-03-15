package main

import (
	"bufio"
	"fmt"
	"os"
)

type TrieNode struct {
}

type Trie struct {
	nodes []*TrieNode
}

// 505. Prefixes and suffixes
// https://codeforces.com/problemsets/acmsguru/problem/99999/505
//
// 给定一个只包含小写字母的字符串集合words.
// 在线处理q个查询：words中有多少个字符串s满足prefix是s的一个前缀, suffix是s的一个后缀.
// n,q,|words[i]|<=1e5
//
// trie+dfs+二维数点:
// !对于给出的每个字符串正着插入字典树A，倒着插入字典树B，
// 对于一个前缀来说，在字典树A上得到的dfs序[st,en]就是所有的匹配串，
// 同理，后缀在字典树B上dfs序[st,en]表示所有的后缀匹配串，
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	// 查询words中有多少个字符串s满足prefix是s的一个前缀, suffix是s的一个后缀.
	query := func(prefix, suffix string) int32 {
		return 0
	}

	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var prefix, suffix string
		fmt.Fscan(in, &prefix, &suffix)
		fmt.Fprintln(out, query(prefix, suffix))
	}
}
