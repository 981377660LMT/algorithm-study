package main

import (
	"bufio"
	"fmt"
	"os"
)

// from functools import lru_cache
// import sys

// sys.setrecursionlimit(int(1e9))
// input = lambda: sys.stdin.readline().rstrip("\r\n")
// MOD = 998244353
// INF = int(4e18)
// # N 個の文字列 S
// # 1
// # ​
// #  ,S
// # 2
// # ​
// #  ,…,S
// # N
// # ​
// #   が与えられます。 S
// # i
// # ​
// #   (1≤i≤N) は英小文字からなる長さ 10 以下の空でない文字列で、互いに異なります。

// # 先手太郎君と後手次郎君がしりとりをします。 このしりとりでは、先手太郎君と後手次郎君の手番が交互に訪れます。 はじめの手番は先手太郎君の手番です。 それぞれのプレイヤーは自分の手番において整数 i (1≤i≤N) を 1 つ選びます。 このとき、i は次の 2 つの条件を満たしていなければなりません。

// # i は、しりとりが開始してからこれまでの 2 人の手番で選ばれたどの整数とも異なる
// # この手番がしりとりの最初の手番であるか、直前に選ばれた整数を j として、S
// # j
// # ​
// #   の最後の文字と S
// # i
// # ​
// #   の最初の文字が等しい
// # 条件を満たす i を選べなくなったプレイヤーの負けで、負けなかったプレイヤーの勝ちです。

// # 2 人が最適に行動したときに勝つのはどちらかを判定してください。
// if __name__ == "__main__":
//     n = int(input())
//     words = [input() for _ in range(n)]

//     @lru_cache(None)
//     def dfs(visited: int, pre: str) -> bool:
//         for i in range(n):
//             if visited & (1 << i) or words[i][0] != pre:
//                 continue
//             if not dfs(visited | (1 << i), words[i][-1]):
//                 return True
//         return False

//     for start in range(n):
//         if dfs(1 << start, words[start][0]):
//             print("First")
//             exit(0)
//     print("Second")

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

	N := (1 << n) * 32
	memo := make([]int, N)
	for i := 0; i < N; i++ {
		memo[i] = -1
	}

	var dfs func(visited int, pre byte) int
	dfs = func(visited int, pre byte) int {
		hash := visited*30 + int(pre)
		if memo[hash] != -1 {
			return memo[hash]
		}

		for i := 0; i < n; i++ {
			if visited&(1<<i) == 0 && (pre == 27 || words[i][0]-'a' == pre) {
				if dfs(visited|(1<<i), words[i][len(words[i])-1]-'a') == 0 {
					memo[hash] = 1
					return 1
				}
			}
		}
		memo[hash] = 0
		return 0
	}

	ok := false
	for start := 0; start < n; start++ {
		if dfs(0, 27) == 1 {
			fmt.Fprintln(out, "First")
			ok = true
			break
		}
	}

	if !ok {
		fmt.Fprintln(out, "Second")
	}
}
