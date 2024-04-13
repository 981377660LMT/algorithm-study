package main

// 给定两个数组，每个数组包含1到n的数字。
// 给定m，对于数组1中的数字x和数组2中的数字y，如果x+y<=m，则它们之间有一条边(x,y)。
// 你需要找到这个二部图上的最大匹配。
type NumberSumLeqMatch struct {
	n, m        int
	used, type1 int
}

func NewNumberSumLeqMatch(n, m int) *NumberSumLeqMatch {
	used := min(n, m-1)
	type1 := m - used - 1
	return &NumberSumLeqMatch{n: n, m: m, used: used, type1: type1}
}

func (nm *NumberSumLeqMatch) MaxMatching() int {
	return nm.used
}

// 如果x没有匹配，返回-1
func (nm *NumberSumLeqMatch) Partner(x int) int {
	if x > nm.used {
		return -1
	}
	if x <= nm.type1 {
		return x
	}
	return nm.m - x
}

// 如果从数组1中删除a和从数组2中删除b，最大匹配数是多少.
func (nm *NumberSumLeqMatch) MaxMatching2(a, b int) int {
	if a >= nm.m && b >= nm.m {
		return nm.used
	}
	if a+b < nm.m && nm.type1 <= 1 {
		return nm.used - 2
	}
	return nm.used - 1
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
