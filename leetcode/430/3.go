package main

import "fmt"

var factorsAll = GetFactorsAll(1e6 + 10)

func GetFactorsAll(max int32) (res [][]int32) {
	res = make([][]int32, max+1)
	for f := int32(1); f <= max; f++ {
		for m := f; m <= max; m += f {
			res[m] = append(res[m], f)
		}
	}
	return
}

func main() {
	sum := 0
	for i := 1; i <= 1e6; i++ {
		sum += len(factorsAll[i])
	}
	fmt.Println(sum)
}
