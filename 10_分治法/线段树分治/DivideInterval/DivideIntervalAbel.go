package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// abc355_e()
	demo()
}

// 交互题
// abc355-E- Guess the Sum-bfs拆分区间
// 给定任意区间[l,r)，求将区间拆分成若干个形如[2^i*j,2^i*(j+1))的区间.
// !操作满足阿贝尔群.
//
// !注意每次输出后调用out.Flush()刷新输出缓冲区.
func abc355_e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	query := func(start, end int32) int {
		i, j := Format(start, end)
		fmt.Println("?", i, j)
		out.Flush()
		var res int
		fmt.Fscan(in, &res)
		return res
	}

	output := func(res int) {
		fmt.Println("!", res)
		out.Flush()
	}

	var LOG, L, R int
	fmt.Fscan(in, &LOG, &L, &R)

	res := 0
	DivideIntervalAbel(L, R+1, func(start, end int, b bool) {
		if b {
			res += query(int32(start), int32(end))
		} else {
			res -= query(int32(start), int32(end))
		}
	})

	res = (res%100 + 100) % 100
	output(res)
}

func Format(start, end int32) (i, j int32) {
	n := end - start
	i = int32(bits.Len(uint(n - 1)))
	j = start >> i
	return
}

func demo() {
	DivideIntervalAbel(0, 10, func(start, end int, b bool) {
		fmt.Println(start, end, b)
	})
	DivideIntervalAbel(3, 11, func(start, end int, b bool) {
		fmt.Println(start, end, b)
	})
}

// 给定任意区间[l,r)，求将区间拆分成若干个形如[2^i*j,2^i*(j+1))的区间，使得拆分的区间个数最少.
func DivideIntervalAbel(l, r int, f func(start, end int, b bool)) {
	if l >= r {
		return
	}
	size := 1
	for l != r {
		if l&1 == 1 {
			if l&2 == 0 && l+1 != r && l+2 != r {
				l--
				start := size * l
				f(start, start+size, false)
			} else {
				start := size * l
				f(start, start+size, true)
				l++
			}
		}
		if r&1 == 1 {
			if r&2 == 0 || l+1 == r {
				r--
				start := size * r
				f(start, start+size, true)
			} else {
				start := size * r
				f(start, start+size, false)
				r++
			}
		}
		l >>= 1
		r >>= 1
		size <<= 1
	}
}
