/*
Package bitarray implements a bit array.  Useful for tracking bool type values in a space
efficient way.  This is *NOT* a threadsafe package.
*/

// api:
// SetBit / GetBit / ClearBit：设置 / 获取 / 清除指定位置的位。
// Capacity()：返回容量或最大下标；
// Count()：统计已置位的数量；
// Or / And / Nand：位运算；
// ToNums()：把已置位的所有位置以 []uint64 形式输出；
// Reset()：将全部位清零；
// Blocks()：返回可迭代的区块迭代器，用于底层遍历；
// 序列化 / 反序列化：Serialize() / Deserialize() 以及提供 Marshal() / Unmarshal() 等等。

package bitarray

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/bits"
	"sort"
	"unsafe"
)

// #region interface

// BitArray represents a structure that can be used to
// quickly check for existence when using a large number
// of items in a very memory efficient way.
type BitArray interface {
	// SetBit sets the bit at the given position.  This
	// function returns an error if the position is out
	// of range.  A sparse bit array never returns an error.
	SetBit(k uint64) error
	// GetBit gets the bit at the given position.  This
	// function returns an error if the position is out
	// of range.  A sparse bit array never returns an error.
	GetBit(k uint64) (bool, error)
	// GetSetBits gets the position of bits set in the array. Will
	// return as many set bits as can fit in the provided buffer
	// starting from the specified position in the array.
	GetSetBits(from uint64, buffer []uint64) []uint64
	// ClearBit clears the bit at the given position.  This
	// function returns an error if the position is out
	// of range.  A sparse bit array never returns an error.
	ClearBit(k uint64) error
	// Reset sets all values to zero.
	Reset()
	// Blocks returns an iterator to be used to iterate
	// over the bit array.
	Blocks() Iterator
	// Equals returns a bool indicating equality between the
	// two bit arrays.
	Equals(other BitArray) bool
	// Intersects returns a bool indicating if the other bit
	// array intersects with this bit array.
	Intersects(other BitArray) bool
	// Capacity returns either the given capacity of the bit array
	// in the case of a dense bit array or the highest possible
	// seen capacity of the sparse array.
	Capacity() uint64
	// Count returns the number of set bits in this array.
	Count() int
	// Or will bitwise or the two bitarrays and return a new bitarray
	// representing the result.
	Or(other BitArray) BitArray
	// And will bitwise and the two bitarrays and return a new bitarray
	// representing the result.
	And(other BitArray) BitArray
	// Nand will bitwise nand the two bitarrays and return a new bitarray
	// representing the result.
	Nand(other BitArray) BitArray
	// ToNums converts this bit array to the list of numbers contained
	// within it.
	ToNums() []uint64
	// IsEmpty checks to see if any values are set on the bitarray
	IsEmpty() bool
}

// Iterator defines methods used to iterate over a bit array.
type Iterator interface {
	// Next moves the pointer to the next block.  Returns
	// false when no blocks remain.
	Next() bool
	// Value returns the next block and its index
	Value() (uint64, block)
}

// #endregion

// #region bitarray

// bitArray is a struct that maintains state of a bit array.
type bitArray struct {
	blocks  []block
	lowest  uint64
	highest uint64
	anyset  bool
}

func getIndexAndRemainder(k uint64) (uint64, uint64) {
	return k / s, k % s
}

func (ba *bitArray) setLowest() {
	for i := uint64(0); i < uint64(len(ba.blocks)); i++ {
		if ba.blocks[i] == 0 {
			continue
		}

		pos := ba.blocks[i].findRightPosition()
		ba.lowest = (i * s) + pos
		ba.anyset = true
		return
	}

	ba.anyset = false
	ba.lowest = 0
	ba.highest = 0
}

func (ba *bitArray) setHighest() {
	for i := len(ba.blocks) - 1; i >= 0; i-- {
		if ba.blocks[i] == 0 {
			continue
		}

		pos := ba.blocks[i].findLeftPosition()
		ba.highest = (uint64(i) * s) + pos
		ba.anyset = true
		return
	}

	ba.anyset = false
	ba.highest = 0
	ba.lowest = 0
}

// capacity returns the total capacity of the bit array.
func (ba *bitArray) Capacity() uint64 {
	return uint64(len(ba.blocks)) * s
}

// ToNums converts this bitarray to a list of numbers contained within it.
func (ba *bitArray) ToNums() []uint64 {
	nums := make([]uint64, 0, ba.highest-ba.lowest/4)
	for i, block := range ba.blocks {
		block.toNums(uint64(i)*s, &nums)
	}

	return nums
}

// SetBit sets a bit at the given index to true.
func (ba *bitArray) SetBit(k uint64) error {
	if k >= ba.Capacity() {
		return OutOfRangeError(k)
	}

	if !ba.anyset {
		ba.lowest = k
		ba.highest = k
		ba.anyset = true
	} else {
		if k < ba.lowest {
			ba.lowest = k
		} else if k > ba.highest {
			ba.highest = k
		}
	}

	i, pos := getIndexAndRemainder(k)
	ba.blocks[i] = ba.blocks[i].insert(pos)
	return nil
}

// GetBit returns a bool indicating if the value at the given
// index has been set.
func (ba *bitArray) GetBit(k uint64) (bool, error) {
	if k >= ba.Capacity() {
		return false, OutOfRangeError(k)
	}

	i, pos := getIndexAndRemainder(k)
	result := ba.blocks[i]&block(1<<pos) != 0
	return result, nil
}

// GetSetBits gets the position of bits set in the array.
func (ba *bitArray) GetSetBits(from uint64, buffer []uint64) []uint64 {
	fromBlockIndex, fromOffset := getIndexAndRemainder(from)
	return getSetBitsInBlocks(
		fromBlockIndex,
		fromOffset,
		ba.blocks[fromBlockIndex:],
		nil,
		buffer,
	)
}

