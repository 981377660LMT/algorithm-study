// 使用方式类似于AC自动机:
// KMP(pattern)：构造函数, pattern为模式串.
// IndexOfAll(s,start): 返回模式串在s中出现的所有位置.
// Move(pos, char): 从当前状态pos沿着char移动到下一个状态, 如果不存在则移动到fail指针指向的状态.
// IsMatched(pos): 判断当前状态pos是否为匹配状态.

package main

import "fmt"

func main() {
	a, b := "ababab", "a"
	fmt.Println(IndexOfAll(a, b, 0))
}

func GetNext(pattern string) []int {
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
func IndexOfAll(longer string, shorter string, position int) []int {
	if len(shorter) == 0 {
		return []int{0}
	}
	if len(longer) < len(shorter) {
		return nil
	}
	res := []int{}
	next := GetNext(shorter)
	hitJ := 0
	for i := position; i < len(longer); i++ {
		for hitJ > 0 && longer[i] != shorter[hitJ] {
			hitJ = next[hitJ-1]
		}
		if longer[i] == shorter[hitJ] {
			hitJ++
		}
		if hitJ == len(shorter) {
			res = append(res, i-len(shorter)+1)
			hitJ = next[hitJ-1] // 不允许重叠时 hitJ = 0
		}
	}
	return res
}

// 单模式串匹配
type KMP struct {
	next    []int
	pattern string
}

func NewKMP(pattern string) *KMP {
	return &KMP{
		next:    GetNext(pattern),
		pattern: pattern,
	}
}

// `o(n+m)`求搜索串 longer 中所有匹配 pattern 的位置.
//  findAll/indexOfAll
func (k *KMP) IndexOfAll(longer string, start int) []int {
	var res []int
	pos := 0
	for i := start; i < len(longer); i++ {
		pos = k.Move(pos, longer[i])
		if k.IsMatched(pos) {
			res = append(res, i-len(k.pattern)+1)
			pos = k.next[pos-1]
		}
	}
	return res
}

func (k *KMP) IndexOf(longer string, start int) int {
	pos := 0
	for i := start; i < len(longer); i++ {
		pos = k.Move(pos, longer[i])
		if k.IsMatched(pos) {
			return i - len(k.pattern) + 1
		}
	}
	return -1
}

func (k *KMP) Move(pos int, char byte) int {
	if pos < 0 || pos >= len(k.pattern) {
		panic("pos out of range")
	}
	for pos > 0 && char != k.pattern[pos] {
		pos = k.next[pos-1]
	}
	if char == k.pattern[pos] {
		pos++
	}
	return pos
}

func (k *KMP) IsMatched(pos int) bool {
	return pos == len(k.pattern)
}

// 求s的前缀[0:i+1)的最小周期.如果不存在,则返回0.
//  0<=i<len(s).
func (k *KMP) Period(i int) int {
	res := i + 1 - k.next[i]
	if res > 0 && (i+1) > res && (i+1)%res == 0 {
		return res
	}
	return 0
}
