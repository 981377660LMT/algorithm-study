from typing import Optional


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional['TreeNode'] = None, right: Optional['TreeNode'] = None
    ):
        self.val = val
        self.left = left
        self.right = right
        self.parent = None


# 每个节点都包含其父节点的引用（指针）。
# p 和 q 存在于树中。

# 思路一：找链表的相交结点
class Solution:
    def lowestCommonAncestor(self, p: 'TreeNode', q: 'TreeNode') -> 'TreeNode':
        n1 = p
        n2 = q
        while n1 != n2:
            n1 = n1.parent or q
            n2 = n2.parent or p
        return n1

