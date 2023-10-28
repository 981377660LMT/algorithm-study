# 移除二叉树叶子节点的方案数
from math import comb
from typing import Tuple


class Tree:
    def __init__(self, val, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def solve(self, root: Tree) -> int:
        def dfs(root: Tree) -> Tuple[int, int]:
            """返回:子树结点数,排序方案数"""
            if not root:
                return 0, 1

            leftSize, leftRes = dfs(root.left)
            rightSize, rightRes = dfs(root.right)
            return (
                leftSize + rightSize + 1,
                leftRes * rightRes * comb(leftSize + rightSize, leftSize),
            )

        return dfs(root)[1]

