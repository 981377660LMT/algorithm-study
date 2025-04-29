// 使用方式类似于AC自动机:
// KMP(pattern)：构造函数, pattern为模式串.
// IndexOfAll(s,start): 返回模式串在s中出现的所有位置.
// Move(pos, char): 从当前状态pos沿着char移动到下一个状态, 如果不存在则移动到fail指针指向的状态.
// IsMatched(pos): 判断当前状态pos是否为匹配状态.

package main

type Str = string

func GetNext(pattern Str) []int32 {
	next := make([]int32, len(pattern))
	j := int32(0)
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

// `halfLinkLength[i]`表示`[:i+1]`这一段子串长度不超过串长一半的最长的border长度.
func GetHalfLinkLength(pattern Str, nexts []int32) (halfLinkLength []int32) {
	n := int32(len(pattern))
	depth := make([]int32, n+1) // fail树结点深度
	for i := int32(1); i <= n; i++ {
		parent := nexts[i-1]
		depth[i] = depth[parent] + 1
	}
	halfLinkLength = make([]int32, n)
	pos := int32(0)
	for i := int32(1); i < n; i++ {
		for pos > 0 && pattern[i] != pattern[pos] {
			pos = nexts[pos-1]
		}
		if pattern[i] == pattern[pos] {
			pos++
		}
		for pos > (i+1)>>1 {
			pos = nexts[pos-1]
		}
		halfLinkLength[i] = depth[pos]
	}
	return
}

// `O(n+m)` 寻找 `shorter` 在 `longer` 中的所有匹配位置.
// nexts 数组为nil时, 会调用GetNext(shorter)求nexts数组.
func IndexOfAll(longer Str, shorter Str, position int32, nexts []int32) []int32 {
	if len(shorter) == 0 {
		return []int32{}
	}
	if len(longer) < len(shorter) {
		return nil
	}
	res := []int32{}
	if nexts == nil {
		nexts = GetNext(shorter)
	}
	hitJ := int32(0)
	for i := position; i < int32(len(longer)); i++ {
		for hitJ > 0 && longer[i] != shorter[hitJ] {
			hitJ = nexts[hitJ-1]
		}
		if longer[i] == shorter[hitJ] {
			hitJ++
		}
		if hitJ == int32(len(shorter)) {
			res = append(res, i-int32(len(shorter))+1)
			hitJ = nexts[hitJ-1] // 不允许重叠时 hitJ = 0
		}
	}
	return res
}

func CountIndexOfAll(longer Str, shorter Str, position int, nexts []int32) int32 {
	if len(shorter) == 0 {
		return 0
	}
	if len(longer) < len(shorter) {
		return 0
	}
	res := int32(0)
	if nexts == nil {
		nexts = GetNext(shorter)
	}
	hitJ := int32(0)
	for i := position; i < len(longer); i++ {
		for hitJ > 0 && longer[i] != shorter[hitJ] {
			hitJ = nexts[hitJ-1]
		}
		if longer[i] == shorter[hitJ] {
			hitJ++
		}
		if hitJ == int32(len(shorter)) {
			res++
			hitJ = nexts[hitJ-1] // 不允许重叠时 hitJ = 0
		}
	}
	return res
}

// 单模式串匹配
type KMP struct {
	next    []int32
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
func (k *KMP) IndexOfAll(longer Str, start int32) []int32 {
	if len(longer) < len(k.pattern) {
		return nil
	}
	var res []int32
	pos := int32(0)
	for i := start; i < int32(len(longer)); i++ {
		pos = k.Move(pos, int32(longer[i]))
		if k.Accept(pos) {
			res = append(res, i-int32(len(k.pattern))+1)
			pos = k.next[pos-1]
		}
	}
	return res
}

// `o(n+m)`求文本串 longer 中所有匹配 pattern 的次数.
//
//	findAll/indexOfAll
func (k *KMP) CountIndexOfAll(longer string, start int32) int32 {
	if len(longer) < len(k.pattern) {
		return 0
	}
	res := int32(0)
	pos := int32(0)
	for i := start; i < int32(len(longer)); i++ {
		pos = k.Move(pos, int32(longer[i]))
		if k.Accept(pos) {
			res++
			pos = k.next[pos-1]
		}
	}
	return res
}

func (k *KMP) IndexOf(longer Str, start int32) int32 {
	pos := int32(0)
	for i := start; i < int32(len(longer)); i++ {
		pos = k.Move(pos, int32(longer[i]))
		if k.Accept(pos) {
			return i - int32(len(k.pattern)) + 1
		}
	}
	return -1
}

// 当前文本后缀能匹配的最长模式的前缀.
func (k *KMP) Move(pos int32, ord int32) int32 {
	if pos < 0 || pos >= int32(len(k.pattern)) {
		panic("pos out of range")
	}
	for pos > 0 && ord != int32(k.pattern[pos]) {
		pos = k.next[pos-1]
	}
	if ord == int32(k.pattern[pos]) {
		pos++
	}
	return pos
}

func (k *KMP) Accept(pos int32) bool {
	return pos == int32(len(k.pattern))
}

// 求s的前缀[0:i+1)的最小周期.如果不存在,则返回0.
//
//	0<=i<len(s).
func (k *KMP) Period(i int32) int32 {
	res := i + 1 - k.next[i]
	if res > 0 && (i+1) > res && (i+1)%res == 0 {
		return res
	}
	return 0
}
