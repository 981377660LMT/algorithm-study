package main

import "fmt"

func main() {
	fmt.Println(largestVariance("abcb"))
	fmt.Println(int(1e18))
}

func largestVariance(s string) (res int) {
	for s1 := 'a'; s1 <= 'z'; s1++ {
		for s2 := 'a'; s2 <= 'z'; s2++ {
			if s1 == s2 {
				continue
			}

			res = max(res, cal(s, s1, s2))
		}
	}

	return
}

func cal(s string, s1, s2 rune) (res int) {
	maxSum, maxSumWithS2 := 0, -int(1e18)
	for _, char := range s {
		if char == s1 {
			maxSum += 1
			maxSumWithS2 += 1
		} else if char == s2 {
			maxSum -= 1
			maxSumWithS2 = maxSum
			maxSum = max(maxSum, 0)
		}

		res = max(res, maxSumWithS2)
	}
	return
}

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}
