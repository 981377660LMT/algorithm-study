// F - Oddly Similar
// https://atcoder.jp/contests/abc348/tasks/abc348_f
// 给定n个长度为m的数组.
// 求有多少对数组满足，对应位置的元素相等的个数为奇数.
// 1 <= n, m <= 2000
// 1 <= a[i][j] <= 999
// https://atcoder.jp/contests/abc348/editorial/9752

// 高速な言語を用いた上でif文の回避や配列へのアクセスの順番を変えるなどの工夫をすることで,
// O(N^2M) 解法が十分高速に動作します。
// 1.消除if分支
// 2.改变数组访问顺序

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

	var n, m int16
	fmt.Fscan(in, &n, &m)

	M := [2005][2005]int16{}
	for i := int16(0); i < n; i++ {
		for j := int16(0); j < m; j++ {
			fmt.Fscan(in, &M[j][i])
		}
	}

	B := [2005][2005]bool{}

	for i := int16(0); i < m; i++ {
		for j := int16(0); j < n; j++ {
			for k := j + 1; k < n; k++ {
				B[j][k] = B[j][k] != (M[i][j] == M[i][k])
			}
		}
	}

	res := int32(0)
	for i := int16(0); i < n; i++ {
		for j := i + 1; j < n; j++ {
			if B[i][j] {
				res++
			}
		}
	}
	fmt.Fprintln(out, res)
}
