// Lyndon 分解
// https://qiita.com/nakashi18/items/66882bd6e0127174267a
// https://oi-wiki.org/string/lyndon/
// https://hitonanode.github.io/cplib-cpp/string/lyndon.hpp

// Lydon 字符串（Lyndon word）:
//   - 如果一个字符串是它所有非空后缀中字典序最小的字符串, 则称这个字符串为 Lyndon 字符串
//     后缀数组中 sa[0] = 0, rank[0] = 0
//    例如: aabab 是 Lyndon 字符串
//!   - Lydon 字符串是字典序最小的循环移位字符串
//!   - Lydon 字符串中不存在 border (border:既是前缀又是后缀的字符串) (结合kmp考虑)
// Lydon 分解（Lyndon factorization）:
//   - 将一个字符串S分解成若干个 Lyndon 字符串 w1w2...wk, 且各个 wi 是字典序不增的
//   - 一个字符串的 Lyndon 分解是唯一的
//   - w1是S的前缀中最长的Lyndon字符串

// Usage:
// LyndonFactorizationString(s) : Lyndon分解, 返回每个lyndon字符串的 [start,end).
// LongestLyndonPrefixes(s, lcp) : 返回s的每个后缀的最长Lyndon前缀长度
// RunEnumerate(s) : 返回s的所有Run 的 [t,l,r]
//   - Run: s[l:r]的最小周期为t,且r-l>=2t,且l,r是满足条件的最大值
//   - abcbcba 的 run 为 [2,1,6]
//   - run可以理解为包含至少出现两次相同子串的字符串
// EnumerateLyndonWords(k, n) : 按照字典序返回长度不超过n的所有Lyndon字符串,字符串中字符的范围是[0,k)

package main

import (
	"fmt"
	"index/suffixarray"
	"math"
	"math/bits"
	"reflect"
	"unsafe"
)

func main() {
	fmt.Println(LyndonFactorizationString("abaababaababaaabbaaaabbaa"))

	_, rank, height := suffixArray("teletelepathy`")
	lcp := LCP(rank, height)
	fmt.Println(LongestLyndonPrefixes("teletelepathy", lcp))

	fmt.Println(RunEnumerate("mississippi"))

	fmt.Println(EnumerateLyndonWords(2, 4))
}

// func main() {
// 	// https://judge.yosupo.jp/submission/127899
// 	in := bufio.NewReader(os.Stdin)
// 	out := bufio.NewWriter(os.Stdout)
// 	defer out.Flush()

// 	var s string
// 	fmt.Fscan(in, &s)
// 	runs := RunEnumerate(s)
// 	fmt.Fprintln(out, len(runs))
// 	for _, run := range runs {
// 		fmt.Fprintln(out, run[0], run[1], run[2])
// 	}
// }

// lydon分解, 返回每个lyndon字符串的 [start,end).
//  O(n)
func LyndonFactorizationString(s string) [][2]int {
	n := len(s)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = int(s[i])
	}
	return LyndonFactorization(nums)
}

// lydon分解, 返回每个lyndon字符串的 [start,end).
//  O(n)
func LyndonFactorization(nums []int) [][2]int {
	n := len(nums)
	res := [][2]int{}
	for l := 0; l < n; {
		i, j := l, l+1
		for j < n && nums[i] <= nums[j] {
			if nums[i] == nums[j] {
				i++
			} else {
				i = l
			}
			j++
		}
		n := (j - l) / (j - i)
		for t := 0; t < n; t++ {
			res = append(res, [2]int{l, l + j - i})
			l += j - i
		}
	}

	return res
}

// 对每个后缀 S[i:] 求出最长lyndon前缀的长度
//   - `teletelepathy` -> [1,4,1,2,1,4,1,2,1,4,1,2,1]
//   - O(n*log(n))
func LongestLyndonPrefixes(s string, lcp func(i, j int) int) []int {
	n := len(s)
	st := [][2]int{{n, n}}
	res := make([]int, n)
	for i, j := n-1, n-1; i >= 0; i, j = i-1, i-1 {
		for len(st) > 1 {
			iv, jv := st[len(st)-1][0], st[len(st)-1][1]
			l := lcp(i, iv)
			if !(iv+l < n && s[i+l] < s[iv+l]) {
				break
			}
			j = jv
			st = st[:len(st)-1]
		}
		st = append(st, [2]int{i, j})
		res[i] = j - i + 1
	}
	return res
}

// 求各个run 的(最小周期,起点,终点)三元组 (c,l,r),按照c升序排序返回.
//  O(n*l)
func RunEnumerate(s string) [][3]int {
	n := len(s)
	if n == 0 {
		return nil
	}
	_, rank1, height1 := suffixArray(s)
	lcp1 := LCP(rank1, height1)
	reverse := func(s string) string {
		ords := []rune(s)
		for i, j := 0, len(ords)-1; i < j; i, j = i+1, j-1 {
			ords[i], ords[j] = ords[j], ords[i]
		}
		return string(ords)
	}

	revs := reverse(s)
	_, rank2, height2 := suffixArray(revs)
	lcp2 := LCP(rank2, height2)

	runs := [][3]int{}
	vis := make(map[[2]int]struct{})
	lst := -1
	for p := 1; p <= n/2; p++ {
		for i := 0; i <= n-p; i += p {
			left := i - lcp2(n-i-p, n-i)
			r := i - p + lcp1(i, i+p)
			if left > r || left == lst {
				continue
			}
			if _, ok := vis[[2]int{left, r + 2*p}]; !ok {
				vis[[2]int{left, r + 2*p}] = struct{}{}
				runs = append(runs, [3]int{p, left, r + 2*p})
			}
			lst = left
		}
	}

	return runs
}

// 求出长度不超过n的所有lyndon字符串,字符由0~k-1组成,按照字典序升序排序
//  (给定字符集,求出所有排列中哪些是Lyndon字符串)
//  k=2, n=4 => [[0,],[0,0,0,1,],[0,0,1,],[0,0,1,1,],[0,1,],[0,1,1,],[0,1,1,1,],[1,],]
func EnumerateLyndonWords(k, n int) [][]int {
	res := [][]int{}
	aux := make([]int, n+1)
	var gen func(t, p int)
	gen = func(t, p int) {
		// t: current length
		// p: current min cycle length
		if t == n {
			slice := make([]int, p)
			copy(slice, aux[1:p+1])
			res = append(res, slice)
		} else {
			t++
			aux[t] = aux[t-p]
			gen(t, p)
			for aux[t]++; aux[t] < k; aux[t]++ {
				gen(t, t)
			}
		}
	}

	gen(0, 1)
	return res
}

// https://github.dev/EndlessCheng/codeforces-go/copypasta/strings.go
// 求两个后缀最长公共前缀 O(nlogn)预处理 O(1)查询
//  _, rank, height := suffixArray("banana")
//  lcp := LCP(rank, height)
//  fmt.Println(lcp(2, 4)) // "nana" "na"  2
func LCP(rank, height []int) func(int, int) int {
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
		if i >= n || j >= n {
			return 0
		}
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

	return lcp

}

// https://github.dev/EndlessCheng/codeforces-go/copypasta/strings.go
func suffixArray(s string) ([]int32, []int, []int) {
	n := len(s)

	sa := *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New([]byte(s))).Elem().FieldByName("sa").Field(0).UnsafeAddr()))

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
