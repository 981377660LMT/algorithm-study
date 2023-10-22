package main

func main() {

}

// 类欧几里得算法
// ∑⌊(ai+b)/m⌋, i in [0,n-1]
// https://oi-wiki.org/math/euclidean/
// todo https://www.luogu.com.cn/blog/AlanWalkerWilson/Akin-Euclidean-algorithm-Basis
//
//	https://www.luogu.com.cn/blog/Shuchong/qian-tan-lei-ou-ji-li-dei-suan-fa
//	万能欧几里得算法 https://www.luogu.com.cn/blog/ILikeDuck/mo-neng-ou-ji-li-dei-suan-fa
//
// 模板题 https://atcoder.jp/contests/practice2/tasks/practice2_c
//
//	https://www.luogu.com.cn/problem/P5170
//	https://loj.ac/p/138
//
// todo https://codeforces.com/problemset/problem/1182/F
//
//	https://codeforces.com/problemset/problem/1098/E
func floorSum(n, m, a, b int) (res int) {
	if a < 0 {
		a2 := a%m + m
		res -= n * (n - 1) / 2 * ((a2 - a) / m)
		a = a2
	}
	if b < 0 {
		b2 := b%m + m
		res -= n * ((b2 - b) / m)
		b = b2
	}
	for {
		if a >= m {
			res += n * (n - 1) / 2 * (a / m)
			a %= m
		}
		if b >= m {
			res += n * (b / m)
			b %= m
		}
		yMax := a*n + b
		if yMax < m {
			break
		}
		n = yMax / m
		b = yMax % m
		m, a = a, m
	}
	return
}
