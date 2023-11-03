# https://leetcode.cn/problems/populating-next-right-pointers-in-each-node-ii/solutions/2510360/san-chong-fang-fa-dfsbfsbfslian-biao-fu-1bmqp/?envType=daily-question&envId=2023-11-03
# 类似 react fiber
# BFS+链表
# 每一层都连接成一个链表了，那么知道链表头，就能访问这一层的所有节点。
# 在 BFS 的时候，可以一边遍历当前层的节点，一边把下一层的节点连接起来。这样就无需存储下一层的节点了，只需要拿到下一层链表的头节点。

from typing import Optional


class Node:
    __slots__ = "val", "left", "right", "next"

    def __init__(
        self,
        val: int = 0,
        left: Optional["Node"] = None,
        right: Optional["Node"] = None,
        next: Optional["Node"] = None,
    ):
        self.val = val
        self.left = left
        self.right = right
        self.next = next


class Solution:
    def connect(self, root: "Node") -> "Node":
        cur = root
        while cur:
            dummy = Node()  # 下一层的链表
            next_ = dummy
            while cur:  # 遍历当前层的链表
                if cur.left:
                    next_.next = cur.left
                    next_ = cur.left
                if cur.right:
                    next_.next = cur.right
                    next_ = cur.right
                cur = cur.next
            cur = dummy.next
        return root
