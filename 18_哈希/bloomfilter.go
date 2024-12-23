/*
Package bloom provides data structures and methods for creating Bloom filters.

A Bloom filter is a representation of a set of _n_ items, where the main
requirement is to make membership queries; _i.e._, whether an item is a
member of a set.

A Bloom filter has two parameters: _m_, a maximum size (typically a reasonably large
multiple of the cardinality of the set to represent) and _k_, the number of hashing
functions on elements of the set. (The actual hashing functions are important, too,
but this is not a parameter for this implementation). A Bloom filter is backed by
a BitSet; a key is represented in the filter by setting the bits at each value of the
hashing functions (modulo _m_). Set membership is done by _testing_ whether the
bits at each value of the hashing functions (again, modulo _m_) are set. If so,
the item is in the set. If the item is actually in the set, a Bloom filter will
never fail (the true positive rate is 1.0); but it is susceptible to false
positives. The art is to choose _k_ and _m_ correctly.

In this implementation, the hashing functions used is murmurhash,
a non-cryptographic hashing function.

This implementation accepts keys for setting as testing as []byte. Thus, to
add a string item, "Love":

	uint n = 1000
	filter := bloom.New(20*n, 5) // load of 20, 5 keys
	filter.Add([]byte("Love"))

Similarly, to test if "Love" is in bloom:

	if filter.Test([]byte("Love"))

For numeric data, I recommend that you look into the binary/encoding library. But,
for example, to add a uint32 to the filter:

	i := uint32(100)
	n1 := make([]byte,4)
	binary.BigEndian.PutUint32(n1,i)
	f.Add(n1)

Finally, there is a method to estimate the false positive rate of a
Bloom filter with _m_ bits and _k_ hashing functions for a set of size _n_:

	if bloom.EstimateFalsePositiveRate(20*n, 5, n) > 0.001 ...

You can use it to validate the computed m, k parameters:

	m, k := bloom.EstimateParameters(n, fp)
	ActualfpRate := bloom.EstimateFalsePositiveRate(m, k, n)

or

	f := bloom.NewWithEstimates(n, fp)
	ActualfpRate := bloom.EstimateFalsePositiveRate(f.m, f.k, n)

You would expect ActualfpRate to be close to the desired fp in these cases.

The EstimateFalsePositiveRate function creates a temporary Bloom filter. It is
also relatively expensive and only meant for validation.
*/
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/bits"
	"strconv"
	"unsafe"
)

func main() {
	bf := NewBloomFilter(1e6, 5)
	bf.Add([]byte("hello"))
	bf.Add([]byte("world"))
	bf.Add([]byte("bloom"))
	bf.Add([]byte("filter"))

	bf.AddString("hello")

	fmt.Println(bf.Test([]byte("hello")), bf.TestString("hello"))
	fmt.Println(EstimateFalsePositiveRate(1e7, 5, 1e5))
}

// A BloomFilter is a representation of a set of _n_ items, where the main
// requirement is to make membership queries; _i.e._, whether an item is a
// member of a set.
type BloomFilter struct {
	m uint // bitset的长度
	k uint // 哈希函数的数量
	b BitSet
}

func max(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}

// New creates a new Bloom filter with _m_ bits and _k_ hashing functions
// We force _m_ and _k_ to be at least one to avoid panics.
func NewBloomFilter(m uint, k uint) *BloomFilter {
	return &BloomFilter{max(1, m), max(1, k), *NewBitSet(m)}
}

// From creates a new Bloom filter with len(_data_) * 64 bits and _k_ hashing
// functions. The data slice is not going to be reset.
func From(data []uint64, k uint) *BloomFilter {
	m := uint(len(data) * 64)
	return FromWithM(data, m, k)
}

// FromWithM creates a new Bloom filter with _m_ length, _k_ hashing functions.
// The data slice is not going to be reset.
func FromWithM(data []uint64, m, k uint) *BloomFilter {
	return &BloomFilter{m, k, *BFrom(data)}
}

// baseHashes returns the four hash values of data that are used to create k
// hashes
func baseHashes(data []byte) [4]uint64 {
	var d digest128 // murmur hashing
	hash1, hash2, hash3, hash4 := d.sum256(data)
	return [4]uint64{
		hash1, hash2, hash3, hash4,
	}
}

// location returns the ith hashed location using the four base hash values
func location(h [4]uint64, i uint) uint64 {
	ii := uint64(i)
	return h[ii%2] + ii*h[2+(((ii+(ii%2))%4)/2)]
}

