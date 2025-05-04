package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	demo()
	// yosupo()
}

// https://judge.yosupo.jp/problem/runenumerate
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	runs := [][3]int{}
	RunEnumerate([]rune(s), func(period, left, right int) {
		runs = append(runs, [3]int{period, left, right})
	})
	sort.Slice(runs, func(i, j int) bool {
		if runs[i][0] != runs[j][0] {
			return runs[i][0] < runs[j][0]
		}
		if runs[i][1] != runs[j][1] {
			return runs[i][1] < runs[j][1]
		}
		return runs[i][2] < runs[j][2]
	})

	fmt.Fprintln(out, len(runs))
	for _, run := range runs {
		fmt.Fprintln(out, run[0], run[1], run[2])
	}
}

func demo() {
	// 示例用法
	s := "aaaaaaaaaa"
	fmt.Println("Runs in string:", s)
	RunEnumerate([]rune(s), func(period, left, right int) {
		fmt.Printf("Period: %d, Range: [%d, %d], Substring: %s\n",
			period, left, right, s[left:right])
	})
}

// RunEnumerate 查找并返回字符串中的所有runs(最大周期子串)
// Run定义为三元组 (period, start, end)，其中：
// - period: 子串的最小周期
// - start: 子串的左边界(包含)
// - end: 子串的右边界(包含)
//
// 具体性质：
//  1. 极大性：S[start:end]具有周期p且长度至少为2p，但S[start-1:end]和S[start:end+1]不满足该周期性
//     (即周期性不能向左或向右扩展，这保证了run的最大性)
//
// 2. 数量限制：长度为n的字符串中，runs的数量不超过n个
//
//  3. 复杂度特性：所有runs的(end-start)/period之和为O(n)
//     (这是runs数量线性的重要理论基础)
func RunEnumerate[S ~[]E, E comparable](s S, f func(period, left, right int)) {
	n := len(s)

	var res []triple

	var dfs func(int, int)
	dfs = func(l, r int) {
		if r-l < 2 {
			return
		}

		mid := (l + r) / 2

		dfs(l, mid)
		dfs(mid, r)

		ns := make(S, r-l)
		copy(ns, s[l:r])
		lenl := mid - l
		lenr := r - mid

		{
			sl := make(S, lenl)
			copy(sl, s[l:l+lenl])

			sr := make(S, lenr+r-l)
			copy(sr, s[mid:mid+lenr])
			copy(sr[lenr:], ns)

			Reverse(sl)

			zl := zAlgorithm(sl)
			zr := zAlgorithm(sr)

			for i := l; i < mid; i++ {
				length := mid - i
				nl, nr := i, mid

				if i > l {
					nl = max(l, nl-zl[length])
				}

				nr = min(r, nr+zr[lenr+(i-l)])

				if nr-nl >= length*2 {
					if nl >= 1 && s[nl-1] == s[nl-1+length] {
						continue
					}
					if nr < n && s[nr] == s[nr-length] {
						continue
					}

					res = append(res, triple{nl, nr, length})
				}
			}
		}

		{
			sr := make(S, lenr)
			copy(sr, s[mid:mid+lenr])

			sl := make(S, lenl+r-l)
			copy(sl, s[l:mid])
			Reverse(sl[:lenl])
			copy(sl[lenl:], ns)
			Reverse(sl[lenl:])

			zl := zAlgorithm(sl)
			zr := zAlgorithm(sr)

			for i := mid + 1; i <= r; i++ {
				length := i - mid
				nl, nr := mid, i

				if i < r {
					nr = min(r, nr+zr[length])
				}

				nl = max(l, nl-zl[lenl+(r-i)])

				if nr-nl >= length*2 {
					if nl >= 1 && s[nl-1] == s[nl-1+length] {
						continue
					}
					if nr < n && s[nr] == s[nr-length] {
						continue
					}
					res = append(res, triple{nl, nr, length})
				}
			}
		}
	}

	dfs(0, n)

	sort.Slice(res, func(i, j int) bool {
		if res[i].First != res[j].First {
			return res[i].First < res[j].First
		}
		if res[i].Second != res[j].Second {
			return res[i].Second < res[j].Second
		}
		return res[i].Third < res[j].Third
	})
	pl, pr := -1, -1
	for _, run := range res {
		l, r, length := run.First, run.Second, run.Third
		if l == pl && r == pr {
			continue
		}
		pl, pr = l, r
		f(length, l, r)
	}
}

type triple struct {
	First, Second, Third int
}

func Reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func zAlgorithm[S ~[]E, E comparable](s S) []int {
	n := len(s)
	res := make([]int, n)

	for i, j := 1, 0; i < n; i++ {
		if i+res[i-j] < j+res[j] {
			res[i] = res[i-j]
		} else {
			res[i] = max(j+res[j]-i, 0)
			for i+res[i] < n && s[i+res[i]] == s[res[i]] {
				res[i]++
			}
			j = i
		}
	}

	res[0] = n
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
