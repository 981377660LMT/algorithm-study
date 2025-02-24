// https://github.com/old-yan/CP-template/blob/f58da5bf9882f44933df6dae5fb10141a6cbb18d/DS/LinearDSU.h#L1

package main

import (
	"fmt"
	"math/bits"
	"strings"
)

const (
	MASK_SIZE  = 64
	MASK_WIDTH = 6 // (64 / 32 + 4) 的结果为 6
)

// Table 实现了类似 C++ 版本的线性不相交集合
type Table struct {
	masks             []uint64 // 每个元素为一个掩码，初值为全 1
	tail              []uint32 // 用于路径压缩的辅助数组
	groupSize         []uint32 // 只有在 maintainGroupSize 为 true 时才维护，每个位置对应所在组的大小
	size              uint32   // 总元素个数
	groupCnt          uint32   // 当前组数
	maintainGroupSize bool     // 是否维护组大小
}

// NewTable 创建一个新的 Table，n 表示元素个数，maintainGroupSize 表示是否维护组大小
func NewTable(n uint32, maintainGroupSize bool) *Table {
	t := &Table{
		maintainGroupSize: maintainGroupSize,
	}
	t.Resize(n)
	return t
}

// Resize 调整 Table 的大小
func (t *Table) Resize(n uint32) {
	t.size = n
	if n == 0 {
		return
	}
	t.groupCnt = n
	// 计算需要多少个 mask 元素
	numMasks := int((n + MASK_SIZE - 1) >> MASK_WIDTH)
	t.masks = make([]uint64, numMasks)
	for i := 0; i < numMasks; i++ {
		t.masks[i] = ^uint64(0) // 所有位都为1
	}
	t.tail = make([]uint32, numMasks)
	for i := 0; i < numMasks; i++ {
		t.tail[i] = uint32(i)
	}
	if t.maintainGroupSize {
		t.groupSize = make([]uint32, n)
		for i := range t.groupSize {
			t.groupSize[i] = 1
		}
	}
}

// findSet 是内部辅助函数，对 tail 数组做路径压缩
func (t *Table) findSet(q uint32) uint32 {
	if t.masks[q] != 0 {
		return q
	}
	t.tail[q] = t.findSet(t.tail[q])
	return t.tail[q]
}

// FindTail 返回下标 i 所在组的尾部下标
func (t *Table) FindTail(i uint32) uint32 {
	quot := i >> MASK_WIDTH                // i / 64
	rem := i & (MASK_SIZE - 1)             // i mod 64
	if ((t.masks[quot] >> rem) & 1) == 1 { // 如果当前位为1，则 i 就是尾部
		return i
	} else {
		// 计算 t.masks[quot] & -(1<<rem)
		shifted := uint64(1) << rem
		mask := t.masks[quot] & uint64(-int64(shifted))
		if mask != 0 {
			// 由于 (quot<<MASK_WIDTH) 下低 MASK_WIDTH 位为0，两者相加等同于位异或
			return (quot << MASK_WIDTH) + uint32(bits.TrailingZeros64(mask))
		} else {
			t.tail[quot] = t.findSet(t.tail[quot])
			return (t.tail[quot] << MASK_WIDTH) + uint32(bits.TrailingZeros64(t.masks[t.tail[quot]]))
		}
	}
}

// FindHead 返回 i 所在组的头部下标（只有 maintainGroupSize 为 true 时可用）
func (t *Table) FindHead(i uint32) uint32 {
	if !t.maintainGroupSize {
		panic("maintainGroupSize 必须为 true 才能调用 FindHead")
	}
	tail := t.FindTail(i)
	return tail - t.groupSize[tail] + 1
}

// FindPrev 返回 i 所在组的前一个位置，即组头的前一个位置
func (t *Table) FindPrev(i uint32) uint32 {
	return t.FindHead(i) - 1
}

// FindNext 返回 i 所在组的下一个位置，即组尾的后一个位置
func (t *Table) FindNext(i uint32) uint32 {
	return t.FindTail(i) + 1
}

// GroupSize 返回 i 所在组的大小（maintainGroupSize 为 true 时有效）
func (t *Table) GroupSize(i uint32) uint32 {
	if !t.maintainGroupSize {
		panic("maintainGroupSize 必须为 true 才能调用 GroupSize")
	}
	tail := t.FindTail(i)
	return t.groupSize[tail]
}

// UniteAfter 尝试将 i 与其后一个元素合并，若成功返回 true
func (t *Table) UniteAfter(i uint32) bool {
	quot := i >> MASK_WIDTH
	rem := i & (MASK_SIZE - 1)
	if ((t.masks[quot] >> rem) & 1) == 1 {
		// 清除该位
		t.masks[quot] &= ^(uint64(1) << rem)
		if rem+1 == MASK_SIZE {
			// 如果已经到了当前 mask 的末尾，则更新 tail
			if int(quot)+1 < len(t.masks) && t.masks[quot+1] != 0 {
				t.tail[quot] = uint32(quot + 1)
			} else if int(quot)+1 < len(t.masks) {
				t.tail[quot] = t.tail[quot+1]
			}
		}
		if t.maintainGroupSize {
			tail := t.FindTail(i)
			t.groupSize[tail] += t.groupSize[i]
		}
		t.groupCnt--
		return true
	}
	return false
}

// InSameGroup 判断 left 和 right 是否在同一组
func (t *Table) InSameGroup(left, right uint32) bool {
	return t.FindTail(left) >= right
}

// IsHead 判断下标 i 是否为一个组的头部
func (t *Table) IsHead(i uint32) bool {
	if i == 0 {
		return true
	}
	return t.IsTail(i - 1)
}

// IsTail 判断下标 i 是否为一个组的尾部
func (t *Table) IsTail(i uint32) bool {
	quot := i >> MASK_WIDTH
	rem := i & (MASK_SIZE - 1)
	return ((t.masks[quot] >> rem) & 1) == 1
}

// Count 返回当前组的个数
func (t *Table) Count() uint32 {
	return t.groupCnt
}

// Heads 返回所有组头的下标集合
func (t *Table) Heads() []uint32 {
	heads := make([]uint32, 0, t.groupCnt)
	var i uint32 = 0
	for i < t.size {
		heads = append(heads, i)
		i = t.FindTail(i) + 1
	}
	return heads
}

// Tails 返回所有组尾的下标集合
func (t *Table) Tails() []uint32 {
	tails := make([]uint32, 0, t.groupCnt)
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

// String 返回 Table 的字符串表示，格式类似于 [[l, r], [l, r], ...]
func (t *Table) String() string {
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
