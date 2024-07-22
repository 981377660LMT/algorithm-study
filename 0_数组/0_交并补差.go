// 两个数组的交集/并集/补集/差集 (交并补差)
//
// Api:
//  IsSubset(a, b []int) bool // a 是否为 b 的子集（相当于 differenceA 为空）.
//  IsDisjoint(a, b []int) bool // 是否为不相交集合（相当于 intersection 为空.

package main

import (
	"fmt"
)

func main() {
	fmt.Println(Intersection([]int{1, 2, 3, 4}, []int{3, 4, 5, 6}))                   // [3 4]
	fmt.Println(Union([]int{1, 2, 3, 4}, []int{3, 4, 5, 6}))                          // [1 2 3 4 5 6]
	fmt.Println(Difference([]int{3, 4, 5, 6}, []int{1, 2, 3, 4}))                     // [5 6]
	fmt.Println(SymmetricDifference([]int{1, 2, 3, 4}, []int{3, 4, 5, 6}))            // [1 2 5 6]
	fmt.Println(SplitDifferenceAndIntersection([]int{1, 2, 3, 4}, []int{3, 4, 5, 6})) // [1 2] [5 6] [3 4]

	fmt.Println(IsSubset([]int{1, 2, 3}, []int{1, 2, 3, 4, 5})) // true
	fmt.Println(IsDisjoint([]int{1, 2, 3, 4}, []int{4, 5, 6}))  // true
}

// 两个有序数组的交集.
func Intersection(a, b []int) (res []int) {
	i, n := 0, len(a)
	j, m := 0, len(b)
	for {
		if i == n || j == m {
			return
		}
		x, y := a[i], b[j]
		if x < y { // 改成 > 为降序
			i++
		} else if x > y { // 改成 < 为降序
			j++
		} else {
			res = append(res, a[i])
			i++
			j++
		}
	}
}

// 两个有序数组的并集(合并两个有序数组).
func Union(a, b []int) []int { return Merge(a, b) }
func Merge(a, b []int) (res []int) {
	i, n := 0, len(a)
	j, m := 0, len(b)
	res = make([]int, 0, n+m)
	for {
		if i == n {
			res = append(res, b[j:]...)
			break
		}
		if j == m {
			res = append(res, a[i:]...)
			break
		}
		x, y := a[i], b[j]
		if x < y { // 改成 > 为降序
			res = append(res, a[i])
			i++
		} else if x > y { // 改成 < 为降序
			res = append(res, b[j])
			j++
		} else {
			res = append(res, a[i])
			i++
			j++
		}
	}
	res = res[:len(res):len(res)] // clip
	return
}

// 两个有序数组的差集 a-b.
// a b 需要是有序的.
func Difference(a, b []int) (res []int) {
	i, n := 0, len(a)
	j, m := 0, len(b)
	for {
		if i == n {
			return
		}
		if j == m {
			res = append(res, a[i:]...)
			return
		}
		x, y := a[i], b[j]
		if x < y { // 改成 > 为降序
			res = append(res, a[i])
			i++
		} else if x > y { // 改成 < 为降序
			j++
		} else {
			i++
			j++
		}
	}
}

// 两个有序数组的对称差集 a▲b.
// a b 需要是有序的.
func SymmetricDifference(a, b []int) (res []int) {
	i, n := 0, len(a)
	j, m := 0, len(b)
	for {
		if i == n {
			res = append(res, b[j:]...)
			return
		}
		if j == m {
			res = append(res, a[i:]...)
			return
		}
		x, y := a[i], b[j]
		if x < y { // 改成 > 为降序
			res = append(res, a[i])
			i++
		} else if x > y { // 改成 < 为降序
			res = append(res, b[j])
			j++
		} else {
			i++
			j++
		}
	}
}

// 求差集 A-B, B-A 和交集 A∩B
// EXTRA: 求并集 union: A∪B = A-B+A∩B = merge(differenceA, intersection) 或 merge(differenceB, intersection)
// EXTRA: 求对称差 symmetric_difference: A▲B = A-B ∪ B-A = merge(differenceA, differenceB)
// a b 必须是有序的（可以为空）
// 与图论结合 https://codeforces.com/problemset/problem/243/B
func SplitDifferenceAndIntersection(a, b []int) (differenceA, differenceB, intersection []int) {
	i, n := 0, len(a)
	j, m := 0, len(b)
	for {
		if i == n {
			differenceB = append(differenceB, b[j:]...)
			return
		}
		if j == m {
			differenceA = append(differenceA, a[i:]...)
			return
		}
		x, y := a[i], b[j]
		if x < y { // 改成 > 为降序
			differenceA = append(differenceA, x)
			i++
		} else if x > y { // 改成 < 为降序
			differenceB = append(differenceB, y)
			j++
		} else {
			intersection = append(intersection, x)
			i++
			j++
		}
	}
}

// a 是否为 b 的子集（相当于 differenceA 为空）.
// a b 需要是有序的.
func IsSubset(a, b []int) bool {
	i, n := 0, len(a)
	j, m := 0, len(b)
	for {
		if i == n {
			return true
		}
		if j == m {
			return false
		}
		x, y := a[i], b[j]
		if x < y { // 改成 > 为降序
			return false
		} else if x > y { // 改成 < 为降序
			j++
		} else {
			i++
			j++
		}
	}
}

// 是否为不相交集合（相当于 intersection 为空.
// a b 需要是有序的.
func IsDisjoint(a, b []int) bool {
	i, n := 0, len(a)
	j, m := 0, len(b)
	for {
		if i == n || j == m {
			return true
		}
		x, y := a[i], b[j]
		if x < y { // 改成 > 为降序
			i++
		} else if x > y { // 改成 < 为降序
			j++
		} else {
			return false
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
