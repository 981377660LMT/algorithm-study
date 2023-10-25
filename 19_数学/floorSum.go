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

	var n, k int
	fmt.Fscan(in, &n, &k)
	fmt.Fprintln(out, ModSum(n, k))
}

// 类欧几里得算法
// 商求和(divSum)
// ∑⌊(ai+b)/m⌋, i in [0,n-1]
// https://oi-wiki.org/math/euclidean/
func FloorSum(n, m, a, b int) (res int) {
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

// 余数求和(ModSum/remainderSum)
// ∑k%i (i in [1,n]), 即 sum(k%i for i in [1,n])
// = ∑k-(k/i)*i
// = n*k-∑(k/i)*i
// 对于 [l,r] 范围内的 i，k/i 不变，此时 ∑(k/i)*i = (k/i)*∑i = (k/i)*(l+r)*(r-l+1)/2
func ModSum(n int, k int) int {
	sum := n * k
	for l, r := 1, 0; l <= n; l = r + 1 {
		h := k / l
		if h > 0 {
			r = min(k/h, n)
		} else {
			r = n
		}
		w := r - l + 1
		s := (l + r) * w / 2
		sum -= h * s
	}
	return sum
}

// 二维整除分块.
// ∑{i=1..min(n,m)} floor(n/i)*floor(m/i)
// https://www.luogu.com.cn/blog/command-block/zheng-chu-fen-kuai-ru-men-xiao-ji
func FloorSum2D(n, m int) (sum int) {
	for l, r := 1, 0; l <= min(n, m); l = r + 1 {
		hn, hm := n/l, m/l
		r = min(n/hn, m/hm)
		w := r - l + 1
		sum += hn * hm * w
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
