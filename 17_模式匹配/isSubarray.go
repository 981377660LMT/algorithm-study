package main

import "fmt"

func main() {
	fmt.Println(IsSubstring("abc", "ab"))
}

// O(n+m)判断`shorter`是否是`longer`的子串.
func IsSubstring(longer, shorter string) bool {
	if len(shorter) > len(longer) {
		return false
	}
	if len(shorter) == 0 {
		return true
	}
	ords1 := make([]int, len(longer))
	for i, c := range longer {
		ords1[i] = int(c)
	}
	ords2 := make([]int, len(shorter))
	for i, c := range shorter {
		ords2[i] = int(c)
	}
	return IsSubarray(ords1, ords2)
}

// O(n+m)判断`shorter`是否是`longer`的子数组.
func IsSubarray(longer, shorter []int) bool {
	if len(shorter) > len(longer) {
		return false
	}
	if len(shorter) == 0 {
		return true
	}
	n, m := len(longer), len(shorter)
	st := make([]int, 0, n+m)
	st = append(st, shorter...)
	st = append(st, longer...)
	z := zAlgo(st)
	for i := m; i < n+m; i++ {
		if z[i] >= m {
			return true
		}
	}
	return false
}

func zAlgo(seq []int) []int {
	n := len(seq)
	if n == 0 {
		return nil
	}
	z := make([]int, n)
	j := 0
	for i := 1; i < n; i++ {
		var k int
		if j+z[j] <= i {
			k = 0
		} else {
			k = min(j+z[j]-i, z[i-j])
		}
		for i+k < n && seq[k] == seq[i+k] {
			k++
		}
		if j+z[j] < i+z[i] {
			j = i
		}
		z[i] = k
	}
	z[0] = n
	return z
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
