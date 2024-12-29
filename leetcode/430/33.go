package main

var factorsAll = GetFactorsAll(1e6 + 10)

func GetFactorsAll(max int32) (res [][]int32) {
	res = make([][]int32, max+1)
	for f := int32(1); f <= max; f++ {
		for m := f; m <= max; m += f {
			res[m] = append(res[m], f)
		}
	}
	return
}

func numberOfSubsequences(nums []int) int64 {
	n := len(nums)
	maxVal := int32(maxs(nums...))
	prefixCount := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		prefixCount[i] = make([]int, maxVal+1)
	}
	for i := 1; i <= n; i++ {
		copy(prefixCount[i], prefixCount[i-1])
		val := nums[i-1]
		prefixCount[i][val]++
	}

	countRange := func(start, end int, value int32) int {
		if start >= end {
			return 0
		}
		return prefixCount[end][value] - prefixCount[start][value]
	}

	res := 0
	for i0 := 0; i0 < n; i0++ {
		for i2 := i0 + 1; i2 < n; i2++ {
			mul := int32(nums[i0] * nums[i2])
			for _, f1 := range factorsAll[mul] {
				if f1 > maxVal {
					break
				}
				f2 := mul / f1
				if f2 > maxVal {
					continue
				}

				res += countRange(i0+2, i2-1, f1) * countRange(i2+2, n, f2)
			}
		}
	}

	return int64(res)
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
