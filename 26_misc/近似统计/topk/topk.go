// https://github.com/dgryski/go-topk/blob/master/topk.go

// Package topk implements the Filtered Space-Saving TopK streaming algorithm
/*

The original Space-Saving algorithm:
https://icmi.cs.ucsb.edu/research/tech_reports/reports/2005-23.pdf

The Filtered Space-Saving enhancement:
http://www.l2f.inesc-id.pt/~fmmb/wiki/uploads/Work/misnis.ref0a.pdf

This implementation follows the algorithm of the FSS paper, but not the
suggested implementation.  Specifically, we use a heap instead of a sorted list
of monitored items, and since we are also using a map to provide O(1) access on
update also don't need the c_i counters in the hash table.

Licensed under the MIT license.

*/

// !使用 Filtered Space-Saving（FSS） 算法（也可视为一种 Space-Saving 变体）来做 Top-K 频繁项统计
// 通过一个额外的计数器结构（这里 alphas[]）来记住“暂时没有进入监控”的元素的累积频次或曾经被踢出的元素的频次。
// 当累积到一定程度超过了堆里最小频度，就把它替换进去。
// 反之，如果数次出现后还是无法超越最小频度，就一直停留在外部计数，不会占用主结构空间。
// !这样的算法非常适合流式或者在线统计，尤其当我们只关心“哪些元素出现次数最多”而不需要精确计数的场合（如日志分析、热门词排名等）

package main

import (
	"bytes"
	"container/heap"
	"encoding/gob"
	"fmt"
	"math/bits"
	"sort"
	"unsafe"
)

func main() {
	// 1. 创建一个 Stream，追踪前 5 个最频繁元素
	tk := NewStream(5)

	// 2. 不断插入数据
	//    模拟处理一批数据（如日志中的 IP，或点击流中的 item 等）
	data := []string{"apple", "banana", "apple", "orange", "banana", "banana", "melon", "apple"}

	for _, item := range data {
		tk.Insert(item, 1) // 每出现一次，就插入计数 1
	}

	// 3. 获取并打印当前最常见的元素（前 5）
	topElements := tk.Keys()
	fmt.Println("Top elements:")
	for _, e := range topElements {
		fmt.Printf("Key=%s  Count=%d  Error=%d\n", e.Key, e.Count, e.Error)
	}

	// 也可以单独查询某个 key 的估计值
	appleEst := tk.Estimate("apple")
	fmt.Printf("Estimate for 'apple': Count=%d  Error=%d\n", appleEst.Count, appleEst.Error)
}

// Element is a TopK item
type Element struct {
	Key   string // 元素的标识（这里用字符串）
	Count int    // 算法估计出的该元素的当前计数
	Error int    // 元素真正加入监控之前的最低可能计数
}

type elementsByCountDescending []Element

func (elts elementsByCountDescending) Len() int { return len(elts) }
func (elts elementsByCountDescending) Less(i, j int) bool {
	return (elts[i].Count > elts[j].Count) || (elts[i].Count == elts[j].Count && elts[i].Key < elts[j].Key)
}
func (elts elementsByCountDescending) Swap(i, j int) { elts[i], elts[j] = elts[j], elts[i] }

// 可删除堆.
type keys struct {
	m    map[string]int // 元素Key => 元素在堆中的索引
	elts []Element      // 最小堆
}

// Implement the container/heap interface
func (tk *keys) Len() int { return len(tk.elts) }
func (tk *keys) Less(i, j int) bool {
	return (tk.elts[i].Count < tk.elts[j].Count) || (tk.elts[i].Count == tk.elts[j].Count && tk.elts[i].Error > tk.elts[j].Error)
}
func (tk *keys) Swap(i, j int) {
	tk.elts[i], tk.elts[j] = tk.elts[j], tk.elts[i]
	tk.m[tk.elts[i].Key] = i
	tk.m[tk.elts[j].Key] = j
}
func (tk *keys) Push(x interface{}) {
	e := x.(Element)
	tk.m[e.Key] = len(tk.elts)
	tk.elts = append(tk.elts, e)
}
func (tk *keys) Pop() interface{} {
	var e Element
	e, tk.elts = tk.elts[len(tk.elts)-1], tk.elts[:len(tk.elts)-1]
	delete(tk.m, e.Key)
	return e
}

// Stream calculates the TopK elements for a stream
type Stream struct {
	n int  // 期望追踪的前 n 大元素
	k keys // 可删除堆用于存放当前在追踪的元素

	// 来记住“暂时没有进入监控”的元素的累积频次或曾经被踢出的元素的频次。当累积到一定程度超过了堆里最小频度，就把它替换进去。
	alphas []int
}

// NewStream returns a Stream estimating the top n most frequent elements
func NewStream(n int) *Stream {
	return &Stream{
		n:      n,
		k:      keys{m: make(map[string]int), elts: make([]Element, 0, n)},
		alphas: make([]int, n*6), // 6 is the multiplicative constant from the paper
	}
}

// 把 64位哈希值低 32 位乘以 n，然后右移 32 位，得到 [0, n) 区间的整数。
func reduce(x uint64, n int) uint32 {
	return uint32(uint64(uint32(x)) * uint64(n) >> 32)
}

