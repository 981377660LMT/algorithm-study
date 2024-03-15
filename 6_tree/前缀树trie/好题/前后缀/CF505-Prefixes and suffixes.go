package main

import (
	"bufio"
	"fmt"
	"os"
)

// 505. Prefixes and suffixes
// https://codeforces.com/problemsets/acmsguru/problem/99999/505
//
// 解法1：按照长度hash.
// 解法2：两棵trie树，一棵存前缀，一棵存后缀.符合条件的字符串在子树内.转换成二维数点问题.
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