// getSetBitsInBlocks fills a buffer with positions of set bits in the provided blocks. Optionally, indices may be
// provided for sparse/non-consecutive blocks.
func getSetBitsInBlocks(
	fromBlockIndex, fromOffset uint64,
	blocks []block,
	indices []uint64,
	buffer []uint64,
) []uint64 {
	bufferCapacity := cap(buffer)
	if bufferCapacity == 0 {
		return buffer[:0]
	}

	results := buffer[:bufferCapacity]
	resultSize := 0

	for i, block := range blocks {
		blockIndex := fromBlockIndex + uint64(i)
		if indices != nil {
			blockIndex = indices[i]
		}

		isFirstBlock := blockIndex == fromBlockIndex
		if isFirstBlock {
			block >>= fromOffset
		}

		for block != 0 {
			trailing := bits.TrailingZeros64(uint64(block))

			if isFirstBlock {
				results[resultSize] = uint64(trailing) + (blockIndex << 6) + fromOffset
			} else {
				results[resultSize] = uint64(trailing) + (blockIndex << 6)
			}
			resultSize++

			if resultSize == cap(results) {
				return results[:resultSize]
			}

			// Clear the bit we just added to the result, which is the last bit set in the block. Ex.:
			//  block                   01001100
			//  ^block                  10110011
			//  (^block) + 1            10110100
			//  block & (^block) + 1    00000100
			//  block ^ mask            01001000
			mask := block & ((^block) + 1)
			block = block ^ mask
		}
	}

	return results[:resultSize]
}

// ClearBit will unset a bit at the given index if it is set.
func (ba *bitArray) ClearBit(k uint64) error {
	if k >= ba.Capacity() {
		return OutOfRangeError(k)
	}

	if !ba.anyset { // nothing is set, might as well bail
		return nil
	}

	i, pos := getIndexAndRemainder(k)
	ba.blocks[i] &^= block(1 << pos)

	if k == ba.highest {
		ba.setHighest()
	} else if k == ba.lowest {
		ba.setLowest()
	}
	return nil
}

// Count returns the number of set bits in this array.
func (ba *bitArray) Count() int {
	count := 0
	for _, block := range ba.blocks {
		count += bits.OnesCount64(uint64(block))
	}
	return count
}

// Or will bitwise or two bit arrays and return a new bit array
// representing the result.
func (ba *bitArray) Or(other BitArray) BitArray {
	if dba, ok := other.(*bitArray); ok {
		return orDenseWithDenseBitArray(ba, dba)
	}

	return orSparseWithDenseBitArray(other.(*sparseBitArray), ba)
}

// And will bitwise and two bit arrays and return a new bit array
// representing the result.
func (ba *bitArray) And(other BitArray) BitArray {
	if dba, ok := other.(*bitArray); ok {
		return andDenseWithDenseBitArray(ba, dba)
	}

	return andSparseWithDenseBitArray(other.(*sparseBitArray), ba)
}

// Nand will return the result of doing a bitwise and not of the bit array
// with the other bit array on each block.
func (ba *bitArray) Nand(other BitArray) BitArray {
	if dba, ok := other.(*bitArray); ok {
		return nandDenseWithDenseBitArray(ba, dba)
	}

	return nandDenseWithSparseBitArray(ba, other.(*sparseBitArray))
}

// Reset clears out the bit array.
func (ba *bitArray) Reset() {
	for i := uint64(0); i < uint64(len(ba.blocks)); i++ {
		ba.blocks[i] &= block(0)
	}
	ba.anyset = false
}

// Equals returns a bool indicating if these two bit arrays are equal.
func (ba *bitArray) Equals(other BitArray) bool {
	if other.Capacity() == 0 && ba.highest > 0 {
		return false
	}

	if other.Capacity() == 0 && !ba.anyset {
		return true
	}

	var selfIndex uint64
	for iter := other.Blocks(); iter.Next(); {
		toIndex, otherBlock := iter.Value()
		if toIndex > selfIndex {
			for i := selfIndex; i < toIndex; i++ {
				if ba.blocks[i] > 0 {
					return false
				}
			}
		}

		selfIndex = toIndex
		if !ba.blocks[selfIndex].equals(otherBlock) {
			return false
		}
		selfIndex++
	}

	lastIndex, _ := getIndexAndRemainder(ba.highest)
	if lastIndex >= selfIndex {
		return false
	}

	return true
}

// Intersects returns a bool indicating if the supplied bitarray intersects
// this bitarray.  This will check for intersection up to the length of the supplied
// bitarray.  If the supplied bitarray is longer than this bitarray, this
// function returns false.
func (ba *bitArray) Intersects(other BitArray) bool {
	if other.Capacity() > ba.Capacity() {
		return false
	}

	if sba, ok := other.(*sparseBitArray); ok {
		return ba.intersectsSparseBitArray(sba)
	}

	return ba.intersectsDenseBitArray(other.(*bitArray))
}

// Blocks will return an iterator over this bit array.
func (ba *bitArray) Blocks() Iterator {
	return newBitArrayIterator(ba)
}

func (ba *bitArray) IsEmpty() bool {
	return !ba.anyset
}

// complement flips all bits in this array.
func (ba *bitArray) complement() {
	for i := uint64(0); i < uint64(len(ba.blocks)); i++ {
		ba.blocks[i] = ^ba.blocks[i]
	}

	ba.setLowest()
	if ba.anyset {
		ba.setHighest()
	}
}

func (ba *bitArray) intersectsSparseBitArray(other *sparseBitArray) bool {
	for i, index := range other.indices {
		if !ba.blocks[index].intersects(other.blocks[i]) {
			return false
		}
	}

	return true
}

