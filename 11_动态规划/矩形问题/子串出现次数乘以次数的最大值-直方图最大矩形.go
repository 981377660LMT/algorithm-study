// https://www.luogu.com.cn/problem/P3804
// 给定一个长度为 n 的只包含小写字母的字符串 s。
// !对于所有 s 的出现次数不为 1 的子串，设其 value值为该 子串出现的次数 × 该子串的长度。
// 请计算，value 的最大值是多少。
// n <= 1e6

// 直方图最大矩形
// lcp范围看成宽,lcp看成高
// https://www.acwing.com/solution/content/25201/

package main

import (
	"bufio"
	"fmt"
	"index/suffixarray"
	"os"
	"reflect"
	"unsafe"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	fmt.Fprintln(out, solve(s))
}

func solve(s string) int {
	if len(s) <= 1 {
		return 0
	}
	_, _, heights := suffixArray(s)
	L, R := getRange(heights, false, false, false) // 求每个元素作为最小值的影响范围(区间)
	res := 0
	for i := 0; i < len(heights); i++ {
		res = max(res, heights[i]*(R[i]-L[i]+2))
	}
	return res
}

// 求每个元素作为最值的影响范围(区间)
func getRange(nums []int, isMax, isLeftStrict, isRightStrict bool) (leftMost, rightMost []int) {
	compareLeft := func(stackValue, curValue int) bool {
		if isLeftStrict && isMax {
			return stackValue <= curValue
		} else if isLeftStrict && !isMax {
			return stackValue >= curValue
		} else if !isLeftStrict && isMax {
			return stackValue < curValue
		} else {
			return stackValue > curValue
		}
	}

	compareRight := func(stackValue, curValue int) bool {
		if isRightStrict && isMax {
			return stackValue <= curValue
		} else if isRightStrict && !isMax {
			return stackValue >= curValue
		} else if !isRightStrict && isMax {
			return stackValue < curValue
		} else {
			return stackValue > curValue
		}
	}

	n := len(nums)
	leftMost, rightMost = make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		rightMost[i] = n - 1
	}

	stack := []int{}
	for i := 0; i < n; i++ {
		for len(stack) > 0 && compareRight(nums[stack[len(stack)-1]], nums[i]) {
			rightMost[stack[len(stack)-1]] = i - 1
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	stack = []int{}
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && compareLeft(nums[stack[len(stack)-1]], nums[i]) {
			leftMost[stack[len(stack)-1]] = i + 1
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// https://github.dev/EndlessCheng/codeforces-go/copypasta/strings.go
func suffixArray(s string) (sa []int32, rank, height []int) {
	n := len(s)
	sa = *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New([]byte(s))).Elem().FieldByName("sa").Field(0).UnsafeAddr()))
	rank = make([]int, n)
	for i := range rank {
		rank[sa[i]] = i
	}
	height = make([]int, n)
	h := 0
	for i, rk := range rank {
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := int(sa[rk-1]); i+h < n && j+h < n && s[i+h] == s[j+h]; h++ {
			}
		}
		height[rk] = h
	}
	return
}
