package main

func main() {

}

// 带权(权重为等差数列)前缀和.
// 权重首项为`first`,公差为`diff`.
// 前缀和为: `first*a0+(first+diff)*a1+...+(first+(k-1)*diff)*ak`.
func WeightedPresum(arr []int, first, diff int) func(start, end int) int {
	preSum1 := make([]int, len(arr)+1)
	preSum2 := make([]int, len(arr)+1)
	for i, v := range arr {
		preSum1[i+1] = preSum1[i] + diff*v
		preSum2[i+1] = preSum2[i] + (first+i*diff)*v
	}
	return func(start, end int) int {
		if start >= end {
			return 0
		}
		return preSum2[end] - preSum2[start] - start*(preSum1[end]-preSum1[start])
	}
}
