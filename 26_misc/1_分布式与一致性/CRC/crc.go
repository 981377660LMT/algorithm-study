package main

import (
	"fmt"
	"hash/crc32"
)

func main() {
	exampleCRC32Simple([]byte("123456789"))
	exampleCRC32Stream([][]byte{[]byte("123"), []byte("456"), []byte("789")})
}

func exampleCRC32Simple(data []byte) {
	// 使用 IEEE 多项式
	sum := crc32.ChecksumIEEE(data)
	fmt.Printf("CRC-32(IEEE) = 0x%08X\n", sum)
}

func exampleCRC32Stream(dataChunks [][]byte) {
	// 准备一个表(IEEE)
	tab := crc32.MakeTable(crc32.IEEE)
	// 创建一个流式哈希对象
	h := crc32.New(tab)

	// 多次写入
	for _, chunk := range dataChunks {
		h.Write(chunk)
	}

	// 最后获取校验
	sum := h.Sum32()
	fmt.Printf("CRC-32(stream) = 0x%08X\n", sum)
}
