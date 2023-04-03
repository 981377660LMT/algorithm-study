// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/bits.go
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/search.go
// https://baobaobear.github.io/post/page/4/

package main

import "math/bits"

func _(x int) {
	// !利用 -v = ^v+1  负数等于反码加1
	lowbit := func(v int64) int64 { return v & -v } // (如果要找最低位的0,先将v取反)
	x = 0b10100                                     // 取反 => ^x == 0b01011
	x = ^x + 1                                      // 补码 => comp == 0b01100

	// 最低位的 1 变 0
	x &= x - 1
	// 最低位连续的所有 1 变 0
	x &= x + 1

	// 最低位的 0 变 1
	x |= x + 1
	// 最低位连续的所有 0 变 1
	x |= x - 1

	// 取右边连续的 1
	x ^= x + 1

	// -1 表示为 -1 = 0b1111...1111
	x &= (^0 << 3)    // 清除从第三位开始右边的所有1 (截断下界)
	x |= (^0 << 3)    // 清除从第三位开始右边的所有0 (截断下界)
	x &= (1 << 3) - 1 // 清除从第三位开始左边的所有1 (截断上界)
	x |= (1 << 3) - 1 // 清除从第三位开始左边的所有0 (截断上界)

	// 有符号整数计算绝对值	( n ^ (n >> 31) ) - (n >> 31)
	abs := func(n int) int { return (n ^ (n >> 31)) - (n >> 31) }

	// x 是 y 的子集
	isSubset := func(x, y int) bool { return x|y == y } // x 和 y 的并集是 y
	isSubset = func(x, y int) bool { return x&y == x }  // x 和 y 的交集是 x

	// 1,2,4,8,...
	isPow2 := func(v int64) bool { return v > 0 && v&(v-1) == 0 }

	// 是否有两个相邻的 1    有 https://oeis.org/A004780 没有 https://oeis.org/A003714
	hasAdjacentOnes := func(v uint) bool { return v>>1&v > 0 }

	// 是否有两个相邻的 0（不考虑前导零）    有 https://oeis.org/A004753 没有 http://oeis.org/A003754
	hasAdjacentZeros := func(v uint) bool {
		v |= v >> 1 // 若没有相邻的 0，则 v 会变成全 1 的数
		return v&(v+1) > 0
	}

	// 快速计算尾随0 trialingZeros: lowbit+bit_length
	tz := func(v int64) int64 { return int64(bits.Len(uint(v&-v)) - 1) }

	// !32位时, clz32 + bit_length = 32

	// !python里用int模仿uint64的行为:	x &((1<<64)-1) 即可

	_ = []interface{}{
		lowbit,
		isSubset,
		isPow2,
		hasAdjacentOnes,
		hasAdjacentZeros,
		abs,
		tz,
	}
}
