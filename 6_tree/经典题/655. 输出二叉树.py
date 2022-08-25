# 输出二叉树
# Definition for a binary tree node.
from typing import List, Optional


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def printTree(self, root: Optional[TreeNode]) -> List[List[str]]:
        def calHeight(root: Optional[TreeNode]) -> int:
            if not root:
                return 0
            return 1 + max(calHeight(root.left), calHeight(root.right))

        def dfs(root: Optional[TreeNode], row: int, col: int) -> None:
            if not root:
                return
            res[row][col] = str(root.val)
            offset = 2 ** (height - row - 2)
            dfs(root.left, row + 1, col - offset)
            dfs(root.right, row + 1, col + offset)

        height = calHeight(root)
        res = [[""] * (2**height - 1) for _ in range(height)]
        dfs(root, 0, (len(res[0]) - 1) // 2)
        return res
