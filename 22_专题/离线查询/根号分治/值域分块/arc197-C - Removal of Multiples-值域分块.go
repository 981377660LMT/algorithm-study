// C - Removal of Multiples
// https://atcoder.jp/contests/arc197/tasks/arc197_c
//
// 给定一个由所有正整数组成的集合 **S**，需要处理 **Q** 个查询。对于每个查询，执行以下两步操作：
//
// 1. 从集合 **S** 中删除所有是给定整数 **Aᵢ** 的倍数的元素。
// 2. 输出集合 **S** 中按升序排列的第 **Bᵢ** 小的元素。
//
// **保证：**
// - 每次查询后，集合 **S** 中的元素数量至少为 **Bᵢ**。
//
// q<=1e5,Ai<=1e9,Bi<=1e5.
//
// !分块维护，每个块维护值域 [K*i, K*i+K) 的元素，

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const N int = 1 << 22 // 答案上界
	const K int = 1 << 11 // 分块大小
	block := [K]int{}     // 值域分块，维护这这块还有多少个元素
	exist := [N]bool{}

	for v := 1; v < N; v++ {
		exist[v] = true
		block[v/K]++
	}

	_removeAt := func(i int) {
		exist[i] = false
		block[i/K]--
	}

	remove := func(mul int) {
		if mul >= N || !exist[mul] {
			return
		}
		for i := mul; i < N; i += mul {
			if exist[i] {
				_removeAt(i)
			}
		}
	}

	query := func(pos int) int {
		for bid := 0; bid < K; bid++ {
			if block[bid] <= pos {
				pos -= block[bid]
				continue
			}

			start, end := bid*K, bid*K+K
			for v := start; v < end; v++ {
				if exist[v] {
					if pos == 0 {
						return v
					}
					pos--
				}
			}
		}
		panic("should not reach here")
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		remove(a)
		res := query(b - 1)
		fmt.Fprintln(out, res)
	}
}
