package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// ===================== 全局常量、帮助函数 ====================== //

// 你可以根据实际需要更改 p (素数) 和 m (哈希表大小)。
// p 需要比可能的输入范围大。如果输入数据是 32 位整数，你可以选
// 一个比 2^32 大的素数，比如下面这个值:
var p = big.NewInt(4294967311) // 这是一个常用的大素数 (~4.294967311e9)
const m = 101                  // 哈希表桶数，可自行修改

// randomInRange 生成 [1, p-1] 范围内的随机大整数
func randomInRange(max *big.Int) *big.Int {
	// crypto/rand 包可以产生更安全的随机数
	// 我们这里保证 1 <= result < p
	// (注意: 也可以选 0 <= b < p; 但通常确保 a!=0)
	for {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		if n.Sign() > 0 { // n > 0
			return n
		}
	}
}

// ===================== Universal Hash 结构定义 ====================== //

type UniversalHash struct {
	a *big.Int
	b *big.Int
}

// NewUniversalHash 随机生成一组 (a, b)，形成一个哈希函数
func NewUniversalHash() *UniversalHash {
	// 随机选 a, b in [1, p-1]
	a := randomInRange(p)
	b := randomInRange(p)
	return &UniversalHash{a: a, b: b}
}

// HashFunc 实际的哈希函数实现： (a*x + b) mod p  再 mod m
func (uh *UniversalHash) HashFunc(x int64) int {
	// 先把 x 转成 *big.Int
	X := big.NewInt(x)

	// compute a*x + b mod p
	// big.Int 提供了大整数运算
	AX := new(big.Int).Mul(uh.a, X)
	AXB := new(big.Int).Add(AX, uh.b)
	AXB.Mod(AXB, p) // 取 mod p

	// 再取 mod m
	hashVal := AXB.Mod(AXB, big.NewInt(m)) // 取 mod m
	return int(hashVal.Int64())
}

// ===================== 简单的哈希表示例 ====================== //

// HashTable 用一个切片 of []int 来存储数据，碰撞时就简单append
type HashTable struct {
	Buckets [][]int64
	Hash    *UniversalHash
}

// NewHashTable 根据指定大小 m，创建空的哈希表，并随机挑一个哈希函数
func NewHashTable() *HashTable {
	ht := &HashTable{
		Buckets: make([][]int64, m),
		Hash:    NewUniversalHash(),
	}
	return ht
}

// Insert 往哈希表里插入一个值 x
func (ht *HashTable) Insert(x int64) {
	idx := ht.Hash.HashFunc(x)
	ht.Buckets[idx] = append(ht.Buckets[idx], x)
}

// Exists 检查是否包含某个值 x
func (ht *HashTable) Exists(x int64) bool {
	idx := ht.Hash.HashFunc(x)
	// 在对应桶里线性搜索
	for _, val := range ht.Buckets[idx] {
		if val == x {
			return true
		}
	}
	return false
}

// PrintAll 打印所有桶的情况（演示用）
func (ht *HashTable) PrintAll() {
	for i, bucket := range ht.Buckets {
		if len(bucket) > 0 {
			fmt.Printf("Bucket[%d]: %v\n", i, bucket)
		}
	}
}

// ===================== 主函数测试 ====================== //

func main() {
	// 构建一个哈希表
	table := NewHashTable()

	// 插入一些数据
	data := []int64{10, 20, 30, 101, 202, 9999, 1234, 7890, 1001, 102, 50}
	for _, d := range data {
		table.Insert(d)
	}

	// 打印当前桶分布
	fmt.Println("Hash Table Buckets distribution:")
	table.PrintAll()
	fmt.Println()

	// 测试查找
	tests := []int64{20, 202, 9999, 777, 50, 101}
	for _, t := range tests {
		if table.Exists(t) {
			fmt.Printf("Value %d is found in Hash Table.\n", t)
		} else {
			fmt.Printf("Value %d is NOT found in Hash Table.\n", t)
		}
	}

	// 你可以多次运行程序，观察到每次选到的 (a,b) 不同，桶分布也可能不同
	// 这就说明了我们随机选了一个通用哈希函数
}
