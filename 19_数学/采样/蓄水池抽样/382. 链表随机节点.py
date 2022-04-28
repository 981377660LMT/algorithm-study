from random import randint
from typing import Optional


class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


# 思想：每个人从[0,i]里面取数，取到0就输了,看最后一个输的人是谁
# https://leetcode.cn/problems/linked-list-random-node/solution/gong-shui-san-xie-xu-shui-chi-chou-yang-1lp9d/
class Solution:
    def __init__(self, head: Optional[ListNode]):
        self.root = head

    def getRandom(self) -> int:
        """从链表中随机选择一个节点并返回该节点的值。链表中所有节点被选中的概率相等
        
        不使用额外空间
        """
        node, res, i = self.root, -1, 0
        while node:
            if not randint(0, i):
                res = node.val
            node, i = node.next, i + 1
        return res