// Insert adds an element to the stream to be tracked
// It returns an estimation for the just inserted element
func (s *Stream) Insert(x string, count int) Element {
	xhash := reduce(Sum64Str(0, 0, x), len(s.alphas))

	// are we tracking this element?
	if idx, ok := s.k.m[x]; ok {
		s.k.elts[idx].Count += count
		e := s.k.elts[idx]
		heap.Fix(&s.k, idx)
		return e
	}

	// can we track more elements?
	if len(s.k.elts) < s.n {
		// there is free space
		e := Element{Key: x, Count: count}
		heap.Push(&s.k, e)
		return e
	}

	if s.alphas[xhash]+count < s.k.elts[0].Count {
		e := Element{
			Key:   x,
			Error: s.alphas[xhash],
			Count: s.alphas[xhash] + count,
		}
		s.alphas[xhash] += count
		return e
	}

	// !只有当 (alphas[xhash]+count) >= 堆顶的最小计数时，才会真正进入主堆。
	// 否则就一直累加在 alphas[xhash]，相当于被“过滤”在外。

	// replace the current minimum element
	minKey := s.k.elts[0].Key

	mkhash := reduce(Sum64Str(0, 0, minKey), len(s.alphas))
	s.alphas[mkhash] = s.k.elts[0].Count

	e := Element{
		Key:   x,
		Error: s.alphas[xhash],
		Count: s.alphas[xhash] + count,
	}
	s.k.elts[0] = e

	// we're not longer monitoring minKey
	delete(s.k.m, minKey)
	// but 'x' is as array position 0
	s.k.m[x] = 0

	heap.Fix(&s.k, 0)
	return e
}

// 返回当前被正式追踪的所有元素（在堆里），并按 Count 降序排序。
// 这是我们最终所说的“Top-K 候选”。
func (s *Stream) Keys() []Element {
	elts := append([]Element(nil), s.k.elts...)
	sort.Sort(elementsByCountDescending(elts))
	return elts
}

// 查询指定元素的近似频率.
func (s *Stream) Estimate(x string) Element {
	xhash := reduce(Sum64Str(0, 0, x), len(s.alphas))

	// are we tracking this element?
	if idx, ok := s.k.m[x]; ok {
		e := s.k.elts[idx]
		return e
	}
	count := s.alphas[xhash]
	e := Element{
		Key:   x,
		Error: count,
		Count: count,
	}
	return e
}

func (s *Stream) GobEncode() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(s.n); err != nil {
		return nil, err
	}
	if err := enc.Encode(s.k.m); err != nil {
		return nil, err
	}
	if err := enc.Encode(s.k.elts); err != nil {
		return nil, err
	}
	if err := enc.Encode(s.alphas); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Stream) GobDecode(b []byte) error {
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	if err := dec.Decode(&s.n); err != nil {
		return err
	}
	if err := dec.Decode(&s.k.m); err != nil {
		return err
	}
	if err := dec.Decode(&s.k.elts); err != nil {
		return err
	}
	if err := dec.Decode(&s.alphas); err != nil {
		return err
	}
	return nil
}

// #region sip13

func Sum64(k0, k1 uint64, b []byte) uint64 {
	return Sum64Str(k0, k1, unsafeB2S(b))
}

func readUint64(b string) uint64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func unsafeB2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Sum64Str(k0, k1 uint64, p string) uint64 {
	var v0, v1, v2, v3 uint64

	v0 = k0 ^ 0x736f6d6570736575
	v1 = k1 ^ 0x646f72616e646f6d
	v2 = k0 ^ 0x6c7967656e657261
	v3 = k1 ^ 0x7465646279746573
	b := uint64(len(p)) << 56

	for len(p) >= 8 {
		m := readUint64(p)
		v3 ^= m
		v0 += v1
		v2 += v3
		v1 = bits.RotateLeft64(v1, 13)
		v3 = bits.RotateLeft64(v3, 16)
		v1 ^= v0
		v3 ^= v2
		v0 = bits.RotateLeft64(v0, 32)
		v2 += v1
		v0 += v3
		v1 = bits.RotateLeft64(v1, 17)
		v3 = bits.RotateLeft64(v3, 21)
		v1 ^= v2
		v3 ^= v0
		v2 = bits.RotateLeft64(v2, 32)
		v0 ^= m
		p = p[8:]
	}

	for _, c := range []byte(p) {
		b = bits.RotateLeft64(b|uint64(c), 56)
	}
	m := bits.RotateLeft64(b, len(p)*8)

	// last block with finalization
	v3 ^= m
	f := uint64(0xff)
	for i := 0; i < 4; i++ {
		v0 += v1
		v2 += v3
		v1 = bits.RotateLeft64(v1, 13)
		v3 = bits.RotateLeft64(v3, 16)
		v1 ^= v0
		v3 ^= v2
		v0 = bits.RotateLeft64(v0, 32)
		v2 += v1
		v0 += v3
		v1 = bits.RotateLeft64(v1, 17)
		v3 = bits.RotateLeft64(v3, 21)
		v1 ^= v2
		v3 ^= v0
		v2 = bits.RotateLeft64(v2, 32)
		v2 ^= f
		v0 ^= m
		m, f = 0, 0 // clear last block and finalization mixins
	}

	return v0 ^ v1 ^ v2 ^ v3
}

// #endregion
