# 2816. 翻倍以链表形式表示的数字-链表倒序处理技巧 - 递归
# 首先 dfs 到最深处，然后在 dfs 返回的时候进行元素处理，这样不就相当于倒序遍历链表了吗！
# https://leetcode.cn/problems/double-a-number-represented-as-a-linked-list/solutions/2385928/lian-biao-dao-xu-chu-li-ji-qiao-di-gui-b-0n1z/


# 给你一个 非空 链表的头节点 head ，表示一个不含前导零的非负数整数。
# 将链表 翻倍 后，返回头节点 head 。


from typing import Optional


class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


class Solution:
    def doubleIt(self, head: Optional[ListNode]) -> Optional[ListNode]:
        def dfs(head: Optional[ListNode]) -> int:
            if head is None:
                return 0
            next_ = dfs(head.next)
            cur = head.val * 2 + next_
            head.val = cur % 10
            return cur // 10

        res = dfs(head)
        return ListNode(res, head) if res else head
