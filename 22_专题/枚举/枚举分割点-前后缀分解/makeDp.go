package main

func main() {

	var nums = []int{1, 2, 3, 4, 5}
	// 前后缀分解

	makeDp := func(arr []int) []int {
		m := len(arr)
		dp := make([]int, m+1)
		for i := 1; i <= m; i++ {
			cur := arr[i-1]
			// ...
		}
		return dp
	}

	preDp := makeDp(nums)
	sufDp := func() []int {
		tmp := append(nums[:0:0], nums...)
		for i, j := 0, len(tmp)-1; i < j; i, j = i+1, j-1 {
			tmp[i], tmp[j] = tmp[j], tmp[i]
		}
		res := makeDp(tmp)
		for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
			res[i], res[j] = res[j], res[i]
		}
		return res
	}()

	res := 0
	for i := 0; i < len(nums); i++ {
		res += preDp[i] * sufDp[i] // [0,i) x [i,n)
	}

}
