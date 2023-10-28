package main

import "math/bits"

func main() {

}

// 在 [low,high] 区间内找两个数字 A B，使其异或值最大且不超过 limit
// 返回值保证 A <= B
// 复杂度 O(log(high))
// https://github.com/981377660LMT/codeforces-go/blob/50efc27004a0864fad32b2070b0f01c87b67b7c1/copypasta/bits.go#L650
func MaxXorWithLimit(low, high, limit int) (int, int) {
	n := bits.Len(uint(high ^ low))
	maxXor := 1<<n - 1
	mid := high&^maxXor | 1<<(n-1)
	if limit >= maxXor { // 无约束，相关题目 https://codeforces.com/problemset/problem/276/D
		return mid - 1, mid
	}
	if limit >= 1<<(n-1) { // A 和 B 能否在第 n-1 位不同的情况下，构造出一个满足要求的解？
		a, b := mid&(mid-1), mid
		for i := n - 2; i >= 0; i-- {
			bt := 1 << i
			if limit&bt > 0 { // a 取 1，b 取 0 总是优于 a 取 0，b 取 1
				a |= bt
			} else if high&(bt<<1-1) > ^low&(bt<<1-1) { // high 侧大，都取 1
				if high&bt == 0 { // b 没法取 1
					goto next
				}
				a |= bt
				b |= bt
			} else { // low 侧大，都取 0
				if low&bt > 0 { // a 没法取 0
					goto next
				}
			}
			if (a^low)&bt > 0 { // a 不受 low 的约束
				a |= limit & (bt - 1)
				break
			}
			if (b^high)&bt > 0 { // b 不受 high 的约束
				a |= bt - 1
				b |= ^limit & (bt - 1)
				break
			}
		}
		return a, b
	}
	// A 和 B 在第 n-1 位上必须相同
next:
	f := func(high int) (int, int) {
		n := bits.Len(uint(high ^ mid))
		maxXor := min(1<<n-1, limit)
		// 只有当 maxXor 为 0 时，返回值才必须相等
		if maxXor == 0 {
			return mid, mid
		}
		// maxXor 的最高位置于 B，其余置于 A
		mb := 1 << (bits.Len(uint(maxXor)) - 1)
		return mid | maxXor&^mb, mid | mb
	}
	if high-mid > mid-1-low { // 选区间长的一侧
		return f(high)
	}
	a, b := f(2*mid - 1 - low) // 对称到 high
	return 2*mid - 1 - b, 2*mid - 1 - a
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