// location returns the ith hashed location using the four base hash values
func (f *BloomFilter) location(h [4]uint64, i uint) uint {
	return uint(location(h, i) % uint64(f.m))
}

// EstimateParameters estimates requirements for m and k.
// Based on https://bitbucket.org/ww/bloom/src/829aa19d01d9/bloom.go
// used with permission.
func EstimateParameters(n uint, p float64) (m uint, k uint) {
	m = uint(math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
	k = uint(math.Ceil(math.Log(2) * float64(m) / float64(n)))
	return
}

// NewWithEstimates creates a new Bloom filter for about n items with fp
// false positive rate
func NewWithEstimates(n uint, fp float64) *BloomFilter {
	m, k := EstimateParameters(n, fp)
	return NewBloomFilter(m, k)
}

// Cap returns the capacity, _m_, of a Bloom filter
func (f *BloomFilter) Cap() uint {
	return f.m
}

// K returns the number of hash functions used in the BloomFilter
func (f *BloomFilter) K() uint {
	return f.k
}

// BitSet returns the underlying bitset for this filter.
func (f *BloomFilter) BitSet() BitSet {
	return f.b
}

// Add data to the Bloom Filter. Returns the filter (allows chaining)
func (f *BloomFilter) Add(data []byte) *BloomFilter {
	h := baseHashes(data)
	for i := uint(0); i < f.k; i++ {
		f.b.Set(f.location(h, i))
	}
	return f
}

// Merge the data from two Bloom Filters.
func (f *BloomFilter) Merge(g *BloomFilter) error {
	// Make sure the m's and k's are the same, otherwise merging has no real use.
	if f.m != g.m {
		return fmt.Errorf("m's don't match: %d != %d", f.m, g.m)
	}

	if f.k != g.k {
		return fmt.Errorf("k's don't match: %d != %d", f.m, g.m)
	}

	f.b.InPlaceUnion(&g.b)
	return nil
}

// Copy creates a copy of a Bloom filter.
func (f *BloomFilter) Copy() *BloomFilter {
	fc := NewBloomFilter(f.m, f.k)
	fc.Merge(f) // #nosec
	return fc
}

// AddString to the Bloom Filter. Returns the filter (allows chaining)
func (f *BloomFilter) AddString(data string) *BloomFilter {
	return f.Add([]byte(data))
}

// Test returns true if the data is in the BloomFilter, false otherwise.
// If true, the result might be a false positive. If false, the data
// is definitely not in the set.
func (f *BloomFilter) Test(data []byte) bool {
	h := baseHashes(data)
	for i := uint(0); i < f.k; i++ {
		if !f.b.Test(f.location(h, i)) {
			return false
		}
	}
	return true
}

// TestString returns true if the string is in the BloomFilter, false otherwise.
// If true, the result might be a false positive. If false, the data
// is definitely not in the set.
func (f *BloomFilter) TestString(data string) bool {
	return f.Test([]byte(data))
}

// TestLocations returns true if all locations are set in the BloomFilter, false
// otherwise.
func (f *BloomFilter) TestLocations(locs []uint64) bool {
	for i := 0; i < len(locs); i++ {
		if !f.b.Test(uint(locs[i] % uint64(f.m))) {
			return false
		}
	}
	return true
}

// TestAndAdd is equivalent to calling Test(data) then Add(data).
// The filter is written to unconditionnally: even if the element is present,
// the corresponding bits are still set. See also TestOrAdd.
// Returns the result of Test.
func (f *BloomFilter) TestAndAdd(data []byte) bool {
	present := true
	h := baseHashes(data)
	for i := uint(0); i < f.k; i++ {
		l := f.location(h, i)
		if !f.b.Test(l) {
			present = false
		}
		f.b.Set(l)
	}
	return present
}

// TestAndAddString is the equivalent to calling Test(string) then Add(string).
// The filter is written to unconditionnally: even if the string is present,
// the corresponding bits are still set. See also TestOrAdd.
// Returns the result of Test.
func (f *BloomFilter) TestAndAddString(data string) bool {
	return f.TestAndAdd([]byte(data))
}

// TestOrAdd is equivalent to calling Test(data) then if not present Add(data).
// If the element is already in the filter, then the filter is unchanged.
// Returns the result of Test.
func (f *BloomFilter) TestOrAdd(data []byte) bool {
	present := true
	h := baseHashes(data)
	for i := uint(0); i < f.k; i++ {
		l := f.location(h, i)
		if !f.b.Test(l) {
			present = false
			f.b.Set(l)
		}
	}
	return present
}

// TestOrAddString is the equivalent to calling Test(string) then if not present Add(string).
// If the string is already in the filter, then the filter is unchanged.
// Returns the result of Test.
func (f *BloomFilter) TestOrAddString(data string) bool {
	return f.TestOrAdd([]byte(data))
}

// ClearAll clears all the data in a Bloom filter, removing all keys
func (f *BloomFilter) ClearAll() *BloomFilter {
	f.b.ClearAll()
	return f
}

// 估算估算布隆过滤器的假阳性率，计算公式  (1 - e^(-k * n / m))^k.
// EstimateFalsePositiveRate returns, for a BloomFilter of m bits
// and k hash functions, an estimation of the false positive rate when
//
//	storing n entries. This is an empirical, relatively slow
//
// test using integers as keys.
// This function is useful to validate the implementation.
func EstimateFalsePositiveRate(m, k, n uint) (fpRate float64) {
	rounds := uint32(100000)
	// We construct a new filter.
	f := NewBloomFilter(m, k)
	n1 := make([]byte, 4)
	// We populate the filter with n values.
	for i := uint32(0); i < uint32(n); i++ {
		binary.BigEndian.PutUint32(n1, i)
		f.Add(n1)
	}
	fp := 0
	// test for number of rounds
	for i := uint32(0); i < rounds; i++ {
		binary.BigEndian.PutUint32(n1, i+uint32(n)+1)
		if f.Test(n1) {
			fp++
		}
	}
	fpRate = float64(fp) / (float64(rounds))
	return
}

// Approximating the number of items
// https://en.wikipedia.org/wiki/Bloom_filter#Approximating_the_number_of_items_in_a_Bloom_filter
func (f *BloomFilter) ApproximatedSize() uint32 {
	x := float64(f.b.Count())
	m := float64(f.Cap())
	k := float64(f.K())
	size := -1 * m / k * math.Log(1-x/m) / math.Log(math.E)
	return uint32(math.Floor(size + 0.5)) // round
}

// bloomFilterJSON is an unexported type for marshaling/unmarshaling BloomFilter struct.
type bloomFilterJSON struct {
	M uint   `json:"m"`
	K uint   `json:"k"`
	B BitSet `json:"b"`
}

// MarshalJSON implements json.Marshaler interface.
func (f BloomFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(bloomFilterJSON{f.m, f.k, f.b})
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (f *BloomFilter) UnmarshalJSON(data []byte) error {
	var j bloomFilterJSON
	err := json.Unmarshal(data, &j)
	if err != nil {
		return err
	}
	f.m = j.M
	f.k = j.K
	f.b = j.B
	return nil
}

// WriteTo writes a binary representation of the BloomFilter to an i/o stream.
// It returns the number of bytes written.
//
// Performance: if this function is used to write to a disk or network
// connection, it might be beneficial to wrap the stream in a bufio.Writer.
// E.g.,
//
//	      f, err := os.Create("myfile")
//		       w := bufio.NewWriter(f)
func (f *BloomFilter) WriteTo(stream io.Writer) (int64, error) {
	err := binary.Write(stream, binary.BigEndian, uint64(f.m))
	if err != nil {
		return 0, err
	}
	err = binary.Write(stream, binary.BigEndian, uint64(f.k))
	if err != nil {
		return 0, err
	}
	numBytes, err := f.b.WriteTo(stream)
	return numBytes + int64(2*binary.Size(uint64(0))), err
}

// ReadFrom reads a binary representation of the BloomFilter (such as might
// have been written by WriteTo()) from an i/o stream. It returns the number
// of bytes read.
//
// Performance: if this function is used to read from a disk or network
// connection, it might be beneficial to wrap the stream in a bufio.Reader.
// E.g.,
//
//	f, err := os.Open("myfile")
//	r := bufio.NewReader(f)
func (f *BloomFilter) ReadFrom(stream io.Reader) (int64, error) {
	var m, k uint64
	err := binary.Read(stream, binary.BigEndian, &m)
	if err != nil {
		return 0, err
	}
	err = binary.Read(stream, binary.BigEndian, &k)
	if err != nil {
		return 0, err
	}
	b := &BitSet{}
	numBytes, err := b.ReadFrom(stream)
	if err != nil {
		return 0, err
	}
	f.m = uint(m)
	f.k = uint(k)
	f.b = *b
	return numBytes + int64(2*binary.Size(uint64(0))), nil
}

// GobEncode implements gob.GobEncoder interface.
func (f *BloomFilter) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	_, err := f.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GobDecode implements gob.GobDecoder interface.
func (f *BloomFilter) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := f.ReadFrom(buf)

	return err
}

// MarshalBinary implements binary.BinaryMarshaler interface.
func (f *BloomFilter) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	_, err := f.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements binary.BinaryUnmarshaler interface.
func (f *BloomFilter) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := f.ReadFrom(buf)

	return err
}

