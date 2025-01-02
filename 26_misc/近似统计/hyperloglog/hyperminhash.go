// https://github.com/axiomhq/hyperminhash
// HyperLogLog 的简单实现（具体来说是 LogLog-Beta）
//
// api:
// - Add(value []byte) 添加一个值
// - Cardinality() uint64 估计去重后元素总数（基数）
// - Merge(other *Sketch) *Sketch 返回合并后的新 Sketch
// - Similarity(other *Sketch) float64 估计两个集合（Sketch）的 Jaccard 相似度
// - Intersection(other *Sketch) uint64 估计两个集合的交集大小
//
// Jaccard 相似度 = 交集数量 / 并集数量

package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
	"strconv"
)

func main() {
	sk1 := NewHyperMinHash()
	sk2 := NewHyperMinHash()

	for i := 0; i < 10000; i++ {
		sk1.Add([]byte(strconv.Itoa(i)))
	}

	fmt.Println(sk1.Cardinality()) // 10001 (should be 10000)

	for i := 3333; i < 23333; i++ {
		sk2.Add([]byte(strconv.Itoa(i)))
	}

	fmt.Println(sk2.Cardinality())     // 19977 (should be 20000)
	fmt.Println(sk1.Similarity(sk2))   // 0.284589082 (should be 0.2857326533)
	fmt.Println(sk1.Intersection(sk2)) // 6623 (should be 6667)

	sk1 = sk1.Merge(sk2)
	fmt.Println(sk1.Cardinality()) // 23271 (should be 23333)
}

const (
	p     = 14
	m     = uint32(1 << p) // 16384
	max   = 64 - p
	maxX  = math.MaxUint64 >> max
	alpha = 0.7213 / (1 + 1.079/float64(m))
	q     = 6  // the number of bits for the LogLog hash
	r     = 10 // number of bits for the bbit hash
	_2q   = 1 << q
	_2r   = 1 << r
	c     = 0.169919487159739093975315012348
)

func beta(ez float64) float64 {
	zl := math.Log(ez + 1)
	return -0.370393911*ez +
		0.070471823*zl +
		0.17393686*math.Pow(zl, 2) +
		0.16339839*math.Pow(zl, 3) +
		-0.09237745*math.Pow(zl, 4) +
		0.03738027*math.Pow(zl, 5) +
		-0.005384159*math.Pow(zl, 6) +
		0.00042419*math.Pow(zl, 7)
}

func regSumAndZeros(registers []register) (float64, float64) {
	var sum, ez float64
	for _, val := range registers {
		lz := val.lz()
		if lz == 0 {
			ez++
		}
		sum += 1 / math.Pow(2, float64(lz))
	}
	return sum, ez
}

type register uint16

func (reg register) lz() uint8 {
	return uint8(uint16(reg) >> (16 - q))
}

func newReg(lz uint8, sig uint16) register {
	return register((uint16(lz) << r) | sig)
}

// Sketch is a sketch for cardinality estimation based on LogLog counting
type Sketch struct {
	reg [m]register
}

// NewHyperMinHash returns a Sketch
func NewHyperMinHash() *Sketch {
	return new(Sketch)
}

// AddHash takes in a "hashed" value (bring your own hashing)
func (sk *Sketch) AddHash(x, y uint64) {
	k := x >> uint(max) // 取哈希值x的最高 p 位，作为寄存器的下标
	lz := uint8(bits.LeadingZeros64((x<<p)^maxX)) + 1
	sig := uint16(y << (64 - r) >> (64 - r))
	reg := newReg(lz, sig)
	if sk.reg[k] < reg {
		sk.reg[k] = reg
	}
}

// Add inserts a value into the sketch
func (sk *Sketch) Add(value []byte) {
	h1, h2 := Hash128(value, 1337)
	sk.AddHash(h1, h2)
}

// Cardinality returns the number of unique elements added to the sketch.
func (sk *Sketch) Cardinality() uint64 {
	sum, ez := regSumAndZeros(sk.reg[:])
	m := float64(m)
	return uint64(alpha * m * (m - ez) / (beta(ez) + sum))
}

// Merge returns a new union sketch of both sk and other
func (sk *Sketch) Merge(other *Sketch) *Sketch {
	m := *sk
	for i := range m.reg {
		if m.reg[i] < other.reg[i] {
			m.reg[i] = other.reg[i]
		}
	}
	return &m
}

// Similarity return a Jaccard Index similarity estimation
func (sk *Sketch) Similarity(other *Sketch) float64 {
	var C, N float64
	for i := range sk.reg {
		if sk.reg[i] != 0 && sk.reg[i] == other.reg[i] {
			C++
		}
		if sk.reg[i] != 0 || other.reg[i] != 0 {
			N++
		}
	}
	if C == 0 {
		return 0
	}

	n := float64(sk.Cardinality())
	m := float64(other.Cardinality())
	ec := sk.approximateExpectedCollisions(n, m)

	//FIXME: must be a better way to predetect this
	if C < ec {
		return 0
	}

	return (C - ec) / N
}

