package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	agc015_d()
}

func demo() {
	fmt.Println(decompose(5))
	fmt.Println(compose([]int{1, 0, 0}))
	fmt.Println(IntervalBitwiseOrExpandGroup(1, 9))
	fmt.Println(IntervalBitwiseOrExpandGroup(1, 5))
}

// https://atcoder.jp/contests/agc015/tasks/agc015_d
func agc015_d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var left, right int
	fmt.Fscan(in, &left, &right)

	fmt.Fprintln(out, IntervalBitwiseOrExpandGroup(left, right))
}

// 区间按位或得到的数的个数.
func IntervalBitwiseOrExpandGroup(left, right int) int {
	if left > right {
		return 0
	}
	if left == right {
		return 1
	}

	bitsOfA := decompose(left)
	bitsOfB := decompose(right)
	since := len(bitsOfA) - 1
	for bitsOfA[since] == bitsOfB[since] {
		bitsOfA[since] = 0
		bitsOfB[since] = 0
		since--
	}
	bitsOfB[since] = 0
	bBits := since
	for bBits > 0 && bitsOfB[bBits-1] == 0 {
		bBits--
	}

	left = compose(bitsOfA)
	if 1<<bBits >= left {
		return (1<<since)*2 - left
	}

	res := (1 << since) - left
	res += 1 << bBits
	res += (1 << since) - left
	return res
}

func decompose(n int) []int {
	res := make([]int, 64)
	for b := 0; n > 0; b++ {
		res[b] = n & 1
		n >>= 1
	}
	return res
}

func compose(bits []int) int {
	res := 0
	for i := len(bits) - 1; i >= 0; i-- {
		res = res<<1 | bits[i]
	}
	return res
}
