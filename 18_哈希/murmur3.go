// https://github.com/shanzi/algo-ds/tree/master/murmur3
// MurmurHash 是一种非加密型哈希函数，适用于一般的哈希检索操作。Murmur3 是该算法的第三个版本，它能以极高的运行效率产生出良好的哈希值分布。

package main

import (
	"hash"
	"unsafe"
)

func main() {
	var test_cases_32 = []string{
		"",
		"hello",
		"hello, world",
		"19 Jan 2038 at 3:14:07 AM",
		"The quick brown fox jumps over the lazy dog.",
	}

	// 测试 New32 和 Sum32 两个函数是否能够对同一组数据生成相同的哈希值

	digest := New32()

	for _, s := range test_cases_32 {
		data := []byte(s)

		digest.Reset()
		digest.Write(data)

		hash := digest.Sum32()
		chash := Sum32(data, 0)

		if hash != chash {
			panic("hash32 mismatch")
		}
	}
}

const (
	c1_32 uint32 = 0xcc9e2d51
	c2_32 uint32 = 0x1b873593
	c3_32 uint32 = 0x85ebca6b
	c4_32 uint32 = 0xc2b2ae35
	c5_32 uint32 = 0xe6546b64
	m_32  uint32 = 5
	r1_32 uint   = 15
	r2_32 uint   = 13
)

const (
	c1_128 uint64 = 0x87c37b91114253d5
	c2_128 uint64 = 0x4cf5ad432745937f
	c3_128 uint64 = 0xff51afd7ed558ccd
	c4_128 uint64 = 0xc4ceb9fe1a85ec53
	c5_128 uint64 = 0x52dce729
	c6_128 uint64 = 0x38495ab5
)

type digest32 struct {
	len  int
	hash uint32
	tail []byte
}

func New32() hash.Hash32 {
	return &digest32{0, 0, nil}
}

func (self *digest32) Size() int {
	return 4
}

func (self *digest32) BlockSize() int {
	return 1
}

func (self *digest32) Reset() {
	self.len = 0
	self.hash = 0
	self.tail = nil
}

func (self *digest32) Write(b []byte) (n int, err error) {
	if len(b) <= 0 {
		return 0, nil
	}

	if self.tail != nil {
		if len(self.tail)+len(b) <= 3 {
			self.tail = append(self.tail, b...)
			self.len += len(b)
		} else {
			r := 4 - len(self.tail)
			bb := append(self.tail, b[:r]...)

			self.len -= len(self.tail)
			self.tail = nil

			self.Write(bb)
			self.Write(b[r:])
		}
		return len(b), nil
	}

	if len(b) <= 3 {
		self.tail = make([]byte, len(b))
		copy(self.tail, b)
		self.len += len(b)
		return len(b), nil
	}

	nblocks := len(b) >> 2
	h := self.hash
	p := unsafe.Pointer(&b[0])
	for i := 0; i < nblocks; i++ {
		k := *(*uint32)(p)

		k *= c1_32
		k = rol32(k, r1_32)
		k *= c2_32

		h ^= k
		h = rol32(h, r2_32)
		h = h*m_32 + c5_32

		p = unsafe.Pointer(uintptr(p) + 4)
	}

	if len(b)&3 != 0 {
		self.tail = make([]byte, len(b)&3)
		copy(self.tail, b[nblocks*4:])
	}

	self.hash = h
	self.len += len(b)

	return len(b), nil
}

func (self *digest32) Sum(b []byte) []byte {
	h := self.Sum32()
	return append(b, byte(h>>24), byte(h>>16), byte(h>>8), byte(h))
}

func (self *digest32) Sum32() uint32 {
	h := self.hash

	if self.tail != nil {
		k := uint32(0)
		for i, b := range self.tail {
			k |= uint32(b) << uint(i*8)
		}

		k *= c1_32
		k = rol32(k, r1_32)
		k *= c2_32

		h ^= k
	}

	h ^= uint32(self.len)
	h ^= (h >> 16)
	h *= c3_32
	h ^= (h >> 13)
	h *= c4_32
	h ^= (h >> 16)
	return h
}

func Sum32(data []byte, seed uint32) uint32 {

	var h1 = seed

	nblocks := len(data) / 4
	var p uintptr
	if len(data) > 0 {
		p = uintptr(unsafe.Pointer(&data[0]))
	}
	p1 := p + uintptr(4*nblocks)
	for ; p < p1; p += 4 {
		k1 := *(*uint32)(unsafe.Pointer(p))

		k1 *= 0xcc9e2d51
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= 0x1b873593

		h1 ^= k1
		h1 = (h1 << 13) | (h1 >> 19)
		h1 = h1*5 + 0xe6546b64
	}

	tail := data[nblocks*4:]

	var k1 uint32
	switch len(tail) & 3 {
	case 3:
		k1 ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(tail[0])
		k1 *= 0xcc9e2d51
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= 0x1b873593
		h1 ^= k1
	}

	h1 ^= uint32(len(data))

	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return h1
}

