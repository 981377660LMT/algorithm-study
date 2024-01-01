package main

import (
	"fmt"
	"math/rand"
)

func main() {
	C := NewAllCountMultipleOfKChecker(3)
	fmt.Println(C.Query())
	C.Add(1)
	fmt.Println(C.Query())
	C.Add(1)
	fmt.Println(C.Query())
	C.Add(1)
	fmt.Println(C.Query())
	C.Add(2)
	fmt.Println(C.Query())
	C.Remove(2)
	fmt.Println(C.Query())
}

type T int

// 判断数据结构中每个数出现的次数是否均k的`倍数`.
type AllCountMultipleOfKChecker struct {
	pool    map[_pair]uint64
	hash    uint64
	k       int
	counter map[T]int
}

func NewAllCountMultipleOfKChecker(k int) *AllCountMultipleOfKChecker {
	res := &AllCountMultipleOfKChecker{}
	res.pool = make(map[_pair]uint64)
	res.k = k
	res.counter = make(map[T]int)
	return res
}

func (c *AllCountMultipleOfKChecker) Add(v T) {
	count := c.counter[v]
	c.hash ^= c._hash(v, count)
	count++
	if count == c.k {
		count = 0
	}
	c.counter[v] = count
	c.hash ^= c._hash(v, count)
}

// 删除前需要保证v在集合中存在.
func (c *AllCountMultipleOfKChecker) Remove(v T) {
	count := c.counter[v]
	c.hash ^= c._hash(v, count)
	count--
	if count == -1 {
		count = c.k - 1
	}
	c.counter[v] = count
	c.hash ^= c._hash(v, count)
}

// 询问数据结构中每个数出现的次数是否均k的倍数.
func (c *AllCountMultipleOfKChecker) Query() bool {
	return c.hash == 0
}

func (c *AllCountMultipleOfKChecker) GetHash() uint64 {
	return c.hash
}

func (c *AllCountMultipleOfKChecker) _hash(v T, countMod int) uint64 {
	p := _pair{v, countMod}
	if res, ok := c.pool[p]; ok {
		return res
	}
	res := c._randUint61()
	c.pool[p] = res
	return res
}

// [1,1<<61-1]
func (c *AllCountMultipleOfKChecker) _randUint61() uint64 {
	return rand.Uint64()%(1<<61-1) + 1
}

type _pair struct {
	value    T
	countMod int
}
