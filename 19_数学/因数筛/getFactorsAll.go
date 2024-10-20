// 预处理多个数的因数/GetAllFactors

package main

import "fmt"

func main() {
	fmt.Println(GetFactorsAll(10))
}

// 预处理 1~max 的所有数的因数.
func GetFactorsAll(max int32) (res [][]int32) {
	res = make([][]int32, max+1)
	for f := int32(1); f <= max; f++ {
		for m := f; m <= max; m += f {
			res[m] = append(res[m], f)
		}
	}
	return
}
