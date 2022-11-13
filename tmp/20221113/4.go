// function maxPalindromes(s: string, k: number): number {
//   const n = s.length
//   const intervals: [number, number][] = []
//   const helper = (left: number, right: number) => {
//     while (left >= 0 && right < n && s[left] === s[right]) {
//       left--
//       right++
//       if (right - left - 1 >= k) {
//         intervals.push([left + 1, right - 1])
//       }
//     }
//   }

//   for (let i = 0; i < n; i++) {
//     helper(i, i)
//     helper(i, i + 1)
//   }

//   intervals.sort((a, b) => a[1] - b[1])
//   let res = 0
//   let preEnd = -1
//   for (const [start, end] of intervals) {
//     if (start > preEnd) {
//       res++
//       preEnd = end
//     }
//   }

//   return res
// }
package main

import "sort"

type interval struct {
	start, end int
}

func maxPalindromes(s string, k int) int {
	n := len(s)
	intervals := make([]interval, 0, n)
	expand := func(left, right int) {
		for left >= 0 && right < n && s[left] == s[right] {
			left--
			right++
			if right-left-1 >= k {
				intervals = append(intervals, interval{left + 1, right - 1})
			}
		}
	}

	for i := 0; i < n; i++ {
		expand(i, i)
		expand(i, i+1)
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].end < intervals[j].end
	})

	res := 0
	preEnd := -1
	for _, interval := range intervals {
		if interval.start > preEnd {
			res++
			preEnd = interval.end
		}
	}
	return res
}
