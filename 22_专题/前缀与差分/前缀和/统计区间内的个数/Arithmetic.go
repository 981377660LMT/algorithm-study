// 区间等差数列常用操作
//
// arithmeticCount
// arithmeticSum
// findFloor
// findCeiling
// findFirst
// findLast

package main

import (
	"fmt"
	"math/rand"
)

func main() {
	check()
}

// 区间[lower, upper]内的奇数个数.
func CountOdds(lower, upper int) int {
	if lower > upper {
		return 0
	}
	return ArithmeticCount(lower, upper, 2, 1)
}

// 区间[lower, upper]内的偶数个数.
func CountEvens(lower, upper int) int {
	if lower > upper {
		return 0
	}
	return ArithmeticCount(lower, upper, 2, 0)
}

// 区间[lower,upper]内形如 `k*x+b` 数的个数.
func ArithmeticCount(lower, upper int, k int, b int) int {
	if lower > upper {
		return 0
	}
	if k == 0 {
		if b >= lower && b <= upper {
			return 1
		}
		return 0
	}
	first, ok1 := FindFirst(lower, upper, k, b)
	if !ok1 {
		return 0
	}
	last, ok2 := FindLast(lower, upper, k, b)
	if !ok2 {
		return 0
	}
	return abs(last-first)/abs(k) + 1
}

// 区间[lower,upper]内形如 `k*x+b` 数的和.
func ArithmeticSum(lower, upper int, k int, b int) int {
	if lower > upper {
		return 0
	}
	if k == 0 {
		if b >= lower && b <= upper {
			return b
		}
		return 0
	}
	first, ok1 := FindFirst(lower, upper, k, b)
	if !ok1 {
		return 0
	}
	last, ok2 := FindLast(lower, upper, k, b)
	if !ok2 {
		return 0
	}
	count := abs(last-first)/abs(k) + 1
	return (first + last) * count / 2
}

// 查找<=x的最大的形如k*x+b的数.
func FindFloor(x int, k, b int) (res int, ok bool) {
	if k == 0 {
		if b <= x {
			return b, true
		}
		return 0, false
	}
	step := abs(k)
	return step*floor(x-b, step) + b, true
}

// 查找>=x的最小的形如k*x+b的数.
func FindCeiling(x int, k, b int) (res int, ok bool) {
	if k == 0 {
		if b >= x {
			return b, true
		}
		return 0, false
	}
	step := abs(k)
	return step*ceil(x-b, step) + b, true
}

// 在区间[lower,upper]内查找第一个形如k*x+b的数.
func FindFirst(lower, upper int, k, b int) (res int, ok bool) {
	if lower > upper {
		return
	}
	if k == 0 {
		if b >= lower && b <= upper {
			return b, true
		}
		return 0, false
	}
	ceiling, ok1 := FindCeiling(lower, k, b)
	if !ok1 {
		return
	}
	if ceiling > upper {
		return
	}
	return ceiling, true
}

// 在区间[lower,upper]内查找最后一个形如k*x+b的数.
func FindLast(lower, upper int, k, b int) (res int, ok bool) {
	if lower > upper {
		return
	}
	if k == 0 {
		if b >= lower && b <= upper {
			return b, true
		}
		return 0, false
	}
	floor, ok1 := FindFloor(upper, k, b)
	if !ok1 {
		return
	}
	if floor < lower {
		return
	}
	return floor, true
}

// -1/2 => -1
func floor(a, b int) int {
	res := a / b
	if a%b != 0 && (a^b) < 0 {
		res--
	}
	return res
}

func ceil(a, b int) int {
	return floor(a+b-1, b)
}

func bmod(a, b int) int {
	return a - b*floor(a, b)
}

func divmod(a, b int) (q, r int) {
	q = floor(a, b)
	r = a - q*b
	return
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func check() {
	findFloorBruteForce := func(x int, k, b int) (res int, ok bool) {
		if k == 0 {
			if b <= x {
				return b, true
			}
			return 0, false
		}
		for cur := x; ; cur-- {
			if (cur-b)%k == 0 {
				return cur, true
			}
		}
	}

	findCeilingBruteForce := func(x int, k, b int) (res int, ok bool) {
		if k == 0 {
			if b >= x {
				return b, true
			}
			return 0, false
		}
		for cur := x; ; cur++ {
			if (cur-b)%k == 0 {
				return cur, true
			}
		}
	}

	findFirstBruteForce := func(lower, upper int, k, b int) (res int, ok bool) {
		if lower > upper {
			return
		}
		if k == 0 {
			if b >= lower && b <= upper {
				return b, true
			}
			return 0, false
		}
		for cur := lower; cur <= upper; cur++ {
			if (cur-b)%k == 0 {
				return cur, true
			}
		}
		return 0, false
	}

	findLastBruteForce := func(lower, upper int, k, b int) (res int, ok bool) {
		if lower > upper {
			return
		}
		if k == 0 {
			if b >= lower && b <= upper {
				return b, true
			}
			return 0, false
		}
		for cur := upper; cur >= lower; cur-- {
			if (cur-b)%k == 0 {
				return cur, true
			}
		}
		return 0, false
	}

	arithmeticCountBruteForce := func(lower, upper int, k, b int) int {
		if k == 0 {
			if b >= lower && b <= upper {
				return 1
			}
			return 0
		}
		res := 0
		for cur := lower; cur <= upper; cur++ {
			if (cur-b)%k == 0 {
				res++
			}
		}
		return res
	}

	arithmeticSumBruteForce := func(lower, upper int, k, b int) int {
		if k == 0 {
			if b >= lower && b <= upper {
				return b
			}
			return 0
		}
		res := 0
		for cur := lower; cur <= upper; cur++ {
			if (cur-b)%k == 0 {
				res += cur
			}
		}
		return res
	}

	for i := 0; i < 100000; i++ {
		x, k, b := -50+rand.Intn(100), -50+rand.Intn(100), -50+rand.Intn(100)
		res1, ok1 := FindFloor(x, k, b)
		res2, ok2 := findFloorBruteForce(x, k, b)
		if ok1 != ok2 || res1 != res2 {
			panic("")
		}
		res3, ok3 := FindCeiling(x, k, b)
		res4, ok4 := findCeilingBruteForce(x, k, b)
		if ok3 != ok4 || res3 != res4 {
			panic("")
		}
	}
	fmt.Println("done1")

	for i := 0; i < 100000; i++ {
		lower, upper := -50+rand.Intn(100), -50+rand.Intn(100)
		k, b := -50+rand.Intn(100), -50+rand.Intn(100)
		res1, ok1 := FindFirst(lower, upper, k, b)
		res2, ok2 := findFirstBruteForce(lower, upper, k, b)
		if ok1 != ok2 || res1 != res2 {
			fmt.Println(ok1, ok2, res1, res2, lower, upper, k, b)
			panic("")
		}
		res3, ok3 := FindLast(lower, upper, k, b)
		res4, ok4 := findLastBruteForce(lower, upper, k, b)
		if ok3 != ok4 || res3 != res4 {
			panic("")
		}
		count1 := ArithmeticCount(lower, upper, k, b)
		count2 := arithmeticCountBruteForce(lower, upper, k, b)
		if count1 != count2 {
			fmt.Println(count1, count2, lower, upper, k, b)
			panic("")
		}
		sum1 := ArithmeticSum(lower, upper, k, b)
		sum2 := arithmeticSumBruteForce(lower, upper, k, b)
		if sum1 != sum2 {
			panic("")
		}
	}
	fmt.Println("done2")

}
