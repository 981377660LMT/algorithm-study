// base64的转换/base64的加密与解密
// 用于压缩源码/压缩打表数据

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	b := NewBase64()
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s := b.Encode(nums)
	fmt.Println(s)
	fmt.Println(b.Decode(s))
}

const _BASE string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

type Base64 struct {
}

func NewBase64() *Base64 {
	return &Base64{}
}

func (base64 *Base64) Encode(nums []int) string {
	x, y := nums[0], nums[0]
	for _, z := range nums {
		x = max(x, z)
		y = min(y, z)
	}
	n := len(nums)
	tmp := 64
	if y >= 0 {
		tmp = bits.Len(uint(x))
	}
	b := max(6, tmp)
	sb := make([]int, (b*n+11)/6)
	sb[0] = b
	for i := 0; i < n; i++ {
		for j := 0; j < b; j++ {
			if nums[i]>>j&1 == 1 {
				sb[(i*b+j)/6+1] |= 1 << ((i*b + j) % 6)
			}
		}
	}

	res := make([]byte, len(sb))
	for i := 0; i < len(sb); i++ {
		res[i] = _BASE[sb[i]]
	}
	return string(res)
}

func (base64 *Base64) Decode(s string) []int {
	sb := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		sb[i] = base64._ibase(s[i])
	}
	b := sb[0]
	m := len(sb) - 1
	nums := make([]int, (6*m)/b)
	for i := 0; i < m; i++ {
		for j := 0; j < 6; j++ {
			if sb[i+1]>>j&1 == 1 {
				nums[(i*6+j)/b] |= 1 << ((i*6 + j) % b)
			}
		}
	}
	return nums
}

func (base64 *Base64) _ibase(c byte) int {
	if c >= 'a' {
		return 0x1A + int(c) - int('a')
	}
	if c >= 'A' {
		return 0x00 + int(c) - int('A')
	}
	if c >= '0' {
		return 0x34 + int(c) - int('0')
	}
	if c == '+' {
		return 0x3E
	}
	if c == '/' {
		return 0x3F
	}
	return 0x40
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
