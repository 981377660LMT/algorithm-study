// bitset 位集
// 参考 https://github.dev/EndlessCheng/codeforces-go/blob/master/misc/atcoder/abc274/e

package main

import (
	"fmt"
	"math/bits"
	"strings"
)

func main() {
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
	bitSet2.Set(1)
	bitSet2.Set(2)
	bitSet2.Set(3)
	fmt.Println(bitSet.OrOnesCount(bitSet2))
	fmt.Println(bitSet.AndOnesCount(bitSet2))
	fmt.Println(bits.TrailingZeros(4503599627370496), 4503599627370496&(-4503599627370496))
}

const UNIT = bits.UintSize

type BitSet struct {
	size   int
	states []uint
}

func NewBitSet(nbits int) *BitSet {
	if nbits < 0 {
		panic("initSize must be non-negative")
	}

	return &BitSet{
		size:   0,
		states: make([]uint, (nbits+UNIT-1)/UNIT), // math.Ceil(nbits/UNIT)
	}

}

func (set *BitSet) Set(index int) {
	set.ensureCapacity(index + 1)
	if set.states[index/UNIT]&(1<<(index%UNIT)) == 0 {
		set.size++
		set.states[index/UNIT] |= 1 << (index % UNIT)
	}
}

func (set *BitSet) Reset(index int) {
	set.ensureCapacity(index + 1)
	if set.states[index/UNIT]&(1<<(index%UNIT)) != 0 {
		set.size--
		set.states[index/UNIT] ^= 1 << (index % UNIT)
	}
}

func (set *BitSet) Flip(index int) {
	set.ensureCapacity(index + 1)
	if set.states[index/UNIT]&(1<<(index%UNIT)) == 0 {
		set.size++
	} else {
		set.size--
	}
	set.states[index/UNIT] ^= 1 << (index % UNIT)
}

func (set *BitSet) Has(index int) bool {
	set.ensureCapacity(index + 1)
	return set.states[index/UNIT]&(1<<(index%UNIT)) != 0
}

func (set *BitSet) OnesCount() int {
	return set.size
}

func (set *BitSet) AndOnesCount(other *BitSet) int {
	res := 0
	minLen := min(len(set.states), len(other.states))
	for i := 0; i < minLen; i++ {
		res += bits.OnesCount(set.states[i] & other.states[i])
	}
	return res
}

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

func (set *BitSet) String() string {
	res := []string{"BitSet{"}
	indexes := []string{}
	for i, v := range set.states {
		for ; v > 0; v &= v - 1 {
			indexes = append(indexes, fmt.Sprintf("%d", i*UNIT+bits.TrailingZeros(v)))
		}
	}

	res = append(res, strings.Join(indexes, ", "), "}")
	return strings.Join(res, "")
}

func (set *BitSet) ensureCapacity(nbits int) {
	if nbits > len(set.states)*UNIT {
		set.resize(nbits)
	}
}

func (set *BitSet) resize(nbits int) {
	newLength := len(set.states) * 2
	if nbits > newLength*UNIT {
		newLength = (nbits + UNIT - 1) / UNIT
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
