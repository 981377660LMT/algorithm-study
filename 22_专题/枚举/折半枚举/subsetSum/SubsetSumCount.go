// 和在[l, r)之间的子集个数

package main

import (
	"fmt"
	"time"
)

func main() {
	// // https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=DPL_4_B
	// in := bufio.NewReader(os.Stdin)
	// out := bufio.NewWriter(os.Stdout)
	// defer out.Flush()

	// var n, k, l, r int
	// fmt.Fscan(in, &n, &k, &l, &r)
	// nums := make([]int, n)
	// for i := 0; i < n; i++ {
	// 	fmt.Fscan(in, &nums[i])
	// }
	// fmt.Fprintln(out, subsetSumCountBySize(nums, l, r+1)[k])

	time1 := time.Now()
	nums := make([]int, 40)
	for i := range nums {
		nums[i] = i
	}
	subsetSumCountBySize(nums, 0, 1000000000)
	time2 := time.Now()
	fmt.Println(time2.Sub(time1)) // 1.5s

}

// 2^n个子集中, 有多少个子集的和在[floor, higher)之间
//
//	O(2^{N/2})
func SubsetSumCount(nums []int, floor, higher int) int {
	n := len(nums)
	getDp := func(nums []int) []int {
		dp := []int{0}
		for _, x := range nums {
			tmp := make([]int, len(dp))
			for i := range tmp {
				tmp[i] = dp[i] + x
			}

			ndp := make([]int, len(dp)+len(tmp))
			i, j, k := 0, 0, 0
			for i < len(dp) && j < len(tmp) {
				if dp[i] < tmp[j] {
					ndp[k] = dp[i]
					i++
				} else {
					ndp[k] = tmp[j]
					j++
				}
				k++
			}
			for i < len(dp) {
				ndp[k] = dp[i]
				i++
				k++
			}
			for j < len(tmp) {
				ndp[k] = tmp[j]
				j++
				k++
			}

			dp = ndp
		}

		return dp
	}

	nums1 := nums[:n/2]
	nums2 := nums[n/2:]
	dp1, dp2 := getDp(nums1), getDp(nums2)
	cal := func(limit int) int {
		res := 0
		right := len(dp2)
		for _, x := range dp1 {
			for right > 0 && x+dp2[right-1] >= limit {
				right--
			}
			res += right
		}
		return res
	}

	return cal(higher) - cal(floor)
}

// 2^n个子集中, 有多少个子集的和在[floor, higher)之间, 求出大小为 0,1,...,n 时的子集的个数
//
//	O(2^{N/2})
func subsetSumCountBySize(nums []int, floor, higher int) []int {
	n := len(nums)
	// dp[i]表示用i个数能组成的所有和
	getDp := func(curNums []int) [][]int {
		dp := [][2]int{{0, 0}}
		for _, x := range curNums {
			tmp := make([][2]int, len(dp))
			copy(tmp, dp)
			for i := range tmp {
				tmp[i][0] += x
				tmp[i][1]++
			}
			ndp := make([][2]int, 0, len(dp)+len(tmp))
			i, j := 0, 0
			for i < len(dp) && j < len(tmp) {
				if dp[i][0] < tmp[j][0] {
					ndp = append(ndp, dp[i])
					i++
				} else {
					ndp = append(ndp, tmp[j])
					j++
				}
			}
			for i < len(dp) {
				ndp = append(ndp, dp[i])
				i++
			}
			for j < len(tmp) {
				ndp = append(ndp, tmp[j])
				j++
			}
			dp = ndp
		}

		res := make([][]int, len(curNums)+1)
		for _, p := range dp {
			res[p[1]] = append(res[p[1]], p[0])
		}
		return res
	}

	nums1 := nums[:n/2]
	nums2 := nums[n/2:]
	dp1, dp2 := getDp(nums1), getDp(nums2)

	cal := func(limit int) []int {
		res := make([]int, n+1)
		for s1, X := range dp1 {
			for s2, Y := range dp2 {
				right := len(Y)
				for _, x := range X {
					for right > 0 && x+Y[right-1] >= limit {
						right--
					}
					res[s1+s2] += right
				}
			}
		}
		return res
	}

	count1 := cal(higher)
	count2 := cal(floor)
	for i := range count1 {
		count1[i] -= count2[i]
	}
	return count1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
