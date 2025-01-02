// HyperLogLog 中 用于稀疏模式下保存哈希值（编码后）的重要部分。
// 通过差分可变长编码 (variable-length encoding) 的方式压缩存储大量 32-bit 整数。
//
// compressedList 的本质是：
//
// - 维护一个递增序列（通常是已排序的 32-bit 整数，且不会乱序插入）；
// - 在插入每个新数 x 时，存储 delta = x - last 的可变长编码到 b 里，并更新 last = x；
// - count 跟踪插入了多少个数。

package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
)

func main() {
	// 1) 创建一个新的 compressedList，预留容量 10
	cl := NewCompressedList(10)

	// 2) 往里面追加若干个 32-bit 整数
	//    这里的使用场景通常是: 事先将它们排序好，然后单调插入
	nums := []uint32{5, 6, 6, 10, 100, 105, 106, 107}
	for _, v := range nums {
		cl.Append(v)
	}

	// 3) 查看当前 compressedList 的元信息
	fmt.Printf("Appended %d items, stored bytes = %d\n", cl.count, cl.Len())

	// 4) 迭代读取看看
	it := cl.Iter()
	fmt.Printf("Reading items in compressedList:\n")
	for it.HasNext() {
		val := it.Next()
		fmt.Printf("  Next = %d\n", val)
	}

	// 5) 查看序列化数据
	data, err := cl.MarshalBinary()
	if err != nil {
		fmt.Println("Marshal error:", err)
		return
	}
	fmt.Printf("Serialized to %d bytes\n", len(data))

	// 6) 反序列化到一个新的 compressedList
	var newCL compressedList
	if err := newCL.UnmarshalBinary(data); err != nil {
		fmt.Println("Unmarshal error:", err)
		return
	}

	// 7) 验证解码后内容一致
	fmt.Printf("After Unmarshal: newCL.count = %d, newCL.len() = %d\n",
		newCL.count, len(newCL.b))

	// 迭代阅读
	it2 := newCL.Iter()
	var recovered []uint32
	for it2.HasNext() {
		recovered = append(recovered, it2.Next())
	}
	fmt.Printf("Recovered items: %v\n", recovered)

	// (可选) 我们可以演示将新的一批数据做归并插入:
	// 假设再插入一批更大的数字 [200, 205, 210], 先要保证非递减才能append
	extra := []uint32{200, 205, 210}
	// compressedList 通常假设数据是有序的, 否则要自己先sort
	sort.Slice(extra, func(i, j int) bool { return extra[i] < extra[j] })

	for _, v := range extra {
		newCL.Append(v)
	}
	fmt.Printf("Appended extra items. newCL.count = %d\n", newCL.count)

	// 再迭代查看
	it3 := newCL.Iter()
	fmt.Println("All items now (including extras):")
	for it3.HasNext() {
		fmt.Printf("  %d\n", it3.Next())
	}
}

// ErrorTooShort is an error that UnmarshalBinary try to parse too short binary.
var ErrorTooShort = errors.New("too short binary")

// compressedList 则是一个有序存储的结构，里面用“前缀差分 + 可变长编码”来节省空间.
type compressedList struct {
	count uint32 // 当前存储的元素数量（逻辑上插入了多少个数）
	last  uint32 // 最后插入的元素值（未差分前的原始值），帮助在差分编码/解码时进行累加
	b     variableLengthList
}

func NewCompressedList(capacity int) *compressedList {
	v := &compressedList{}
	v.b = make(variableLengthList, 0, capacity)
	return v
}

func (v *compressedList) Clone() *compressedList {
	if v == nil {
		return nil
	}

	newV := &compressedList{
		count: v.count,
		last:  v.last,
	}

	newV.b = make(variableLengthList, len(v.b))
	copy(newV.b, v.b)
	return newV
}

func (v *compressedList) MarshalBinary() (data []byte, err error) {
	// Marshal the variableLengthList
	bdata, err := v.b.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// At least 4 bytes for the two fixed sized values plus the size of bdata.
	data = make([]byte, 0, 4+4+len(bdata))

	// Marshal the count and last values.
	data = append(data, []byte{
		// Number of items in the list.
		byte(v.count >> 24),
		byte(v.count >> 16),
		byte(v.count >> 8),
		byte(v.count),
		// The last item in the list.
		byte(v.last >> 24),
		byte(v.last >> 16),
		byte(v.last >> 8),
		byte(v.last),
	}...)

	// Append the list
	return append(data, bdata...), nil
}

