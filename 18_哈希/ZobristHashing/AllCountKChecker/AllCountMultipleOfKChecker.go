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

type Value = int

type hashKey = struct {
	value    Value
	modCount int
}

// 生成[min,max]范围内的随机数,并保证同一个value对应的随机数是固定的.
func RandomHash(min, max uint64) func(key hashKey) uint64 {
	pool := make(map[hashKey]uint64)
	f := func(key hashKey) uint64 {
		if hash, ok := pool[key]; ok {
			return hash
		}
		rand := rand.Uint64()%(max-min+1) + min
		pool[key] = rand
		return rand
	}
	return f
}

// 判断数据结构中每个数出现的次数是否均k的`倍数`.
// 如果为空集合,则返回True.
type AllCountMultipleOfKChecker struct {
	hash       uint64
	modCounter map[Value]int
	k          int
	randomHash func(key hashKey) uint64
}

func NewAllCountMultipleOfKChecker(k int) *AllCountMultipleOfKChecker {
	return &AllCountMultipleOfKChecker{
		hash:       0,
		modCounter: make(map[Value]int),
		k:          k,
		randomHash: RandomHash(1, (1<<61)-1),
	}
}

func (c *AllCountMultipleOfKChecker) Add(v Value) {
	count := c.modCounter[v]
	random := c.randomHash(hashKey{v, count})
	c.hash ^= random
	count++
	if count == c.k {
		count = 0
	}
	c.hash ^= c.randomHash(hashKey{v, count})
	c.modCounter[v] = count
}

// 删除前需要保证v在集合中存在.
func (c *AllCountMultipleOfKChecker) Remove(v Value) {
	count := c.modCounter[v]
	c.hash ^= c.randomHash(hashKey{v, count})
	count--
	if count == -1 {
		count = c.k - 1
	}
	c.hash ^= c.randomHash(hashKey{v, count})
	c.modCounter[v] = count
}

// 询问数据结构中每个数出现的次数是否均k的倍数.
func (c *AllCountMultipleOfKChecker) Query() bool {
	return c.hash == 0
}

func (c *AllCountMultipleOfKChecker) GetHash() uint64 {
	return c.hash
}

func (c *AllCountMultipleOfKChecker) Clear() {
	c.hash = 0
	c.modCounter = make(map[Value]int)
}
