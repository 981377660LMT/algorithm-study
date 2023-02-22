// https://yukicoder.me/problems/1016

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// !分割竹子
	// 有一根长L的竹子,轮流分成三段,要保证L1,L2,L3都不相等长且不为0
	// 且 max(L1,L2,L3)-min(L1,L2,L3)<=D
	// 不能继续操作的人输
	// 问先手是否必胜
	// 1<=D<=L<=500

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var L, D int
	fmt.Fscan(in, &L, &D)

	memo := make([]int, L+10)
	for i := range memo {
		memo[i] = -1
	}

	var grundy func(int) int
	grundy = func(len_ int) int {
		if memo[len_] != -1 {
			return memo[len_]
		}

		// 枚举分割长度
		nextStates := make(map[int]struct{})
		for len1 := 1; len1 < len_; len1++ {
			for len2 := 1; len2 < len_; len2++ {
				len3 := len_ - len1 - len2
				if len1 < len2 && len2 < len3 && len3-len1 <= D {
					nextStates[grundy(len1)^grundy(len2)^grundy(len3)] = struct{}{}
				}
			}
		}

		mex := 0
		for {
			if _, ok := nextStates[mex]; !ok {
				break
			}
			mex++
		}
		memo[len_] = mex
		return mex
	}

	// 先手はkadoくん、後手はmatsuくんです。
	if grundy(L) == 0 {
		fmt.Fprintln(out, "matsu") // lose
	} else {
		fmt.Fprintln(out, "kado") // win
	}
}
