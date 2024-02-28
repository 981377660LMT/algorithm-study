// 部分可持久化KMP/动态KMP/可撤销KMP
// 这里的kmp说的是一个可以查询border的数据结构，修改就是往后添加一个字符
// 历史版本可以查询，最新版本可以修改和查询
// https://zhuanlan.zhihu.com/p/527154301
// https://www.cnblogs.com/Blue233333/p/8241503.html

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	CF1721E()
}

// Prefix Function Queries
// https://www.luogu.com.cn/problem/CF1721E
// 给定字符串s以及q个串ti，求将s分别与每个ti拼接起来后，最靠右的|ti|个前缀的 border 长度。询问间相互独立。。
func CF1721E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	K := NewKmpUndoable(len(s), func(i int) byte { return s[i] })

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var t string
		fmt.Fscan(in, &t)

		state := K.GetState()
		for i := 0; i < len(t); i++ {
			K.Append(t[i])
			fmt.Fprint(out, K.Query(), " ")
		}
		fmt.Fprintln(out)

		K.Rollback(state)
	}
}

type Int interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

type KmpUndoable[S Int] struct {
	State int32 // 状态即为当前字符长度
	Nums  []S
	fail  []int32
}

// Kmp 初始状态.
func NewKmpUndoable[S Int](n int, f func(i int) S) *KmpUndoable[S] {
	nums := make([]S, n)
	for i := 0; i < n; i++ {
		nums[i] = f(i)
	}
	fail := make([]int32, n)
	j := int32(0)
	for i := 1; i < n; i++ {
		vi := nums[i]
		for j > 0 && vi != nums[j] {
			j = fail[j-1]
		}
		if vi == nums[j] {
			j++
		}
		fail[i] = j
	}
	return &KmpUndoable[S]{State: int32(n), Nums: nums, fail: fail}
}

// 在当前字符串末尾追加一个字符.
func (kmp *KmpUndoable[S]) Append(c S) {
	kmp.expand(int(kmp.State) + 1)
	kmp.Nums[kmp.State] = c
	kmp.State++
	kmp.updateFail(kmp.State-1, kmp.State)
}

// 查询当前字符串的border长度.
func (kmp *KmpUndoable[S]) Query() int {
	return int(kmp.fail[max32(kmp.State-1, 0)])
}

func (kmp *KmpUndoable[S]) Move(pos int, ord S) int {
	if pos < 0 || pos >= int(kmp.State) {
		panic("pos out of range")
	}
	pos32 := int32(pos)
	for pos32 > 0 && ord != kmp.Nums[pos32] {
		pos32 = kmp.fail[pos32-1]
	}
	if ord == kmp.Nums[pos32] {
		pos32++
	}
	return int(pos32)
}

func (kmp *KmpUndoable[S]) Accept(pos int) bool {
	return pos == int(kmp.State)
}

func (kmp *KmpUndoable[S]) GetState() int {
	return int(kmp.State)
}

func (kmp *KmpUndoable[S]) Undo() bool {
	if kmp.State == 0 {
		return false
	}
	kmp.State--
	return true
}

func (kmp *KmpUndoable[S]) Rollback(state int) bool {
	if state < 0 || state > int(kmp.State) {
		return false
	} else {
		kmp.State = int32(state)
		return true
	}
}

// 求s的前缀[0:i+1)的最小周期.如果不存在,则返回0.
//
//	0<=i<len(s).
func (kmp *KmpUndoable[S]) Period(i int) int {
	res := i + 1 - int(kmp.fail[i])
	if res > 0 && (i+1) > res && (i+1)%res == 0 {
		return res
	}
	return 0
}

// 计算字符串从rawLen开始到targetLen的前缀函数，前提是<rawLen的部分已经求出.
func (kmp *KmpUndoable[S]) updateFail(rawLen, targetLen int32) {
	j := int32(0)
	if rawLen > 0 {
		j = kmp.fail[rawLen-1]
	}
	nums, fail := kmp.Nums, kmp.fail
	for i := max32(rawLen, 1); i < targetLen; i++ {
		for j > 0 && nums[j] != nums[i] {
			if (fail[j-1] <= j>>1) || nums[i] == nums[fail[j-1]] {
				j = fail[j-1]
			} else {
				period := j - fail[j-1]
				j = fail[j%period+period-1]
			}
		}
		if j > 0 || nums[j] == nums[i] {
			j++
			fail[i] = j
		} else {
			fail[i] = 0
		}
	}
}

func (kmp *KmpUndoable[S]) expand(size int) {
	for len(kmp.Nums) < size {
		kmp.Nums = append(kmp.Nums, 0)
		kmp.fail = append(kmp.fail, 0)
	}
}

func max(a, b int) int {
	if a > b {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
