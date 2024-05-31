package main

import "fmt"

func main() {
	var TABLE = [...]int32{
		0, 1, 1, 1, 2, 3, 6, 11, 23, 47, 106, 235, 551, 1301, 3159, 7741, 19320, 48629,
		123867, 317955, 823065, 2144505, 5623756, 14828074, 39299897, 104636890,
	}

	test := func(n int32) {
		count := int32(0)
		EnumerateUnLabeledRootedTree(n, func(edges [][2]int32) {
			count++
		})
		if count != TABLE[n] {
			panic("assertion failed")
		}
	}

	for i := int32(0); i < int32(20); i++ {
		test(i)
	}
	fmt.Println("pass")
}

// `无标号有根树`同型枚举
// https://oeis.org/A000055
func EnumerateUnLabeledRootedTree(n int32, f func(edges [][2]int32)) {
	if n == 0 {
		return
	}

	m := n / 2
	dat := make([][]uint32, m+1)
	if m >= 1 {
		dat[1] = append(dat[1], 0)
	}

	for r := int32(2); r <= m; r++ {
		var dfs func(m, k int32, now uint32, nowE int32)
		dfs = func(m, k int32, now uint32, nowE int32) {
			if nowE == r-1 {
				dat[r] = append(dat[r], now)
				return
			}
			if nowE+m >= r {
				m = r - 1 - nowE
				k = 0
			}
			if m == 0 {
				return
			}
			for i := k; i < int32(len(dat[m])); i++ {
				x := dat[m][i]
				x = (x << 1) | 1
				dfs(m, i, now|x<<(2*nowE), nowE+m)
			}
			dfs(m-1, 0, now, nowE)
		}
		dfs(r-1, 0, 0, 0)
	}

	decode := func(x uint64) [][2]int32 {
		var edge [][2]int32
		path := []int32{0}
		p := int32(0)
		for i := int32(0); i < 2*n-2; i++ {
			if x>>i&1 == 1 {
				edge = append(edge, [2]int32{path[len(path)-1], p + 1})
				path = append(path, p+1)
				p++
			} else {
				path = path[:len(path)-1]
			}
		}
		return edge
	}

	var dfs func(m, k int32, now uint64, nowE int32)
	dfs = func(m, k int32, now uint64, nowE int32) {
		if nowE == n-1 {
			f(decode(now))
			return
		}
		if nowE+m >= n {
			m = n - 1 - nowE
			k = 0
		}
		if m == 0 {
			return
		}
		for i := k; i < int32(len(dat[m])); i++ {
			x := uint64(dat[m][i])
			x = (x << 1) | 1
			dfs(m, i, now|x<<(2*nowE), nowE+m)
		}
		dfs(m-1, 0, now, nowE)
	}
	dfs((n-1)/2, 0, 0, 0)

	if 2*m == n {
		for i := int32(0); i < int32(len(dat[m])); i++ {
			for j := int32(0); j <= i; j++ {
				x, y := uint64(dat[m][i]), uint64(dat[m][j])
				y = (y << 1) | 1
				f(decode(x | (y << (2 * (m - 1)))))
			}
		}
	}
}
