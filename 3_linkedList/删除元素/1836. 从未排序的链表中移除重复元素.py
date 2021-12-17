from collections import Counter


class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


# 统计频率再删除
class Solution:
    def deleteDuplicatesUnsorted(self, head: ListNode) -> ListNode:
        counter = Counter()
        headP = head
        while headP:
            counter[headP.val] += 1
            headP = headP.next

        dummy = ListNode(0, head)
        pre = dummy
        headP = head
        while headP:
            # 删除多余一次的
            while headP and counter[headP.val] > 1:
                pre.next = headP.next
                headP = headP.next
            if headP:
                pre = headP
                headP = headP.next

        return dummy.next