func (ba *bitArray) intersectsDenseBitArray(other *bitArray) bool {
	for i, block := range other.blocks {
		if !ba.blocks[i].intersects(block) {
			return false
		}
	}

	return true
}

func (ba *bitArray) copy() BitArray {
	blocks := make(blocks, len(ba.blocks))
	copy(blocks, ba.blocks)
	return &bitArray{
		blocks:  blocks,
		lowest:  ba.lowest,
		highest: ba.highest,
		anyset:  ba.anyset,
	}
}

// newBitArray returns a new dense BitArray at the specified size. This is a
// separate private constructor so unit tests don't have to constantly cast the
// BitArray interface to the concrete type.
func newBitArray(size uint64, args ...bool) *bitArray {
	i, r := getIndexAndRemainder(size)
	if r > 0 {
		i++
	}

	ba := &bitArray{
		blocks: make([]block, i),
		anyset: false,
	}

	if len(args) > 0 && args[0] == true {
		for i := uint64(0); i < uint64(len(ba.blocks)); i++ {
			ba.blocks[i] = maximumBlock
		}

		ba.lowest = 0
		ba.highest = i*s - 1
		ba.anyset = true
	}

	return ba
}

// NewBitArray returns a new BitArray at the specified size.  The
// optional arg denotes whether this bitarray should be set to the
// bitwise complement of the empty array, ie. sets all bits.
func NewBitArray(size uint64, args ...bool) BitArray {
	return newBitArray(size, args...)
}

// #endregion

// #region bitmap

// Bitmap32 tracks 32 bool values within a uint32
type Bitmap32 uint32

// SetBit returns a Bitmap32 with the bit at the given position set to 1
func (b Bitmap32) SetBit(pos uint) Bitmap32 {
	return b | (1 << pos)
}

// ClearBit returns a Bitmap32 with the bit at the given position set to 0
func (b Bitmap32) ClearBit(pos uint) Bitmap32 {
	return b & ^(1 << pos)
}

// GetBit returns true if the bit at the given position in the Bitmap32 is 1
func (b Bitmap32) GetBit(pos uint) bool {
	return (b & (1 << pos)) != 0
}

// PopCount returns the amount of bits set to 1 in the Bitmap32
func (b Bitmap32) PopCount() int {
	// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
	b -= (b >> 1) & 0x55555555
	b = (b>>2)&0x33333333 + b&0x33333333
	b += b >> 4
	b &= 0x0f0f0f0f
	b *= 0x01010101
	return int(byte(b >> 24))
}

// Bitmap64 tracks 64 bool values within a uint64
type Bitmap64 uint64

// SetBit returns a Bitmap64 with the bit at the given position set to 1
func (b Bitmap64) SetBit(pos uint) Bitmap64 {
	return b | (1 << pos)
}

// ClearBit returns a Bitmap64 with the bit at the given position set to 0
func (b Bitmap64) ClearBit(pos uint) Bitmap64 {
	return b & ^(1 << pos)
}

// GetBit returns true if the bit at the given position in the Bitmap64 is 1
func (b Bitmap64) GetBit(pos uint) bool {
	return (b & (1 << pos)) != 0
}

// PopCount returns the amount of bits set to 1 in the Bitmap64
func (b Bitmap64) PopCount() int {
	// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
	b -= (b >> 1) & 0x5555555555555555
	b = (b>>2)&0x3333333333333333 + b&0x3333333333333333
	b += b >> 4
	b &= 0x0f0f0f0f0f0f0f0f
	b *= 0x0101010101010101
	return int(byte(b >> 56))
}

// #endregion

// #region block

// block defines how we split apart the bit array. This also determines the size
// of s. This can be changed to any unsigned integer type: uint8, uint16,
// uint32, and so on.
type block uint64

// s denotes the size of any element in the block array.
// For a block of uint64, s will be equal to 64
// For a block of uint32, s will be equal to 32
// and so on...
const s = uint64(unsafe.Sizeof(block(0)) * 8)

// maximumBlock represents a block of all 1s and is used in the constructors.
const maximumBlock = block(0) | ^block(0)

func (b block) toNums(offset uint64, nums *[]uint64) {
	for i := uint64(0); i < s; i++ {
		if b&block(1<<i) > 0 {
			*nums = append(*nums, i+offset)
		}
	}
}

func (b block) findLeftPosition() uint64 {
	for i := s - 1; i < s; i-- {
		test := block(1 << i)
		if b&test == test {
			return i
		}
	}

	return s
}

func (b block) findRightPosition() uint64 {
	for i := uint64(0); i < s; i++ {
		test := block(1 << i)
		if b&test == test {
			return i
		}
	}

	return s
}

func (b block) insert(position uint64) block {
	return b | block(1<<position)
}

func (b block) remove(position uint64) block {
	return b & ^block(1<<position)
}

func (b block) or(other block) block {
	return b | other
}

func (b block) and(other block) block {
	return b & other
}

func (b block) nand(other block) block {
	return b &^ other
}

func (b block) get(position uint64) bool {
	return b&block(1<<position) != 0
}

func (b block) equals(other block) bool {
	return b == other
}

func (b block) intersects(other block) bool {
	return b&other == other
}

func (b block) String() string {
	return fmt.Sprintf(fmt.Sprintf("%%0%db", s), uint64(b))
}

// #endregion

// #region and

