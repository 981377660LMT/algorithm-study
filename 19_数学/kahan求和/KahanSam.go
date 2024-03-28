package main

import "fmt"

func main() {
	nums := []float64{0.1, 0.2, 0.3, 0.4, 0.5}
	sum := KahanSummation(len(nums), func(i int) float64 { return nums[i] })
	fmt.Println(sum)
}

// https://oi-wiki.org/misc/kahan-summation/
// Kahan求和算法可以用于一些计算浮点数的和，其造成的误差是 O(1)，即与加总的次数无关.
func KahanSummation(n int, f func(int) float64) float64 {
	sum := 0.0
	err := 0.0
	for i := 0; i < n; i++ {
		y := f(i) - err
		t := sum + y
		err = (t - sum) - y
		sum = t
	}
	return sum
}
