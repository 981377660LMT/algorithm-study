// 还是不如内置的map快

package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	// mp := NewHashMapUint64()
	// mp.Set(1, 2)
	// mp.Set(2, 3)
	// mp.EnumerateAll(func(key uint64, value uint64) {
	// 	fmt.Println(key, value)
	// })
	// fmt.Println(mp.Get(1))
	yosupo()
}

// https://judge.yosupo.jp/problem/associative_array
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)
	mp := NewHashMapUint64WithCapicity(int(1e6))
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var pos, val uint64
			fmt.Fscan(in, &pos, &val)
			mp.Set(pos, val)
		} else if op == 1 {
			var pos uint64
			fmt.Fscan(in, &pos)
			fmt.Fprintln(out, mp.GetOrDefault(pos, 0))
		}
	}
}

type T = uint64
type pair = struct {
	first  uint64
	second T
}

// ma
type HashMapUint64 struct {
	capicity    uint32 // 2的幂次
	mask        uint64
	keys        []uint64
	values      []T
	used        []bool
	fixedRandom uint64
}

func NewHashMapUint64() *HashMapUint64 {
	return NewHashMapUint64WithCapicity(0)
}

func NewHashMapUint64WithCapicity(initCapicity int) *HashMapUint64 {
	res := &HashMapUint64{fixedRandom: uint64(time.Now().UnixNano()/2 + 1)}
	res.build(initCapicity)
	return res
}

func (mp *HashMapUint64) Clear() {
	mp.build(0)
}

func (mp *HashMapUint64) Size() int {
	return len(mp.used) - int(mp.capicity)
}

func (mp *HashMapUint64) Get(x uint64) T {
	i := mp.index(x)
	if mp.used[i] {
		return mp.values[i]
	}
	panic("not found")
}

func (mp *HashMapUint64) GetOrDefault(x uint64, defaultValue T) T {
	i := mp.index(x)
	if mp.used[i] {
		return mp.values[i]
	}
	return defaultValue
}

func (mp *HashMapUint64) Has(x uint64) bool {
	i := mp.index(x)
	return mp.used[i] && mp.keys[i] == x
}

func (mp *HashMapUint64) EnumerateAll(f func(key uint64, value T)) {
	for i, v := range mp.used {
		if v {
			f(mp.keys[i], mp.values[i])
		}
	}
}

func (mp *HashMapUint64) Set(x uint64, v T) {
	if mp.capicity == 0 {
		mp.extend()
	}
	i := mp.index(x)
	mp.used[i] = true
	mp.keys[i] = x
	mp.values[i] = v
	mp.capicity--
}

func (mp *HashMapUint64) build(size int) {
	k := 8
	target := size * 2
	for k < target {
		k *= 2
	}
	mp.capicity = uint32(k / 2)
	mp.mask = uint64(k - 1)
	mp.keys = make([]uint64, k)
	mp.values = make([]T, k)
	mp.used = make([]bool, k)
}

func (mp *HashMapUint64) hash(x uint64) uint64 {
	x += mp.fixedRandom
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	return (x ^ (x >> 31)) & mp.mask
}

// 迁移数据.
func (mp *HashMapUint64) extend() {
	data := make([]pair, len(mp.used)-int(mp.capicity))
	for i, v := range mp.used {
		if v {
			data = append(data, pair{mp.keys[i], mp.values[i]})
		}
	}
	mp.build(2 * len(data))
	for _, v := range data {
		mp.Set(v.first, v.second)
	}
}

// 链地址法.
func (mp *HashMapUint64) index(k uint64) int {
	i := uint64(0)
	for i = mp.hash(k); mp.used[i] && mp.keys[i] != k; i = (i + 1) & mp.mask {
	}
	return int(i)
}