func andSparseWithSparseBitArray(sba, other *sparseBitArray) BitArray {
	max := maxInt64(int64(len(sba.indices)), int64(len(other.indices)))
	indices := make(uintSlice, 0, max)
	blocks := make(blocks, 0, max)

	selfIndex := 0
	otherIndex := 0
	var resultBlock block

	// move through the array and compare the blocks if they happen to
	// intersect
	for {
		if selfIndex == len(sba.indices) || otherIndex == len(other.indices) {
			// One of the arrays has been exhausted. We don't need
			// to compare anything else for a bitwise and; the
			// operation is complete.
			break
		}

		selfValue := sba.indices[selfIndex]
		otherValue := other.indices[otherIndex]

		switch {
		case otherValue < selfValue:
			// The `sba` bitarray has a block with a index position
			// greater than us. We want to compare with that block
			// if possible, so move our `other` index closer to that
			// block's index.
			otherIndex++

		case otherValue > selfValue:
			// This is the exact logical inverse of the above case.
			selfIndex++

		default:
			// Here, our indices match for both `sba` and `other`.
			// Time to do the bitwise AND operation and add a block
			// to our result list if the block has values in it.
			resultBlock = sba.blocks[selfIndex].and(other.blocks[otherIndex])
			if resultBlock > 0 {
				indices = append(indices, selfValue)
				blocks = append(blocks, resultBlock)
			}
			selfIndex++
			otherIndex++
		}
	}

	return &sparseBitArray{
		indices: indices,
		blocks:  blocks,
	}
}

func andSparseWithDenseBitArray(sba *sparseBitArray, other *bitArray) BitArray {
	if other.IsEmpty() {
		return newSparseBitArray()
	}

	// Use a duplicate of the sparse array to store the results of the
	// bitwise and. More memory-efficient than allocating a new dense bit
	// array.
	//
	// NOTE: this could be faster if we didn't copy the values as well
	// (since they are overwritten), but I don't want this method to know
	// too much about the internals of sparseBitArray. The performance hit
	// should be minor anyway.
	ba := sba.copy()

	// Run through the sparse array and attempt comparisons wherever
	// possible against the dense bit array.
	for selfIndex, selfValue := range ba.indices {

		if selfValue >= uint64(len(other.blocks)) {
			// The dense bit array has been exhausted. This is the
			// annoying case because we have to trim the sparse
			// array to the size of the dense array.
			ba.blocks = ba.blocks[:selfIndex-1]
			ba.indices = ba.indices[:selfIndex-1]

			// once this is done, there are no more comparisons.
			// We're ready to return
			break
		}
		ba.blocks[selfIndex] = ba.blocks[selfIndex].and(
			other.blocks[selfValue])

	}

	// Ensure any zero'd blocks in the resulting sparse
	// array are deleted
	for i := 0; i < len(ba.blocks); i++ {
		if ba.blocks[i] == 0 {
			ba.blocks.deleteAtIndex(int64(i))
			ba.indices.deleteAtIndex(int64(i))
			i--
		}
	}

	return ba
}

func andDenseWithDenseBitArray(dba, other *bitArray) BitArray {
	min := minUint64(uint64(len(dba.blocks)), uint64(len(other.blocks)))

	ba := newBitArray(min * s)

	for i := uint64(0); i < min; i++ {
		ba.blocks[i] = dba.blocks[i].and(other.blocks[i])
	}

	ba.setLowest()
	ba.setHighest()

	return ba
}

// #endregion

// #region nadn

func nandSparseWithSparseBitArray(sba, other *sparseBitArray) BitArray {
	// nand is an operation on the incoming array only, so the size will never
	// be more than the incoming array, regardless of the size of the other
	max := len(sba.indices)
	indices := make(uintSlice, 0, max)
	blocks := make(blocks, 0, max)

	selfIndex := 0
	otherIndex := 0
	var resultBlock block

	// move through the array and compare the blocks if they happen to
	// intersect
	for {
		if selfIndex == len(sba.indices) {
			// The bitarray being operated on is exhausted, so just return
			break
		} else if otherIndex == len(other.indices) {
			// The other array is exhausted. In this case, we assume that we
			// are calling nand on empty bit arrays, which is the same as just
			// copying the value in the sba array
			indices = append(indices, sba.indices[selfIndex])
			blocks = append(blocks, sba.blocks[selfIndex])
			selfIndex++
			continue
		}

		selfValue := sba.indices[selfIndex]
		otherValue := other.indices[otherIndex]

		switch {
		case otherValue < selfValue:
			// The `sba` bitarray has a block with a index position
			// greater than us. We want to compare with that block
			// if possible, so move our `other` index closer to that
			// block's index.
			otherIndex++

		case otherValue > selfValue:
			// Here, the sba array has blocks that the other array doesn't
			// have. In this case, we just copy exactly the sba array values
			indices = append(indices, selfValue)
			blocks = append(blocks, sba.blocks[selfIndex])

			// This is the exact logical inverse of the above case.
			selfIndex++

		default:
			// Here, our indices match for both `sba` and `other`.
			// Time to do the bitwise AND operation and add a block
			// to our result list if the block has values in it.
			resultBlock = sba.blocks[selfIndex].nand(other.blocks[otherIndex])
			if resultBlock > 0 {
				indices = append(indices, selfValue)
				blocks = append(blocks, resultBlock)
			}
			selfIndex++
			otherIndex++
		}
	}

	return &sparseBitArray{
		indices: indices,
		blocks:  blocks,
	}
}

func nandSparseWithDenseBitArray(sba *sparseBitArray, other *bitArray) BitArray {
	// Since nand is non-commutative, the resulting array should be sparse,
	// and the same length or less than the sparse array
	indices := make(uintSlice, 0, len(sba.indices))
	blocks := make(blocks, 0, len(sba.indices))

	var resultBlock block

	// Loop through the sparse array and match it with the dense array.
	for selfIndex, selfValue := range sba.indices {
		if selfValue >= uint64(len(other.blocks)) {
			// Since the dense array is exhausted, just copy over the data
			// from the sparse array
			resultBlock = sba.blocks[selfIndex]
			indices = append(indices, selfValue)
			blocks = append(blocks, resultBlock)
			continue
		}

		resultBlock = sba.blocks[selfIndex].nand(other.blocks[selfValue])
		if resultBlock > 0 {
			indices = append(indices, selfValue)
			blocks = append(blocks, resultBlock)
		}
	}

	return &sparseBitArray{
		indices: indices,
		blocks:  blocks,
	}
}

