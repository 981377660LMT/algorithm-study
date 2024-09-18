package main

import "fmt"

func main() {
	EnumerateMultiset(3, 2, func(a []int) {
		fmt.Println(a)
	})
}

// N个元素，k种类型，频率之和为sum
// A[0]+...+A[N-1] == sum
func EnumerateMultiset(n, sum int, f func([]int)) {
	res := make([]int, n)
	var dfs func(int, int)
	dfs = func(index, remain int) {
		if index == n {
			if remain == 0 {
				f(res)
			}
			return
		}
		for x := 0; x <= remain; x++ {
			res[index] = x
			dfs(index+1, remain-x)
		}
	}
	dfs(0, sum)
}
