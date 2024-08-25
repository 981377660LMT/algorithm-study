// FastDigitSum
// 计算一个数字各位digit之和
// !不如直接计算

package main

import (
	"fmt"
	"time"
)

func main() {
	d := NewDigitSum(-1)
	time1 := time.Now()
	for i := 0; i < 1e8; i++ {
		d.Sum(i)
	}
	time2 := time.Now()
	fmt.Println(time2.Sub(time1)) // 1.58807s

	time1 = time.Now()
	for i := 0; i < 1e8; i++ {
		DigitSumNaive(i)
	}
	time2 = time.Now()
	fmt.Println(time2.Sub(time1)) // 734.4214ms
}

func DigitSumNaive(x int) int {
	res := 0
	for x > 0 {
		res += x % 10
		x /= 10
	}
	return res
}

type DigitSum struct {
	mod int
	dp  []int // 长为10^step的数组, dp[x]表示x的各位数字之和
}

// step = -1 表示使用默认值 5.
func NewDigitSum(step int) *DigitSum {
	if step == -1 {
		step = 5
	}

	if step < 4 {
		step = 4
	}
	if step > 8 {
		step = 8
	}

	mod := 1
	for i := 0; i < step; i++ {
		mod *= 10
	}

	d := &DigitSum{mod: mod, dp: make([]int, mod)}
	for x := 1; x < mod; x++ {
		d.dp[x] = d.dp[x/10] + x%10
	}
	return d
}

func (d *DigitSum) Sum(x int) int {
	res := 0
	for x > 0 {
		res += d.dp[x%d.mod]
		x /= d.mod
	}
	return res
}