// Equal tests for the equality of two Bloom filters
func (f *BloomFilter) Equal(g *BloomFilter) bool {
	return f.m == g.m && f.k == g.k && f.b.Equal(&g.b)
}

// Locations returns a list of hash locations representing a data item.
func Locations(data []byte, k uint) []uint64 {
	locs := make([]uint64, k)

	// calculate locations
	h := baseHashes(data)
	for i := uint(0); i < k; i++ {
		locs[i] = location(h, i)
	}

	return locs
}

// #region murmur3.go

const (
	c1_128     = 0x87c37b91114253d5
	c2_128     = 0x4cf5ad432745937f
	block_size = 16
)

// digest128 represents a partial evaluation of a 128 bites hash.
type digest128 struct {
	h1 uint64 // Unfinalized running hash part 1.
	h2 uint64 // Unfinalized running hash part 2.
}

// bmix will hash blocks (16 bytes)
func (d *digest128) bmix(p []byte) {
	nblocks := len(p) / block_size
	for i := 0; i < nblocks; i++ {
		b := (*[16]byte)(unsafe.Pointer(&p[i*block_size]))
		k1, k2 := binary.LittleEndian.Uint64(b[:8]), binary.LittleEndian.Uint64(b[8:])
		d.bmix_words(k1, k2)
	}
}

