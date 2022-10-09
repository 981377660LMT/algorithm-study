package main

import (
	"bufio"
	"fmt"
	"index/suffixarray"
	"math"
	"math/bits"
	"os"
	"reflect"
	"sort"
	"unsafe"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var s, t []byte
	fmt.Fscan(in, &n, &s, &t)

	s2 := append(s, s[:n-1]...)
	t2 := append(t, t[:n-1]...)
	_, rank1, _ := suffixArray(s2)
	_, rank2, _ := suffixArray(t2)
	rank1 = rank1[:n]
	rank2 = rank2[:n]
	starts1 := make([]int, n)
	for i := range starts1 {
		starts1[i] = i
	}
	starts2 := make([]int, n)
	for i := range starts2 {
		starts2[i] = i
	}

	sort.Slice(starts1, func(i, j int) bool { return rank1[starts1[i]] < rank1[starts1[j]] })
	sort.Slice(starts2, func(i, j int) bool { return rank2[starts2[i]] < rank2[starts2[j]] })
	// # !辞書順で f(S,i) が f(T,j) より小さいか同じであるものの個数を求めてください。

	newS := append(s2, t2...) // !s+s[:-1]+t+t[:-1]

	offset := 2*n - 1
	i, j := 0, 0
	res := 0
	_, newRank, newHeight := suffixArray(newS)
	compareFunc := useCompareSub(newRank, newHeight)
	for i < n && j < n {
		if compareFunc(starts1[i], starts1[i]+n, offset+starts2[j], offset+starts2[j]+n) <= 0 {
			res += n - j
			i++
		} else {
			j++
		}
	}

	fmt.Fprintln(out, res)
}

// 比较两个子串，返回 strings.Compare(s[l1:r1], s[l2:r2])，注意这里是左闭右开区间
func useCompareSub(rank, height []int) func(l1, r1, l2, r2 int) int {
	n := len(rank)

	max := int(math.Ceil(math.Log2(float64(n)))) + 1
	st := make([][]int, n)
	for i := range st {
		st[i] = make([]int, max)
	}

	for i, v := range height {
		st[i][0] = v
	}
	for j := 1; 1<<j <= n; j++ {
		for i := 0; i+1<<j <= n; i++ {
			st[i][j] = min(st[i][j-1], st[i+1<<(j-1)][j-1])
		}
	}

	_q := func(l, r int) int { k := bits.Len(uint(r-l)) - 1; return min(st[l][k], st[r-1<<k][k]) }
	lcp := func(i, j int) int {
		if i == j {
			return n - i
		}
		// 将 s[i:] 和 s[j:] 通过 rank 数组映射为 height 的下标
		ri, rj := rank[i], rank[j]
		if ri > rj {
			ri, rj = rj, ri
		}
		return _q(ri+1, rj+1)
	}

	compareFunc := func(l1, r1, l2, r2 int) int {
		len1, len2 := r1-l1, r2-l2
		l := lcp(l1, l2)
		if len1 == len2 && l >= len1 {
			return 0
		}
		if l >= len1 || l >= len2 { // 一个是另一个的前缀
			if len1 < len2 {
				return -1
			}
			return 1
		}
		if rank[l1] < rank[l2] { // 或者 s[l1+l] < s[l2+l]
			return -1
		}
		return 1
	}

	return compareFunc
}

// https://github.dev/EndlessCheng/codeforces-go/copypasta/strings.go
func suffixArray(s []byte) ([]int32, []int, []int) {
	n := len(s)

	sa := *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New(s)).Elem().FieldByName("sa").Field(0).UnsafeAddr()))

	rank := make([]int, n)
	for i := range rank {
		rank[sa[i]] = i
	}

	height := make([]int, n)
	h := 0
	for i, rk := range rank {
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := int(sa[rk-1]); i+h < n && j+h < n && s[i+h] == s[j+h]; h++ {
			}
		}
		height[rk] = h
	}

	return sa, rank, height
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
