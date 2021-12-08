from typing import Optional
from collections import defaultdict


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional['TreeNode'] = None, right: Optional['TreeNode'] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# 请你找出层内元素之和 最大 的那几层（可能只有一层）的层号，并返回其中 最小 的那个。
class Solution:
    def maxLevelSum(self, root: TreeNode) -> int:
        def inorder(node, level):
            if node:
                level_sum[level] += node.val
                inorder(node.left, level + 1)
                inorder(node.right, level + 1)

        level_sum = defaultdict(int)
        inorder(root, 1)
        return max(level_sum, key=level_sum.get)


# 直接存储:level=>sum