func nandDenseWithSparseBitArray(sba *bitArray, other *sparseBitArray) BitArray {
	// Since nand is non-commutative, the resulting array should be dense,
	// and the same length or less than the dense array
	tmp := sba.copy()
	ret := tmp.(*bitArray)

	// Loop through the other array and match it with the sba array.
	for otherIndex, otherValue := range other.indices {
		if otherValue >= uint64(len(ret.blocks)) {
			break
		}

		ret.blocks[otherValue] = sba.blocks[otherValue].nand(other.blocks[otherIndex])
	}

	ret.setLowest()
	ret.setHighest()

	return ret
}

func nandDenseWithDenseBitArray(dba, other *bitArray) BitArray {
	min := uint64(len(dba.blocks))

	ba := newBitArray(min * s)

	for i := uint64(0); i < min; i++ {
		ba.blocks[i] = dba.blocks[i].nand(other.blocks[i])
	}

	ba.setLowest()
	ba.setHighest()

	return ba
}

// #endregion

// #region or
func orSparseWithSparseBitArray(sba *sparseBitArray,
	other *sparseBitArray) BitArray {

	if len(other.indices) == 0 {
		return sba.copy()
	}

	if len(sba.indices) == 0 {
		return other.copy()
	}

	max := maxInt64(int64(len(sba.indices)), int64(len(other.indices)))
	indices := make(uintSlice, 0, max)
	blocks := make(blocks, 0, max)

	selfIndex := 0
	otherIndex := 0
	for {
		// last comparison was a real or, we are both exhausted now
		if selfIndex == len(sba.indices) && otherIndex == len(other.indices) {
			break
		} else if selfIndex == len(sba.indices) {
			indices = append(indices, other.indices[otherIndex:]...)
			blocks = append(blocks, other.blocks[otherIndex:]...)
			break
		} else if otherIndex == len(other.indices) {
			indices = append(indices, sba.indices[selfIndex:]...)
			blocks = append(blocks, sba.blocks[selfIndex:]...)
			break
		}

		selfValue := sba.indices[selfIndex]
		otherValue := other.indices[otherIndex]

		switch diff := int(otherValue) - int(selfValue); {
		case diff > 0:
			indices = append(indices, selfValue)
			blocks = append(blocks, sba.blocks[selfIndex])
			selfIndex++
		case diff < 0:
			indices = append(indices, otherValue)
			blocks = append(blocks, other.blocks[otherIndex])
			otherIndex++
		default:
			indices = append(indices, otherValue)
			blocks = append(blocks, sba.blocks[selfIndex].or(other.blocks[otherIndex]))
			selfIndex++
			otherIndex++
		}
	}

	return &sparseBitArray{
		indices: indices,
		blocks:  blocks,
	}
}

func orSparseWithDenseBitArray(sba *sparseBitArray, other *bitArray) BitArray {
	if other.Capacity() == 0 || !other.anyset {
		return sba.copy()
	}

	if sba.Capacity() == 0 {
		return other.copy()
	}

	max := maxUint64(uint64(sba.Capacity()), uint64(other.Capacity()))

	ba := newBitArray(max * s)
	selfIndex := 0
	otherIndex := 0
	for {
		if selfIndex == len(sba.indices) && otherIndex == len(other.blocks) {
			break
		} else if selfIndex == len(sba.indices) {
			copy(ba.blocks[otherIndex:], other.blocks[otherIndex:])
			break
		} else if otherIndex == len(other.blocks) {
			for i, value := range sba.indices[selfIndex:] {
				ba.blocks[value] = sba.blocks[i+selfIndex]
			}
			break
		}

		selfValue := sba.indices[selfIndex]
		if selfValue == uint64(otherIndex) {
			ba.blocks[otherIndex] = sba.blocks[selfIndex].or(other.blocks[otherIndex])
			selfIndex++
			otherIndex++
			continue
		}

		ba.blocks[otherIndex] = other.blocks[otherIndex]
		otherIndex++
	}

	ba.setHighest()
	ba.setLowest()

	return ba
}

func orDenseWithDenseBitArray(dba *bitArray, other *bitArray) BitArray {
	if dba.Capacity() == 0 || !dba.anyset {
		return other.copy()
	}

	if other.Capacity() == 0 || !other.anyset {
		return dba.copy()
	}

	max := maxUint64(uint64(len(dba.blocks)), uint64(len(other.blocks)))

	ba := newBitArray(max * s)

	for i := uint64(0); i < max; i++ {
		if i == uint64(len(dba.blocks)) {
			copy(ba.blocks[i:], other.blocks[i:])
			break
		}

		if i == uint64(len(other.blocks)) {
			copy(ba.blocks[i:], dba.blocks[i:])
			break
		}

		ba.blocks[i] = dba.blocks[i].or(other.blocks[i])
	}

	ba.setLowest()
	ba.setHighest()

	return ba
}

// #endregion

// #region sparse_bitarray
// uintSlice is an alias for a slice of ints.  Len, Swap, and Less
// are exported to fulfill an interface needed for the search
// function in the sort library.
type uintSlice []uint64

// Len returns the length of the slice.
func (u uintSlice) Len() int64 {
	return int64(len(u))
}

