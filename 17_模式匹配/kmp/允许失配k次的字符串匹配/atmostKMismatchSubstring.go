package main

func minStartingIndex(s string, pattern string) int {
	res, _ := AtmostKMismatchSubstring(s, pattern, 1)
	return res
}

// 最多允许失配k次的子串匹配.
func AtmostKMismatchSubstring(longer, shorter string, k int) (int, int) {
	if len(shorter) > len(longer) {
		return -1, -1
	}

	nums1, nums2 := make([]int, len(longer)), make([]int, len(shorter))
	for i := range longer {
		nums1[i] = int(longer[i])
	}
	for i := range shorter {
		nums2[i] = int(shorter[i])
	}
	return AtmostKMismatchSubarray(nums1, nums2, k)
}

// 最多允许失配k次的子数组匹配.
func AtmostKMismatchSubarray(longer, shorter []int, k int) (int, int) {
	if len(shorter) > len(longer) {
		return -1, -1
	}

	n1, n2 := len(longer), len(shorter)
	dp1 := func() []int {
		arr := append(shorter, longer...)
		return ZAlgo(len(arr), func(i, j int) bool { return arr[i] == arr[j] })
	}()
	dp2 := func() []int {
		shorter = append(shorter[:0:0], shorter...)
		Reverse(shorter)
		longer = append(longer[:0:0], longer...)
		Reverse(longer)
		arr := append(shorter, longer...)
		res := ZAlgo(len(arr), func(i, j int) bool { return arr[i] == arr[j] })
		Reverse(res)
		return res
	}()

	for i := n2; i < n1+1; i++ {
		if dp1[i]+dp2[i-1] >= n2-k {
			return i - n2, i
		}
	}
	return -1, -1
}

func ZAlgo(n int, equal func(i, j int) bool) []int {
	z := make([]int, n)
	left, right := 0, 0
	for i := 1; i < n; i++ {
		z[i] = max(min(z[i-left], right-i+1), 0)
		for i+z[i] < n && equal(z[i], i+z[i]) {
			left, right = i, i+z[i]
			z[i]++
		}
	}
	return z
}

func Reverse[S ~[]E, E any](arr S) {
	n := len(arr)
	for i := 0; i < n/2; i++ {
		arr[i], arr[n-1-i] = arr[n-1-i], arr[i]
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
