package main

// 合并K个有序数据结构.
func MergeKSorted[E any](sortedItems []E, merge func(E, E) E) (res E) {
	n := len(sortedItems)
	if n == 0 {
		return
	}
	if n == 1 {
		return sortedItems[0]
	}
	if n == 2 {
		return merge(sortedItems[0], sortedItems[1])
	}

	var f func(start, end int) E
	f = func(start, end int) E {
		if end-start == 1 {
			return sortedItems[start]
		}
		mid := (start + end) >> 1
		return merge(f(start, mid), f(mid, end))
	}
	return f(0, n)
}

// https://leetcode.cn/problems/merge-k-sorted-lists/
type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeKLists(lists []*ListNode) *ListNode {
	merge := func(a, b *ListNode) *ListNode {
		dummy := &ListNode{}
		cur := dummy
		for a != nil && b != nil {
			if a.Val < b.Val {
				cur.Next = a
				a = a.Next
			} else {
				cur.Next = b
				b = b.Next
			}
			cur = cur.Next
		}
		if a != nil {
			cur.Next = a
		}
		if b != nil {
			cur.Next = b
		}
		return dummy.Next
	}
	return MergeKSorted(lists, merge)
}