type Hash128 interface {
	hash.Hash

	Sum128() (uint64, uint64)
}

type digest128 struct {
	len   int
	hash1 uint64
	hash2 uint64
	tail  []byte
}

func New128() Hash128 {
	return &digest128{0, 0, 0, nil}
}

func (self *digest128) Size() int {
	return 8
}

func (self *digest128) BlockSize() int {
	return 1
}

func (self *digest128) Reset() {
	self.len = 0
	self.hash1 = 0
	self.hash2 = 0
	self.tail = nil
}

func (self *digest128) Write(b []byte) (int, error) {
	if len(b) <= 0 {
		return 0, nil
	}

	if self.tail != nil {
		if len(self.tail)+len(b) <= 15 {
			self.tail = append(self.tail, b...)
			self.len += len(b)
		} else {
			r := 16 - len(self.tail)
			bb := append(self.tail, b[:r]...)

			self.len -= len(self.tail)
			self.tail = nil

			self.Write(bb)
			self.Write(b[r:])
		}
		return len(b), nil
	}

	if len(b) <= 15 {
		self.tail = make([]byte, len(b))
		copy(self.tail, b)
		self.len += len(b)
		return len(b), nil
	}

	nblocks := len(b) >> 4
	h1, h2 := self.hash1, self.hash2
	p := unsafe.Pointer(&b[0])

	for i := 0; i < nblocks; i++ {
		t := (*[2]uint64)(p)
		k1, k2 := t[0], t[1]

		k1 *= c1_128
		k1 = rol64(k1, 31)
		k1 *= c2_128
		h1 ^= k1

		h1 = rol64(h1, 27)
		h1 += h2
		h1 = h1*5 + c5_128

		k2 *= c2_128
		k2 = rol64(k2, 33)
		k2 *= c1_128
		h2 ^= k2

		h2 = rol64(h2, 31)
		h2 += h1
		h2 = h2*5 + c6_128

		p = unsafe.Pointer(uintptr(p) + 16)
	}

	if len(b)&15 != 0 {
		self.tail = make([]byte, len(b)&15)
		copy(self.tail, b[nblocks*16:])
	}

	self.hash1 = h1
	self.hash2 = h2
	self.len += len(b)

	return len(b), nil
}

func (self *digest128) Sum(b []byte) []byte {
	h1, h2 := self.Sum128()

	for i := 56; i >= 0; i -= 8 {
		b = append(b, byte(h1>>uint(i)))
	}

	for i := 56; i >= 0; i -= 8 {
		b = append(b, byte(h2>>uint(i)))
	}

	return b
}

func (self *digest128) Sum128() (uint64, uint64) {
	h1, h2 := self.hash1, self.hash2

	if len(self.tail) > 8 {
		k2 := fetch64(self.tail[8:])
		k2 *= c2_128
		k2 = rol64(k2, 33)
		k2 *= c1_128
		h2 ^= k2
	}

	if len(self.tail) > 0 {
		l := len(self.tail)
		if l > 8 {
			l = 8
		}

		k1 := fetch64(self.tail[:l])
		k1 *= c1_128
		k1 = rol64(k1, 31)
		k1 *= c2_128
		h1 ^= k1
	}

	h1 ^= uint64(self.len)
	h2 ^= uint64(self.len)

	h1 += h2
	h2 += h1

	h1 = fmix64(h1)
	h2 = fmix64(h2)

	h1 += h2
	h2 += h1
	return h1, h2
}

func rol32(v uint32, l uint) uint32 {
	l &= 31
	return (v << l) | (v >> (32 - l))
}

// 循环左移
func rol64(v uint64, r uint) uint64 {
	r &= 63
	return (v << r) | (v >> (64 - r))
}

// 从字节切片中获取64位无符号整数
func fetch64(b []byte) uint64 {
	k := uint64(0)
	for i := len(b) - 1; i >= 0; i-- {
		k |= uint64(b[i]) << uint(i*8)
	}
	return k
}

func fmix64(h uint64) uint64 {
	h ^= h >> 33
	h *= c3_128
	h ^= h >> 33
	h *= c4_128
	h ^= h >> 33
	return h
}