func (sk *Sketch) approximateExpectedCollisions(n, m float64) float64 {
	if n < m {
		n, m = m, n
	}
	if n > math.Pow(2, math.Pow(2, q)+r) {
		return math.MaxUint64
	} else if n > math.Pow(2, p+5) {
		d := (4 * n / m) / math.Pow((1+n)/m, 2)
		return c*math.Pow(2, p-r)*d + 0.5
	} else {
		return sk.expectedCollision(n, m) / float64(p)
	}
}

func (sk *Sketch) expectedCollision(n, m float64) float64 {
	var x, b1, b2 float64
	for i := 1.0; i <= _2q; i++ {
		for j := 1.0; j <= _2r; j++ {
			if i != _2q {
				den := math.Pow(2, p+r+i)
				b1 = (_2r + j) / den
				b2 = (_2r + j + 1) / den
			} else {
				den := math.Pow(2, p+r+i-1)
				b1 = j / den
				b2 = (j + 1) / den
			}
			prx := math.Pow(1-b2, n) - math.Pow(1-b1, n)
			pry := math.Pow(1-b2, m) - math.Pow(1-b1, m)
			x += (prx * pry)
		}
	}
	return (x * float64(p)) + 0.5
}

// Intersection returns number of intersections between sk and other
func (sk *Sketch) Intersection(other *Sketch) uint64 {
	sim := sk.Similarity(other)
	return uint64((sim*float64(sk.Merge(other).Cardinality()) + 0.5))
}

func Hash128(buffer []byte, seed uint64) (uint64, uint64) {

	const (
		k0 = 0xC83A91E1
		k1 = 0x8648DBDB
		k2 = 0x7BDEC03B
		k3 = 0x2F5870A5
	)

	ptr := buffer

	var v [4]uint64

	v[0] = (seed - k0) * k3
	v[1] = (seed + k1) * k2

	if len(ptr) >= 32 {
		v[2] = (seed + k0) * k2
		v[3] = (seed - k1) * k3

		for len(ptr) >= 32 {
			v[0] += binary.LittleEndian.Uint64(ptr) * k0
			ptr = ptr[8:]
			v[0] = rotate_right(v[0], 29) + v[2]
			v[1] += binary.LittleEndian.Uint64(ptr) * k1
			ptr = ptr[8:]
			v[1] = rotate_right(v[1], 29) + v[3]
			v[2] += binary.LittleEndian.Uint64(ptr) * k2
			ptr = ptr[8:]
			v[2] = rotate_right(v[2], 29) + v[0]
			v[3] += binary.LittleEndian.Uint64(ptr) * k3
			ptr = ptr[8:]
			v[3] = rotate_right(v[3], 29) + v[1]
		}

		v[2] ^= rotate_right(((v[0]+v[3])*k0)+v[1], 21) * k1
		v[3] ^= rotate_right(((v[1]+v[2])*k1)+v[0], 21) * k0
		v[0] ^= rotate_right(((v[0]+v[2])*k0)+v[3], 21) * k1
		v[1] ^= rotate_right(((v[1]+v[3])*k1)+v[2], 21) * k0
	}

	if len(ptr) >= 16 {
		v[0] += binary.LittleEndian.Uint64(ptr) * k2
		ptr = ptr[8:]
		v[0] = rotate_right(v[0], 33) * k3
		v[1] += binary.LittleEndian.Uint64(ptr) * k2
		ptr = ptr[8:]
		v[1] = rotate_right(v[1], 33) * k3
		v[0] ^= rotate_right((v[0]*k2)+v[1], 45) * k1
		v[1] ^= rotate_right((v[1]*k3)+v[0], 45) * k0
	}

	if len(ptr) >= 8 {
		v[0] += binary.LittleEndian.Uint64(ptr) * k2
		ptr = ptr[8:]
		v[0] = rotate_right(v[0], 33) * k3
		v[0] ^= rotate_right((v[0]*k2)+v[1], 27) * k1
	}

	if len(ptr) >= 4 {
		v[1] += uint64(binary.LittleEndian.Uint32(ptr)) * k2
		ptr = ptr[4:]
		v[1] = rotate_right(v[1], 33) * k3
		v[1] ^= rotate_right((v[1]*k3)+v[0], 46) * k0
	}

	if len(ptr) >= 2 {
		v[0] += uint64(binary.LittleEndian.Uint16(ptr)) * k2
		ptr = ptr[2:]
		v[0] = rotate_right(v[0], 33) * k3
		v[0] ^= rotate_right((v[0]*k2)+v[1], 22) * k1
	}

	if len(ptr) >= 1 {
		v[1] += uint64(ptr[0]) * k2
		v[1] = rotate_right(v[1], 33) * k3
		v[1] ^= rotate_right((v[1]*k3)+v[0], 58) * k0
	}

	v[0] += rotate_right((v[0]*k0)+v[1], 13)
	v[1] += rotate_right((v[1]*k1)+v[0], 37)
	v[0] += rotate_right((v[0]*k2)+v[1], 13)
	v[1] += rotate_right((v[1]*k3)+v[0], 37)

	return v[0], v[1]
}

func rotate_right(v uint64, k uint) uint64 {
	return (v >> k) | (v << (64 - k))
}
