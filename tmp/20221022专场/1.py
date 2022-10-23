from typing import Optional


MOD = int(1e9 + 7)
INF = int(1e20)


# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


class Solution:
    def numberEvenListNode(self, head: Optional[ListNode]) -> int:
        res = 0
        headP = head
        while headP:
            res += headP.val & 1
            headP = headP.next
        return res
