package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://yukicoder.me/problems/13
	// !将石头堆拆分成多个堆,当前堆的grundy数等于子状态的grundy数的异或和
	// n個の石が積まれた山が1つある。
	// A君とB君が交互に石を分けるゲームを行う。
	// !分けるときに石を2つの山か3つの山に分ける。
	// eg:
	// 2x -> x,x
	// 2x+1 -> x,x+1
	// 3x -> x,x,x
	// 3x+1 -> x,x,x+1
	// 3x+2 -> x,x+1,x+1

	// !ゲームは石を最後に分けられなくなったほうが負けである。
	// よって、この最初の石が5個のゲームの場合には、
	// ケース１のように先手のA君がまず石を3つに分ければA君が必ず勝てる。
	// A君が先手でA君もB君も勝つために最善を尽くすとき、
	// 最初のNによってA君が勝つかB君が勝つかを判定せよ。
	// 1<=n<=100

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	memo := make([]int, 110)
	for i := range memo {
		memo[i] = -1
	}

	var grundy func(state int) int
	grundy = func(state int) int {
		if memo[state] != -1 {
			return memo[state]
		}
		if state == 0 || state == 1 {
			memo[state] = 0
			return 0
		}

		nextStates := make(map[int]struct{})
		if state%2 == 0 {
			nextStates[grundy(state/2)^grundy(state/2)] = struct{}{}
		}
		if state%2 == 1 {
			nextStates[grundy(state/2)^grundy(state/2+1)] = struct{}{}
		}
		if state%3 == 0 {
			nextStates[grundy(state/3)^grundy(state/3)^grundy(state/3)] = struct{}{}
		}
		if state%3 == 1 {
			nextStates[grundy(state/3)^grundy(state/3)^grundy(state/3+1)] = struct{}{}
		}
		if state%3 == 2 {
			nextStates[grundy(state/3)^grundy(state/3+1)^grundy(state/3+1)] = struct{}{}
		}

		mex := 0
		for {
			if _, ok := nextStates[mex]; !ok {
				break
			}
			mex++
		}
		memo[state] = mex
		return mex
	}

	if grundy(n) == 0 {
		fmt.Fprintln(out, "B")
	} else {
		fmt.Fprintln(out, "A")
	}
}
