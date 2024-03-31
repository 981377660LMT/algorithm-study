package main

import "math"

// https://leetcode.cn/problems/linked-list-cycle-ii/description/
type ListNode struct {
	Val  int
	Next *ListNode
}

func detectCycle(head *ListNode) *ListNode {
	pos := [3]*ListNode{head, head, head}
	next := func(i int) {
		if pos[i] != nil {
			pos[i] = pos[i].Next
		}
	}
	equal := func(i, j int) bool { return pos[i] == pos[j] }
	start, _, hasCycle := FloydSearchCircleWithMaxRound(next, equal, -1)
	if !hasCycle {
		return nil
	}
	cur := head
	for i := 0; i < start; i++ {
		cur = cur.Next
	}
	return cur
}

// 三路探环法.
//
// next(i) 表示第 i 个参与者移动到下一个位置(0<=i<=2)
// equal(i, j) 表示第i个参与者与第j个参与者是否在相同位置(0<=i,j<=2)
// maxRound 表示最大探测轮数，如果 start+period > maxRound 则视为不存在环.-1表示无限轮数.
//
// 时间复杂度: O(2*(period+start))
// 空间复杂度: O(1)
func FloydSearchCircleWithMaxRound(
	next func(i int),
	equal func(i, j int) bool,
	maxRound int,
) (start, period int, hasCycle bool) {
	if maxRound == -1 {
		maxRound = math.MaxInt64
	}
	next(0)
	next(0)
	next(1)
	maxRound--
	for !equal(0, 1) && maxRound > 0 {
		next(0)
		next(0)
		next(1)
		maxRound--
	}
	if !equal(0, 1) {
		return
	}
	start = 0
	for !equal(0, 2) {
		next(0)
		next(2)
		start++
	}
	period = 1
	next(0)
	for !equal(0, 2) {
		next(0)
		period++
	}
	return start, period, true
}
