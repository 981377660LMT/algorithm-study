package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var A, B, C string
	fmt.Fscan(in, &A, &B, &C)
	res := solveHukumenzan(A, B, C)
	if len(res) == 0 {
		fmt.Fprintln(out, "UNSOLVABLE")
		return
	}
	a, b, c := res[0][0], res[0][1], res[0][2]
	fmt.Fprintln(out, a, b, c)
}

func solveHukumenzan(A, B, C string) [][]int {
	res := [][]int{}
	set := make(map[byte]struct{})
	for _, v := range A + B + C {
		set[byte(v)] = struct{}{}
	}
	if len(set) > 10 {
		return res
	}
	chars := make([]byte, 0, len(set))
	for k := range set {
		chars = append(chars, k)
	}

	order := make([]int, 10)
	for i := range order {
		order[i] = i
	}

	for order, ok := nextPermutation(order, true); ok; order, ok = nextPermutation(order, true) {
		mp := make(map[byte]int)
		for i := range chars {
			mp[chars[i]] = order[i]
		}
		if mp[A[0]] == 0 || mp[B[0]] == 0 || mp[C[0]] == 0 {
			continue
		}
		a, b, c := 0, 0, 0
		for i := 0; i < len(A); i++ {
			a = 10*a + mp[A[i]]
		}
		for i := 0; i < len(B); i++ {
			b = 10*b + mp[B[i]]
		}
		for i := 0; i < len(C); i++ {
			c = 10*c + mp[C[i]]
		}

		if a+b == c {
			res = append(res, []int{a, b, c})
		}
	}

	return res
}

// 返回下一个字典序的排列.
func nextPermutation(nums []int, inPlace bool) (res []int, ok bool) {
	if !inPlace {
		nums = append(nums[:0:0], nums...)
	}
	left, right := len(nums)-1, len(nums)-1
	for left > 0 && nums[left-1] >= nums[left] {
		left--
	}
	if left == 0 {
		return
	}
	last := left - 1
	for nums[right] <= nums[last] {
		right--
	}
	nums[last], nums[right] = nums[right], nums[last]
	for i, j := last+1, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
	return nums, true
}
