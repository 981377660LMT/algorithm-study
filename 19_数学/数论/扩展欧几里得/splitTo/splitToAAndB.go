package main

// 2910. 合法分组的最少组数
// https://leetcode.cn/problems/minimum-number-of-groups-to-create-a-valid-assignment/
func minGroupsForValidAssignment(nums []int) int {
	n := len(nums)
	tmpCounter := make(map[int]int)
	for _, v := range nums {
		tmpCounter[v]++
	}
	freq := make([]int, 0, len(tmpCounter))
	for _, v := range tmpCounter {
		freq = append(freq, v)
	}
	freqCounter := make(map[int]int)
	for _, v := range freq {
		freqCounter[v]++
	}

	res := n
	for size := 1; size < n; size++ {
		ok := true
		cand := 0
		for k := range freqCounter {
			count1, count2, ok_ := SplitToAAndB(k, size, size+1, true)
			if !ok_ {
				ok = false
				break
			}
			cand += (count1 + count2) * freqCounter[k]
		}
		if ok {
			res = min(res, cand)
		}
	}

	return res
}

// 将 num 拆分成 a 和 b 的和，使得拆分的个数最(多/少).
//
//	num: 正整数.
//	a: 正整数.
//	b: 正整数.
//	minimize: 是否使得拆分的个数最少. 默认为最少(true).
//
//	返回值: countA和countB分别是拆分成a和b的个数，ok表示是否可以拆分.
func SplitToAAndB(num, a, b int, minimize bool) (count1, count2 int, ok bool) {
	n, x1, y1, x2, y2 := SolveLinearEquation(a, b, num, true)
	if n < 0 {
		ok = false
		return
	}
	if n > 0 {
		res1Smaller := x1+y1 <= x2+y2
		if res1Smaller == minimize {
			return x1, y1, true
		} else {
			return x2, y2, true
		}
	}

	modA, modB := num%a, num%b
	if modA != 0 && modB != 0 {
		ok = false
		return
	}
	if modA != 0 {
		return 0, num / b, true
	}
	if modB != 0 {
		return num / a, 0, true
	}
	div1, div2 := num/a, num/b
	res1Smaller := div1 <= div2
	if res1Smaller == minimize {
		return div1, 0, true
	} else {
		return 0, div2, true
	}
}

// 解一元一次方程 ax + by = c
func SolveLinearEquation(a, b, c int, allowZero bool) (n, x1, y1, x2, y2 int) {
	g, x0, y0 := exgcd(a, b)

	// 无解
	if c%g != 0 {
		n = -1
		return
	}

	a /= g
	b /= g
	c /= g
	x0 *= c
	y0 *= c

	x1 = x0 % b
	if allowZero {
		if x1 < 0 {
			x1 += b
		}
	} else {
		if x1 <= 0 {
			x1 += b
		}
	}

	k1 := (x1 - x0) / b
	y1 = y0 - k1*a

	y2 = y0 % a
	if allowZero {
		if y2 < 0 {
			y2 += a
		}
	} else {
		if y2 <= 0 {
			y2 += a
		}
	}

	k2 := (y0 - y2) / a
	x2 = x0 + k2*b

	if y1 <= 0 {
		return
	}

	n = k2 - k1 + 1
	return
}

func exgcd(a, b int) (gcd, x, y int) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, y, x = exgcd(b, a%b)
	y -= a / b * x
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
