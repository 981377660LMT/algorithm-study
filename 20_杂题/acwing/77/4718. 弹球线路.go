/*
给定一个 n×m 的方格矩阵。

将一个弹球放置到其中的任意一个方格之中，
并使其沿四个对角方向之一（左上、左下、右上、右下）进行无限移动。
当弹球移动至矩阵边缘的方格时，它可以发生反弹并继续移动。

更具体地说，当它到达矩阵边（但并非角）上的方格时，
它可以将移动方向改变 90 度并继续移动，当它到达矩阵角上的方格时，
它可以将移动方向改变 180 度并继续移动。

不难发现，弹球会沿着一个固定的线路进行无限循环移动。
我们认为两个线路不同，当且仅当两个线路包含的方格不完全相同（而不考虑弹球的具体移动方向）。
弹球在不同的初始位置沿不同的初始方向进行移动，其移动线路可能不同，也可能相同。
!现在，请你计算弹球一共可能有多少种不同的移动线路。
n,m<=1e6

!推公式(光线反射/镜面反射)/并查集
*/

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	const INF int = int(1e18)
	const MOD int = 998244353

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var ROW, COL int
	fmt.Fscan(in, &ROW, &COL)

	res := gcd(ROW-1, COL-1) + 1 // !公式
	fmt.Fprintln(out, res)
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}
