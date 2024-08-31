// 字符串重新排列不含前导0的数字的个数

package main

var fac [21]int

func init() {
	fac[0] = 1
	for i := 1; i < len(fac); i++ {
		fac[i] = fac[i-1] * i
	}
}

func reArrangeDigits(s string) int {
	n := len(s)
	counter := [10]int{}
	for _, c := range s {
		counter[c-'0']++
	}

	res := (n - counter[0]) * fac[n-1]
	for i := 0; i < 10; i++ {
		res /= fac[counter[i]]
	}
	return res
}
