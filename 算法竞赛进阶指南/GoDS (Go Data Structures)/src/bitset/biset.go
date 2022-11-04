// bitset 位集
// 参考 https://github.dev/EndlessCheng/codeforces-go/blob/master/misc/atcoder/abc274/e

package bitset

import (
	"fmt"
	"math/bits"
	"strings"
)

func demo() {
	bitSet := NewBitSet(64)
	bitSet.Set(1)
	fmt.Println(bitSet.OnesCount())
	bitSet.Set(2)
	bitSet.Set(64)
	fmt.Println(bitSet)
	bitSet.Reset(2)
	fmt.Println(bitSet)
	bitSet.Flip(64)
	fmt.Println(bitSet)
	fmt.Println(bitSet.Has(64))

	bitSet2 := NewBitSet(64)
	bitSet2.Set(12)
	bitSet2.Set(23)
	bitSet2.Set(3)
	bitSet2.Set(30)
	fmt.Println(bitSet.OrOnesCount(bitSet2))
	fmt.Println(bitSet.AndOnesCount(bitSet2))

	bitSet.Or(bitSet2)
	fmt.Println(bitSet, bitSet.OnesCount())
	bitSet.And(bitSet2)
	fmt.Println(bitSet, bitSet.OnesCount())

	fmt.Println(bitSet.IsSubset(bitSet2))
	fmt.Println(bitSet2.IsSubset(bitSet))
}

func NewBitSet(nbits int) *BitSet {
	if nbits < 0 {
		panic("initSize must be non-negative")
	}

	return &BitSet{
		size:   0,
		states: make([]uint, (nbits+unit-1)/unit), // math.Ceil(nbits/unit)
	}
}

const unit = bits.UintSize

type BitSet struct {
	size   int
	states []uint
}

func (set *BitSet) Set(index int) {
	set.ensureCapacity(index + 1)
	if set.states[index/unit]&(1<<(index%unit)) == 0 {
		set.size++
		set.states[index/unit] |= 1 << (index % unit)
	}
}

func (set *BitSet) Reset(index int) {
	set.ensureCapacity(index + 1)
	if set.states[index/unit]&(1<<(index%unit)) != 0 {
		set.size--
		set.states[index/unit] ^= 1 << (index % unit)
	}
}

func (set *BitSet) Flip(index int) {
	set.ensureCapacity(index + 1)
	if set.states[index/unit]&(1<<(index%unit)) == 0 {
		set.size++
	} else {
		set.size--
	}
	set.states[index/unit] ^= 1 << (index % unit)
}

func (set *BitSet) Has(index int) bool {
	set.ensureCapacity(index + 1)
	return set.states[index/unit]&(1<<(index%unit)) != 0
}

func (set *BitSet) OnesCount() int {
	return set.size
}

// 当前位集 & other 位集
func (set *BitSet) And(other *BitSet) {
	minLen := min(len(set.states), len(other.states))
	newCount := 0
	for i := 0; i < minLen; i++ {
		set.states[i] &= other.states[i]
		newCount += bits.OnesCount(set.states[i])
	}
	for i := minLen; i < len(set.states); i++ {
		set.states[i] = 0
	}

	set.size = newCount
}

// 交集的元素个数
func (set *BitSet) AndOnesCount(other *BitSet) int {
	res := 0
	minLen := min(len(set.states), len(other.states))
	for i := 0; i < minLen; i++ {
		res += bits.OnesCount(set.states[i] & other.states[i])
	}
	return res
}

// 当前位集 | other 位集
func (set *BitSet) Or(other *BitSet) {
	maxLen := max(len(set.states), len(other.states))
	set.ensureCapacity(maxLen * unit)
	newCount := 0
	for i := 0; i < maxLen; i++ {
		if i < len(set.states) && i < len(other.states) {
			set.states[i] |= other.states[i]
		} else if i < len(set.states) {
			continue
		} else {
			set.states[i] = other.states[i]
		}
		newCount += bits.OnesCount(set.states[i])
	}

	set.size = newCount
}

// 并集的元素个数
func (set *BitSet) OrOnesCount(other *BitSet) int {
	res := 0
	maxLen := max(len(set.states), len(other.states))
	for i := 0; i < maxLen; i++ {
		if i < len(set.states) && i < len(other.states) {
			res += bits.OnesCount(set.states[i] | other.states[i])
		} else if i < len(set.states) {
			res += bits.OnesCount(set.states[i])
		} else {
			res += bits.OnesCount(other.states[i])
		}
	}
	return res
}

func (set *BitSet) IsSubset(other *BitSet) bool {
	maxLen := max(len(set.states), len(other.states))
	for i := 0; i < maxLen; i++ {
		if i < len(set.states) && i < len(other.states) {
			if set.states[i]&other.states[i] != set.states[i] {
				return false
			}
		} else if i < len(set.states) {
			if set.states[i] != 0 {
				return false
			}
		}
	}

	return true
}

func (set *BitSet) String() string {
	res := []string{"BitSet{"}
	indexes := []string{}
	for i, v := range set.states {
		for ; v > 0; v &= v - 1 {
			indexes = append(indexes, fmt.Sprintf("%d", i*unit+bits.TrailingZeros(v)))
		}
	}

	res = append(res, strings.Join(indexes, ", "), "}")
	return strings.Join(res, "")
}

func (set *BitSet) ensureCapacity(nbits int) {
	if nbits > len(set.states)*unit {
		set.resize(nbits)
	}
}

func (set *BitSet) resize(nbits int) {
	newLength := len(set.states) * 2
	if nbits > newLength*unit {
		newLength = (nbits + unit - 1) / unit
	}

	newStates := make([]uint, newLength)
	copy(newStates, set.states)
	set.states = newStates
}

func min(a, b int) int {
	if a < b {
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
