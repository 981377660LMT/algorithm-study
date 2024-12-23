package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/bits"
	"strconv"
)

// api:
// func New(length uint) (bset *BitSet)
// func MustNew(length uint) (bset *BitSet)
// func From(buf []uint64) *BitSet
// func FromWithLength(length uint, set []uint64) *BitSet
//
// func (b *BitSet) Len() uint
// func Cap() uint
//
// func (b *BitSet) Set(i uint) *BitSet
// func (b *BitSet) Clear(i uint) *BitSet
// func (b *BitSet) SetTo(i uint, value bool) *BitSet
// func (b *BitSet) Flip(i uint) *BitSet
// func (b *BitSet) Test(i uint) bool
//
// func (b *BitSet) Union(c *BitSet) *BitSet
// func (b *BitSet) Intersection(c *BitSet) *BitSet
// func (b *BitSet) Difference(c *BitSet) *BitSet
// func (b *BitSet) SymmetricDifference(c *BitSet) *BitSet
//
// func (b *BitSet) Count() uint
// func (b *BitSet) All() bool
// func (b *BitSet) Any() bool
// func (b *BitSet) None() bool
// func (b *BitSet) Equal(c *BitSet) bool
// func (b *BitSet) IsSuperSet(other *BitSet) bool
// func (b *BitSet) IsStrictSuperSet(other *BitSet) bool
//
// func (b *BitSet) NextSet(i uint) (uint, bool)
// func (b *BitSet) NextSetMany(i uint, buffer []uint) (uint, []uint)
//
// func (b *BitSet) Rank(index uint) uint
// func (b *BitSet) Select(index uint) uint
//
// func (b *BitSet) ShiftLeft(bits uint)
// func (b *BitSet) ShiftRight(bits uint)
//
// func (b *BitSet) WriteTo(stream io.Writer) (int64, error)
// func (b *BitSet) ReadFrom(stream io.Reader) (int64, error)
// func (b *BitSet) MarshalBinary() ([]byte, error)
// func (b *BitSet) UnmarshalBinary(data []byte) error
// func (b BitSet) MarshalJSON() ([]byte, error)
// func (b *BitSet) UnmarshalJSON(data []byte) error
//
// func (b *BitSet) Count() uint

func main() {
	// 创建一个新的 BitSet，长度为 100 位
	bitset := New(100)
	fmt.Println("Created BitSet with length 100.")

	// 设置一些位
	bitset.Set(10).Set(20).Set(30)
	fmt.Println("After setting bits 10, 20, 30:", bitset.String())

	// 测试位
	fmt.Println("Test bit 10:", bitset.Test(10)) // true
	fmt.Println("Test bit 15:", bitset.Test(15)) // false

	// 清除位 20
	bitset.Clear(20)
	fmt.Println("After clearing bit 20:", bitset.String())

	// 翻转位 30 和 40
	bitset.Flip(30).Flip(40)
	fmt.Println("After flipping bits 30 and 40:", bitset.String())

	// 计数被设置的位
	fmt.Println("Number of set bits:", bitset.Count()) // 2

	// 创建另一个 BitSet 并设置一些位
	other := New(100)
	other.Set(40).Set(50).Set(60)
	fmt.Println("Other BitSet:", other.String())

	// 并集
	union := bitset.Union(other)
	fmt.Println("Union of bitset and other:", union.String())

	// 交集
	intersection := bitset.Intersection(other)
	fmt.Println("Intersection of bitset and other:", intersection.String())

	// 差集
	difference := bitset.Difference(other)
	fmt.Println("Difference (bitset ∖ other):", difference.String())

	// 对称差集
	symmetricDiff := bitset.SymmetricDifference(other)
	fmt.Println("Symmetric Difference:", symmetricDiff.String())

	// 序列化为 JSON
	jsonData, err := json.Marshal(bitset)
	if err != nil {
		log.Fatalf("Failed to marshal BitSet to JSON: %v", err)
	}
	fmt.Println("JSON Encoded BitSet:", string(jsonData))

	// 反序列化 JSON
	var decoded BitSet
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON to BitSet: %v", err)
	}
	fmt.Println("Decoded BitSet from JSON:", decoded.String())

	// 使用 NextSet 迭代所有被设置的位
	fmt.Println("Iterating set bits using NextSet:")
	for i, ok := decoded.NextSet(0); ok; i, ok = decoded.NextSet(i + 1) {
		fmt.Printf("Bit %d is set.\n", i)
	}

	// 创建一个补集 BitSet
	complement := decoded.Complement()
	fmt.Println("Complement of decoded BitSet:", complement.String())
}

// the wordSize of a bit set
const wordSize = 64

// the wordSize of a bit set in bytes
const wordBytes = wordSize / 8

// wordMask is wordSize-1, used for bit indexing in a word
const wordMask = wordSize - 1

// log2WordSize is lg(wordSize)
const log2WordSize = 6

// allBits has every bit set
const allBits uint64 = 0xffffffffffffffff

// default binary BigEndian
// 用于二进制编码/解码的字节序，默认为 大端序
var binaryOrder binary.ByteOrder = binary.BigEndian

// default json encoding base64.URLEncoding
// 用于 JSON 编码/解码的 Base64 编码方式，默认为 URL 编码
var base64Encoding = base64.URLEncoding

// Base64StdEncoding Marshal/Unmarshal BitSet with base64.StdEncoding(Default: base64.URLEncoding)
func Base64StdEncoding() { base64Encoding = base64.StdEncoding }

// LittleEndian sets Marshal/Unmarshal Binary as Little Endian (Default: binary.BigEndian)
func LittleEndian() { binaryOrder = binary.LittleEndian }

// BigEndian sets Marshal/Unmarshal Binary as Big Endian (Default: binary.BigEndian)
func BigEndian() { binaryOrder = binary.BigEndian }

// BinaryOrder returns the current binary order, see also LittleEndian()
// and BigEndian() to change the order.
func BinaryOrder() binary.ByteOrder { return binaryOrder }

// A BitSet is a set of bits. The zero value of a BitSet is an empty set of length 0.
type BitSet struct {
	length uint
	set    []uint64
}

// Error is used to distinguish errors (panics) generated in this package.
type Error string

// safeSet will fixup b.set to be non-nil and return the field value
func (b *BitSet) safeSet() []uint64 {
	if b.set == nil {
		b.set = make([]uint64, wordsNeeded(0))
	}
	return b.set
}

// SetBitsetFrom fills the bitset with an array of integers without creating a new BitSet instance
func (b *BitSet) SetBitsetFrom(buf []uint64) {
	b.length = uint(len(buf)) * 64
	b.set = buf
}

// From is a constructor used to create a BitSet from an array of words
func From(buf []uint64) *BitSet {
	return FromWithLength(uint(len(buf))*64, buf)
}

// FromWithLength constructs from an array of words and length in bits.
// This function is for advanced users, most users should prefer
// the From function.
// As a user of FromWithLength, you are responsible for ensuring
// that the length is correct: your slice should have length at
// least (length+63)/64 in 64-bit words.
func FromWithLength(length uint, set []uint64) *BitSet {
	if len(set) < wordsNeeded(length) {
		panic("BitSet.FromWithLength: slice is too short")
	}
	return &BitSet{length, set}
}

// Bytes returns the bitset as array of 64-bit words, giving direct access to the internal representation.
// It is not a copy, so changes to the returned slice will affect the bitset.
// It is meant for advanced users.
//
// Deprecated: Bytes is deprecated. Use [BitSet.Words] instead.
func (b *BitSet) Bytes() []uint64 {
	return b.set
}

// Words returns the bitset as array of 64-bit words, giving direct access to the internal representation.
// It is not a copy, so changes to the returned slice will affect the bitset.
// It is meant for advanced users.
func (b *BitSet) Words() []uint64 {
	return b.set
}

// wordsNeeded calculates the number of words needed for i bits
func wordsNeeded(i uint) int {
	if i > (Cap() - wordMask) {
		return int(Cap() >> log2WordSize)
	}
	return int((i + wordMask) >> log2WordSize)
}

// wordsNeededUnbound calculates the number of words needed for i bits, possibly exceeding the capacity.
// This function is useful if you know that the capacity cannot be exceeded (e.g., you have an existing BitSet).
func wordsNeededUnbound(i uint) int {
	return (int(i) + wordMask) >> log2WordSize
}

// wordsIndex calculates the index of words in a `uint64`
func wordsIndex(i uint) uint {
	return i & wordMask
}

// New creates a new BitSet with a hint that length bits will be required.
// The memory usage is at least length/8 bytes.
// In case of allocation failure, the function will return a BitSet with zero
// capacity.
func New(length uint) (bset *BitSet) {
	defer func() {
		if r := recover(); r != nil {
			bset = &BitSet{
				0,
				make([]uint64, 0),
			}
		}
	}()

	bset = &BitSet{
		length,
		make([]uint64, wordsNeeded(length)),
	}

	return bset
}

// MustNew creates a new BitSet with the given length bits.
// It panics if length exceeds the possible capacity or by a lack of memory.
func MustNew(length uint) (bset *BitSet) {
	if length >= Cap() {
		panic("You are exceeding the capacity")
	}

	return &BitSet{
		length,
		make([]uint64, wordsNeeded(length)), // may panic on lack of memory
	}
}

// Cap returns the total possible capacity, or number of bits
// that can be stored in the BitSet theoretically. Under 32-bit system,
// it is 4294967295 and under 64-bit system, it is 18446744073709551615.
// Note that this is further limited by the maximum allocation size in Go,
// and your available memory, as any Go data structure.
func Cap() uint {
	return ^uint(0)
}

// Len returns the number of bits in the BitSet.
// Note that it differ from Count function.
func (b *BitSet) Len() uint {
	return b.length
}

// extendSet adds additional words to incorporate new bits if needed
func (b *BitSet) extendSet(i uint) {
	if i >= Cap() {
		panic("You are exceeding the capacity")
	}
	nsize := wordsNeeded(i + 1)
	if b.set == nil {
		b.set = make([]uint64, nsize)
	} else if cap(b.set) >= nsize {
		b.set = b.set[:nsize] // fast resize
	} else if len(b.set) < nsize {
		newset := make([]uint64, nsize, 2*nsize) // increase capacity 2x
		copy(newset, b.set)
		b.set = newset
	}
	b.length = i + 1
}

// Test whether bit i is set.
func (b *BitSet) Test(i uint) bool {
	if i >= b.length {
		return false
	}
	return b.set[i>>log2WordSize]&(1<<wordsIndex(i)) != 0
}

// GetWord64AtBit retrieves bits i through i+63 as a single uint64 value
func (b *BitSet) GetWord64AtBit(i uint) uint64 {
	firstWordIndex := int(i >> log2WordSize)
	subWordIndex := wordsIndex(i)

	// The word that the index falls within, shifted so the index is at bit 0
	var firstWord, secondWord uint64
	if firstWordIndex < len(b.set) {
		firstWord = b.set[firstWordIndex] >> subWordIndex
	}

	// The next word, masked to only include the necessary bits and shifted to cover the
	// top of the word
	if (firstWordIndex + 1) < len(b.set) {
		secondWord = b.set[firstWordIndex+1] << uint64(wordSize-subWordIndex)
	}

	return firstWord | secondWord
}

// Set bit i to 1, the capacity of the bitset is automatically
// increased accordingly.
// Warning: using a very large value for 'i'
// may lead to a memory shortage and a panic: the caller is responsible
// for providing sensible parameters in line with their memory capacity.
// The memory usage is at least slightly over i/8 bytes.
func (b *BitSet) Set(i uint) *BitSet {
	if i >= b.length { // if we need more bits, make 'em
		b.extendSet(i)
	}
	b.set[i>>log2WordSize] |= 1 << wordsIndex(i)
	return b
}

// Clear bit i to 0. This never cause a memory allocation. It is always safe.
func (b *BitSet) Clear(i uint) *BitSet {
	if i >= b.length {
		return b
	}
	b.set[i>>log2WordSize] &^= 1 << wordsIndex(i)
	return b
}

// SetTo sets bit i to value.
// Warning: using a very large value for 'i'
// may lead to a memory shortage and a panic: the caller is responsible
// for providing sensible parameters in line with their memory capacity.
func (b *BitSet) SetTo(i uint, value bool) *BitSet {
	if value {
		return b.Set(i)
	}
	return b.Clear(i)
}

// Flip bit at i.
// Warning: using a very large value for 'i'
// may lead to a memory shortage and a panic: the caller is responsible
// for providing sensible parameters in line with their memory capacity.
func (b *BitSet) Flip(i uint) *BitSet {
	if i >= b.length {
		return b.Set(i)
	}
	b.set[i>>log2WordSize] ^= 1 << wordsIndex(i)
	return b
}

// FlipRange bit in [start, end).
// Warning: using a very large value for 'end'
// may lead to a memory shortage and a panic: the caller is responsible
// for providing sensible parameters in line with their memory capacity.
func (b *BitSet) FlipRange(start, end uint) *BitSet {
	if start >= end {
		return b
	}

	if end-1 >= b.length { // if we need more bits, make 'em
		b.extendSet(end - 1)
	}

	startWord := int(start >> log2WordSize)
	endWord := int(end >> log2WordSize)

	// b.set[startWord] ^= ^(^uint64(0) << wordsIndex(start))
	//  e.g:
	//  start = 71,
	//  startWord = 1
	//  wordsIndex(start) = 71 % 64 = 7
	//   (^uint64(0) << 7) = 0b111111....11110000000
	//
	//  mask = ^(^uint64(0) << 7) = 0b000000....00001111111
	//
	// flips the first 7 bits in b.set[1] and
	// in the range loop, the b.set[1] gets again flipped
	// so the two expressions flip results in a flip
	// in b.set[1] from [7,63]
	//
	// handle startWord special, get's reflipped in range loop
	b.set[startWord] ^= ^(^uint64(0) << wordsIndex(start))

	for idx := range b.set[startWord:endWord] {
		b.set[startWord+idx] = ^b.set[startWord+idx]
	}

	// handle endWord special
	//  e.g.
	// end = 135
	//  endWord = 2
	//
	//  wordsIndex(-7) = 57
	//  see the golang spec:
	//   "For unsigned integer values, the operations +, -, *, and << are computed
	//   modulo 2n, where n is the bit width of the unsigned integer's type."
	//
	//   mask = ^uint64(0) >> 57 = 0b00000....0001111111
	//
	// flips in b.set[2] from [0,7]
	//
	// is end at word boundary?
	if idx := wordsIndex(-end); idx != 0 {
		b.set[endWord] ^= ^uint64(0) >> wordsIndex(idx)
	}

	return b
}

// Shrink shrinks BitSet so that the provided value is the last possible
// set value. It clears all bits > the provided index and reduces the size
// and length of the set.
//
// Note that the parameter value is not the new length in bits: it is the
// maximal value that can be stored in the bitset after the function call.
// The new length in bits is the parameter value + 1. Thus it is not possible
// to use this function to set the length to 0, the minimal value of the length
// after this function call is 1.
//
// A new slice is allocated to store the new bits, so you may see an increase in
// memory usage until the GC runs. Normally this should not be a problem, but if you
// have an extremely large BitSet its important to understand that the old BitSet will
// remain in memory until the GC frees it.
// If you are memory constrained, this function may cause a panic.
func (b *BitSet) Shrink(lastbitindex uint) *BitSet {
	length := lastbitindex + 1
	idx := wordsNeeded(length)
	if idx > len(b.set) {
		return b
	}
	shrunk := make([]uint64, idx)
	copy(shrunk, b.set[:idx])
	b.set = shrunk
	b.length = length
	lastWordUsedBits := length % 64
	if lastWordUsedBits != 0 {
		b.set[idx-1] &= allBits >> uint64(64-wordsIndex(lastWordUsedBits))
	}
	return b
}

// Compact shrinks BitSet to so that we preserve all set bits, while minimizing
// memory usage. Compact calls Shrink.
// A new slice is allocated to store the new bits, so you may see an increase in
// memory usage until the GC runs. Normally this should not be a problem, but if you
// have an extremely large BitSet its important to understand that the old BitSet will
// remain in memory until the GC frees it.
// If you are memory constrained, this function may cause a panic.
func (b *BitSet) Compact() *BitSet {
	idx := len(b.set) - 1
	for ; idx >= 0 && b.set[idx] == 0; idx-- {
	}
	newlength := uint((idx + 1) << log2WordSize)
	if newlength >= b.length {
		return b // nothing to do
	}
	if newlength > 0 {
		return b.Shrink(newlength - 1)
	}
	// We preserve one word
	return b.Shrink(63)
}

// InsertAt takes an index which indicates where a bit should be
// inserted. Then it shifts all the bits in the set to the left by 1, starting
// from the given index position, and sets the index position to 0.
//
// Depending on the size of your BitSet, and where you are inserting the new entry,
// this method could be extremely slow and in some cases might cause the entire BitSet
// to be recopied.
func (b *BitSet) InsertAt(idx uint) *BitSet {
	insertAtElement := idx >> log2WordSize

	// if length of set is a multiple of wordSize we need to allocate more space first
	if b.isLenExactMultiple() {
		b.set = append(b.set, uint64(0))
	}

	var i uint
	for i = uint(len(b.set) - 1); i > insertAtElement; i-- {
		// all elements above the position where we want to insert can simply by shifted
		b.set[i] <<= 1

		// we take the most significant bit of the previous element and set it as
		// the least significant bit of the current element
		b.set[i] |= (b.set[i-1] & 0x8000000000000000) >> 63
	}

	// generate a mask to extract the data that we need to shift left
	// within the element where we insert a bit
	dataMask := uint64(1)<<uint64(wordsIndex(idx)) - 1

	// extract that data that we'll shift
	data := b.set[i] & (^dataMask)

	// set the positions of the data mask to 0 in the element where we insert
	b.set[i] &= dataMask

	// shift data mask to the left and insert its data to the slice element
	b.set[i] |= data << 1

	// add 1 to length of BitSet
	b.length++

	return b
}

// String creates a string representation of the BitSet. It is only intended for
// human-readable output and not for serialization.
func (b *BitSet) String() string {
	// follows code from https://github.com/RoaringBitmap/roaring
	var buffer bytes.Buffer
	start := []byte("{")
	buffer.Write(start)
	counter := 0
	i, e := b.NextSet(0)
	for e {
		counter = counter + 1
		// to avoid exhausting the memory
		if counter > 0x40000 {
			buffer.WriteString("...")
			break
		}
		buffer.WriteString(strconv.FormatInt(int64(i), 10))
		i, e = b.NextSet(i + 1)
		if e {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")
	return buffer.String()
}

// DeleteAt deletes the bit at the given index position from
// within the bitset
// All the bits residing on the left of the deleted bit get
// shifted right by 1
// The running time of this operation may potentially be
// relatively slow, O(length)
func (b *BitSet) DeleteAt(i uint) *BitSet {
	// the index of the slice element where we'll delete a bit
	deleteAtElement := i >> log2WordSize

	// generate a mask for the data that needs to be shifted right
	// within that slice element that gets modified
	dataMask := ^((uint64(1) << wordsIndex(i)) - 1)

	// extract the data that we'll shift right from the slice element
	data := b.set[deleteAtElement] & dataMask

	// set the masked area to 0 while leaving the rest as it is
	b.set[deleteAtElement] &= ^dataMask

	// shift the previously extracted data to the right and then
	// set it in the previously masked area
	b.set[deleteAtElement] |= (data >> 1) & dataMask

	// loop over all the consecutive slice elements to copy each
	// lowest bit into the highest position of the previous element,
	// then shift the entire content to the right by 1
	for i := int(deleteAtElement) + 1; i < len(b.set); i++ {
		b.set[i-1] |= (b.set[i] & 1) << 63
		b.set[i] >>= 1
	}

	b.length = b.length - 1

	return b
}

// AppendTo appends all set bits to buf and returns the (maybe extended) buf.
// In case of allocation failure, the function will panic.
//
// See also [BitSet.AsSlice] and [BitSet.NextSetMany].
func (b *BitSet) AppendTo(buf []uint) []uint {
	// In theory, we could overflow uint, but in practice, we will not.
	for idx, word := range b.set {
		for word != 0 {
			// In theory idx<<log2WordSize could overflow, but it will not overflow
			// in practice.
			buf = append(buf, uint(idx<<log2WordSize+bits.TrailingZeros64(word)))

			// clear the rightmost set bit
			word &= word - 1
		}
	}

	return buf
}

// AsSlice returns all set bits as slice.
// It panics if the capacity of buf is < b.Count()
//
// See also [BitSet.AppendTo] and [BitSet.NextSetMany].
func (b *BitSet) AsSlice(buf []uint) []uint {
	buf = buf[:cap(buf)] // len = cap

	size := 0
	for idx, word := range b.set {
		for ; word != 0; size++ {
			// panics if capacity of buf is exceeded.
			// In theory idx<<log2WordSize could overflow, but it will not overflow
			// in practice.
			buf[size] = uint(idx<<log2WordSize + bits.TrailingZeros64(word))

			// clear the rightmost set bit
			word &= word - 1
		}
	}

	buf = buf[:size]
	return buf
}

// NextSet returns the next bit set from the specified index,
// including possibly the current index
// along with an error code (true = valid, false = no set bit found)
// for i,e := v.NextSet(0); e; i,e = v.NextSet(i + 1) {...}
//
// Users concerned with performance may want to use NextSetMany to
// retrieve several values at once.
func (b *BitSet) NextSet(i uint) (uint, bool) {
	x := int(i >> log2WordSize)
	if x >= len(b.set) {
		return 0, false
	}

	// process first (partial) word
	word := b.set[x] >> wordsIndex(i)
	if word != 0 {
		return i + uint(bits.TrailingZeros64(word)), true
	}

	// process the following full words until next bit is set
	// x < len(b.set), no out-of-bounds panic in following slice expression
	x++
	for idx, word := range b.set[x:] {
		if word != 0 {
			return uint((x+idx)<<log2WordSize + bits.TrailingZeros64(word)), true
		}
	}

	return 0, false
}

// NextSetMany returns many next bit sets from the specified index,
// including possibly the current index and up to cap(buffer).
// If the returned slice has len zero, then no more set bits were found
//
//	buffer := make([]uint, 256) // this should be reused
//	j := uint(0)
//	j, buffer = bitmap.NextSetMany(j, buffer)
//	for ; len(buffer) > 0; j, buffer = bitmap.NextSetMany(j,buffer) {
//	 for k := range buffer {
//	  do something with buffer[k]
//	 }
//	 j += 1
//	}
//
// It is possible to retrieve all set bits as follow:
//
//	indices := make([]uint, bitmap.Count())
//	bitmap.NextSetMany(0, indices)
//
// It is also possible to retrieve all set bits with [BitSet.AppendTo]
// or [BitSet.AsSlice].
//
// However if Count() is large, it might be preferable to
// use several calls to NextSetMany for memory reasons.
func (b *BitSet) NextSetMany(i uint, buffer []uint) (uint, []uint) {
	// In theory, we could overflow uint, but in practice, we will not.
	capacity := cap(buffer)
	result := buffer[:capacity]

	x := int(i >> log2WordSize)
	if x >= len(b.set) || capacity == 0 {
		return 0, result[:0]
	}

	// process first (partial) word
	word := b.set[x] >> wordsIndex(i)

	size := 0
	for word != 0 {
		result[size] = i + uint(bits.TrailingZeros64(word))

		size++
		if size == capacity {
			return result[size-1], result[:size]
		}

		// clear the rightmost set bit
		word &= word - 1
	}

	// process the following full words
	// x < len(b.set), no out-of-bounds panic in following slice expression
	x++
	for idx, word := range b.set[x:] {
		for word != 0 {
			result[size] = uint((x+idx)<<log2WordSize + bits.TrailingZeros64(word))

			size++
			if size == capacity {
				return result[size-1], result[:size]
			}

			// clear the rightmost set bit
			word &= word - 1
		}
	}

	if size > 0 {
		return result[size-1], result[:size]
	}
	return 0, result[:0]
}

// NextClear returns the next clear bit from the specified index,
// including possibly the current index
// along with an error code (true = valid, false = no bit found i.e. all bits are set)
func (b *BitSet) NextClear(i uint) (uint, bool) {
	x := int(i >> log2WordSize)
	if x >= len(b.set) {
		return 0, false
	}

	// process first (maybe partial) word
	word := b.set[x]
	word = word >> wordsIndex(i)
	wordAll := allBits >> wordsIndex(i)

	index := i + uint(bits.TrailingZeros64(^word))
	if word != wordAll && index < b.length {
		return index, true
	}

	// process the following full words until next bit is cleared
	// x < len(b.set), no out-of-bounds panic in following slice expression
	x++
	for idx, word := range b.set[x:] {
		if word != allBits {
			index = uint((x+idx)*wordSize + bits.TrailingZeros64(^word))
			if index < b.length {
				return index, true
			}
		}
	}

	return 0, false
}

// PreviousSet returns the previous set bit from the specified index,
// including possibly the current index
// along with an error code (true = valid, false = no bit found i.e. all bits are clear)
func (b *BitSet) PreviousSet(i uint) (uint, bool) {
	x := int(i >> log2WordSize)
	if x >= len(b.set) {
		return 0, false
	}
	word := b.set[x]

	// Clear the bits above the index
	word = word & ((1 << (wordsIndex(i) + 1)) - 1)
	if word != 0 {
		return uint(x<<log2WordSize+bits.Len64(word)) - 1, true
	}

	for x--; x >= 0; x-- {
		word = b.set[x]
		if word != 0 {
			return uint(x<<log2WordSize+bits.Len64(word)) - 1, true
		}
	}
	return 0, false
}

// PreviousClear returns the previous clear bit from the specified index,
// including possibly the current index
// along with an error code (true = valid, false = no clear bit found i.e. all bits are set)
func (b *BitSet) PreviousClear(i uint) (uint, bool) {
	x := int(i >> log2WordSize)
	if x >= len(b.set) {
		return 0, false
	}
	word := b.set[x]

	// Flip all bits and find the highest one bit
	word = ^word

	// Clear the bits above the index
	word = word & ((1 << (wordsIndex(i) + 1)) - 1)

	if word != 0 {
		return uint(x<<log2WordSize+bits.Len64(word)) - 1, true
	}

	for x--; x >= 0; x-- {
		word = b.set[x]
		word = ^word
		if word != 0 {
			return uint(x<<log2WordSize+bits.Len64(word)) - 1, true
		}
	}
	return 0, false
}

// ClearAll clears the entire BitSet.
// It does not free the memory.
func (b *BitSet) ClearAll() *BitSet {
	if b != nil && b.set != nil {
		for i := range b.set {
			b.set[i] = 0
		}
	}
	return b
}

// SetAll sets the entire BitSet
func (b *BitSet) SetAll() *BitSet {
	if b != nil && b.set != nil {
		for i := range b.set {
			b.set[i] = allBits
		}

		b.cleanLastWord()
	}
	return b
}

// wordCount returns the number of words used in a bit set
func (b *BitSet) wordCount() int {
	return wordsNeededUnbound(b.length)
}

// Clone this BitSet, returning a new BitSet that has the same bits set.
// In case of allocation failure, the function will return an empty BitSet.
func (b *BitSet) Clone() *BitSet {
	c := New(b.length)
	if b.set != nil { // Clone should not modify current object
		copy(c.set, b.set)
	}
	return c
}

// Copy into a destination BitSet using the Go array copy semantics:
// the number of bits copied is the minimum of the number of bits in the current
// BitSet (Len()) and the destination Bitset.
// We return the number of bits copied in the destination BitSet.
func (b *BitSet) Copy(c *BitSet) (count uint) {
	if c == nil {
		return
	}
	if b.set != nil { // Copy should not modify current object
		copy(c.set, b.set)
	}
	count = c.length
	if b.length < c.length {
		count = b.length
	}
	// Cleaning the last word is needed to keep the invariant that other functions, such as Count, require
	// that any bits in the last word that would exceed the length of the bitmask are set to 0.
	c.cleanLastWord()
	return
}

// CopyFull copies into a destination BitSet such that the destination is
// identical to the source after the operation, allocating memory if necessary.
func (b *BitSet) CopyFull(c *BitSet) {
	if c == nil {
		return
	}
	c.length = b.length
	if len(b.set) == 0 {
		if c.set != nil {
			c.set = c.set[:0]
		}
	} else {
		if cap(c.set) < len(b.set) {
			c.set = make([]uint64, len(b.set))
		} else {
			c.set = c.set[:len(b.set)]
		}
		copy(c.set, b.set)
	}
}

// Count (number of set bits).
// Also known as "popcount" or "population count".
func (b *BitSet) Count() uint {
	if b != nil && b.set != nil {
		return uint(popcntSlice(b.set))
	}
	return 0
}

// Equal tests the equivalence of two BitSets.
// False if they are of different sizes, otherwise true
// only if all the same bits are set
func (b *BitSet) Equal(c *BitSet) bool {
	if c == nil || b == nil {
		return c == b
	}
	if b.length != c.length {
		return false
	}
	if b.length == 0 { // if they have both length == 0, then could have nil set
		return true
	}
	wn := b.wordCount()
	// bounds check elimination
	if wn <= 0 {
		return true
	}
	_ = b.set[wn-1]
	_ = c.set[wn-1]
	for p := 0; p < wn; p++ {
		if c.set[p] != b.set[p] {
			return false
		}
	}
	return true
}

func panicIfNull(b *BitSet) {
	if b == nil {
		panic(Error("BitSet must not be null"))
	}
}

// Difference of base set and other set
// This is the BitSet equivalent of &^ (and not)
func (b *BitSet) Difference(compare *BitSet) (result *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	result = b.Clone() // clone b (in case b is bigger than compare)
	l := compare.wordCount()
	if l > b.wordCount() {
		l = b.wordCount()
	}
	for i := 0; i < l; i++ {
		result.set[i] = b.set[i] &^ compare.set[i]
	}
	return
}

// DifferenceCardinality computes the cardinality of the difference
func (b *BitSet) DifferenceCardinality(compare *BitSet) uint {
	panicIfNull(b)
	panicIfNull(compare)
	l := compare.wordCount()
	if l > b.wordCount() {
		l = b.wordCount()
	}
	cnt := uint64(0)
	cnt += popcntMaskSlice(b.set[:l], compare.set[:l])
	cnt += popcntSlice(b.set[l:])
	return uint(cnt)
}

// InPlaceDifference computes the difference of base set and other set
// This is the BitSet equivalent of &^ (and not)
func (b *BitSet) InPlaceDifference(compare *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	l := compare.wordCount()
	if l > b.wordCount() {
		l = b.wordCount()
	}
	if l <= 0 {
		return
	}
	// bounds check elimination
	data, cmpData := b.set, compare.set
	_ = data[l-1]
	_ = cmpData[l-1]
	for i := 0; i < l; i++ {
		data[i] &^= cmpData[i]
	}
}

// Convenience function: return two bitsets ordered by
// increasing length. Note: neither can be nil
func sortByLength(a *BitSet, b *BitSet) (ap *BitSet, bp *BitSet) {
	if a.length <= b.length {
		ap, bp = a, b
	} else {
		ap, bp = b, a
	}
	return
}

// Intersection of base set and other set
// This is the BitSet equivalent of & (and)
// In case of allocation failure, the function will return an empty BitSet.
func (b *BitSet) Intersection(compare *BitSet) (result *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	result = New(b.length)
	for i, word := range b.set {
		result.set[i] = word & compare.set[i]
	}
	return
}

// IntersectionCardinality computes the cardinality of the intersection
func (b *BitSet) IntersectionCardinality(compare *BitSet) uint {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	cnt := popcntAndSlice(b.set, compare.set)
	return uint(cnt)
}

// InPlaceIntersection destructively computes the intersection of
// base set and the compare set.
// This is the BitSet equivalent of & (and)
func (b *BitSet) InPlaceIntersection(compare *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	l := compare.wordCount()
	if l > b.wordCount() {
		l = b.wordCount()
	}
	if l > 0 {
		// bounds check elimination
		data, cmpData := b.set, compare.set
		_ = data[l-1]
		_ = cmpData[l-1]

		for i := 0; i < l; i++ {
			data[i] &= cmpData[i]
		}
	}
	if l >= 0 {
		for i := l; i < len(b.set); i++ {
			b.set[i] = 0
		}
	}
	if compare.length > 0 {
		if compare.length-1 >= b.length {
			b.extendSet(compare.length - 1)
		}
	}
}

// Union of base set and other set
// This is the BitSet equivalent of | (or)
func (b *BitSet) Union(compare *BitSet) (result *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	result = compare.Clone()
	for i, word := range b.set {
		result.set[i] = word | compare.set[i]
	}
	return
}

// UnionCardinality computes the cardinality of the uniton of the base set
// and the compare set.
func (b *BitSet) UnionCardinality(compare *BitSet) uint {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	cnt := popcntOrSlice(b.set, compare.set)
	if len(compare.set) > len(b.set) {
		cnt += popcntSlice(compare.set[len(b.set):])
	}
	return uint(cnt)
}

// InPlaceUnion creates the destructive union of base set and compare set.
// This is the BitSet equivalent of | (or).
func (b *BitSet) InPlaceUnion(compare *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	l := compare.wordCount()
	if l > b.wordCount() {
		l = b.wordCount()
	}
	if compare.length > 0 && compare.length-1 >= b.length {
		b.extendSet(compare.length - 1)
	}
	if l > 0 {
		// bounds check elimination
		data, cmpData := b.set, compare.set
		_ = data[l-1]
		_ = cmpData[l-1]

		for i := 0; i < l; i++ {
			data[i] |= cmpData[i]
		}
	}
	if len(compare.set) > l {
		for i := l; i < len(compare.set); i++ {
			b.set[i] = compare.set[i]
		}
	}
}

// SymmetricDifference of base set and other set
// This is the BitSet equivalent of ^ (xor)
func (b *BitSet) SymmetricDifference(compare *BitSet) (result *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	// compare is bigger, so clone it
	result = compare.Clone()
	for i, word := range b.set {
		result.set[i] = word ^ compare.set[i]
	}
	return
}

// SymmetricDifferenceCardinality computes the cardinality of the symmetric difference
func (b *BitSet) SymmetricDifferenceCardinality(compare *BitSet) uint {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	cnt := popcntXorSlice(b.set, compare.set)
	if len(compare.set) > len(b.set) {
		cnt += popcntSlice(compare.set[len(b.set):])
	}
	return uint(cnt)
}

// InPlaceSymmetricDifference creates the destructive SymmetricDifference of base set and other set
// This is the BitSet equivalent of ^ (xor)
func (b *BitSet) InPlaceSymmetricDifference(compare *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	l := compare.wordCount()
	if l > b.wordCount() {
		l = b.wordCount()
	}
	if compare.length > 0 && compare.length-1 >= b.length {
		b.extendSet(compare.length - 1)
	}
	if l > 0 {
		// bounds check elimination
		data, cmpData := b.set, compare.set
		_ = data[l-1]
		_ = cmpData[l-1]
		for i := 0; i < l; i++ {
			data[i] ^= cmpData[i]
		}
	}
	if len(compare.set) > l {
		for i := l; i < len(compare.set); i++ {
			b.set[i] = compare.set[i]
		}
	}
}

// Is the length an exact multiple of word sizes?
func (b *BitSet) isLenExactMultiple() bool {
	return wordsIndex(b.length) == 0
}

// Clean last word by setting unused bits to 0
func (b *BitSet) cleanLastWord() {
	if !b.isLenExactMultiple() {
		b.set[len(b.set)-1] &= allBits >> (wordSize - wordsIndex(b.length))
	}
}

// Complement computes the (local) complement of a bitset (up to length bits)
// In case of allocation failure, the function will return an empty BitSet.
func (b *BitSet) Complement() (result *BitSet) {
	panicIfNull(b)
	result = New(b.length)
	for i, word := range b.set {
		result.set[i] = ^word
	}
	result.cleanLastWord()
	return
}

// All returns true if all bits are set, false otherwise. Returns true for
// empty sets.
func (b *BitSet) All() bool {
	panicIfNull(b)
	return b.Count() == b.length
}

// None returns true if no bit is set, false otherwise. Returns true for
// empty sets.
func (b *BitSet) None() bool {
	panicIfNull(b)
	if b != nil && b.set != nil {
		for _, word := range b.set {
			if word > 0 {
				return false
			}
		}
	}
	return true
}

// Any returns true if any bit is set, false otherwise
func (b *BitSet) Any() bool {
	panicIfNull(b)
	return !b.None()
}

// IsSuperSet returns true if this is a superset of the other set
func (b *BitSet) IsSuperSet(other *BitSet) bool {
	l := other.wordCount()
	if b.wordCount() < l {
		l = b.wordCount()
	}
	for i, word := range other.set[:l] {
		if b.set[i]&word != word {
			return false
		}
	}
	return popcntSlice(other.set[l:]) == 0
}

// IsStrictSuperSet returns true if this is a strict superset of the other set
func (b *BitSet) IsStrictSuperSet(other *BitSet) bool {
	return b.Count() > other.Count() && b.IsSuperSet(other)
}

// DumpAsBits dumps a bit set as a string of bits. Following the usual convention in Go,
// the least significant bits are printed last (index 0 is at the end of the string).
// This is useful for debugging and testing. It is not suitable for serialization.
func (b *BitSet) DumpAsBits() string {
	if b.set == nil {
		return "."
	}
	buffer := bytes.NewBufferString("")
	i := len(b.set) - 1
	for ; i >= 0; i-- {
		fmt.Fprintf(buffer, "%064b.", b.set[i])
	}
	return buffer.String()
}

// BinaryStorageSize returns the binary storage requirements (see WriteTo) in bytes.
func (b *BitSet) BinaryStorageSize() int {
	return wordBytes + wordBytes*b.wordCount()
}

func readUint64Array(reader io.Reader, data []uint64) error {
	length := len(data)
	bufferSize := 128
	buffer := make([]byte, bufferSize*wordBytes)
	for i := 0; i < length; i += bufferSize {
		end := i + bufferSize
		if end > length {
			end = length
			buffer = buffer[:wordBytes*(end-i)]
		}
		chunk := data[i:end]
		if _, err := io.ReadFull(reader, buffer); err != nil {
			return err
		}
		for i := range chunk {
			chunk[i] = uint64(binaryOrder.Uint64(buffer[8*i:]))
		}
	}
	return nil
}

func writeUint64Array(writer io.Writer, data []uint64) error {
	bufferSize := 128
	buffer := make([]byte, bufferSize*wordBytes)
	for i := 0; i < len(data); i += bufferSize {
		end := i + bufferSize
		if end > len(data) {
			end = len(data)
			buffer = buffer[:wordBytes*(end-i)]
		}
		chunk := data[i:end]
		for i, x := range chunk {
			binaryOrder.PutUint64(buffer[8*i:], x)
		}
		_, err := writer.Write(buffer)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteTo writes a BitSet to a stream. The format is:
// 1. uint64 length
// 2. []uint64 set
// The length is the number of bits in the BitSet.
//
// The set is a slice of uint64s containing between length and length + 63 bits.
// It is interpreted as a big-endian array of uint64s by default (see BinaryOrder())
// meaning that the first 8 bits are stored at byte index 7, the next 8 bits are stored
// at byte index 6... the bits 64 to 71 are stored at byte index 8, etc.
// If you change the binary order, you need to do so for both reading and writing.
// We recommend using the default binary order.
//
// Upon success, the number of bytes written is returned.
//
// Performance: if this function is used to write to a disk or network
// connection, it might be beneficial to wrap the stream in a bufio.Writer.
// E.g.,
//
//	      f, err := os.Create("myfile")
//		       w := bufio.NewWriter(f)
func (b *BitSet) WriteTo(stream io.Writer) (int64, error) {
	length := uint64(b.length)
	// Write length
	err := binary.Write(stream, binaryOrder, &length)
	if err != nil {
		// Upon failure, we do not guarantee that we
		// return the number of bytes written.
		return int64(0), err
	}
	err = writeUint64Array(stream, b.set[:b.wordCount()])
	if err != nil {
		// Upon failure, we do not guarantee that we
		// return the number of bytes written.
		return int64(wordBytes), err
	}
	return int64(b.BinaryStorageSize()), nil
}

// ReadFrom reads a BitSet from a stream written using WriteTo
// The format is:
// 1. uint64 length
// 2. []uint64 set
// See WriteTo for details.
// Upon success, the number of bytes read is returned.
// If the current BitSet is not large enough to hold the data,
// it is extended. In case of error, the BitSet is either
// left unchanged or made empty if the error occurs too late
// to preserve the content.
//
// Performance: if this function is used to read from a disk or network
// connection, it might be beneficial to wrap the stream in a bufio.Reader.
// E.g.,
//
//	f, err := os.Open("myfile")
//	r := bufio.NewReader(f)
func (b *BitSet) ReadFrom(stream io.Reader) (int64, error) {
	var length uint64
	err := binary.Read(stream, binaryOrder, &length)
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return 0, err
	}
	newlength := uint(length)

	if uint64(newlength) != length {
		return 0, errors.New("unmarshalling error: type mismatch")
	}
	nWords := wordsNeeded(uint(newlength))
	if cap(b.set) >= nWords {
		b.set = b.set[:nWords]
	} else {
		b.set = make([]uint64, nWords)
	}

	b.length = newlength

	err = readUint64Array(stream, b.set)
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		// We do not want to leave the BitSet partially filled as
		// it is error prone.
		b.set = b.set[:0]
		b.length = 0
		return 0, err
	}

	return int64(b.BinaryStorageSize()), nil
}

// MarshalBinary encodes a BitSet into a binary form and returns the result.
// Please see WriteTo for details.
func (b *BitSet) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	_, err := b.WriteTo(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), err
}

// UnmarshalBinary decodes the binary form generated by MarshalBinary.
// Please see WriteTo for details.
func (b *BitSet) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)
	_, err := b.ReadFrom(buf)
	return err
}

// MarshalJSON marshals a BitSet as a JSON structure
func (b BitSet) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, b.BinaryStorageSize()))
	_, err := b.WriteTo(buffer)
	if err != nil {
		return nil, err
	}

	// URLEncode all bytes
	return json.Marshal(base64Encoding.EncodeToString(buffer.Bytes()))
}

// UnmarshalJSON unmarshals a BitSet from JSON created using MarshalJSON
func (b *BitSet) UnmarshalJSON(data []byte) error {
	// Unmarshal as string
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	// URLDecode string
	buf, err := base64Encoding.DecodeString(s)
	if err != nil {
		return err
	}

	_, err = b.ReadFrom(bytes.NewReader(buf))
	return err
}

// Rank returns the number of set bits up to and including the index
// that are set in the bitset.
// See https://en.wikipedia.org/wiki/Ranking#Ranking_in_statistics
func (b *BitSet) Rank(index uint) uint {
	if index >= b.length {
		return b.Count()
	}
	leftover := (index + 1) & 63
	answer := uint(popcntSlice(b.set[:(index+1)>>6]))
	if leftover != 0 {
		answer += uint(bits.OnesCount64(b.set[(index+1)>>6] << (64 - leftover)))
	}
	return answer
}

// Select returns the index of the jth set bit, where j is the argument.
// The caller is responsible to ensure that 0 <= j < Count(): when j is
// out of range, the function returns the length of the bitset (b.length).
//
// Note that this function differs in convention from the Rank function which
// returns 1 when ranking the smallest value. We follow the conventional
// textbook definition of Select and Rank.
func (b *BitSet) Select(index uint) uint {
	leftover := index
	for idx, word := range b.set {
		w := uint(bits.OnesCount64(word))
		if w > leftover {
			return uint(idx)*64 + select64(word, leftover)
		}
		leftover -= w
	}
	return b.length
}

// top detects the top bit set
func (b *BitSet) top() (uint, bool) {
	for idx := len(b.set) - 1; idx >= 0; idx-- {
		if word := b.set[idx]; word != 0 {
			return uint(idx<<log2WordSize+bits.Len64(word)) - 1, true
		}
	}

	return 0, false
}

// ShiftLeft shifts the bitset like << operation would do.
//
// Left shift may require bitset size extension. We try to avoid the
// unnecessary memory operations by detecting the leftmost set bit.
// The function will panic if shift causes excess of capacity.
func (b *BitSet) ShiftLeft(bits uint) {
	panicIfNull(b)

	if bits == 0 {
		return
	}

	top, ok := b.top()
	if !ok {
		return
	}

	// capacity check
	if top+bits < bits {
		panic("You are exceeding the capacity")
	}

	// destination set
	dst := b.set

	// not using extendSet() to avoid unneeded data copying
	nsize := wordsNeeded(top + bits)
	if len(b.set) < nsize {
		dst = make([]uint64, nsize)
	}
	if top+bits >= b.length {
		b.length = top + bits + 1
	}

	pad, idx := top%wordSize, top>>log2WordSize
	shift, pages := bits%wordSize, bits>>log2WordSize
	if bits%wordSize == 0 { // happy case: just add pages
		copy(dst[pages:nsize], b.set)
	} else {
		if pad+shift >= wordSize {
			dst[idx+pages+1] = b.set[idx] >> (wordSize - shift)
		}

		for i := int(idx); i >= 0; i-- {
			if i > 0 {
				dst[i+int(pages)] = (b.set[i] << shift) | (b.set[i-1] >> (wordSize - shift))
			} else {
				dst[i+int(pages)] = b.set[i] << shift
			}
		}
	}

	// zeroing extra pages
	for i := 0; i < int(pages); i++ {
		dst[i] = 0
	}

	b.set = dst
}

// ShiftRight shifts the bitset like >> operation would do.
func (b *BitSet) ShiftRight(bits uint) {
	panicIfNull(b)

	if bits == 0 {
		return
	}

	top, ok := b.top()
	if !ok {
		return
	}

	if bits >= top {
		b.set = make([]uint64, wordsNeeded(b.length))
		return
	}

	pad, idx := top%wordSize, top>>log2WordSize
	shift, pages := bits%wordSize, bits>>log2WordSize
	if bits%wordSize == 0 { // happy case: just clear pages
		b.set = b.set[pages:]
		b.length -= pages * wordSize
	} else {
		for i := 0; i <= int(idx-pages); i++ {
			if i < int(idx-pages) {
				b.set[i] = (b.set[i+int(pages)] >> shift) | (b.set[i+int(pages)+1] << (wordSize - shift))
			} else {
				b.set[i] = b.set[i+int(pages)] >> shift
			}
		}

		if pad < shift {
			b.set[int(idx-pages)] = 0
		}
	}

	for i := int(idx-pages) + 1; i <= int(idx); i++ {
		b.set[i] = 0
	}
}

// OnesBetween returns the number of set bits in the range [from, to).
// The range is inclusive of 'from' and exclusive of 'to'.
// Returns 0 if from >= to.
func (b *BitSet) OnesBetween(from, to uint) uint {
	panicIfNull(b)

	if from >= to {
		return 0
	}

	// Calculate indices and masks for the starting and ending words
	startWord := from >> log2WordSize // Divide by wordSize
	endWord := to >> log2WordSize
	startOffset := from & wordMask // Mod wordSize
	endOffset := to & wordMask

	// Case 1: Bits lie within a single word
	if startWord == endWord {
		// Create mask for bits between from and to
		mask := uint64((1<<endOffset)-1) &^ ((1 << startOffset) - 1)
		return uint(bits.OnesCount64(b.set[startWord] & mask))
	}

	var count uint

	// Case 2: Bits span multiple words
	// 2a: Count bits in first word (from startOffset to end of word)
	startMask := ^uint64((1 << startOffset) - 1) // Mask for bits >= startOffset
	count = uint(bits.OnesCount64(b.set[startWord] & startMask))

	// 2b: Count all bits in complete words between start and end
	if endWord > startWord+1 {
		count += uint(popcntSlice(b.set[startWord+1 : endWord]))
	}

	// 2c: Count bits in last word (from start of word to endOffset)
	if endOffset > 0 {
		endMask := uint64(1<<endOffset) - 1 // Mask for bits < endOffset
		count += uint(bits.OnesCount64(b.set[endWord] & endMask))
	}

	return count
}

// Extract extracts bits according to a mask and returns the result
// in a new BitSet. See ExtractTo for details.
func (b *BitSet) Extract(mask *BitSet) *BitSet {
	dst := New(mask.Count())
	b.ExtractTo(mask, dst)
	return dst
}

// ExtractTo copies bits from the BitSet using positions specified in mask
// into a compacted form in dst. The number of set bits in mask determines
// the number of bits that will be extracted.
//
// For example, if mask has bits set at positions 1,4,5, then ExtractTo will
// take bits at those positions from the source BitSet and pack them into
// consecutive positions 0,1,2 in the destination BitSet.
func (b *BitSet) ExtractTo(mask *BitSet, dst *BitSet) {
	panicIfNull(b)
	panicIfNull(mask)
	panicIfNull(dst)

	if len(mask.set) == 0 || len(b.set) == 0 {
		return
	}

	// Ensure destination has enough space for extracted bits
	resultBits := uint(popcntSlice(mask.set))
	if dst.length < resultBits {
		dst.extendSet(resultBits - 1)
	}

	outPos := uint(0)
	length := len(mask.set)
	if len(b.set) < length {
		length = len(b.set)
	}

	// Process each word
	for i := 0; i < length; i++ {
		if mask.set[i] == 0 {
			continue // Skip words with no bits to extract
		}

		// Extract and compact bits according to mask
		extracted := pext(b.set[i], mask.set[i])
		bitsExtracted := uint(bits.OnesCount64(mask.set[i]))

		// Calculate destination position
		wordIdx := outPos >> log2WordSize
		bitOffset := outPos & wordMask

		// Write extracted bits, handling word boundary crossing
		dst.set[wordIdx] |= extracted << bitOffset
		if bitOffset+bitsExtracted > wordSize {
			dst.set[wordIdx+1] = extracted >> (wordSize - bitOffset)
		}

		outPos += bitsExtracted
	}
}

// Deposit creates a new BitSet and deposits bits according to a mask.
// See DepositTo for details.
func (b *BitSet) Deposit(mask *BitSet) *BitSet {
	dst := New(mask.length)
	b.DepositTo(mask, dst)
	return dst
}

// DepositTo spreads bits from a compacted form in the BitSet into positions
// specified by mask in dst. This is the inverse operation of Extract.
//
// For example, if mask has bits set at positions 1,4,5, then DepositTo will
// take consecutive bits 0,1,2 from the source BitSet and place them into
// positions 1,4,5 in the destination BitSet.
func (b *BitSet) DepositTo(mask *BitSet, dst *BitSet) {
	panicIfNull(b)
	panicIfNull(mask)
	panicIfNull(dst)

	if len(dst.set) == 0 || len(mask.set) == 0 || len(b.set) == 0 {
		return
	}

	inPos := uint(0)
	length := len(mask.set)
	if len(dst.set) < length {
		length = len(dst.set)
	}

	// Process each word
	for i := 0; i < length; i++ {
		if mask.set[i] == 0 {
			continue // Skip words with no bits to deposit
		}

		// Calculate source word index
		wordIdx := inPos >> log2WordSize
		if wordIdx >= uint(len(b.set)) {
			break // No more source bits available
		}

		// Get source bits, handling word boundary crossing
		sourceBits := b.set[wordIdx]
		bitOffset := inPos & wordMask
		if wordIdx+1 < uint(len(b.set)) && bitOffset != 0 {
			// Combine bits from current and next word
			sourceBits = (sourceBits >> bitOffset) |
				(b.set[wordIdx+1] << (wordSize - bitOffset))
		} else {
			sourceBits >>= bitOffset
		}

		// Deposit bits according to mask
		dst.set[i] = (dst.set[i] &^ mask.set[i]) | pdep(sourceBits, mask.set[i])
		inPos += uint(bits.OnesCount64(mask.set[i]))
	}
}

func pext(w, m uint64) (result uint64) {
	var outPos uint

	// Process byte by byte
	for i := 0; i < 8; i++ {
		shift := i << 3 // i * 8 using bit shift
		b := uint8(w >> shift)
		mask := uint8(m >> shift)

		extracted := pextLUT[b][mask]
		bits := popLUT[mask]

		result |= uint64(extracted) << outPos
		outPos += uint(bits)
	}

	return result
}

func pdep(w, m uint64) (result uint64) {
	var inPos uint

	// Process byte by byte
	for i := 0; i < 8; i++ {
		shift := i << 3 // i * 8 using bit shift
		mask := uint8(m >> shift)
		bits := popLUT[mask]

		// Get the bits we'll deposit from the source
		b := uint8(w >> inPos)

		// Deposit them according to the mask for this byte
		deposited := pdepLUT[b][mask]

		// Add to result
		result |= uint64(deposited) << shift
		inPos += uint(bits)
	}

	return result
}

func popcntSlice(s []uint64) uint64 {
	var cnt int
	for _, x := range s {
		cnt += bits.OnesCount64(x)
	}
	return uint64(cnt)
}

func popcntMaskSlice(s, m []uint64) uint64 {
	var cnt int
	// this explicit check eliminates a bounds check in the loop
	if len(m) < len(s) {
		panic("mask slice is too short")
	}
	for i := range s {
		cnt += bits.OnesCount64(s[i] &^ m[i])
	}
	return uint64(cnt)
}

func popcntAndSlice(s, m []uint64) uint64 {
	var cnt int
	// this explicit check eliminates a bounds check in the loop
	if len(m) < len(s) {
		panic("mask slice is too short")
	}
	for i := range s {
		cnt += bits.OnesCount64(s[i] & m[i])
	}
	return uint64(cnt)
}

func popcntOrSlice(s, m []uint64) uint64 {
	var cnt int
	// this explicit check eliminates a bounds check in the loop
	if len(m) < len(s) {
		panic("mask slice is too short")
	}
	for i := range s {
		cnt += bits.OnesCount64(s[i] | m[i])
	}
	return uint64(cnt)
}

func popcntXorSlice(s, m []uint64) uint64 {
	var cnt int
	// this explicit check eliminates a bounds check in the loop
	if len(m) < len(s) {
		panic("mask slice is too short")
	}
	for i := range s {
		cnt += bits.OnesCount64(s[i] ^ m[i])
	}
	return uint64(cnt)
}

func select64(w uint64, j uint) uint {
	seen := 0
	// Divide 64bit
	part := w & 0xFFFFFFFF
	n := uint(bits.OnesCount64(part))
	if n <= j {
		part = w >> 32
		seen += 32
		j -= n
	}
	ww := part

	// Divide 32bit
	part = ww & 0xFFFF

	n = uint(bits.OnesCount64(part))
	if n <= j {
		part = ww >> 16
		seen += 16
		j -= n
	}
	ww = part

	// Divide 16bit
	part = ww & 0xFF
	n = uint(bits.OnesCount64(part))
	if n <= j {
		part = ww >> 8
		seen += 8
		j -= n
	}
	ww = part

	// Lookup in final byte
	counter := 0
	for ; counter < 8; counter++ {
		j -= uint((ww >> counter) & 1)
		if j+1 == 0 {
			break
		}
	}
	return uint(seen + counter)
}

// pextLUT contains pre-computed parallel bit extraction results
var pextLUT = [256][256]uint8{
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		4, 8, 8, 16, 8, 16, 16, 32, 8, 16, 16, 32, 16, 32, 32, 64,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		4, 8, 8, 16, 8, 16, 16, 32, 8, 16, 16, 32, 16, 32, 32, 64,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		4, 9, 8, 17, 8, 17, 16, 33, 8, 17, 16, 33, 16, 33, 32, 65,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		4, 9, 8, 17, 8, 17, 16, 33, 8, 17, 16, 33, 16, 33, 32, 65,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		4, 8, 9, 18, 8, 16, 17, 34, 8, 16, 17, 34, 16, 32, 33, 66,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		4, 8, 9, 18, 8, 16, 17, 34, 8, 16, 17, 34, 16, 32, 33, 66,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		4, 9, 9, 19, 8, 17, 17, 35, 8, 17, 17, 35, 16, 33, 33, 67,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		4, 9, 9, 19, 8, 17, 17, 35, 8, 17, 17, 35, 16, 33, 33, 67,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		4, 8, 8, 16, 9, 18, 18, 36, 8, 16, 16, 32, 17, 34, 34, 68,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		4, 8, 8, 16, 9, 18, 18, 36, 8, 16, 16, 32, 17, 34, 34, 68,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		4, 9, 8, 17, 9, 19, 18, 37, 8, 17, 16, 33, 17, 35, 34, 69,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		4, 9, 8, 17, 9, 19, 18, 37, 8, 17, 16, 33, 17, 35, 34, 69,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		4, 8, 9, 18, 9, 18, 19, 38, 8, 16, 17, 34, 17, 34, 35, 70,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		4, 8, 9, 18, 9, 18, 19, 38, 8, 16, 17, 34, 17, 34, 35, 70,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		4, 9, 9, 19, 9, 19, 19, 39, 8, 17, 17, 35, 17, 35, 35, 71,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		4, 9, 9, 19, 9, 19, 19, 39, 8, 17, 17, 35, 17, 35, 35, 71,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		4, 8, 8, 16, 8, 16, 16, 32, 9, 18, 18, 36, 18, 36, 36, 72,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		4, 8, 8, 16, 8, 16, 16, 32, 9, 18, 18, 36, 18, 36, 36, 72,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		4, 9, 8, 17, 8, 17, 16, 33, 9, 19, 18, 37, 18, 37, 36, 73,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		4, 9, 8, 17, 8, 17, 16, 33, 9, 19, 18, 37, 18, 37, 36, 73,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		4, 8, 9, 18, 8, 16, 17, 34, 9, 18, 19, 38, 18, 36, 37, 74,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		4, 8, 9, 18, 8, 16, 17, 34, 9, 18, 19, 38, 18, 36, 37, 74,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		4, 9, 9, 19, 8, 17, 17, 35, 9, 19, 19, 39, 18, 37, 37, 75,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		4, 9, 9, 19, 8, 17, 17, 35, 9, 19, 19, 39, 18, 37, 37, 75,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		4, 8, 8, 16, 9, 18, 18, 36, 9, 18, 18, 36, 19, 38, 38, 76,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		4, 8, 8, 16, 9, 18, 18, 36, 9, 18, 18, 36, 19, 38, 38, 76,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		4, 9, 8, 17, 9, 19, 18, 37, 9, 19, 18, 37, 19, 39, 38, 77,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		4, 9, 8, 17, 9, 19, 18, 37, 9, 19, 18, 37, 19, 39, 38, 77,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		4, 8, 9, 18, 9, 18, 19, 38, 9, 18, 19, 38, 19, 38, 39, 78,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		4, 8, 9, 18, 9, 18, 19, 38, 9, 18, 19, 38, 19, 38, 39, 78,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		4, 9, 9, 19, 9, 19, 19, 39, 9, 19, 19, 39, 19, 39, 39, 79,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		4, 9, 9, 19, 9, 19, 19, 39, 9, 19, 19, 39, 19, 39, 39, 79,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		5, 10, 10, 20, 10, 20, 20, 40, 10, 20, 20, 40, 20, 40, 40, 80,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		5, 10, 10, 20, 10, 20, 20, 40, 10, 20, 20, 40, 20, 40, 40, 80,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		5, 11, 10, 21, 10, 21, 20, 41, 10, 21, 20, 41, 20, 41, 40, 81,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		5, 11, 10, 21, 10, 21, 20, 41, 10, 21, 20, 41, 20, 41, 40, 81,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		5, 10, 11, 22, 10, 20, 21, 42, 10, 20, 21, 42, 20, 40, 41, 82,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		5, 10, 11, 22, 10, 20, 21, 42, 10, 20, 21, 42, 20, 40, 41, 82,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		5, 11, 11, 23, 10, 21, 21, 43, 10, 21, 21, 43, 20, 41, 41, 83,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		5, 11, 11, 23, 10, 21, 21, 43, 10, 21, 21, 43, 20, 41, 41, 83,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		5, 10, 10, 20, 11, 22, 22, 44, 10, 20, 20, 40, 21, 42, 42, 84,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		5, 10, 10, 20, 11, 22, 22, 44, 10, 20, 20, 40, 21, 42, 42, 84,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		5, 11, 10, 21, 11, 23, 22, 45, 10, 21, 20, 41, 21, 43, 42, 85,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		5, 11, 10, 21, 11, 23, 22, 45, 10, 21, 20, 41, 21, 43, 42, 85,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		5, 10, 11, 22, 11, 22, 23, 46, 10, 20, 21, 42, 21, 42, 43, 86,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		5, 10, 11, 22, 11, 22, 23, 46, 10, 20, 21, 42, 21, 42, 43, 86,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		5, 11, 11, 23, 11, 23, 23, 47, 10, 21, 21, 43, 21, 43, 43, 87,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		5, 11, 11, 23, 11, 23, 23, 47, 10, 21, 21, 43, 21, 43, 43, 87,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		5, 10, 10, 20, 10, 20, 20, 40, 11, 22, 22, 44, 22, 44, 44, 88,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		5, 10, 10, 20, 10, 20, 20, 40, 11, 22, 22, 44, 22, 44, 44, 88,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		5, 11, 10, 21, 10, 21, 20, 41, 11, 23, 22, 45, 22, 45, 44, 89,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		5, 11, 10, 21, 10, 21, 20, 41, 11, 23, 22, 45, 22, 45, 44, 89,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		5, 10, 11, 22, 10, 20, 21, 42, 11, 22, 23, 46, 22, 44, 45, 90,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		5, 10, 11, 22, 10, 20, 21, 42, 11, 22, 23, 46, 22, 44, 45, 90,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		5, 11, 11, 23, 10, 21, 21, 43, 11, 23, 23, 47, 22, 45, 45, 91,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		5, 11, 11, 23, 10, 21, 21, 43, 11, 23, 23, 47, 22, 45, 45, 91,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		5, 10, 10, 20, 11, 22, 22, 44, 11, 22, 22, 44, 23, 46, 46, 92,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		5, 10, 10, 20, 11, 22, 22, 44, 11, 22, 22, 44, 23, 46, 46, 92,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		5, 11, 10, 21, 11, 23, 22, 45, 11, 23, 22, 45, 23, 47, 46, 93,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		5, 11, 10, 21, 11, 23, 22, 45, 11, 23, 22, 45, 23, 47, 46, 93,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		5, 10, 11, 22, 11, 22, 23, 46, 11, 22, 23, 46, 23, 46, 47, 94,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		5, 10, 11, 22, 11, 22, 23, 46, 11, 22, 23, 46, 23, 46, 47, 94,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		5, 11, 11, 23, 11, 23, 23, 47, 11, 23, 23, 47, 23, 47, 47, 95,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		5, 11, 11, 23, 11, 23, 23, 47, 11, 23, 23, 47, 23, 47, 47, 95,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		6, 12, 12, 24, 12, 24, 24, 48, 12, 24, 24, 48, 24, 48, 48, 96,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		6, 12, 12, 24, 12, 24, 24, 48, 12, 24, 24, 48, 24, 48, 48, 96,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		6, 13, 12, 25, 12, 25, 24, 49, 12, 25, 24, 49, 24, 49, 48, 97,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		6, 13, 12, 25, 12, 25, 24, 49, 12, 25, 24, 49, 24, 49, 48, 97,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		6, 12, 13, 26, 12, 24, 25, 50, 12, 24, 25, 50, 24, 48, 49, 98,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		6, 12, 13, 26, 12, 24, 25, 50, 12, 24, 25, 50, 24, 48, 49, 98,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		6, 13, 13, 27, 12, 25, 25, 51, 12, 25, 25, 51, 24, 49, 49, 99,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		6, 13, 13, 27, 12, 25, 25, 51, 12, 25, 25, 51, 24, 49, 49, 99,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		6, 12, 12, 24, 13, 26, 26, 52, 12, 24, 24, 48, 25, 50, 50, 100,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		6, 12, 12, 24, 13, 26, 26, 52, 12, 24, 24, 48, 25, 50, 50, 100,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		6, 13, 12, 25, 13, 27, 26, 53, 12, 25, 24, 49, 25, 51, 50, 101,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		6, 13, 12, 25, 13, 27, 26, 53, 12, 25, 24, 49, 25, 51, 50, 101,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		6, 12, 13, 26, 13, 26, 27, 54, 12, 24, 25, 50, 25, 50, 51, 102,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		6, 12, 13, 26, 13, 26, 27, 54, 12, 24, 25, 50, 25, 50, 51, 102,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		6, 13, 13, 27, 13, 27, 27, 55, 12, 25, 25, 51, 25, 51, 51, 103,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		6, 13, 13, 27, 13, 27, 27, 55, 12, 25, 25, 51, 25, 51, 51, 103,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		6, 12, 12, 24, 12, 24, 24, 48, 13, 26, 26, 52, 26, 52, 52, 104,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		6, 12, 12, 24, 12, 24, 24, 48, 13, 26, 26, 52, 26, 52, 52, 104,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		6, 13, 12, 25, 12, 25, 24, 49, 13, 27, 26, 53, 26, 53, 52, 105,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		6, 13, 12, 25, 12, 25, 24, 49, 13, 27, 26, 53, 26, 53, 52, 105,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		6, 12, 13, 26, 12, 24, 25, 50, 13, 26, 27, 54, 26, 52, 53, 106,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		6, 12, 13, 26, 12, 24, 25, 50, 13, 26, 27, 54, 26, 52, 53, 106,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		6, 13, 13, 27, 12, 25, 25, 51, 13, 27, 27, 55, 26, 53, 53, 107,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		6, 13, 13, 27, 12, 25, 25, 51, 13, 27, 27, 55, 26, 53, 53, 107,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		6, 12, 12, 24, 13, 26, 26, 52, 13, 26, 26, 52, 27, 54, 54, 108,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		6, 12, 12, 24, 13, 26, 26, 52, 13, 26, 26, 52, 27, 54, 54, 108,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		6, 13, 12, 25, 13, 27, 26, 53, 13, 27, 26, 53, 27, 55, 54, 109,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		6, 13, 12, 25, 13, 27, 26, 53, 13, 27, 26, 53, 27, 55, 54, 109,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		6, 12, 13, 26, 13, 26, 27, 54, 13, 26, 27, 54, 27, 54, 55, 110,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		6, 12, 13, 26, 13, 26, 27, 54, 13, 26, 27, 54, 27, 54, 55, 110,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		6, 13, 13, 27, 13, 27, 27, 55, 13, 27, 27, 55, 27, 55, 55, 111,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		6, 13, 13, 27, 13, 27, 27, 55, 13, 27, 27, 55, 27, 55, 55, 111,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		7, 14, 14, 28, 14, 28, 28, 56, 14, 28, 28, 56, 28, 56, 56, 112,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		7, 14, 14, 28, 14, 28, 28, 56, 14, 28, 28, 56, 28, 56, 56, 112,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		7, 15, 14, 29, 14, 29, 28, 57, 14, 29, 28, 57, 28, 57, 56, 113,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		7, 15, 14, 29, 14, 29, 28, 57, 14, 29, 28, 57, 28, 57, 56, 113,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		7, 14, 15, 30, 14, 28, 29, 58, 14, 28, 29, 58, 28, 56, 57, 114,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		7, 14, 15, 30, 14, 28, 29, 58, 14, 28, 29, 58, 28, 56, 57, 114,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		7, 15, 15, 31, 14, 29, 29, 59, 14, 29, 29, 59, 28, 57, 57, 115,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		7, 15, 15, 31, 14, 29, 29, 59, 14, 29, 29, 59, 28, 57, 57, 115,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		7, 14, 14, 28, 15, 30, 30, 60, 14, 28, 28, 56, 29, 58, 58, 116,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		7, 14, 14, 28, 15, 30, 30, 60, 14, 28, 28, 56, 29, 58, 58, 116,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		7, 15, 14, 29, 15, 31, 30, 61, 14, 29, 28, 57, 29, 59, 58, 117,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		7, 15, 14, 29, 15, 31, 30, 61, 14, 29, 28, 57, 29, 59, 58, 117,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		7, 14, 15, 30, 15, 30, 31, 62, 14, 28, 29, 58, 29, 58, 59, 118,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		7, 14, 15, 30, 15, 30, 31, 62, 14, 28, 29, 58, 29, 58, 59, 118,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		7, 15, 15, 31, 15, 31, 31, 63, 14, 29, 29, 59, 29, 59, 59, 119,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		7, 15, 15, 31, 15, 31, 31, 63, 14, 29, 29, 59, 29, 59, 59, 119,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		7, 14, 14, 28, 14, 28, 28, 56, 15, 30, 30, 60, 30, 60, 60, 120,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		7, 14, 14, 28, 14, 28, 28, 56, 15, 30, 30, 60, 30, 60, 60, 120,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		7, 15, 14, 29, 14, 29, 28, 57, 15, 31, 30, 61, 30, 61, 60, 121,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		7, 15, 14, 29, 14, 29, 28, 57, 15, 31, 30, 61, 30, 61, 60, 121,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		7, 14, 15, 30, 14, 28, 29, 58, 15, 30, 31, 62, 30, 60, 61, 122,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		7, 14, 15, 30, 14, 28, 29, 58, 15, 30, 31, 62, 30, 60, 61, 122,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		7, 15, 15, 31, 14, 29, 29, 59, 15, 31, 31, 63, 30, 61, 61, 123,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		7, 15, 15, 31, 14, 29, 29, 59, 15, 31, 31, 63, 30, 61, 61, 123,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		7, 14, 14, 28, 15, 30, 30, 60, 15, 30, 30, 60, 31, 62, 62, 124,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		7, 14, 14, 28, 15, 30, 30, 60, 15, 30, 30, 60, 31, 62, 62, 124,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		7, 15, 14, 29, 15, 31, 30, 61, 15, 31, 30, 61, 31, 63, 62, 125,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		7, 15, 14, 29, 15, 31, 30, 61, 15, 31, 30, 61, 31, 63, 62, 125,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		7, 14, 15, 30, 15, 30, 31, 62, 15, 30, 31, 62, 31, 62, 63, 126,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		7, 14, 15, 30, 15, 30, 31, 62, 15, 30, 31, 62, 31, 62, 63, 126,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		7, 15, 15, 31, 15, 31, 31, 63, 15, 31, 31, 63, 31, 63, 63, 127,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		7, 15, 15, 31, 15, 31, 31, 63, 15, 31, 31, 63, 31, 63, 63, 127,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		4, 8, 8, 16, 8, 16, 16, 32, 8, 16, 16, 32, 16, 32, 32, 64,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		4, 8, 8, 16, 8, 16, 16, 32, 8, 16, 16, 32, 16, 32, 32, 64,
		4, 8, 8, 16, 8, 16, 16, 32, 8, 16, 16, 32, 16, 32, 32, 64,
		8, 16, 16, 32, 16, 32, 32, 64, 16, 32, 32, 64, 32, 64, 64, 128,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		4, 9, 8, 17, 8, 17, 16, 33, 8, 17, 16, 33, 16, 33, 32, 65,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		4, 9, 8, 17, 8, 17, 16, 33, 8, 17, 16, 33, 16, 33, 32, 65,
		4, 9, 8, 17, 8, 17, 16, 33, 8, 17, 16, 33, 16, 33, 32, 65,
		8, 17, 16, 33, 16, 33, 32, 65, 16, 33, 32, 65, 32, 65, 64, 129,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		4, 8, 9, 18, 8, 16, 17, 34, 8, 16, 17, 34, 16, 32, 33, 66,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		4, 8, 9, 18, 8, 16, 17, 34, 8, 16, 17, 34, 16, 32, 33, 66,
		4, 8, 9, 18, 8, 16, 17, 34, 8, 16, 17, 34, 16, 32, 33, 66,
		8, 16, 17, 34, 16, 32, 33, 66, 16, 32, 33, 66, 32, 64, 65, 130,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		4, 9, 9, 19, 8, 17, 17, 35, 8, 17, 17, 35, 16, 33, 33, 67,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		4, 9, 9, 19, 8, 17, 17, 35, 8, 17, 17, 35, 16, 33, 33, 67,
		4, 9, 9, 19, 8, 17, 17, 35, 8, 17, 17, 35, 16, 33, 33, 67,
		8, 17, 17, 35, 16, 33, 33, 67, 16, 33, 33, 67, 32, 65, 65, 131,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		4, 8, 8, 16, 9, 18, 18, 36, 8, 16, 16, 32, 17, 34, 34, 68,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		4, 8, 8, 16, 9, 18, 18, 36, 8, 16, 16, 32, 17, 34, 34, 68,
		4, 8, 8, 16, 9, 18, 18, 36, 8, 16, 16, 32, 17, 34, 34, 68,
		8, 16, 16, 32, 17, 34, 34, 68, 16, 32, 32, 64, 33, 66, 66, 132,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		4, 9, 8, 17, 9, 19, 18, 37, 8, 17, 16, 33, 17, 35, 34, 69,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		4, 9, 8, 17, 9, 19, 18, 37, 8, 17, 16, 33, 17, 35, 34, 69,
		4, 9, 8, 17, 9, 19, 18, 37, 8, 17, 16, 33, 17, 35, 34, 69,
		8, 17, 16, 33, 17, 35, 34, 69, 16, 33, 32, 65, 33, 67, 66, 133,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		4, 8, 9, 18, 9, 18, 19, 38, 8, 16, 17, 34, 17, 34, 35, 70,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		4, 8, 9, 18, 9, 18, 19, 38, 8, 16, 17, 34, 17, 34, 35, 70,
		4, 8, 9, 18, 9, 18, 19, 38, 8, 16, 17, 34, 17, 34, 35, 70,
		8, 16, 17, 34, 17, 34, 35, 70, 16, 32, 33, 66, 33, 66, 67, 134,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		4, 9, 9, 19, 9, 19, 19, 39, 8, 17, 17, 35, 17, 35, 35, 71,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		4, 9, 9, 19, 9, 19, 19, 39, 8, 17, 17, 35, 17, 35, 35, 71,
		4, 9, 9, 19, 9, 19, 19, 39, 8, 17, 17, 35, 17, 35, 35, 71,
		8, 17, 17, 35, 17, 35, 35, 71, 16, 33, 33, 67, 33, 67, 67, 135,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		4, 8, 8, 16, 8, 16, 16, 32, 9, 18, 18, 36, 18, 36, 36, 72,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		4, 8, 8, 16, 8, 16, 16, 32, 9, 18, 18, 36, 18, 36, 36, 72,
		4, 8, 8, 16, 8, 16, 16, 32, 9, 18, 18, 36, 18, 36, 36, 72,
		8, 16, 16, 32, 16, 32, 32, 64, 17, 34, 34, 68, 34, 68, 68, 136,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		4, 9, 8, 17, 8, 17, 16, 33, 9, 19, 18, 37, 18, 37, 36, 73,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		4, 9, 8, 17, 8, 17, 16, 33, 9, 19, 18, 37, 18, 37, 36, 73,
		4, 9, 8, 17, 8, 17, 16, 33, 9, 19, 18, 37, 18, 37, 36, 73,
		8, 17, 16, 33, 16, 33, 32, 65, 17, 35, 34, 69, 34, 69, 68, 137,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		4, 8, 9, 18, 8, 16, 17, 34, 9, 18, 19, 38, 18, 36, 37, 74,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		4, 8, 9, 18, 8, 16, 17, 34, 9, 18, 19, 38, 18, 36, 37, 74,
		4, 8, 9, 18, 8, 16, 17, 34, 9, 18, 19, 38, 18, 36, 37, 74,
		8, 16, 17, 34, 16, 32, 33, 66, 17, 34, 35, 70, 34, 68, 69, 138,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		4, 9, 9, 19, 8, 17, 17, 35, 9, 19, 19, 39, 18, 37, 37, 75,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		4, 9, 9, 19, 8, 17, 17, 35, 9, 19, 19, 39, 18, 37, 37, 75,
		4, 9, 9, 19, 8, 17, 17, 35, 9, 19, 19, 39, 18, 37, 37, 75,
		8, 17, 17, 35, 16, 33, 33, 67, 17, 35, 35, 71, 34, 69, 69, 139,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		4, 8, 8, 16, 9, 18, 18, 36, 9, 18, 18, 36, 19, 38, 38, 76,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		4, 8, 8, 16, 9, 18, 18, 36, 9, 18, 18, 36, 19, 38, 38, 76,
		4, 8, 8, 16, 9, 18, 18, 36, 9, 18, 18, 36, 19, 38, 38, 76,
		8, 16, 16, 32, 17, 34, 34, 68, 17, 34, 34, 68, 35, 70, 70, 140,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		4, 9, 8, 17, 9, 19, 18, 37, 9, 19, 18, 37, 19, 39, 38, 77,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		4, 9, 8, 17, 9, 19, 18, 37, 9, 19, 18, 37, 19, 39, 38, 77,
		4, 9, 8, 17, 9, 19, 18, 37, 9, 19, 18, 37, 19, 39, 38, 77,
		8, 17, 16, 33, 17, 35, 34, 69, 17, 35, 34, 69, 35, 71, 70, 141,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		4, 8, 9, 18, 9, 18, 19, 38, 9, 18, 19, 38, 19, 38, 39, 78,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		4, 8, 9, 18, 9, 18, 19, 38, 9, 18, 19, 38, 19, 38, 39, 78,
		4, 8, 9, 18, 9, 18, 19, 38, 9, 18, 19, 38, 19, 38, 39, 78,
		8, 16, 17, 34, 17, 34, 35, 70, 17, 34, 35, 70, 35, 70, 71, 142,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		4, 9, 9, 19, 9, 19, 19, 39, 9, 19, 19, 39, 19, 39, 39, 79,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		4, 9, 9, 19, 9, 19, 19, 39, 9, 19, 19, 39, 19, 39, 39, 79,
		4, 9, 9, 19, 9, 19, 19, 39, 9, 19, 19, 39, 19, 39, 39, 79,
		8, 17, 17, 35, 17, 35, 35, 71, 17, 35, 35, 71, 35, 71, 71, 143,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		5, 10, 10, 20, 10, 20, 20, 40, 10, 20, 20, 40, 20, 40, 40, 80,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		5, 10, 10, 20, 10, 20, 20, 40, 10, 20, 20, 40, 20, 40, 40, 80,
		4, 8, 8, 16, 8, 16, 16, 32, 8, 16, 16, 32, 16, 32, 32, 64,
		9, 18, 18, 36, 18, 36, 36, 72, 18, 36, 36, 72, 36, 72, 72, 144,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		5, 11, 10, 21, 10, 21, 20, 41, 10, 21, 20, 41, 20, 41, 40, 81,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		5, 11, 10, 21, 10, 21, 20, 41, 10, 21, 20, 41, 20, 41, 40, 81,
		4, 9, 8, 17, 8, 17, 16, 33, 8, 17, 16, 33, 16, 33, 32, 65,
		9, 19, 18, 37, 18, 37, 36, 73, 18, 37, 36, 73, 36, 73, 72, 145,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		5, 10, 11, 22, 10, 20, 21, 42, 10, 20, 21, 42, 20, 40, 41, 82,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		5, 10, 11, 22, 10, 20, 21, 42, 10, 20, 21, 42, 20, 40, 41, 82,
		4, 8, 9, 18, 8, 16, 17, 34, 8, 16, 17, 34, 16, 32, 33, 66,
		9, 18, 19, 38, 18, 36, 37, 74, 18, 36, 37, 74, 36, 72, 73, 146,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		5, 11, 11, 23, 10, 21, 21, 43, 10, 21, 21, 43, 20, 41, 41, 83,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		5, 11, 11, 23, 10, 21, 21, 43, 10, 21, 21, 43, 20, 41, 41, 83,
		4, 9, 9, 19, 8, 17, 17, 35, 8, 17, 17, 35, 16, 33, 33, 67,
		9, 19, 19, 39, 18, 37, 37, 75, 18, 37, 37, 75, 36, 73, 73, 147,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		5, 10, 10, 20, 11, 22, 22, 44, 10, 20, 20, 40, 21, 42, 42, 84,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		5, 10, 10, 20, 11, 22, 22, 44, 10, 20, 20, 40, 21, 42, 42, 84,
		4, 8, 8, 16, 9, 18, 18, 36, 8, 16, 16, 32, 17, 34, 34, 68,
		9, 18, 18, 36, 19, 38, 38, 76, 18, 36, 36, 72, 37, 74, 74, 148,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		5, 11, 10, 21, 11, 23, 22, 45, 10, 21, 20, 41, 21, 43, 42, 85,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		5, 11, 10, 21, 11, 23, 22, 45, 10, 21, 20, 41, 21, 43, 42, 85,
		4, 9, 8, 17, 9, 19, 18, 37, 8, 17, 16, 33, 17, 35, 34, 69,
		9, 19, 18, 37, 19, 39, 38, 77, 18, 37, 36, 73, 37, 75, 74, 149,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		5, 10, 11, 22, 11, 22, 23, 46, 10, 20, 21, 42, 21, 42, 43, 86,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		5, 10, 11, 22, 11, 22, 23, 46, 10, 20, 21, 42, 21, 42, 43, 86,
		4, 8, 9, 18, 9, 18, 19, 38, 8, 16, 17, 34, 17, 34, 35, 70,
		9, 18, 19, 38, 19, 38, 39, 78, 18, 36, 37, 74, 37, 74, 75, 150,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		5, 11, 11, 23, 11, 23, 23, 47, 10, 21, 21, 43, 21, 43, 43, 87,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		5, 11, 11, 23, 11, 23, 23, 47, 10, 21, 21, 43, 21, 43, 43, 87,
		4, 9, 9, 19, 9, 19, 19, 39, 8, 17, 17, 35, 17, 35, 35, 71,
		9, 19, 19, 39, 19, 39, 39, 79, 18, 37, 37, 75, 37, 75, 75, 151,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		5, 10, 10, 20, 10, 20, 20, 40, 11, 22, 22, 44, 22, 44, 44, 88,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		5, 10, 10, 20, 10, 20, 20, 40, 11, 22, 22, 44, 22, 44, 44, 88,
		4, 8, 8, 16, 8, 16, 16, 32, 9, 18, 18, 36, 18, 36, 36, 72,
		9, 18, 18, 36, 18, 36, 36, 72, 19, 38, 38, 76, 38, 76, 76, 152,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		5, 11, 10, 21, 10, 21, 20, 41, 11, 23, 22, 45, 22, 45, 44, 89,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		5, 11, 10, 21, 10, 21, 20, 41, 11, 23, 22, 45, 22, 45, 44, 89,
		4, 9, 8, 17, 8, 17, 16, 33, 9, 19, 18, 37, 18, 37, 36, 73,
		9, 19, 18, 37, 18, 37, 36, 73, 19, 39, 38, 77, 38, 77, 76, 153,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		5, 10, 11, 22, 10, 20, 21, 42, 11, 22, 23, 46, 22, 44, 45, 90,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		5, 10, 11, 22, 10, 20, 21, 42, 11, 22, 23, 46, 22, 44, 45, 90,
		4, 8, 9, 18, 8, 16, 17, 34, 9, 18, 19, 38, 18, 36, 37, 74,
		9, 18, 19, 38, 18, 36, 37, 74, 19, 38, 39, 78, 38, 76, 77, 154,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		5, 11, 11, 23, 10, 21, 21, 43, 11, 23, 23, 47, 22, 45, 45, 91,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		5, 11, 11, 23, 10, 21, 21, 43, 11, 23, 23, 47, 22, 45, 45, 91,
		4, 9, 9, 19, 8, 17, 17, 35, 9, 19, 19, 39, 18, 37, 37, 75,
		9, 19, 19, 39, 18, 37, 37, 75, 19, 39, 39, 79, 38, 77, 77, 155,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		5, 10, 10, 20, 11, 22, 22, 44, 11, 22, 22, 44, 23, 46, 46, 92,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		5, 10, 10, 20, 11, 22, 22, 44, 11, 22, 22, 44, 23, 46, 46, 92,
		4, 8, 8, 16, 9, 18, 18, 36, 9, 18, 18, 36, 19, 38, 38, 76,
		9, 18, 18, 36, 19, 38, 38, 76, 19, 38, 38, 76, 39, 78, 78, 156,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		5, 11, 10, 21, 11, 23, 22, 45, 11, 23, 22, 45, 23, 47, 46, 93,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		5, 11, 10, 21, 11, 23, 22, 45, 11, 23, 22, 45, 23, 47, 46, 93,
		4, 9, 8, 17, 9, 19, 18, 37, 9, 19, 18, 37, 19, 39, 38, 77,
		9, 19, 18, 37, 19, 39, 38, 77, 19, 39, 38, 77, 39, 79, 78, 157,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		5, 10, 11, 22, 11, 22, 23, 46, 11, 22, 23, 46, 23, 46, 47, 94,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		5, 10, 11, 22, 11, 22, 23, 46, 11, 22, 23, 46, 23, 46, 47, 94,
		4, 8, 9, 18, 9, 18, 19, 38, 9, 18, 19, 38, 19, 38, 39, 78,
		9, 18, 19, 38, 19, 38, 39, 78, 19, 38, 39, 78, 39, 78, 79, 158,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		5, 11, 11, 23, 11, 23, 23, 47, 11, 23, 23, 47, 23, 47, 47, 95,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		5, 11, 11, 23, 11, 23, 23, 47, 11, 23, 23, 47, 23, 47, 47, 95,
		4, 9, 9, 19, 9, 19, 19, 39, 9, 19, 19, 39, 19, 39, 39, 79,
		9, 19, 19, 39, 19, 39, 39, 79, 19, 39, 39, 79, 39, 79, 79, 159,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		6, 12, 12, 24, 12, 24, 24, 48, 12, 24, 24, 48, 24, 48, 48, 96,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		4, 8, 8, 16, 8, 16, 16, 32, 8, 16, 16, 32, 16, 32, 32, 64,
		5, 10, 10, 20, 10, 20, 20, 40, 10, 20, 20, 40, 20, 40, 40, 80,
		10, 20, 20, 40, 20, 40, 40, 80, 20, 40, 40, 80, 40, 80, 80, 160,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		6, 13, 12, 25, 12, 25, 24, 49, 12, 25, 24, 49, 24, 49, 48, 97,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		4, 9, 8, 17, 8, 17, 16, 33, 8, 17, 16, 33, 16, 33, 32, 65,
		5, 11, 10, 21, 10, 21, 20, 41, 10, 21, 20, 41, 20, 41, 40, 81,
		10, 21, 20, 41, 20, 41, 40, 81, 20, 41, 40, 81, 40, 81, 80, 161,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		6, 12, 13, 26, 12, 24, 25, 50, 12, 24, 25, 50, 24, 48, 49, 98,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		4, 8, 9, 18, 8, 16, 17, 34, 8, 16, 17, 34, 16, 32, 33, 66,
		5, 10, 11, 22, 10, 20, 21, 42, 10, 20, 21, 42, 20, 40, 41, 82,
		10, 20, 21, 42, 20, 40, 41, 82, 20, 40, 41, 82, 40, 80, 81, 162,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		6, 13, 13, 27, 12, 25, 25, 51, 12, 25, 25, 51, 24, 49, 49, 99,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		4, 9, 9, 19, 8, 17, 17, 35, 8, 17, 17, 35, 16, 33, 33, 67,
		5, 11, 11, 23, 10, 21, 21, 43, 10, 21, 21, 43, 20, 41, 41, 83,
		10, 21, 21, 43, 20, 41, 41, 83, 20, 41, 41, 83, 40, 81, 81, 163,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		6, 12, 12, 24, 13, 26, 26, 52, 12, 24, 24, 48, 25, 50, 50, 100,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		4, 8, 8, 16, 9, 18, 18, 36, 8, 16, 16, 32, 17, 34, 34, 68,
		5, 10, 10, 20, 11, 22, 22, 44, 10, 20, 20, 40, 21, 42, 42, 84,
		10, 20, 20, 40, 21, 42, 42, 84, 20, 40, 40, 80, 41, 82, 82, 164,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		6, 13, 12, 25, 13, 27, 26, 53, 12, 25, 24, 49, 25, 51, 50, 101,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		4, 9, 8, 17, 9, 19, 18, 37, 8, 17, 16, 33, 17, 35, 34, 69,
		5, 11, 10, 21, 11, 23, 22, 45, 10, 21, 20, 41, 21, 43, 42, 85,
		10, 21, 20, 41, 21, 43, 42, 85, 20, 41, 40, 81, 41, 83, 82, 165,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		6, 12, 13, 26, 13, 26, 27, 54, 12, 24, 25, 50, 25, 50, 51, 102,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		4, 8, 9, 18, 9, 18, 19, 38, 8, 16, 17, 34, 17, 34, 35, 70,
		5, 10, 11, 22, 11, 22, 23, 46, 10, 20, 21, 42, 21, 42, 43, 86,
		10, 20, 21, 42, 21, 42, 43, 86, 20, 40, 41, 82, 41, 82, 83, 166,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		6, 13, 13, 27, 13, 27, 27, 55, 12, 25, 25, 51, 25, 51, 51, 103,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		4, 9, 9, 19, 9, 19, 19, 39, 8, 17, 17, 35, 17, 35, 35, 71,
		5, 11, 11, 23, 11, 23, 23, 47, 10, 21, 21, 43, 21, 43, 43, 87,
		10, 21, 21, 43, 21, 43, 43, 87, 20, 41, 41, 83, 41, 83, 83, 167,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		6, 12, 12, 24, 12, 24, 24, 48, 13, 26, 26, 52, 26, 52, 52, 104,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		4, 8, 8, 16, 8, 16, 16, 32, 9, 18, 18, 36, 18, 36, 36, 72,
		5, 10, 10, 20, 10, 20, 20, 40, 11, 22, 22, 44, 22, 44, 44, 88,
		10, 20, 20, 40, 20, 40, 40, 80, 21, 42, 42, 84, 42, 84, 84, 168,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		6, 13, 12, 25, 12, 25, 24, 49, 13, 27, 26, 53, 26, 53, 52, 105,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		4, 9, 8, 17, 8, 17, 16, 33, 9, 19, 18, 37, 18, 37, 36, 73,
		5, 11, 10, 21, 10, 21, 20, 41, 11, 23, 22, 45, 22, 45, 44, 89,
		10, 21, 20, 41, 20, 41, 40, 81, 21, 43, 42, 85, 42, 85, 84, 169,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		6, 12, 13, 26, 12, 24, 25, 50, 13, 26, 27, 54, 26, 52, 53, 106,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		4, 8, 9, 18, 8, 16, 17, 34, 9, 18, 19, 38, 18, 36, 37, 74,
		5, 10, 11, 22, 10, 20, 21, 42, 11, 22, 23, 46, 22, 44, 45, 90,
		10, 20, 21, 42, 20, 40, 41, 82, 21, 42, 43, 86, 42, 84, 85, 170,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		6, 13, 13, 27, 12, 25, 25, 51, 13, 27, 27, 55, 26, 53, 53, 107,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		4, 9, 9, 19, 8, 17, 17, 35, 9, 19, 19, 39, 18, 37, 37, 75,
		5, 11, 11, 23, 10, 21, 21, 43, 11, 23, 23, 47, 22, 45, 45, 91,
		10, 21, 21, 43, 20, 41, 41, 83, 21, 43, 43, 87, 42, 85, 85, 171,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		6, 12, 12, 24, 13, 26, 26, 52, 13, 26, 26, 52, 27, 54, 54, 108,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		4, 8, 8, 16, 9, 18, 18, 36, 9, 18, 18, 36, 19, 38, 38, 76,
		5, 10, 10, 20, 11, 22, 22, 44, 11, 22, 22, 44, 23, 46, 46, 92,
		10, 20, 20, 40, 21, 42, 42, 84, 21, 42, 42, 84, 43, 86, 86, 172,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		6, 13, 12, 25, 13, 27, 26, 53, 13, 27, 26, 53, 27, 55, 54, 109,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		4, 9, 8, 17, 9, 19, 18, 37, 9, 19, 18, 37, 19, 39, 38, 77,
		5, 11, 10, 21, 11, 23, 22, 45, 11, 23, 22, 45, 23, 47, 46, 93,
		10, 21, 20, 41, 21, 43, 42, 85, 21, 43, 42, 85, 43, 87, 86, 173,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		6, 12, 13, 26, 13, 26, 27, 54, 13, 26, 27, 54, 27, 54, 55, 110,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		4, 8, 9, 18, 9, 18, 19, 38, 9, 18, 19, 38, 19, 38, 39, 78,
		5, 10, 11, 22, 11, 22, 23, 46, 11, 22, 23, 46, 23, 46, 47, 94,
		10, 20, 21, 42, 21, 42, 43, 86, 21, 42, 43, 86, 43, 86, 87, 174,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		6, 13, 13, 27, 13, 27, 27, 55, 13, 27, 27, 55, 27, 55, 55, 111,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		4, 9, 9, 19, 9, 19, 19, 39, 9, 19, 19, 39, 19, 39, 39, 79,
		5, 11, 11, 23, 11, 23, 23, 47, 11, 23, 23, 47, 23, 47, 47, 95,
		10, 21, 21, 43, 21, 43, 43, 87, 21, 43, 43, 87, 43, 87, 87, 175,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		7, 14, 14, 28, 14, 28, 28, 56, 14, 28, 28, 56, 28, 56, 56, 112,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		5, 10, 10, 20, 10, 20, 20, 40, 10, 20, 20, 40, 20, 40, 40, 80,
		5, 10, 10, 20, 10, 20, 20, 40, 10, 20, 20, 40, 20, 40, 40, 80,
		11, 22, 22, 44, 22, 44, 44, 88, 22, 44, 44, 88, 44, 88, 88, 176,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		7, 15, 14, 29, 14, 29, 28, 57, 14, 29, 28, 57, 28, 57, 56, 113,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		5, 11, 10, 21, 10, 21, 20, 41, 10, 21, 20, 41, 20, 41, 40, 81,
		5, 11, 10, 21, 10, 21, 20, 41, 10, 21, 20, 41, 20, 41, 40, 81,
		11, 23, 22, 45, 22, 45, 44, 89, 22, 45, 44, 89, 44, 89, 88, 177,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		7, 14, 15, 30, 14, 28, 29, 58, 14, 28, 29, 58, 28, 56, 57, 114,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		5, 10, 11, 22, 10, 20, 21, 42, 10, 20, 21, 42, 20, 40, 41, 82,
		5, 10, 11, 22, 10, 20, 21, 42, 10, 20, 21, 42, 20, 40, 41, 82,
		11, 22, 23, 46, 22, 44, 45, 90, 22, 44, 45, 90, 44, 88, 89, 178,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		7, 15, 15, 31, 14, 29, 29, 59, 14, 29, 29, 59, 28, 57, 57, 115,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		5, 11, 11, 23, 10, 21, 21, 43, 10, 21, 21, 43, 20, 41, 41, 83,
		5, 11, 11, 23, 10, 21, 21, 43, 10, 21, 21, 43, 20, 41, 41, 83,
		11, 23, 23, 47, 22, 45, 45, 91, 22, 45, 45, 91, 44, 89, 89, 179,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		7, 14, 14, 28, 15, 30, 30, 60, 14, 28, 28, 56, 29, 58, 58, 116,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		5, 10, 10, 20, 11, 22, 22, 44, 10, 20, 20, 40, 21, 42, 42, 84,
		5, 10, 10, 20, 11, 22, 22, 44, 10, 20, 20, 40, 21, 42, 42, 84,
		11, 22, 22, 44, 23, 46, 46, 92, 22, 44, 44, 88, 45, 90, 90, 180,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		7, 15, 14, 29, 15, 31, 30, 61, 14, 29, 28, 57, 29, 59, 58, 117,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		5, 11, 10, 21, 11, 23, 22, 45, 10, 21, 20, 41, 21, 43, 42, 85,
		5, 11, 10, 21, 11, 23, 22, 45, 10, 21, 20, 41, 21, 43, 42, 85,
		11, 23, 22, 45, 23, 47, 46, 93, 22, 45, 44, 89, 45, 91, 90, 181,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		7, 14, 15, 30, 15, 30, 31, 62, 14, 28, 29, 58, 29, 58, 59, 118,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		5, 10, 11, 22, 11, 22, 23, 46, 10, 20, 21, 42, 21, 42, 43, 86,
		5, 10, 11, 22, 11, 22, 23, 46, 10, 20, 21, 42, 21, 42, 43, 86,
		11, 22, 23, 46, 23, 46, 47, 94, 22, 44, 45, 90, 45, 90, 91, 182,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		7, 15, 15, 31, 15, 31, 31, 63, 14, 29, 29, 59, 29, 59, 59, 119,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		5, 11, 11, 23, 11, 23, 23, 47, 10, 21, 21, 43, 21, 43, 43, 87,
		5, 11, 11, 23, 11, 23, 23, 47, 10, 21, 21, 43, 21, 43, 43, 87,
		11, 23, 23, 47, 23, 47, 47, 95, 22, 45, 45, 91, 45, 91, 91, 183,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		7, 14, 14, 28, 14, 28, 28, 56, 15, 30, 30, 60, 30, 60, 60, 120,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		5, 10, 10, 20, 10, 20, 20, 40, 11, 22, 22, 44, 22, 44, 44, 88,
		5, 10, 10, 20, 10, 20, 20, 40, 11, 22, 22, 44, 22, 44, 44, 88,
		11, 22, 22, 44, 22, 44, 44, 88, 23, 46, 46, 92, 46, 92, 92, 184,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		7, 15, 14, 29, 14, 29, 28, 57, 15, 31, 30, 61, 30, 61, 60, 121,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		5, 11, 10, 21, 10, 21, 20, 41, 11, 23, 22, 45, 22, 45, 44, 89,
		5, 11, 10, 21, 10, 21, 20, 41, 11, 23, 22, 45, 22, 45, 44, 89,
		11, 23, 22, 45, 22, 45, 44, 89, 23, 47, 46, 93, 46, 93, 92, 185,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		7, 14, 15, 30, 14, 28, 29, 58, 15, 30, 31, 62, 30, 60, 61, 122,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		5, 10, 11, 22, 10, 20, 21, 42, 11, 22, 23, 46, 22, 44, 45, 90,
		5, 10, 11, 22, 10, 20, 21, 42, 11, 22, 23, 46, 22, 44, 45, 90,
		11, 22, 23, 46, 22, 44, 45, 90, 23, 46, 47, 94, 46, 92, 93, 186,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		7, 15, 15, 31, 14, 29, 29, 59, 15, 31, 31, 63, 30, 61, 61, 123,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		5, 11, 11, 23, 10, 21, 21, 43, 11, 23, 23, 47, 22, 45, 45, 91,
		5, 11, 11, 23, 10, 21, 21, 43, 11, 23, 23, 47, 22, 45, 45, 91,
		11, 23, 23, 47, 22, 45, 45, 91, 23, 47, 47, 95, 46, 93, 93, 187,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		7, 14, 14, 28, 15, 30, 30, 60, 15, 30, 30, 60, 31, 62, 62, 124,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		5, 10, 10, 20, 11, 22, 22, 44, 11, 22, 22, 44, 23, 46, 46, 92,
		5, 10, 10, 20, 11, 22, 22, 44, 11, 22, 22, 44, 23, 46, 46, 92,
		11, 22, 22, 44, 23, 46, 46, 92, 23, 46, 46, 92, 47, 94, 94, 188,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		7, 15, 14, 29, 15, 31, 30, 61, 15, 31, 30, 61, 31, 63, 62, 125,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		5, 11, 10, 21, 11, 23, 22, 45, 11, 23, 22, 45, 23, 47, 46, 93,
		5, 11, 10, 21, 11, 23, 22, 45, 11, 23, 22, 45, 23, 47, 46, 93,
		11, 23, 22, 45, 23, 47, 46, 93, 23, 47, 46, 93, 47, 95, 94, 189,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		7, 14, 15, 30, 15, 30, 31, 62, 15, 30, 31, 62, 31, 62, 63, 126,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		5, 10, 11, 22, 11, 22, 23, 46, 11, 22, 23, 46, 23, 46, 47, 94,
		5, 10, 11, 22, 11, 22, 23, 46, 11, 22, 23, 46, 23, 46, 47, 94,
		11, 22, 23, 46, 23, 46, 47, 94, 23, 46, 47, 94, 47, 94, 95, 190,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		7, 15, 15, 31, 15, 31, 31, 63, 15, 31, 31, 63, 31, 63, 63, 127,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		5, 11, 11, 23, 11, 23, 23, 47, 11, 23, 23, 47, 23, 47, 47, 95,
		5, 11, 11, 23, 11, 23, 23, 47, 11, 23, 23, 47, 23, 47, 47, 95,
		11, 23, 23, 47, 23, 47, 47, 95, 23, 47, 47, 95, 47, 95, 95, 191,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		4, 8, 8, 16, 8, 16, 16, 32, 8, 16, 16, 32, 16, 32, 32, 64,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		4, 8, 8, 16, 8, 16, 16, 32, 8, 16, 16, 32, 16, 32, 32, 64,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		6, 12, 12, 24, 12, 24, 24, 48, 12, 24, 24, 48, 24, 48, 48, 96,
		6, 12, 12, 24, 12, 24, 24, 48, 12, 24, 24, 48, 24, 48, 48, 96,
		12, 24, 24, 48, 24, 48, 48, 96, 24, 48, 48, 96, 48, 96, 96, 192,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		4, 9, 8, 17, 8, 17, 16, 33, 8, 17, 16, 33, 16, 33, 32, 65,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		4, 9, 8, 17, 8, 17, 16, 33, 8, 17, 16, 33, 16, 33, 32, 65,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		6, 13, 12, 25, 12, 25, 24, 49, 12, 25, 24, 49, 24, 49, 48, 97,
		6, 13, 12, 25, 12, 25, 24, 49, 12, 25, 24, 49, 24, 49, 48, 97,
		12, 25, 24, 49, 24, 49, 48, 97, 24, 49, 48, 97, 48, 97, 96, 193,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		4, 8, 9, 18, 8, 16, 17, 34, 8, 16, 17, 34, 16, 32, 33, 66,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		4, 8, 9, 18, 8, 16, 17, 34, 8, 16, 17, 34, 16, 32, 33, 66,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		6, 12, 13, 26, 12, 24, 25, 50, 12, 24, 25, 50, 24, 48, 49, 98,
		6, 12, 13, 26, 12, 24, 25, 50, 12, 24, 25, 50, 24, 48, 49, 98,
		12, 24, 25, 50, 24, 48, 49, 98, 24, 48, 49, 98, 48, 96, 97, 194,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		4, 9, 9, 19, 8, 17, 17, 35, 8, 17, 17, 35, 16, 33, 33, 67,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		4, 9, 9, 19, 8, 17, 17, 35, 8, 17, 17, 35, 16, 33, 33, 67,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		6, 13, 13, 27, 12, 25, 25, 51, 12, 25, 25, 51, 24, 49, 49, 99,
		6, 13, 13, 27, 12, 25, 25, 51, 12, 25, 25, 51, 24, 49, 49, 99,
		12, 25, 25, 51, 24, 49, 49, 99, 24, 49, 49, 99, 48, 97, 97, 195,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		4, 8, 8, 16, 9, 18, 18, 36, 8, 16, 16, 32, 17, 34, 34, 68,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		4, 8, 8, 16, 9, 18, 18, 36, 8, 16, 16, 32, 17, 34, 34, 68,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		6, 12, 12, 24, 13, 26, 26, 52, 12, 24, 24, 48, 25, 50, 50, 100,
		6, 12, 12, 24, 13, 26, 26, 52, 12, 24, 24, 48, 25, 50, 50, 100,
		12, 24, 24, 48, 25, 50, 50, 100, 24, 48, 48, 96, 49, 98, 98, 196,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		4, 9, 8, 17, 9, 19, 18, 37, 8, 17, 16, 33, 17, 35, 34, 69,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		4, 9, 8, 17, 9, 19, 18, 37, 8, 17, 16, 33, 17, 35, 34, 69,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		6, 13, 12, 25, 13, 27, 26, 53, 12, 25, 24, 49, 25, 51, 50, 101,
		6, 13, 12, 25, 13, 27, 26, 53, 12, 25, 24, 49, 25, 51, 50, 101,
		12, 25, 24, 49, 25, 51, 50, 101, 24, 49, 48, 97, 49, 99, 98, 197,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		4, 8, 9, 18, 9, 18, 19, 38, 8, 16, 17, 34, 17, 34, 35, 70,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		4, 8, 9, 18, 9, 18, 19, 38, 8, 16, 17, 34, 17, 34, 35, 70,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		6, 12, 13, 26, 13, 26, 27, 54, 12, 24, 25, 50, 25, 50, 51, 102,
		6, 12, 13, 26, 13, 26, 27, 54, 12, 24, 25, 50, 25, 50, 51, 102,
		12, 24, 25, 50, 25, 50, 51, 102, 24, 48, 49, 98, 49, 98, 99, 198,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		4, 9, 9, 19, 9, 19, 19, 39, 8, 17, 17, 35, 17, 35, 35, 71,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		4, 9, 9, 19, 9, 19, 19, 39, 8, 17, 17, 35, 17, 35, 35, 71,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		6, 13, 13, 27, 13, 27, 27, 55, 12, 25, 25, 51, 25, 51, 51, 103,
		6, 13, 13, 27, 13, 27, 27, 55, 12, 25, 25, 51, 25, 51, 51, 103,
		12, 25, 25, 51, 25, 51, 51, 103, 24, 49, 49, 99, 49, 99, 99, 199,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		4, 8, 8, 16, 8, 16, 16, 32, 9, 18, 18, 36, 18, 36, 36, 72,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		4, 8, 8, 16, 8, 16, 16, 32, 9, 18, 18, 36, 18, 36, 36, 72,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		6, 12, 12, 24, 12, 24, 24, 48, 13, 26, 26, 52, 26, 52, 52, 104,
		6, 12, 12, 24, 12, 24, 24, 48, 13, 26, 26, 52, 26, 52, 52, 104,
		12, 24, 24, 48, 24, 48, 48, 96, 25, 50, 50, 100, 50, 100, 100, 200,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		4, 9, 8, 17, 8, 17, 16, 33, 9, 19, 18, 37, 18, 37, 36, 73,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		4, 9, 8, 17, 8, 17, 16, 33, 9, 19, 18, 37, 18, 37, 36, 73,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		6, 13, 12, 25, 12, 25, 24, 49, 13, 27, 26, 53, 26, 53, 52, 105,
		6, 13, 12, 25, 12, 25, 24, 49, 13, 27, 26, 53, 26, 53, 52, 105,
		12, 25, 24, 49, 24, 49, 48, 97, 25, 51, 50, 101, 50, 101, 100, 201,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		4, 8, 9, 18, 8, 16, 17, 34, 9, 18, 19, 38, 18, 36, 37, 74,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		4, 8, 9, 18, 8, 16, 17, 34, 9, 18, 19, 38, 18, 36, 37, 74,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		6, 12, 13, 26, 12, 24, 25, 50, 13, 26, 27, 54, 26, 52, 53, 106,
		6, 12, 13, 26, 12, 24, 25, 50, 13, 26, 27, 54, 26, 52, 53, 106,
		12, 24, 25, 50, 24, 48, 49, 98, 25, 50, 51, 102, 50, 100, 101, 202,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		4, 9, 9, 19, 8, 17, 17, 35, 9, 19, 19, 39, 18, 37, 37, 75,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		4, 9, 9, 19, 8, 17, 17, 35, 9, 19, 19, 39, 18, 37, 37, 75,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		6, 13, 13, 27, 12, 25, 25, 51, 13, 27, 27, 55, 26, 53, 53, 107,
		6, 13, 13, 27, 12, 25, 25, 51, 13, 27, 27, 55, 26, 53, 53, 107,
		12, 25, 25, 51, 24, 49, 49, 99, 25, 51, 51, 103, 50, 101, 101, 203,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		4, 8, 8, 16, 9, 18, 18, 36, 9, 18, 18, 36, 19, 38, 38, 76,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		4, 8, 8, 16, 9, 18, 18, 36, 9, 18, 18, 36, 19, 38, 38, 76,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		6, 12, 12, 24, 13, 26, 26, 52, 13, 26, 26, 52, 27, 54, 54, 108,
		6, 12, 12, 24, 13, 26, 26, 52, 13, 26, 26, 52, 27, 54, 54, 108,
		12, 24, 24, 48, 25, 50, 50, 100, 25, 50, 50, 100, 51, 102, 102, 204,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		4, 9, 8, 17, 9, 19, 18, 37, 9, 19, 18, 37, 19, 39, 38, 77,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		4, 9, 8, 17, 9, 19, 18, 37, 9, 19, 18, 37, 19, 39, 38, 77,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		6, 13, 12, 25, 13, 27, 26, 53, 13, 27, 26, 53, 27, 55, 54, 109,
		6, 13, 12, 25, 13, 27, 26, 53, 13, 27, 26, 53, 27, 55, 54, 109,
		12, 25, 24, 49, 25, 51, 50, 101, 25, 51, 50, 101, 51, 103, 102, 205,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		4, 8, 9, 18, 9, 18, 19, 38, 9, 18, 19, 38, 19, 38, 39, 78,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		4, 8, 9, 18, 9, 18, 19, 38, 9, 18, 19, 38, 19, 38, 39, 78,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		6, 12, 13, 26, 13, 26, 27, 54, 13, 26, 27, 54, 27, 54, 55, 110,
		6, 12, 13, 26, 13, 26, 27, 54, 13, 26, 27, 54, 27, 54, 55, 110,
		12, 24, 25, 50, 25, 50, 51, 102, 25, 50, 51, 102, 51, 102, 103, 206,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		4, 9, 9, 19, 9, 19, 19, 39, 9, 19, 19, 39, 19, 39, 39, 79,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		4, 9, 9, 19, 9, 19, 19, 39, 9, 19, 19, 39, 19, 39, 39, 79,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		6, 13, 13, 27, 13, 27, 27, 55, 13, 27, 27, 55, 27, 55, 55, 111,
		6, 13, 13, 27, 13, 27, 27, 55, 13, 27, 27, 55, 27, 55, 55, 111,
		12, 25, 25, 51, 25, 51, 51, 103, 25, 51, 51, 103, 51, 103, 103, 207,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		5, 10, 10, 20, 10, 20, 20, 40, 10, 20, 20, 40, 20, 40, 40, 80,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		5, 10, 10, 20, 10, 20, 20, 40, 10, 20, 20, 40, 20, 40, 40, 80,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		7, 14, 14, 28, 14, 28, 28, 56, 14, 28, 28, 56, 28, 56, 56, 112,
		6, 12, 12, 24, 12, 24, 24, 48, 12, 24, 24, 48, 24, 48, 48, 96,
		13, 26, 26, 52, 26, 52, 52, 104, 26, 52, 52, 104, 52, 104, 104, 208,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		5, 11, 10, 21, 10, 21, 20, 41, 10, 21, 20, 41, 20, 41, 40, 81,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		5, 11, 10, 21, 10, 21, 20, 41, 10, 21, 20, 41, 20, 41, 40, 81,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		7, 15, 14, 29, 14, 29, 28, 57, 14, 29, 28, 57, 28, 57, 56, 113,
		6, 13, 12, 25, 12, 25, 24, 49, 12, 25, 24, 49, 24, 49, 48, 97,
		13, 27, 26, 53, 26, 53, 52, 105, 26, 53, 52, 105, 52, 105, 104, 209,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		5, 10, 11, 22, 10, 20, 21, 42, 10, 20, 21, 42, 20, 40, 41, 82,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		5, 10, 11, 22, 10, 20, 21, 42, 10, 20, 21, 42, 20, 40, 41, 82,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		7, 14, 15, 30, 14, 28, 29, 58, 14, 28, 29, 58, 28, 56, 57, 114,
		6, 12, 13, 26, 12, 24, 25, 50, 12, 24, 25, 50, 24, 48, 49, 98,
		13, 26, 27, 54, 26, 52, 53, 106, 26, 52, 53, 106, 52, 104, 105, 210,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		5, 11, 11, 23, 10, 21, 21, 43, 10, 21, 21, 43, 20, 41, 41, 83,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		5, 11, 11, 23, 10, 21, 21, 43, 10, 21, 21, 43, 20, 41, 41, 83,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		7, 15, 15, 31, 14, 29, 29, 59, 14, 29, 29, 59, 28, 57, 57, 115,
		6, 13, 13, 27, 12, 25, 25, 51, 12, 25, 25, 51, 24, 49, 49, 99,
		13, 27, 27, 55, 26, 53, 53, 107, 26, 53, 53, 107, 52, 105, 105, 211,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		5, 10, 10, 20, 11, 22, 22, 44, 10, 20, 20, 40, 21, 42, 42, 84,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		5, 10, 10, 20, 11, 22, 22, 44, 10, 20, 20, 40, 21, 42, 42, 84,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		7, 14, 14, 28, 15, 30, 30, 60, 14, 28, 28, 56, 29, 58, 58, 116,
		6, 12, 12, 24, 13, 26, 26, 52, 12, 24, 24, 48, 25, 50, 50, 100,
		13, 26, 26, 52, 27, 54, 54, 108, 26, 52, 52, 104, 53, 106, 106, 212,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		5, 11, 10, 21, 11, 23, 22, 45, 10, 21, 20, 41, 21, 43, 42, 85,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		5, 11, 10, 21, 11, 23, 22, 45, 10, 21, 20, 41, 21, 43, 42, 85,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		7, 15, 14, 29, 15, 31, 30, 61, 14, 29, 28, 57, 29, 59, 58, 117,
		6, 13, 12, 25, 13, 27, 26, 53, 12, 25, 24, 49, 25, 51, 50, 101,
		13, 27, 26, 53, 27, 55, 54, 109, 26, 53, 52, 105, 53, 107, 106, 213,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		5, 10, 11, 22, 11, 22, 23, 46, 10, 20, 21, 42, 21, 42, 43, 86,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		5, 10, 11, 22, 11, 22, 23, 46, 10, 20, 21, 42, 21, 42, 43, 86,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		7, 14, 15, 30, 15, 30, 31, 62, 14, 28, 29, 58, 29, 58, 59, 118,
		6, 12, 13, 26, 13, 26, 27, 54, 12, 24, 25, 50, 25, 50, 51, 102,
		13, 26, 27, 54, 27, 54, 55, 110, 26, 52, 53, 106, 53, 106, 107, 214,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		5, 11, 11, 23, 11, 23, 23, 47, 10, 21, 21, 43, 21, 43, 43, 87,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		5, 11, 11, 23, 11, 23, 23, 47, 10, 21, 21, 43, 21, 43, 43, 87,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		7, 15, 15, 31, 15, 31, 31, 63, 14, 29, 29, 59, 29, 59, 59, 119,
		6, 13, 13, 27, 13, 27, 27, 55, 12, 25, 25, 51, 25, 51, 51, 103,
		13, 27, 27, 55, 27, 55, 55, 111, 26, 53, 53, 107, 53, 107, 107, 215,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		5, 10, 10, 20, 10, 20, 20, 40, 11, 22, 22, 44, 22, 44, 44, 88,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		5, 10, 10, 20, 10, 20, 20, 40, 11, 22, 22, 44, 22, 44, 44, 88,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		7, 14, 14, 28, 14, 28, 28, 56, 15, 30, 30, 60, 30, 60, 60, 120,
		6, 12, 12, 24, 12, 24, 24, 48, 13, 26, 26, 52, 26, 52, 52, 104,
		13, 26, 26, 52, 26, 52, 52, 104, 27, 54, 54, 108, 54, 108, 108, 216,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		5, 11, 10, 21, 10, 21, 20, 41, 11, 23, 22, 45, 22, 45, 44, 89,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		5, 11, 10, 21, 10, 21, 20, 41, 11, 23, 22, 45, 22, 45, 44, 89,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		7, 15, 14, 29, 14, 29, 28, 57, 15, 31, 30, 61, 30, 61, 60, 121,
		6, 13, 12, 25, 12, 25, 24, 49, 13, 27, 26, 53, 26, 53, 52, 105,
		13, 27, 26, 53, 26, 53, 52, 105, 27, 55, 54, 109, 54, 109, 108, 217,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		5, 10, 11, 22, 10, 20, 21, 42, 11, 22, 23, 46, 22, 44, 45, 90,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		5, 10, 11, 22, 10, 20, 21, 42, 11, 22, 23, 46, 22, 44, 45, 90,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		7, 14, 15, 30, 14, 28, 29, 58, 15, 30, 31, 62, 30, 60, 61, 122,
		6, 12, 13, 26, 12, 24, 25, 50, 13, 26, 27, 54, 26, 52, 53, 106,
		13, 26, 27, 54, 26, 52, 53, 106, 27, 54, 55, 110, 54, 108, 109, 218,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		5, 11, 11, 23, 10, 21, 21, 43, 11, 23, 23, 47, 22, 45, 45, 91,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		5, 11, 11, 23, 10, 21, 21, 43, 11, 23, 23, 47, 22, 45, 45, 91,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		7, 15, 15, 31, 14, 29, 29, 59, 15, 31, 31, 63, 30, 61, 61, 123,
		6, 13, 13, 27, 12, 25, 25, 51, 13, 27, 27, 55, 26, 53, 53, 107,
		13, 27, 27, 55, 26, 53, 53, 107, 27, 55, 55, 111, 54, 109, 109, 219,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		5, 10, 10, 20, 11, 22, 22, 44, 11, 22, 22, 44, 23, 46, 46, 92,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		5, 10, 10, 20, 11, 22, 22, 44, 11, 22, 22, 44, 23, 46, 46, 92,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		7, 14, 14, 28, 15, 30, 30, 60, 15, 30, 30, 60, 31, 62, 62, 124,
		6, 12, 12, 24, 13, 26, 26, 52, 13, 26, 26, 52, 27, 54, 54, 108,
		13, 26, 26, 52, 27, 54, 54, 108, 27, 54, 54, 108, 55, 110, 110, 220,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		5, 11, 10, 21, 11, 23, 22, 45, 11, 23, 22, 45, 23, 47, 46, 93,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		5, 11, 10, 21, 11, 23, 22, 45, 11, 23, 22, 45, 23, 47, 46, 93,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		7, 15, 14, 29, 15, 31, 30, 61, 15, 31, 30, 61, 31, 63, 62, 125,
		6, 13, 12, 25, 13, 27, 26, 53, 13, 27, 26, 53, 27, 55, 54, 109,
		13, 27, 26, 53, 27, 55, 54, 109, 27, 55, 54, 109, 55, 111, 110, 221,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		5, 10, 11, 22, 11, 22, 23, 46, 11, 22, 23, 46, 23, 46, 47, 94,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		5, 10, 11, 22, 11, 22, 23, 46, 11, 22, 23, 46, 23, 46, 47, 94,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		7, 14, 15, 30, 15, 30, 31, 62, 15, 30, 31, 62, 31, 62, 63, 126,
		6, 12, 13, 26, 13, 26, 27, 54, 13, 26, 27, 54, 27, 54, 55, 110,
		13, 26, 27, 54, 27, 54, 55, 110, 27, 54, 55, 110, 55, 110, 111, 222,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		5, 11, 11, 23, 11, 23, 23, 47, 11, 23, 23, 47, 23, 47, 47, 95,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		5, 11, 11, 23, 11, 23, 23, 47, 11, 23, 23, 47, 23, 47, 47, 95,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		7, 15, 15, 31, 15, 31, 31, 63, 15, 31, 31, 63, 31, 63, 63, 127,
		6, 13, 13, 27, 13, 27, 27, 55, 13, 27, 27, 55, 27, 55, 55, 111,
		13, 27, 27, 55, 27, 55, 55, 111, 27, 55, 55, 111, 55, 111, 111, 223,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		6, 12, 12, 24, 12, 24, 24, 48, 12, 24, 24, 48, 24, 48, 48, 96,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		2, 4, 4, 8, 4, 8, 8, 16, 4, 8, 8, 16, 8, 16, 16, 32,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		6, 12, 12, 24, 12, 24, 24, 48, 12, 24, 24, 48, 24, 48, 48, 96,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		6, 12, 12, 24, 12, 24, 24, 48, 12, 24, 24, 48, 24, 48, 48, 96,
		7, 14, 14, 28, 14, 28, 28, 56, 14, 28, 28, 56, 28, 56, 56, 112,
		14, 28, 28, 56, 28, 56, 56, 112, 28, 56, 56, 112, 56, 112, 112, 224,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		6, 13, 12, 25, 12, 25, 24, 49, 12, 25, 24, 49, 24, 49, 48, 97,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		2, 5, 4, 9, 4, 9, 8, 17, 4, 9, 8, 17, 8, 17, 16, 33,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		6, 13, 12, 25, 12, 25, 24, 49, 12, 25, 24, 49, 24, 49, 48, 97,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		6, 13, 12, 25, 12, 25, 24, 49, 12, 25, 24, 49, 24, 49, 48, 97,
		7, 15, 14, 29, 14, 29, 28, 57, 14, 29, 28, 57, 28, 57, 56, 113,
		14, 29, 28, 57, 28, 57, 56, 113, 28, 57, 56, 113, 56, 113, 112, 225,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		6, 12, 13, 26, 12, 24, 25, 50, 12, 24, 25, 50, 24, 48, 49, 98,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		2, 4, 5, 10, 4, 8, 9, 18, 4, 8, 9, 18, 8, 16, 17, 34,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		6, 12, 13, 26, 12, 24, 25, 50, 12, 24, 25, 50, 24, 48, 49, 98,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		6, 12, 13, 26, 12, 24, 25, 50, 12, 24, 25, 50, 24, 48, 49, 98,
		7, 14, 15, 30, 14, 28, 29, 58, 14, 28, 29, 58, 28, 56, 57, 114,
		14, 28, 29, 58, 28, 56, 57, 114, 28, 56, 57, 114, 56, 112, 113, 226,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		6, 13, 13, 27, 12, 25, 25, 51, 12, 25, 25, 51, 24, 49, 49, 99,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		2, 5, 5, 11, 4, 9, 9, 19, 4, 9, 9, 19, 8, 17, 17, 35,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		6, 13, 13, 27, 12, 25, 25, 51, 12, 25, 25, 51, 24, 49, 49, 99,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		6, 13, 13, 27, 12, 25, 25, 51, 12, 25, 25, 51, 24, 49, 49, 99,
		7, 15, 15, 31, 14, 29, 29, 59, 14, 29, 29, 59, 28, 57, 57, 115,
		14, 29, 29, 59, 28, 57, 57, 115, 28, 57, 57, 115, 56, 113, 113, 227,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		6, 12, 12, 24, 13, 26, 26, 52, 12, 24, 24, 48, 25, 50, 50, 100,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		2, 4, 4, 8, 5, 10, 10, 20, 4, 8, 8, 16, 9, 18, 18, 36,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		6, 12, 12, 24, 13, 26, 26, 52, 12, 24, 24, 48, 25, 50, 50, 100,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		6, 12, 12, 24, 13, 26, 26, 52, 12, 24, 24, 48, 25, 50, 50, 100,
		7, 14, 14, 28, 15, 30, 30, 60, 14, 28, 28, 56, 29, 58, 58, 116,
		14, 28, 28, 56, 29, 58, 58, 116, 28, 56, 56, 112, 57, 114, 114, 228,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		6, 13, 12, 25, 13, 27, 26, 53, 12, 25, 24, 49, 25, 51, 50, 101,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		2, 5, 4, 9, 5, 11, 10, 21, 4, 9, 8, 17, 9, 19, 18, 37,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		6, 13, 12, 25, 13, 27, 26, 53, 12, 25, 24, 49, 25, 51, 50, 101,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		6, 13, 12, 25, 13, 27, 26, 53, 12, 25, 24, 49, 25, 51, 50, 101,
		7, 15, 14, 29, 15, 31, 30, 61, 14, 29, 28, 57, 29, 59, 58, 117,
		14, 29, 28, 57, 29, 59, 58, 117, 28, 57, 56, 113, 57, 115, 114, 229,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		6, 12, 13, 26, 13, 26, 27, 54, 12, 24, 25, 50, 25, 50, 51, 102,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		2, 4, 5, 10, 5, 10, 11, 22, 4, 8, 9, 18, 9, 18, 19, 38,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		6, 12, 13, 26, 13, 26, 27, 54, 12, 24, 25, 50, 25, 50, 51, 102,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		6, 12, 13, 26, 13, 26, 27, 54, 12, 24, 25, 50, 25, 50, 51, 102,
		7, 14, 15, 30, 15, 30, 31, 62, 14, 28, 29, 58, 29, 58, 59, 118,
		14, 28, 29, 58, 29, 58, 59, 118, 28, 56, 57, 114, 57, 114, 115, 230,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		6, 13, 13, 27, 13, 27, 27, 55, 12, 25, 25, 51, 25, 51, 51, 103,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		2, 5, 5, 11, 5, 11, 11, 23, 4, 9, 9, 19, 9, 19, 19, 39,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		6, 13, 13, 27, 13, 27, 27, 55, 12, 25, 25, 51, 25, 51, 51, 103,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		6, 13, 13, 27, 13, 27, 27, 55, 12, 25, 25, 51, 25, 51, 51, 103,
		7, 15, 15, 31, 15, 31, 31, 63, 14, 29, 29, 59, 29, 59, 59, 119,
		14, 29, 29, 59, 29, 59, 59, 119, 28, 57, 57, 115, 57, 115, 115, 231,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		6, 12, 12, 24, 12, 24, 24, 48, 13, 26, 26, 52, 26, 52, 52, 104,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		2, 4, 4, 8, 4, 8, 8, 16, 5, 10, 10, 20, 10, 20, 20, 40,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		6, 12, 12, 24, 12, 24, 24, 48, 13, 26, 26, 52, 26, 52, 52, 104,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		6, 12, 12, 24, 12, 24, 24, 48, 13, 26, 26, 52, 26, 52, 52, 104,
		7, 14, 14, 28, 14, 28, 28, 56, 15, 30, 30, 60, 30, 60, 60, 120,
		14, 28, 28, 56, 28, 56, 56, 112, 29, 58, 58, 116, 58, 116, 116, 232,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		6, 13, 12, 25, 12, 25, 24, 49, 13, 27, 26, 53, 26, 53, 52, 105,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		2, 5, 4, 9, 4, 9, 8, 17, 5, 11, 10, 21, 10, 21, 20, 41,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		6, 13, 12, 25, 12, 25, 24, 49, 13, 27, 26, 53, 26, 53, 52, 105,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		6, 13, 12, 25, 12, 25, 24, 49, 13, 27, 26, 53, 26, 53, 52, 105,
		7, 15, 14, 29, 14, 29, 28, 57, 15, 31, 30, 61, 30, 61, 60, 121,
		14, 29, 28, 57, 28, 57, 56, 113, 29, 59, 58, 117, 58, 117, 116, 233,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		6, 12, 13, 26, 12, 24, 25, 50, 13, 26, 27, 54, 26, 52, 53, 106,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		2, 4, 5, 10, 4, 8, 9, 18, 5, 10, 11, 22, 10, 20, 21, 42,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		6, 12, 13, 26, 12, 24, 25, 50, 13, 26, 27, 54, 26, 52, 53, 106,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		6, 12, 13, 26, 12, 24, 25, 50, 13, 26, 27, 54, 26, 52, 53, 106,
		7, 14, 15, 30, 14, 28, 29, 58, 15, 30, 31, 62, 30, 60, 61, 122,
		14, 28, 29, 58, 28, 56, 57, 114, 29, 58, 59, 118, 58, 116, 117, 234,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		6, 13, 13, 27, 12, 25, 25, 51, 13, 27, 27, 55, 26, 53, 53, 107,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		2, 5, 5, 11, 4, 9, 9, 19, 5, 11, 11, 23, 10, 21, 21, 43,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		6, 13, 13, 27, 12, 25, 25, 51, 13, 27, 27, 55, 26, 53, 53, 107,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		6, 13, 13, 27, 12, 25, 25, 51, 13, 27, 27, 55, 26, 53, 53, 107,
		7, 15, 15, 31, 14, 29, 29, 59, 15, 31, 31, 63, 30, 61, 61, 123,
		14, 29, 29, 59, 28, 57, 57, 115, 29, 59, 59, 119, 58, 117, 117, 235,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		6, 12, 12, 24, 13, 26, 26, 52, 13, 26, 26, 52, 27, 54, 54, 108,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		2, 4, 4, 8, 5, 10, 10, 20, 5, 10, 10, 20, 11, 22, 22, 44,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		6, 12, 12, 24, 13, 26, 26, 52, 13, 26, 26, 52, 27, 54, 54, 108,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		6, 12, 12, 24, 13, 26, 26, 52, 13, 26, 26, 52, 27, 54, 54, 108,
		7, 14, 14, 28, 15, 30, 30, 60, 15, 30, 30, 60, 31, 62, 62, 124,
		14, 28, 28, 56, 29, 58, 58, 116, 29, 58, 58, 116, 59, 118, 118, 236,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		6, 13, 12, 25, 13, 27, 26, 53, 13, 27, 26, 53, 27, 55, 54, 109,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		2, 5, 4, 9, 5, 11, 10, 21, 5, 11, 10, 21, 11, 23, 22, 45,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		6, 13, 12, 25, 13, 27, 26, 53, 13, 27, 26, 53, 27, 55, 54, 109,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		6, 13, 12, 25, 13, 27, 26, 53, 13, 27, 26, 53, 27, 55, 54, 109,
		7, 15, 14, 29, 15, 31, 30, 61, 15, 31, 30, 61, 31, 63, 62, 125,
		14, 29, 28, 57, 29, 59, 58, 117, 29, 59, 58, 117, 59, 119, 118, 237,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		6, 12, 13, 26, 13, 26, 27, 54, 13, 26, 27, 54, 27, 54, 55, 110,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		2, 4, 5, 10, 5, 10, 11, 22, 5, 10, 11, 22, 11, 22, 23, 46,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		6, 12, 13, 26, 13, 26, 27, 54, 13, 26, 27, 54, 27, 54, 55, 110,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		6, 12, 13, 26, 13, 26, 27, 54, 13, 26, 27, 54, 27, 54, 55, 110,
		7, 14, 15, 30, 15, 30, 31, 62, 15, 30, 31, 62, 31, 62, 63, 126,
		14, 28, 29, 58, 29, 58, 59, 118, 29, 58, 59, 118, 59, 118, 119, 238,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		6, 13, 13, 27, 13, 27, 27, 55, 13, 27, 27, 55, 27, 55, 55, 111,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		2, 5, 5, 11, 5, 11, 11, 23, 5, 11, 11, 23, 11, 23, 23, 47,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		6, 13, 13, 27, 13, 27, 27, 55, 13, 27, 27, 55, 27, 55, 55, 111,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		6, 13, 13, 27, 13, 27, 27, 55, 13, 27, 27, 55, 27, 55, 55, 111,
		7, 15, 15, 31, 15, 31, 31, 63, 15, 31, 31, 63, 31, 63, 63, 127,
		14, 29, 29, 59, 29, 59, 59, 119, 29, 59, 59, 119, 59, 119, 119, 239,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		7, 14, 14, 28, 14, 28, 28, 56, 14, 28, 28, 56, 28, 56, 56, 112,
		1, 2, 2, 4, 2, 4, 4, 8, 2, 4, 4, 8, 4, 8, 8, 16,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		7, 14, 14, 28, 14, 28, 28, 56, 14, 28, 28, 56, 28, 56, 56, 112,
		3, 6, 6, 12, 6, 12, 12, 24, 6, 12, 12, 24, 12, 24, 24, 48,
		7, 14, 14, 28, 14, 28, 28, 56, 14, 28, 28, 56, 28, 56, 56, 112,
		7, 14, 14, 28, 14, 28, 28, 56, 14, 28, 28, 56, 28, 56, 56, 112,
		15, 30, 30, 60, 30, 60, 60, 120, 30, 60, 60, 120, 60, 120, 120, 240,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		7, 15, 14, 29, 14, 29, 28, 57, 14, 29, 28, 57, 28, 57, 56, 113,
		1, 3, 2, 5, 2, 5, 4, 9, 2, 5, 4, 9, 4, 9, 8, 17,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		7, 15, 14, 29, 14, 29, 28, 57, 14, 29, 28, 57, 28, 57, 56, 113,
		3, 7, 6, 13, 6, 13, 12, 25, 6, 13, 12, 25, 12, 25, 24, 49,
		7, 15, 14, 29, 14, 29, 28, 57, 14, 29, 28, 57, 28, 57, 56, 113,
		7, 15, 14, 29, 14, 29, 28, 57, 14, 29, 28, 57, 28, 57, 56, 113,
		15, 31, 30, 61, 30, 61, 60, 121, 30, 61, 60, 121, 60, 121, 120, 241,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2, 0, 0, 1, 2,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		7, 14, 15, 30, 14, 28, 29, 58, 14, 28, 29, 58, 28, 56, 57, 114,
		1, 2, 3, 6, 2, 4, 5, 10, 2, 4, 5, 10, 4, 8, 9, 18,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		7, 14, 15, 30, 14, 28, 29, 58, 14, 28, 29, 58, 28, 56, 57, 114,
		3, 6, 7, 14, 6, 12, 13, 26, 6, 12, 13, 26, 12, 24, 25, 50,
		7, 14, 15, 30, 14, 28, 29, 58, 14, 28, 29, 58, 28, 56, 57, 114,
		7, 14, 15, 30, 14, 28, 29, 58, 14, 28, 29, 58, 28, 56, 57, 114,
		15, 30, 31, 62, 30, 60, 61, 122, 30, 60, 61, 122, 60, 120, 121, 242,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3, 0, 1, 1, 3,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		7, 15, 15, 31, 14, 29, 29, 59, 14, 29, 29, 59, 28, 57, 57, 115,
		1, 3, 3, 7, 2, 5, 5, 11, 2, 5, 5, 11, 4, 9, 9, 19,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		7, 15, 15, 31, 14, 29, 29, 59, 14, 29, 29, 59, 28, 57, 57, 115,
		3, 7, 7, 15, 6, 13, 13, 27, 6, 13, 13, 27, 12, 25, 25, 51,
		7, 15, 15, 31, 14, 29, 29, 59, 14, 29, 29, 59, 28, 57, 57, 115,
		7, 15, 15, 31, 14, 29, 29, 59, 14, 29, 29, 59, 28, 57, 57, 115,
		15, 31, 31, 63, 30, 61, 61, 123, 30, 61, 61, 123, 60, 121, 121, 243,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 0, 0, 0, 0, 1, 2, 2, 4,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		7, 14, 14, 28, 15, 30, 30, 60, 14, 28, 28, 56, 29, 58, 58, 116,
		1, 2, 2, 4, 3, 6, 6, 12, 2, 4, 4, 8, 5, 10, 10, 20,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		7, 14, 14, 28, 15, 30, 30, 60, 14, 28, 28, 56, 29, 58, 58, 116,
		3, 6, 6, 12, 7, 14, 14, 28, 6, 12, 12, 24, 13, 26, 26, 52,
		7, 14, 14, 28, 15, 30, 30, 60, 14, 28, 28, 56, 29, 58, 58, 116,
		7, 14, 14, 28, 15, 30, 30, 60, 14, 28, 28, 56, 29, 58, 58, 116,
		15, 30, 30, 60, 31, 62, 62, 124, 30, 60, 60, 120, 61, 122, 122, 244,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 0, 1, 0, 1, 1, 3, 2, 5,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		7, 15, 14, 29, 15, 31, 30, 61, 14, 29, 28, 57, 29, 59, 58, 117,
		1, 3, 2, 5, 3, 7, 6, 13, 2, 5, 4, 9, 5, 11, 10, 21,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		7, 15, 14, 29, 15, 31, 30, 61, 14, 29, 28, 57, 29, 59, 58, 117,
		3, 7, 6, 13, 7, 15, 14, 29, 6, 13, 12, 25, 13, 27, 26, 53,
		7, 15, 14, 29, 15, 31, 30, 61, 14, 29, 28, 57, 29, 59, 58, 117,
		7, 15, 14, 29, 15, 31, 30, 61, 14, 29, 28, 57, 29, 59, 58, 117,
		15, 31, 30, 61, 31, 63, 62, 125, 30, 61, 60, 121, 61, 123, 122, 245,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 0, 0, 1, 2, 1, 2, 3, 6,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		7, 14, 15, 30, 15, 30, 31, 62, 14, 28, 29, 58, 29, 58, 59, 118,
		1, 2, 3, 6, 3, 6, 7, 14, 2, 4, 5, 10, 5, 10, 11, 22,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		7, 14, 15, 30, 15, 30, 31, 62, 14, 28, 29, 58, 29, 58, 59, 118,
		3, 6, 7, 14, 7, 14, 15, 30, 6, 12, 13, 26, 13, 26, 27, 54,
		7, 14, 15, 30, 15, 30, 31, 62, 14, 28, 29, 58, 29, 58, 59, 118,
		7, 14, 15, 30, 15, 30, 31, 62, 14, 28, 29, 58, 29, 58, 59, 118,
		15, 30, 31, 62, 31, 62, 63, 126, 30, 60, 61, 122, 61, 122, 123, 246,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 0, 1, 1, 3, 1, 3, 3, 7,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		7, 15, 15, 31, 15, 31, 31, 63, 14, 29, 29, 59, 29, 59, 59, 119,
		1, 3, 3, 7, 3, 7, 7, 15, 2, 5, 5, 11, 5, 11, 11, 23,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		7, 15, 15, 31, 15, 31, 31, 63, 14, 29, 29, 59, 29, 59, 59, 119,
		3, 7, 7, 15, 7, 15, 15, 31, 6, 13, 13, 27, 13, 27, 27, 55,
		7, 15, 15, 31, 15, 31, 31, 63, 14, 29, 29, 59, 29, 59, 59, 119,
		7, 15, 15, 31, 15, 31, 31, 63, 14, 29, 29, 59, 29, 59, 59, 119,
		15, 31, 31, 63, 31, 63, 63, 127, 30, 61, 61, 123, 61, 123, 123, 247,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 4, 2, 4, 4, 8,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		7, 14, 14, 28, 14, 28, 28, 56, 15, 30, 30, 60, 30, 60, 60, 120,
		1, 2, 2, 4, 2, 4, 4, 8, 3, 6, 6, 12, 6, 12, 12, 24,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		7, 14, 14, 28, 14, 28, 28, 56, 15, 30, 30, 60, 30, 60, 60, 120,
		3, 6, 6, 12, 6, 12, 12, 24, 7, 14, 14, 28, 14, 28, 28, 56,
		7, 14, 14, 28, 14, 28, 28, 56, 15, 30, 30, 60, 30, 60, 60, 120,
		7, 14, 14, 28, 14, 28, 28, 56, 15, 30, 30, 60, 30, 60, 60, 120,
		15, 30, 30, 60, 30, 60, 60, 120, 31, 62, 62, 124, 62, 124, 124, 248,
	}, {
		0, 1, 0, 1, 0, 1, 0, 1, 1, 3, 2, 5, 2, 5, 4, 9,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		7, 15, 14, 29, 14, 29, 28, 57, 15, 31, 30, 61, 30, 61, 60, 121,
		1, 3, 2, 5, 2, 5, 4, 9, 3, 7, 6, 13, 6, 13, 12, 25,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		7, 15, 14, 29, 14, 29, 28, 57, 15, 31, 30, 61, 30, 61, 60, 121,
		3, 7, 6, 13, 6, 13, 12, 25, 7, 15, 14, 29, 14, 29, 28, 57,
		7, 15, 14, 29, 14, 29, 28, 57, 15, 31, 30, 61, 30, 61, 60, 121,
		7, 15, 14, 29, 14, 29, 28, 57, 15, 31, 30, 61, 30, 61, 60, 121,
		15, 31, 30, 61, 30, 61, 60, 121, 31, 63, 62, 125, 62, 125, 124, 249,
	}, {
		0, 0, 1, 2, 0, 0, 1, 2, 1, 2, 3, 6, 2, 4, 5, 10,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		7, 14, 15, 30, 14, 28, 29, 58, 15, 30, 31, 62, 30, 60, 61, 122,
		1, 2, 3, 6, 2, 4, 5, 10, 3, 6, 7, 14, 6, 12, 13, 26,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		7, 14, 15, 30, 14, 28, 29, 58, 15, 30, 31, 62, 30, 60, 61, 122,
		3, 6, 7, 14, 6, 12, 13, 26, 7, 14, 15, 30, 14, 28, 29, 58,
		7, 14, 15, 30, 14, 28, 29, 58, 15, 30, 31, 62, 30, 60, 61, 122,
		7, 14, 15, 30, 14, 28, 29, 58, 15, 30, 31, 62, 30, 60, 61, 122,
		15, 30, 31, 62, 30, 60, 61, 122, 31, 62, 63, 126, 62, 124, 125, 250,
	}, {
		0, 1, 1, 3, 0, 1, 1, 3, 1, 3, 3, 7, 2, 5, 5, 11,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		7, 15, 15, 31, 14, 29, 29, 59, 15, 31, 31, 63, 30, 61, 61, 123,
		1, 3, 3, 7, 2, 5, 5, 11, 3, 7, 7, 15, 6, 13, 13, 27,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		7, 15, 15, 31, 14, 29, 29, 59, 15, 31, 31, 63, 30, 61, 61, 123,
		3, 7, 7, 15, 6, 13, 13, 27, 7, 15, 15, 31, 14, 29, 29, 59,
		7, 15, 15, 31, 14, 29, 29, 59, 15, 31, 31, 63, 30, 61, 61, 123,
		7, 15, 15, 31, 14, 29, 29, 59, 15, 31, 31, 63, 30, 61, 61, 123,
		15, 31, 31, 63, 30, 61, 61, 123, 31, 63, 63, 127, 62, 125, 125, 251,
	},
	{
		0, 0, 0, 0, 1, 2, 2, 4, 1, 2, 2, 4, 3, 6, 6, 12,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		7, 14, 14, 28, 15, 30, 30, 60, 15, 30, 30, 60, 31, 62, 62, 124,
		1, 2, 2, 4, 3, 6, 6, 12, 3, 6, 6, 12, 7, 14, 14, 28,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		7, 14, 14, 28, 15, 30, 30, 60, 15, 30, 30, 60, 31, 62, 62, 124,
		3, 6, 6, 12, 7, 14, 14, 28, 7, 14, 14, 28, 15, 30, 30, 60,
		7, 14, 14, 28, 15, 30, 30, 60, 15, 30, 30, 60, 31, 62, 62, 124,
		7, 14, 14, 28, 15, 30, 30, 60, 15, 30, 30, 60, 31, 62, 62, 124,
		15, 30, 30, 60, 31, 62, 62, 124, 31, 62, 62, 124, 63, 126, 126, 252,
	}, {
		0, 1, 0, 1, 1, 3, 2, 5, 1, 3, 2, 5, 3, 7, 6, 13,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		7, 15, 14, 29, 15, 31, 30, 61, 15, 31, 30, 61, 31, 63, 62, 125,
		1, 3, 2, 5, 3, 7, 6, 13, 3, 7, 6, 13, 7, 15, 14, 29,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		7, 15, 14, 29, 15, 31, 30, 61, 15, 31, 30, 61, 31, 63, 62, 125,
		3, 7, 6, 13, 7, 15, 14, 29, 7, 15, 14, 29, 15, 31, 30, 61,
		7, 15, 14, 29, 15, 31, 30, 61, 15, 31, 30, 61, 31, 63, 62, 125,
		7, 15, 14, 29, 15, 31, 30, 61, 15, 31, 30, 61, 31, 63, 62, 125,
		15, 31, 30, 61, 31, 63, 62, 125, 31, 63, 62, 125, 63, 127, 126, 253,
	}, {
		0, 0, 1, 2, 1, 2, 3, 6, 1, 2, 3, 6, 3, 6, 7, 14,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		7, 14, 15, 30, 15, 30, 31, 62, 15, 30, 31, 62, 31, 62, 63, 126,
		1, 2, 3, 6, 3, 6, 7, 14, 3, 6, 7, 14, 7, 14, 15, 30,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		7, 14, 15, 30, 15, 30, 31, 62, 15, 30, 31, 62, 31, 62, 63, 126,
		3, 6, 7, 14, 7, 14, 15, 30, 7, 14, 15, 30, 15, 30, 31, 62,
		7, 14, 15, 30, 15, 30, 31, 62, 15, 30, 31, 62, 31, 62, 63, 126,
		7, 14, 15, 30, 15, 30, 31, 62, 15, 30, 31, 62, 31, 62, 63, 126,
		15, 30, 31, 62, 31, 62, 63, 126, 31, 62, 63, 126, 63, 126, 127, 254,
	}, {
		0, 1, 1, 3, 1, 3, 3, 7, 1, 3, 3, 7, 3, 7, 7, 15,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		7, 15, 15, 31, 15, 31, 31, 63, 15, 31, 31, 63, 31, 63, 63, 127,
		1, 3, 3, 7, 3, 7, 7, 15, 3, 7, 7, 15, 7, 15, 15, 31,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		7, 15, 15, 31, 15, 31, 31, 63, 15, 31, 31, 63, 31, 63, 63, 127,
		3, 7, 7, 15, 7, 15, 15, 31, 7, 15, 15, 31, 15, 31, 31, 63,
		7, 15, 15, 31, 15, 31, 31, 63, 15, 31, 31, 63, 31, 63, 63, 127,
		7, 15, 15, 31, 15, 31, 31, 63, 15, 31, 31, 63, 31, 63, 63, 127,
		15, 31, 31, 63, 31, 63, 63, 127, 31, 63, 63, 127, 63, 127, 127, 255,
	},
}

// pdepLUT contains pre-computed parallel bit deposit results
var pdepLUT = [256][256]uint8{
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 2,
		64, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		64, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 2,
		128, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		128, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		128, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 2,
		64, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		64, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 3,
		80, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		96, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 3,
		144, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		160, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		192, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 3,
		80, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		96, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 4,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 4,
		0, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 4,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 4,
		0, 64, 64, 16, 64, 16, 16, 4, 64, 16, 16, 8, 16, 8, 8, 4,
		0, 64, 64, 32, 64, 32, 32, 4, 64, 32, 32, 8, 32, 8, 8, 4,
		64, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 4,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 4,
		0, 128, 128, 16, 128, 16, 16, 4, 128, 16, 16, 8, 16, 8, 8, 4,
		0, 128, 128, 32, 128, 32, 32, 4, 128, 32, 32, 8, 32, 8, 8, 4,
		128, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 4,
		0, 128, 128, 64, 128, 64, 64, 4, 128, 64, 64, 8, 64, 8, 8, 4,
		128, 64, 64, 16, 64, 16, 16, 4, 64, 16, 16, 8, 16, 8, 8, 4,
		128, 64, 64, 32, 64, 32, 32, 4, 64, 32, 32, 8, 32, 8, 8, 4,
		64, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 4,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 5,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 5,
		16, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 5,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 5,
		16, 65, 66, 17, 68, 17, 18, 5, 72, 17, 18, 9, 20, 9, 10, 5,
		32, 65, 66, 33, 68, 33, 34, 5, 72, 33, 34, 9, 36, 9, 10, 5,
		80, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 5,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 5,
		16, 129, 130, 17, 132, 17, 18, 5, 136, 17, 18, 9, 20, 9, 10, 5,
		32, 129, 130, 33, 132, 33, 34, 5, 136, 33, 34, 9, 36, 9, 10, 5,
		144, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 5,
		64, 129, 130, 65, 132, 65, 66, 5, 136, 65, 66, 9, 68, 9, 10, 5,
		144, 65, 66, 17, 68, 17, 18, 5, 72, 17, 18, 9, 20, 9, 10, 5,
		160, 65, 66, 33, 68, 33, 34, 5, 72, 33, 34, 9, 36, 9, 10, 5,
		80, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 5,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 6,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 6,
		32, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 6,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 6,
		64, 80, 80, 18, 80, 20, 20, 6, 80, 24, 24, 10, 24, 12, 12, 6,
		64, 96, 96, 34, 96, 36, 36, 6, 96, 40, 40, 10, 40, 12, 12, 6,
		96, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 6,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 6,
		128, 144, 144, 18, 144, 20, 20, 6, 144, 24, 24, 10, 24, 12, 12, 6,
		128, 160, 160, 34, 160, 36, 36, 6, 160, 40, 40, 10, 40, 12, 12, 6,
		160, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 6,
		128, 192, 192, 66, 192, 68, 68, 6, 192, 72, 72, 10, 72, 12, 12, 6,
		192, 80, 80, 18, 80, 20, 20, 6, 80, 24, 24, 10, 24, 12, 12, 6,
		192, 96, 96, 34, 96, 36, 36, 6, 96, 40, 40, 10, 40, 12, 12, 6,
		96, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 6,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 7,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 7,
		48, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 7,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 7,
		80, 81, 82, 19, 84, 21, 22, 7, 88, 25, 26, 11, 28, 13, 14, 7,
		96, 97, 98, 35, 100, 37, 38, 7, 104, 41, 42, 11, 44, 13, 14, 7,
		112, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 7,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 7,
		144, 145, 146, 19, 148, 21, 22, 7, 152, 25, 26, 11, 28, 13, 14, 7,
		160, 161, 162, 35, 164, 37, 38, 7, 168, 41, 42, 11, 44, 13, 14, 7,
		176, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 7,
		192, 193, 194, 67, 196, 69, 70, 7, 200, 73, 74, 11, 76, 13, 14, 7,
		208, 81, 82, 19, 84, 21, 22, 7, 88, 25, 26, 11, 28, 13, 14, 7,
		224, 97, 98, 35, 100, 37, 38, 7, 104, 41, 42, 11, 44, 13, 14, 7,
		112, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 7,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 8,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 8,
		0, 0, 0, 32, 0, 32, 32, 16, 0, 32, 32, 16, 32, 16, 16, 8,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 8,
		0, 0, 0, 64, 0, 64, 64, 16, 0, 64, 64, 16, 64, 16, 16, 8,
		0, 0, 0, 64, 0, 64, 64, 32, 0, 64, 64, 32, 64, 32, 32, 8,
		0, 64, 64, 32, 64, 32, 32, 16, 64, 32, 32, 16, 32, 16, 16, 8,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 8,
		0, 0, 0, 128, 0, 128, 128, 16, 0, 128, 128, 16, 128, 16, 16, 8,
		0, 0, 0, 128, 0, 128, 128, 32, 0, 128, 128, 32, 128, 32, 32, 8,
		0, 128, 128, 32, 128, 32, 32, 16, 128, 32, 32, 16, 32, 16, 16, 8,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 8,
		0, 128, 128, 64, 128, 64, 64, 16, 128, 64, 64, 16, 64, 16, 16, 8,
		0, 128, 128, 64, 128, 64, 64, 32, 128, 64, 64, 32, 64, 32, 32, 8,
		128, 64, 64, 32, 64, 32, 32, 16, 64, 32, 32, 16, 32, 16, 16, 8,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 9,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 9,
		16, 1, 2, 33, 4, 33, 34, 17, 8, 33, 34, 17, 36, 17, 18, 9,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 9,
		16, 1, 2, 65, 4, 65, 66, 17, 8, 65, 66, 17, 68, 17, 18, 9,
		32, 1, 2, 65, 4, 65, 66, 33, 8, 65, 66, 33, 68, 33, 34, 9,
		16, 65, 66, 33, 68, 33, 34, 17, 72, 33, 34, 17, 36, 17, 18, 9,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 9,
		16, 1, 2, 129, 4, 129, 130, 17, 8, 129, 130, 17, 132, 17, 18, 9,
		32, 1, 2, 129, 4, 129, 130, 33, 8, 129, 130, 33, 132, 33, 34, 9,
		16, 129, 130, 33, 132, 33, 34, 17, 136, 33, 34, 17, 36, 17, 18, 9,
		64, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 9,
		16, 129, 130, 65, 132, 65, 66, 17, 136, 65, 66, 17, 68, 17, 18, 9,
		32, 129, 130, 65, 132, 65, 66, 33, 136, 65, 66, 33, 68, 33, 34, 9,
		144, 65, 66, 33, 68, 33, 34, 17, 72, 33, 34, 17, 36, 17, 18, 9,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 10,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 10,
		32, 16, 16, 34, 16, 36, 36, 18, 16, 40, 40, 18, 40, 20, 20, 10,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 10,
		64, 16, 16, 66, 16, 68, 68, 18, 16, 72, 72, 18, 72, 20, 20, 10,
		64, 32, 32, 66, 32, 68, 68, 34, 32, 72, 72, 34, 72, 36, 36, 10,
		32, 80, 80, 34, 80, 36, 36, 18, 80, 40, 40, 18, 40, 20, 20, 10,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 10,
		128, 16, 16, 130, 16, 132, 132, 18, 16, 136, 136, 18, 136, 20, 20, 10,
		128, 32, 32, 130, 32, 132, 132, 34, 32, 136, 136, 34, 136, 36, 36, 10,
		32, 144, 144, 34, 144, 36, 36, 18, 144, 40, 40, 18, 40, 20, 20, 10,
		128, 64, 64, 130, 64, 132, 132, 66, 64, 136, 136, 66, 136, 68, 68, 10,
		64, 144, 144, 66, 144, 68, 68, 18, 144, 72, 72, 18, 72, 20, 20, 10,
		64, 160, 160, 66, 160, 68, 68, 34, 160, 72, 72, 34, 72, 36, 36, 10,
		160, 80, 80, 34, 80, 36, 36, 18, 80, 40, 40, 18, 40, 20, 20, 10,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 11,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 11,
		48, 17, 18, 35, 20, 37, 38, 19, 24, 41, 42, 19, 44, 21, 22, 11,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 11,
		80, 17, 18, 67, 20, 69, 70, 19, 24, 73, 74, 19, 76, 21, 22, 11,
		96, 33, 34, 67, 36, 69, 70, 35, 40, 73, 74, 35, 76, 37, 38, 11,
		48, 81, 82, 35, 84, 37, 38, 19, 88, 41, 42, 19, 44, 21, 22, 11,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 11,
		144, 17, 18, 131, 20, 133, 134, 19, 24, 137, 138, 19, 140, 21, 22, 11,
		160, 33, 34, 131, 36, 133, 134, 35, 40, 137, 138, 35, 140, 37, 38, 11,
		48, 145, 146, 35, 148, 37, 38, 19, 152, 41, 42, 19, 44, 21, 22, 11,
		192, 65, 66, 131, 68, 133, 134, 67, 72, 137, 138, 67, 140, 69, 70, 11,
		80, 145, 146, 67, 148, 69, 70, 19, 152, 73, 74, 19, 76, 21, 22, 11,
		96, 161, 162, 67, 164, 69, 70, 35, 168, 73, 74, 35, 76, 37, 38, 11,
		176, 81, 82, 35, 84, 37, 38, 19, 88, 41, 42, 19, 44, 21, 22, 11,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 12,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 12,
		0, 32, 32, 48, 32, 48, 48, 20, 32, 48, 48, 24, 48, 24, 24, 12,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 12,
		0, 64, 64, 80, 64, 80, 80, 20, 64, 80, 80, 24, 80, 24, 24, 12,
		0, 64, 64, 96, 64, 96, 96, 36, 64, 96, 96, 40, 96, 40, 40, 12,
		64, 96, 96, 48, 96, 48, 48, 20, 96, 48, 48, 24, 48, 24, 24, 12,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 12,
		0, 128, 128, 144, 128, 144, 144, 20, 128, 144, 144, 24, 144, 24, 24, 12,
		0, 128, 128, 160, 128, 160, 160, 36, 128, 160, 160, 40, 160, 40, 40, 12,
		128, 160, 160, 48, 160, 48, 48, 20, 160, 48, 48, 24, 48, 24, 24, 12,
		0, 128, 128, 192, 128, 192, 192, 68, 128, 192, 192, 72, 192, 72, 72, 12,
		128, 192, 192, 80, 192, 80, 80, 20, 192, 80, 80, 24, 80, 24, 24, 12,
		128, 192, 192, 96, 192, 96, 96, 36, 192, 96, 96, 40, 96, 40, 40, 12,
		192, 96, 96, 48, 96, 48, 48, 20, 96, 48, 48, 24, 48, 24, 24, 12,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 13,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 13,
		16, 33, 34, 49, 36, 49, 50, 21, 40, 49, 50, 25, 52, 25, 26, 13,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 13,
		16, 65, 66, 81, 68, 81, 82, 21, 72, 81, 82, 25, 84, 25, 26, 13,
		32, 65, 66, 97, 68, 97, 98, 37, 72, 97, 98, 41, 100, 41, 42, 13,
		80, 97, 98, 49, 100, 49, 50, 21, 104, 49, 50, 25, 52, 25, 26, 13,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 13,
		16, 129, 130, 145, 132, 145, 146, 21, 136, 145, 146, 25, 148, 25, 26, 13,
		32, 129, 130, 161, 132, 161, 162, 37, 136, 161, 162, 41, 164, 41, 42, 13,
		144, 161, 162, 49, 164, 49, 50, 21, 168, 49, 50, 25, 52, 25, 26, 13,
		64, 129, 130, 193, 132, 193, 194, 69, 136, 193, 194, 73, 196, 73, 74, 13,
		144, 193, 194, 81, 196, 81, 82, 21, 200, 81, 82, 25, 84, 25, 26, 13,
		160, 193, 194, 97, 196, 97, 98, 37, 200, 97, 98, 41, 100, 41, 42, 13,
		208, 97, 98, 49, 100, 49, 50, 21, 104, 49, 50, 25, 52, 25, 26, 13,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 14,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 14,
		32, 48, 48, 50, 48, 52, 52, 22, 48, 56, 56, 26, 56, 28, 28, 14,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 14,
		64, 80, 80, 82, 80, 84, 84, 22, 80, 88, 88, 26, 88, 28, 28, 14,
		64, 96, 96, 98, 96, 100, 100, 38, 96, 104, 104, 42, 104, 44, 44, 14,
		96, 112, 112, 50, 112, 52, 52, 22, 112, 56, 56, 26, 56, 28, 28, 14,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 14,
		128, 144, 144, 146, 144, 148, 148, 22, 144, 152, 152, 26, 152, 28, 28, 14,
		128, 160, 160, 162, 160, 164, 164, 38, 160, 168, 168, 42, 168, 44, 44, 14,
		160, 176, 176, 50, 176, 52, 52, 22, 176, 56, 56, 26, 56, 28, 28, 14,
		128, 192, 192, 194, 192, 196, 196, 70, 192, 200, 200, 74, 200, 76, 76, 14,
		192, 208, 208, 82, 208, 84, 84, 22, 208, 88, 88, 26, 88, 28, 28, 14,
		192, 224, 224, 98, 224, 100, 100, 38, 224, 104, 104, 42, 104, 44, 44, 14,
		224, 112, 112, 50, 112, 52, 52, 22, 112, 56, 56, 26, 56, 28, 28, 14,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 15,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 15,
		48, 49, 50, 51, 52, 53, 54, 23, 56, 57, 58, 27, 60, 29, 30, 15,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 15,
		80, 81, 82, 83, 84, 85, 86, 23, 88, 89, 90, 27, 92, 29, 30, 15,
		96, 97, 98, 99, 100, 101, 102, 39, 104, 105, 106, 43, 108, 45, 46, 15,
		112, 113, 114, 51, 116, 53, 54, 23, 120, 57, 58, 27, 60, 29, 30, 15,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 15,
		144, 145, 146, 147, 148, 149, 150, 23, 152, 153, 154, 27, 156, 29, 30, 15,
		160, 161, 162, 163, 164, 165, 166, 39, 168, 169, 170, 43, 172, 45, 46, 15,
		176, 177, 178, 51, 180, 53, 54, 23, 184, 57, 58, 27, 60, 29, 30, 15,
		192, 193, 194, 195, 196, 197, 198, 71, 200, 201, 202, 75, 204, 77, 78, 15,
		208, 209, 210, 83, 212, 85, 86, 23, 216, 89, 90, 27, 92, 29, 30, 15,
		224, 225, 226, 99, 228, 101, 102, 39, 232, 105, 106, 43, 108, 45, 46, 15,
		240, 113, 114, 51, 116, 53, 54, 23, 120, 57, 58, 27, 60, 29, 30, 15,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 16,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 32,
		0, 0, 0, 64, 0, 64, 64, 32, 0, 64, 64, 32, 64, 32, 32, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 16,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 32,
		0, 0, 0, 128, 0, 128, 128, 32, 0, 128, 128, 32, 128, 32, 32, 16,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 64,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 16,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 32,
		0, 128, 128, 64, 128, 64, 64, 32, 128, 64, 64, 32, 64, 32, 32, 16,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 17,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 33,
		16, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 17,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		16, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 17,
		32, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 33,
		16, 1, 2, 65, 4, 65, 66, 33, 8, 65, 66, 33, 68, 33, 34, 17,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 17,
		32, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 33,
		16, 1, 2, 129, 4, 129, 130, 33, 8, 129, 130, 33, 132, 33, 34, 17,
		64, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 65,
		16, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 17,
		32, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 33,
		16, 129, 130, 65, 132, 65, 66, 33, 136, 65, 66, 33, 68, 33, 34, 17,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 18,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 34,
		32, 16, 16, 2, 16, 4, 4, 34, 16, 8, 8, 34, 8, 36, 36, 18,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 66,
		64, 16, 16, 2, 16, 4, 4, 66, 16, 8, 8, 66, 8, 68, 68, 18,
		64, 32, 32, 2, 32, 4, 4, 66, 32, 8, 8, 66, 8, 68, 68, 34,
		32, 16, 16, 66, 16, 68, 68, 34, 16, 72, 72, 34, 72, 36, 36, 18,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 130,
		128, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 18,
		128, 32, 32, 2, 32, 4, 4, 130, 32, 8, 8, 130, 8, 132, 132, 34,
		32, 16, 16, 130, 16, 132, 132, 34, 16, 136, 136, 34, 136, 36, 36, 18,
		128, 64, 64, 2, 64, 4, 4, 130, 64, 8, 8, 130, 8, 132, 132, 66,
		64, 16, 16, 130, 16, 132, 132, 66, 16, 136, 136, 66, 136, 68, 68, 18,
		64, 32, 32, 130, 32, 132, 132, 66, 32, 136, 136, 66, 136, 68, 68, 34,
		32, 144, 144, 66, 144, 68, 68, 34, 144, 72, 72, 34, 72, 36, 36, 18,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 19,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 35,
		48, 17, 18, 3, 20, 5, 6, 35, 24, 9, 10, 35, 12, 37, 38, 19,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 67,
		80, 17, 18, 3, 20, 5, 6, 67, 24, 9, 10, 67, 12, 69, 70, 19,
		96, 33, 34, 3, 36, 5, 6, 67, 40, 9, 10, 67, 12, 69, 70, 35,
		48, 17, 18, 67, 20, 69, 70, 35, 24, 73, 74, 35, 76, 37, 38, 19,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 131,
		144, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 19,
		160, 33, 34, 3, 36, 5, 6, 131, 40, 9, 10, 131, 12, 133, 134, 35,
		48, 17, 18, 131, 20, 133, 134, 35, 24, 137, 138, 35, 140, 37, 38, 19,
		192, 65, 66, 3, 68, 5, 6, 131, 72, 9, 10, 131, 12, 133, 134, 67,
		80, 17, 18, 131, 20, 133, 134, 67, 24, 137, 138, 67, 140, 69, 70, 19,
		96, 33, 34, 131, 36, 133, 134, 67, 40, 137, 138, 67, 140, 69, 70, 35,
		48, 145, 146, 67, 148, 69, 70, 35, 152, 73, 74, 35, 76, 37, 38, 19,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 20,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 36,
		0, 32, 32, 16, 32, 16, 16, 36, 32, 16, 16, 40, 16, 40, 40, 20,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 68,
		0, 64, 64, 16, 64, 16, 16, 68, 64, 16, 16, 72, 16, 72, 72, 20,
		0, 64, 64, 32, 64, 32, 32, 68, 64, 32, 32, 72, 32, 72, 72, 36,
		64, 32, 32, 80, 32, 80, 80, 36, 32, 80, 80, 40, 80, 40, 40, 20,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 132,
		0, 128, 128, 16, 128, 16, 16, 132, 128, 16, 16, 136, 16, 136, 136, 20,
		0, 128, 128, 32, 128, 32, 32, 132, 128, 32, 32, 136, 32, 136, 136, 36,
		128, 32, 32, 144, 32, 144, 144, 36, 32, 144, 144, 40, 144, 40, 40, 20,
		0, 128, 128, 64, 128, 64, 64, 132, 128, 64, 64, 136, 64, 136, 136, 68,
		128, 64, 64, 144, 64, 144, 144, 68, 64, 144, 144, 72, 144, 72, 72, 20,
		128, 64, 64, 160, 64, 160, 160, 68, 64, 160, 160, 72, 160, 72, 72, 36,
		64, 160, 160, 80, 160, 80, 80, 36, 160, 80, 80, 40, 80, 40, 40, 20,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 21,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 37,
		16, 33, 34, 17, 36, 17, 18, 37, 40, 17, 18, 41, 20, 41, 42, 21,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 69,
		16, 65, 66, 17, 68, 17, 18, 69, 72, 17, 18, 73, 20, 73, 74, 21,
		32, 65, 66, 33, 68, 33, 34, 69, 72, 33, 34, 73, 36, 73, 74, 37,
		80, 33, 34, 81, 36, 81, 82, 37, 40, 81, 82, 41, 84, 41, 42, 21,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 133,
		16, 129, 130, 17, 132, 17, 18, 133, 136, 17, 18, 137, 20, 137, 138, 21,
		32, 129, 130, 33, 132, 33, 34, 133, 136, 33, 34, 137, 36, 137, 138, 37,
		144, 33, 34, 145, 36, 145, 146, 37, 40, 145, 146, 41, 148, 41, 42, 21,
		64, 129, 130, 65, 132, 65, 66, 133, 136, 65, 66, 137, 68, 137, 138, 69,
		144, 65, 66, 145, 68, 145, 146, 69, 72, 145, 146, 73, 148, 73, 74, 21,
		160, 65, 66, 161, 68, 161, 162, 69, 72, 161, 162, 73, 164, 73, 74, 37,
		80, 161, 162, 81, 164, 81, 82, 37, 168, 81, 82, 41, 84, 41, 42, 21,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 22,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 38,
		32, 48, 48, 18, 48, 20, 20, 38, 48, 24, 24, 42, 24, 44, 44, 22,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 70,
		64, 80, 80, 18, 80, 20, 20, 70, 80, 24, 24, 74, 24, 76, 76, 22,
		64, 96, 96, 34, 96, 36, 36, 70, 96, 40, 40, 74, 40, 76, 76, 38,
		96, 48, 48, 82, 48, 84, 84, 38, 48, 88, 88, 42, 88, 44, 44, 22,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 134,
		128, 144, 144, 18, 144, 20, 20, 134, 144, 24, 24, 138, 24, 140, 140, 22,
		128, 160, 160, 34, 160, 36, 36, 134, 160, 40, 40, 138, 40, 140, 140, 38,
		160, 48, 48, 146, 48, 148, 148, 38, 48, 152, 152, 42, 152, 44, 44, 22,
		128, 192, 192, 66, 192, 68, 68, 134, 192, 72, 72, 138, 72, 140, 140, 70,
		192, 80, 80, 146, 80, 148, 148, 70, 80, 152, 152, 74, 152, 76, 76, 22,
		192, 96, 96, 162, 96, 164, 164, 70, 96, 168, 168, 74, 168, 76, 76, 38,
		96, 176, 176, 82, 176, 84, 84, 38, 176, 88, 88, 42, 88, 44, 44, 22,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 23,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 39,
		48, 49, 50, 19, 52, 21, 22, 39, 56, 25, 26, 43, 28, 45, 46, 23,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 71,
		80, 81, 82, 19, 84, 21, 22, 71, 88, 25, 26, 75, 28, 77, 78, 23,
		96, 97, 98, 35, 100, 37, 38, 71, 104, 41, 42, 75, 44, 77, 78, 39,
		112, 49, 50, 83, 52, 85, 86, 39, 56, 89, 90, 43, 92, 45, 46, 23,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 135,
		144, 145, 146, 19, 148, 21, 22, 135, 152, 25, 26, 139, 28, 141, 142, 23,
		160, 161, 162, 35, 164, 37, 38, 135, 168, 41, 42, 139, 44, 141, 142, 39,
		176, 49, 50, 147, 52, 149, 150, 39, 56, 153, 154, 43, 156, 45, 46, 23,
		192, 193, 194, 67, 196, 69, 70, 135, 200, 73, 74, 139, 76, 141, 142, 71,
		208, 81, 82, 147, 84, 149, 150, 71, 88, 153, 154, 75, 156, 77, 78, 23,
		224, 97, 98, 163, 100, 165, 166, 71, 104, 169, 170, 75, 172, 77, 78, 39,
		112, 177, 178, 83, 180, 85, 86, 39, 184, 89, 90, 43, 92, 45, 46, 23,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 24,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 40,
		0, 0, 0, 32, 0, 32, 32, 48, 0, 32, 32, 48, 32, 48, 48, 24,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 72,
		0, 0, 0, 64, 0, 64, 64, 80, 0, 64, 64, 80, 64, 80, 80, 24,
		0, 0, 0, 64, 0, 64, 64, 96, 0, 64, 64, 96, 64, 96, 96, 40,
		0, 64, 64, 96, 64, 96, 96, 48, 64, 96, 96, 48, 96, 48, 48, 24,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 136,
		0, 0, 0, 128, 0, 128, 128, 144, 0, 128, 128, 144, 128, 144, 144, 24,
		0, 0, 0, 128, 0, 128, 128, 160, 0, 128, 128, 160, 128, 160, 160, 40,
		0, 128, 128, 160, 128, 160, 160, 48, 128, 160, 160, 48, 160, 48, 48, 24,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 72,
		0, 128, 128, 192, 128, 192, 192, 80, 128, 192, 192, 80, 192, 80, 80, 24,
		0, 128, 128, 192, 128, 192, 192, 96, 128, 192, 192, 96, 192, 96, 96, 40,
		128, 192, 192, 96, 192, 96, 96, 48, 192, 96, 96, 48, 96, 48, 48, 24,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 25,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 41,
		16, 1, 2, 33, 4, 33, 34, 49, 8, 33, 34, 49, 36, 49, 50, 25,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 73,
		16, 1, 2, 65, 4, 65, 66, 81, 8, 65, 66, 81, 68, 81, 82, 25,
		32, 1, 2, 65, 4, 65, 66, 97, 8, 65, 66, 97, 68, 97, 98, 41,
		16, 65, 66, 97, 68, 97, 98, 49, 72, 97, 98, 49, 100, 49, 50, 25,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 137,
		16, 1, 2, 129, 4, 129, 130, 145, 8, 129, 130, 145, 132, 145, 146, 25,
		32, 1, 2, 129, 4, 129, 130, 161, 8, 129, 130, 161, 132, 161, 162, 41,
		16, 129, 130, 161, 132, 161, 162, 49, 136, 161, 162, 49, 164, 49, 50, 25,
		64, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 73,
		16, 129, 130, 193, 132, 193, 194, 81, 136, 193, 194, 81, 196, 81, 82, 25,
		32, 129, 130, 193, 132, 193, 194, 97, 136, 193, 194, 97, 196, 97, 98, 41,
		144, 193, 194, 97, 196, 97, 98, 49, 200, 97, 98, 49, 100, 49, 50, 25,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 26,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 42,
		32, 16, 16, 34, 16, 36, 36, 50, 16, 40, 40, 50, 40, 52, 52, 26,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 74,
		64, 16, 16, 66, 16, 68, 68, 82, 16, 72, 72, 82, 72, 84, 84, 26,
		64, 32, 32, 66, 32, 68, 68, 98, 32, 72, 72, 98, 72, 100, 100, 42,
		32, 80, 80, 98, 80, 100, 100, 50, 80, 104, 104, 50, 104, 52, 52, 26,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 138,
		128, 16, 16, 130, 16, 132, 132, 146, 16, 136, 136, 146, 136, 148, 148, 26,
		128, 32, 32, 130, 32, 132, 132, 162, 32, 136, 136, 162, 136, 164, 164, 42,
		32, 144, 144, 162, 144, 164, 164, 50, 144, 168, 168, 50, 168, 52, 52, 26,
		128, 64, 64, 130, 64, 132, 132, 194, 64, 136, 136, 194, 136, 196, 196, 74,
		64, 144, 144, 194, 144, 196, 196, 82, 144, 200, 200, 82, 200, 84, 84, 26,
		64, 160, 160, 194, 160, 196, 196, 98, 160, 200, 200, 98, 200, 100, 100, 42,
		160, 208, 208, 98, 208, 100, 100, 50, 208, 104, 104, 50, 104, 52, 52, 26,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 27,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 43,
		48, 17, 18, 35, 20, 37, 38, 51, 24, 41, 42, 51, 44, 53, 54, 27,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 75,
		80, 17, 18, 67, 20, 69, 70, 83, 24, 73, 74, 83, 76, 85, 86, 27,
		96, 33, 34, 67, 36, 69, 70, 99, 40, 73, 74, 99, 76, 101, 102, 43,
		48, 81, 82, 99, 84, 101, 102, 51, 88, 105, 106, 51, 108, 53, 54, 27,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 139,
		144, 17, 18, 131, 20, 133, 134, 147, 24, 137, 138, 147, 140, 149, 150, 27,
		160, 33, 34, 131, 36, 133, 134, 163, 40, 137, 138, 163, 140, 165, 166, 43,
		48, 145, 146, 163, 148, 165, 166, 51, 152, 169, 170, 51, 172, 53, 54, 27,
		192, 65, 66, 131, 68, 133, 134, 195, 72, 137, 138, 195, 140, 197, 198, 75,
		80, 145, 146, 195, 148, 197, 198, 83, 152, 201, 202, 83, 204, 85, 86, 27,
		96, 161, 162, 195, 164, 197, 198, 99, 168, 201, 202, 99, 204, 101, 102, 43,
		176, 209, 210, 99, 212, 101, 102, 51, 216, 105, 106, 51, 108, 53, 54, 27,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 28,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 44,
		0, 32, 32, 48, 32, 48, 48, 52, 32, 48, 48, 56, 48, 56, 56, 28,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 76,
		0, 64, 64, 80, 64, 80, 80, 84, 64, 80, 80, 88, 80, 88, 88, 28,
		0, 64, 64, 96, 64, 96, 96, 100, 64, 96, 96, 104, 96, 104, 104, 44,
		64, 96, 96, 112, 96, 112, 112, 52, 96, 112, 112, 56, 112, 56, 56, 28,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 140,
		0, 128, 128, 144, 128, 144, 144, 148, 128, 144, 144, 152, 144, 152, 152, 28,
		0, 128, 128, 160, 128, 160, 160, 164, 128, 160, 160, 168, 160, 168, 168, 44,
		128, 160, 160, 176, 160, 176, 176, 52, 160, 176, 176, 56, 176, 56, 56, 28,
		0, 128, 128, 192, 128, 192, 192, 196, 128, 192, 192, 200, 192, 200, 200, 76,
		128, 192, 192, 208, 192, 208, 208, 84, 192, 208, 208, 88, 208, 88, 88, 28,
		128, 192, 192, 224, 192, 224, 224, 100, 192, 224, 224, 104, 224, 104, 104, 44,
		192, 224, 224, 112, 224, 112, 112, 52, 224, 112, 112, 56, 112, 56, 56, 28,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 29,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 45,
		16, 33, 34, 49, 36, 49, 50, 53, 40, 49, 50, 57, 52, 57, 58, 29,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 77,
		16, 65, 66, 81, 68, 81, 82, 85, 72, 81, 82, 89, 84, 89, 90, 29,
		32, 65, 66, 97, 68, 97, 98, 101, 72, 97, 98, 105, 100, 105, 106, 45,
		80, 97, 98, 113, 100, 113, 114, 53, 104, 113, 114, 57, 116, 57, 58, 29,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 141,
		16, 129, 130, 145, 132, 145, 146, 149, 136, 145, 146, 153, 148, 153, 154, 29,
		32, 129, 130, 161, 132, 161, 162, 165, 136, 161, 162, 169, 164, 169, 170, 45,
		144, 161, 162, 177, 164, 177, 178, 53, 168, 177, 178, 57, 180, 57, 58, 29,
		64, 129, 130, 193, 132, 193, 194, 197, 136, 193, 194, 201, 196, 201, 202, 77,
		144, 193, 194, 209, 196, 209, 210, 85, 200, 209, 210, 89, 212, 89, 90, 29,
		160, 193, 194, 225, 196, 225, 226, 101, 200, 225, 226, 105, 228, 105, 106, 45,
		208, 225, 226, 113, 228, 113, 114, 53, 232, 113, 114, 57, 116, 57, 58, 29,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 30,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 46,
		32, 48, 48, 50, 48, 52, 52, 54, 48, 56, 56, 58, 56, 60, 60, 30,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 78,
		64, 80, 80, 82, 80, 84, 84, 86, 80, 88, 88, 90, 88, 92, 92, 30,
		64, 96, 96, 98, 96, 100, 100, 102, 96, 104, 104, 106, 104, 108, 108, 46,
		96, 112, 112, 114, 112, 116, 116, 54, 112, 120, 120, 58, 120, 60, 60, 30,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 142,
		128, 144, 144, 146, 144, 148, 148, 150, 144, 152, 152, 154, 152, 156, 156, 30,
		128, 160, 160, 162, 160, 164, 164, 166, 160, 168, 168, 170, 168, 172, 172, 46,
		160, 176, 176, 178, 176, 180, 180, 54, 176, 184, 184, 58, 184, 60, 60, 30,
		128, 192, 192, 194, 192, 196, 196, 198, 192, 200, 200, 202, 200, 204, 204, 78,
		192, 208, 208, 210, 208, 212, 212, 86, 208, 216, 216, 90, 216, 92, 92, 30,
		192, 224, 224, 226, 224, 228, 228, 102, 224, 232, 232, 106, 232, 108, 108, 46,
		224, 240, 240, 114, 240, 116, 116, 54, 240, 120, 120, 58, 120, 60, 60, 30,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 31,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
		80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 31,
		96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 47,
		112, 113, 114, 115, 116, 117, 118, 55, 120, 121, 122, 59, 124, 61, 62, 31,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143,
		144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 31,
		160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 47,
		176, 177, 178, 179, 180, 181, 182, 55, 184, 185, 186, 59, 188, 61, 62, 31,
		192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 79,
		208, 209, 210, 211, 212, 213, 214, 87, 216, 217, 218, 91, 220, 93, 94, 31,
		224, 225, 226, 227, 228, 229, 230, 103, 232, 233, 234, 107, 236, 109, 110, 47,
		240, 241, 242, 115, 244, 117, 118, 55, 248, 121, 122, 59, 124, 61, 62, 31,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 32,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 32,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 64,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 64,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 32,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 33,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		16, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 33,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 33,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 65,
		32, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 65,
		16, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 33,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 34,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 2,
		64, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 66,
		64, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 66,
		32, 16, 16, 2, 16, 4, 4, 66, 16, 8, 8, 66, 8, 68, 68, 34,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 2,
		128, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 130,
		128, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 130,
		32, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 34,
		128, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 130,
		64, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 66,
		64, 32, 32, 2, 32, 4, 4, 130, 32, 8, 8, 130, 8, 132, 132, 66,
		32, 16, 16, 130, 16, 132, 132, 66, 16, 136, 136, 66, 136, 68, 68, 34,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 35,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 3,
		80, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 67,
		96, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 67,
		48, 17, 18, 3, 20, 5, 6, 67, 24, 9, 10, 67, 12, 69, 70, 35,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 3,
		144, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 131,
		160, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 131,
		48, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 35,
		192, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 131,
		80, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 67,
		96, 33, 34, 3, 36, 5, 6, 131, 40, 9, 10, 131, 12, 133, 134, 67,
		48, 17, 18, 131, 20, 133, 134, 67, 24, 137, 138, 67, 140, 69, 70, 35,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 4,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 4,
		0, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 36,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 4,
		0, 64, 64, 16, 64, 16, 16, 4, 64, 16, 16, 8, 16, 8, 8, 68,
		0, 64, 64, 32, 64, 32, 32, 4, 64, 32, 32, 8, 32, 8, 8, 68,
		64, 32, 32, 16, 32, 16, 16, 68, 32, 16, 16, 72, 16, 72, 72, 36,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 4,
		0, 128, 128, 16, 128, 16, 16, 4, 128, 16, 16, 8, 16, 8, 8, 132,
		0, 128, 128, 32, 128, 32, 32, 4, 128, 32, 32, 8, 32, 8, 8, 132,
		128, 32, 32, 16, 32, 16, 16, 132, 32, 16, 16, 136, 16, 136, 136, 36,
		0, 128, 128, 64, 128, 64, 64, 4, 128, 64, 64, 8, 64, 8, 8, 132,
		128, 64, 64, 16, 64, 16, 16, 132, 64, 16, 16, 136, 16, 136, 136, 68,
		128, 64, 64, 32, 64, 32, 32, 132, 64, 32, 32, 136, 32, 136, 136, 68,
		64, 32, 32, 144, 32, 144, 144, 68, 32, 144, 144, 72, 144, 72, 72, 36,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 5,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 5,
		16, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 37,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 5,
		16, 65, 66, 17, 68, 17, 18, 5, 72, 17, 18, 9, 20, 9, 10, 69,
		32, 65, 66, 33, 68, 33, 34, 5, 72, 33, 34, 9, 36, 9, 10, 69,
		80, 33, 34, 17, 36, 17, 18, 69, 40, 17, 18, 73, 20, 73, 74, 37,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 5,
		16, 129, 130, 17, 132, 17, 18, 5, 136, 17, 18, 9, 20, 9, 10, 133,
		32, 129, 130, 33, 132, 33, 34, 5, 136, 33, 34, 9, 36, 9, 10, 133,
		144, 33, 34, 17, 36, 17, 18, 133, 40, 17, 18, 137, 20, 137, 138, 37,
		64, 129, 130, 65, 132, 65, 66, 5, 136, 65, 66, 9, 68, 9, 10, 133,
		144, 65, 66, 17, 68, 17, 18, 133, 72, 17, 18, 137, 20, 137, 138, 69,
		160, 65, 66, 33, 68, 33, 34, 133, 72, 33, 34, 137, 36, 137, 138, 69,
		80, 33, 34, 145, 36, 145, 146, 69, 40, 145, 146, 73, 148, 73, 74, 37,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 6,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 6,
		32, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 38,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 6,
		64, 80, 80, 18, 80, 20, 20, 6, 80, 24, 24, 10, 24, 12, 12, 70,
		64, 96, 96, 34, 96, 36, 36, 6, 96, 40, 40, 10, 40, 12, 12, 70,
		96, 48, 48, 18, 48, 20, 20, 70, 48, 24, 24, 74, 24, 76, 76, 38,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 6,
		128, 144, 144, 18, 144, 20, 20, 6, 144, 24, 24, 10, 24, 12, 12, 134,
		128, 160, 160, 34, 160, 36, 36, 6, 160, 40, 40, 10, 40, 12, 12, 134,
		160, 48, 48, 18, 48, 20, 20, 134, 48, 24, 24, 138, 24, 140, 140, 38,
		128, 192, 192, 66, 192, 68, 68, 6, 192, 72, 72, 10, 72, 12, 12, 134,
		192, 80, 80, 18, 80, 20, 20, 134, 80, 24, 24, 138, 24, 140, 140, 70,
		192, 96, 96, 34, 96, 36, 36, 134, 96, 40, 40, 138, 40, 140, 140, 70,
		96, 48, 48, 146, 48, 148, 148, 70, 48, 152, 152, 74, 152, 76, 76, 38,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 7,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 7,
		48, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 39,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 7,
		80, 81, 82, 19, 84, 21, 22, 7, 88, 25, 26, 11, 28, 13, 14, 71,
		96, 97, 98, 35, 100, 37, 38, 7, 104, 41, 42, 11, 44, 13, 14, 71,
		112, 49, 50, 19, 52, 21, 22, 71, 56, 25, 26, 75, 28, 77, 78, 39,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 7,
		144, 145, 146, 19, 148, 21, 22, 7, 152, 25, 26, 11, 28, 13, 14, 135,
		160, 161, 162, 35, 164, 37, 38, 7, 168, 41, 42, 11, 44, 13, 14, 135,
		176, 49, 50, 19, 52, 21, 22, 135, 56, 25, 26, 139, 28, 141, 142, 39,
		192, 193, 194, 67, 196, 69, 70, 7, 200, 73, 74, 11, 76, 13, 14, 135,
		208, 81, 82, 19, 84, 21, 22, 135, 88, 25, 26, 139, 28, 141, 142, 71,
		224, 97, 98, 35, 100, 37, 38, 135, 104, 41, 42, 139, 44, 141, 142, 71,
		112, 49, 50, 147, 52, 149, 150, 71, 56, 153, 154, 75, 156, 77, 78, 39,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 8,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 8,
		0, 0, 0, 32, 0, 32, 32, 16, 0, 32, 32, 16, 32, 16, 16, 40,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 8,
		0, 0, 0, 64, 0, 64, 64, 16, 0, 64, 64, 16, 64, 16, 16, 72,
		0, 0, 0, 64, 0, 64, 64, 32, 0, 64, 64, 32, 64, 32, 32, 72,
		0, 64, 64, 32, 64, 32, 32, 80, 64, 32, 32, 80, 32, 80, 80, 40,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 8,
		0, 0, 0, 128, 0, 128, 128, 16, 0, 128, 128, 16, 128, 16, 16, 136,
		0, 0, 0, 128, 0, 128, 128, 32, 0, 128, 128, 32, 128, 32, 32, 136,
		0, 128, 128, 32, 128, 32, 32, 144, 128, 32, 32, 144, 32, 144, 144, 40,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 136,
		0, 128, 128, 64, 128, 64, 64, 144, 128, 64, 64, 144, 64, 144, 144, 72,
		0, 128, 128, 64, 128, 64, 64, 160, 128, 64, 64, 160, 64, 160, 160, 72,
		128, 64, 64, 160, 64, 160, 160, 80, 64, 160, 160, 80, 160, 80, 80, 40,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 9,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 9,
		16, 1, 2, 33, 4, 33, 34, 17, 8, 33, 34, 17, 36, 17, 18, 41,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 9,
		16, 1, 2, 65, 4, 65, 66, 17, 8, 65, 66, 17, 68, 17, 18, 73,
		32, 1, 2, 65, 4, 65, 66, 33, 8, 65, 66, 33, 68, 33, 34, 73,
		16, 65, 66, 33, 68, 33, 34, 81, 72, 33, 34, 81, 36, 81, 82, 41,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 9,
		16, 1, 2, 129, 4, 129, 130, 17, 8, 129, 130, 17, 132, 17, 18, 137,
		32, 1, 2, 129, 4, 129, 130, 33, 8, 129, 130, 33, 132, 33, 34, 137,
		16, 129, 130, 33, 132, 33, 34, 145, 136, 33, 34, 145, 36, 145, 146, 41,
		64, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 137,
		16, 129, 130, 65, 132, 65, 66, 145, 136, 65, 66, 145, 68, 145, 146, 73,
		32, 129, 130, 65, 132, 65, 66, 161, 136, 65, 66, 161, 68, 161, 162, 73,
		144, 65, 66, 161, 68, 161, 162, 81, 72, 161, 162, 81, 164, 81, 82, 41,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 10,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 10,
		32, 16, 16, 34, 16, 36, 36, 18, 16, 40, 40, 18, 40, 20, 20, 42,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 10,
		64, 16, 16, 66, 16, 68, 68, 18, 16, 72, 72, 18, 72, 20, 20, 74,
		64, 32, 32, 66, 32, 68, 68, 34, 32, 72, 72, 34, 72, 36, 36, 74,
		32, 80, 80, 34, 80, 36, 36, 82, 80, 40, 40, 82, 40, 84, 84, 42,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 10,
		128, 16, 16, 130, 16, 132, 132, 18, 16, 136, 136, 18, 136, 20, 20, 138,
		128, 32, 32, 130, 32, 132, 132, 34, 32, 136, 136, 34, 136, 36, 36, 138,
		32, 144, 144, 34, 144, 36, 36, 146, 144, 40, 40, 146, 40, 148, 148, 42,
		128, 64, 64, 130, 64, 132, 132, 66, 64, 136, 136, 66, 136, 68, 68, 138,
		64, 144, 144, 66, 144, 68, 68, 146, 144, 72, 72, 146, 72, 148, 148, 74,
		64, 160, 160, 66, 160, 68, 68, 162, 160, 72, 72, 162, 72, 164, 164, 74,
		160, 80, 80, 162, 80, 164, 164, 82, 80, 168, 168, 82, 168, 84, 84, 42,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 11,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 11,
		48, 17, 18, 35, 20, 37, 38, 19, 24, 41, 42, 19, 44, 21, 22, 43,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 11,
		80, 17, 18, 67, 20, 69, 70, 19, 24, 73, 74, 19, 76, 21, 22, 75,
		96, 33, 34, 67, 36, 69, 70, 35, 40, 73, 74, 35, 76, 37, 38, 75,
		48, 81, 82, 35, 84, 37, 38, 83, 88, 41, 42, 83, 44, 85, 86, 43,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 11,
		144, 17, 18, 131, 20, 133, 134, 19, 24, 137, 138, 19, 140, 21, 22, 139,
		160, 33, 34, 131, 36, 133, 134, 35, 40, 137, 138, 35, 140, 37, 38, 139,
		48, 145, 146, 35, 148, 37, 38, 147, 152, 41, 42, 147, 44, 149, 150, 43,
		192, 65, 66, 131, 68, 133, 134, 67, 72, 137, 138, 67, 140, 69, 70, 139,
		80, 145, 146, 67, 148, 69, 70, 147, 152, 73, 74, 147, 76, 149, 150, 75,
		96, 161, 162, 67, 164, 69, 70, 163, 168, 73, 74, 163, 76, 165, 166, 75,
		176, 81, 82, 163, 84, 165, 166, 83, 88, 169, 170, 83, 172, 85, 86, 43,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 12,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 12,
		0, 32, 32, 48, 32, 48, 48, 20, 32, 48, 48, 24, 48, 24, 24, 44,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 12,
		0, 64, 64, 80, 64, 80, 80, 20, 64, 80, 80, 24, 80, 24, 24, 76,
		0, 64, 64, 96, 64, 96, 96, 36, 64, 96, 96, 40, 96, 40, 40, 76,
		64, 96, 96, 48, 96, 48, 48, 84, 96, 48, 48, 88, 48, 88, 88, 44,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 12,
		0, 128, 128, 144, 128, 144, 144, 20, 128, 144, 144, 24, 144, 24, 24, 140,
		0, 128, 128, 160, 128, 160, 160, 36, 128, 160, 160, 40, 160, 40, 40, 140,
		128, 160, 160, 48, 160, 48, 48, 148, 160, 48, 48, 152, 48, 152, 152, 44,
		0, 128, 128, 192, 128, 192, 192, 68, 128, 192, 192, 72, 192, 72, 72, 140,
		128, 192, 192, 80, 192, 80, 80, 148, 192, 80, 80, 152, 80, 152, 152, 76,
		128, 192, 192, 96, 192, 96, 96, 164, 192, 96, 96, 168, 96, 168, 168, 76,
		192, 96, 96, 176, 96, 176, 176, 84, 96, 176, 176, 88, 176, 88, 88, 44,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 13,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 13,
		16, 33, 34, 49, 36, 49, 50, 21, 40, 49, 50, 25, 52, 25, 26, 45,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 13,
		16, 65, 66, 81, 68, 81, 82, 21, 72, 81, 82, 25, 84, 25, 26, 77,
		32, 65, 66, 97, 68, 97, 98, 37, 72, 97, 98, 41, 100, 41, 42, 77,
		80, 97, 98, 49, 100, 49, 50, 85, 104, 49, 50, 89, 52, 89, 90, 45,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 13,
		16, 129, 130, 145, 132, 145, 146, 21, 136, 145, 146, 25, 148, 25, 26, 141,
		32, 129, 130, 161, 132, 161, 162, 37, 136, 161, 162, 41, 164, 41, 42, 141,
		144, 161, 162, 49, 164, 49, 50, 149, 168, 49, 50, 153, 52, 153, 154, 45,
		64, 129, 130, 193, 132, 193, 194, 69, 136, 193, 194, 73, 196, 73, 74, 141,
		144, 193, 194, 81, 196, 81, 82, 149, 200, 81, 82, 153, 84, 153, 154, 77,
		160, 193, 194, 97, 196, 97, 98, 165, 200, 97, 98, 169, 100, 169, 170, 77,
		208, 97, 98, 177, 100, 177, 178, 85, 104, 177, 178, 89, 180, 89, 90, 45,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 14,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 14,
		32, 48, 48, 50, 48, 52, 52, 22, 48, 56, 56, 26, 56, 28, 28, 46,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 14,
		64, 80, 80, 82, 80, 84, 84, 22, 80, 88, 88, 26, 88, 28, 28, 78,
		64, 96, 96, 98, 96, 100, 100, 38, 96, 104, 104, 42, 104, 44, 44, 78,
		96, 112, 112, 50, 112, 52, 52, 86, 112, 56, 56, 90, 56, 92, 92, 46,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 14,
		128, 144, 144, 146, 144, 148, 148, 22, 144, 152, 152, 26, 152, 28, 28, 142,
		128, 160, 160, 162, 160, 164, 164, 38, 160, 168, 168, 42, 168, 44, 44, 142,
		160, 176, 176, 50, 176, 52, 52, 150, 176, 56, 56, 154, 56, 156, 156, 46,
		128, 192, 192, 194, 192, 196, 196, 70, 192, 200, 200, 74, 200, 76, 76, 142,
		192, 208, 208, 82, 208, 84, 84, 150, 208, 88, 88, 154, 88, 156, 156, 78,
		192, 224, 224, 98, 224, 100, 100, 166, 224, 104, 104, 170, 104, 172, 172, 78,
		224, 112, 112, 178, 112, 180, 180, 86, 112, 184, 184, 90, 184, 92, 92, 46,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 15,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 15,
		48, 49, 50, 51, 52, 53, 54, 23, 56, 57, 58, 27, 60, 29, 30, 47,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 15,
		80, 81, 82, 83, 84, 85, 86, 23, 88, 89, 90, 27, 92, 29, 30, 79,
		96, 97, 98, 99, 100, 101, 102, 39, 104, 105, 106, 43, 108, 45, 46, 79,
		112, 113, 114, 51, 116, 53, 54, 87, 120, 57, 58, 91, 60, 93, 94, 47,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 15,
		144, 145, 146, 147, 148, 149, 150, 23, 152, 153, 154, 27, 156, 29, 30, 143,
		160, 161, 162, 163, 164, 165, 166, 39, 168, 169, 170, 43, 172, 45, 46, 143,
		176, 177, 178, 51, 180, 53, 54, 151, 184, 57, 58, 155, 60, 157, 158, 47,
		192, 193, 194, 195, 196, 197, 198, 71, 200, 201, 202, 75, 204, 77, 78, 143,
		208, 209, 210, 83, 212, 85, 86, 151, 216, 89, 90, 155, 92, 157, 158, 79,
		224, 225, 226, 99, 228, 101, 102, 167, 232, 105, 106, 171, 108, 173, 174, 79,
		240, 113, 114, 179, 116, 181, 182, 87, 120, 185, 186, 91, 188, 93, 94, 47,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 48,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 80,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 96,
		0, 0, 0, 64, 0, 64, 64, 96, 0, 64, 64, 96, 64, 96, 96, 48,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 144,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 160,
		0, 0, 0, 128, 0, 128, 128, 160, 0, 128, 128, 160, 128, 160, 160, 48,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 192,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 80,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 96,
		0, 128, 128, 192, 128, 192, 192, 96, 128, 192, 192, 96, 192, 96, 96, 48,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 17,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 33,
		16, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 49,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		16, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 81,
		32, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 97,
		16, 1, 2, 65, 4, 65, 66, 97, 8, 65, 66, 97, 68, 97, 98, 49,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 145,
		32, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 161,
		16, 1, 2, 129, 4, 129, 130, 161, 8, 129, 130, 161, 132, 161, 162, 49,
		64, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 193,
		16, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 81,
		32, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 97,
		16, 129, 130, 193, 132, 193, 194, 97, 136, 193, 194, 97, 196, 97, 98, 49,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 18,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 34,
		32, 16, 16, 2, 16, 4, 4, 34, 16, 8, 8, 34, 8, 36, 36, 50,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 66,
		64, 16, 16, 2, 16, 4, 4, 66, 16, 8, 8, 66, 8, 68, 68, 82,
		64, 32, 32, 2, 32, 4, 4, 66, 32, 8, 8, 66, 8, 68, 68, 98,
		32, 16, 16, 66, 16, 68, 68, 98, 16, 72, 72, 98, 72, 100, 100, 50,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 130,
		128, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 146,
		128, 32, 32, 2, 32, 4, 4, 130, 32, 8, 8, 130, 8, 132, 132, 162,
		32, 16, 16, 130, 16, 132, 132, 162, 16, 136, 136, 162, 136, 164, 164, 50,
		128, 64, 64, 2, 64, 4, 4, 130, 64, 8, 8, 130, 8, 132, 132, 194,
		64, 16, 16, 130, 16, 132, 132, 194, 16, 136, 136, 194, 136, 196, 196, 82,
		64, 32, 32, 130, 32, 132, 132, 194, 32, 136, 136, 194, 136, 196, 196, 98,
		32, 144, 144, 194, 144, 196, 196, 98, 144, 200, 200, 98, 200, 100, 100, 50,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 19,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 35,
		48, 17, 18, 3, 20, 5, 6, 35, 24, 9, 10, 35, 12, 37, 38, 51,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 67,
		80, 17, 18, 3, 20, 5, 6, 67, 24, 9, 10, 67, 12, 69, 70, 83,
		96, 33, 34, 3, 36, 5, 6, 67, 40, 9, 10, 67, 12, 69, 70, 99,
		48, 17, 18, 67, 20, 69, 70, 99, 24, 73, 74, 99, 76, 101, 102, 51,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 131,
		144, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 147,
		160, 33, 34, 3, 36, 5, 6, 131, 40, 9, 10, 131, 12, 133, 134, 163,
		48, 17, 18, 131, 20, 133, 134, 163, 24, 137, 138, 163, 140, 165, 166, 51,
		192, 65, 66, 3, 68, 5, 6, 131, 72, 9, 10, 131, 12, 133, 134, 195,
		80, 17, 18, 131, 20, 133, 134, 195, 24, 137, 138, 195, 140, 197, 198, 83,
		96, 33, 34, 131, 36, 133, 134, 195, 40, 137, 138, 195, 140, 197, 198, 99,
		48, 145, 146, 195, 148, 197, 198, 99, 152, 201, 202, 99, 204, 101, 102, 51,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 20,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 36,
		0, 32, 32, 16, 32, 16, 16, 36, 32, 16, 16, 40, 16, 40, 40, 52,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 68,
		0, 64, 64, 16, 64, 16, 16, 68, 64, 16, 16, 72, 16, 72, 72, 84,
		0, 64, 64, 32, 64, 32, 32, 68, 64, 32, 32, 72, 32, 72, 72, 100,
		64, 32, 32, 80, 32, 80, 80, 100, 32, 80, 80, 104, 80, 104, 104, 52,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 132,
		0, 128, 128, 16, 128, 16, 16, 132, 128, 16, 16, 136, 16, 136, 136, 148,
		0, 128, 128, 32, 128, 32, 32, 132, 128, 32, 32, 136, 32, 136, 136, 164,
		128, 32, 32, 144, 32, 144, 144, 164, 32, 144, 144, 168, 144, 168, 168, 52,
		0, 128, 128, 64, 128, 64, 64, 132, 128, 64, 64, 136, 64, 136, 136, 196,
		128, 64, 64, 144, 64, 144, 144, 196, 64, 144, 144, 200, 144, 200, 200, 84,
		128, 64, 64, 160, 64, 160, 160, 196, 64, 160, 160, 200, 160, 200, 200, 100,
		64, 160, 160, 208, 160, 208, 208, 100, 160, 208, 208, 104, 208, 104, 104, 52,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 21,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 37,
		16, 33, 34, 17, 36, 17, 18, 37, 40, 17, 18, 41, 20, 41, 42, 53,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 69,
		16, 65, 66, 17, 68, 17, 18, 69, 72, 17, 18, 73, 20, 73, 74, 85,
		32, 65, 66, 33, 68, 33, 34, 69, 72, 33, 34, 73, 36, 73, 74, 101,
		80, 33, 34, 81, 36, 81, 82, 101, 40, 81, 82, 105, 84, 105, 106, 53,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 133,
		16, 129, 130, 17, 132, 17, 18, 133, 136, 17, 18, 137, 20, 137, 138, 149,
		32, 129, 130, 33, 132, 33, 34, 133, 136, 33, 34, 137, 36, 137, 138, 165,
		144, 33, 34, 145, 36, 145, 146, 165, 40, 145, 146, 169, 148, 169, 170, 53,
		64, 129, 130, 65, 132, 65, 66, 133, 136, 65, 66, 137, 68, 137, 138, 197,
		144, 65, 66, 145, 68, 145, 146, 197, 72, 145, 146, 201, 148, 201, 202, 85,
		160, 65, 66, 161, 68, 161, 162, 197, 72, 161, 162, 201, 164, 201, 202, 101,
		80, 161, 162, 209, 164, 209, 210, 101, 168, 209, 210, 105, 212, 105, 106, 53,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 22,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 38,
		32, 48, 48, 18, 48, 20, 20, 38, 48, 24, 24, 42, 24, 44, 44, 54,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 70,
		64, 80, 80, 18, 80, 20, 20, 70, 80, 24, 24, 74, 24, 76, 76, 86,
		64, 96, 96, 34, 96, 36, 36, 70, 96, 40, 40, 74, 40, 76, 76, 102,
		96, 48, 48, 82, 48, 84, 84, 102, 48, 88, 88, 106, 88, 108, 108, 54,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 134,
		128, 144, 144, 18, 144, 20, 20, 134, 144, 24, 24, 138, 24, 140, 140, 150,
		128, 160, 160, 34, 160, 36, 36, 134, 160, 40, 40, 138, 40, 140, 140, 166,
		160, 48, 48, 146, 48, 148, 148, 166, 48, 152, 152, 170, 152, 172, 172, 54,
		128, 192, 192, 66, 192, 68, 68, 134, 192, 72, 72, 138, 72, 140, 140, 198,
		192, 80, 80, 146, 80, 148, 148, 198, 80, 152, 152, 202, 152, 204, 204, 86,
		192, 96, 96, 162, 96, 164, 164, 198, 96, 168, 168, 202, 168, 204, 204, 102,
		96, 176, 176, 210, 176, 212, 212, 102, 176, 216, 216, 106, 216, 108, 108, 54,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 23,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 39,
		48, 49, 50, 19, 52, 21, 22, 39, 56, 25, 26, 43, 28, 45, 46, 55,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 71,
		80, 81, 82, 19, 84, 21, 22, 71, 88, 25, 26, 75, 28, 77, 78, 87,
		96, 97, 98, 35, 100, 37, 38, 71, 104, 41, 42, 75, 44, 77, 78, 103,
		112, 49, 50, 83, 52, 85, 86, 103, 56, 89, 90, 107, 92, 109, 110, 55,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 135,
		144, 145, 146, 19, 148, 21, 22, 135, 152, 25, 26, 139, 28, 141, 142, 151,
		160, 161, 162, 35, 164, 37, 38, 135, 168, 41, 42, 139, 44, 141, 142, 167,
		176, 49, 50, 147, 52, 149, 150, 167, 56, 153, 154, 171, 156, 173, 174, 55,
		192, 193, 194, 67, 196, 69, 70, 135, 200, 73, 74, 139, 76, 141, 142, 199,
		208, 81, 82, 147, 84, 149, 150, 199, 88, 153, 154, 203, 156, 205, 206, 87,
		224, 97, 98, 163, 100, 165, 166, 199, 104, 169, 170, 203, 172, 205, 206, 103,
		112, 177, 178, 211, 180, 213, 214, 103, 184, 217, 218, 107, 220, 109, 110, 55,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 24,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 40,
		0, 0, 0, 32, 0, 32, 32, 48, 0, 32, 32, 48, 32, 48, 48, 56,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 72,
		0, 0, 0, 64, 0, 64, 64, 80, 0, 64, 64, 80, 64, 80, 80, 88,
		0, 0, 0, 64, 0, 64, 64, 96, 0, 64, 64, 96, 64, 96, 96, 104,
		0, 64, 64, 96, 64, 96, 96, 112, 64, 96, 96, 112, 96, 112, 112, 56,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 136,
		0, 0, 0, 128, 0, 128, 128, 144, 0, 128, 128, 144, 128, 144, 144, 152,
		0, 0, 0, 128, 0, 128, 128, 160, 0, 128, 128, 160, 128, 160, 160, 168,
		0, 128, 128, 160, 128, 160, 160, 176, 128, 160, 160, 176, 160, 176, 176, 56,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 200,
		0, 128, 128, 192, 128, 192, 192, 208, 128, 192, 192, 208, 192, 208, 208, 88,
		0, 128, 128, 192, 128, 192, 192, 224, 128, 192, 192, 224, 192, 224, 224, 104,
		128, 192, 192, 224, 192, 224, 224, 112, 192, 224, 224, 112, 224, 112, 112, 56,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 25,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 41,
		16, 1, 2, 33, 4, 33, 34, 49, 8, 33, 34, 49, 36, 49, 50, 57,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 73,
		16, 1, 2, 65, 4, 65, 66, 81, 8, 65, 66, 81, 68, 81, 82, 89,
		32, 1, 2, 65, 4, 65, 66, 97, 8, 65, 66, 97, 68, 97, 98, 105,
		16, 65, 66, 97, 68, 97, 98, 113, 72, 97, 98, 113, 100, 113, 114, 57,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 137,
		16, 1, 2, 129, 4, 129, 130, 145, 8, 129, 130, 145, 132, 145, 146, 153,
		32, 1, 2, 129, 4, 129, 130, 161, 8, 129, 130, 161, 132, 161, 162, 169,
		16, 129, 130, 161, 132, 161, 162, 177, 136, 161, 162, 177, 164, 177, 178, 57,
		64, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 201,
		16, 129, 130, 193, 132, 193, 194, 209, 136, 193, 194, 209, 196, 209, 210, 89,
		32, 129, 130, 193, 132, 193, 194, 225, 136, 193, 194, 225, 196, 225, 226, 105,
		144, 193, 194, 225, 196, 225, 226, 113, 200, 225, 226, 113, 228, 113, 114, 57,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 26,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 42,
		32, 16, 16, 34, 16, 36, 36, 50, 16, 40, 40, 50, 40, 52, 52, 58,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 74,
		64, 16, 16, 66, 16, 68, 68, 82, 16, 72, 72, 82, 72, 84, 84, 90,
		64, 32, 32, 66, 32, 68, 68, 98, 32, 72, 72, 98, 72, 100, 100, 106,
		32, 80, 80, 98, 80, 100, 100, 114, 80, 104, 104, 114, 104, 116, 116, 58,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 138,
		128, 16, 16, 130, 16, 132, 132, 146, 16, 136, 136, 146, 136, 148, 148, 154,
		128, 32, 32, 130, 32, 132, 132, 162, 32, 136, 136, 162, 136, 164, 164, 170,
		32, 144, 144, 162, 144, 164, 164, 178, 144, 168, 168, 178, 168, 180, 180, 58,
		128, 64, 64, 130, 64, 132, 132, 194, 64, 136, 136, 194, 136, 196, 196, 202,
		64, 144, 144, 194, 144, 196, 196, 210, 144, 200, 200, 210, 200, 212, 212, 90,
		64, 160, 160, 194, 160, 196, 196, 226, 160, 200, 200, 226, 200, 228, 228, 106,
		160, 208, 208, 226, 208, 228, 228, 114, 208, 232, 232, 114, 232, 116, 116, 58,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 27,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 43,
		48, 17, 18, 35, 20, 37, 38, 51, 24, 41, 42, 51, 44, 53, 54, 59,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 75,
		80, 17, 18, 67, 20, 69, 70, 83, 24, 73, 74, 83, 76, 85, 86, 91,
		96, 33, 34, 67, 36, 69, 70, 99, 40, 73, 74, 99, 76, 101, 102, 107,
		48, 81, 82, 99, 84, 101, 102, 115, 88, 105, 106, 115, 108, 117, 118, 59,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 139,
		144, 17, 18, 131, 20, 133, 134, 147, 24, 137, 138, 147, 140, 149, 150, 155,
		160, 33, 34, 131, 36, 133, 134, 163, 40, 137, 138, 163, 140, 165, 166, 171,
		48, 145, 146, 163, 148, 165, 166, 179, 152, 169, 170, 179, 172, 181, 182, 59,
		192, 65, 66, 131, 68, 133, 134, 195, 72, 137, 138, 195, 140, 197, 198, 203,
		80, 145, 146, 195, 148, 197, 198, 211, 152, 201, 202, 211, 204, 213, 214, 91,
		96, 161, 162, 195, 164, 197, 198, 227, 168, 201, 202, 227, 204, 229, 230, 107,
		176, 209, 210, 227, 212, 229, 230, 115, 216, 233, 234, 115, 236, 117, 118, 59,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 28,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 44,
		0, 32, 32, 48, 32, 48, 48, 52, 32, 48, 48, 56, 48, 56, 56, 60,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 76,
		0, 64, 64, 80, 64, 80, 80, 84, 64, 80, 80, 88, 80, 88, 88, 92,
		0, 64, 64, 96, 64, 96, 96, 100, 64, 96, 96, 104, 96, 104, 104, 108,
		64, 96, 96, 112, 96, 112, 112, 116, 96, 112, 112, 120, 112, 120, 120, 60,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 140,
		0, 128, 128, 144, 128, 144, 144, 148, 128, 144, 144, 152, 144, 152, 152, 156,
		0, 128, 128, 160, 128, 160, 160, 164, 128, 160, 160, 168, 160, 168, 168, 172,
		128, 160, 160, 176, 160, 176, 176, 180, 160, 176, 176, 184, 176, 184, 184, 60,
		0, 128, 128, 192, 128, 192, 192, 196, 128, 192, 192, 200, 192, 200, 200, 204,
		128, 192, 192, 208, 192, 208, 208, 212, 192, 208, 208, 216, 208, 216, 216, 92,
		128, 192, 192, 224, 192, 224, 224, 228, 192, 224, 224, 232, 224, 232, 232, 108,
		192, 224, 224, 240, 224, 240, 240, 116, 224, 240, 240, 120, 240, 120, 120, 60,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 29,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 45,
		16, 33, 34, 49, 36, 49, 50, 53, 40, 49, 50, 57, 52, 57, 58, 61,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 77,
		16, 65, 66, 81, 68, 81, 82, 85, 72, 81, 82, 89, 84, 89, 90, 93,
		32, 65, 66, 97, 68, 97, 98, 101, 72, 97, 98, 105, 100, 105, 106, 109,
		80, 97, 98, 113, 100, 113, 114, 117, 104, 113, 114, 121, 116, 121, 122, 61,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 141,
		16, 129, 130, 145, 132, 145, 146, 149, 136, 145, 146, 153, 148, 153, 154, 157,
		32, 129, 130, 161, 132, 161, 162, 165, 136, 161, 162, 169, 164, 169, 170, 173,
		144, 161, 162, 177, 164, 177, 178, 181, 168, 177, 178, 185, 180, 185, 186, 61,
		64, 129, 130, 193, 132, 193, 194, 197, 136, 193, 194, 201, 196, 201, 202, 205,
		144, 193, 194, 209, 196, 209, 210, 213, 200, 209, 210, 217, 212, 217, 218, 93,
		160, 193, 194, 225, 196, 225, 226, 229, 200, 225, 226, 233, 228, 233, 234, 109,
		208, 225, 226, 241, 228, 241, 242, 117, 232, 241, 242, 121, 244, 121, 122, 61,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 30,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 46,
		32, 48, 48, 50, 48, 52, 52, 54, 48, 56, 56, 58, 56, 60, 60, 62,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 78,
		64, 80, 80, 82, 80, 84, 84, 86, 80, 88, 88, 90, 88, 92, 92, 94,
		64, 96, 96, 98, 96, 100, 100, 102, 96, 104, 104, 106, 104, 108, 108, 110,
		96, 112, 112, 114, 112, 116, 116, 118, 112, 120, 120, 122, 120, 124, 124, 62,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 142,
		128, 144, 144, 146, 144, 148, 148, 150, 144, 152, 152, 154, 152, 156, 156, 158,
		128, 160, 160, 162, 160, 164, 164, 166, 160, 168, 168, 170, 168, 172, 172, 174,
		160, 176, 176, 178, 176, 180, 180, 182, 176, 184, 184, 186, 184, 188, 188, 62,
		128, 192, 192, 194, 192, 196, 196, 198, 192, 200, 200, 202, 200, 204, 204, 206,
		192, 208, 208, 210, 208, 212, 212, 214, 208, 216, 216, 218, 216, 220, 220, 94,
		192, 224, 224, 226, 224, 228, 228, 230, 224, 232, 232, 234, 232, 236, 236, 110,
		224, 240, 240, 242, 240, 244, 244, 118, 240, 248, 248, 122, 248, 124, 124, 62,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
		80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95,
		96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111,
		112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 63,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143,
		144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159,
		160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175,
		176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 63,
		192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207,
		208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 95,
		224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 111,
		240, 241, 242, 243, 244, 245, 246, 119, 248, 249, 250, 123, 252, 125, 126, 63,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 64,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 65,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 2,
		64, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		64, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 66,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 2,
		128, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		128, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 130,
		128, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 2,
		64, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 130,
		64, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 130,
		32, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 66,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 3,
		80, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		96, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 67,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 3,
		144, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		160, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 131,
		192, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 3,
		80, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 131,
		96, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 131,
		48, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 67,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 4,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 4,
		0, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 4,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 4,
		0, 64, 64, 16, 64, 16, 16, 4, 64, 16, 16, 8, 16, 8, 8, 4,
		0, 64, 64, 32, 64, 32, 32, 4, 64, 32, 32, 8, 32, 8, 8, 4,
		64, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 68,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 4,
		0, 128, 128, 16, 128, 16, 16, 4, 128, 16, 16, 8, 16, 8, 8, 4,
		0, 128, 128, 32, 128, 32, 32, 4, 128, 32, 32, 8, 32, 8, 8, 4,
		128, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 132,
		0, 128, 128, 64, 128, 64, 64, 4, 128, 64, 64, 8, 64, 8, 8, 4,
		128, 64, 64, 16, 64, 16, 16, 4, 64, 16, 16, 8, 16, 8, 8, 132,
		128, 64, 64, 32, 64, 32, 32, 4, 64, 32, 32, 8, 32, 8, 8, 132,
		64, 32, 32, 16, 32, 16, 16, 132, 32, 16, 16, 136, 16, 136, 136, 68,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 5,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 5,
		16, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 5,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 5,
		16, 65, 66, 17, 68, 17, 18, 5, 72, 17, 18, 9, 20, 9, 10, 5,
		32, 65, 66, 33, 68, 33, 34, 5, 72, 33, 34, 9, 36, 9, 10, 5,
		80, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 69,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 5,
		16, 129, 130, 17, 132, 17, 18, 5, 136, 17, 18, 9, 20, 9, 10, 5,
		32, 129, 130, 33, 132, 33, 34, 5, 136, 33, 34, 9, 36, 9, 10, 5,
		144, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 133,
		64, 129, 130, 65, 132, 65, 66, 5, 136, 65, 66, 9, 68, 9, 10, 5,
		144, 65, 66, 17, 68, 17, 18, 5, 72, 17, 18, 9, 20, 9, 10, 133,
		160, 65, 66, 33, 68, 33, 34, 5, 72, 33, 34, 9, 36, 9, 10, 133,
		80, 33, 34, 17, 36, 17, 18, 133, 40, 17, 18, 137, 20, 137, 138, 69,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 6,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 6,
		32, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 6,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 6,
		64, 80, 80, 18, 80, 20, 20, 6, 80, 24, 24, 10, 24, 12, 12, 6,
		64, 96, 96, 34, 96, 36, 36, 6, 96, 40, 40, 10, 40, 12, 12, 6,
		96, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 70,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 6,
		128, 144, 144, 18, 144, 20, 20, 6, 144, 24, 24, 10, 24, 12, 12, 6,
		128, 160, 160, 34, 160, 36, 36, 6, 160, 40, 40, 10, 40, 12, 12, 6,
		160, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 134,
		128, 192, 192, 66, 192, 68, 68, 6, 192, 72, 72, 10, 72, 12, 12, 6,
		192, 80, 80, 18, 80, 20, 20, 6, 80, 24, 24, 10, 24, 12, 12, 134,
		192, 96, 96, 34, 96, 36, 36, 6, 96, 40, 40, 10, 40, 12, 12, 134,
		96, 48, 48, 18, 48, 20, 20, 134, 48, 24, 24, 138, 24, 140, 140, 70,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 7,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 7,
		48, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 7,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 7,
		80, 81, 82, 19, 84, 21, 22, 7, 88, 25, 26, 11, 28, 13, 14, 7,
		96, 97, 98, 35, 100, 37, 38, 7, 104, 41, 42, 11, 44, 13, 14, 7,
		112, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 71,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 7,
		144, 145, 146, 19, 148, 21, 22, 7, 152, 25, 26, 11, 28, 13, 14, 7,
		160, 161, 162, 35, 164, 37, 38, 7, 168, 41, 42, 11, 44, 13, 14, 7,
		176, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 135,
		192, 193, 194, 67, 196, 69, 70, 7, 200, 73, 74, 11, 76, 13, 14, 7,
		208, 81, 82, 19, 84, 21, 22, 7, 88, 25, 26, 11, 28, 13, 14, 135,
		224, 97, 98, 35, 100, 37, 38, 7, 104, 41, 42, 11, 44, 13, 14, 135,
		112, 49, 50, 19, 52, 21, 22, 135, 56, 25, 26, 139, 28, 141, 142, 71,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 8,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 8,
		0, 0, 0, 32, 0, 32, 32, 16, 0, 32, 32, 16, 32, 16, 16, 8,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 8,
		0, 0, 0, 64, 0, 64, 64, 16, 0, 64, 64, 16, 64, 16, 16, 8,
		0, 0, 0, 64, 0, 64, 64, 32, 0, 64, 64, 32, 64, 32, 32, 8,
		0, 64, 64, 32, 64, 32, 32, 16, 64, 32, 32, 16, 32, 16, 16, 72,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 8,
		0, 0, 0, 128, 0, 128, 128, 16, 0, 128, 128, 16, 128, 16, 16, 8,
		0, 0, 0, 128, 0, 128, 128, 32, 0, 128, 128, 32, 128, 32, 32, 8,
		0, 128, 128, 32, 128, 32, 32, 16, 128, 32, 32, 16, 32, 16, 16, 136,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 8,
		0, 128, 128, 64, 128, 64, 64, 16, 128, 64, 64, 16, 64, 16, 16, 136,
		0, 128, 128, 64, 128, 64, 64, 32, 128, 64, 64, 32, 64, 32, 32, 136,
		128, 64, 64, 32, 64, 32, 32, 144, 64, 32, 32, 144, 32, 144, 144, 72,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 9,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 9,
		16, 1, 2, 33, 4, 33, 34, 17, 8, 33, 34, 17, 36, 17, 18, 9,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 9,
		16, 1, 2, 65, 4, 65, 66, 17, 8, 65, 66, 17, 68, 17, 18, 9,
		32, 1, 2, 65, 4, 65, 66, 33, 8, 65, 66, 33, 68, 33, 34, 9,
		16, 65, 66, 33, 68, 33, 34, 17, 72, 33, 34, 17, 36, 17, 18, 73,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 9,
		16, 1, 2, 129, 4, 129, 130, 17, 8, 129, 130, 17, 132, 17, 18, 9,
		32, 1, 2, 129, 4, 129, 130, 33, 8, 129, 130, 33, 132, 33, 34, 9,
		16, 129, 130, 33, 132, 33, 34, 17, 136, 33, 34, 17, 36, 17, 18, 137,
		64, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 9,
		16, 129, 130, 65, 132, 65, 66, 17, 136, 65, 66, 17, 68, 17, 18, 137,
		32, 129, 130, 65, 132, 65, 66, 33, 136, 65, 66, 33, 68, 33, 34, 137,
		144, 65, 66, 33, 68, 33, 34, 145, 72, 33, 34, 145, 36, 145, 146, 73,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 10,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 10,
		32, 16, 16, 34, 16, 36, 36, 18, 16, 40, 40, 18, 40, 20, 20, 10,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 10,
		64, 16, 16, 66, 16, 68, 68, 18, 16, 72, 72, 18, 72, 20, 20, 10,
		64, 32, 32, 66, 32, 68, 68, 34, 32, 72, 72, 34, 72, 36, 36, 10,
		32, 80, 80, 34, 80, 36, 36, 18, 80, 40, 40, 18, 40, 20, 20, 74,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 10,
		128, 16, 16, 130, 16, 132, 132, 18, 16, 136, 136, 18, 136, 20, 20, 10,
		128, 32, 32, 130, 32, 132, 132, 34, 32, 136, 136, 34, 136, 36, 36, 10,
		32, 144, 144, 34, 144, 36, 36, 18, 144, 40, 40, 18, 40, 20, 20, 138,
		128, 64, 64, 130, 64, 132, 132, 66, 64, 136, 136, 66, 136, 68, 68, 10,
		64, 144, 144, 66, 144, 68, 68, 18, 144, 72, 72, 18, 72, 20, 20, 138,
		64, 160, 160, 66, 160, 68, 68, 34, 160, 72, 72, 34, 72, 36, 36, 138,
		160, 80, 80, 34, 80, 36, 36, 146, 80, 40, 40, 146, 40, 148, 148, 74,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 11,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 11,
		48, 17, 18, 35, 20, 37, 38, 19, 24, 41, 42, 19, 44, 21, 22, 11,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 11,
		80, 17, 18, 67, 20, 69, 70, 19, 24, 73, 74, 19, 76, 21, 22, 11,
		96, 33, 34, 67, 36, 69, 70, 35, 40, 73, 74, 35, 76, 37, 38, 11,
		48, 81, 82, 35, 84, 37, 38, 19, 88, 41, 42, 19, 44, 21, 22, 75,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 11,
		144, 17, 18, 131, 20, 133, 134, 19, 24, 137, 138, 19, 140, 21, 22, 11,
		160, 33, 34, 131, 36, 133, 134, 35, 40, 137, 138, 35, 140, 37, 38, 11,
		48, 145, 146, 35, 148, 37, 38, 19, 152, 41, 42, 19, 44, 21, 22, 139,
		192, 65, 66, 131, 68, 133, 134, 67, 72, 137, 138, 67, 140, 69, 70, 11,
		80, 145, 146, 67, 148, 69, 70, 19, 152, 73, 74, 19, 76, 21, 22, 139,
		96, 161, 162, 67, 164, 69, 70, 35, 168, 73, 74, 35, 76, 37, 38, 139,
		176, 81, 82, 35, 84, 37, 38, 147, 88, 41, 42, 147, 44, 149, 150, 75,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 12,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 12,
		0, 32, 32, 48, 32, 48, 48, 20, 32, 48, 48, 24, 48, 24, 24, 12,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 12,
		0, 64, 64, 80, 64, 80, 80, 20, 64, 80, 80, 24, 80, 24, 24, 12,
		0, 64, 64, 96, 64, 96, 96, 36, 64, 96, 96, 40, 96, 40, 40, 12,
		64, 96, 96, 48, 96, 48, 48, 20, 96, 48, 48, 24, 48, 24, 24, 76,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 12,
		0, 128, 128, 144, 128, 144, 144, 20, 128, 144, 144, 24, 144, 24, 24, 12,
		0, 128, 128, 160, 128, 160, 160, 36, 128, 160, 160, 40, 160, 40, 40, 12,
		128, 160, 160, 48, 160, 48, 48, 20, 160, 48, 48, 24, 48, 24, 24, 140,
		0, 128, 128, 192, 128, 192, 192, 68, 128, 192, 192, 72, 192, 72, 72, 12,
		128, 192, 192, 80, 192, 80, 80, 20, 192, 80, 80, 24, 80, 24, 24, 140,
		128, 192, 192, 96, 192, 96, 96, 36, 192, 96, 96, 40, 96, 40, 40, 140,
		192, 96, 96, 48, 96, 48, 48, 148, 96, 48, 48, 152, 48, 152, 152, 76,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 13,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 13,
		16, 33, 34, 49, 36, 49, 50, 21, 40, 49, 50, 25, 52, 25, 26, 13,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 13,
		16, 65, 66, 81, 68, 81, 82, 21, 72, 81, 82, 25, 84, 25, 26, 13,
		32, 65, 66, 97, 68, 97, 98, 37, 72, 97, 98, 41, 100, 41, 42, 13,
		80, 97, 98, 49, 100, 49, 50, 21, 104, 49, 50, 25, 52, 25, 26, 77,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 13,
		16, 129, 130, 145, 132, 145, 146, 21, 136, 145, 146, 25, 148, 25, 26, 13,
		32, 129, 130, 161, 132, 161, 162, 37, 136, 161, 162, 41, 164, 41, 42, 13,
		144, 161, 162, 49, 164, 49, 50, 21, 168, 49, 50, 25, 52, 25, 26, 141,
		64, 129, 130, 193, 132, 193, 194, 69, 136, 193, 194, 73, 196, 73, 74, 13,
		144, 193, 194, 81, 196, 81, 82, 21, 200, 81, 82, 25, 84, 25, 26, 141,
		160, 193, 194, 97, 196, 97, 98, 37, 200, 97, 98, 41, 100, 41, 42, 141,
		208, 97, 98, 49, 100, 49, 50, 149, 104, 49, 50, 153, 52, 153, 154, 77,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 14,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 14,
		32, 48, 48, 50, 48, 52, 52, 22, 48, 56, 56, 26, 56, 28, 28, 14,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 14,
		64, 80, 80, 82, 80, 84, 84, 22, 80, 88, 88, 26, 88, 28, 28, 14,
		64, 96, 96, 98, 96, 100, 100, 38, 96, 104, 104, 42, 104, 44, 44, 14,
		96, 112, 112, 50, 112, 52, 52, 22, 112, 56, 56, 26, 56, 28, 28, 78,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 14,
		128, 144, 144, 146, 144, 148, 148, 22, 144, 152, 152, 26, 152, 28, 28, 14,
		128, 160, 160, 162, 160, 164, 164, 38, 160, 168, 168, 42, 168, 44, 44, 14,
		160, 176, 176, 50, 176, 52, 52, 22, 176, 56, 56, 26, 56, 28, 28, 142,
		128, 192, 192, 194, 192, 196, 196, 70, 192, 200, 200, 74, 200, 76, 76, 14,
		192, 208, 208, 82, 208, 84, 84, 22, 208, 88, 88, 26, 88, 28, 28, 142,
		192, 224, 224, 98, 224, 100, 100, 38, 224, 104, 104, 42, 104, 44, 44, 142,
		224, 112, 112, 50, 112, 52, 52, 150, 112, 56, 56, 154, 56, 156, 156, 78,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 15,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 15,
		48, 49, 50, 51, 52, 53, 54, 23, 56, 57, 58, 27, 60, 29, 30, 15,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 15,
		80, 81, 82, 83, 84, 85, 86, 23, 88, 89, 90, 27, 92, 29, 30, 15,
		96, 97, 98, 99, 100, 101, 102, 39, 104, 105, 106, 43, 108, 45, 46, 15,
		112, 113, 114, 51, 116, 53, 54, 23, 120, 57, 58, 27, 60, 29, 30, 79,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 15,
		144, 145, 146, 147, 148, 149, 150, 23, 152, 153, 154, 27, 156, 29, 30, 15,
		160, 161, 162, 163, 164, 165, 166, 39, 168, 169, 170, 43, 172, 45, 46, 15,
		176, 177, 178, 51, 180, 53, 54, 23, 184, 57, 58, 27, 60, 29, 30, 143,
		192, 193, 194, 195, 196, 197, 198, 71, 200, 201, 202, 75, 204, 77, 78, 15,
		208, 209, 210, 83, 212, 85, 86, 23, 216, 89, 90, 27, 92, 29, 30, 143,
		224, 225, 226, 99, 228, 101, 102, 39, 232, 105, 106, 43, 108, 45, 46, 143,
		240, 113, 114, 51, 116, 53, 54, 151, 120, 57, 58, 155, 60, 157, 158, 79,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 16,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 32,
		0, 0, 0, 64, 0, 64, 64, 32, 0, 64, 64, 32, 64, 32, 32, 80,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 16,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 32,
		0, 0, 0, 128, 0, 128, 128, 32, 0, 128, 128, 32, 128, 32, 32, 144,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 64,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 144,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 160,
		0, 128, 128, 64, 128, 64, 64, 160, 128, 64, 64, 160, 64, 160, 160, 80,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 17,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 33,
		16, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 17,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		16, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 17,
		32, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 33,
		16, 1, 2, 65, 4, 65, 66, 33, 8, 65, 66, 33, 68, 33, 34, 81,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 17,
		32, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 33,
		16, 1, 2, 129, 4, 129, 130, 33, 8, 129, 130, 33, 132, 33, 34, 145,
		64, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 65,
		16, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 145,
		32, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 161,
		16, 129, 130, 65, 132, 65, 66, 161, 136, 65, 66, 161, 68, 161, 162, 81,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 18,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 34,
		32, 16, 16, 2, 16, 4, 4, 34, 16, 8, 8, 34, 8, 36, 36, 18,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 66,
		64, 16, 16, 2, 16, 4, 4, 66, 16, 8, 8, 66, 8, 68, 68, 18,
		64, 32, 32, 2, 32, 4, 4, 66, 32, 8, 8, 66, 8, 68, 68, 34,
		32, 16, 16, 66, 16, 68, 68, 34, 16, 72, 72, 34, 72, 36, 36, 82,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 130,
		128, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 18,
		128, 32, 32, 2, 32, 4, 4, 130, 32, 8, 8, 130, 8, 132, 132, 34,
		32, 16, 16, 130, 16, 132, 132, 34, 16, 136, 136, 34, 136, 36, 36, 146,
		128, 64, 64, 2, 64, 4, 4, 130, 64, 8, 8, 130, 8, 132, 132, 66,
		64, 16, 16, 130, 16, 132, 132, 66, 16, 136, 136, 66, 136, 68, 68, 146,
		64, 32, 32, 130, 32, 132, 132, 66, 32, 136, 136, 66, 136, 68, 68, 162,
		32, 144, 144, 66, 144, 68, 68, 162, 144, 72, 72, 162, 72, 164, 164, 82,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 19,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 35,
		48, 17, 18, 3, 20, 5, 6, 35, 24, 9, 10, 35, 12, 37, 38, 19,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 67,
		80, 17, 18, 3, 20, 5, 6, 67, 24, 9, 10, 67, 12, 69, 70, 19,
		96, 33, 34, 3, 36, 5, 6, 67, 40, 9, 10, 67, 12, 69, 70, 35,
		48, 17, 18, 67, 20, 69, 70, 35, 24, 73, 74, 35, 76, 37, 38, 83,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 131,
		144, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 19,
		160, 33, 34, 3, 36, 5, 6, 131, 40, 9, 10, 131, 12, 133, 134, 35,
		48, 17, 18, 131, 20, 133, 134, 35, 24, 137, 138, 35, 140, 37, 38, 147,
		192, 65, 66, 3, 68, 5, 6, 131, 72, 9, 10, 131, 12, 133, 134, 67,
		80, 17, 18, 131, 20, 133, 134, 67, 24, 137, 138, 67, 140, 69, 70, 147,
		96, 33, 34, 131, 36, 133, 134, 67, 40, 137, 138, 67, 140, 69, 70, 163,
		48, 145, 146, 67, 148, 69, 70, 163, 152, 73, 74, 163, 76, 165, 166, 83,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 20,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 36,
		0, 32, 32, 16, 32, 16, 16, 36, 32, 16, 16, 40, 16, 40, 40, 20,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 68,
		0, 64, 64, 16, 64, 16, 16, 68, 64, 16, 16, 72, 16, 72, 72, 20,
		0, 64, 64, 32, 64, 32, 32, 68, 64, 32, 32, 72, 32, 72, 72, 36,
		64, 32, 32, 80, 32, 80, 80, 36, 32, 80, 80, 40, 80, 40, 40, 84,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 132,
		0, 128, 128, 16, 128, 16, 16, 132, 128, 16, 16, 136, 16, 136, 136, 20,
		0, 128, 128, 32, 128, 32, 32, 132, 128, 32, 32, 136, 32, 136, 136, 36,
		128, 32, 32, 144, 32, 144, 144, 36, 32, 144, 144, 40, 144, 40, 40, 148,
		0, 128, 128, 64, 128, 64, 64, 132, 128, 64, 64, 136, 64, 136, 136, 68,
		128, 64, 64, 144, 64, 144, 144, 68, 64, 144, 144, 72, 144, 72, 72, 148,
		128, 64, 64, 160, 64, 160, 160, 68, 64, 160, 160, 72, 160, 72, 72, 164,
		64, 160, 160, 80, 160, 80, 80, 164, 160, 80, 80, 168, 80, 168, 168, 84,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 21,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 37,
		16, 33, 34, 17, 36, 17, 18, 37, 40, 17, 18, 41, 20, 41, 42, 21,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 69,
		16, 65, 66, 17, 68, 17, 18, 69, 72, 17, 18, 73, 20, 73, 74, 21,
		32, 65, 66, 33, 68, 33, 34, 69, 72, 33, 34, 73, 36, 73, 74, 37,
		80, 33, 34, 81, 36, 81, 82, 37, 40, 81, 82, 41, 84, 41, 42, 85,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 133,
		16, 129, 130, 17, 132, 17, 18, 133, 136, 17, 18, 137, 20, 137, 138, 21,
		32, 129, 130, 33, 132, 33, 34, 133, 136, 33, 34, 137, 36, 137, 138, 37,
		144, 33, 34, 145, 36, 145, 146, 37, 40, 145, 146, 41, 148, 41, 42, 149,
		64, 129, 130, 65, 132, 65, 66, 133, 136, 65, 66, 137, 68, 137, 138, 69,
		144, 65, 66, 145, 68, 145, 146, 69, 72, 145, 146, 73, 148, 73, 74, 149,
		160, 65, 66, 161, 68, 161, 162, 69, 72, 161, 162, 73, 164, 73, 74, 165,
		80, 161, 162, 81, 164, 81, 82, 165, 168, 81, 82, 169, 84, 169, 170, 85,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 22,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 38,
		32, 48, 48, 18, 48, 20, 20, 38, 48, 24, 24, 42, 24, 44, 44, 22,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 70,
		64, 80, 80, 18, 80, 20, 20, 70, 80, 24, 24, 74, 24, 76, 76, 22,
		64, 96, 96, 34, 96, 36, 36, 70, 96, 40, 40, 74, 40, 76, 76, 38,
		96, 48, 48, 82, 48, 84, 84, 38, 48, 88, 88, 42, 88, 44, 44, 86,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 134,
		128, 144, 144, 18, 144, 20, 20, 134, 144, 24, 24, 138, 24, 140, 140, 22,
		128, 160, 160, 34, 160, 36, 36, 134, 160, 40, 40, 138, 40, 140, 140, 38,
		160, 48, 48, 146, 48, 148, 148, 38, 48, 152, 152, 42, 152, 44, 44, 150,
		128, 192, 192, 66, 192, 68, 68, 134, 192, 72, 72, 138, 72, 140, 140, 70,
		192, 80, 80, 146, 80, 148, 148, 70, 80, 152, 152, 74, 152, 76, 76, 150,
		192, 96, 96, 162, 96, 164, 164, 70, 96, 168, 168, 74, 168, 76, 76, 166,
		96, 176, 176, 82, 176, 84, 84, 166, 176, 88, 88, 170, 88, 172, 172, 86,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 23,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 39,
		48, 49, 50, 19, 52, 21, 22, 39, 56, 25, 26, 43, 28, 45, 46, 23,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 71,
		80, 81, 82, 19, 84, 21, 22, 71, 88, 25, 26, 75, 28, 77, 78, 23,
		96, 97, 98, 35, 100, 37, 38, 71, 104, 41, 42, 75, 44, 77, 78, 39,
		112, 49, 50, 83, 52, 85, 86, 39, 56, 89, 90, 43, 92, 45, 46, 87,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 135,
		144, 145, 146, 19, 148, 21, 22, 135, 152, 25, 26, 139, 28, 141, 142, 23,
		160, 161, 162, 35, 164, 37, 38, 135, 168, 41, 42, 139, 44, 141, 142, 39,
		176, 49, 50, 147, 52, 149, 150, 39, 56, 153, 154, 43, 156, 45, 46, 151,
		192, 193, 194, 67, 196, 69, 70, 135, 200, 73, 74, 139, 76, 141, 142, 71,
		208, 81, 82, 147, 84, 149, 150, 71, 88, 153, 154, 75, 156, 77, 78, 151,
		224, 97, 98, 163, 100, 165, 166, 71, 104, 169, 170, 75, 172, 77, 78, 167,
		112, 177, 178, 83, 180, 85, 86, 167, 184, 89, 90, 171, 92, 173, 174, 87,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 24,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 40,
		0, 0, 0, 32, 0, 32, 32, 48, 0, 32, 32, 48, 32, 48, 48, 24,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 72,
		0, 0, 0, 64, 0, 64, 64, 80, 0, 64, 64, 80, 64, 80, 80, 24,
		0, 0, 0, 64, 0, 64, 64, 96, 0, 64, 64, 96, 64, 96, 96, 40,
		0, 64, 64, 96, 64, 96, 96, 48, 64, 96, 96, 48, 96, 48, 48, 88,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 136,
		0, 0, 0, 128, 0, 128, 128, 144, 0, 128, 128, 144, 128, 144, 144, 24,
		0, 0, 0, 128, 0, 128, 128, 160, 0, 128, 128, 160, 128, 160, 160, 40,
		0, 128, 128, 160, 128, 160, 160, 48, 128, 160, 160, 48, 160, 48, 48, 152,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 72,
		0, 128, 128, 192, 128, 192, 192, 80, 128, 192, 192, 80, 192, 80, 80, 152,
		0, 128, 128, 192, 128, 192, 192, 96, 128, 192, 192, 96, 192, 96, 96, 168,
		128, 192, 192, 96, 192, 96, 96, 176, 192, 96, 96, 176, 96, 176, 176, 88,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 25,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 41,
		16, 1, 2, 33, 4, 33, 34, 49, 8, 33, 34, 49, 36, 49, 50, 25,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 73,
		16, 1, 2, 65, 4, 65, 66, 81, 8, 65, 66, 81, 68, 81, 82, 25,
		32, 1, 2, 65, 4, 65, 66, 97, 8, 65, 66, 97, 68, 97, 98, 41,
		16, 65, 66, 97, 68, 97, 98, 49, 72, 97, 98, 49, 100, 49, 50, 89,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 137,
		16, 1, 2, 129, 4, 129, 130, 145, 8, 129, 130, 145, 132, 145, 146, 25,
		32, 1, 2, 129, 4, 129, 130, 161, 8, 129, 130, 161, 132, 161, 162, 41,
		16, 129, 130, 161, 132, 161, 162, 49, 136, 161, 162, 49, 164, 49, 50, 153,
		64, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 73,
		16, 129, 130, 193, 132, 193, 194, 81, 136, 193, 194, 81, 196, 81, 82, 153,
		32, 129, 130, 193, 132, 193, 194, 97, 136, 193, 194, 97, 196, 97, 98, 169,
		144, 193, 194, 97, 196, 97, 98, 177, 200, 97, 98, 177, 100, 177, 178, 89,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 26,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 42,
		32, 16, 16, 34, 16, 36, 36, 50, 16, 40, 40, 50, 40, 52, 52, 26,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 74,
		64, 16, 16, 66, 16, 68, 68, 82, 16, 72, 72, 82, 72, 84, 84, 26,
		64, 32, 32, 66, 32, 68, 68, 98, 32, 72, 72, 98, 72, 100, 100, 42,
		32, 80, 80, 98, 80, 100, 100, 50, 80, 104, 104, 50, 104, 52, 52, 90,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 138,
		128, 16, 16, 130, 16, 132, 132, 146, 16, 136, 136, 146, 136, 148, 148, 26,
		128, 32, 32, 130, 32, 132, 132, 162, 32, 136, 136, 162, 136, 164, 164, 42,
		32, 144, 144, 162, 144, 164, 164, 50, 144, 168, 168, 50, 168, 52, 52, 154,
		128, 64, 64, 130, 64, 132, 132, 194, 64, 136, 136, 194, 136, 196, 196, 74,
		64, 144, 144, 194, 144, 196, 196, 82, 144, 200, 200, 82, 200, 84, 84, 154,
		64, 160, 160, 194, 160, 196, 196, 98, 160, 200, 200, 98, 200, 100, 100, 170,
		160, 208, 208, 98, 208, 100, 100, 178, 208, 104, 104, 178, 104, 180, 180, 90,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 27,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 43,
		48, 17, 18, 35, 20, 37, 38, 51, 24, 41, 42, 51, 44, 53, 54, 27,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 75,
		80, 17, 18, 67, 20, 69, 70, 83, 24, 73, 74, 83, 76, 85, 86, 27,
		96, 33, 34, 67, 36, 69, 70, 99, 40, 73, 74, 99, 76, 101, 102, 43,
		48, 81, 82, 99, 84, 101, 102, 51, 88, 105, 106, 51, 108, 53, 54, 91,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 139,
		144, 17, 18, 131, 20, 133, 134, 147, 24, 137, 138, 147, 140, 149, 150, 27,
		160, 33, 34, 131, 36, 133, 134, 163, 40, 137, 138, 163, 140, 165, 166, 43,
		48, 145, 146, 163, 148, 165, 166, 51, 152, 169, 170, 51, 172, 53, 54, 155,
		192, 65, 66, 131, 68, 133, 134, 195, 72, 137, 138, 195, 140, 197, 198, 75,
		80, 145, 146, 195, 148, 197, 198, 83, 152, 201, 202, 83, 204, 85, 86, 155,
		96, 161, 162, 195, 164, 197, 198, 99, 168, 201, 202, 99, 204, 101, 102, 171,
		176, 209, 210, 99, 212, 101, 102, 179, 216, 105, 106, 179, 108, 181, 182, 91,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 28,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 44,
		0, 32, 32, 48, 32, 48, 48, 52, 32, 48, 48, 56, 48, 56, 56, 28,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 76,
		0, 64, 64, 80, 64, 80, 80, 84, 64, 80, 80, 88, 80, 88, 88, 28,
		0, 64, 64, 96, 64, 96, 96, 100, 64, 96, 96, 104, 96, 104, 104, 44,
		64, 96, 96, 112, 96, 112, 112, 52, 96, 112, 112, 56, 112, 56, 56, 92,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 140,
		0, 128, 128, 144, 128, 144, 144, 148, 128, 144, 144, 152, 144, 152, 152, 28,
		0, 128, 128, 160, 128, 160, 160, 164, 128, 160, 160, 168, 160, 168, 168, 44,
		128, 160, 160, 176, 160, 176, 176, 52, 160, 176, 176, 56, 176, 56, 56, 156,
		0, 128, 128, 192, 128, 192, 192, 196, 128, 192, 192, 200, 192, 200, 200, 76,
		128, 192, 192, 208, 192, 208, 208, 84, 192, 208, 208, 88, 208, 88, 88, 156,
		128, 192, 192, 224, 192, 224, 224, 100, 192, 224, 224, 104, 224, 104, 104, 172,
		192, 224, 224, 112, 224, 112, 112, 180, 224, 112, 112, 184, 112, 184, 184, 92,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 29,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 45,
		16, 33, 34, 49, 36, 49, 50, 53, 40, 49, 50, 57, 52, 57, 58, 29,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 77,
		16, 65, 66, 81, 68, 81, 82, 85, 72, 81, 82, 89, 84, 89, 90, 29,
		32, 65, 66, 97, 68, 97, 98, 101, 72, 97, 98, 105, 100, 105, 106, 45,
		80, 97, 98, 113, 100, 113, 114, 53, 104, 113, 114, 57, 116, 57, 58, 93,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 141,
		16, 129, 130, 145, 132, 145, 146, 149, 136, 145, 146, 153, 148, 153, 154, 29,
		32, 129, 130, 161, 132, 161, 162, 165, 136, 161, 162, 169, 164, 169, 170, 45,
		144, 161, 162, 177, 164, 177, 178, 53, 168, 177, 178, 57, 180, 57, 58, 157,
		64, 129, 130, 193, 132, 193, 194, 197, 136, 193, 194, 201, 196, 201, 202, 77,
		144, 193, 194, 209, 196, 209, 210, 85, 200, 209, 210, 89, 212, 89, 90, 157,
		160, 193, 194, 225, 196, 225, 226, 101, 200, 225, 226, 105, 228, 105, 106, 173,
		208, 225, 226, 113, 228, 113, 114, 181, 232, 113, 114, 185, 116, 185, 186, 93,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 30,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 46,
		32, 48, 48, 50, 48, 52, 52, 54, 48, 56, 56, 58, 56, 60, 60, 30,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 78,
		64, 80, 80, 82, 80, 84, 84, 86, 80, 88, 88, 90, 88, 92, 92, 30,
		64, 96, 96, 98, 96, 100, 100, 102, 96, 104, 104, 106, 104, 108, 108, 46,
		96, 112, 112, 114, 112, 116, 116, 54, 112, 120, 120, 58, 120, 60, 60, 94,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 142,
		128, 144, 144, 146, 144, 148, 148, 150, 144, 152, 152, 154, 152, 156, 156, 30,
		128, 160, 160, 162, 160, 164, 164, 166, 160, 168, 168, 170, 168, 172, 172, 46,
		160, 176, 176, 178, 176, 180, 180, 54, 176, 184, 184, 58, 184, 60, 60, 158,
		128, 192, 192, 194, 192, 196, 196, 198, 192, 200, 200, 202, 200, 204, 204, 78,
		192, 208, 208, 210, 208, 212, 212, 86, 208, 216, 216, 90, 216, 92, 92, 158,
		192, 224, 224, 226, 224, 228, 228, 102, 224, 232, 232, 106, 232, 108, 108, 174,
		224, 240, 240, 114, 240, 116, 116, 182, 240, 120, 120, 186, 120, 188, 188, 94,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 31,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
		80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 31,
		96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 47,
		112, 113, 114, 115, 116, 117, 118, 55, 120, 121, 122, 59, 124, 61, 62, 95,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143,
		144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 31,
		160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 47,
		176, 177, 178, 179, 180, 181, 182, 55, 184, 185, 186, 59, 188, 61, 62, 159,
		192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 79,
		208, 209, 210, 211, 212, 213, 214, 87, 216, 217, 218, 91, 220, 93, 94, 159,
		224, 225, 226, 227, 228, 229, 230, 103, 232, 233, 234, 107, 236, 109, 110, 175,
		240, 241, 242, 115, 244, 117, 118, 183, 248, 121, 122, 187, 124, 189, 190, 95,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 96,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 160,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 192,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 192,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 96,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 33,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		16, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 97,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 161,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 193,
		32, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 193,
		16, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 97,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 34,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 2,
		64, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 66,
		64, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 66,
		32, 16, 16, 2, 16, 4, 4, 66, 16, 8, 8, 66, 8, 68, 68, 98,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 2,
		128, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 130,
		128, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 130,
		32, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 162,
		128, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 130,
		64, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 194,
		64, 32, 32, 2, 32, 4, 4, 130, 32, 8, 8, 130, 8, 132, 132, 194,
		32, 16, 16, 130, 16, 132, 132, 194, 16, 136, 136, 194, 136, 196, 196, 98,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 35,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 3,
		80, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 67,
		96, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 67,
		48, 17, 18, 3, 20, 5, 6, 67, 24, 9, 10, 67, 12, 69, 70, 99,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 3,
		144, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 131,
		160, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 131,
		48, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 163,
		192, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 131,
		80, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 195,
		96, 33, 34, 3, 36, 5, 6, 131, 40, 9, 10, 131, 12, 133, 134, 195,
		48, 17, 18, 131, 20, 133, 134, 195, 24, 137, 138, 195, 140, 197, 198, 99,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 4,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 4,
		0, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 36,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 4,
		0, 64, 64, 16, 64, 16, 16, 4, 64, 16, 16, 8, 16, 8, 8, 68,
		0, 64, 64, 32, 64, 32, 32, 4, 64, 32, 32, 8, 32, 8, 8, 68,
		64, 32, 32, 16, 32, 16, 16, 68, 32, 16, 16, 72, 16, 72, 72, 100,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 4,
		0, 128, 128, 16, 128, 16, 16, 4, 128, 16, 16, 8, 16, 8, 8, 132,
		0, 128, 128, 32, 128, 32, 32, 4, 128, 32, 32, 8, 32, 8, 8, 132,
		128, 32, 32, 16, 32, 16, 16, 132, 32, 16, 16, 136, 16, 136, 136, 164,
		0, 128, 128, 64, 128, 64, 64, 4, 128, 64, 64, 8, 64, 8, 8, 132,
		128, 64, 64, 16, 64, 16, 16, 132, 64, 16, 16, 136, 16, 136, 136, 196,
		128, 64, 64, 32, 64, 32, 32, 132, 64, 32, 32, 136, 32, 136, 136, 196,
		64, 32, 32, 144, 32, 144, 144, 196, 32, 144, 144, 200, 144, 200, 200, 100,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 5,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 5,
		16, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 37,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 5,
		16, 65, 66, 17, 68, 17, 18, 5, 72, 17, 18, 9, 20, 9, 10, 69,
		32, 65, 66, 33, 68, 33, 34, 5, 72, 33, 34, 9, 36, 9, 10, 69,
		80, 33, 34, 17, 36, 17, 18, 69, 40, 17, 18, 73, 20, 73, 74, 101,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 5,
		16, 129, 130, 17, 132, 17, 18, 5, 136, 17, 18, 9, 20, 9, 10, 133,
		32, 129, 130, 33, 132, 33, 34, 5, 136, 33, 34, 9, 36, 9, 10, 133,
		144, 33, 34, 17, 36, 17, 18, 133, 40, 17, 18, 137, 20, 137, 138, 165,
		64, 129, 130, 65, 132, 65, 66, 5, 136, 65, 66, 9, 68, 9, 10, 133,
		144, 65, 66, 17, 68, 17, 18, 133, 72, 17, 18, 137, 20, 137, 138, 197,
		160, 65, 66, 33, 68, 33, 34, 133, 72, 33, 34, 137, 36, 137, 138, 197,
		80, 33, 34, 145, 36, 145, 146, 197, 40, 145, 146, 201, 148, 201, 202, 101,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 6,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 6,
		32, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 38,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 6,
		64, 80, 80, 18, 80, 20, 20, 6, 80, 24, 24, 10, 24, 12, 12, 70,
		64, 96, 96, 34, 96, 36, 36, 6, 96, 40, 40, 10, 40, 12, 12, 70,
		96, 48, 48, 18, 48, 20, 20, 70, 48, 24, 24, 74, 24, 76, 76, 102,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 6,
		128, 144, 144, 18, 144, 20, 20, 6, 144, 24, 24, 10, 24, 12, 12, 134,
		128, 160, 160, 34, 160, 36, 36, 6, 160, 40, 40, 10, 40, 12, 12, 134,
		160, 48, 48, 18, 48, 20, 20, 134, 48, 24, 24, 138, 24, 140, 140, 166,
		128, 192, 192, 66, 192, 68, 68, 6, 192, 72, 72, 10, 72, 12, 12, 134,
		192, 80, 80, 18, 80, 20, 20, 134, 80, 24, 24, 138, 24, 140, 140, 198,
		192, 96, 96, 34, 96, 36, 36, 134, 96, 40, 40, 138, 40, 140, 140, 198,
		96, 48, 48, 146, 48, 148, 148, 198, 48, 152, 152, 202, 152, 204, 204, 102,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 7,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 7,
		48, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 39,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 7,
		80, 81, 82, 19, 84, 21, 22, 7, 88, 25, 26, 11, 28, 13, 14, 71,
		96, 97, 98, 35, 100, 37, 38, 7, 104, 41, 42, 11, 44, 13, 14, 71,
		112, 49, 50, 19, 52, 21, 22, 71, 56, 25, 26, 75, 28, 77, 78, 103,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 7,
		144, 145, 146, 19, 148, 21, 22, 7, 152, 25, 26, 11, 28, 13, 14, 135,
		160, 161, 162, 35, 164, 37, 38, 7, 168, 41, 42, 11, 44, 13, 14, 135,
		176, 49, 50, 19, 52, 21, 22, 135, 56, 25, 26, 139, 28, 141, 142, 167,
		192, 193, 194, 67, 196, 69, 70, 7, 200, 73, 74, 11, 76, 13, 14, 135,
		208, 81, 82, 19, 84, 21, 22, 135, 88, 25, 26, 139, 28, 141, 142, 199,
		224, 97, 98, 35, 100, 37, 38, 135, 104, 41, 42, 139, 44, 141, 142, 199,
		112, 49, 50, 147, 52, 149, 150, 199, 56, 153, 154, 203, 156, 205, 206, 103,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 8,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 8,
		0, 0, 0, 32, 0, 32, 32, 16, 0, 32, 32, 16, 32, 16, 16, 40,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 8,
		0, 0, 0, 64, 0, 64, 64, 16, 0, 64, 64, 16, 64, 16, 16, 72,
		0, 0, 0, 64, 0, 64, 64, 32, 0, 64, 64, 32, 64, 32, 32, 72,
		0, 64, 64, 32, 64, 32, 32, 80, 64, 32, 32, 80, 32, 80, 80, 104,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 8,
		0, 0, 0, 128, 0, 128, 128, 16, 0, 128, 128, 16, 128, 16, 16, 136,
		0, 0, 0, 128, 0, 128, 128, 32, 0, 128, 128, 32, 128, 32, 32, 136,
		0, 128, 128, 32, 128, 32, 32, 144, 128, 32, 32, 144, 32, 144, 144, 168,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 136,
		0, 128, 128, 64, 128, 64, 64, 144, 128, 64, 64, 144, 64, 144, 144, 200,
		0, 128, 128, 64, 128, 64, 64, 160, 128, 64, 64, 160, 64, 160, 160, 200,
		128, 64, 64, 160, 64, 160, 160, 208, 64, 160, 160, 208, 160, 208, 208, 104,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 9,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 9,
		16, 1, 2, 33, 4, 33, 34, 17, 8, 33, 34, 17, 36, 17, 18, 41,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 9,
		16, 1, 2, 65, 4, 65, 66, 17, 8, 65, 66, 17, 68, 17, 18, 73,
		32, 1, 2, 65, 4, 65, 66, 33, 8, 65, 66, 33, 68, 33, 34, 73,
		16, 65, 66, 33, 68, 33, 34, 81, 72, 33, 34, 81, 36, 81, 82, 105,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 9,
		16, 1, 2, 129, 4, 129, 130, 17, 8, 129, 130, 17, 132, 17, 18, 137,
		32, 1, 2, 129, 4, 129, 130, 33, 8, 129, 130, 33, 132, 33, 34, 137,
		16, 129, 130, 33, 132, 33, 34, 145, 136, 33, 34, 145, 36, 145, 146, 169,
		64, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 137,
		16, 129, 130, 65, 132, 65, 66, 145, 136, 65, 66, 145, 68, 145, 146, 201,
		32, 129, 130, 65, 132, 65, 66, 161, 136, 65, 66, 161, 68, 161, 162, 201,
		144, 65, 66, 161, 68, 161, 162, 209, 72, 161, 162, 209, 164, 209, 210, 105,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 10,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 10,
		32, 16, 16, 34, 16, 36, 36, 18, 16, 40, 40, 18, 40, 20, 20, 42,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 10,
		64, 16, 16, 66, 16, 68, 68, 18, 16, 72, 72, 18, 72, 20, 20, 74,
		64, 32, 32, 66, 32, 68, 68, 34, 32, 72, 72, 34, 72, 36, 36, 74,
		32, 80, 80, 34, 80, 36, 36, 82, 80, 40, 40, 82, 40, 84, 84, 106,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 10,
		128, 16, 16, 130, 16, 132, 132, 18, 16, 136, 136, 18, 136, 20, 20, 138,
		128, 32, 32, 130, 32, 132, 132, 34, 32, 136, 136, 34, 136, 36, 36, 138,
		32, 144, 144, 34, 144, 36, 36, 146, 144, 40, 40, 146, 40, 148, 148, 170,
		128, 64, 64, 130, 64, 132, 132, 66, 64, 136, 136, 66, 136, 68, 68, 138,
		64, 144, 144, 66, 144, 68, 68, 146, 144, 72, 72, 146, 72, 148, 148, 202,
		64, 160, 160, 66, 160, 68, 68, 162, 160, 72, 72, 162, 72, 164, 164, 202,
		160, 80, 80, 162, 80, 164, 164, 210, 80, 168, 168, 210, 168, 212, 212, 106,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 11,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 11,
		48, 17, 18, 35, 20, 37, 38, 19, 24, 41, 42, 19, 44, 21, 22, 43,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 11,
		80, 17, 18, 67, 20, 69, 70, 19, 24, 73, 74, 19, 76, 21, 22, 75,
		96, 33, 34, 67, 36, 69, 70, 35, 40, 73, 74, 35, 76, 37, 38, 75,
		48, 81, 82, 35, 84, 37, 38, 83, 88, 41, 42, 83, 44, 85, 86, 107,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 11,
		144, 17, 18, 131, 20, 133, 134, 19, 24, 137, 138, 19, 140, 21, 22, 139,
		160, 33, 34, 131, 36, 133, 134, 35, 40, 137, 138, 35, 140, 37, 38, 139,
		48, 145, 146, 35, 148, 37, 38, 147, 152, 41, 42, 147, 44, 149, 150, 171,
		192, 65, 66, 131, 68, 133, 134, 67, 72, 137, 138, 67, 140, 69, 70, 139,
		80, 145, 146, 67, 148, 69, 70, 147, 152, 73, 74, 147, 76, 149, 150, 203,
		96, 161, 162, 67, 164, 69, 70, 163, 168, 73, 74, 163, 76, 165, 166, 203,
		176, 81, 82, 163, 84, 165, 166, 211, 88, 169, 170, 211, 172, 213, 214, 107,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 12,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 12,
		0, 32, 32, 48, 32, 48, 48, 20, 32, 48, 48, 24, 48, 24, 24, 44,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 12,
		0, 64, 64, 80, 64, 80, 80, 20, 64, 80, 80, 24, 80, 24, 24, 76,
		0, 64, 64, 96, 64, 96, 96, 36, 64, 96, 96, 40, 96, 40, 40, 76,
		64, 96, 96, 48, 96, 48, 48, 84, 96, 48, 48, 88, 48, 88, 88, 108,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 12,
		0, 128, 128, 144, 128, 144, 144, 20, 128, 144, 144, 24, 144, 24, 24, 140,
		0, 128, 128, 160, 128, 160, 160, 36, 128, 160, 160, 40, 160, 40, 40, 140,
		128, 160, 160, 48, 160, 48, 48, 148, 160, 48, 48, 152, 48, 152, 152, 172,
		0, 128, 128, 192, 128, 192, 192, 68, 128, 192, 192, 72, 192, 72, 72, 140,
		128, 192, 192, 80, 192, 80, 80, 148, 192, 80, 80, 152, 80, 152, 152, 204,
		128, 192, 192, 96, 192, 96, 96, 164, 192, 96, 96, 168, 96, 168, 168, 204,
		192, 96, 96, 176, 96, 176, 176, 212, 96, 176, 176, 216, 176, 216, 216, 108,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 13,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 13,
		16, 33, 34, 49, 36, 49, 50, 21, 40, 49, 50, 25, 52, 25, 26, 45,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 13,
		16, 65, 66, 81, 68, 81, 82, 21, 72, 81, 82, 25, 84, 25, 26, 77,
		32, 65, 66, 97, 68, 97, 98, 37, 72, 97, 98, 41, 100, 41, 42, 77,
		80, 97, 98, 49, 100, 49, 50, 85, 104, 49, 50, 89, 52, 89, 90, 109,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 13,
		16, 129, 130, 145, 132, 145, 146, 21, 136, 145, 146, 25, 148, 25, 26, 141,
		32, 129, 130, 161, 132, 161, 162, 37, 136, 161, 162, 41, 164, 41, 42, 141,
		144, 161, 162, 49, 164, 49, 50, 149, 168, 49, 50, 153, 52, 153, 154, 173,
		64, 129, 130, 193, 132, 193, 194, 69, 136, 193, 194, 73, 196, 73, 74, 141,
		144, 193, 194, 81, 196, 81, 82, 149, 200, 81, 82, 153, 84, 153, 154, 205,
		160, 193, 194, 97, 196, 97, 98, 165, 200, 97, 98, 169, 100, 169, 170, 205,
		208, 97, 98, 177, 100, 177, 178, 213, 104, 177, 178, 217, 180, 217, 218, 109,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 14,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 14,
		32, 48, 48, 50, 48, 52, 52, 22, 48, 56, 56, 26, 56, 28, 28, 46,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 14,
		64, 80, 80, 82, 80, 84, 84, 22, 80, 88, 88, 26, 88, 28, 28, 78,
		64, 96, 96, 98, 96, 100, 100, 38, 96, 104, 104, 42, 104, 44, 44, 78,
		96, 112, 112, 50, 112, 52, 52, 86, 112, 56, 56, 90, 56, 92, 92, 110,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 14,
		128, 144, 144, 146, 144, 148, 148, 22, 144, 152, 152, 26, 152, 28, 28, 142,
		128, 160, 160, 162, 160, 164, 164, 38, 160, 168, 168, 42, 168, 44, 44, 142,
		160, 176, 176, 50, 176, 52, 52, 150, 176, 56, 56, 154, 56, 156, 156, 174,
		128, 192, 192, 194, 192, 196, 196, 70, 192, 200, 200, 74, 200, 76, 76, 142,
		192, 208, 208, 82, 208, 84, 84, 150, 208, 88, 88, 154, 88, 156, 156, 206,
		192, 224, 224, 98, 224, 100, 100, 166, 224, 104, 104, 170, 104, 172, 172, 206,
		224, 112, 112, 178, 112, 180, 180, 214, 112, 184, 184, 218, 184, 220, 220, 110,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 15,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 15,
		48, 49, 50, 51, 52, 53, 54, 23, 56, 57, 58, 27, 60, 29, 30, 47,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 15,
		80, 81, 82, 83, 84, 85, 86, 23, 88, 89, 90, 27, 92, 29, 30, 79,
		96, 97, 98, 99, 100, 101, 102, 39, 104, 105, 106, 43, 108, 45, 46, 79,
		112, 113, 114, 51, 116, 53, 54, 87, 120, 57, 58, 91, 60, 93, 94, 111,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 15,
		144, 145, 146, 147, 148, 149, 150, 23, 152, 153, 154, 27, 156, 29, 30, 143,
		160, 161, 162, 163, 164, 165, 166, 39, 168, 169, 170, 43, 172, 45, 46, 143,
		176, 177, 178, 51, 180, 53, 54, 151, 184, 57, 58, 155, 60, 157, 158, 175,
		192, 193, 194, 195, 196, 197, 198, 71, 200, 201, 202, 75, 204, 77, 78, 143,
		208, 209, 210, 83, 212, 85, 86, 151, 216, 89, 90, 155, 92, 157, 158, 207,
		224, 225, 226, 99, 228, 101, 102, 167, 232, 105, 106, 171, 108, 173, 174, 207,
		240, 113, 114, 179, 116, 181, 182, 215, 120, 185, 186, 219, 188, 221, 222, 111,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 48,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 80,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 96,
		0, 0, 0, 64, 0, 64, 64, 96, 0, 64, 64, 96, 64, 96, 96, 112,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 144,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 160,
		0, 0, 0, 128, 0, 128, 128, 160, 0, 128, 128, 160, 128, 160, 160, 176,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 192,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 208,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 224,
		0, 128, 128, 192, 128, 192, 192, 224, 128, 192, 192, 224, 192, 224, 224, 112,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 17,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 33,
		16, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 49,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		16, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 81,
		32, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 97,
		16, 1, 2, 65, 4, 65, 66, 97, 8, 65, 66, 97, 68, 97, 98, 113,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 145,
		32, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 161,
		16, 1, 2, 129, 4, 129, 130, 161, 8, 129, 130, 161, 132, 161, 162, 177,
		64, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 193,
		16, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 209,
		32, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 225,
		16, 129, 130, 193, 132, 193, 194, 225, 136, 193, 194, 225, 196, 225, 226, 113,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 18,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 34,
		32, 16, 16, 2, 16, 4, 4, 34, 16, 8, 8, 34, 8, 36, 36, 50,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 66,
		64, 16, 16, 2, 16, 4, 4, 66, 16, 8, 8, 66, 8, 68, 68, 82,
		64, 32, 32, 2, 32, 4, 4, 66, 32, 8, 8, 66, 8, 68, 68, 98,
		32, 16, 16, 66, 16, 68, 68, 98, 16, 72, 72, 98, 72, 100, 100, 114,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 130,
		128, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 146,
		128, 32, 32, 2, 32, 4, 4, 130, 32, 8, 8, 130, 8, 132, 132, 162,
		32, 16, 16, 130, 16, 132, 132, 162, 16, 136, 136, 162, 136, 164, 164, 178,
		128, 64, 64, 2, 64, 4, 4, 130, 64, 8, 8, 130, 8, 132, 132, 194,
		64, 16, 16, 130, 16, 132, 132, 194, 16, 136, 136, 194, 136, 196, 196, 210,
		64, 32, 32, 130, 32, 132, 132, 194, 32, 136, 136, 194, 136, 196, 196, 226,
		32, 144, 144, 194, 144, 196, 196, 226, 144, 200, 200, 226, 200, 228, 228, 114,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 19,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 35,
		48, 17, 18, 3, 20, 5, 6, 35, 24, 9, 10, 35, 12, 37, 38, 51,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 67,
		80, 17, 18, 3, 20, 5, 6, 67, 24, 9, 10, 67, 12, 69, 70, 83,
		96, 33, 34, 3, 36, 5, 6, 67, 40, 9, 10, 67, 12, 69, 70, 99,
		48, 17, 18, 67, 20, 69, 70, 99, 24, 73, 74, 99, 76, 101, 102, 115,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 131,
		144, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 147,
		160, 33, 34, 3, 36, 5, 6, 131, 40, 9, 10, 131, 12, 133, 134, 163,
		48, 17, 18, 131, 20, 133, 134, 163, 24, 137, 138, 163, 140, 165, 166, 179,
		192, 65, 66, 3, 68, 5, 6, 131, 72, 9, 10, 131, 12, 133, 134, 195,
		80, 17, 18, 131, 20, 133, 134, 195, 24, 137, 138, 195, 140, 197, 198, 211,
		96, 33, 34, 131, 36, 133, 134, 195, 40, 137, 138, 195, 140, 197, 198, 227,
		48, 145, 146, 195, 148, 197, 198, 227, 152, 201, 202, 227, 204, 229, 230, 115,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 20,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 36,
		0, 32, 32, 16, 32, 16, 16, 36, 32, 16, 16, 40, 16, 40, 40, 52,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 68,
		0, 64, 64, 16, 64, 16, 16, 68, 64, 16, 16, 72, 16, 72, 72, 84,
		0, 64, 64, 32, 64, 32, 32, 68, 64, 32, 32, 72, 32, 72, 72, 100,
		64, 32, 32, 80, 32, 80, 80, 100, 32, 80, 80, 104, 80, 104, 104, 116,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 132,
		0, 128, 128, 16, 128, 16, 16, 132, 128, 16, 16, 136, 16, 136, 136, 148,
		0, 128, 128, 32, 128, 32, 32, 132, 128, 32, 32, 136, 32, 136, 136, 164,
		128, 32, 32, 144, 32, 144, 144, 164, 32, 144, 144, 168, 144, 168, 168, 180,
		0, 128, 128, 64, 128, 64, 64, 132, 128, 64, 64, 136, 64, 136, 136, 196,
		128, 64, 64, 144, 64, 144, 144, 196, 64, 144, 144, 200, 144, 200, 200, 212,
		128, 64, 64, 160, 64, 160, 160, 196, 64, 160, 160, 200, 160, 200, 200, 228,
		64, 160, 160, 208, 160, 208, 208, 228, 160, 208, 208, 232, 208, 232, 232, 116,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 21,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 37,
		16, 33, 34, 17, 36, 17, 18, 37, 40, 17, 18, 41, 20, 41, 42, 53,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 69,
		16, 65, 66, 17, 68, 17, 18, 69, 72, 17, 18, 73, 20, 73, 74, 85,
		32, 65, 66, 33, 68, 33, 34, 69, 72, 33, 34, 73, 36, 73, 74, 101,
		80, 33, 34, 81, 36, 81, 82, 101, 40, 81, 82, 105, 84, 105, 106, 117,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 133,
		16, 129, 130, 17, 132, 17, 18, 133, 136, 17, 18, 137, 20, 137, 138, 149,
		32, 129, 130, 33, 132, 33, 34, 133, 136, 33, 34, 137, 36, 137, 138, 165,
		144, 33, 34, 145, 36, 145, 146, 165, 40, 145, 146, 169, 148, 169, 170, 181,
		64, 129, 130, 65, 132, 65, 66, 133, 136, 65, 66, 137, 68, 137, 138, 197,
		144, 65, 66, 145, 68, 145, 146, 197, 72, 145, 146, 201, 148, 201, 202, 213,
		160, 65, 66, 161, 68, 161, 162, 197, 72, 161, 162, 201, 164, 201, 202, 229,
		80, 161, 162, 209, 164, 209, 210, 229, 168, 209, 210, 233, 212, 233, 234, 117,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 22,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 38,
		32, 48, 48, 18, 48, 20, 20, 38, 48, 24, 24, 42, 24, 44, 44, 54,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 70,
		64, 80, 80, 18, 80, 20, 20, 70, 80, 24, 24, 74, 24, 76, 76, 86,
		64, 96, 96, 34, 96, 36, 36, 70, 96, 40, 40, 74, 40, 76, 76, 102,
		96, 48, 48, 82, 48, 84, 84, 102, 48, 88, 88, 106, 88, 108, 108, 118,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 134,
		128, 144, 144, 18, 144, 20, 20, 134, 144, 24, 24, 138, 24, 140, 140, 150,
		128, 160, 160, 34, 160, 36, 36, 134, 160, 40, 40, 138, 40, 140, 140, 166,
		160, 48, 48, 146, 48, 148, 148, 166, 48, 152, 152, 170, 152, 172, 172, 182,
		128, 192, 192, 66, 192, 68, 68, 134, 192, 72, 72, 138, 72, 140, 140, 198,
		192, 80, 80, 146, 80, 148, 148, 198, 80, 152, 152, 202, 152, 204, 204, 214,
		192, 96, 96, 162, 96, 164, 164, 198, 96, 168, 168, 202, 168, 204, 204, 230,
		96, 176, 176, 210, 176, 212, 212, 230, 176, 216, 216, 234, 216, 236, 236, 118,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 23,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 39,
		48, 49, 50, 19, 52, 21, 22, 39, 56, 25, 26, 43, 28, 45, 46, 55,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 71,
		80, 81, 82, 19, 84, 21, 22, 71, 88, 25, 26, 75, 28, 77, 78, 87,
		96, 97, 98, 35, 100, 37, 38, 71, 104, 41, 42, 75, 44, 77, 78, 103,
		112, 49, 50, 83, 52, 85, 86, 103, 56, 89, 90, 107, 92, 109, 110, 119,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 135,
		144, 145, 146, 19, 148, 21, 22, 135, 152, 25, 26, 139, 28, 141, 142, 151,
		160, 161, 162, 35, 164, 37, 38, 135, 168, 41, 42, 139, 44, 141, 142, 167,
		176, 49, 50, 147, 52, 149, 150, 167, 56, 153, 154, 171, 156, 173, 174, 183,
		192, 193, 194, 67, 196, 69, 70, 135, 200, 73, 74, 139, 76, 141, 142, 199,
		208, 81, 82, 147, 84, 149, 150, 199, 88, 153, 154, 203, 156, 205, 206, 215,
		224, 97, 98, 163, 100, 165, 166, 199, 104, 169, 170, 203, 172, 205, 206, 231,
		112, 177, 178, 211, 180, 213, 214, 231, 184, 217, 218, 235, 220, 237, 238, 119,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 24,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 40,
		0, 0, 0, 32, 0, 32, 32, 48, 0, 32, 32, 48, 32, 48, 48, 56,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 72,
		0, 0, 0, 64, 0, 64, 64, 80, 0, 64, 64, 80, 64, 80, 80, 88,
		0, 0, 0, 64, 0, 64, 64, 96, 0, 64, 64, 96, 64, 96, 96, 104,
		0, 64, 64, 96, 64, 96, 96, 112, 64, 96, 96, 112, 96, 112, 112, 120,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 136,
		0, 0, 0, 128, 0, 128, 128, 144, 0, 128, 128, 144, 128, 144, 144, 152,
		0, 0, 0, 128, 0, 128, 128, 160, 0, 128, 128, 160, 128, 160, 160, 168,
		0, 128, 128, 160, 128, 160, 160, 176, 128, 160, 160, 176, 160, 176, 176, 184,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 200,
		0, 128, 128, 192, 128, 192, 192, 208, 128, 192, 192, 208, 192, 208, 208, 216,
		0, 128, 128, 192, 128, 192, 192, 224, 128, 192, 192, 224, 192, 224, 224, 232,
		128, 192, 192, 224, 192, 224, 224, 240, 192, 224, 224, 240, 224, 240, 240, 120,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 25,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 41,
		16, 1, 2, 33, 4, 33, 34, 49, 8, 33, 34, 49, 36, 49, 50, 57,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 73,
		16, 1, 2, 65, 4, 65, 66, 81, 8, 65, 66, 81, 68, 81, 82, 89,
		32, 1, 2, 65, 4, 65, 66, 97, 8, 65, 66, 97, 68, 97, 98, 105,
		16, 65, 66, 97, 68, 97, 98, 113, 72, 97, 98, 113, 100, 113, 114, 121,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 137,
		16, 1, 2, 129, 4, 129, 130, 145, 8, 129, 130, 145, 132, 145, 146, 153,
		32, 1, 2, 129, 4, 129, 130, 161, 8, 129, 130, 161, 132, 161, 162, 169,
		16, 129, 130, 161, 132, 161, 162, 177, 136, 161, 162, 177, 164, 177, 178, 185,
		64, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 201,
		16, 129, 130, 193, 132, 193, 194, 209, 136, 193, 194, 209, 196, 209, 210, 217,
		32, 129, 130, 193, 132, 193, 194, 225, 136, 193, 194, 225, 196, 225, 226, 233,
		144, 193, 194, 225, 196, 225, 226, 241, 200, 225, 226, 241, 228, 241, 242, 121,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 26,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 42,
		32, 16, 16, 34, 16, 36, 36, 50, 16, 40, 40, 50, 40, 52, 52, 58,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 74,
		64, 16, 16, 66, 16, 68, 68, 82, 16, 72, 72, 82, 72, 84, 84, 90,
		64, 32, 32, 66, 32, 68, 68, 98, 32, 72, 72, 98, 72, 100, 100, 106,
		32, 80, 80, 98, 80, 100, 100, 114, 80, 104, 104, 114, 104, 116, 116, 122,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 138,
		128, 16, 16, 130, 16, 132, 132, 146, 16, 136, 136, 146, 136, 148, 148, 154,
		128, 32, 32, 130, 32, 132, 132, 162, 32, 136, 136, 162, 136, 164, 164, 170,
		32, 144, 144, 162, 144, 164, 164, 178, 144, 168, 168, 178, 168, 180, 180, 186,
		128, 64, 64, 130, 64, 132, 132, 194, 64, 136, 136, 194, 136, 196, 196, 202,
		64, 144, 144, 194, 144, 196, 196, 210, 144, 200, 200, 210, 200, 212, 212, 218,
		64, 160, 160, 194, 160, 196, 196, 226, 160, 200, 200, 226, 200, 228, 228, 234,
		160, 208, 208, 226, 208, 228, 228, 242, 208, 232, 232, 242, 232, 244, 244, 122,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 27,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 43,
		48, 17, 18, 35, 20, 37, 38, 51, 24, 41, 42, 51, 44, 53, 54, 59,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 75,
		80, 17, 18, 67, 20, 69, 70, 83, 24, 73, 74, 83, 76, 85, 86, 91,
		96, 33, 34, 67, 36, 69, 70, 99, 40, 73, 74, 99, 76, 101, 102, 107,
		48, 81, 82, 99, 84, 101, 102, 115, 88, 105, 106, 115, 108, 117, 118, 123,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 139,
		144, 17, 18, 131, 20, 133, 134, 147, 24, 137, 138, 147, 140, 149, 150, 155,
		160, 33, 34, 131, 36, 133, 134, 163, 40, 137, 138, 163, 140, 165, 166, 171,
		48, 145, 146, 163, 148, 165, 166, 179, 152, 169, 170, 179, 172, 181, 182, 187,
		192, 65, 66, 131, 68, 133, 134, 195, 72, 137, 138, 195, 140, 197, 198, 203,
		80, 145, 146, 195, 148, 197, 198, 211, 152, 201, 202, 211, 204, 213, 214, 219,
		96, 161, 162, 195, 164, 197, 198, 227, 168, 201, 202, 227, 204, 229, 230, 235,
		176, 209, 210, 227, 212, 229, 230, 243, 216, 233, 234, 243, 236, 245, 246, 123,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 28,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 44,
		0, 32, 32, 48, 32, 48, 48, 52, 32, 48, 48, 56, 48, 56, 56, 60,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 76,
		0, 64, 64, 80, 64, 80, 80, 84, 64, 80, 80, 88, 80, 88, 88, 92,
		0, 64, 64, 96, 64, 96, 96, 100, 64, 96, 96, 104, 96, 104, 104, 108,
		64, 96, 96, 112, 96, 112, 112, 116, 96, 112, 112, 120, 112, 120, 120, 124,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 140,
		0, 128, 128, 144, 128, 144, 144, 148, 128, 144, 144, 152, 144, 152, 152, 156,
		0, 128, 128, 160, 128, 160, 160, 164, 128, 160, 160, 168, 160, 168, 168, 172,
		128, 160, 160, 176, 160, 176, 176, 180, 160, 176, 176, 184, 176, 184, 184, 188,
		0, 128, 128, 192, 128, 192, 192, 196, 128, 192, 192, 200, 192, 200, 200, 204,
		128, 192, 192, 208, 192, 208, 208, 212, 192, 208, 208, 216, 208, 216, 216, 220,
		128, 192, 192, 224, 192, 224, 224, 228, 192, 224, 224, 232, 224, 232, 232, 236,
		192, 224, 224, 240, 224, 240, 240, 244, 224, 240, 240, 248, 240, 248, 248, 124,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 29,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 45,
		16, 33, 34, 49, 36, 49, 50, 53, 40, 49, 50, 57, 52, 57, 58, 61,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 77,
		16, 65, 66, 81, 68, 81, 82, 85, 72, 81, 82, 89, 84, 89, 90, 93,
		32, 65, 66, 97, 68, 97, 98, 101, 72, 97, 98, 105, 100, 105, 106, 109,
		80, 97, 98, 113, 100, 113, 114, 117, 104, 113, 114, 121, 116, 121, 122, 125,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 141,
		16, 129, 130, 145, 132, 145, 146, 149, 136, 145, 146, 153, 148, 153, 154, 157,
		32, 129, 130, 161, 132, 161, 162, 165, 136, 161, 162, 169, 164, 169, 170, 173,
		144, 161, 162, 177, 164, 177, 178, 181, 168, 177, 178, 185, 180, 185, 186, 189,
		64, 129, 130, 193, 132, 193, 194, 197, 136, 193, 194, 201, 196, 201, 202, 205,
		144, 193, 194, 209, 196, 209, 210, 213, 200, 209, 210, 217, 212, 217, 218, 221,
		160, 193, 194, 225, 196, 225, 226, 229, 200, 225, 226, 233, 228, 233, 234, 237,
		208, 225, 226, 241, 228, 241, 242, 245, 232, 241, 242, 249, 244, 249, 250, 125,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 30,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 46,
		32, 48, 48, 50, 48, 52, 52, 54, 48, 56, 56, 58, 56, 60, 60, 62,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 78,
		64, 80, 80, 82, 80, 84, 84, 86, 80, 88, 88, 90, 88, 92, 92, 94,
		64, 96, 96, 98, 96, 100, 100, 102, 96, 104, 104, 106, 104, 108, 108, 110,
		96, 112, 112, 114, 112, 116, 116, 118, 112, 120, 120, 122, 120, 124, 124, 126,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 142,
		128, 144, 144, 146, 144, 148, 148, 150, 144, 152, 152, 154, 152, 156, 156, 158,
		128, 160, 160, 162, 160, 164, 164, 166, 160, 168, 168, 170, 168, 172, 172, 174,
		160, 176, 176, 178, 176, 180, 180, 182, 176, 184, 184, 186, 184, 188, 188, 190,
		128, 192, 192, 194, 192, 196, 196, 198, 192, 200, 200, 202, 200, 204, 204, 206,
		192, 208, 208, 210, 208, 212, 212, 214, 208, 216, 216, 218, 216, 220, 220, 222,
		192, 224, 224, 226, 224, 228, 228, 230, 224, 232, 232, 234, 232, 236, 236, 238,
		224, 240, 240, 242, 240, 244, 244, 246, 240, 248, 248, 250, 248, 252, 252, 126,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
		80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95,
		96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111,
		112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143,
		144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159,
		160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175,
		176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191,
		192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207,
		208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223,
		224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239,
		240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 127,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 2,
		64, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		64, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 2,
		128, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		128, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		128, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 2,
		64, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		64, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 130,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 3,
		80, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		96, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 3,
		144, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		160, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		192, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 3,
		80, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		96, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 131,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 4,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 4,
		0, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 4,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 4,
		0, 64, 64, 16, 64, 16, 16, 4, 64, 16, 16, 8, 16, 8, 8, 4,
		0, 64, 64, 32, 64, 32, 32, 4, 64, 32, 32, 8, 32, 8, 8, 4,
		64, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 4,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 4,
		0, 128, 128, 16, 128, 16, 16, 4, 128, 16, 16, 8, 16, 8, 8, 4,
		0, 128, 128, 32, 128, 32, 32, 4, 128, 32, 32, 8, 32, 8, 8, 4,
		128, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 4,
		0, 128, 128, 64, 128, 64, 64, 4, 128, 64, 64, 8, 64, 8, 8, 4,
		128, 64, 64, 16, 64, 16, 16, 4, 64, 16, 16, 8, 16, 8, 8, 4,
		128, 64, 64, 32, 64, 32, 32, 4, 64, 32, 32, 8, 32, 8, 8, 4,
		64, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 132,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 5,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 5,
		16, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 5,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 5,
		16, 65, 66, 17, 68, 17, 18, 5, 72, 17, 18, 9, 20, 9, 10, 5,
		32, 65, 66, 33, 68, 33, 34, 5, 72, 33, 34, 9, 36, 9, 10, 5,
		80, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 5,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 5,
		16, 129, 130, 17, 132, 17, 18, 5, 136, 17, 18, 9, 20, 9, 10, 5,
		32, 129, 130, 33, 132, 33, 34, 5, 136, 33, 34, 9, 36, 9, 10, 5,
		144, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 5,
		64, 129, 130, 65, 132, 65, 66, 5, 136, 65, 66, 9, 68, 9, 10, 5,
		144, 65, 66, 17, 68, 17, 18, 5, 72, 17, 18, 9, 20, 9, 10, 5,
		160, 65, 66, 33, 68, 33, 34, 5, 72, 33, 34, 9, 36, 9, 10, 5,
		80, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 133,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 6,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 6,
		32, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 6,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 6,
		64, 80, 80, 18, 80, 20, 20, 6, 80, 24, 24, 10, 24, 12, 12, 6,
		64, 96, 96, 34, 96, 36, 36, 6, 96, 40, 40, 10, 40, 12, 12, 6,
		96, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 6,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 6,
		128, 144, 144, 18, 144, 20, 20, 6, 144, 24, 24, 10, 24, 12, 12, 6,
		128, 160, 160, 34, 160, 36, 36, 6, 160, 40, 40, 10, 40, 12, 12, 6,
		160, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 6,
		128, 192, 192, 66, 192, 68, 68, 6, 192, 72, 72, 10, 72, 12, 12, 6,
		192, 80, 80, 18, 80, 20, 20, 6, 80, 24, 24, 10, 24, 12, 12, 6,
		192, 96, 96, 34, 96, 36, 36, 6, 96, 40, 40, 10, 40, 12, 12, 6,
		96, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 134,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 7,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 7,
		48, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 7,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 7,
		80, 81, 82, 19, 84, 21, 22, 7, 88, 25, 26, 11, 28, 13, 14, 7,
		96, 97, 98, 35, 100, 37, 38, 7, 104, 41, 42, 11, 44, 13, 14, 7,
		112, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 7,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 7,
		144, 145, 146, 19, 148, 21, 22, 7, 152, 25, 26, 11, 28, 13, 14, 7,
		160, 161, 162, 35, 164, 37, 38, 7, 168, 41, 42, 11, 44, 13, 14, 7,
		176, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 7,
		192, 193, 194, 67, 196, 69, 70, 7, 200, 73, 74, 11, 76, 13, 14, 7,
		208, 81, 82, 19, 84, 21, 22, 7, 88, 25, 26, 11, 28, 13, 14, 7,
		224, 97, 98, 35, 100, 37, 38, 7, 104, 41, 42, 11, 44, 13, 14, 7,
		112, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 135,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 8,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 8,
		0, 0, 0, 32, 0, 32, 32, 16, 0, 32, 32, 16, 32, 16, 16, 8,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 8,
		0, 0, 0, 64, 0, 64, 64, 16, 0, 64, 64, 16, 64, 16, 16, 8,
		0, 0, 0, 64, 0, 64, 64, 32, 0, 64, 64, 32, 64, 32, 32, 8,
		0, 64, 64, 32, 64, 32, 32, 16, 64, 32, 32, 16, 32, 16, 16, 8,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 8,
		0, 0, 0, 128, 0, 128, 128, 16, 0, 128, 128, 16, 128, 16, 16, 8,
		0, 0, 0, 128, 0, 128, 128, 32, 0, 128, 128, 32, 128, 32, 32, 8,
		0, 128, 128, 32, 128, 32, 32, 16, 128, 32, 32, 16, 32, 16, 16, 8,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 8,
		0, 128, 128, 64, 128, 64, 64, 16, 128, 64, 64, 16, 64, 16, 16, 8,
		0, 128, 128, 64, 128, 64, 64, 32, 128, 64, 64, 32, 64, 32, 32, 8,
		128, 64, 64, 32, 64, 32, 32, 16, 64, 32, 32, 16, 32, 16, 16, 136,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 9,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 9,
		16, 1, 2, 33, 4, 33, 34, 17, 8, 33, 34, 17, 36, 17, 18, 9,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 9,
		16, 1, 2, 65, 4, 65, 66, 17, 8, 65, 66, 17, 68, 17, 18, 9,
		32, 1, 2, 65, 4, 65, 66, 33, 8, 65, 66, 33, 68, 33, 34, 9,
		16, 65, 66, 33, 68, 33, 34, 17, 72, 33, 34, 17, 36, 17, 18, 9,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 9,
		16, 1, 2, 129, 4, 129, 130, 17, 8, 129, 130, 17, 132, 17, 18, 9,
		32, 1, 2, 129, 4, 129, 130, 33, 8, 129, 130, 33, 132, 33, 34, 9,
		16, 129, 130, 33, 132, 33, 34, 17, 136, 33, 34, 17, 36, 17, 18, 9,
		64, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 9,
		16, 129, 130, 65, 132, 65, 66, 17, 136, 65, 66, 17, 68, 17, 18, 9,
		32, 129, 130, 65, 132, 65, 66, 33, 136, 65, 66, 33, 68, 33, 34, 9,
		144, 65, 66, 33, 68, 33, 34, 17, 72, 33, 34, 17, 36, 17, 18, 137,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 10,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 10,
		32, 16, 16, 34, 16, 36, 36, 18, 16, 40, 40, 18, 40, 20, 20, 10,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 10,
		64, 16, 16, 66, 16, 68, 68, 18, 16, 72, 72, 18, 72, 20, 20, 10,
		64, 32, 32, 66, 32, 68, 68, 34, 32, 72, 72, 34, 72, 36, 36, 10,
		32, 80, 80, 34, 80, 36, 36, 18, 80, 40, 40, 18, 40, 20, 20, 10,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 10,
		128, 16, 16, 130, 16, 132, 132, 18, 16, 136, 136, 18, 136, 20, 20, 10,
		128, 32, 32, 130, 32, 132, 132, 34, 32, 136, 136, 34, 136, 36, 36, 10,
		32, 144, 144, 34, 144, 36, 36, 18, 144, 40, 40, 18, 40, 20, 20, 10,
		128, 64, 64, 130, 64, 132, 132, 66, 64, 136, 136, 66, 136, 68, 68, 10,
		64, 144, 144, 66, 144, 68, 68, 18, 144, 72, 72, 18, 72, 20, 20, 10,
		64, 160, 160, 66, 160, 68, 68, 34, 160, 72, 72, 34, 72, 36, 36, 10,
		160, 80, 80, 34, 80, 36, 36, 18, 80, 40, 40, 18, 40, 20, 20, 138,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 11,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 11,
		48, 17, 18, 35, 20, 37, 38, 19, 24, 41, 42, 19, 44, 21, 22, 11,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 11,
		80, 17, 18, 67, 20, 69, 70, 19, 24, 73, 74, 19, 76, 21, 22, 11,
		96, 33, 34, 67, 36, 69, 70, 35, 40, 73, 74, 35, 76, 37, 38, 11,
		48, 81, 82, 35, 84, 37, 38, 19, 88, 41, 42, 19, 44, 21, 22, 11,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 11,
		144, 17, 18, 131, 20, 133, 134, 19, 24, 137, 138, 19, 140, 21, 22, 11,
		160, 33, 34, 131, 36, 133, 134, 35, 40, 137, 138, 35, 140, 37, 38, 11,
		48, 145, 146, 35, 148, 37, 38, 19, 152, 41, 42, 19, 44, 21, 22, 11,
		192, 65, 66, 131, 68, 133, 134, 67, 72, 137, 138, 67, 140, 69, 70, 11,
		80, 145, 146, 67, 148, 69, 70, 19, 152, 73, 74, 19, 76, 21, 22, 11,
		96, 161, 162, 67, 164, 69, 70, 35, 168, 73, 74, 35, 76, 37, 38, 11,
		176, 81, 82, 35, 84, 37, 38, 19, 88, 41, 42, 19, 44, 21, 22, 139,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 12,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 12,
		0, 32, 32, 48, 32, 48, 48, 20, 32, 48, 48, 24, 48, 24, 24, 12,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 12,
		0, 64, 64, 80, 64, 80, 80, 20, 64, 80, 80, 24, 80, 24, 24, 12,
		0, 64, 64, 96, 64, 96, 96, 36, 64, 96, 96, 40, 96, 40, 40, 12,
		64, 96, 96, 48, 96, 48, 48, 20, 96, 48, 48, 24, 48, 24, 24, 12,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 12,
		0, 128, 128, 144, 128, 144, 144, 20, 128, 144, 144, 24, 144, 24, 24, 12,
		0, 128, 128, 160, 128, 160, 160, 36, 128, 160, 160, 40, 160, 40, 40, 12,
		128, 160, 160, 48, 160, 48, 48, 20, 160, 48, 48, 24, 48, 24, 24, 12,
		0, 128, 128, 192, 128, 192, 192, 68, 128, 192, 192, 72, 192, 72, 72, 12,
		128, 192, 192, 80, 192, 80, 80, 20, 192, 80, 80, 24, 80, 24, 24, 12,
		128, 192, 192, 96, 192, 96, 96, 36, 192, 96, 96, 40, 96, 40, 40, 12,
		192, 96, 96, 48, 96, 48, 48, 20, 96, 48, 48, 24, 48, 24, 24, 140,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 13,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 13,
		16, 33, 34, 49, 36, 49, 50, 21, 40, 49, 50, 25, 52, 25, 26, 13,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 13,
		16, 65, 66, 81, 68, 81, 82, 21, 72, 81, 82, 25, 84, 25, 26, 13,
		32, 65, 66, 97, 68, 97, 98, 37, 72, 97, 98, 41, 100, 41, 42, 13,
		80, 97, 98, 49, 100, 49, 50, 21, 104, 49, 50, 25, 52, 25, 26, 13,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 13,
		16, 129, 130, 145, 132, 145, 146, 21, 136, 145, 146, 25, 148, 25, 26, 13,
		32, 129, 130, 161, 132, 161, 162, 37, 136, 161, 162, 41, 164, 41, 42, 13,
		144, 161, 162, 49, 164, 49, 50, 21, 168, 49, 50, 25, 52, 25, 26, 13,
		64, 129, 130, 193, 132, 193, 194, 69, 136, 193, 194, 73, 196, 73, 74, 13,
		144, 193, 194, 81, 196, 81, 82, 21, 200, 81, 82, 25, 84, 25, 26, 13,
		160, 193, 194, 97, 196, 97, 98, 37, 200, 97, 98, 41, 100, 41, 42, 13,
		208, 97, 98, 49, 100, 49, 50, 21, 104, 49, 50, 25, 52, 25, 26, 141,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 14,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 14,
		32, 48, 48, 50, 48, 52, 52, 22, 48, 56, 56, 26, 56, 28, 28, 14,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 14,
		64, 80, 80, 82, 80, 84, 84, 22, 80, 88, 88, 26, 88, 28, 28, 14,
		64, 96, 96, 98, 96, 100, 100, 38, 96, 104, 104, 42, 104, 44, 44, 14,
		96, 112, 112, 50, 112, 52, 52, 22, 112, 56, 56, 26, 56, 28, 28, 14,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 14,
		128, 144, 144, 146, 144, 148, 148, 22, 144, 152, 152, 26, 152, 28, 28, 14,
		128, 160, 160, 162, 160, 164, 164, 38, 160, 168, 168, 42, 168, 44, 44, 14,
		160, 176, 176, 50, 176, 52, 52, 22, 176, 56, 56, 26, 56, 28, 28, 14,
		128, 192, 192, 194, 192, 196, 196, 70, 192, 200, 200, 74, 200, 76, 76, 14,
		192, 208, 208, 82, 208, 84, 84, 22, 208, 88, 88, 26, 88, 28, 28, 14,
		192, 224, 224, 98, 224, 100, 100, 38, 224, 104, 104, 42, 104, 44, 44, 14,
		224, 112, 112, 50, 112, 52, 52, 22, 112, 56, 56, 26, 56, 28, 28, 142,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 15,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 15,
		48, 49, 50, 51, 52, 53, 54, 23, 56, 57, 58, 27, 60, 29, 30, 15,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 15,
		80, 81, 82, 83, 84, 85, 86, 23, 88, 89, 90, 27, 92, 29, 30, 15,
		96, 97, 98, 99, 100, 101, 102, 39, 104, 105, 106, 43, 108, 45, 46, 15,
		112, 113, 114, 51, 116, 53, 54, 23, 120, 57, 58, 27, 60, 29, 30, 15,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 15,
		144, 145, 146, 147, 148, 149, 150, 23, 152, 153, 154, 27, 156, 29, 30, 15,
		160, 161, 162, 163, 164, 165, 166, 39, 168, 169, 170, 43, 172, 45, 46, 15,
		176, 177, 178, 51, 180, 53, 54, 23, 184, 57, 58, 27, 60, 29, 30, 15,
		192, 193, 194, 195, 196, 197, 198, 71, 200, 201, 202, 75, 204, 77, 78, 15,
		208, 209, 210, 83, 212, 85, 86, 23, 216, 89, 90, 27, 92, 29, 30, 15,
		224, 225, 226, 99, 228, 101, 102, 39, 232, 105, 106, 43, 108, 45, 46, 15,
		240, 113, 114, 51, 116, 53, 54, 23, 120, 57, 58, 27, 60, 29, 30, 143,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 16,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 32,
		0, 0, 0, 64, 0, 64, 64, 32, 0, 64, 64, 32, 64, 32, 32, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 16,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 32,
		0, 0, 0, 128, 0, 128, 128, 32, 0, 128, 128, 32, 128, 32, 32, 16,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 64,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 16,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 32,
		0, 128, 128, 64, 128, 64, 64, 32, 128, 64, 64, 32, 64, 32, 32, 144,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 17,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 33,
		16, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 17,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		16, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 17,
		32, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 33,
		16, 1, 2, 65, 4, 65, 66, 33, 8, 65, 66, 33, 68, 33, 34, 17,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 17,
		32, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 33,
		16, 1, 2, 129, 4, 129, 130, 33, 8, 129, 130, 33, 132, 33, 34, 17,
		64, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 65,
		16, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 17,
		32, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 33,
		16, 129, 130, 65, 132, 65, 66, 33, 136, 65, 66, 33, 68, 33, 34, 145,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 18,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 34,
		32, 16, 16, 2, 16, 4, 4, 34, 16, 8, 8, 34, 8, 36, 36, 18,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 66,
		64, 16, 16, 2, 16, 4, 4, 66, 16, 8, 8, 66, 8, 68, 68, 18,
		64, 32, 32, 2, 32, 4, 4, 66, 32, 8, 8, 66, 8, 68, 68, 34,
		32, 16, 16, 66, 16, 68, 68, 34, 16, 72, 72, 34, 72, 36, 36, 18,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 130,
		128, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 18,
		128, 32, 32, 2, 32, 4, 4, 130, 32, 8, 8, 130, 8, 132, 132, 34,
		32, 16, 16, 130, 16, 132, 132, 34, 16, 136, 136, 34, 136, 36, 36, 18,
		128, 64, 64, 2, 64, 4, 4, 130, 64, 8, 8, 130, 8, 132, 132, 66,
		64, 16, 16, 130, 16, 132, 132, 66, 16, 136, 136, 66, 136, 68, 68, 18,
		64, 32, 32, 130, 32, 132, 132, 66, 32, 136, 136, 66, 136, 68, 68, 34,
		32, 144, 144, 66, 144, 68, 68, 34, 144, 72, 72, 34, 72, 36, 36, 146,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 19,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 35,
		48, 17, 18, 3, 20, 5, 6, 35, 24, 9, 10, 35, 12, 37, 38, 19,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 67,
		80, 17, 18, 3, 20, 5, 6, 67, 24, 9, 10, 67, 12, 69, 70, 19,
		96, 33, 34, 3, 36, 5, 6, 67, 40, 9, 10, 67, 12, 69, 70, 35,
		48, 17, 18, 67, 20, 69, 70, 35, 24, 73, 74, 35, 76, 37, 38, 19,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 131,
		144, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 19,
		160, 33, 34, 3, 36, 5, 6, 131, 40, 9, 10, 131, 12, 133, 134, 35,
		48, 17, 18, 131, 20, 133, 134, 35, 24, 137, 138, 35, 140, 37, 38, 19,
		192, 65, 66, 3, 68, 5, 6, 131, 72, 9, 10, 131, 12, 133, 134, 67,
		80, 17, 18, 131, 20, 133, 134, 67, 24, 137, 138, 67, 140, 69, 70, 19,
		96, 33, 34, 131, 36, 133, 134, 67, 40, 137, 138, 67, 140, 69, 70, 35,
		48, 145, 146, 67, 148, 69, 70, 35, 152, 73, 74, 35, 76, 37, 38, 147,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 20,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 36,
		0, 32, 32, 16, 32, 16, 16, 36, 32, 16, 16, 40, 16, 40, 40, 20,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 68,
		0, 64, 64, 16, 64, 16, 16, 68, 64, 16, 16, 72, 16, 72, 72, 20,
		0, 64, 64, 32, 64, 32, 32, 68, 64, 32, 32, 72, 32, 72, 72, 36,
		64, 32, 32, 80, 32, 80, 80, 36, 32, 80, 80, 40, 80, 40, 40, 20,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 132,
		0, 128, 128, 16, 128, 16, 16, 132, 128, 16, 16, 136, 16, 136, 136, 20,
		0, 128, 128, 32, 128, 32, 32, 132, 128, 32, 32, 136, 32, 136, 136, 36,
		128, 32, 32, 144, 32, 144, 144, 36, 32, 144, 144, 40, 144, 40, 40, 20,
		0, 128, 128, 64, 128, 64, 64, 132, 128, 64, 64, 136, 64, 136, 136, 68,
		128, 64, 64, 144, 64, 144, 144, 68, 64, 144, 144, 72, 144, 72, 72, 20,
		128, 64, 64, 160, 64, 160, 160, 68, 64, 160, 160, 72, 160, 72, 72, 36,
		64, 160, 160, 80, 160, 80, 80, 36, 160, 80, 80, 40, 80, 40, 40, 148,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 21,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 37,
		16, 33, 34, 17, 36, 17, 18, 37, 40, 17, 18, 41, 20, 41, 42, 21,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 69,
		16, 65, 66, 17, 68, 17, 18, 69, 72, 17, 18, 73, 20, 73, 74, 21,
		32, 65, 66, 33, 68, 33, 34, 69, 72, 33, 34, 73, 36, 73, 74, 37,
		80, 33, 34, 81, 36, 81, 82, 37, 40, 81, 82, 41, 84, 41, 42, 21,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 133,
		16, 129, 130, 17, 132, 17, 18, 133, 136, 17, 18, 137, 20, 137, 138, 21,
		32, 129, 130, 33, 132, 33, 34, 133, 136, 33, 34, 137, 36, 137, 138, 37,
		144, 33, 34, 145, 36, 145, 146, 37, 40, 145, 146, 41, 148, 41, 42, 21,
		64, 129, 130, 65, 132, 65, 66, 133, 136, 65, 66, 137, 68, 137, 138, 69,
		144, 65, 66, 145, 68, 145, 146, 69, 72, 145, 146, 73, 148, 73, 74, 21,
		160, 65, 66, 161, 68, 161, 162, 69, 72, 161, 162, 73, 164, 73, 74, 37,
		80, 161, 162, 81, 164, 81, 82, 37, 168, 81, 82, 41, 84, 41, 42, 149,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 22,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 38,
		32, 48, 48, 18, 48, 20, 20, 38, 48, 24, 24, 42, 24, 44, 44, 22,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 70,
		64, 80, 80, 18, 80, 20, 20, 70, 80, 24, 24, 74, 24, 76, 76, 22,
		64, 96, 96, 34, 96, 36, 36, 70, 96, 40, 40, 74, 40, 76, 76, 38,
		96, 48, 48, 82, 48, 84, 84, 38, 48, 88, 88, 42, 88, 44, 44, 22,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 134,
		128, 144, 144, 18, 144, 20, 20, 134, 144, 24, 24, 138, 24, 140, 140, 22,
		128, 160, 160, 34, 160, 36, 36, 134, 160, 40, 40, 138, 40, 140, 140, 38,
		160, 48, 48, 146, 48, 148, 148, 38, 48, 152, 152, 42, 152, 44, 44, 22,
		128, 192, 192, 66, 192, 68, 68, 134, 192, 72, 72, 138, 72, 140, 140, 70,
		192, 80, 80, 146, 80, 148, 148, 70, 80, 152, 152, 74, 152, 76, 76, 22,
		192, 96, 96, 162, 96, 164, 164, 70, 96, 168, 168, 74, 168, 76, 76, 38,
		96, 176, 176, 82, 176, 84, 84, 38, 176, 88, 88, 42, 88, 44, 44, 150,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 23,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 39,
		48, 49, 50, 19, 52, 21, 22, 39, 56, 25, 26, 43, 28, 45, 46, 23,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 71,
		80, 81, 82, 19, 84, 21, 22, 71, 88, 25, 26, 75, 28, 77, 78, 23,
		96, 97, 98, 35, 100, 37, 38, 71, 104, 41, 42, 75, 44, 77, 78, 39,
		112, 49, 50, 83, 52, 85, 86, 39, 56, 89, 90, 43, 92, 45, 46, 23,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 135,
		144, 145, 146, 19, 148, 21, 22, 135, 152, 25, 26, 139, 28, 141, 142, 23,
		160, 161, 162, 35, 164, 37, 38, 135, 168, 41, 42, 139, 44, 141, 142, 39,
		176, 49, 50, 147, 52, 149, 150, 39, 56, 153, 154, 43, 156, 45, 46, 23,
		192, 193, 194, 67, 196, 69, 70, 135, 200, 73, 74, 139, 76, 141, 142, 71,
		208, 81, 82, 147, 84, 149, 150, 71, 88, 153, 154, 75, 156, 77, 78, 23,
		224, 97, 98, 163, 100, 165, 166, 71, 104, 169, 170, 75, 172, 77, 78, 39,
		112, 177, 178, 83, 180, 85, 86, 39, 184, 89, 90, 43, 92, 45, 46, 151,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 24,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 40,
		0, 0, 0, 32, 0, 32, 32, 48, 0, 32, 32, 48, 32, 48, 48, 24,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 72,
		0, 0, 0, 64, 0, 64, 64, 80, 0, 64, 64, 80, 64, 80, 80, 24,
		0, 0, 0, 64, 0, 64, 64, 96, 0, 64, 64, 96, 64, 96, 96, 40,
		0, 64, 64, 96, 64, 96, 96, 48, 64, 96, 96, 48, 96, 48, 48, 24,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 136,
		0, 0, 0, 128, 0, 128, 128, 144, 0, 128, 128, 144, 128, 144, 144, 24,
		0, 0, 0, 128, 0, 128, 128, 160, 0, 128, 128, 160, 128, 160, 160, 40,
		0, 128, 128, 160, 128, 160, 160, 48, 128, 160, 160, 48, 160, 48, 48, 24,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 72,
		0, 128, 128, 192, 128, 192, 192, 80, 128, 192, 192, 80, 192, 80, 80, 24,
		0, 128, 128, 192, 128, 192, 192, 96, 128, 192, 192, 96, 192, 96, 96, 40,
		128, 192, 192, 96, 192, 96, 96, 48, 192, 96, 96, 48, 96, 48, 48, 152,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 25,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 41,
		16, 1, 2, 33, 4, 33, 34, 49, 8, 33, 34, 49, 36, 49, 50, 25,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 73,
		16, 1, 2, 65, 4, 65, 66, 81, 8, 65, 66, 81, 68, 81, 82, 25,
		32, 1, 2, 65, 4, 65, 66, 97, 8, 65, 66, 97, 68, 97, 98, 41,
		16, 65, 66, 97, 68, 97, 98, 49, 72, 97, 98, 49, 100, 49, 50, 25,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 137,
		16, 1, 2, 129, 4, 129, 130, 145, 8, 129, 130, 145, 132, 145, 146, 25,
		32, 1, 2, 129, 4, 129, 130, 161, 8, 129, 130, 161, 132, 161, 162, 41,
		16, 129, 130, 161, 132, 161, 162, 49, 136, 161, 162, 49, 164, 49, 50, 25,
		64, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 73,
		16, 129, 130, 193, 132, 193, 194, 81, 136, 193, 194, 81, 196, 81, 82, 25,
		32, 129, 130, 193, 132, 193, 194, 97, 136, 193, 194, 97, 196, 97, 98, 41,
		144, 193, 194, 97, 196, 97, 98, 49, 200, 97, 98, 49, 100, 49, 50, 153,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 26,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 42,
		32, 16, 16, 34, 16, 36, 36, 50, 16, 40, 40, 50, 40, 52, 52, 26,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 74,
		64, 16, 16, 66, 16, 68, 68, 82, 16, 72, 72, 82, 72, 84, 84, 26,
		64, 32, 32, 66, 32, 68, 68, 98, 32, 72, 72, 98, 72, 100, 100, 42,
		32, 80, 80, 98, 80, 100, 100, 50, 80, 104, 104, 50, 104, 52, 52, 26,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 138,
		128, 16, 16, 130, 16, 132, 132, 146, 16, 136, 136, 146, 136, 148, 148, 26,
		128, 32, 32, 130, 32, 132, 132, 162, 32, 136, 136, 162, 136, 164, 164, 42,
		32, 144, 144, 162, 144, 164, 164, 50, 144, 168, 168, 50, 168, 52, 52, 26,
		128, 64, 64, 130, 64, 132, 132, 194, 64, 136, 136, 194, 136, 196, 196, 74,
		64, 144, 144, 194, 144, 196, 196, 82, 144, 200, 200, 82, 200, 84, 84, 26,
		64, 160, 160, 194, 160, 196, 196, 98, 160, 200, 200, 98, 200, 100, 100, 42,
		160, 208, 208, 98, 208, 100, 100, 50, 208, 104, 104, 50, 104, 52, 52, 154,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 27,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 43,
		48, 17, 18, 35, 20, 37, 38, 51, 24, 41, 42, 51, 44, 53, 54, 27,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 75,
		80, 17, 18, 67, 20, 69, 70, 83, 24, 73, 74, 83, 76, 85, 86, 27,
		96, 33, 34, 67, 36, 69, 70, 99, 40, 73, 74, 99, 76, 101, 102, 43,
		48, 81, 82, 99, 84, 101, 102, 51, 88, 105, 106, 51, 108, 53, 54, 27,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 139,
		144, 17, 18, 131, 20, 133, 134, 147, 24, 137, 138, 147, 140, 149, 150, 27,
		160, 33, 34, 131, 36, 133, 134, 163, 40, 137, 138, 163, 140, 165, 166, 43,
		48, 145, 146, 163, 148, 165, 166, 51, 152, 169, 170, 51, 172, 53, 54, 27,
		192, 65, 66, 131, 68, 133, 134, 195, 72, 137, 138, 195, 140, 197, 198, 75,
		80, 145, 146, 195, 148, 197, 198, 83, 152, 201, 202, 83, 204, 85, 86, 27,
		96, 161, 162, 195, 164, 197, 198, 99, 168, 201, 202, 99, 204, 101, 102, 43,
		176, 209, 210, 99, 212, 101, 102, 51, 216, 105, 106, 51, 108, 53, 54, 155,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 28,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 44,
		0, 32, 32, 48, 32, 48, 48, 52, 32, 48, 48, 56, 48, 56, 56, 28,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 76,
		0, 64, 64, 80, 64, 80, 80, 84, 64, 80, 80, 88, 80, 88, 88, 28,
		0, 64, 64, 96, 64, 96, 96, 100, 64, 96, 96, 104, 96, 104, 104, 44,
		64, 96, 96, 112, 96, 112, 112, 52, 96, 112, 112, 56, 112, 56, 56, 28,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 140,
		0, 128, 128, 144, 128, 144, 144, 148, 128, 144, 144, 152, 144, 152, 152, 28,
		0, 128, 128, 160, 128, 160, 160, 164, 128, 160, 160, 168, 160, 168, 168, 44,
		128, 160, 160, 176, 160, 176, 176, 52, 160, 176, 176, 56, 176, 56, 56, 28,
		0, 128, 128, 192, 128, 192, 192, 196, 128, 192, 192, 200, 192, 200, 200, 76,
		128, 192, 192, 208, 192, 208, 208, 84, 192, 208, 208, 88, 208, 88, 88, 28,
		128, 192, 192, 224, 192, 224, 224, 100, 192, 224, 224, 104, 224, 104, 104, 44,
		192, 224, 224, 112, 224, 112, 112, 52, 224, 112, 112, 56, 112, 56, 56, 156,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 29,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 45,
		16, 33, 34, 49, 36, 49, 50, 53, 40, 49, 50, 57, 52, 57, 58, 29,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 77,
		16, 65, 66, 81, 68, 81, 82, 85, 72, 81, 82, 89, 84, 89, 90, 29,
		32, 65, 66, 97, 68, 97, 98, 101, 72, 97, 98, 105, 100, 105, 106, 45,
		80, 97, 98, 113, 100, 113, 114, 53, 104, 113, 114, 57, 116, 57, 58, 29,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 141,
		16, 129, 130, 145, 132, 145, 146, 149, 136, 145, 146, 153, 148, 153, 154, 29,
		32, 129, 130, 161, 132, 161, 162, 165, 136, 161, 162, 169, 164, 169, 170, 45,
		144, 161, 162, 177, 164, 177, 178, 53, 168, 177, 178, 57, 180, 57, 58, 29,
		64, 129, 130, 193, 132, 193, 194, 197, 136, 193, 194, 201, 196, 201, 202, 77,
		144, 193, 194, 209, 196, 209, 210, 85, 200, 209, 210, 89, 212, 89, 90, 29,
		160, 193, 194, 225, 196, 225, 226, 101, 200, 225, 226, 105, 228, 105, 106, 45,
		208, 225, 226, 113, 228, 113, 114, 53, 232, 113, 114, 57, 116, 57, 58, 157,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 30,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 46,
		32, 48, 48, 50, 48, 52, 52, 54, 48, 56, 56, 58, 56, 60, 60, 30,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 78,
		64, 80, 80, 82, 80, 84, 84, 86, 80, 88, 88, 90, 88, 92, 92, 30,
		64, 96, 96, 98, 96, 100, 100, 102, 96, 104, 104, 106, 104, 108, 108, 46,
		96, 112, 112, 114, 112, 116, 116, 54, 112, 120, 120, 58, 120, 60, 60, 30,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 142,
		128, 144, 144, 146, 144, 148, 148, 150, 144, 152, 152, 154, 152, 156, 156, 30,
		128, 160, 160, 162, 160, 164, 164, 166, 160, 168, 168, 170, 168, 172, 172, 46,
		160, 176, 176, 178, 176, 180, 180, 54, 176, 184, 184, 58, 184, 60, 60, 30,
		128, 192, 192, 194, 192, 196, 196, 198, 192, 200, 200, 202, 200, 204, 204, 78,
		192, 208, 208, 210, 208, 212, 212, 86, 208, 216, 216, 90, 216, 92, 92, 30,
		192, 224, 224, 226, 224, 228, 228, 102, 224, 232, 232, 106, 232, 108, 108, 46,
		224, 240, 240, 114, 240, 116, 116, 54, 240, 120, 120, 58, 120, 60, 60, 158,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 31,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
		80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 31,
		96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 47,
		112, 113, 114, 115, 116, 117, 118, 55, 120, 121, 122, 59, 124, 61, 62, 31,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143,
		144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 31,
		160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 47,
		176, 177, 178, 179, 180, 181, 182, 55, 184, 185, 186, 59, 188, 61, 62, 31,
		192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 79,
		208, 209, 210, 211, 212, 213, 214, 87, 216, 217, 218, 91, 220, 93, 94, 31,
		224, 225, 226, 227, 228, 229, 230, 103, 232, 233, 234, 107, 236, 109, 110, 47,
		240, 241, 242, 115, 244, 117, 118, 55, 248, 121, 122, 59, 124, 61, 62, 159,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 32,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 32,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 64,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 64,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 160,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 33,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		16, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 33,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 33,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 65,
		32, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 65,
		16, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 161,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 34,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 2,
		64, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 66,
		64, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 66,
		32, 16, 16, 2, 16, 4, 4, 66, 16, 8, 8, 66, 8, 68, 68, 34,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 2,
		128, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 130,
		128, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 130,
		32, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 34,
		128, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 130,
		64, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 66,
		64, 32, 32, 2, 32, 4, 4, 130, 32, 8, 8, 130, 8, 132, 132, 66,
		32, 16, 16, 130, 16, 132, 132, 66, 16, 136, 136, 66, 136, 68, 68, 162,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 35,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 3,
		80, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 67,
		96, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 67,
		48, 17, 18, 3, 20, 5, 6, 67, 24, 9, 10, 67, 12, 69, 70, 35,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 3,
		144, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 131,
		160, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 131,
		48, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 35,
		192, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 131,
		80, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 67,
		96, 33, 34, 3, 36, 5, 6, 131, 40, 9, 10, 131, 12, 133, 134, 67,
		48, 17, 18, 131, 20, 133, 134, 67, 24, 137, 138, 67, 140, 69, 70, 163,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 4,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 4,
		0, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 36,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 4,
		0, 64, 64, 16, 64, 16, 16, 4, 64, 16, 16, 8, 16, 8, 8, 68,
		0, 64, 64, 32, 64, 32, 32, 4, 64, 32, 32, 8, 32, 8, 8, 68,
		64, 32, 32, 16, 32, 16, 16, 68, 32, 16, 16, 72, 16, 72, 72, 36,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 4,
		0, 128, 128, 16, 128, 16, 16, 4, 128, 16, 16, 8, 16, 8, 8, 132,
		0, 128, 128, 32, 128, 32, 32, 4, 128, 32, 32, 8, 32, 8, 8, 132,
		128, 32, 32, 16, 32, 16, 16, 132, 32, 16, 16, 136, 16, 136, 136, 36,
		0, 128, 128, 64, 128, 64, 64, 4, 128, 64, 64, 8, 64, 8, 8, 132,
		128, 64, 64, 16, 64, 16, 16, 132, 64, 16, 16, 136, 16, 136, 136, 68,
		128, 64, 64, 32, 64, 32, 32, 132, 64, 32, 32, 136, 32, 136, 136, 68,
		64, 32, 32, 144, 32, 144, 144, 68, 32, 144, 144, 72, 144, 72, 72, 164,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 5,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 5,
		16, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 37,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 5,
		16, 65, 66, 17, 68, 17, 18, 5, 72, 17, 18, 9, 20, 9, 10, 69,
		32, 65, 66, 33, 68, 33, 34, 5, 72, 33, 34, 9, 36, 9, 10, 69,
		80, 33, 34, 17, 36, 17, 18, 69, 40, 17, 18, 73, 20, 73, 74, 37,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 5,
		16, 129, 130, 17, 132, 17, 18, 5, 136, 17, 18, 9, 20, 9, 10, 133,
		32, 129, 130, 33, 132, 33, 34, 5, 136, 33, 34, 9, 36, 9, 10, 133,
		144, 33, 34, 17, 36, 17, 18, 133, 40, 17, 18, 137, 20, 137, 138, 37,
		64, 129, 130, 65, 132, 65, 66, 5, 136, 65, 66, 9, 68, 9, 10, 133,
		144, 65, 66, 17, 68, 17, 18, 133, 72, 17, 18, 137, 20, 137, 138, 69,
		160, 65, 66, 33, 68, 33, 34, 133, 72, 33, 34, 137, 36, 137, 138, 69,
		80, 33, 34, 145, 36, 145, 146, 69, 40, 145, 146, 73, 148, 73, 74, 165,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 6,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 6,
		32, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 38,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 6,
		64, 80, 80, 18, 80, 20, 20, 6, 80, 24, 24, 10, 24, 12, 12, 70,
		64, 96, 96, 34, 96, 36, 36, 6, 96, 40, 40, 10, 40, 12, 12, 70,
		96, 48, 48, 18, 48, 20, 20, 70, 48, 24, 24, 74, 24, 76, 76, 38,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 6,
		128, 144, 144, 18, 144, 20, 20, 6, 144, 24, 24, 10, 24, 12, 12, 134,
		128, 160, 160, 34, 160, 36, 36, 6, 160, 40, 40, 10, 40, 12, 12, 134,
		160, 48, 48, 18, 48, 20, 20, 134, 48, 24, 24, 138, 24, 140, 140, 38,
		128, 192, 192, 66, 192, 68, 68, 6, 192, 72, 72, 10, 72, 12, 12, 134,
		192, 80, 80, 18, 80, 20, 20, 134, 80, 24, 24, 138, 24, 140, 140, 70,
		192, 96, 96, 34, 96, 36, 36, 134, 96, 40, 40, 138, 40, 140, 140, 70,
		96, 48, 48, 146, 48, 148, 148, 70, 48, 152, 152, 74, 152, 76, 76, 166,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 7,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 7,
		48, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 39,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 7,
		80, 81, 82, 19, 84, 21, 22, 7, 88, 25, 26, 11, 28, 13, 14, 71,
		96, 97, 98, 35, 100, 37, 38, 7, 104, 41, 42, 11, 44, 13, 14, 71,
		112, 49, 50, 19, 52, 21, 22, 71, 56, 25, 26, 75, 28, 77, 78, 39,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 7,
		144, 145, 146, 19, 148, 21, 22, 7, 152, 25, 26, 11, 28, 13, 14, 135,
		160, 161, 162, 35, 164, 37, 38, 7, 168, 41, 42, 11, 44, 13, 14, 135,
		176, 49, 50, 19, 52, 21, 22, 135, 56, 25, 26, 139, 28, 141, 142, 39,
		192, 193, 194, 67, 196, 69, 70, 7, 200, 73, 74, 11, 76, 13, 14, 135,
		208, 81, 82, 19, 84, 21, 22, 135, 88, 25, 26, 139, 28, 141, 142, 71,
		224, 97, 98, 35, 100, 37, 38, 135, 104, 41, 42, 139, 44, 141, 142, 71,
		112, 49, 50, 147, 52, 149, 150, 71, 56, 153, 154, 75, 156, 77, 78, 167,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 8,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 8,
		0, 0, 0, 32, 0, 32, 32, 16, 0, 32, 32, 16, 32, 16, 16, 40,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 8,
		0, 0, 0, 64, 0, 64, 64, 16, 0, 64, 64, 16, 64, 16, 16, 72,
		0, 0, 0, 64, 0, 64, 64, 32, 0, 64, 64, 32, 64, 32, 32, 72,
		0, 64, 64, 32, 64, 32, 32, 80, 64, 32, 32, 80, 32, 80, 80, 40,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 8,
		0, 0, 0, 128, 0, 128, 128, 16, 0, 128, 128, 16, 128, 16, 16, 136,
		0, 0, 0, 128, 0, 128, 128, 32, 0, 128, 128, 32, 128, 32, 32, 136,
		0, 128, 128, 32, 128, 32, 32, 144, 128, 32, 32, 144, 32, 144, 144, 40,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 136,
		0, 128, 128, 64, 128, 64, 64, 144, 128, 64, 64, 144, 64, 144, 144, 72,
		0, 128, 128, 64, 128, 64, 64, 160, 128, 64, 64, 160, 64, 160, 160, 72,
		128, 64, 64, 160, 64, 160, 160, 80, 64, 160, 160, 80, 160, 80, 80, 168,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 9,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 9,
		16, 1, 2, 33, 4, 33, 34, 17, 8, 33, 34, 17, 36, 17, 18, 41,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 9,
		16, 1, 2, 65, 4, 65, 66, 17, 8, 65, 66, 17, 68, 17, 18, 73,
		32, 1, 2, 65, 4, 65, 66, 33, 8, 65, 66, 33, 68, 33, 34, 73,
		16, 65, 66, 33, 68, 33, 34, 81, 72, 33, 34, 81, 36, 81, 82, 41,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 9,
		16, 1, 2, 129, 4, 129, 130, 17, 8, 129, 130, 17, 132, 17, 18, 137,
		32, 1, 2, 129, 4, 129, 130, 33, 8, 129, 130, 33, 132, 33, 34, 137,
		16, 129, 130, 33, 132, 33, 34, 145, 136, 33, 34, 145, 36, 145, 146, 41,
		64, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 137,
		16, 129, 130, 65, 132, 65, 66, 145, 136, 65, 66, 145, 68, 145, 146, 73,
		32, 129, 130, 65, 132, 65, 66, 161, 136, 65, 66, 161, 68, 161, 162, 73,
		144, 65, 66, 161, 68, 161, 162, 81, 72, 161, 162, 81, 164, 81, 82, 169,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 10,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 10,
		32, 16, 16, 34, 16, 36, 36, 18, 16, 40, 40, 18, 40, 20, 20, 42,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 10,
		64, 16, 16, 66, 16, 68, 68, 18, 16, 72, 72, 18, 72, 20, 20, 74,
		64, 32, 32, 66, 32, 68, 68, 34, 32, 72, 72, 34, 72, 36, 36, 74,
		32, 80, 80, 34, 80, 36, 36, 82, 80, 40, 40, 82, 40, 84, 84, 42,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 10,
		128, 16, 16, 130, 16, 132, 132, 18, 16, 136, 136, 18, 136, 20, 20, 138,
		128, 32, 32, 130, 32, 132, 132, 34, 32, 136, 136, 34, 136, 36, 36, 138,
		32, 144, 144, 34, 144, 36, 36, 146, 144, 40, 40, 146, 40, 148, 148, 42,
		128, 64, 64, 130, 64, 132, 132, 66, 64, 136, 136, 66, 136, 68, 68, 138,
		64, 144, 144, 66, 144, 68, 68, 146, 144, 72, 72, 146, 72, 148, 148, 74,
		64, 160, 160, 66, 160, 68, 68, 162, 160, 72, 72, 162, 72, 164, 164, 74,
		160, 80, 80, 162, 80, 164, 164, 82, 80, 168, 168, 82, 168, 84, 84, 170,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 11,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 11,
		48, 17, 18, 35, 20, 37, 38, 19, 24, 41, 42, 19, 44, 21, 22, 43,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 11,
		80, 17, 18, 67, 20, 69, 70, 19, 24, 73, 74, 19, 76, 21, 22, 75,
		96, 33, 34, 67, 36, 69, 70, 35, 40, 73, 74, 35, 76, 37, 38, 75,
		48, 81, 82, 35, 84, 37, 38, 83, 88, 41, 42, 83, 44, 85, 86, 43,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 11,
		144, 17, 18, 131, 20, 133, 134, 19, 24, 137, 138, 19, 140, 21, 22, 139,
		160, 33, 34, 131, 36, 133, 134, 35, 40, 137, 138, 35, 140, 37, 38, 139,
		48, 145, 146, 35, 148, 37, 38, 147, 152, 41, 42, 147, 44, 149, 150, 43,
		192, 65, 66, 131, 68, 133, 134, 67, 72, 137, 138, 67, 140, 69, 70, 139,
		80, 145, 146, 67, 148, 69, 70, 147, 152, 73, 74, 147, 76, 149, 150, 75,
		96, 161, 162, 67, 164, 69, 70, 163, 168, 73, 74, 163, 76, 165, 166, 75,
		176, 81, 82, 163, 84, 165, 166, 83, 88, 169, 170, 83, 172, 85, 86, 171,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 12,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 12,
		0, 32, 32, 48, 32, 48, 48, 20, 32, 48, 48, 24, 48, 24, 24, 44,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 12,
		0, 64, 64, 80, 64, 80, 80, 20, 64, 80, 80, 24, 80, 24, 24, 76,
		0, 64, 64, 96, 64, 96, 96, 36, 64, 96, 96, 40, 96, 40, 40, 76,
		64, 96, 96, 48, 96, 48, 48, 84, 96, 48, 48, 88, 48, 88, 88, 44,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 12,
		0, 128, 128, 144, 128, 144, 144, 20, 128, 144, 144, 24, 144, 24, 24, 140,
		0, 128, 128, 160, 128, 160, 160, 36, 128, 160, 160, 40, 160, 40, 40, 140,
		128, 160, 160, 48, 160, 48, 48, 148, 160, 48, 48, 152, 48, 152, 152, 44,
		0, 128, 128, 192, 128, 192, 192, 68, 128, 192, 192, 72, 192, 72, 72, 140,
		128, 192, 192, 80, 192, 80, 80, 148, 192, 80, 80, 152, 80, 152, 152, 76,
		128, 192, 192, 96, 192, 96, 96, 164, 192, 96, 96, 168, 96, 168, 168, 76,
		192, 96, 96, 176, 96, 176, 176, 84, 96, 176, 176, 88, 176, 88, 88, 172,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 13,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 13,
		16, 33, 34, 49, 36, 49, 50, 21, 40, 49, 50, 25, 52, 25, 26, 45,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 13,
		16, 65, 66, 81, 68, 81, 82, 21, 72, 81, 82, 25, 84, 25, 26, 77,
		32, 65, 66, 97, 68, 97, 98, 37, 72, 97, 98, 41, 100, 41, 42, 77,
		80, 97, 98, 49, 100, 49, 50, 85, 104, 49, 50, 89, 52, 89, 90, 45,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 13,
		16, 129, 130, 145, 132, 145, 146, 21, 136, 145, 146, 25, 148, 25, 26, 141,
		32, 129, 130, 161, 132, 161, 162, 37, 136, 161, 162, 41, 164, 41, 42, 141,
		144, 161, 162, 49, 164, 49, 50, 149, 168, 49, 50, 153, 52, 153, 154, 45,
		64, 129, 130, 193, 132, 193, 194, 69, 136, 193, 194, 73, 196, 73, 74, 141,
		144, 193, 194, 81, 196, 81, 82, 149, 200, 81, 82, 153, 84, 153, 154, 77,
		160, 193, 194, 97, 196, 97, 98, 165, 200, 97, 98, 169, 100, 169, 170, 77,
		208, 97, 98, 177, 100, 177, 178, 85, 104, 177, 178, 89, 180, 89, 90, 173,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 14,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 14,
		32, 48, 48, 50, 48, 52, 52, 22, 48, 56, 56, 26, 56, 28, 28, 46,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 14,
		64, 80, 80, 82, 80, 84, 84, 22, 80, 88, 88, 26, 88, 28, 28, 78,
		64, 96, 96, 98, 96, 100, 100, 38, 96, 104, 104, 42, 104, 44, 44, 78,
		96, 112, 112, 50, 112, 52, 52, 86, 112, 56, 56, 90, 56, 92, 92, 46,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 14,
		128, 144, 144, 146, 144, 148, 148, 22, 144, 152, 152, 26, 152, 28, 28, 142,
		128, 160, 160, 162, 160, 164, 164, 38, 160, 168, 168, 42, 168, 44, 44, 142,
		160, 176, 176, 50, 176, 52, 52, 150, 176, 56, 56, 154, 56, 156, 156, 46,
		128, 192, 192, 194, 192, 196, 196, 70, 192, 200, 200, 74, 200, 76, 76, 142,
		192, 208, 208, 82, 208, 84, 84, 150, 208, 88, 88, 154, 88, 156, 156, 78,
		192, 224, 224, 98, 224, 100, 100, 166, 224, 104, 104, 170, 104, 172, 172, 78,
		224, 112, 112, 178, 112, 180, 180, 86, 112, 184, 184, 90, 184, 92, 92, 174,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 15,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 15,
		48, 49, 50, 51, 52, 53, 54, 23, 56, 57, 58, 27, 60, 29, 30, 47,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 15,
		80, 81, 82, 83, 84, 85, 86, 23, 88, 89, 90, 27, 92, 29, 30, 79,
		96, 97, 98, 99, 100, 101, 102, 39, 104, 105, 106, 43, 108, 45, 46, 79,
		112, 113, 114, 51, 116, 53, 54, 87, 120, 57, 58, 91, 60, 93, 94, 47,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 15,
		144, 145, 146, 147, 148, 149, 150, 23, 152, 153, 154, 27, 156, 29, 30, 143,
		160, 161, 162, 163, 164, 165, 166, 39, 168, 169, 170, 43, 172, 45, 46, 143,
		176, 177, 178, 51, 180, 53, 54, 151, 184, 57, 58, 155, 60, 157, 158, 47,
		192, 193, 194, 195, 196, 197, 198, 71, 200, 201, 202, 75, 204, 77, 78, 143,
		208, 209, 210, 83, 212, 85, 86, 151, 216, 89, 90, 155, 92, 157, 158, 79,
		224, 225, 226, 99, 228, 101, 102, 167, 232, 105, 106, 171, 108, 173, 174, 79,
		240, 113, 114, 179, 116, 181, 182, 87, 120, 185, 186, 91, 188, 93, 94, 175,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 48,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 80,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 96,
		0, 0, 0, 64, 0, 64, 64, 96, 0, 64, 64, 96, 64, 96, 96, 48,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 144,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 160,
		0, 0, 0, 128, 0, 128, 128, 160, 0, 128, 128, 160, 128, 160, 160, 48,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 192,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 80,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 96,
		0, 128, 128, 192, 128, 192, 192, 96, 128, 192, 192, 96, 192, 96, 96, 176,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 17,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 33,
		16, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 49,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		16, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 81,
		32, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 97,
		16, 1, 2, 65, 4, 65, 66, 97, 8, 65, 66, 97, 68, 97, 98, 49,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 145,
		32, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 161,
		16, 1, 2, 129, 4, 129, 130, 161, 8, 129, 130, 161, 132, 161, 162, 49,
		64, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 193,
		16, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 81,
		32, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 97,
		16, 129, 130, 193, 132, 193, 194, 97, 136, 193, 194, 97, 196, 97, 98, 177,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 18,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 34,
		32, 16, 16, 2, 16, 4, 4, 34, 16, 8, 8, 34, 8, 36, 36, 50,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 66,
		64, 16, 16, 2, 16, 4, 4, 66, 16, 8, 8, 66, 8, 68, 68, 82,
		64, 32, 32, 2, 32, 4, 4, 66, 32, 8, 8, 66, 8, 68, 68, 98,
		32, 16, 16, 66, 16, 68, 68, 98, 16, 72, 72, 98, 72, 100, 100, 50,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 130,
		128, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 146,
		128, 32, 32, 2, 32, 4, 4, 130, 32, 8, 8, 130, 8, 132, 132, 162,
		32, 16, 16, 130, 16, 132, 132, 162, 16, 136, 136, 162, 136, 164, 164, 50,
		128, 64, 64, 2, 64, 4, 4, 130, 64, 8, 8, 130, 8, 132, 132, 194,
		64, 16, 16, 130, 16, 132, 132, 194, 16, 136, 136, 194, 136, 196, 196, 82,
		64, 32, 32, 130, 32, 132, 132, 194, 32, 136, 136, 194, 136, 196, 196, 98,
		32, 144, 144, 194, 144, 196, 196, 98, 144, 200, 200, 98, 200, 100, 100, 178,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 19,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 35,
		48, 17, 18, 3, 20, 5, 6, 35, 24, 9, 10, 35, 12, 37, 38, 51,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 67,
		80, 17, 18, 3, 20, 5, 6, 67, 24, 9, 10, 67, 12, 69, 70, 83,
		96, 33, 34, 3, 36, 5, 6, 67, 40, 9, 10, 67, 12, 69, 70, 99,
		48, 17, 18, 67, 20, 69, 70, 99, 24, 73, 74, 99, 76, 101, 102, 51,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 131,
		144, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 147,
		160, 33, 34, 3, 36, 5, 6, 131, 40, 9, 10, 131, 12, 133, 134, 163,
		48, 17, 18, 131, 20, 133, 134, 163, 24, 137, 138, 163, 140, 165, 166, 51,
		192, 65, 66, 3, 68, 5, 6, 131, 72, 9, 10, 131, 12, 133, 134, 195,
		80, 17, 18, 131, 20, 133, 134, 195, 24, 137, 138, 195, 140, 197, 198, 83,
		96, 33, 34, 131, 36, 133, 134, 195, 40, 137, 138, 195, 140, 197, 198, 99,
		48, 145, 146, 195, 148, 197, 198, 99, 152, 201, 202, 99, 204, 101, 102, 179,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 20,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 36,
		0, 32, 32, 16, 32, 16, 16, 36, 32, 16, 16, 40, 16, 40, 40, 52,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 68,
		0, 64, 64, 16, 64, 16, 16, 68, 64, 16, 16, 72, 16, 72, 72, 84,
		0, 64, 64, 32, 64, 32, 32, 68, 64, 32, 32, 72, 32, 72, 72, 100,
		64, 32, 32, 80, 32, 80, 80, 100, 32, 80, 80, 104, 80, 104, 104, 52,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 132,
		0, 128, 128, 16, 128, 16, 16, 132, 128, 16, 16, 136, 16, 136, 136, 148,
		0, 128, 128, 32, 128, 32, 32, 132, 128, 32, 32, 136, 32, 136, 136, 164,
		128, 32, 32, 144, 32, 144, 144, 164, 32, 144, 144, 168, 144, 168, 168, 52,
		0, 128, 128, 64, 128, 64, 64, 132, 128, 64, 64, 136, 64, 136, 136, 196,
		128, 64, 64, 144, 64, 144, 144, 196, 64, 144, 144, 200, 144, 200, 200, 84,
		128, 64, 64, 160, 64, 160, 160, 196, 64, 160, 160, 200, 160, 200, 200, 100,
		64, 160, 160, 208, 160, 208, 208, 100, 160, 208, 208, 104, 208, 104, 104, 180,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 21,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 37,
		16, 33, 34, 17, 36, 17, 18, 37, 40, 17, 18, 41, 20, 41, 42, 53,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 69,
		16, 65, 66, 17, 68, 17, 18, 69, 72, 17, 18, 73, 20, 73, 74, 85,
		32, 65, 66, 33, 68, 33, 34, 69, 72, 33, 34, 73, 36, 73, 74, 101,
		80, 33, 34, 81, 36, 81, 82, 101, 40, 81, 82, 105, 84, 105, 106, 53,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 133,
		16, 129, 130, 17, 132, 17, 18, 133, 136, 17, 18, 137, 20, 137, 138, 149,
		32, 129, 130, 33, 132, 33, 34, 133, 136, 33, 34, 137, 36, 137, 138, 165,
		144, 33, 34, 145, 36, 145, 146, 165, 40, 145, 146, 169, 148, 169, 170, 53,
		64, 129, 130, 65, 132, 65, 66, 133, 136, 65, 66, 137, 68, 137, 138, 197,
		144, 65, 66, 145, 68, 145, 146, 197, 72, 145, 146, 201, 148, 201, 202, 85,
		160, 65, 66, 161, 68, 161, 162, 197, 72, 161, 162, 201, 164, 201, 202, 101,
		80, 161, 162, 209, 164, 209, 210, 101, 168, 209, 210, 105, 212, 105, 106, 181,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 22,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 38,
		32, 48, 48, 18, 48, 20, 20, 38, 48, 24, 24, 42, 24, 44, 44, 54,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 70,
		64, 80, 80, 18, 80, 20, 20, 70, 80, 24, 24, 74, 24, 76, 76, 86,
		64, 96, 96, 34, 96, 36, 36, 70, 96, 40, 40, 74, 40, 76, 76, 102,
		96, 48, 48, 82, 48, 84, 84, 102, 48, 88, 88, 106, 88, 108, 108, 54,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 134,
		128, 144, 144, 18, 144, 20, 20, 134, 144, 24, 24, 138, 24, 140, 140, 150,
		128, 160, 160, 34, 160, 36, 36, 134, 160, 40, 40, 138, 40, 140, 140, 166,
		160, 48, 48, 146, 48, 148, 148, 166, 48, 152, 152, 170, 152, 172, 172, 54,
		128, 192, 192, 66, 192, 68, 68, 134, 192, 72, 72, 138, 72, 140, 140, 198,
		192, 80, 80, 146, 80, 148, 148, 198, 80, 152, 152, 202, 152, 204, 204, 86,
		192, 96, 96, 162, 96, 164, 164, 198, 96, 168, 168, 202, 168, 204, 204, 102,
		96, 176, 176, 210, 176, 212, 212, 102, 176, 216, 216, 106, 216, 108, 108, 182,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 23,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 39,
		48, 49, 50, 19, 52, 21, 22, 39, 56, 25, 26, 43, 28, 45, 46, 55,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 71,
		80, 81, 82, 19, 84, 21, 22, 71, 88, 25, 26, 75, 28, 77, 78, 87,
		96, 97, 98, 35, 100, 37, 38, 71, 104, 41, 42, 75, 44, 77, 78, 103,
		112, 49, 50, 83, 52, 85, 86, 103, 56, 89, 90, 107, 92, 109, 110, 55,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 135,
		144, 145, 146, 19, 148, 21, 22, 135, 152, 25, 26, 139, 28, 141, 142, 151,
		160, 161, 162, 35, 164, 37, 38, 135, 168, 41, 42, 139, 44, 141, 142, 167,
		176, 49, 50, 147, 52, 149, 150, 167, 56, 153, 154, 171, 156, 173, 174, 55,
		192, 193, 194, 67, 196, 69, 70, 135, 200, 73, 74, 139, 76, 141, 142, 199,
		208, 81, 82, 147, 84, 149, 150, 199, 88, 153, 154, 203, 156, 205, 206, 87,
		224, 97, 98, 163, 100, 165, 166, 199, 104, 169, 170, 203, 172, 205, 206, 103,
		112, 177, 178, 211, 180, 213, 214, 103, 184, 217, 218, 107, 220, 109, 110, 183,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 24,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 40,
		0, 0, 0, 32, 0, 32, 32, 48, 0, 32, 32, 48, 32, 48, 48, 56,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 72,
		0, 0, 0, 64, 0, 64, 64, 80, 0, 64, 64, 80, 64, 80, 80, 88,
		0, 0, 0, 64, 0, 64, 64, 96, 0, 64, 64, 96, 64, 96, 96, 104,
		0, 64, 64, 96, 64, 96, 96, 112, 64, 96, 96, 112, 96, 112, 112, 56,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 136,
		0, 0, 0, 128, 0, 128, 128, 144, 0, 128, 128, 144, 128, 144, 144, 152,
		0, 0, 0, 128, 0, 128, 128, 160, 0, 128, 128, 160, 128, 160, 160, 168,
		0, 128, 128, 160, 128, 160, 160, 176, 128, 160, 160, 176, 160, 176, 176, 56,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 200,
		0, 128, 128, 192, 128, 192, 192, 208, 128, 192, 192, 208, 192, 208, 208, 88,
		0, 128, 128, 192, 128, 192, 192, 224, 128, 192, 192, 224, 192, 224, 224, 104,
		128, 192, 192, 224, 192, 224, 224, 112, 192, 224, 224, 112, 224, 112, 112, 184,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 25,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 41,
		16, 1, 2, 33, 4, 33, 34, 49, 8, 33, 34, 49, 36, 49, 50, 57,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 73,
		16, 1, 2, 65, 4, 65, 66, 81, 8, 65, 66, 81, 68, 81, 82, 89,
		32, 1, 2, 65, 4, 65, 66, 97, 8, 65, 66, 97, 68, 97, 98, 105,
		16, 65, 66, 97, 68, 97, 98, 113, 72, 97, 98, 113, 100, 113, 114, 57,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 137,
		16, 1, 2, 129, 4, 129, 130, 145, 8, 129, 130, 145, 132, 145, 146, 153,
		32, 1, 2, 129, 4, 129, 130, 161, 8, 129, 130, 161, 132, 161, 162, 169,
		16, 129, 130, 161, 132, 161, 162, 177, 136, 161, 162, 177, 164, 177, 178, 57,
		64, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 201,
		16, 129, 130, 193, 132, 193, 194, 209, 136, 193, 194, 209, 196, 209, 210, 89,
		32, 129, 130, 193, 132, 193, 194, 225, 136, 193, 194, 225, 196, 225, 226, 105,
		144, 193, 194, 225, 196, 225, 226, 113, 200, 225, 226, 113, 228, 113, 114, 185,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 26,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 42,
		32, 16, 16, 34, 16, 36, 36, 50, 16, 40, 40, 50, 40, 52, 52, 58,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 74,
		64, 16, 16, 66, 16, 68, 68, 82, 16, 72, 72, 82, 72, 84, 84, 90,
		64, 32, 32, 66, 32, 68, 68, 98, 32, 72, 72, 98, 72, 100, 100, 106,
		32, 80, 80, 98, 80, 100, 100, 114, 80, 104, 104, 114, 104, 116, 116, 58,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 138,
		128, 16, 16, 130, 16, 132, 132, 146, 16, 136, 136, 146, 136, 148, 148, 154,
		128, 32, 32, 130, 32, 132, 132, 162, 32, 136, 136, 162, 136, 164, 164, 170,
		32, 144, 144, 162, 144, 164, 164, 178, 144, 168, 168, 178, 168, 180, 180, 58,
		128, 64, 64, 130, 64, 132, 132, 194, 64, 136, 136, 194, 136, 196, 196, 202,
		64, 144, 144, 194, 144, 196, 196, 210, 144, 200, 200, 210, 200, 212, 212, 90,
		64, 160, 160, 194, 160, 196, 196, 226, 160, 200, 200, 226, 200, 228, 228, 106,
		160, 208, 208, 226, 208, 228, 228, 114, 208, 232, 232, 114, 232, 116, 116, 186,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 27,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 43,
		48, 17, 18, 35, 20, 37, 38, 51, 24, 41, 42, 51, 44, 53, 54, 59,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 75,
		80, 17, 18, 67, 20, 69, 70, 83, 24, 73, 74, 83, 76, 85, 86, 91,
		96, 33, 34, 67, 36, 69, 70, 99, 40, 73, 74, 99, 76, 101, 102, 107,
		48, 81, 82, 99, 84, 101, 102, 115, 88, 105, 106, 115, 108, 117, 118, 59,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 139,
		144, 17, 18, 131, 20, 133, 134, 147, 24, 137, 138, 147, 140, 149, 150, 155,
		160, 33, 34, 131, 36, 133, 134, 163, 40, 137, 138, 163, 140, 165, 166, 171,
		48, 145, 146, 163, 148, 165, 166, 179, 152, 169, 170, 179, 172, 181, 182, 59,
		192, 65, 66, 131, 68, 133, 134, 195, 72, 137, 138, 195, 140, 197, 198, 203,
		80, 145, 146, 195, 148, 197, 198, 211, 152, 201, 202, 211, 204, 213, 214, 91,
		96, 161, 162, 195, 164, 197, 198, 227, 168, 201, 202, 227, 204, 229, 230, 107,
		176, 209, 210, 227, 212, 229, 230, 115, 216, 233, 234, 115, 236, 117, 118, 187,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 28,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 44,
		0, 32, 32, 48, 32, 48, 48, 52, 32, 48, 48, 56, 48, 56, 56, 60,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 76,
		0, 64, 64, 80, 64, 80, 80, 84, 64, 80, 80, 88, 80, 88, 88, 92,
		0, 64, 64, 96, 64, 96, 96, 100, 64, 96, 96, 104, 96, 104, 104, 108,
		64, 96, 96, 112, 96, 112, 112, 116, 96, 112, 112, 120, 112, 120, 120, 60,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 140,
		0, 128, 128, 144, 128, 144, 144, 148, 128, 144, 144, 152, 144, 152, 152, 156,
		0, 128, 128, 160, 128, 160, 160, 164, 128, 160, 160, 168, 160, 168, 168, 172,
		128, 160, 160, 176, 160, 176, 176, 180, 160, 176, 176, 184, 176, 184, 184, 60,
		0, 128, 128, 192, 128, 192, 192, 196, 128, 192, 192, 200, 192, 200, 200, 204,
		128, 192, 192, 208, 192, 208, 208, 212, 192, 208, 208, 216, 208, 216, 216, 92,
		128, 192, 192, 224, 192, 224, 224, 228, 192, 224, 224, 232, 224, 232, 232, 108,
		192, 224, 224, 240, 224, 240, 240, 116, 224, 240, 240, 120, 240, 120, 120, 188,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 29,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 45,
		16, 33, 34, 49, 36, 49, 50, 53, 40, 49, 50, 57, 52, 57, 58, 61,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 77,
		16, 65, 66, 81, 68, 81, 82, 85, 72, 81, 82, 89, 84, 89, 90, 93,
		32, 65, 66, 97, 68, 97, 98, 101, 72, 97, 98, 105, 100, 105, 106, 109,
		80, 97, 98, 113, 100, 113, 114, 117, 104, 113, 114, 121, 116, 121, 122, 61,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 141,
		16, 129, 130, 145, 132, 145, 146, 149, 136, 145, 146, 153, 148, 153, 154, 157,
		32, 129, 130, 161, 132, 161, 162, 165, 136, 161, 162, 169, 164, 169, 170, 173,
		144, 161, 162, 177, 164, 177, 178, 181, 168, 177, 178, 185, 180, 185, 186, 61,
		64, 129, 130, 193, 132, 193, 194, 197, 136, 193, 194, 201, 196, 201, 202, 205,
		144, 193, 194, 209, 196, 209, 210, 213, 200, 209, 210, 217, 212, 217, 218, 93,
		160, 193, 194, 225, 196, 225, 226, 229, 200, 225, 226, 233, 228, 233, 234, 109,
		208, 225, 226, 241, 228, 241, 242, 117, 232, 241, 242, 121, 244, 121, 122, 189,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 30,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 46,
		32, 48, 48, 50, 48, 52, 52, 54, 48, 56, 56, 58, 56, 60, 60, 62,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 78,
		64, 80, 80, 82, 80, 84, 84, 86, 80, 88, 88, 90, 88, 92, 92, 94,
		64, 96, 96, 98, 96, 100, 100, 102, 96, 104, 104, 106, 104, 108, 108, 110,
		96, 112, 112, 114, 112, 116, 116, 118, 112, 120, 120, 122, 120, 124, 124, 62,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 142,
		128, 144, 144, 146, 144, 148, 148, 150, 144, 152, 152, 154, 152, 156, 156, 158,
		128, 160, 160, 162, 160, 164, 164, 166, 160, 168, 168, 170, 168, 172, 172, 174,
		160, 176, 176, 178, 176, 180, 180, 182, 176, 184, 184, 186, 184, 188, 188, 62,
		128, 192, 192, 194, 192, 196, 196, 198, 192, 200, 200, 202, 200, 204, 204, 206,
		192, 208, 208, 210, 208, 212, 212, 214, 208, 216, 216, 218, 216, 220, 220, 94,
		192, 224, 224, 226, 224, 228, 228, 230, 224, 232, 232, 234, 232, 236, 236, 110,
		224, 240, 240, 242, 240, 244, 244, 118, 240, 248, 248, 122, 248, 124, 124, 190,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
		80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95,
		96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111,
		112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 63,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143,
		144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159,
		160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175,
		176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 63,
		192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207,
		208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 95,
		224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 111,
		240, 241, 242, 243, 244, 245, 246, 119, 248, 249, 250, 123, 252, 125, 126, 191,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 192,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 193,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 2,
		64, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		64, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 66,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 2,
		128, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		128, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 130,
		128, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 2,
		64, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 130,
		64, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 130,
		32, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 194,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 3,
		80, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		96, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 67,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 3,
		144, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		160, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 131,
		192, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 3,
		80, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 131,
		96, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 131,
		48, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 195,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 4,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 4,
		0, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 4,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 4,
		0, 64, 64, 16, 64, 16, 16, 4, 64, 16, 16, 8, 16, 8, 8, 4,
		0, 64, 64, 32, 64, 32, 32, 4, 64, 32, 32, 8, 32, 8, 8, 4,
		64, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 68,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 4,
		0, 128, 128, 16, 128, 16, 16, 4, 128, 16, 16, 8, 16, 8, 8, 4,
		0, 128, 128, 32, 128, 32, 32, 4, 128, 32, 32, 8, 32, 8, 8, 4,
		128, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 132,
		0, 128, 128, 64, 128, 64, 64, 4, 128, 64, 64, 8, 64, 8, 8, 4,
		128, 64, 64, 16, 64, 16, 16, 4, 64, 16, 16, 8, 16, 8, 8, 132,
		128, 64, 64, 32, 64, 32, 32, 4, 64, 32, 32, 8, 32, 8, 8, 132,
		64, 32, 32, 16, 32, 16, 16, 132, 32, 16, 16, 136, 16, 136, 136, 196,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 5,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 5,
		16, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 5,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 5,
		16, 65, 66, 17, 68, 17, 18, 5, 72, 17, 18, 9, 20, 9, 10, 5,
		32, 65, 66, 33, 68, 33, 34, 5, 72, 33, 34, 9, 36, 9, 10, 5,
		80, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 69,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 5,
		16, 129, 130, 17, 132, 17, 18, 5, 136, 17, 18, 9, 20, 9, 10, 5,
		32, 129, 130, 33, 132, 33, 34, 5, 136, 33, 34, 9, 36, 9, 10, 5,
		144, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 133,
		64, 129, 130, 65, 132, 65, 66, 5, 136, 65, 66, 9, 68, 9, 10, 5,
		144, 65, 66, 17, 68, 17, 18, 5, 72, 17, 18, 9, 20, 9, 10, 133,
		160, 65, 66, 33, 68, 33, 34, 5, 72, 33, 34, 9, 36, 9, 10, 133,
		80, 33, 34, 17, 36, 17, 18, 133, 40, 17, 18, 137, 20, 137, 138, 197,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 6,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 6,
		32, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 6,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 6,
		64, 80, 80, 18, 80, 20, 20, 6, 80, 24, 24, 10, 24, 12, 12, 6,
		64, 96, 96, 34, 96, 36, 36, 6, 96, 40, 40, 10, 40, 12, 12, 6,
		96, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 70,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 6,
		128, 144, 144, 18, 144, 20, 20, 6, 144, 24, 24, 10, 24, 12, 12, 6,
		128, 160, 160, 34, 160, 36, 36, 6, 160, 40, 40, 10, 40, 12, 12, 6,
		160, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 134,
		128, 192, 192, 66, 192, 68, 68, 6, 192, 72, 72, 10, 72, 12, 12, 6,
		192, 80, 80, 18, 80, 20, 20, 6, 80, 24, 24, 10, 24, 12, 12, 134,
		192, 96, 96, 34, 96, 36, 36, 6, 96, 40, 40, 10, 40, 12, 12, 134,
		96, 48, 48, 18, 48, 20, 20, 134, 48, 24, 24, 138, 24, 140, 140, 198,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 7,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 7,
		48, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 7,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 7,
		80, 81, 82, 19, 84, 21, 22, 7, 88, 25, 26, 11, 28, 13, 14, 7,
		96, 97, 98, 35, 100, 37, 38, 7, 104, 41, 42, 11, 44, 13, 14, 7,
		112, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 71,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 7,
		144, 145, 146, 19, 148, 21, 22, 7, 152, 25, 26, 11, 28, 13, 14, 7,
		160, 161, 162, 35, 164, 37, 38, 7, 168, 41, 42, 11, 44, 13, 14, 7,
		176, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 135,
		192, 193, 194, 67, 196, 69, 70, 7, 200, 73, 74, 11, 76, 13, 14, 7,
		208, 81, 82, 19, 84, 21, 22, 7, 88, 25, 26, 11, 28, 13, 14, 135,
		224, 97, 98, 35, 100, 37, 38, 7, 104, 41, 42, 11, 44, 13, 14, 135,
		112, 49, 50, 19, 52, 21, 22, 135, 56, 25, 26, 139, 28, 141, 142, 199,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 8,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 8,
		0, 0, 0, 32, 0, 32, 32, 16, 0, 32, 32, 16, 32, 16, 16, 8,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 8,
		0, 0, 0, 64, 0, 64, 64, 16, 0, 64, 64, 16, 64, 16, 16, 8,
		0, 0, 0, 64, 0, 64, 64, 32, 0, 64, 64, 32, 64, 32, 32, 8,
		0, 64, 64, 32, 64, 32, 32, 16, 64, 32, 32, 16, 32, 16, 16, 72,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 8,
		0, 0, 0, 128, 0, 128, 128, 16, 0, 128, 128, 16, 128, 16, 16, 8,
		0, 0, 0, 128, 0, 128, 128, 32, 0, 128, 128, 32, 128, 32, 32, 8,
		0, 128, 128, 32, 128, 32, 32, 16, 128, 32, 32, 16, 32, 16, 16, 136,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 8,
		0, 128, 128, 64, 128, 64, 64, 16, 128, 64, 64, 16, 64, 16, 16, 136,
		0, 128, 128, 64, 128, 64, 64, 32, 128, 64, 64, 32, 64, 32, 32, 136,
		128, 64, 64, 32, 64, 32, 32, 144, 64, 32, 32, 144, 32, 144, 144, 200,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 9,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 9,
		16, 1, 2, 33, 4, 33, 34, 17, 8, 33, 34, 17, 36, 17, 18, 9,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 9,
		16, 1, 2, 65, 4, 65, 66, 17, 8, 65, 66, 17, 68, 17, 18, 9,
		32, 1, 2, 65, 4, 65, 66, 33, 8, 65, 66, 33, 68, 33, 34, 9,
		16, 65, 66, 33, 68, 33, 34, 17, 72, 33, 34, 17, 36, 17, 18, 73,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 9,
		16, 1, 2, 129, 4, 129, 130, 17, 8, 129, 130, 17, 132, 17, 18, 9,
		32, 1, 2, 129, 4, 129, 130, 33, 8, 129, 130, 33, 132, 33, 34, 9,
		16, 129, 130, 33, 132, 33, 34, 17, 136, 33, 34, 17, 36, 17, 18, 137,
		64, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 9,
		16, 129, 130, 65, 132, 65, 66, 17, 136, 65, 66, 17, 68, 17, 18, 137,
		32, 129, 130, 65, 132, 65, 66, 33, 136, 65, 66, 33, 68, 33, 34, 137,
		144, 65, 66, 33, 68, 33, 34, 145, 72, 33, 34, 145, 36, 145, 146, 201,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 10,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 10,
		32, 16, 16, 34, 16, 36, 36, 18, 16, 40, 40, 18, 40, 20, 20, 10,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 10,
		64, 16, 16, 66, 16, 68, 68, 18, 16, 72, 72, 18, 72, 20, 20, 10,
		64, 32, 32, 66, 32, 68, 68, 34, 32, 72, 72, 34, 72, 36, 36, 10,
		32, 80, 80, 34, 80, 36, 36, 18, 80, 40, 40, 18, 40, 20, 20, 74,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 10,
		128, 16, 16, 130, 16, 132, 132, 18, 16, 136, 136, 18, 136, 20, 20, 10,
		128, 32, 32, 130, 32, 132, 132, 34, 32, 136, 136, 34, 136, 36, 36, 10,
		32, 144, 144, 34, 144, 36, 36, 18, 144, 40, 40, 18, 40, 20, 20, 138,
		128, 64, 64, 130, 64, 132, 132, 66, 64, 136, 136, 66, 136, 68, 68, 10,
		64, 144, 144, 66, 144, 68, 68, 18, 144, 72, 72, 18, 72, 20, 20, 138,
		64, 160, 160, 66, 160, 68, 68, 34, 160, 72, 72, 34, 72, 36, 36, 138,
		160, 80, 80, 34, 80, 36, 36, 146, 80, 40, 40, 146, 40, 148, 148, 202,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 11,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 11,
		48, 17, 18, 35, 20, 37, 38, 19, 24, 41, 42, 19, 44, 21, 22, 11,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 11,
		80, 17, 18, 67, 20, 69, 70, 19, 24, 73, 74, 19, 76, 21, 22, 11,
		96, 33, 34, 67, 36, 69, 70, 35, 40, 73, 74, 35, 76, 37, 38, 11,
		48, 81, 82, 35, 84, 37, 38, 19, 88, 41, 42, 19, 44, 21, 22, 75,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 11,
		144, 17, 18, 131, 20, 133, 134, 19, 24, 137, 138, 19, 140, 21, 22, 11,
		160, 33, 34, 131, 36, 133, 134, 35, 40, 137, 138, 35, 140, 37, 38, 11,
		48, 145, 146, 35, 148, 37, 38, 19, 152, 41, 42, 19, 44, 21, 22, 139,
		192, 65, 66, 131, 68, 133, 134, 67, 72, 137, 138, 67, 140, 69, 70, 11,
		80, 145, 146, 67, 148, 69, 70, 19, 152, 73, 74, 19, 76, 21, 22, 139,
		96, 161, 162, 67, 164, 69, 70, 35, 168, 73, 74, 35, 76, 37, 38, 139,
		176, 81, 82, 35, 84, 37, 38, 147, 88, 41, 42, 147, 44, 149, 150, 203,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 12,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 12,
		0, 32, 32, 48, 32, 48, 48, 20, 32, 48, 48, 24, 48, 24, 24, 12,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 12,
		0, 64, 64, 80, 64, 80, 80, 20, 64, 80, 80, 24, 80, 24, 24, 12,
		0, 64, 64, 96, 64, 96, 96, 36, 64, 96, 96, 40, 96, 40, 40, 12,
		64, 96, 96, 48, 96, 48, 48, 20, 96, 48, 48, 24, 48, 24, 24, 76,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 12,
		0, 128, 128, 144, 128, 144, 144, 20, 128, 144, 144, 24, 144, 24, 24, 12,
		0, 128, 128, 160, 128, 160, 160, 36, 128, 160, 160, 40, 160, 40, 40, 12,
		128, 160, 160, 48, 160, 48, 48, 20, 160, 48, 48, 24, 48, 24, 24, 140,
		0, 128, 128, 192, 128, 192, 192, 68, 128, 192, 192, 72, 192, 72, 72, 12,
		128, 192, 192, 80, 192, 80, 80, 20, 192, 80, 80, 24, 80, 24, 24, 140,
		128, 192, 192, 96, 192, 96, 96, 36, 192, 96, 96, 40, 96, 40, 40, 140,
		192, 96, 96, 48, 96, 48, 48, 148, 96, 48, 48, 152, 48, 152, 152, 204,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 13,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 13,
		16, 33, 34, 49, 36, 49, 50, 21, 40, 49, 50, 25, 52, 25, 26, 13,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 13,
		16, 65, 66, 81, 68, 81, 82, 21, 72, 81, 82, 25, 84, 25, 26, 13,
		32, 65, 66, 97, 68, 97, 98, 37, 72, 97, 98, 41, 100, 41, 42, 13,
		80, 97, 98, 49, 100, 49, 50, 21, 104, 49, 50, 25, 52, 25, 26, 77,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 13,
		16, 129, 130, 145, 132, 145, 146, 21, 136, 145, 146, 25, 148, 25, 26, 13,
		32, 129, 130, 161, 132, 161, 162, 37, 136, 161, 162, 41, 164, 41, 42, 13,
		144, 161, 162, 49, 164, 49, 50, 21, 168, 49, 50, 25, 52, 25, 26, 141,
		64, 129, 130, 193, 132, 193, 194, 69, 136, 193, 194, 73, 196, 73, 74, 13,
		144, 193, 194, 81, 196, 81, 82, 21, 200, 81, 82, 25, 84, 25, 26, 141,
		160, 193, 194, 97, 196, 97, 98, 37, 200, 97, 98, 41, 100, 41, 42, 141,
		208, 97, 98, 49, 100, 49, 50, 149, 104, 49, 50, 153, 52, 153, 154, 205,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 14,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 14,
		32, 48, 48, 50, 48, 52, 52, 22, 48, 56, 56, 26, 56, 28, 28, 14,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 14,
		64, 80, 80, 82, 80, 84, 84, 22, 80, 88, 88, 26, 88, 28, 28, 14,
		64, 96, 96, 98, 96, 100, 100, 38, 96, 104, 104, 42, 104, 44, 44, 14,
		96, 112, 112, 50, 112, 52, 52, 22, 112, 56, 56, 26, 56, 28, 28, 78,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 14,
		128, 144, 144, 146, 144, 148, 148, 22, 144, 152, 152, 26, 152, 28, 28, 14,
		128, 160, 160, 162, 160, 164, 164, 38, 160, 168, 168, 42, 168, 44, 44, 14,
		160, 176, 176, 50, 176, 52, 52, 22, 176, 56, 56, 26, 56, 28, 28, 142,
		128, 192, 192, 194, 192, 196, 196, 70, 192, 200, 200, 74, 200, 76, 76, 14,
		192, 208, 208, 82, 208, 84, 84, 22, 208, 88, 88, 26, 88, 28, 28, 142,
		192, 224, 224, 98, 224, 100, 100, 38, 224, 104, 104, 42, 104, 44, 44, 142,
		224, 112, 112, 50, 112, 52, 52, 150, 112, 56, 56, 154, 56, 156, 156, 206,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 15,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 15,
		48, 49, 50, 51, 52, 53, 54, 23, 56, 57, 58, 27, 60, 29, 30, 15,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 15,
		80, 81, 82, 83, 84, 85, 86, 23, 88, 89, 90, 27, 92, 29, 30, 15,
		96, 97, 98, 99, 100, 101, 102, 39, 104, 105, 106, 43, 108, 45, 46, 15,
		112, 113, 114, 51, 116, 53, 54, 23, 120, 57, 58, 27, 60, 29, 30, 79,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 15,
		144, 145, 146, 147, 148, 149, 150, 23, 152, 153, 154, 27, 156, 29, 30, 15,
		160, 161, 162, 163, 164, 165, 166, 39, 168, 169, 170, 43, 172, 45, 46, 15,
		176, 177, 178, 51, 180, 53, 54, 23, 184, 57, 58, 27, 60, 29, 30, 143,
		192, 193, 194, 195, 196, 197, 198, 71, 200, 201, 202, 75, 204, 77, 78, 15,
		208, 209, 210, 83, 212, 85, 86, 23, 216, 89, 90, 27, 92, 29, 30, 143,
		224, 225, 226, 99, 228, 101, 102, 39, 232, 105, 106, 43, 108, 45, 46, 143,
		240, 113, 114, 51, 116, 53, 54, 151, 120, 57, 58, 155, 60, 157, 158, 207,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 16,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 32,
		0, 0, 0, 64, 0, 64, 64, 32, 0, 64, 64, 32, 64, 32, 32, 80,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 16,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 32,
		0, 0, 0, 128, 0, 128, 128, 32, 0, 128, 128, 32, 128, 32, 32, 144,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 64,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 144,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 160,
		0, 128, 128, 64, 128, 64, 64, 160, 128, 64, 64, 160, 64, 160, 160, 208,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 17,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 33,
		16, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 17,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		16, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 17,
		32, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 33,
		16, 1, 2, 65, 4, 65, 66, 33, 8, 65, 66, 33, 68, 33, 34, 81,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 17,
		32, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 33,
		16, 1, 2, 129, 4, 129, 130, 33, 8, 129, 130, 33, 132, 33, 34, 145,
		64, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 65,
		16, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 145,
		32, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 161,
		16, 129, 130, 65, 132, 65, 66, 161, 136, 65, 66, 161, 68, 161, 162, 209,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 18,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 34,
		32, 16, 16, 2, 16, 4, 4, 34, 16, 8, 8, 34, 8, 36, 36, 18,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 66,
		64, 16, 16, 2, 16, 4, 4, 66, 16, 8, 8, 66, 8, 68, 68, 18,
		64, 32, 32, 2, 32, 4, 4, 66, 32, 8, 8, 66, 8, 68, 68, 34,
		32, 16, 16, 66, 16, 68, 68, 34, 16, 72, 72, 34, 72, 36, 36, 82,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 130,
		128, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 18,
		128, 32, 32, 2, 32, 4, 4, 130, 32, 8, 8, 130, 8, 132, 132, 34,
		32, 16, 16, 130, 16, 132, 132, 34, 16, 136, 136, 34, 136, 36, 36, 146,
		128, 64, 64, 2, 64, 4, 4, 130, 64, 8, 8, 130, 8, 132, 132, 66,
		64, 16, 16, 130, 16, 132, 132, 66, 16, 136, 136, 66, 136, 68, 68, 146,
		64, 32, 32, 130, 32, 132, 132, 66, 32, 136, 136, 66, 136, 68, 68, 162,
		32, 144, 144, 66, 144, 68, 68, 162, 144, 72, 72, 162, 72, 164, 164, 210,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 19,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 35,
		48, 17, 18, 3, 20, 5, 6, 35, 24, 9, 10, 35, 12, 37, 38, 19,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 67,
		80, 17, 18, 3, 20, 5, 6, 67, 24, 9, 10, 67, 12, 69, 70, 19,
		96, 33, 34, 3, 36, 5, 6, 67, 40, 9, 10, 67, 12, 69, 70, 35,
		48, 17, 18, 67, 20, 69, 70, 35, 24, 73, 74, 35, 76, 37, 38, 83,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 131,
		144, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 19,
		160, 33, 34, 3, 36, 5, 6, 131, 40, 9, 10, 131, 12, 133, 134, 35,
		48, 17, 18, 131, 20, 133, 134, 35, 24, 137, 138, 35, 140, 37, 38, 147,
		192, 65, 66, 3, 68, 5, 6, 131, 72, 9, 10, 131, 12, 133, 134, 67,
		80, 17, 18, 131, 20, 133, 134, 67, 24, 137, 138, 67, 140, 69, 70, 147,
		96, 33, 34, 131, 36, 133, 134, 67, 40, 137, 138, 67, 140, 69, 70, 163,
		48, 145, 146, 67, 148, 69, 70, 163, 152, 73, 74, 163, 76, 165, 166, 211,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 20,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 36,
		0, 32, 32, 16, 32, 16, 16, 36, 32, 16, 16, 40, 16, 40, 40, 20,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 68,
		0, 64, 64, 16, 64, 16, 16, 68, 64, 16, 16, 72, 16, 72, 72, 20,
		0, 64, 64, 32, 64, 32, 32, 68, 64, 32, 32, 72, 32, 72, 72, 36,
		64, 32, 32, 80, 32, 80, 80, 36, 32, 80, 80, 40, 80, 40, 40, 84,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 132,
		0, 128, 128, 16, 128, 16, 16, 132, 128, 16, 16, 136, 16, 136, 136, 20,
		0, 128, 128, 32, 128, 32, 32, 132, 128, 32, 32, 136, 32, 136, 136, 36,
		128, 32, 32, 144, 32, 144, 144, 36, 32, 144, 144, 40, 144, 40, 40, 148,
		0, 128, 128, 64, 128, 64, 64, 132, 128, 64, 64, 136, 64, 136, 136, 68,
		128, 64, 64, 144, 64, 144, 144, 68, 64, 144, 144, 72, 144, 72, 72, 148,
		128, 64, 64, 160, 64, 160, 160, 68, 64, 160, 160, 72, 160, 72, 72, 164,
		64, 160, 160, 80, 160, 80, 80, 164, 160, 80, 80, 168, 80, 168, 168, 212,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 21,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 37,
		16, 33, 34, 17, 36, 17, 18, 37, 40, 17, 18, 41, 20, 41, 42, 21,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 69,
		16, 65, 66, 17, 68, 17, 18, 69, 72, 17, 18, 73, 20, 73, 74, 21,
		32, 65, 66, 33, 68, 33, 34, 69, 72, 33, 34, 73, 36, 73, 74, 37,
		80, 33, 34, 81, 36, 81, 82, 37, 40, 81, 82, 41, 84, 41, 42, 85,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 133,
		16, 129, 130, 17, 132, 17, 18, 133, 136, 17, 18, 137, 20, 137, 138, 21,
		32, 129, 130, 33, 132, 33, 34, 133, 136, 33, 34, 137, 36, 137, 138, 37,
		144, 33, 34, 145, 36, 145, 146, 37, 40, 145, 146, 41, 148, 41, 42, 149,
		64, 129, 130, 65, 132, 65, 66, 133, 136, 65, 66, 137, 68, 137, 138, 69,
		144, 65, 66, 145, 68, 145, 146, 69, 72, 145, 146, 73, 148, 73, 74, 149,
		160, 65, 66, 161, 68, 161, 162, 69, 72, 161, 162, 73, 164, 73, 74, 165,
		80, 161, 162, 81, 164, 81, 82, 165, 168, 81, 82, 169, 84, 169, 170, 213,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 22,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 38,
		32, 48, 48, 18, 48, 20, 20, 38, 48, 24, 24, 42, 24, 44, 44, 22,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 70,
		64, 80, 80, 18, 80, 20, 20, 70, 80, 24, 24, 74, 24, 76, 76, 22,
		64, 96, 96, 34, 96, 36, 36, 70, 96, 40, 40, 74, 40, 76, 76, 38,
		96, 48, 48, 82, 48, 84, 84, 38, 48, 88, 88, 42, 88, 44, 44, 86,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 134,
		128, 144, 144, 18, 144, 20, 20, 134, 144, 24, 24, 138, 24, 140, 140, 22,
		128, 160, 160, 34, 160, 36, 36, 134, 160, 40, 40, 138, 40, 140, 140, 38,
		160, 48, 48, 146, 48, 148, 148, 38, 48, 152, 152, 42, 152, 44, 44, 150,
		128, 192, 192, 66, 192, 68, 68, 134, 192, 72, 72, 138, 72, 140, 140, 70,
		192, 80, 80, 146, 80, 148, 148, 70, 80, 152, 152, 74, 152, 76, 76, 150,
		192, 96, 96, 162, 96, 164, 164, 70, 96, 168, 168, 74, 168, 76, 76, 166,
		96, 176, 176, 82, 176, 84, 84, 166, 176, 88, 88, 170, 88, 172, 172, 214,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 23,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 39,
		48, 49, 50, 19, 52, 21, 22, 39, 56, 25, 26, 43, 28, 45, 46, 23,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 71,
		80, 81, 82, 19, 84, 21, 22, 71, 88, 25, 26, 75, 28, 77, 78, 23,
		96, 97, 98, 35, 100, 37, 38, 71, 104, 41, 42, 75, 44, 77, 78, 39,
		112, 49, 50, 83, 52, 85, 86, 39, 56, 89, 90, 43, 92, 45, 46, 87,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 135,
		144, 145, 146, 19, 148, 21, 22, 135, 152, 25, 26, 139, 28, 141, 142, 23,
		160, 161, 162, 35, 164, 37, 38, 135, 168, 41, 42, 139, 44, 141, 142, 39,
		176, 49, 50, 147, 52, 149, 150, 39, 56, 153, 154, 43, 156, 45, 46, 151,
		192, 193, 194, 67, 196, 69, 70, 135, 200, 73, 74, 139, 76, 141, 142, 71,
		208, 81, 82, 147, 84, 149, 150, 71, 88, 153, 154, 75, 156, 77, 78, 151,
		224, 97, 98, 163, 100, 165, 166, 71, 104, 169, 170, 75, 172, 77, 78, 167,
		112, 177, 178, 83, 180, 85, 86, 167, 184, 89, 90, 171, 92, 173, 174, 215,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 24,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 40,
		0, 0, 0, 32, 0, 32, 32, 48, 0, 32, 32, 48, 32, 48, 48, 24,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 72,
		0, 0, 0, 64, 0, 64, 64, 80, 0, 64, 64, 80, 64, 80, 80, 24,
		0, 0, 0, 64, 0, 64, 64, 96, 0, 64, 64, 96, 64, 96, 96, 40,
		0, 64, 64, 96, 64, 96, 96, 48, 64, 96, 96, 48, 96, 48, 48, 88,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 136,
		0, 0, 0, 128, 0, 128, 128, 144, 0, 128, 128, 144, 128, 144, 144, 24,
		0, 0, 0, 128, 0, 128, 128, 160, 0, 128, 128, 160, 128, 160, 160, 40,
		0, 128, 128, 160, 128, 160, 160, 48, 128, 160, 160, 48, 160, 48, 48, 152,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 72,
		0, 128, 128, 192, 128, 192, 192, 80, 128, 192, 192, 80, 192, 80, 80, 152,
		0, 128, 128, 192, 128, 192, 192, 96, 128, 192, 192, 96, 192, 96, 96, 168,
		128, 192, 192, 96, 192, 96, 96, 176, 192, 96, 96, 176, 96, 176, 176, 216,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 25,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 41,
		16, 1, 2, 33, 4, 33, 34, 49, 8, 33, 34, 49, 36, 49, 50, 25,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 73,
		16, 1, 2, 65, 4, 65, 66, 81, 8, 65, 66, 81, 68, 81, 82, 25,
		32, 1, 2, 65, 4, 65, 66, 97, 8, 65, 66, 97, 68, 97, 98, 41,
		16, 65, 66, 97, 68, 97, 98, 49, 72, 97, 98, 49, 100, 49, 50, 89,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 137,
		16, 1, 2, 129, 4, 129, 130, 145, 8, 129, 130, 145, 132, 145, 146, 25,
		32, 1, 2, 129, 4, 129, 130, 161, 8, 129, 130, 161, 132, 161, 162, 41,
		16, 129, 130, 161, 132, 161, 162, 49, 136, 161, 162, 49, 164, 49, 50, 153,
		64, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 73,
		16, 129, 130, 193, 132, 193, 194, 81, 136, 193, 194, 81, 196, 81, 82, 153,
		32, 129, 130, 193, 132, 193, 194, 97, 136, 193, 194, 97, 196, 97, 98, 169,
		144, 193, 194, 97, 196, 97, 98, 177, 200, 97, 98, 177, 100, 177, 178, 217,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 26,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 42,
		32, 16, 16, 34, 16, 36, 36, 50, 16, 40, 40, 50, 40, 52, 52, 26,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 74,
		64, 16, 16, 66, 16, 68, 68, 82, 16, 72, 72, 82, 72, 84, 84, 26,
		64, 32, 32, 66, 32, 68, 68, 98, 32, 72, 72, 98, 72, 100, 100, 42,
		32, 80, 80, 98, 80, 100, 100, 50, 80, 104, 104, 50, 104, 52, 52, 90,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 138,
		128, 16, 16, 130, 16, 132, 132, 146, 16, 136, 136, 146, 136, 148, 148, 26,
		128, 32, 32, 130, 32, 132, 132, 162, 32, 136, 136, 162, 136, 164, 164, 42,
		32, 144, 144, 162, 144, 164, 164, 50, 144, 168, 168, 50, 168, 52, 52, 154,
		128, 64, 64, 130, 64, 132, 132, 194, 64, 136, 136, 194, 136, 196, 196, 74,
		64, 144, 144, 194, 144, 196, 196, 82, 144, 200, 200, 82, 200, 84, 84, 154,
		64, 160, 160, 194, 160, 196, 196, 98, 160, 200, 200, 98, 200, 100, 100, 170,
		160, 208, 208, 98, 208, 100, 100, 178, 208, 104, 104, 178, 104, 180, 180, 218,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 27,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 43,
		48, 17, 18, 35, 20, 37, 38, 51, 24, 41, 42, 51, 44, 53, 54, 27,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 75,
		80, 17, 18, 67, 20, 69, 70, 83, 24, 73, 74, 83, 76, 85, 86, 27,
		96, 33, 34, 67, 36, 69, 70, 99, 40, 73, 74, 99, 76, 101, 102, 43,
		48, 81, 82, 99, 84, 101, 102, 51, 88, 105, 106, 51, 108, 53, 54, 91,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 139,
		144, 17, 18, 131, 20, 133, 134, 147, 24, 137, 138, 147, 140, 149, 150, 27,
		160, 33, 34, 131, 36, 133, 134, 163, 40, 137, 138, 163, 140, 165, 166, 43,
		48, 145, 146, 163, 148, 165, 166, 51, 152, 169, 170, 51, 172, 53, 54, 155,
		192, 65, 66, 131, 68, 133, 134, 195, 72, 137, 138, 195, 140, 197, 198, 75,
		80, 145, 146, 195, 148, 197, 198, 83, 152, 201, 202, 83, 204, 85, 86, 155,
		96, 161, 162, 195, 164, 197, 198, 99, 168, 201, 202, 99, 204, 101, 102, 171,
		176, 209, 210, 99, 212, 101, 102, 179, 216, 105, 106, 179, 108, 181, 182, 219,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 28,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 44,
		0, 32, 32, 48, 32, 48, 48, 52, 32, 48, 48, 56, 48, 56, 56, 28,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 76,
		0, 64, 64, 80, 64, 80, 80, 84, 64, 80, 80, 88, 80, 88, 88, 28,
		0, 64, 64, 96, 64, 96, 96, 100, 64, 96, 96, 104, 96, 104, 104, 44,
		64, 96, 96, 112, 96, 112, 112, 52, 96, 112, 112, 56, 112, 56, 56, 92,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 140,
		0, 128, 128, 144, 128, 144, 144, 148, 128, 144, 144, 152, 144, 152, 152, 28,
		0, 128, 128, 160, 128, 160, 160, 164, 128, 160, 160, 168, 160, 168, 168, 44,
		128, 160, 160, 176, 160, 176, 176, 52, 160, 176, 176, 56, 176, 56, 56, 156,
		0, 128, 128, 192, 128, 192, 192, 196, 128, 192, 192, 200, 192, 200, 200, 76,
		128, 192, 192, 208, 192, 208, 208, 84, 192, 208, 208, 88, 208, 88, 88, 156,
		128, 192, 192, 224, 192, 224, 224, 100, 192, 224, 224, 104, 224, 104, 104, 172,
		192, 224, 224, 112, 224, 112, 112, 180, 224, 112, 112, 184, 112, 184, 184, 220,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 29,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 45,
		16, 33, 34, 49, 36, 49, 50, 53, 40, 49, 50, 57, 52, 57, 58, 29,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 77,
		16, 65, 66, 81, 68, 81, 82, 85, 72, 81, 82, 89, 84, 89, 90, 29,
		32, 65, 66, 97, 68, 97, 98, 101, 72, 97, 98, 105, 100, 105, 106, 45,
		80, 97, 98, 113, 100, 113, 114, 53, 104, 113, 114, 57, 116, 57, 58, 93,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 141,
		16, 129, 130, 145, 132, 145, 146, 149, 136, 145, 146, 153, 148, 153, 154, 29,
		32, 129, 130, 161, 132, 161, 162, 165, 136, 161, 162, 169, 164, 169, 170, 45,
		144, 161, 162, 177, 164, 177, 178, 53, 168, 177, 178, 57, 180, 57, 58, 157,
		64, 129, 130, 193, 132, 193, 194, 197, 136, 193, 194, 201, 196, 201, 202, 77,
		144, 193, 194, 209, 196, 209, 210, 85, 200, 209, 210, 89, 212, 89, 90, 157,
		160, 193, 194, 225, 196, 225, 226, 101, 200, 225, 226, 105, 228, 105, 106, 173,
		208, 225, 226, 113, 228, 113, 114, 181, 232, 113, 114, 185, 116, 185, 186, 221,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 30,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 46,
		32, 48, 48, 50, 48, 52, 52, 54, 48, 56, 56, 58, 56, 60, 60, 30,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 78,
		64, 80, 80, 82, 80, 84, 84, 86, 80, 88, 88, 90, 88, 92, 92, 30,
		64, 96, 96, 98, 96, 100, 100, 102, 96, 104, 104, 106, 104, 108, 108, 46,
		96, 112, 112, 114, 112, 116, 116, 54, 112, 120, 120, 58, 120, 60, 60, 94,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 142,
		128, 144, 144, 146, 144, 148, 148, 150, 144, 152, 152, 154, 152, 156, 156, 30,
		128, 160, 160, 162, 160, 164, 164, 166, 160, 168, 168, 170, 168, 172, 172, 46,
		160, 176, 176, 178, 176, 180, 180, 54, 176, 184, 184, 58, 184, 60, 60, 158,
		128, 192, 192, 194, 192, 196, 196, 198, 192, 200, 200, 202, 200, 204, 204, 78,
		192, 208, 208, 210, 208, 212, 212, 86, 208, 216, 216, 90, 216, 92, 92, 158,
		192, 224, 224, 226, 224, 228, 228, 102, 224, 232, 232, 106, 232, 108, 108, 174,
		224, 240, 240, 114, 240, 116, 116, 182, 240, 120, 120, 186, 120, 188, 188, 222,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 31,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
		80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 31,
		96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 47,
		112, 113, 114, 115, 116, 117, 118, 55, 120, 121, 122, 59, 124, 61, 62, 95,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143,
		144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 31,
		160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 47,
		176, 177, 178, 179, 180, 181, 182, 55, 184, 185, 186, 59, 188, 61, 62, 159,
		192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 79,
		208, 209, 210, 211, 212, 213, 214, 87, 216, 217, 218, 91, 220, 93, 94, 159,
		224, 225, 226, 227, 228, 229, 230, 103, 232, 233, 234, 107, 236, 109, 110, 175,
		240, 241, 242, 115, 244, 117, 118, 183, 248, 121, 122, 187, 124, 189, 190, 223,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 96,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 160,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 192,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 192,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 224,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 33,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		16, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 97,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 161,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 193,
		32, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 193,
		16, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 225,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 2,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 2,
		32, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 34,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 2,
		64, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 66,
		64, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 66,
		32, 16, 16, 2, 16, 4, 4, 66, 16, 8, 8, 66, 8, 68, 68, 98,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 2,
		128, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 130,
		128, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 130,
		32, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 162,
		128, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 130,
		64, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 194,
		64, 32, 32, 2, 32, 4, 4, 130, 32, 8, 8, 130, 8, 132, 132, 194,
		32, 16, 16, 130, 16, 132, 132, 194, 16, 136, 136, 194, 136, 196, 196, 226,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 3,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 3,
		48, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 35,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 3,
		80, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 67,
		96, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 67,
		48, 17, 18, 3, 20, 5, 6, 67, 24, 9, 10, 67, 12, 69, 70, 99,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 3,
		144, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 131,
		160, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 131,
		48, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 163,
		192, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 131,
		80, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 195,
		96, 33, 34, 3, 36, 5, 6, 131, 40, 9, 10, 131, 12, 133, 134, 195,
		48, 17, 18, 131, 20, 133, 134, 195, 24, 137, 138, 195, 140, 197, 198, 227,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 4,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 4,
		0, 32, 32, 16, 32, 16, 16, 4, 32, 16, 16, 8, 16, 8, 8, 36,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 4,
		0, 64, 64, 16, 64, 16, 16, 4, 64, 16, 16, 8, 16, 8, 8, 68,
		0, 64, 64, 32, 64, 32, 32, 4, 64, 32, 32, 8, 32, 8, 8, 68,
		64, 32, 32, 16, 32, 16, 16, 68, 32, 16, 16, 72, 16, 72, 72, 100,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 4,
		0, 128, 128, 16, 128, 16, 16, 4, 128, 16, 16, 8, 16, 8, 8, 132,
		0, 128, 128, 32, 128, 32, 32, 4, 128, 32, 32, 8, 32, 8, 8, 132,
		128, 32, 32, 16, 32, 16, 16, 132, 32, 16, 16, 136, 16, 136, 136, 164,
		0, 128, 128, 64, 128, 64, 64, 4, 128, 64, 64, 8, 64, 8, 8, 132,
		128, 64, 64, 16, 64, 16, 16, 132, 64, 16, 16, 136, 16, 136, 136, 196,
		128, 64, 64, 32, 64, 32, 32, 132, 64, 32, 32, 136, 32, 136, 136, 196,
		64, 32, 32, 144, 32, 144, 144, 196, 32, 144, 144, 200, 144, 200, 200, 228,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 5,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 5,
		16, 33, 34, 17, 36, 17, 18, 5, 40, 17, 18, 9, 20, 9, 10, 37,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 5,
		16, 65, 66, 17, 68, 17, 18, 5, 72, 17, 18, 9, 20, 9, 10, 69,
		32, 65, 66, 33, 68, 33, 34, 5, 72, 33, 34, 9, 36, 9, 10, 69,
		80, 33, 34, 17, 36, 17, 18, 69, 40, 17, 18, 73, 20, 73, 74, 101,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 5,
		16, 129, 130, 17, 132, 17, 18, 5, 136, 17, 18, 9, 20, 9, 10, 133,
		32, 129, 130, 33, 132, 33, 34, 5, 136, 33, 34, 9, 36, 9, 10, 133,
		144, 33, 34, 17, 36, 17, 18, 133, 40, 17, 18, 137, 20, 137, 138, 165,
		64, 129, 130, 65, 132, 65, 66, 5, 136, 65, 66, 9, 68, 9, 10, 133,
		144, 65, 66, 17, 68, 17, 18, 133, 72, 17, 18, 137, 20, 137, 138, 197,
		160, 65, 66, 33, 68, 33, 34, 133, 72, 33, 34, 137, 36, 137, 138, 197,
		80, 33, 34, 145, 36, 145, 146, 197, 40, 145, 146, 201, 148, 201, 202, 229,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 6,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 6,
		32, 48, 48, 18, 48, 20, 20, 6, 48, 24, 24, 10, 24, 12, 12, 38,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 6,
		64, 80, 80, 18, 80, 20, 20, 6, 80, 24, 24, 10, 24, 12, 12, 70,
		64, 96, 96, 34, 96, 36, 36, 6, 96, 40, 40, 10, 40, 12, 12, 70,
		96, 48, 48, 18, 48, 20, 20, 70, 48, 24, 24, 74, 24, 76, 76, 102,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 6,
		128, 144, 144, 18, 144, 20, 20, 6, 144, 24, 24, 10, 24, 12, 12, 134,
		128, 160, 160, 34, 160, 36, 36, 6, 160, 40, 40, 10, 40, 12, 12, 134,
		160, 48, 48, 18, 48, 20, 20, 134, 48, 24, 24, 138, 24, 140, 140, 166,
		128, 192, 192, 66, 192, 68, 68, 6, 192, 72, 72, 10, 72, 12, 12, 134,
		192, 80, 80, 18, 80, 20, 20, 134, 80, 24, 24, 138, 24, 140, 140, 198,
		192, 96, 96, 34, 96, 36, 36, 134, 96, 40, 40, 138, 40, 140, 140, 198,
		96, 48, 48, 146, 48, 148, 148, 198, 48, 152, 152, 202, 152, 204, 204, 230,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 7,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 7,
		48, 49, 50, 19, 52, 21, 22, 7, 56, 25, 26, 11, 28, 13, 14, 39,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 7,
		80, 81, 82, 19, 84, 21, 22, 7, 88, 25, 26, 11, 28, 13, 14, 71,
		96, 97, 98, 35, 100, 37, 38, 7, 104, 41, 42, 11, 44, 13, 14, 71,
		112, 49, 50, 19, 52, 21, 22, 71, 56, 25, 26, 75, 28, 77, 78, 103,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 7,
		144, 145, 146, 19, 148, 21, 22, 7, 152, 25, 26, 11, 28, 13, 14, 135,
		160, 161, 162, 35, 164, 37, 38, 7, 168, 41, 42, 11, 44, 13, 14, 135,
		176, 49, 50, 19, 52, 21, 22, 135, 56, 25, 26, 139, 28, 141, 142, 167,
		192, 193, 194, 67, 196, 69, 70, 7, 200, 73, 74, 11, 76, 13, 14, 135,
		208, 81, 82, 19, 84, 21, 22, 135, 88, 25, 26, 139, 28, 141, 142, 199,
		224, 97, 98, 35, 100, 37, 38, 135, 104, 41, 42, 139, 44, 141, 142, 199,
		112, 49, 50, 147, 52, 149, 150, 199, 56, 153, 154, 203, 156, 205, 206, 231,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 8,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 8,
		0, 0, 0, 32, 0, 32, 32, 16, 0, 32, 32, 16, 32, 16, 16, 40,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 8,
		0, 0, 0, 64, 0, 64, 64, 16, 0, 64, 64, 16, 64, 16, 16, 72,
		0, 0, 0, 64, 0, 64, 64, 32, 0, 64, 64, 32, 64, 32, 32, 72,
		0, 64, 64, 32, 64, 32, 32, 80, 64, 32, 32, 80, 32, 80, 80, 104,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 8,
		0, 0, 0, 128, 0, 128, 128, 16, 0, 128, 128, 16, 128, 16, 16, 136,
		0, 0, 0, 128, 0, 128, 128, 32, 0, 128, 128, 32, 128, 32, 32, 136,
		0, 128, 128, 32, 128, 32, 32, 144, 128, 32, 32, 144, 32, 144, 144, 168,
		0, 0, 0, 128, 0, 128, 128, 64, 0, 128, 128, 64, 128, 64, 64, 136,
		0, 128, 128, 64, 128, 64, 64, 144, 128, 64, 64, 144, 64, 144, 144, 200,
		0, 128, 128, 64, 128, 64, 64, 160, 128, 64, 64, 160, 64, 160, 160, 200,
		128, 64, 64, 160, 64, 160, 160, 208, 64, 160, 160, 208, 160, 208, 208, 232,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 9,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 9,
		16, 1, 2, 33, 4, 33, 34, 17, 8, 33, 34, 17, 36, 17, 18, 41,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 9,
		16, 1, 2, 65, 4, 65, 66, 17, 8, 65, 66, 17, 68, 17, 18, 73,
		32, 1, 2, 65, 4, 65, 66, 33, 8, 65, 66, 33, 68, 33, 34, 73,
		16, 65, 66, 33, 68, 33, 34, 81, 72, 33, 34, 81, 36, 81, 82, 105,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 9,
		16, 1, 2, 129, 4, 129, 130, 17, 8, 129, 130, 17, 132, 17, 18, 137,
		32, 1, 2, 129, 4, 129, 130, 33, 8, 129, 130, 33, 132, 33, 34, 137,
		16, 129, 130, 33, 132, 33, 34, 145, 136, 33, 34, 145, 36, 145, 146, 169,
		64, 1, 2, 129, 4, 129, 130, 65, 8, 129, 130, 65, 132, 65, 66, 137,
		16, 129, 130, 65, 132, 65, 66, 145, 136, 65, 66, 145, 68, 145, 146, 201,
		32, 129, 130, 65, 132, 65, 66, 161, 136, 65, 66, 161, 68, 161, 162, 201,
		144, 65, 66, 161, 68, 161, 162, 209, 72, 161, 162, 209, 164, 209, 210, 233,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 10,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 10,
		32, 16, 16, 34, 16, 36, 36, 18, 16, 40, 40, 18, 40, 20, 20, 42,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 10,
		64, 16, 16, 66, 16, 68, 68, 18, 16, 72, 72, 18, 72, 20, 20, 74,
		64, 32, 32, 66, 32, 68, 68, 34, 32, 72, 72, 34, 72, 36, 36, 74,
		32, 80, 80, 34, 80, 36, 36, 82, 80, 40, 40, 82, 40, 84, 84, 106,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 10,
		128, 16, 16, 130, 16, 132, 132, 18, 16, 136, 136, 18, 136, 20, 20, 138,
		128, 32, 32, 130, 32, 132, 132, 34, 32, 136, 136, 34, 136, 36, 36, 138,
		32, 144, 144, 34, 144, 36, 36, 146, 144, 40, 40, 146, 40, 148, 148, 170,
		128, 64, 64, 130, 64, 132, 132, 66, 64, 136, 136, 66, 136, 68, 68, 138,
		64, 144, 144, 66, 144, 68, 68, 146, 144, 72, 72, 146, 72, 148, 148, 202,
		64, 160, 160, 66, 160, 68, 68, 162, 160, 72, 72, 162, 72, 164, 164, 202,
		160, 80, 80, 162, 80, 164, 164, 210, 80, 168, 168, 210, 168, 212, 212, 234,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 11,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 11,
		48, 17, 18, 35, 20, 37, 38, 19, 24, 41, 42, 19, 44, 21, 22, 43,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 11,
		80, 17, 18, 67, 20, 69, 70, 19, 24, 73, 74, 19, 76, 21, 22, 75,
		96, 33, 34, 67, 36, 69, 70, 35, 40, 73, 74, 35, 76, 37, 38, 75,
		48, 81, 82, 35, 84, 37, 38, 83, 88, 41, 42, 83, 44, 85, 86, 107,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 11,
		144, 17, 18, 131, 20, 133, 134, 19, 24, 137, 138, 19, 140, 21, 22, 139,
		160, 33, 34, 131, 36, 133, 134, 35, 40, 137, 138, 35, 140, 37, 38, 139,
		48, 145, 146, 35, 148, 37, 38, 147, 152, 41, 42, 147, 44, 149, 150, 171,
		192, 65, 66, 131, 68, 133, 134, 67, 72, 137, 138, 67, 140, 69, 70, 139,
		80, 145, 146, 67, 148, 69, 70, 147, 152, 73, 74, 147, 76, 149, 150, 203,
		96, 161, 162, 67, 164, 69, 70, 163, 168, 73, 74, 163, 76, 165, 166, 203,
		176, 81, 82, 163, 84, 165, 166, 211, 88, 169, 170, 211, 172, 213, 214, 235,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 12,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 12,
		0, 32, 32, 48, 32, 48, 48, 20, 32, 48, 48, 24, 48, 24, 24, 44,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 12,
		0, 64, 64, 80, 64, 80, 80, 20, 64, 80, 80, 24, 80, 24, 24, 76,
		0, 64, 64, 96, 64, 96, 96, 36, 64, 96, 96, 40, 96, 40, 40, 76,
		64, 96, 96, 48, 96, 48, 48, 84, 96, 48, 48, 88, 48, 88, 88, 108,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 12,
		0, 128, 128, 144, 128, 144, 144, 20, 128, 144, 144, 24, 144, 24, 24, 140,
		0, 128, 128, 160, 128, 160, 160, 36, 128, 160, 160, 40, 160, 40, 40, 140,
		128, 160, 160, 48, 160, 48, 48, 148, 160, 48, 48, 152, 48, 152, 152, 172,
		0, 128, 128, 192, 128, 192, 192, 68, 128, 192, 192, 72, 192, 72, 72, 140,
		128, 192, 192, 80, 192, 80, 80, 148, 192, 80, 80, 152, 80, 152, 152, 204,
		128, 192, 192, 96, 192, 96, 96, 164, 192, 96, 96, 168, 96, 168, 168, 204,
		192, 96, 96, 176, 96, 176, 176, 212, 96, 176, 176, 216, 176, 216, 216, 236,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 13,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 13,
		16, 33, 34, 49, 36, 49, 50, 21, 40, 49, 50, 25, 52, 25, 26, 45,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 13,
		16, 65, 66, 81, 68, 81, 82, 21, 72, 81, 82, 25, 84, 25, 26, 77,
		32, 65, 66, 97, 68, 97, 98, 37, 72, 97, 98, 41, 100, 41, 42, 77,
		80, 97, 98, 49, 100, 49, 50, 85, 104, 49, 50, 89, 52, 89, 90, 109,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 13,
		16, 129, 130, 145, 132, 145, 146, 21, 136, 145, 146, 25, 148, 25, 26, 141,
		32, 129, 130, 161, 132, 161, 162, 37, 136, 161, 162, 41, 164, 41, 42, 141,
		144, 161, 162, 49, 164, 49, 50, 149, 168, 49, 50, 153, 52, 153, 154, 173,
		64, 129, 130, 193, 132, 193, 194, 69, 136, 193, 194, 73, 196, 73, 74, 141,
		144, 193, 194, 81, 196, 81, 82, 149, 200, 81, 82, 153, 84, 153, 154, 205,
		160, 193, 194, 97, 196, 97, 98, 165, 200, 97, 98, 169, 100, 169, 170, 205,
		208, 97, 98, 177, 100, 177, 178, 213, 104, 177, 178, 217, 180, 217, 218, 237,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 14,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 14,
		32, 48, 48, 50, 48, 52, 52, 22, 48, 56, 56, 26, 56, 28, 28, 46,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 14,
		64, 80, 80, 82, 80, 84, 84, 22, 80, 88, 88, 26, 88, 28, 28, 78,
		64, 96, 96, 98, 96, 100, 100, 38, 96, 104, 104, 42, 104, 44, 44, 78,
		96, 112, 112, 50, 112, 52, 52, 86, 112, 56, 56, 90, 56, 92, 92, 110,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 14,
		128, 144, 144, 146, 144, 148, 148, 22, 144, 152, 152, 26, 152, 28, 28, 142,
		128, 160, 160, 162, 160, 164, 164, 38, 160, 168, 168, 42, 168, 44, 44, 142,
		160, 176, 176, 50, 176, 52, 52, 150, 176, 56, 56, 154, 56, 156, 156, 174,
		128, 192, 192, 194, 192, 196, 196, 70, 192, 200, 200, 74, 200, 76, 76, 142,
		192, 208, 208, 82, 208, 84, 84, 150, 208, 88, 88, 154, 88, 156, 156, 206,
		192, 224, 224, 98, 224, 100, 100, 166, 224, 104, 104, 170, 104, 172, 172, 206,
		224, 112, 112, 178, 112, 180, 180, 214, 112, 184, 184, 218, 184, 220, 220, 238,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 15,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 15,
		48, 49, 50, 51, 52, 53, 54, 23, 56, 57, 58, 27, 60, 29, 30, 47,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 15,
		80, 81, 82, 83, 84, 85, 86, 23, 88, 89, 90, 27, 92, 29, 30, 79,
		96, 97, 98, 99, 100, 101, 102, 39, 104, 105, 106, 43, 108, 45, 46, 79,
		112, 113, 114, 51, 116, 53, 54, 87, 120, 57, 58, 91, 60, 93, 94, 111,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 15,
		144, 145, 146, 147, 148, 149, 150, 23, 152, 153, 154, 27, 156, 29, 30, 143,
		160, 161, 162, 163, 164, 165, 166, 39, 168, 169, 170, 43, 172, 45, 46, 143,
		176, 177, 178, 51, 180, 53, 54, 151, 184, 57, 58, 155, 60, 157, 158, 175,
		192, 193, 194, 195, 196, 197, 198, 71, 200, 201, 202, 75, 204, 77, 78, 143,
		208, 209, 210, 83, 212, 85, 86, 151, 216, 89, 90, 155, 92, 157, 158, 207,
		224, 225, 226, 99, 228, 101, 102, 167, 232, 105, 106, 171, 108, 173, 174, 207,
		240, 113, 114, 179, 116, 181, 182, 215, 120, 185, 186, 219, 188, 221, 222, 239,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 48,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 80,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 96,
		0, 0, 0, 64, 0, 64, 64, 96, 0, 64, 64, 96, 64, 96, 96, 112,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 144,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 160,
		0, 0, 0, 128, 0, 128, 128, 160, 0, 128, 128, 160, 128, 160, 160, 176,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 192,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 208,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 224,
		0, 128, 128, 192, 128, 192, 192, 224, 128, 192, 192, 224, 192, 224, 224, 240,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 1,
		16, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 17,
		32, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 33,
		16, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 49,
		64, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 65,
		16, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 81,
		32, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 97,
		16, 1, 2, 65, 4, 65, 66, 97, 8, 65, 66, 97, 68, 97, 98, 113,
		128, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 129,
		16, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 145,
		32, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 161,
		16, 1, 2, 129, 4, 129, 130, 161, 8, 129, 130, 161, 132, 161, 162, 177,
		64, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 193,
		16, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 209,
		32, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 225,
		16, 129, 130, 193, 132, 193, 194, 225, 136, 193, 194, 225, 196, 225, 226, 241,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 2,
		0, 16, 16, 2, 16, 4, 4, 2, 16, 8, 8, 2, 8, 4, 4, 18,
		0, 32, 32, 2, 32, 4, 4, 2, 32, 8, 8, 2, 8, 4, 4, 34,
		32, 16, 16, 2, 16, 4, 4, 34, 16, 8, 8, 34, 8, 36, 36, 50,
		0, 64, 64, 2, 64, 4, 4, 2, 64, 8, 8, 2, 8, 4, 4, 66,
		64, 16, 16, 2, 16, 4, 4, 66, 16, 8, 8, 66, 8, 68, 68, 82,
		64, 32, 32, 2, 32, 4, 4, 66, 32, 8, 8, 66, 8, 68, 68, 98,
		32, 16, 16, 66, 16, 68, 68, 98, 16, 72, 72, 98, 72, 100, 100, 114,
		0, 128, 128, 2, 128, 4, 4, 2, 128, 8, 8, 2, 8, 4, 4, 130,
		128, 16, 16, 2, 16, 4, 4, 130, 16, 8, 8, 130, 8, 132, 132, 146,
		128, 32, 32, 2, 32, 4, 4, 130, 32, 8, 8, 130, 8, 132, 132, 162,
		32, 16, 16, 130, 16, 132, 132, 162, 16, 136, 136, 162, 136, 164, 164, 178,
		128, 64, 64, 2, 64, 4, 4, 130, 64, 8, 8, 130, 8, 132, 132, 194,
		64, 16, 16, 130, 16, 132, 132, 194, 16, 136, 136, 194, 136, 196, 196, 210,
		64, 32, 32, 130, 32, 132, 132, 194, 32, 136, 136, 194, 136, 196, 196, 226,
		32, 144, 144, 194, 144, 196, 196, 226, 144, 200, 200, 226, 200, 228, 228, 242,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 3,
		16, 17, 18, 3, 20, 5, 6, 3, 24, 9, 10, 3, 12, 5, 6, 19,
		32, 33, 34, 3, 36, 5, 6, 3, 40, 9, 10, 3, 12, 5, 6, 35,
		48, 17, 18, 3, 20, 5, 6, 35, 24, 9, 10, 35, 12, 37, 38, 51,
		64, 65, 66, 3, 68, 5, 6, 3, 72, 9, 10, 3, 12, 5, 6, 67,
		80, 17, 18, 3, 20, 5, 6, 67, 24, 9, 10, 67, 12, 69, 70, 83,
		96, 33, 34, 3, 36, 5, 6, 67, 40, 9, 10, 67, 12, 69, 70, 99,
		48, 17, 18, 67, 20, 69, 70, 99, 24, 73, 74, 99, 76, 101, 102, 115,
		128, 129, 130, 3, 132, 5, 6, 3, 136, 9, 10, 3, 12, 5, 6, 131,
		144, 17, 18, 3, 20, 5, 6, 131, 24, 9, 10, 131, 12, 133, 134, 147,
		160, 33, 34, 3, 36, 5, 6, 131, 40, 9, 10, 131, 12, 133, 134, 163,
		48, 17, 18, 131, 20, 133, 134, 163, 24, 137, 138, 163, 140, 165, 166, 179,
		192, 65, 66, 3, 68, 5, 6, 131, 72, 9, 10, 131, 12, 133, 134, 195,
		80, 17, 18, 131, 20, 133, 134, 195, 24, 137, 138, 195, 140, 197, 198, 211,
		96, 33, 34, 131, 36, 133, 134, 195, 40, 137, 138, 195, 140, 197, 198, 227,
		48, 145, 146, 195, 148, 197, 198, 227, 152, 201, 202, 227, 204, 229, 230, 243,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 4,
		0, 0, 0, 16, 0, 16, 16, 4, 0, 16, 16, 8, 16, 8, 8, 20,
		0, 0, 0, 32, 0, 32, 32, 4, 0, 32, 32, 8, 32, 8, 8, 36,
		0, 32, 32, 16, 32, 16, 16, 36, 32, 16, 16, 40, 16, 40, 40, 52,
		0, 0, 0, 64, 0, 64, 64, 4, 0, 64, 64, 8, 64, 8, 8, 68,
		0, 64, 64, 16, 64, 16, 16, 68, 64, 16, 16, 72, 16, 72, 72, 84,
		0, 64, 64, 32, 64, 32, 32, 68, 64, 32, 32, 72, 32, 72, 72, 100,
		64, 32, 32, 80, 32, 80, 80, 100, 32, 80, 80, 104, 80, 104, 104, 116,
		0, 0, 0, 128, 0, 128, 128, 4, 0, 128, 128, 8, 128, 8, 8, 132,
		0, 128, 128, 16, 128, 16, 16, 132, 128, 16, 16, 136, 16, 136, 136, 148,
		0, 128, 128, 32, 128, 32, 32, 132, 128, 32, 32, 136, 32, 136, 136, 164,
		128, 32, 32, 144, 32, 144, 144, 164, 32, 144, 144, 168, 144, 168, 168, 180,
		0, 128, 128, 64, 128, 64, 64, 132, 128, 64, 64, 136, 64, 136, 136, 196,
		128, 64, 64, 144, 64, 144, 144, 196, 64, 144, 144, 200, 144, 200, 200, 212,
		128, 64, 64, 160, 64, 160, 160, 196, 64, 160, 160, 200, 160, 200, 200, 228,
		64, 160, 160, 208, 160, 208, 208, 228, 160, 208, 208, 232, 208, 232, 232, 244,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 5,
		16, 1, 2, 17, 4, 17, 18, 5, 8, 17, 18, 9, 20, 9, 10, 21,
		32, 1, 2, 33, 4, 33, 34, 5, 8, 33, 34, 9, 36, 9, 10, 37,
		16, 33, 34, 17, 36, 17, 18, 37, 40, 17, 18, 41, 20, 41, 42, 53,
		64, 1, 2, 65, 4, 65, 66, 5, 8, 65, 66, 9, 68, 9, 10, 69,
		16, 65, 66, 17, 68, 17, 18, 69, 72, 17, 18, 73, 20, 73, 74, 85,
		32, 65, 66, 33, 68, 33, 34, 69, 72, 33, 34, 73, 36, 73, 74, 101,
		80, 33, 34, 81, 36, 81, 82, 101, 40, 81, 82, 105, 84, 105, 106, 117,
		128, 1, 2, 129, 4, 129, 130, 5, 8, 129, 130, 9, 132, 9, 10, 133,
		16, 129, 130, 17, 132, 17, 18, 133, 136, 17, 18, 137, 20, 137, 138, 149,
		32, 129, 130, 33, 132, 33, 34, 133, 136, 33, 34, 137, 36, 137, 138, 165,
		144, 33, 34, 145, 36, 145, 146, 165, 40, 145, 146, 169, 148, 169, 170, 181,
		64, 129, 130, 65, 132, 65, 66, 133, 136, 65, 66, 137, 68, 137, 138, 197,
		144, 65, 66, 145, 68, 145, 146, 197, 72, 145, 146, 201, 148, 201, 202, 213,
		160, 65, 66, 161, 68, 161, 162, 197, 72, 161, 162, 201, 164, 201, 202, 229,
		80, 161, 162, 209, 164, 209, 210, 229, 168, 209, 210, 233, 212, 233, 234, 245,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 6,
		0, 16, 16, 18, 16, 20, 20, 6, 16, 24, 24, 10, 24, 12, 12, 22,
		0, 32, 32, 34, 32, 36, 36, 6, 32, 40, 40, 10, 40, 12, 12, 38,
		32, 48, 48, 18, 48, 20, 20, 38, 48, 24, 24, 42, 24, 44, 44, 54,
		0, 64, 64, 66, 64, 68, 68, 6, 64, 72, 72, 10, 72, 12, 12, 70,
		64, 80, 80, 18, 80, 20, 20, 70, 80, 24, 24, 74, 24, 76, 76, 86,
		64, 96, 96, 34, 96, 36, 36, 70, 96, 40, 40, 74, 40, 76, 76, 102,
		96, 48, 48, 82, 48, 84, 84, 102, 48, 88, 88, 106, 88, 108, 108, 118,
		0, 128, 128, 130, 128, 132, 132, 6, 128, 136, 136, 10, 136, 12, 12, 134,
		128, 144, 144, 18, 144, 20, 20, 134, 144, 24, 24, 138, 24, 140, 140, 150,
		128, 160, 160, 34, 160, 36, 36, 134, 160, 40, 40, 138, 40, 140, 140, 166,
		160, 48, 48, 146, 48, 148, 148, 166, 48, 152, 152, 170, 152, 172, 172, 182,
		128, 192, 192, 66, 192, 68, 68, 134, 192, 72, 72, 138, 72, 140, 140, 198,
		192, 80, 80, 146, 80, 148, 148, 198, 80, 152, 152, 202, 152, 204, 204, 214,
		192, 96, 96, 162, 96, 164, 164, 198, 96, 168, 168, 202, 168, 204, 204, 230,
		96, 176, 176, 210, 176, 212, 212, 230, 176, 216, 216, 234, 216, 236, 236, 246,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 7,
		16, 17, 18, 19, 20, 21, 22, 7, 24, 25, 26, 11, 28, 13, 14, 23,
		32, 33, 34, 35, 36, 37, 38, 7, 40, 41, 42, 11, 44, 13, 14, 39,
		48, 49, 50, 19, 52, 21, 22, 39, 56, 25, 26, 43, 28, 45, 46, 55,
		64, 65, 66, 67, 68, 69, 70, 7, 72, 73, 74, 11, 76, 13, 14, 71,
		80, 81, 82, 19, 84, 21, 22, 71, 88, 25, 26, 75, 28, 77, 78, 87,
		96, 97, 98, 35, 100, 37, 38, 71, 104, 41, 42, 75, 44, 77, 78, 103,
		112, 49, 50, 83, 52, 85, 86, 103, 56, 89, 90, 107, 92, 109, 110, 119,
		128, 129, 130, 131, 132, 133, 134, 7, 136, 137, 138, 11, 140, 13, 14, 135,
		144, 145, 146, 19, 148, 21, 22, 135, 152, 25, 26, 139, 28, 141, 142, 151,
		160, 161, 162, 35, 164, 37, 38, 135, 168, 41, 42, 139, 44, 141, 142, 167,
		176, 49, 50, 147, 52, 149, 150, 167, 56, 153, 154, 171, 156, 173, 174, 183,
		192, 193, 194, 67, 196, 69, 70, 135, 200, 73, 74, 139, 76, 141, 142, 199,
		208, 81, 82, 147, 84, 149, 150, 199, 88, 153, 154, 203, 156, 205, 206, 215,
		224, 97, 98, 163, 100, 165, 166, 199, 104, 169, 170, 203, 172, 205, 206, 231,
		112, 177, 178, 211, 180, 213, 214, 231, 184, 217, 218, 235, 220, 237, 238, 247,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 16, 0, 16, 16, 24,
		0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 32, 0, 32, 32, 40,
		0, 0, 0, 32, 0, 32, 32, 48, 0, 32, 32, 48, 32, 48, 48, 56,
		0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 64, 0, 64, 64, 72,
		0, 0, 0, 64, 0, 64, 64, 80, 0, 64, 64, 80, 64, 80, 80, 88,
		0, 0, 0, 64, 0, 64, 64, 96, 0, 64, 64, 96, 64, 96, 96, 104,
		0, 64, 64, 96, 64, 96, 96, 112, 64, 96, 96, 112, 96, 112, 112, 120,
		0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 128, 128, 136,
		0, 0, 0, 128, 0, 128, 128, 144, 0, 128, 128, 144, 128, 144, 144, 152,
		0, 0, 0, 128, 0, 128, 128, 160, 0, 128, 128, 160, 128, 160, 160, 168,
		0, 128, 128, 160, 128, 160, 160, 176, 128, 160, 160, 176, 160, 176, 176, 184,
		0, 0, 0, 128, 0, 128, 128, 192, 0, 128, 128, 192, 128, 192, 192, 200,
		0, 128, 128, 192, 128, 192, 192, 208, 128, 192, 192, 208, 192, 208, 208, 216,
		0, 128, 128, 192, 128, 192, 192, 224, 128, 192, 192, 224, 192, 224, 224, 232,
		128, 192, 192, 224, 192, 224, 224, 240, 192, 224, 224, 240, 224, 240, 240, 248,
	}, {
		0, 1, 2, 1, 4, 1, 2, 1, 8, 1, 2, 1, 4, 1, 2, 9,
		16, 1, 2, 1, 4, 1, 2, 17, 8, 1, 2, 17, 4, 17, 18, 25,
		32, 1, 2, 1, 4, 1, 2, 33, 8, 1, 2, 33, 4, 33, 34, 41,
		16, 1, 2, 33, 4, 33, 34, 49, 8, 33, 34, 49, 36, 49, 50, 57,
		64, 1, 2, 1, 4, 1, 2, 65, 8, 1, 2, 65, 4, 65, 66, 73,
		16, 1, 2, 65, 4, 65, 66, 81, 8, 65, 66, 81, 68, 81, 82, 89,
		32, 1, 2, 65, 4, 65, 66, 97, 8, 65, 66, 97, 68, 97, 98, 105,
		16, 65, 66, 97, 68, 97, 98, 113, 72, 97, 98, 113, 100, 113, 114, 121,
		128, 1, 2, 1, 4, 1, 2, 129, 8, 1, 2, 129, 4, 129, 130, 137,
		16, 1, 2, 129, 4, 129, 130, 145, 8, 129, 130, 145, 132, 145, 146, 153,
		32, 1, 2, 129, 4, 129, 130, 161, 8, 129, 130, 161, 132, 161, 162, 169,
		16, 129, 130, 161, 132, 161, 162, 177, 136, 161, 162, 177, 164, 177, 178, 185,
		64, 1, 2, 129, 4, 129, 130, 193, 8, 129, 130, 193, 132, 193, 194, 201,
		16, 129, 130, 193, 132, 193, 194, 209, 136, 193, 194, 209, 196, 209, 210, 217,
		32, 129, 130, 193, 132, 193, 194, 225, 136, 193, 194, 225, 196, 225, 226, 233,
		144, 193, 194, 225, 196, 225, 226, 241, 200, 225, 226, 241, 228, 241, 242, 249,
	}, {
		0, 0, 0, 2, 0, 4, 4, 2, 0, 8, 8, 2, 8, 4, 4, 10,
		0, 16, 16, 2, 16, 4, 4, 18, 16, 8, 8, 18, 8, 20, 20, 26,
		0, 32, 32, 2, 32, 4, 4, 34, 32, 8, 8, 34, 8, 36, 36, 42,
		32, 16, 16, 34, 16, 36, 36, 50, 16, 40, 40, 50, 40, 52, 52, 58,
		0, 64, 64, 2, 64, 4, 4, 66, 64, 8, 8, 66, 8, 68, 68, 74,
		64, 16, 16, 66, 16, 68, 68, 82, 16, 72, 72, 82, 72, 84, 84, 90,
		64, 32, 32, 66, 32, 68, 68, 98, 32, 72, 72, 98, 72, 100, 100, 106,
		32, 80, 80, 98, 80, 100, 100, 114, 80, 104, 104, 114, 104, 116, 116, 122,
		0, 128, 128, 2, 128, 4, 4, 130, 128, 8, 8, 130, 8, 132, 132, 138,
		128, 16, 16, 130, 16, 132, 132, 146, 16, 136, 136, 146, 136, 148, 148, 154,
		128, 32, 32, 130, 32, 132, 132, 162, 32, 136, 136, 162, 136, 164, 164, 170,
		32, 144, 144, 162, 144, 164, 164, 178, 144, 168, 168, 178, 168, 180, 180, 186,
		128, 64, 64, 130, 64, 132, 132, 194, 64, 136, 136, 194, 136, 196, 196, 202,
		64, 144, 144, 194, 144, 196, 196, 210, 144, 200, 200, 210, 200, 212, 212, 218,
		64, 160, 160, 194, 160, 196, 196, 226, 160, 200, 200, 226, 200, 228, 228, 234,
		160, 208, 208, 226, 208, 228, 228, 242, 208, 232, 232, 242, 232, 244, 244, 250,
	}, {
		0, 1, 2, 3, 4, 5, 6, 3, 8, 9, 10, 3, 12, 5, 6, 11,
		16, 17, 18, 3, 20, 5, 6, 19, 24, 9, 10, 19, 12, 21, 22, 27,
		32, 33, 34, 3, 36, 5, 6, 35, 40, 9, 10, 35, 12, 37, 38, 43,
		48, 17, 18, 35, 20, 37, 38, 51, 24, 41, 42, 51, 44, 53, 54, 59,
		64, 65, 66, 3, 68, 5, 6, 67, 72, 9, 10, 67, 12, 69, 70, 75,
		80, 17, 18, 67, 20, 69, 70, 83, 24, 73, 74, 83, 76, 85, 86, 91,
		96, 33, 34, 67, 36, 69, 70, 99, 40, 73, 74, 99, 76, 101, 102, 107,
		48, 81, 82, 99, 84, 101, 102, 115, 88, 105, 106, 115, 108, 117, 118, 123,
		128, 129, 130, 3, 132, 5, 6, 131, 136, 9, 10, 131, 12, 133, 134, 139,
		144, 17, 18, 131, 20, 133, 134, 147, 24, 137, 138, 147, 140, 149, 150, 155,
		160, 33, 34, 131, 36, 133, 134, 163, 40, 137, 138, 163, 140, 165, 166, 171,
		48, 145, 146, 163, 148, 165, 166, 179, 152, 169, 170, 179, 172, 181, 182, 187,
		192, 65, 66, 131, 68, 133, 134, 195, 72, 137, 138, 195, 140, 197, 198, 203,
		80, 145, 146, 195, 148, 197, 198, 211, 152, 201, 202, 211, 204, 213, 214, 219,
		96, 161, 162, 195, 164, 197, 198, 227, 168, 201, 202, 227, 204, 229, 230, 235,
		176, 209, 210, 227, 212, 229, 230, 243, 216, 233, 234, 243, 236, 245, 246, 251,
	},
	{
		0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 8, 0, 8, 8, 12,
		0, 0, 0, 16, 0, 16, 16, 20, 0, 16, 16, 24, 16, 24, 24, 28,
		0, 0, 0, 32, 0, 32, 32, 36, 0, 32, 32, 40, 32, 40, 40, 44,
		0, 32, 32, 48, 32, 48, 48, 52, 32, 48, 48, 56, 48, 56, 56, 60,
		0, 0, 0, 64, 0, 64, 64, 68, 0, 64, 64, 72, 64, 72, 72, 76,
		0, 64, 64, 80, 64, 80, 80, 84, 64, 80, 80, 88, 80, 88, 88, 92,
		0, 64, 64, 96, 64, 96, 96, 100, 64, 96, 96, 104, 96, 104, 104, 108,
		64, 96, 96, 112, 96, 112, 112, 116, 96, 112, 112, 120, 112, 120, 120, 124,
		0, 0, 0, 128, 0, 128, 128, 132, 0, 128, 128, 136, 128, 136, 136, 140,
		0, 128, 128, 144, 128, 144, 144, 148, 128, 144, 144, 152, 144, 152, 152, 156,
		0, 128, 128, 160, 128, 160, 160, 164, 128, 160, 160, 168, 160, 168, 168, 172,
		128, 160, 160, 176, 160, 176, 176, 180, 160, 176, 176, 184, 176, 184, 184, 188,
		0, 128, 128, 192, 128, 192, 192, 196, 128, 192, 192, 200, 192, 200, 200, 204,
		128, 192, 192, 208, 192, 208, 208, 212, 192, 208, 208, 216, 208, 216, 216, 220,
		128, 192, 192, 224, 192, 224, 224, 228, 192, 224, 224, 232, 224, 232, 232, 236,
		192, 224, 224, 240, 224, 240, 240, 244, 224, 240, 240, 248, 240, 248, 248, 252,
	}, {
		0, 1, 2, 1, 4, 1, 2, 5, 8, 1, 2, 9, 4, 9, 10, 13,
		16, 1, 2, 17, 4, 17, 18, 21, 8, 17, 18, 25, 20, 25, 26, 29,
		32, 1, 2, 33, 4, 33, 34, 37, 8, 33, 34, 41, 36, 41, 42, 45,
		16, 33, 34, 49, 36, 49, 50, 53, 40, 49, 50, 57, 52, 57, 58, 61,
		64, 1, 2, 65, 4, 65, 66, 69, 8, 65, 66, 73, 68, 73, 74, 77,
		16, 65, 66, 81, 68, 81, 82, 85, 72, 81, 82, 89, 84, 89, 90, 93,
		32, 65, 66, 97, 68, 97, 98, 101, 72, 97, 98, 105, 100, 105, 106, 109,
		80, 97, 98, 113, 100, 113, 114, 117, 104, 113, 114, 121, 116, 121, 122, 125,
		128, 1, 2, 129, 4, 129, 130, 133, 8, 129, 130, 137, 132, 137, 138, 141,
		16, 129, 130, 145, 132, 145, 146, 149, 136, 145, 146, 153, 148, 153, 154, 157,
		32, 129, 130, 161, 132, 161, 162, 165, 136, 161, 162, 169, 164, 169, 170, 173,
		144, 161, 162, 177, 164, 177, 178, 181, 168, 177, 178, 185, 180, 185, 186, 189,
		64, 129, 130, 193, 132, 193, 194, 197, 136, 193, 194, 201, 196, 201, 202, 205,
		144, 193, 194, 209, 196, 209, 210, 213, 200, 209, 210, 217, 212, 217, 218, 221,
		160, 193, 194, 225, 196, 225, 226, 229, 200, 225, 226, 233, 228, 233, 234, 237,
		208, 225, 226, 241, 228, 241, 242, 245, 232, 241, 242, 249, 244, 249, 250, 253,
	}, {
		0, 0, 0, 2, 0, 4, 4, 6, 0, 8, 8, 10, 8, 12, 12, 14,
		0, 16, 16, 18, 16, 20, 20, 22, 16, 24, 24, 26, 24, 28, 28, 30,
		0, 32, 32, 34, 32, 36, 36, 38, 32, 40, 40, 42, 40, 44, 44, 46,
		32, 48, 48, 50, 48, 52, 52, 54, 48, 56, 56, 58, 56, 60, 60, 62,
		0, 64, 64, 66, 64, 68, 68, 70, 64, 72, 72, 74, 72, 76, 76, 78,
		64, 80, 80, 82, 80, 84, 84, 86, 80, 88, 88, 90, 88, 92, 92, 94,
		64, 96, 96, 98, 96, 100, 100, 102, 96, 104, 104, 106, 104, 108, 108, 110,
		96, 112, 112, 114, 112, 116, 116, 118, 112, 120, 120, 122, 120, 124, 124, 126,
		0, 128, 128, 130, 128, 132, 132, 134, 128, 136, 136, 138, 136, 140, 140, 142,
		128, 144, 144, 146, 144, 148, 148, 150, 144, 152, 152, 154, 152, 156, 156, 158,
		128, 160, 160, 162, 160, 164, 164, 166, 160, 168, 168, 170, 168, 172, 172, 174,
		160, 176, 176, 178, 176, 180, 180, 182, 176, 184, 184, 186, 184, 188, 188, 190,
		128, 192, 192, 194, 192, 196, 196, 198, 192, 200, 200, 202, 200, 204, 204, 206,
		192, 208, 208, 210, 208, 212, 212, 214, 208, 216, 216, 218, 216, 220, 220, 222,
		192, 224, 224, 226, 224, 228, 228, 230, 224, 232, 232, 234, 232, 236, 236, 238,
		224, 240, 240, 242, 240, 244, 244, 246, 240, 248, 248, 250, 248, 252, 252, 254,
	}, {
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
		80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95,
		96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111,
		112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127,
		128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143,
		144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159,
		160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175,
		176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191,
		192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207,
		208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223,
		224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239,
		240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255,
	},
}

// popLUT contains pre-computed population counts
var popLUT = [256]uint8{
	0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	4, 5, 5, 6, 5, 6, 6, 7, 5, 6, 6, 7, 6, 7, 7, 8,
}