// Swap swaps values in this slice at the positions given.
func (u uintSlice) Swap(i, j int64) {
	u[i], u[j] = u[j], u[i]
}

// Less returns a bool indicating if the value at position i is
// less than position j.
func (u uintSlice) Less(i, j int64) bool {
	return u[i] < u[j]
}

func (u uintSlice) search(x uint64) int64 {
	return int64(sort.Search(len(u), func(i int) bool { return uint64(u[i]) >= x }))
}

func (u *uintSlice) insert(x uint64) (int64, bool) {
	i := u.search(x)

	if i == int64(len(*u)) {
		*u = append(*u, x)
		return i, true
	}

	if (*u)[i] == x {
		return i, false
	}

	*u = append(*u, 0)
	copy((*u)[i+1:], (*u)[i:])
	(*u)[i] = x
	return i, true
}

func (u *uintSlice) deleteAtIndex(i int64) {
	copy((*u)[i:], (*u)[i+1:])
	(*u)[len(*u)-1] = 0
	*u = (*u)[:len(*u)-1]
}

func (u uintSlice) get(x uint64) int64 {
	i := u.search(x)
	if i == int64(len(u)) {
		return -1
	}

	if u[i] == x {
		return i
	}

	return -1
}

type blocks []block

func (b *blocks) insert(index int64) {
	if index == int64(len(*b)) {
		*b = append(*b, block(0))
		return
	}

	*b = append(*b, block(0))
	copy((*b)[index+1:], (*b)[index:])
	(*b)[index] = block(0)
}

func (b *blocks) deleteAtIndex(i int64) {
	copy((*b)[i:], (*b)[i+1:])
	(*b)[len(*b)-1] = block(0)
	*b = (*b)[:len(*b)-1]
}

type sparseBitArray struct {
	blocks  blocks
	indices uintSlice
}

// SetBit sets the bit at the given position.
func (sba *sparseBitArray) SetBit(k uint64) error {
	index, position := getIndexAndRemainder(k)
	i, inserted := sba.indices.insert(index)

	if inserted {
		sba.blocks.insert(i)
	}
	sba.blocks[i] = sba.blocks[i].insert(position)
	return nil
}

// GetBit gets the bit at the given position.
func (sba *sparseBitArray) GetBit(k uint64) (bool, error) {
	index, position := getIndexAndRemainder(k)
	i := sba.indices.get(index)
	if i == -1 {
		return false, nil
	}

	return sba.blocks[i].get(position), nil
}

// GetSetBits gets the position of bits set in the array.
func (sba *sparseBitArray) GetSetBits(from uint64, buffer []uint64) []uint64 {
	fromBlockIndex, fromOffset := getIndexAndRemainder(from)

	fromBlockLocation := sba.indices.search(fromBlockIndex)
	if int(fromBlockLocation) == len(sba.indices) {
		return buffer[:0]
	}

	return getSetBitsInBlocks(
		fromBlockIndex,
		fromOffset,
		sba.blocks[fromBlockLocation:],
		sba.indices[fromBlockLocation:],
		buffer,
	)
}

// ToNums converts this sparse bitarray to a list of numbers contained
// within it.
func (sba *sparseBitArray) ToNums() []uint64 {
	if len(sba.indices) == 0 {
		return nil
	}

	diff := uint64(len(sba.indices)) * s
	nums := make([]uint64, 0, diff/4)

	for i, offset := range sba.indices {
		sba.blocks[i].toNums(offset*s, &nums)
	}

	return nums
}

// ClearBit clears the bit at the given position.
func (sba *sparseBitArray) ClearBit(k uint64) error {
	index, position := getIndexAndRemainder(k)
	i := sba.indices.get(index)
	if i == -1 {
		return nil
	}

	sba.blocks[i] = sba.blocks[i].remove(position)
	if sba.blocks[i] == 0 {
		sba.blocks.deleteAtIndex(i)
		sba.indices.deleteAtIndex(i)
	}

	return nil
}

// Reset erases all values from this bitarray.
func (sba *sparseBitArray) Reset() {
	sba.blocks = sba.blocks[:0]
	sba.indices = sba.indices[:0]
}

// Blocks returns an iterator to iterator of this bitarray's blocks.
func (sba *sparseBitArray) Blocks() Iterator {
	return newCompressedBitArrayIterator(sba)
}

// Capacity returns the value of the highest possible *seen* value
// in this sparse bitarray.
func (sba *sparseBitArray) Capacity() uint64 {
	if len(sba.indices) == 0 {
		return 0
	}

	return (sba.indices[len(sba.indices)-1] + 1) * s
}

// Equals returns a bool indicating if the provided bit array
// equals this bitarray.
func (sba *sparseBitArray) Equals(other BitArray) bool {
	if other.Capacity() == 0 && sba.Capacity() > 0 {
		return false
	}

	var selfIndex uint64
	for iter := other.Blocks(); iter.Next(); {
		otherIndex, otherBlock := iter.Value()
		if len(sba.indices) == 0 {
			if otherBlock > 0 {
				return false
			}

			continue
		}

		if selfIndex >= uint64(len(sba.indices)) {
			return false
		}

		if otherIndex < sba.indices[selfIndex] {
			if otherBlock > 0 {
				return false
			}
			continue
		}

		if otherIndex > sba.indices[selfIndex] {
			return false
		}

		if !sba.blocks[selfIndex].equals(otherBlock) {
			return false
		}

		selfIndex++
	}

	return true
}

// Count returns the number of set bits in this array.
func (sba *sparseBitArray) Count() int {
	count := 0
	for _, block := range sba.blocks {
		count += bits.OnesCount64(uint64(block))
	}
	return count
}

