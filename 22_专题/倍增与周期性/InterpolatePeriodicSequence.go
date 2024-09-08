// InterpolatePeriodicSequence
// 周期序列插值(一阶/二阶)
// !用于发现周期性序列的循环节
// 例如 123[456][456][456]...
// 适合打表找规律的场合.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	yuki2858()
	// demo()
}

func demo() {
	{
		// 零阶周期序列插值
		arr := []int{1, 2, 4, 5, 6, 8, 9, 8, 9}
		seq1 := NewInterpolatePeriodicSequence(arr)
		for i := 0; i < 20; i++ {
			fmt.Println(seq1.Get(i))
		}
	}

	{
		// 一阶周期序列插值
		arr := []int{0, 1, 2, 4, 5, 6, 7, 8, 9}
		seq2 := NewInterpolateDifferencePeriodicSequence(arr)
		for i := 0; i < 20; i++ {
			fmt.Println(seq2.Get(i))
		}
	}
}

// No.2858 Make a Palindrome (repeatStringToMakePalindrome)
// https://yukicoder.me/problems/no/2858
// 给定一个字符串s.
// 可以重复s任意次，使得新字符串的最长回文子串长度大于等于k.
// 求最少重复几次.如果无解输出-1.
//
// 注意到连接一定字符后，增长的回文长度是一个周期序列.
func yuki2858() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(s string, k int) int {
		lens := []int{0}
		for i := 1; i <= 4; i++ {
			curS := strings.Repeat(s, i)
			maxLen := int32(0)
			LongestPalindromes(int32(len(curS)), func(i, j int32) bool { return curS[i] == curS[j] }, func(start, end int32) {
				maxLen = max32(maxLen, end-start)
			})
			lens = append(lens, int(maxLen))
		}

		S := NewInterpolateDifferencePeriodicSequence(lens)
		res := MinLeft(k+1, func(n int) bool { return S.Get(n) >= k }, 0)
		if res == k+1 {
			return -1
		}
		return res
	}

	var T int32
	fmt.Fscan(in, &T)
	for i := int32(0); i < T; i++ {
		var n, k int
		var s string
		fmt.Fscan(in, &n, &k, &s)
		fmt.Fprintln(out, solve(s, k))
	}
}

// 一阶周期序列插值(差分形如 012[345][345][345]...的周期序列).
// Diff1.
type InterpolateDifferencePeriodicSequence struct {
	offset int
	diff   int
	seq    []int
}

func NewInterpolateDifferencePeriodicSequence(seq []int) *InterpolateDifferencePeriodicSequence {
	seq = append(seq[:0:0], seq...)
	diff := make([]int, 0, len(seq)-1)
	for i := 0; i < len(seq)-1; i++ {
		diff = append(diff, seq[i+1]-seq[i])
	}
	for i, j := 0, len(diff)-1; i < j; i, j = i+1, j-1 {
		diff[i], diff[j] = diff[j], diff[i]
	}
	z := ZAlgoNums(diff)
	z[0] = 0
	max_, maxIndex := int32(-1), -1
	for i, v := range z {
		if v > max_ {
			max_ = v
			maxIndex = i
		}
	}
	n := len(seq)
	d := seq[n-1] - seq[n-maxIndex-1]
	return &InterpolateDifferencePeriodicSequence{offset: maxIndex, diff: d, seq: seq}
}

func (ids *InterpolateDifferencePeriodicSequence) Get(index int) int {
	if index < len(ids.seq) {
		return ids.seq[index]
	}
	if ids.offset == 0 {
		panic("invalid sequence")
	}
	k := (index - (len(ids.seq) - 1) + ids.offset - 1) / ids.offset
	index -= k * ids.offset
	return ids.seq[index] + k*ids.diff
}

// 零阶周期序列插值(形如 123[456][456][456]...的周期序列).
// Diff0.
type InterpolatePeriodicSequence struct {
	offset int
	seq    []int
}

func NewInterpolatePeriodicSequence(seq []int) *InterpolatePeriodicSequence {
	seq = append(seq[:0:0], seq...)
	revSeq := make([]int, len(seq))
	for i := range seq {
		revSeq[i] = seq[len(seq)-1-i]
	}
	z := ZAlgoNums(revSeq)
	z[0] = 0
	max_, maxIndex := int32(-1), -1
	for i, v := range z {
		if v > max_ {
			max_ = v
			maxIndex = i
		}
	}
	return &InterpolatePeriodicSequence{offset: maxIndex, seq: seq}
}

func (ips *InterpolatePeriodicSequence) Get(index int) int {
	if index < len(ips.seq) {
		return ips.seq[index]
	}
	if ips.offset == 0 {
		panic("invalid sequence")
	}
	k := (index - (len(ips.seq) - 1) + ips.offset - 1) / ips.offset
	index -= k * ips.offset
	return ips.seq[index]
}

func ZAlgoNums(arr []int) []int32 {
	n := int32(len(arr))
	z := make([]int32, n)
	if n == 0 {
		return z
	}
	for i, j := int32(1), int32(0); i < n; i++ {
		if j+z[j] <= i {
			z[i] = 0
		} else {
			z[i] = min32(j+z[j]-i, z[i-j])
		}
		for i+z[i] < n && arr[z[i]] == arr[i+z[i]] {
			z[i]++
		}
		if j+z[j] < i+z[i] {
			j = i
		}
	}
	z[0] = n
	return z
}

// 给定一个字符串，返回极长回文子串的区间.这样的极长回文子串最多有 2n-1 个.
// ManacherSimple.
func LongestPalindromes(n int32, equals func(i, j int32) bool, consumer func(start, end int32)) {
	f := func(i, j int32) bool {
		if i > j {
			return false
		}
		if i&1 == 1 {
			return true
		}
		return equals(i>>1, j>>1)
	}
	dp := make([]int32, 2*n-1)
	i, j := int32(0), int32(0)
	for i < 2*n-1 {
		for i-j >= 0 && i+j < 2*n-1 && f(i-j, i+j) {
			j++
		}
		dp[i] = j
		k := int32(1)
		for i-k >= 0 && i+k < 2*n-1 && k < j && dp[i-k] != j-k {
			dp[i+k] = min32(j-k, dp[i-k])
			k++
		}
		i += k
		j = max32(j-k, 0)
	}

	for i := int32(0); i < int32(len(dp)); i++ {
		if dp[i] == 0 {
			continue
		}
		l := (i - dp[i] + 2) >> 1
		r := (i + dp[i] - 1) >> 1
		if l <= r {
			consumer(l, r+1)
		}
	}
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
// left>=lower.
func MinLeft(right int, check func(left int) bool, lower int) int {
	ok, ng := right, lower-1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}
