// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/bits.go
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/search.go
// https://baobaobear.github.io/post/page/4/
// https://github.dev/old-yan/CP-template/blob/a07b6fe0092e9ee890a0e35ada6ea1bb2c83ba05/MATH/BitwiseHelper.h#L1

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	foo(0)
}

func foo(x int) {
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

	// 最低的len位取反
	flipLow := func(x, len int) int {
		return x ^ ((1 << len) - 1)
	}

	// 差集(A\B),从A中去掉B中的元素
	sub := func(A, B int) int { return A &^ B }

	// 全集 (0-n-1)
	all_ := func(n int) int { return ^(-1 << n) }

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

	// !对二的整数幂取模可以换成与运算
	a := 1200
	fmt.Println(a%32 == a&31)

	// https://noshi91.github.io/Library/other/popcount.cpp
	popCount64 := func(x uint64) int {
		x -= x >> 1 & 0x5555555555555555
		x = (x & 0x3333333333333333) + (x >> 2 & 0x3333333333333333)
		x = x + (x>>4)&0x0F0F0F0F0F0F0F0F
		return int(x * 0x0101010101010101 >> 56 & 0x7f)
	}
	fmt.Println(popCount64(13) == bits.OnesCount64(13))

	a = 5
	fmt.Println(1<<a - 1)   // 全1掩码
	fmt.Println(^(^0 << a)) // 全1(-1)左移a位,再取反

	// 利用位运算消除特定分支的方法:
	// !1.将if条件映射到 -1/0/1,先转变为加法/乘法
	// !2.-1/0/1的乘法相当于与运算,加法相当于或运算
	sum := 0
	flag := 1
	sum += flag & 128 // 消除分支: if(flag) sum += 128

	_ = []interface{}{
		lowbit,
		isSubset,
		isPow2,
		hasAdjacentOnes,
		hasAdjacentZeros,
		abs,
		tz,
		flipLow,
		all_,
		sub,
	}
}

// 掩码操作
func mask() {
	// !假定这里的位数为64

	// 生成全1掩码
	makeMask := func() int { return -1 }

	// 生成末尾有k个1的掩码
	makeBackOnes := func(k int) int {
		if k >= 64 {
			return -1
		}
		return (1 << k) - 1
	}

	// !生成某一段为1的掩码([l,r]闭区间)
	makeMaskRange := func(l, r int) int {
		return makeBackOnes(r+1) ^ makeBackOnes(l)
	}

	// 向上取整到2的幂
	getCeil := func(x int) int {
		if x == 0 {
			return 0
		}
		return 1 << (bits.Len(uint(x - 1)))
	}

	// 向下取整到2的幂
	getFloor := func(x int) int {
		if x == 0 {
			return 0
		}
		return 1 << (bits.Len(uint(x)) - 1)
	}

}