// Or will perform a bitwise or operation with the provided bitarray and
// return a new result bitarray.
func (sba *sparseBitArray) Or(other BitArray) BitArray {
	if ba, ok := other.(*sparseBitArray); ok {
		return orSparseWithSparseBitArray(sba, ba)
	}

	return orSparseWithDenseBitArray(sba, other.(*bitArray))
}

// And will perform a bitwise and operation with the provided bitarray and
// return a new result bitarray.
func (sba *sparseBitArray) And(other BitArray) BitArray {
	if ba, ok := other.(*sparseBitArray); ok {
		return andSparseWithSparseBitArray(sba, ba)
	}

	return andSparseWithDenseBitArray(sba, other.(*bitArray))
}

// Nand will return the result of doing a bitwise and not of the bit array
// with the other bit array on each block.
func (sba *sparseBitArray) Nand(other BitArray) BitArray {
	if ba, ok := other.(*sparseBitArray); ok {
		return nandSparseWithSparseBitArray(sba, ba)
	}

	return nandSparseWithDenseBitArray(sba, other.(*bitArray))
}

func (sba *sparseBitArray) IsEmpty() bool {
	// This works because the and, nand and delete functions only
	// keep values that have a non-zero block.
	return len(sba.indices) == 0
}

func (sba *sparseBitArray) copy() *sparseBitArray {
	blocks := make(blocks, len(sba.blocks))
	copy(blocks, sba.blocks)
	indices := make(uintSlice, len(sba.indices))
	copy(indices, sba.indices)
	return &sparseBitArray{
		blocks:  blocks,
		indices: indices,
	}
}

// Intersects returns a bool indicating if the provided bit array
// intersects with this bitarray.
func (sba *sparseBitArray) Intersects(other BitArray) bool {
	if other.Capacity() == 0 {
		return true
	}

	var selfIndex int64
	for iter := other.Blocks(); iter.Next(); {
		otherI, otherBlock := iter.Value()
		if len(sba.indices) == 0 {
			if otherBlock > 0 {
				return false
			}
			continue
		}
		// here we grab where the block should live in ourselves
		i := uintSlice(sba.indices[selfIndex:]).search(otherI)
		// this is a block we don't have, doesn't intersect
		if i == int64(len(sba.indices)) {
			return false
		}

		if sba.indices[i] != otherI {
			return false
		}

		if !sba.blocks[i].intersects(otherBlock) {
			return false
		}

		selfIndex = i
	}

	return true
}

func (sba *sparseBitArray) IntersectsBetween(other BitArray, start, stop uint64) bool {
	return true
}

func newSparseBitArray() *sparseBitArray {
	return &sparseBitArray{}
}

// NewSparseBitArray will create a bit array that consumes a great
// deal less memory at the expense of longer sets and gets.
func NewSparseBitArray() BitArray {
	return newSparseBitArray()
}

// #endregion

// #region encoding

// Marshal takes a dense or sparse bit array and serializes it to a
// byte slice.
func Marshal(ba BitArray) ([]byte, error) {
	if eba, ok := ba.(*bitArray); ok {
		return eba.Serialize()
	} else if sba, ok := ba.(*sparseBitArray); ok {
		return sba.Serialize()
	} else {
		return nil, errors.New("not a valid BitArray")
	}
}

// Unmarshal takes a byte slice, of the same format produced by Marshal,
// and returns a BitArray.
func Unmarshal(input []byte) (BitArray, error) {
	if len(input) == 0 {
		return nil, errors.New("no data in input")
	}
	if input[0] == 'B' {
		ret := newBitArray(0)
		err := ret.Deserialize(input)
		if err != nil {
			return nil, err
		}
		return ret, nil
	} else if input[0] == 'S' {
		ret := newSparseBitArray()
		err := ret.Deserialize(input)
		if err != nil {
			return nil, err
		}
		return ret, nil
	} else {
		return nil, errors.New("unrecognized encoding")
	}
}

