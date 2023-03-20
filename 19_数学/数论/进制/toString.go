package main

import "fmt"

func main() {
	fmt.Println(toString(100, 21))
}

// 十进制转其他进制
func toString(x, base int) (res []int) {
	for ; x > 0; x /= base {
		res = append(res, x%base)
	}
	for i, n := 0, len(res); i < n/2; i++ {
		res[i], res[n-1-i] = res[n-1-i], res[i]
	}
	return
}
