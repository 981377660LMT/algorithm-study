// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/bits.go
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/search.go

package main

func _(x int) {
	// 利用 -v = ^v+1
	lowbit := func(v int64) int64 { return v & -v }

	// 最低位的 1 变 0
	x &= x - 1

	// 最低位的 0 变 1
	x |= x + 1

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

	_ = []interface{}{
		lowbit,
		isSubset,
		isPow2,
		hasAdjacentOnes,
		hasAdjacentZeros,
	}
}
