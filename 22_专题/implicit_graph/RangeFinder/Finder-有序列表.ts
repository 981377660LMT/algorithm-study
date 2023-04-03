// // 寻找前驱后继/区间删除
// package main

// import (
// 	"fmt"
// 	"sort"
// 	"strings"
// 	"time"
// )

// const INF int = 1e18

// // 2612. 最少翻转操作数
// // https://leetcode.cn/problems/minimum-reverse-operations/
// func minReverseOperations(n int, p int, banned []int, k int) []int {
// 	finder := [2]*Finder{
// 		NewFinder(func(a, b int) int { return a - b }, n/2),
// 		NewFinder(func(a, b int) int { return a - b }, n/2),
// 	}

// 	for i := 0; i < n; i++ {
// 		finder[i&1].Insert(i)
// 	}
// 	for _, i := range banned {
// 		finder[i&1].Erase(i)
// 	}

// 	getRange := func(i int) (int, int) {
// 		return max(i-k+1, k-i-1), min(i+k-1, 2*n-k-i-1)
// 	}
// 	setUsed := func(u int) {
// 		finder[u&1].Erase(u)
// 	}

// 	findUnused := func(u int) int {
// 		left, right := getRange(u)
// 		pre, ok := finder[(u+k+1)&1].Prev(right)
// 		if ok && left <= pre && pre <= right {
// 			return pre
// 		}
// 		next, ok := finder[(u+k+1)&1].Next(left)
// 		if ok && left <= next && next <= right {
// 			return next
// 		}
// 		return -1
// 	}

// 	dist := OnlineBfs(n, p, setUsed, findUnused)
// 	res := make([]int, n)
// 	for i, d := range dist {
// 		if d == INF {
// 			res[i] = -1
// 		} else {
// 			res[i] = d
// 		}
// 	}
// 	return res
// }
