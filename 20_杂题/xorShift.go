// https://ei1333.github.io/library/other/xor-shift.hpp
// 生成伪随机数的算法 xorshift
// Usage:
//  NewXorShift(seed uint64) : 指定seed初始化
//  Get64() : 生成[0, 2^64)的随机数
//  Get64Range(l, r uint64) : 生成[l, r)的随机数
//  Get32() : 生成[0, 2^32)的随机数
//  Get32Range(l, r uint32) : 生成[l, r)的随机数
//  Probability() : 生成[0.0, 1.0)的随机数

package main

import "fmt"

const (
	R    float64 = 1.0 / 0xffffffff
	SEED uint64  = 88172645463325252
)

type XorShift struct {
	seed uint64
}

func NewXorShift(seed uint64) *XorShift {
	return &XorShift{seed: seed}
}

func (xs *XorShift) Get64() uint64 {
	xs.seed ^= xs.seed << 7
	xs.seed ^= xs.seed >> 9
	return xs.seed
}

func (xs *XorShift) Get64Range(l, r uint64) uint64 {
	return l + xs.Get64()%(r-l)
}

func (xs *XorShift) Get32() uint32 {
	xs.seed ^= xs.seed << 7
	xs.seed ^= xs.seed >> 9
	return uint32(xs.seed)
}

func (xs *XorShift) Get32Range(l, r uint32) uint32 {
	return l + xs.Get32()%(r-l)
}

func (xs *XorShift) Probability() float64 {
	return float64(xs.Get32()) * R
}

func main() {
	xorShift := NewXorShift(88172645463325252)
	fmt.Println(xorShift.Get32())
	fmt.Println(xorShift.Get64())
	fmt.Println(xorShift.Probability())
	fmt.Println(xorShift.Get32Range(1, 10))
	fmt.Println(xorShift.Get64Range(1, 10))
}
