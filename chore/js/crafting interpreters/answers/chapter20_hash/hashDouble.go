package main

import "math"

func hashDouble(value float64) uint32 {
	value += 1.0
	// 将float64转换为uint64位的二进制表示
	bits := math.Float64bits(value)
	// 提取低32位和高32位
	low := uint32(bits & 0xFFFFFFFF)
	high := uint32(bits >> 32)
	// 返回低32位和高32位的和作为哈希码
	return low + high
}