// Serialize converts the sparseBitArray to a byte slice
func (ba *sparseBitArray) Serialize() ([]byte, error) {
	w := new(bytes.Buffer)

	var identifier uint8 = 'S'
	err := binary.Write(w, binary.LittleEndian, identifier)
	if err != nil {
		return nil, err
	}

	blocksLen := uint64(len(ba.blocks))
	indexLen := uint64(len(ba.indices))

	err = binary.Write(w, binary.LittleEndian, blocksLen)
	if err != nil {
		return nil, err
	}

	err = binary.Write(w, binary.LittleEndian, ba.blocks)
	if err != nil {
		return nil, err
	}

	err = binary.Write(w, binary.LittleEndian, indexLen)
	if err != nil {
		return nil, err
	}

	err = binary.Write(w, binary.LittleEndian, ba.indices)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// This function is a copy from the binary package, with some added error
// checking to avoid panics. The function will return the value, and the number
// of bytes read from the buffer. If the number of bytes is negative, then
// not enough bytes were passed in and the return value will be zero.
func Uint64FromBytes(b []byte) (uint64, int) {
	if len(b) < 8 {
		return 0, -1
	}

	val := uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	return val, 8
}

// Deserialize takes the incoming byte slice, and populates the sparseBitArray
// with data in the bytes. Note that this will overwrite any capacity
// specified when creating the sparseBitArray. Also note that if an error
// is returned, the sparseBitArray this is called on might be populated
// with partial data.
func (ret *sparseBitArray) Deserialize(incoming []byte) error {
	var intsize = uint64(s / 8)
	var curLoc = uint64(1) // Ignore the identifier byte

	var intsToRead uint64
	var bytesRead int
	intsToRead, bytesRead = Uint64FromBytes(incoming[curLoc : curLoc+intsize])
	if bytesRead < 0 {
		return errors.New("Invalid data for BitArray")
	}
	curLoc += intsize

	var nextblock uint64
	ret.blocks = make([]block, intsToRead)
	for i := uint64(0); i < intsToRead; i++ {
		nextblock, bytesRead = Uint64FromBytes(incoming[curLoc : curLoc+intsize])
		if bytesRead < 0 {
			return errors.New("Invalid data for BitArray")
		}
		ret.blocks[i] = block(nextblock)
		curLoc += intsize
	}

	intsToRead, bytesRead = Uint64FromBytes(incoming[curLoc : curLoc+intsize])
	if bytesRead < 0 {
		return errors.New("Invalid data for BitArray")
	}
	curLoc += intsize

	var nextuint uint64
	ret.indices = make(uintSlice, intsToRead)
	for i := uint64(0); i < intsToRead; i++ {
		nextuint, bytesRead = Uint64FromBytes(incoming[curLoc : curLoc+intsize])
		if bytesRead < 0 {
			return errors.New("Invalid data for BitArray")
		}
		ret.indices[i] = nextuint
		curLoc += intsize
	}
	return nil
}

// Serialize converts the bitArray to a byte slice.
func (ba *bitArray) Serialize() ([]byte, error) {
	w := new(bytes.Buffer)

	var identifier uint8 = 'B'
	err := binary.Write(w, binary.LittleEndian, identifier)
	if err != nil {
		return nil, err
	}

	err = binary.Write(w, binary.LittleEndian, ba.lowest)
	if err != nil {
		return nil, err
	}
	err = binary.Write(w, binary.LittleEndian, ba.highest)
	if err != nil {
		return nil, err
	}

	var encodedanyset uint8
	if ba.anyset {
		encodedanyset = 1
	} else {
		encodedanyset = 0
	}
	err = binary.Write(w, binary.LittleEndian, encodedanyset)
	if err != nil {
		return nil, err
	}

	err = binary.Write(w, binary.LittleEndian, ba.blocks)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// Deserialize takes the incoming byte slice, and populates the bitArray
// with data in the bytes. Note that this will overwrite any capacity
// specified when creating the bitArray. Also note that if an error is returned,
// the bitArray this is called on might be populated with partial data.
func (ret *bitArray) Deserialize(incoming []byte) error {
	r := bytes.NewReader(incoming[1:]) // Discard identifier

	err := binary.Read(r, binary.LittleEndian, &ret.lowest)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.LittleEndian, &ret.highest)
	if err != nil {
		return err
	}

	var encodedanyset uint8
	err = binary.Read(r, binary.LittleEndian, &encodedanyset)
	if err != nil {
		return err
	}

	// anyset defaults to false so we don't need an else statement
	if encodedanyset == 1 {
		ret.anyset = true
	}

	var nextblock block
	err = binary.Read(r, binary.LittleEndian, &nextblock)
	for err == nil {
		ret.blocks = append(ret.blocks, nextblock)
		err = binary.Read(r, binary.LittleEndian, &nextblock)
	}
	if err != io.EOF {
		return err
	}
	return nil
}

// #endregion

// #region error

// OutOfRangeError is an error caused by trying to access a bitarray past the end of its
// capacity.
type OutOfRangeError uint64

// Error returns a human readable description of the out-of-range error.
func (err OutOfRangeError) Error() string {
	return fmt.Sprintf(`Index %d is out of range.`, err)
}

// #endregion

// #region iterator

type sparseBitArrayIterator struct {
	index int64
	sba   *sparseBitArray
}

// Next increments the index and returns a bool indicating
// if any further items exist.
func (iter *sparseBitArrayIterator) Next() bool {
	iter.index++
	return iter.index < int64(len(iter.sba.indices))
}

// Value returns the index and block at the given index.
func (iter *sparseBitArrayIterator) Value() (uint64, block) {
	return iter.sba.indices[iter.index], iter.sba.blocks[iter.index]
}

func newCompressedBitArrayIterator(sba *sparseBitArray) *sparseBitArrayIterator {
	return &sparseBitArrayIterator{
		sba:   sba,
		index: -1,
	}
}

type bitArrayIterator struct {
	index     int64
	stopIndex uint64
	ba        *bitArray
}

// Next increments the index and returns a bool indicating if any further
// items exist.
func (iter *bitArrayIterator) Next() bool {
	iter.index++
	return uint64(iter.index) <= iter.stopIndex
}

// Value returns an index and the block at this index.
func (iter *bitArrayIterator) Value() (uint64, block) {
	return uint64(iter.index), iter.ba.blocks[iter.index]
}

func newBitArrayIterator(ba *bitArray) *bitArrayIterator {
	stop, _ := getIndexAndRemainder(ba.highest)
	start, _ := getIndexAndRemainder(ba.lowest)
	return &bitArrayIterator{
		ba:        ba,
		index:     int64(start) - 1,
		stopIndex: stop,
	}
}

// #endregion

// #region util

// maxInt64 returns the highest integer in the provided list of int64s
func maxInt64(ints ...int64) int64 {
	maxInt := ints[0]
	for i := 1; i < len(ints); i++ {
		if ints[i] > maxInt {
			maxInt = ints[i]
		}
	}

	return maxInt
}

// maxUint64 returns the highest integer in the provided list of uint64s
func maxUint64(ints ...uint64) uint64 {
	maxInt := ints[0]
	for i := 1; i < len(ints); i++ {
		if ints[i] > maxInt {
			maxInt = ints[i]
		}
	}

	return maxInt
}

// minUint64 returns the lowest integer in the provided list of int32s
func minUint64(ints ...uint64) uint64 {
	minInt := ints[0]
	for i := 1; i < len(ints); i++ {
		if ints[i] < minInt {
			minInt = ints[i]
		}
	}

	return minInt
}

// #endregion
