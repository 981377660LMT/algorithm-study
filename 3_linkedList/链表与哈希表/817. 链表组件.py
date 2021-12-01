from typing import List, Optional

# 返回列表 G 中组件的个数，这里对组件的定义为：链表中一段最长连续结点的值（该值必须在列表 G 中）构成的集合。
# a linked list containing unique integer values


class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


# 总结:考虑分界(前面在，后面不在)，分界处加1
class Solution:
    def numComponents(self, head: Optional['ListNode'], nums: List[int]) -> int:
        numSet = set(nums)
        res = 0
        while head:
            if head.val in numSet and (head.next == None or head.next.val not in numSet):
                res += 1
            head = head.next

        return res


# head: 0->1->2->3->4
# G = [0, 3, 1, 4]
# 输出: 2
# 解释:
# 链表中，0 和 1 是相连接的，3 和 4 是相连接的，所以 [0, 1] 和 [3, 4] 是两个组件，故返回 2。
