// 使用方式类似于AC自动机:
// KMP(pattern)：构造函数, pattern为模式串.
// IndexOfAll(s,start): 返回模式串在s中出现的所有位置.
// Move(pos, char): 从当前状态pos沿着char移动到下一个状态, 如果不存在则移动到fail指针指向的状态.
// IsMatched(pos): 判断当前状态pos是否为匹配状态.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// a, b := "ababab", "a"
	// fmt.Println(IndexOfAll(a, b, 0))

	// P4824()

}

// 面试题 17.17. 多次搜索
// https://leetcode.cn/problems/multi-search-lcci/
func multiSearch(big string, smalls []string) [][]int {
	res := make([][]int, len(smalls))
	for i, small := range smalls {
		res[i] = IndexOfAll(big, small, 0, nil)
	}
	return res
}

// https://www.luogu.com.cn/problem/P4824
// 在longer中不断删除shorter，求剩下的字符串.
func P4824() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var longer, shorter string
	fmt.Fscan(in, &longer, &shorter)

	kmp := NewKMP(shorter)
	pos := 0
	stack := make([]int, 0, len(longer))
	posRecord := make([]int, len(longer))
	for i := range longer {
		pos = kmp.Move(pos, int(longer[i]))
		posRecord[i] = pos
		stack = append(stack, i)
		if kmp.Accept(pos) {
			stack = stack[:len(stack)-len(shorter)]
			if len(stack) > 0 {
				pos = posRecord[stack[len(stack)-1]]
			} else {
				pos = 0
			}
		}
	}

	res := make([]byte, 0, len(stack))
	for _, v := range stack {
		res = append(res, longer[v])
	}
	fmt.Fprintln(out, string(res))
}

type Str = string

func GetNext(pattern Str) []int {
	next := make([]int, len(pattern))
	j := 0
	for i := 1; i < len(pattern); i++ {
		for j > 0 && pattern[i] != pattern[j] {
			j = next[j-1]
		}
		if pattern[i] == pattern[j] {
			j++
		}
		next[i] = j
	}
	return next
}

// `O(n+m)` 寻找 `shorter` 在 `longer` 中的所有匹配位置.
// nexts 数组为nil时, 会调用GetNext(shorter)求nexts数组.
func IndexOfAll(longer Str, shorter Str, position int, nexts []int) []int {
	if len(shorter) == 0 {
		return []int{}
	}
	if len(longer) < len(shorter) {
		return nil
	}
	res := []int{}
	if nexts == nil {
		nexts = GetNext(shorter)
	}
	hitJ := 0
	for i := position; i < len(longer); i++ {
		for hitJ > 0 && longer[i] != shorter[hitJ] {
			hitJ = nexts[hitJ-1]
		}
		if longer[i] == shorter[hitJ] {
			hitJ++
		}
		if hitJ == len(shorter) {
			res = append(res, i-len(shorter)+1)
			hitJ = nexts[hitJ-1] // 不允许重叠时 hitJ = 0
		}
	}
	return res
}

func CountIndexOfAll(longer Str, shorter Str, position int, nexts []int) int {
	if len(shorter) == 0 {
		return 0
	}
	if len(longer) < len(shorter) {
		return 0
	}
	res := 0
	if nexts == nil {
		nexts = GetNext(shorter)
	}
	hitJ := 0
	for i := position; i < len(longer); i++ {
		for hitJ > 0 && longer[i] != shorter[hitJ] {
			hitJ = nexts[hitJ-1]
		}
		if longer[i] == shorter[hitJ] {
			hitJ++
		}
		if hitJ == len(shorter) {
			res++
			hitJ = nexts[hitJ-1] // 不允许重叠时 hitJ = 0
		}
	}
	return res
}

// 单模式串匹配
type KMP struct {
	next    []int
	pattern Str
}

func NewKMP(pattern Str) *KMP {
	return &KMP{
		next:    GetNext(pattern),
		pattern: pattern,
	}
}

// `o(n+m)`求搜索串 longer 中所有匹配 pattern 的位置.
//
//	findAll/indexOfAll
func (k *KMP) IndexOfAll(longer Str, start int) []int {
	if len(longer) < len(k.pattern) {
		return nil
	}
	var res []int
	pos := 0
	for i := start; i < len(longer); i++ {
		pos = k.Move(pos, int(longer[i]))
		if k.Accept(pos) {
			res = append(res, i-len(k.pattern)+1)
			pos = k.next[pos-1]
		}
	}
	return res
}

// `o(n+m)`求文本串 longer 中所有匹配 pattern 的次数.
//
//	findAll/indexOfAll
func (k *KMP) CountIndexOfAll(longer string, start int) int {
	if len(longer) < len(k.pattern) {
		return 0
	}
	res := 0
	pos := 0
	for i := start; i < len(longer); i++ {
		pos = k.Move(pos, int(longer[i]))
		if k.Accept(pos) {
			res++
			pos = k.next[pos-1]
		}
	}
	return res
}

func (k *KMP) IndexOf(longer Str, start int) int {
	pos := 0
	for i := start; i < len(longer); i++ {
		pos = k.Move(pos, int(longer[i]))
		if k.Accept(pos) {
			return i - len(k.pattern) + 1
		}
	}
	return -1
}

func (k *KMP) Move(pos int, ord int) int {
	if pos < 0 || pos >= len(k.pattern) {
		panic("pos out of range")
	}
	for pos > 0 && ord != int(k.pattern[pos]) {
		pos = k.next[pos-1]
	}
	if ord == int(k.pattern[pos]) {
		pos++
	}
	return pos
}

func (k *KMP) Accept(pos int) bool {
	return pos == len(k.pattern)
}

// 求s的前缀[0:i+1)的最小周期.如果不存在,则返回0.
//
//	0<=i<len(s).
func (k *KMP) Period(i int) int {
	res := i + 1 - k.next[i]
	if res > 0 && (i+1) > res && (i+1)%res == 0 {
		return res
	}
	return 0
}
