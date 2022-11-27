from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个链表的头节点 head 。

# 对于列表中的每个节点 node ，如果其右侧存在一个具有 严格更大 值的节点，则移除 node 。

# 返回修改后链表的头节点 head 。
# Definition for singly-linked list.


class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


def arrayToLinkedList(nums):
    dummy = ListNode(0)
    p = dummy
    for num in nums:
        p.next = ListNode(num)
        p = p.next
    return dummy.next


def linkedListToArray(head):
    nums = []
    p = head
    while p:
        nums.append(p.val)
        p = p.next
    return nums


class Solution:
    def removeNodes(self, head: Optional[ListNode]) -> Optional[ListNode]:
        nums = linkedListToArray(head)
        stack = []
        for num in nums:
            while stack and stack[-1] < num:
                stack.pop()
            stack.append(num)
        return arrayToLinkedList(stack)
