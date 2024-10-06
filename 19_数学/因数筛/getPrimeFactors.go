package main

import "fmt"

func main() {
	GetPrimeFactors(100, func(v, c int) {
		fmt.Println(v, c)
	})
}

// 质因数分解.
func GetPrimeFactors(n int, f func(v, c int)) {
	upper := n
	i := 2
	for i*i <= upper {
		if upper%i == 0 {
			c := 0
			for upper%i == 0 {
				c++
				upper /= i
			}
			f(i, c)
		}
		i++
	}
	if upper != 1 {
		f(upper, 1)
	}
}
