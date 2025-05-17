// ZigZag编码/解码：
// 通过位运算将有符号整数映射为无符号整数，使小的负数也能被高效编码
// 提供了32位和64位两个版本

// 变长编码/解码 (类似于Protocol Buffers的Varint):
// 对数字按7位分组，每组用一个字节表示
// 字节的最高位(MSB)标记是否还有后续字节
// 小整数使用较少字节，大整数使用较多字节

// 组合编码/解码：
// 先ZigZag将有符号转为无符号，再变长编码
// 结合两种方式的优势，高效编码正负数

package main

import (
	"bytes"
	"fmt"
)

// ZigZag编码 - 将有符号整数映射为无符号整数
func EncodeZigZag32(n int32) uint32 {
	return uint32((n << 1) ^ (n >> 31))
}

// ZigZag解码 - 将ZigZag编码的无符号整数还原为有符号整数
func DecodeZigZag32(n uint32) int32 {
	return int32((n >> 1) ^ uint32((int32(n&1)<<31)>>31))
}

// ZigZag编码 - 64位版本
func EncodeZigZag64(n int64) uint64 {
	return uint64((n << 1) ^ (n >> 63))
}

// ZigZag解码 - 64位版本
func DecodeZigZag64(n uint64) int64 {
	return int64((n >> 1) ^ uint64((int64(n&1)<<63)>>63))
}

// 变长编码 - 将整数编码为变长字节序列
func EncodeVarint(x uint64) []byte {
	var buf bytes.Buffer
	for {
		b := byte(x & 0x7F) // 取低7位
		x >>= 7             // 右移7位
		if x != 0 {
			b |= 0x80 // 如果后面还有数据，设置最高位为1
		}
		buf.WriteByte(b)
		if x == 0 {
			break
		}
	}
	return buf.Bytes()
}

// 变长解码 - 将变长字节序列解码为整数
func DecodeVarint(data []byte) (uint64, int) {
	var result uint64 = 0
	var shift uint = 0

	for i, b := range data {
		result |= uint64(b&0x7F) << shift
		if b < 0x80 {
			return result, i + 1 // 返回值和已读取的字节数
		}
		shift += 7

		// 超过64位，无效
		if shift >= 64 {
			return 0, 0
		}
	}

	// 数据不完整
	return 0, 0
}

// 结合ZigZag和变长编码处理有符号整数
func EncodeSigned32(n int32) []byte {
	zigzag := EncodeZigZag32(n)
	return EncodeVarint(uint64(zigzag))
}

func DecodeSigned32(data []byte) (int32, int) {
	v, n := DecodeVarint(data)
	if n == 0 {
		return 0, 0
	}
	return DecodeZigZag32(uint32(v)), n
}

func main() {
	// 测试ZigZag编码
	values := []int32{0, 1, -1, 2, -2, 127, -128, 32767, -32768}

	fmt.Println("===== ZigZag编码测试 =====")
	for _, v := range values {
		encoded := EncodeZigZag32(v)
		decoded := DecodeZigZag32(encoded)
		fmt.Printf("原始: %6d -> 编码: %6d -> 解码: %6d\n", v, encoded, decoded)
	}

	fmt.Println("\n===== 变长编码测试 =====")
	for i := uint64(0); i < 300; i += 50 {
		encoded := EncodeVarint(i)
		decoded, _ := DecodeVarint(encoded)
		fmt.Printf("原始: %3d -> 编码: % X (%d字节) -> 解码: %3d\n",
			i, encoded, len(encoded), decoded)
	}

	fmt.Println("\n===== ZigZag+变长编码测试 =====")
	for _, v := range values {
		encoded := EncodeSigned32(v)
		decoded, _ := DecodeSigned32(encoded)
		fmt.Printf("原始: %6d -> 编码: % X (%d字节) -> 解码: %6d\n",
			v, encoded, len(encoded), decoded)
	}
}
