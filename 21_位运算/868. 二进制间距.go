// 868. 二进制间距
// !给定一个正整数 n，找到并返回 n 的二进制表示中两个 相邻 1 之间的 最长距离 。如果不存在两个相邻的 1，返回 0 。
// 位运算，O(log log n)
// https://leetcode.cn/problems/binary-gap/solutions/837523/wei-yun-suan-olog-log-n-by-hqztrue-aflh/
// 沈洋《回归本源——位运算及其应用》，OI国家集训队2014论文集 这里有更多的类似技巧，或者更早一点的某篇集训队文章。

package main

import (
	"fmt"
	"math/bits"
)

// x <= 2^31
func binaryGap(x int) int {
	x >>= bits.TrailingZeros(uint(x))
	if x == 1 {
		return 0
	}
	y, s := x, 0
	y |= y >> 1
	y |= y >> 2
	y |= y >> 4
	y |= y >> 8
	y |= y >> 16
	if x == y {
		return 1
	}

	if (x | (x >> 1)) != y {
		x |= x >> 1
		s += 1
	} else {
		goto l1
	}

	if (x | (x >> 2)) != y {
		x |= x >> 2
		s += 2
	} else {
		goto l2
	}

	if (x | (x >> 4)) != y {
		x |= x >> 4
		s += 4
	} else {
		goto l3
	}

	if (x | (x >> 8)) != y {
		x |= x >> 8
		s += 8
	} else {
		goto l4
	}

	if (x | (x >> 8)) != y {
		x |= x >> 8
		s += 8
	}
l4:
	if (x | (x >> 4)) != y {
		x |= x >> 4
		s += 4
	}
l3:
	if (x | (x >> 2)) != y {
		x |= x >> 2
		s += 2
	}
l2:
	if (x | (x >> 1)) != y {
		x |= x >> 1
		s += 1
	}
l1:
	return s + 2
}

func main() {
	// Example test case
	fmt.Println(binaryGap(22)) // Example output for the test case
}
