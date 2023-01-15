# 合并k个有序链表
# 1. 分治法
# 2. 优先队列/有序集合
# 时间复杂度：O(nlogk)


from typing import List, Optional
from sortedcontainers import SortedList


class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


class Solution:
    def mergeKLists1(self, nodes: List[Optional[ListNode]]) -> Optional[ListNode]:
        sl = SortedList(key=lambda x: x.val)
        for node in nodes:
            if node:
                sl.add(node)

        dummy = ListNode()
        cur = dummy
        while sl:
            min_ = sl.pop(0)
            cur.next = min_
            cur = cur.next
            if min_.next:
                sl.add(min_.next)

        return dummy.next

    def mergeKLists2(self, nodes: List[Optional[ListNode]]) -> Optional[ListNode]:
        def mergeTwo(node1: Optional[ListNode], node2: Optional[ListNode]) -> Optional[ListNode]:
            dummy = ListNode()
            cur = dummy
            while node1 and node2:
                if node1.val < node2.val:
                    cur.next = node1
                    node1 = node1.next
                else:
                    cur.next = node2
                    node2 = node2.next
                cur = cur.next

            cur.next = node1 if node1 else node2
            return dummy.next

        def merge(left: int, right: int) -> Optional[ListNode]:
            if left > right:
                return None
            if left == right:
                return nodes[left - 1]
            if left + 1 == right:
                return mergeTwo(nodes[left - 1], nodes[right - 1])
            mid = (left + right) // 2
            return mergeTwo(merge(left, mid), merge(mid + 1, right))

        return merge(1, len(nodes))  # !类似线段树,范围为[1,n]
