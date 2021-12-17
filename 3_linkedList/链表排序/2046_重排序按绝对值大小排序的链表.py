from typing import Optional


class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


# 负数插到头部即可
class Solution:
    def sortLinkedList(self, head: Optional[ListNode]) -> Optional[ListNode]:
        dummy = ListNode(-1, head)
        neg = ListNode(-1)
        negP, dummyP = neg, dummy

        while dummyP.next:
            if dummyP.next.val < 0:
                next = dummyP.next
                dummyP.next = next.next

                # 将负数结点插入negP.next
                negNext = negP.next
                negP.next = next
                next.next = negNext
            else:
                dummyP = dummyP.next

        negP = neg
        while negP.next:
            negP = negP.next
        negP.next = dummy.next
        return neg.next


# Input: head = [0,2,-5,5,10,-10]
# Output: [-10,-5,0,2,5,10]
