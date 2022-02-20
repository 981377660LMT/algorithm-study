from typing import List, Optional, Tuple

MOD = int(1e9 + 7)

# Definition for singly-linked list.
class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


class Solution:
    def mergeNodes(self, head: Optional[ListNode]) -> Optional[ListNode]:
        headP = head
        while headP:
            if headP.val == 0:
                headP = headP.next
                continue
            cur = headP.next
            while cur and cur.val != 0:
                headP.val += cur.val
                cur = cur.next
            headP.next = cur
            headP = headP.next

        dummy = ListNode(0, head)
        headP = dummy
        while headP:
            headP.next = headP.next.next if headP.next else None
            headP = headP.next

        return dummy.next


cur = ListNode(1)
curP = cur
for num in [0, 3, 1, 0, 4, 5, 2, 0]:
    curP.next = ListNode(num)
    curP = curP.next
print(Solution().mergeNodes(cur))
