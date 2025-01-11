package main

import (
	"container/heap"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// HashFunction 定义哈希函数类型，返回 [0,1) 的浮点数
type HashFunction func(string) float64

// NewSHA256HashFunction 返回一个基于 SHA256 的哈希函数
func NewSHA256HashFunction() HashFunction {
	return func(s string) float64 {
		hash := sha256.Sum256([]byte(s))
		// 取前8字节转换为 uint64
		u := binary.BigEndian.Uint64(hash[:8])
		// 将 uint64 映射到 [0,1)
		return float64(u) / float64(math.MaxUint64)
	}
}

// MinHeap 实现一个最大堆，用于维护最小 k 个哈希值
type MinHeap []float64

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] > h[j] } // 最大堆
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push 添加元素到堆中
func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(float64))
}

// Pop 从堆中移除并返回最大元素
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// KMV 结构体
type KMV struct {
	k           int
	hashFunc    HashFunction
	minHeap     *MinHeap
	initialized bool
}

// NewKMV 创建一个新的 KMV 实例
func NewKMV(k int, hashFunc HashFunction) *KMV {
	h := &MinHeap{}
	heap.Init(h)
	return &KMV{
		k:        k,
		hashFunc: hashFunc,
		minHeap:  h,
	}
}

// Insert 插入一个元素，维护最小 k 个哈希值
func (m *KMV) Insert(element string) {
	hashVal := m.hashFunc(element)
	if m.minHeap.Len() < m.k {
		heap.Push(m.minHeap, hashVal)
	} else if hashVal < (*m.minHeap)[0] {
		heap.Pop(m.minHeap)
		heap.Push(m.minHeap, hashVal)
	}
}

// EstimateCardinality 估计基数
func (m *KMV) EstimateCardinality() float64 {
	if m.minHeap.Len() == 0 {
		return 0
	}
	// 取堆顶，即第 k 小的哈希值
	v_k := (*m.minHeap)[0]
	// 估计公式 n ≈ k / v_k
	return float64(m.k) / v_k
}

// Merge 合并另一个 KMV 到当前 KMV
func (m *KMV) Merge(other *KMV) {
	for _, hashVal := range *other.minHeap {
		m.InsertByValue(hashVal)
	}
}

// InsertByValue 插入一个已知的哈希值，用于合并操作
func (m *KMV) InsertByValue(hashVal float64) {
	if m.minHeap.Len() < m.k {
		heap.Push(m.minHeap, hashVal)
	} else if hashVal < (*m.minHeap)[0] {
		heap.Pop(m.minHeap)
		heap.Push(m.minHeap, hashVal)
	}
}

// ------------------------ 测试与示例 ------------------------ //

func main() {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 创建哈希函数
	hashFunc := NewSHA256HashFunction()

	// 设置 k 值
	k := 100

	// 创建 KMV 实例
	kmv := NewKMV(k, hashFunc)

	// 模拟插入数据
	n := 1000000 // 实际基数
	for i := 0; i < n; i++ {
		element := fmt.Sprintf("element_%d", i)
		kmv.Insert(element)
	}

	// 估计基数
	estimated := kmv.EstimateCardinality()
	fmt.Printf("Actual Cardinality: %d\n", n)
	fmt.Printf("Estimated Cardinality: %.0f\n", estimated)
	fmt.Printf("Error: %.2f%%\n", math.Abs(float64(n)-estimated)/float64(n)*100)

	// 示例合并两个 KMV 实例
	kmv1 := NewKMV(k, hashFunc)
	kmv2 := NewKMV(k, hashFunc)

	// 插入部分数据到 kmv1
	for i := 0; i < n/2; i++ {
		element := fmt.Sprintf("element_%d", i)
		kmv1.Insert(element)
	}

	// 插入剩余数据到 kmv2
	for i := n / 2; i < n; i++ {
		element := fmt.Sprintf("element_%d", i)
		kmv2.Insert(element)
	}

	// 合并 kmv2 到 kmv1
	kmv1.Merge(kmv2)

	// 估计合并后的基数
	mergedEstimated := kmv1.EstimateCardinality()
	fmt.Printf("\nAfter Merging:\n")
	fmt.Printf("Actual Cardinality: %d\n", n)
	fmt.Printf("Estimated Cardinality: %.0f\n", mergedEstimated)
	fmt.Printf("Error: %.2f%%\n", math.Abs(float64(n)-mergedEstimated)/float64(n)*100)
}
