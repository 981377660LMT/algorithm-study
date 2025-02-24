// https://github.com/old-yan/CP-template/blob/f58da5bf9882f44933df6dae5fb10141a6cbb18d/DS/LinearDSU.h#L1
// ​数据结构：线性并查集
//
// ​练习题目：
//
// 1. [1488. 避免洪水泛滥](https://leetcode.cn/problems/avoid-flood-in-the-city/)
// 2. [P4145 上帝造题的七分钟 2 / 花神游历各国](https://www.luogu.com.cn/problem/P4145)
// 3. [小苯的蓄水池（hard）](https://ac.nowcoder.com/acm/problem/281338)
//
// 当处理一般问题时，并查集的复杂度往往并非线性，而是额外带有一个对数或者阿克曼函数；当然，运行速度并不慢。
// 在只有相邻元素才会发生合并的场景下，并查集的复杂度可以优化到线性。
// !本模板与普通并查集的一大不同是，本模板的每个连通块一定是一个连续不间断的区间。
//
// api:
// - NewLinearDSU(n uint32, maintainGroupSize bool) *LinearDSU
// !- (t *LinearDSU) UniteAfter(i uint32) bool  // 尝试将 i 与其后一个元素合并，若成功返回 true
// - (t *LinearDSU) FindHead(i uint32) uint32  // 查询分组首元素
// - (t *LinearDSU) FindTail(i uint32) uint32  // 查询分组尾元素
// - (t *LinearDSU) FindPrev(i uint32) uint32  // 寻找元素所在分组的上一个分组的尾元素
// - (t *LinearDSU) FindNext(i uint32) uint32  // 寻找元素所在分组的下一个分组的首元素
// - (t *LinearDSU) Size(i uint32) uint32 // 查询分组大小
// - (t *LinearDSU) InSameGroup(left, right uint32) bool
// - (t *LinearDSU) IsHead(i uint32) bool
// - (t *LinearDSU) IsTail(i uint32) bool
// - (t *LinearDSU) GroupCount() uint32
// - (t *LinearDSU) Heads() []uint32
// - (t *LinearDSU) Tails() []uint32
// - (t *LinearDSU) String() string

package main

import (
	"fmt"
	"math/bits"
	"strings"
)

func main() {
	dsu := NewLinearDSU(10)
	fmt.Println(dsu)

	// 查询 3 和 6 的关系
	fmt.Println("3 and 6 in same group?  ", dsu.InSameGroup(3, 6))

	dsu.UniteAfter(3)
	dsu.UniteAfter(5)
	dsu.UniteAfter(8)
	dsu.UniteAfter(4)
	dsu.UniteAfter(6)

	fmt.Println(dsu)

	// 查询 6 所在的分组首元素和尾元素
	fmt.Println("6 is now in which group:", dsu.FindHead(6), "~", dsu.FindTail(6))
	// 查询 3 所在的分组首领
	fmt.Println("3 is now in which group:", dsu.FindHead(3), "~", dsu.FindTail(3))
	// 查询 3 和 6 的关系
	fmt.Println("3 and 6 in same group?  ", dsu.InSameGroup(3, 6))

	heads := dsu.Heads()
	for _, a := range heads {
		fmt.Println(a, "is a head")
	}

	tails := dsu.Tails()
	for _, a := range tails {
		fmt.Println(a, "is a tail")
	}
}

const (
	MASK_SIZE  = 64
	MASK_WIDTH = 6
)

type LinearDSU struct {
	tail       []uint32
	masks      []uint64
	groupSize  []uint32
	size       uint32
	groupCount uint32
}

func NewLinearDSU(n uint32) *LinearDSU {
	t := &LinearDSU{}
	t.build(n)
	return t
}

func (t *LinearDSU) Size(i uint32) uint32 {
	tail := t.FindTail(i)
	return t.groupSize[tail]
}

// FindPrev 寻找元素所在分组的上一个分组的尾元素.
// 如果不存在，返回 -1.
func (t *LinearDSU) FindPrev(i uint32) uint32 {
	return t.FindHead(i) - 1
}