// bmix_words will hash two 64-bit words (16 bytes)
func (d *digest128) bmix_words(k1, k2 uint64) {
	h1, h2 := d.h1, d.h2

	k1 *= c1_128
	k1 = bits.RotateLeft64(k1, 31)
	k1 *= c2_128
	h1 ^= k1

	h1 = bits.RotateLeft64(h1, 27)
	h1 += h2
	h1 = h1*5 + 0x52dce729

	k2 *= c2_128
	k2 = bits.RotateLeft64(k2, 33)
	k2 *= c1_128
	h2 ^= k2

	h2 = bits.RotateLeft64(h2, 31)
	h2 += h1
	h2 = h2*5 + 0x38495ab5
	d.h1, d.h2 = h1, h2
}

// sum128 computers two 64-bit hash value. It is assumed that
// bmix was first called on the data to process complete blocks
// of 16 bytes. The 'tail' is a slice representing the 'tail' (leftover
// elements, fewer than 16). If pad_tail is true, we make it seem like
// there is an extra element with value 1 appended to the tail.
// The length parameter represents the full length of the data (including
// the blocks of 16 bytes, and, if pad_tail is true, an extra byte).
func (d *digest128) sum128(pad_tail bool, length uint, tail []byte) (h1, h2 uint64) {
	h1, h2 = d.h1, d.h2

	var k1, k2 uint64
	if pad_tail {
		switch (len(tail) + 1) & 15 {
		case 15:
			k2 ^= uint64(1) << 48
			break
		case 14:
			k2 ^= uint64(1) << 40
			break
		case 13:
			k2 ^= uint64(1) << 32
			break
		case 12:
			k2 ^= uint64(1) << 24
			break
		case 11:
			k2 ^= uint64(1) << 16
			break
		case 10:
			k2 ^= uint64(1) << 8
			break
		case 9:
			k2 ^= uint64(1) << 0

			k2 *= c2_128
			k2 = bits.RotateLeft64(k2, 33)
			k2 *= c1_128
			h2 ^= k2

			break

		case 8:
			k1 ^= uint64(1) << 56
			break
		case 7:
			k1 ^= uint64(1) << 48
			break
		case 6:
			k1 ^= uint64(1) << 40
			break
		case 5:
			k1 ^= uint64(1) << 32
			break
		case 4:
			k1 ^= uint64(1) << 24
			break
		case 3:
			k1 ^= uint64(1) << 16
			break
		case 2:
			k1 ^= uint64(1) << 8
			break
		case 1:
			k1 ^= uint64(1) << 0
			k1 *= c1_128
			k1 = bits.RotateLeft64(k1, 31)
			k1 *= c2_128
			h1 ^= k1
		}

	}
	switch len(tail) & 15 {
	case 15:
		k2 ^= uint64(tail[14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(tail[13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(tail[12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(tail[11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(tail[10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(tail[9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(tail[8]) << 0

		k2 *= c2_128
		k2 = bits.RotateLeft64(k2, 33)
		k2 *= c1_128
		h2 ^= k2

		fallthrough

	case 8:
		k1 ^= uint64(tail[7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(tail[0]) << 0
		k1 *= c1_128
		k1 = bits.RotateLeft64(k1, 31)
		k1 *= c2_128
		h1 ^= k1
	}

	h1 ^= uint64(length)
	h2 ^= uint64(length)

	h1 += h2
	h2 += h1

	h1 = fmix64(h1)
	h2 = fmix64(h2)

	h1 += h2
	h2 += h1

	return h1, h2
}

func fmix64(k uint64) uint64 {
	k ^= k >> 33
	k *= 0xff51afd7ed558ccd
	k ^= k >> 33
	k *= 0xc4ceb9fe1a85ec53
	k ^= k >> 33
	return k
}

// sum256 will compute 4 64-bit hash values from the input.
// It is designed to never allocate memory on the heap. So it
// works without any byte buffer whatsoever.
// It is designed to be strictly equivalent to
//
//				a1 := []byte{1}
//	         hasher := murmur3.New128()
//	         hasher.Write(data) // #nosec
//	         v1, v2 := hasher.Sum128()
//	         hasher.Write(a1) // #nosec
//	         v3, v4 := hasher.Sum128()
//
// See TestHashRandom.
func (d *digest128) sum256(data []byte) (hash1, hash2, hash3, hash4 uint64) {
	// We always start from zero.
	d.h1, d.h2 = 0, 0
	// Process as many bytes as possible.
	d.bmix(data)
	// We have enough to compute the first two 64-bit numbers
	length := uint(len(data))
	tail_length := length % block_size
	tail := data[length-tail_length:]
	hash1, hash2 = d.sum128(false, length, tail)
	// Next we want to 'virtually' append 1 to the input, but,
	// we do not want to append to an actual array!!!
	if tail_length+1 == block_size {
		// We are left with no tail!!!
		word1 := binary.LittleEndian.Uint64(tail[:8])
		word2 := uint64(binary.LittleEndian.Uint32(tail[8 : 8+4]))
		word2 = word2 | (uint64(tail[12]) << 32) | (uint64(tail[13]) << 40) | (uint64(tail[14]) << 48)
		// We append 1.
		word2 = word2 | (uint64(1) << 56)
		// We process the resulting 2 words.
		d.bmix_words(word1, word2)
		tail := data[length:] // empty slice, deliberate.
		hash3, hash4 = d.sum128(false, length+1, tail)
	} else {
		// We still have a tail (fewer than 15 bytes) but we
		// need to append '1' to it.
		hash3, hash4 = d.sum128(true, length+1, tail)
	}

	return hash1, hash2, hash3, hash4
}

// #endregion

// #region bitset

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
func BFrom(buf []uint64) *BitSet {
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
// It is not a copy, so changes to the returned slice will affect the
// It is meant for advanced users.
//
// Deprecated: Bytes is deprecated. Use [BitSet.Words] instead.
func (b *BitSet) Bytes() []uint64 {
	return b.set
}

// Words returns the bitset as array of 64-bit words, giving direct access to the internal representation.
// It is not a copy, so changes to the returned slice will affect the
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
func NewBitSet(length uint) (bset *BitSet) {
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
func MustNewBitSet(length uint) (bset *BitSet) {
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

// Len returns the number of bits in the
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

// String creates a string representation of the It is only intended for
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

// ClearAll clears the entire
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
// In case of allocation failure, the function will return an empty
func (b *BitSet) Clone() *BitSet {
	c := NewBitSet(b.length)
	if b.set != nil { // Clone should not modify current object
		copy(c.set, b.set)
	}
	return c
}

// Copy into a destination BitSet using the Go array copy semantics:
// the number of bits copied is the minimum of the number of bits in the current
// BitSet (Len()) and the destination
// We return the number of bits copied in the destination
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
// In case of allocation failure, the function will return an empty
func (b *BitSet) Intersection(compare *BitSet) (result *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	result = NewBitSet(b.length)
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
// In case of allocation failure, the function will return an empty
func (b *BitSet) Complement() (result *BitSet) {
	panicIfNull(b)
	result = NewBitSet(b.length)
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
// The length is the number of bits in the
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
// that are set in the
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

// #endregion
