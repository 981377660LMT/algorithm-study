// https://yukicoder.me/problems/no/886
// 给定ROW*COL的矩阵(ROW,COL<=3e6)
// 求连接两点且不经过其他任何点的连接数(类似安卓屏幕解锁的轨迹)

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

	var ROW, COL int
	fmt.Fscan(in, &ROW, &COL)

	res := ROW*(COL-1) + COL*(ROW-1) // 原来相邻的线段数
	n := max(ROW, COL)
	A, B := make([]int, n+2), make([]int, n+2)
	for i := 0; i < ROW; i++ {
		A[i] = ROW - i - 1
	}
	for i := 0; i < COL; i++ {
		B[i] = COL - i - 1
	}

	C := GcdConvolution(A, B)
	fmt.Fprintln(out, (res+C[0]*2)%MOD) // pair of gcd(row,col)=1

}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

const MOD int = 1e9 + 7

// c[k-1] = ∑a[i]*b[j] mod MOD, gcd(i,j)=k (k=1,2,...)
func GcdConvolution(nums1, nums2 []int) []int {
	n := len(nums1)
	pf := make([]int, n+1)
	copy1, copy2 := make([]int, n+1), make([]int, n+1)
	for i := 0; i < n; i++ {
		copy1[i+1] = nums1[i]
		copy2[i+1] = nums2[i]
	}

	for i := 2; i < n+1; i++ {
		if pf[i] == 0 {
			for j := n / i; j > 0; j-- {
				pf[j*i] = 1
				copy1[j] = (copy1[j] + copy1[j*i]) % MOD
				copy2[j] = (copy2[j] + copy2[j*i]) % MOD
			}
			pf[i] = 0
		}
	}

	res := make([]int, n+1)
	for i := 0; i < n+1; i++ {
		res[i] = copy1[i] * copy2[i] % MOD
	}

	for i := 2; i < n+1; i++ {
		if pf[i] == 0 {
			for j := 1; j < n/i+1; j++ {
				res[j] = (res[j] - res[j*i] + MOD) % MOD
			}
		}
	}

	return res[1:]
}