// FindNext 寻找元素所在分组的下一个分组的首元素.
// 如果不存在，返回 n.
func (t *LinearDSU) FindNext(i uint32) uint32 {
	return t.FindTail(i) + 1
}

// 返回 i 所在分组的首元素.
func (t *LinearDSU) FindHead(i uint32) uint32 {
	tail := t.FindTail(i)
	return tail - t.groupSize[tail] + 1
}

// 返回 i 所在分组的尾元素.
func (t *LinearDSU) FindTail(i uint32) uint32 {
	quot := i >> MASK_WIDTH
	rem := i & (MASK_SIZE - 1)
	if ((t.masks[quot] >> rem) & 1) == 1 {
		return i
	}
	mask := t.masks[quot] & ^((1 << rem) - 1)
	if mask != 0 {
		return (quot << MASK_WIDTH) ^ uint32(bits.TrailingZeros64(mask))
	}
	t.tail[quot] = t.findBlock(t.tail[quot])
	return (t.tail[quot] << MASK_WIDTH) ^ uint32(bits.TrailingZeros64(t.masks[t.tail[quot]]))
}

// UniteAfter 尝试将 i 与其后一个元素合并，若成功返回 true.
func (t *LinearDSU) UniteAfter(i uint32) bool {
	quot := i >> MASK_WIDTH
	rem := i & (MASK_SIZE - 1)
	if (t.masks[quot]>>rem)&1 != 0 {
		t.masks[quot] &= ^(uint64(1) << rem)
		if rem+1 == MASK_SIZE {
			if t.masks[quot+1] != 0 {
				t.tail[quot] = quot + 1
			} else {
				t.tail[quot] = t.tail[quot+1]
			}
		}
		t.groupSize[t.FindTail(i)] += t.groupSize[i]
		t.groupCount--
		return true
	}
	return false
}

func (t *LinearDSU) InSameGroup(left, right uint32) bool {
	return t.FindTail(left) >= right
}

// IsHead 判断 i 是否为分组的首元素.
func (t *LinearDSU) IsHead(i uint32) bool {
	if i == 0 {
		return true
	}
	return t.IsTail(i - 1)
}

// IsTail 判断 i 是否为分组的尾元素.
func (t *LinearDSU) IsTail(i uint32) bool {
	quot := i >> MASK_WIDTH
	rem := i & (MASK_SIZE - 1)
	return ((t.masks[quot] >> rem) & 1) == 1
}

func (t *LinearDSU) GroupCount() uint32 {
	return t.groupCount
}

// Heads 返回所有分组的首元素.
func (t *LinearDSU) Heads() []uint32 {
	heads := make([]uint32, 0, t.groupCount)
	var i uint32 = 0
	for i < t.size {
		heads = append(heads, i)
		i = t.FindTail(i) + 1
	}
	return heads
}

// Tails 返回所有分组的尾元素.
func (t *LinearDSU) Tails() []uint32 {
	tails := make([]uint32, 0, t.groupCount)
	i := t.FindTail(0)
	for {
		tails = append(tails, i)
		next := t.FindTail(i + 1)
		if next >= t.size {
			break
		}
		i = next
	}
	return tails
}

func (t *LinearDSU) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	first := true
	for l := uint32(0); l < t.size; {
		r := t.FindTail(l)
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString("[")
		sb.WriteString(fmt.Sprintf("%d, %d", l, r))
		sb.WriteString("]")
		l = r + 1
	}
	sb.WriteString("]")
	return sb.String()
}

func (t *LinearDSU) build(n uint32) {
	t.size = n
	if n == 0 {
		return
	}
	blockCount := (n + MASK_SIZE - 1) >> MASK_WIDTH
	t.groupCount = n
	t.masks = make([]uint64, blockCount)
	t.tail = make([]uint32, blockCount)
	for i := range t.masks {
		t.masks[i] = ^uint64(0)
		t.tail[i] = uint32(i)
	}
	t.groupSize = make([]uint32, n)
	for i := range t.groupSize {
		t.groupSize[i] = 1
	}
}

func (t *LinearDSU) findBlock(q uint32) uint32 {
	if t.masks[q] != 0 {
		return q
	}
	t.tail[q] = t.findBlock(t.tail[q])
	return t.tail[q]
}
