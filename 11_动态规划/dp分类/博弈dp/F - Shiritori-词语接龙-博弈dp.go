// shiritori 词语接龙/成语接龙
// 不能用用过的词 且开头必须和上一个词的结尾相同(如果有的话)
// !问先手是否必胜

package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(words []string) bool {
	n := len(words)
	N := (1 << n) * 30
	memo := make([]int, N)
	for i := 0; i < N; i++ {
		memo[i] = -1
	}

	var dfs func(visited int, pre byte) int
	dfs = func(visited int, pre byte) int {
		hash := visited*27 + int(pre)
		if memo[hash] != -1 {
			return memo[hash]
		}

		for i := 0; i < n; i++ {
			if visited&(1<<i) == 0 && (pre == 26 || words[i][0]-'a' == pre) {
				if dfs(visited|(1<<i), words[i][len(words[i])-1]-'a') == 0 {
					memo[hash] = 1
					return 1
				}
			}
		}
		memo[hash] = 0
		return 0
	}

	return dfs(0, 26) == 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	res := solve(words)
	if res {
		fmt.Fprintln(out, "First")
	} else {
		fmt.Fprintln(out, "Second")
	}
}
