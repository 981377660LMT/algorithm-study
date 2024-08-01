package main

import (
	"fmt"
	"math"
)

func main() {
	// 浮点数可以用 Nextafter 算出 > x 的下一个浮点数
	// !Nextafter returns the next representable float64 value after x `towards` y.
	fmt.Println(math.Nextafter(float64(1), math.MaxFloat64))  // 1.0000000000000002
	fmt.Println(math.Nextafter(float64(1), -math.MaxFloat64)) // 0.9999999999999999
}

func NextAfter(x float64) float64 {
	return math.Nextafter(x, math.MaxFloat64)
}

func NextBefore(x float64) float64 {
	return math.Nextafter(x, -math.MaxFloat64)
}
