// 给定一个区间 [start,end)
// 将区间表示为形如[2^i*j,2^i*(j+1)).
// 保证答案存在.
package main

import (
	"fmt"
	"math/bits"
)

func main() {
	fmt.Println(Format(3, 4))
	fmt.Println(Format(4, 8))
	fmt.Println(Format(8, 16))
	fmt.Println(Format(16, 18))
	fmt.Println(Format(18, 19))
}

func Format(start, end int) (i, j int) {
	n := end - start
	i = bits.Len(uint(n - 1))
	j = start >> i
	return
}
