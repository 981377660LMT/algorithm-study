package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(NextFloat64(1.0)) // 1.0000000000000002
	fmt.Println(PrevFloat64(1.0)) // 0.9999999999999999
}

// > x 的下一个浮点数
func PrevFloat64(x float64) float64 {
	return math.Nextafter(x, -math.MaxFloat64)
}

// < x 的下一个浮点数
func NextFloat64(x float64) float64 {
	return math.Nextafter(x, math.MaxFloat64)
}
