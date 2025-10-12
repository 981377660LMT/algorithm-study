package main

func GetSquareFreeKernelList(maxN int) []int {
	maxN++
	res := make([]int, maxN)
	for i := 1; i < maxN; i++ {
		if res[i] == 0 {
			for j := 1; i*j*j < maxN; j++ {
				res[i*j*j] = i
			}
		}
	}
	return res
}
