from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数数组 nums 和一个链表的头节点 head。从链表中移除所有存在于 nums 中的节点后，返回修改后的链表的头节点。

# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next


# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next


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
    def modifiedList(self, nums: List[int], head: Optional[ListNode]) -> Optional[ListNode]:
        s = set(nums)
        arr = linkedListToArray(head)
        arr = [x for x in arr if x not in s]
        return arrayToLinkedList(arr)
