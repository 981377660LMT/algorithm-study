// 求小于等于n的哈沙德数(Harshad number)的个数
// harshad number: n能被n的各位数码之和整除
// n<=1e14
// https://atcoder.jp/contests/abc336/editorial/9055
// !需要枚举数码之和 1~9*m 才能确定当前数模 digitSum 的模数

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	// !list(map(int, str(n))
	nums := []int{}
	for n > 0 {
		nums = append(nums, n%10)
		n /= 10
	}
	for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
	m := len(nums)

	// !各位数之和为 targetDigitSum 时, [1, n] 中符合条件的数的个数
	cal := func(targetDigitSum int) int {
		memo := [15][2][2][150][150]int{} // pos, hasLeadingZero, isLimit, digitSum, curMod
		for i1 := 0; i1 < 15; i1++ {
			for i2 := 0; i2 < 2; i2++ {
				for i3 := 0; i3 < 2; i3++ {
					for i4 := 0; i4 < 150; i4++ {
						for i5 := 0; i5 < 150; i5++ {
							memo[i1][i2][i3][i4][i5] = -1
						}
					}
				}
			}
		}

		var dfs func(pos int, hasLeadingZero bool, isLimit bool, digitSum int, curMod int) int
		dfs = func(pos int, hasLeadingZero bool, isLimit bool, digitSum int, curMod int) int {
			if pos == m {
				return BoolToInt(!hasLeadingZero && curMod == 0 && digitSum == targetDigitSum)
			}

			ptr := &memo[pos][BoolToInt(hasLeadingZero)][BoolToInt(isLimit)][digitSum][curMod]
			if *ptr != -1 {
				return *ptr
			}

			res := 0
			up := 9
			if isLimit {
				up = nums[pos]
			}
			for cur := 0; cur <= up; cur++ {
				if hasLeadingZero && cur == 0 {
					res += dfs(pos+1, true, (isLimit && cur == up), digitSum, curMod)
				} else {
					nextDigitSum := digitSum + cur
					nextMod := (curMod*10 + cur) % targetDigitSum
					res += dfs(pos+1, false, (isLimit && cur == up), nextDigitSum, nextMod)
				}
			}
			*ptr = res
			return res
		}

		return dfs(0, true, true, 0, 0)
	}

	res := 0
	for digitSum := 1; digitSum <= 9*m; digitSum++ { // 枚举各位数之和
		res += cal(digitSum)
	}
	fmt.Fprintln(out, res)
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
