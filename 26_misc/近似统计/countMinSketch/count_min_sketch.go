// https://github.com/Tv0ridobro/data-structure/blob/v0.2.1/count-min-sketch/count-min-sketch.go

// Package countminsketch implements a Count-Min Sketch: a probabilistic data
// structure that serves as a frequency table of events in a stream of data.
// See https://en.wikipedia.org/wiki/Count%E2%80%93min_sketch for more details.
package main

import (
	"fmt"
	"math"
	"math/rand"
	"unsafe"
)

func main() {
	// 例如，选取误差率 1%（0.01）和置信度 99%（0.99）
	cms := NewCountMinSketch[string](0.01, 0.99)

	// 1) 一次性插入
	cms.Insert("apple")

	// 2) 插入多次（一次性增加计数）
	cms.InsertN("banana", 5)

	// 3) 在循环里插入（模拟数据流）
	for i := 0; i < 100; i++ {
		cms.Insert("orange")
	}

	countApple := cms.Count("apple")
	countBanana := cms.Count("banana")
	countOrange := cms.Count("orange")

	fmt.Printf("apple  ~ %d\n", countApple)
	fmt.Printf("banana ~ %d\n", countBanana)
	fmt.Printf("orange ~ %d\n", countOrange)

	// 插入更多的 "banana"
	cms.InsertN("banana", 2)
	fmt.Println("Count(banana) =", cms.Count("banana")) // ~7

	// 查询不存在的元素
	fmt.Println("Count(grape)  =", cms.Count("grape")) // ~0 or 正常为0

}

type CountMinSketch[T comparable] struct {
	matrix  [][]uint64
	hashers []Hasher[T]
}

// NewCountMinSketch is a constructor function that creates a new count-min sketch
// with desired error rate and confidence.
//
// errorRate：允许出现的最大相对误差。值越小，需要的内存越大，估计越精确。
// confidence：估计结果（不超过最大误差）的概率。值越接近 1，需要的哈希行数（d）越多。
func NewCountMinSketch[T comparable](errorRate, confidence float64) CountMinSketch[T] {
	w := int(math.Ceil(math.E / errorRate))
	d := int(math.Ceil(math.Log(1 / confidence)))

	hashers := make([]Hasher[T], d)
	for i := 0; i < d; i++ {
		hashers[i] = NewHasher[T]()
	}

	// 用 d 行, 每行 w 个桶
	matrix := make([][]uint64, d)
	for i := 0; i < d; i++ {
		matrix[i] = make([]uint64, w)
	}

	return CountMinSketch[T]{
		matrix:  matrix,
		hashers: hashers,
	}
}

// Insert adds an element to the count-min sketch with a count of 1.
func (c CountMinSketch[T]) Insert(elem T) {
	c.InsertN(elem, 1)
}

// InsertN adds an element to the count-min sketch with a given count.
func (c CountMinSketch[T]) InsertN(elem T, count uint64) {
	for i, hasher := range c.hashers {
		hash := hasher.Hash(elem)
		c.matrix[i][hash%uint64(len(c.matrix[i]))] += count
	}
}

// !Count 方法返回对该元素的上界频次估计.
// Count returns the approximate count of an element in the count-min sketch.
func (c CountMinSketch[T]) Count(elem T) uint64 {
	var min uint64 = math.MaxUint64
	for i, hasher := range c.hashers {
		hash := hasher.Hash(elem)
		if value := c.matrix[i][hash%uint64(len(c.matrix[i]))]; value < min {
			min = value
		}
	}
	return min
}

// #region hasher

// !Go 语言 运行时自带的哈希函数（AES-based hashing）
// Hasher hashes values of type K.
// Uses runtime AES-based hashing.
type Hasher[K comparable] struct {
	hash hashfn
	seed uintptr
}

// NewHasher creates a new Hasher[K] with a random seed.
func NewHasher[K comparable]() Hasher[K] {
	return Hasher[K]{
		hash: getRuntimeHasher[K](),
		seed: newHashSeed(),
	}
}

// NewSeed returns a copy of |h| with a new hash seed.
func NewSeed[K comparable](h Hasher[K]) Hasher[K] {
	return Hasher[K]{
		hash: h.hash,
		seed: newHashSeed(),
	}
}

// Hash hashes |key|.
func (h Hasher[K]) Hash(key K) uint64 {
	return uint64(h.Hash2(key))
}

// Hash2 hashes |key| as more flexible uintptr.
func (h Hasher[K]) Hash2(key K) uintptr {
	// promise to the compiler that pointer
	// |p| does not escape the stack.
	p := noescape(unsafe.Pointer(&key))
	return h.hash(p, h.seed)
}

type hashfn func(unsafe.Pointer, uintptr) uintptr

func getRuntimeHasher[K comparable]() (h hashfn) {
	a := any(make(map[K]struct{}))
	i := (*mapiface)(unsafe.Pointer(&a))
	h = i.typ.hasher
	return
}

func newHashSeed() uintptr {
	return uintptr(rand.Int())
}

// !把指针先变成 uintptr 再变回 unsafe.Pointer，编译器暂时就无法准确追踪这个指针是否会逃逸到堆上.
// noescape hides a pointer from escape analysis. It is the identity function
// but escape analysis doesn't think the output depends on the input.
// noescape is inlined and currently compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime (via pkg "strings"); see issues 23382 and 7921.
//
//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

type mapiface struct {
	typ *maptype
	val *hmap
}

// go/src/runtime/type.go
type maptype struct {
	typ    _type
	key    *_type
	elem   *_type
	bucket *_type
	// function for hashing keys (ptr to key, seed) -> hash
	hasher     func(unsafe.Pointer, uintptr) uintptr
	keysize    uint8
	elemsize   uint8
	bucketsize uint16
	flags      uint32
}

// go/src/runtime/map.go
type hmap struct {
	count     int
	flags     uint8
	B         uint8
	noverflow uint16
	// hash seed
	hash0      uint32
	buckets    unsafe.Pointer
	oldbuckets unsafe.Pointer
	nevacuate  uintptr
	// true type is *mapextra
	// but we don't need this data
	extra unsafe.Pointer
}

// go/src/runtime/type.go
type tflag uint8
type nameOff int32
type typeOff int32

// go/src/runtime/type.go
type _type struct {
	size       uintptr
	ptrdata    uintptr
	hash       uint32
	tflag      tflag
	align      uint8
	fieldAlign uint8
	kind       uint8
	equal      func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata     *byte
	str        nameOff
	ptrToThis  typeOff
}

// #endregion
