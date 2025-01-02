// Go 在标准库 encoding/binary 包中提供了对 Varint 的支持

package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	var buf [binary.MaxVarintLen64]byte

	// 1. 将 123456 这个 int64 用变长编码写入 buf
	n := binary.PutVarint(buf[:], 123456)
	fmt.Println("Encoded bytes:", buf[:n]) // 输出写入的具体字节

	// 2. 解码：从 buf[:n] 中解析出 Varint，得到原始数值
	x, readBytes := binary.Varint(buf[:n])
	fmt.Println("Decoded value:", x, "Read bytes:", readBytes)
}

// 手写一个最简 Varint (无符号) 编码器
// 这段代码没有处理溢出和异常，仅用于教学
// EncodeUvarint returns the varint-encoded bytes of x (uint64).
func EncodeUvarint(x uint64) []byte {
	var buf []byte
	for {
		b := byte(x & 0x7F) // 取 7 位
		x >>= 7
		if x != 0 {
			b |= 0x80 // 最高位设为1，表示后面还有
			buf = append(buf, b)
		} else {
			// 最后一次：最高位为0，并退出
			buf = append(buf, b)
			break
		}
	}
	return buf
}

// DecodeUvarint 解出一个 uint64 值，以及实际使用的字节数
func DecodeUvarint(buf []byte) (uint64, int) {
	var x uint64
	var s uint
	for i, b := range buf {
		// 取出低 7 位，左移 s
		x |= uint64(b&0x7F) << s
		s += 7
		// 若最高位 == 0，说明结束
		if (b & 0x80) == 0 {
			return x, i + 1
		}
	}
	// 若循环完还没遇到最高位=0，说明数据不完整或错误
	return 0, 0
}
