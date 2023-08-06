from math import gcd
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


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
    def insertGreatestCommonDivisors(self, head: Optional[ListNode]) -> Optional[ListNode]:
        arr = linkedListToArray(head)
        res = []
        for i in range(len(arr)):
            res.append(arr[i])
            if i < len(arr) - 1:
                res.append(gcd(arr[i], arr[i + 1]))
        return arrayToLinkedList(res)
