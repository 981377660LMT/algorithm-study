// https://www.luogu.com.cn/problem/U248425

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

	var n, m, k int
	var s, t string
	fmt.Fscan(in, &n)
	fmt.Fscan(in, &s)
	fmt.Fscan(in, &m)
	fmt.Fscan(in, &t)
	fmt.Fscan(in, &k)

	fmt.Fprintln(out, MatchWithFail(s, t, k))
}

const MOD1, MOD2 int = 1e8 + 7, 1e9 + 7
const BASE1, BASE2 int = 131, 13131
const OFFSET1, OFFSET2 int = 0, 0

// 允许失配k次的字符串匹配.k<=5.
func MatchWithFail(long, short string, k int) int {
	n1, n2 := len(long), len(short)
	if n1 < n2 {
		return 0
	}
	ords1, ords2 := make([]int, n1), make([]int, n2)
	for i, v := range long {
		ords1[i] = int(v)
	}
	for i, v := range short {
		ords2[i] = int(v)
	}
	h1 := StringHasher(ords1, MOD1, MOD2, BASE1)
	h2 := StringHasher(ords2, MOD1, MOD2, BASE1)

	// 从long的下标start开始,找到long[offset:offset+n2)与short第一个不同的位置.
	// 如果不存在,返回-1.
	indexOfDiff := func(start, offset int) int {
		if start >= offset+n2 {
			return -1
		}
		left, right := start, offset+n2
		for left <= right {
			mid := (left + right) >> 1
			if h1(offset, mid) == h2(0, mid-offset) {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		if left < offset+n2 {
			return left
		}
		return -1
	}

	res := 0
	for start := 0; start < n1-n2+1; start++ {
		cur := start
		ok := false
		for i := 0; i < k+1; i++ {
			nextDiff := indexOfDiff(cur, start)
			if nextDiff == -1 {
				ok = true
				break
			}
			cur = nextDiff + 1
		}
		if ok {
			res++
		}
	}
	return res
}

func StringHasher(ords []int, mod int, base int, offset int) func(left int, right int) int {
	prePow := make([]int, len(ords)+1)
	prePow[0] = 1
	preHash := make([]int, len(ords)+1)
	for i, v := range ords {
		prePow[i+1] = (prePow[i] * base) % mod
		preHash[i+1] = (preHash[i]*base + v - offset) % mod
	}

	sliceHash := func(left, right int) int {
		if left >= right {
			return 0
		}
		return (preHash[right] - preHash[left]*prePow[right-left]%mod + mod) % mod
	}

	return sliceHash
}