func (v *compressedList) Len() int {
	return len(v.b)
}

func (v *compressedList) decode(i int, last uint32) (uint32, int) {
	n, i := v.b.decode(i, last)
	return n + last, i
}

func (v *compressedList) Append(x uint32) {
	v.count++
	v.b = v.b.Append(x - v.last)
	v.last = x
}

func (v *compressedList) Iter() *iterator {
	return &iterator{0, 0, v}
}

type variableLengthList []uint8

func (v variableLengthList) MarshalBinary() (data []byte, err error) {
	// 4 bytes for the size of the list, and a byte for each element in the
	// list.
	data = make([]byte, 0, 4+v.Len())

	// Length of the list. We only need 32 bits because the size of the set
	// couldn't exceed that on 32 bit architectures.
	sz := v.Len()
	data = append(data, []byte{
		byte(sz >> 24),
		byte(sz >> 16),
		byte(sz >> 8),
		byte(sz),
	}...)

	// Marshal each element in the list.
	for i := 0; i < sz; i++ {
		data = append(data, v[i])
	}

	return data, nil
}

func (v *compressedList) UnmarshalBinary(data []byte) error {
	if len(data) < 12 {
		return ErrorTooShort
	}

	// Set the count.
	v.count, data = binary.BigEndian.Uint32(data[:4]), data[4:]

	// Set the last value.
	v.last, data = binary.BigEndian.Uint32(data[:4]), data[4:]

	// Set the list.
	sz, data := binary.BigEndian.Uint32(data[:4]), data[4:]
	v.b = make([]uint8, sz)
	if uint32(len(data)) < sz {
		return ErrorTooShort
	}
	for i := uint32(0); i < sz; i++ {
		v.b[i] = data[i]
	}
	return nil
}

func (v variableLengthList) Len() int {
	return len(v)
}

func (v *variableLengthList) Iter() *iterator {
	return &iterator{0, 0, v}
}

// 返回 (新元素值, 下一次读取起点)。
func (v variableLengthList) decode(i int, last uint32) (uint32, int) {
	// 从 v[i] 开始，读连续字节，每个字节低 7 位拼起来，高位 bit=1 表示下一字节
	var x uint32
	j := i
	for ; v[j]&0x80 != 0; j++ {
		x |= uint32(v[j]&0x7f) << (uint(j-i) * 7)
	}
	// 处理最后一个字节
	x |= uint32(v[j]) << (uint(j-i) * 7)

	return x, j + 1
}

// 若 x <= 127, 只写 1 字节；
// 若大于 127, 则拆分为多个 7-bit 块，每个块最高位 0x80 表示“还有后续字节”。直到最后一个块把最高位清零 (0x80 位置为 0)。
func (v variableLengthList) Append(x uint32) variableLengthList {
	// 当 x 超过 7 bit，就在高位继续分割
	// 最后以一个 0x7F 范围内的字节结束
	for x&0xffffff80 != 0 {
		v = append(v, uint8((x&0x7f)|0x80))
		x >>= 7
	}
	// 剩余的低 7 bit
	return append(v, uint8(x&0x7f))
}

// Original author of this file is github.com/clarkduvall/hyperloglog
type iterable interface {
	decode(i int, last uint32) (uint32, int)
	Len() int
	Iter() *iterator
}

type iterator struct {
	i    int      // 当前在 variableLengthList 的索引
	last uint32   // 用于累计恢复原值
	v    iterable // 这里 v = compressedList
}

func (iter *iterator) Next() uint32 {
	n, i := iter.v.decode(iter.i, iter.last)
	iter.last = n
	iter.i = i
	return n
}

func (iter *iterator) Peek() uint32 {
	n, _ := iter.v.decode(iter.i, iter.last)
	return n
}

func (iter iterator) HasNext() bool {
	return iter.i < iter.v.Len()
}
