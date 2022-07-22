# 返回移除了所有不包含 1 的子树的原二叉树。
from typing import Optional


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def pruneTree(self, root: Optional["TreeNode"]) -> Optional["TreeNode"]:
        """返回移除了所有不包含 1 的子树的原二叉树"""

        def dfs(root: Optional["TreeNode"]) -> bool:
            if not root:
                return False
            hasLeft = dfs(root.left)
            hasRight = dfs(root.right)
            if not hasLeft:
                root.left = None
            if not hasRight:
                root.right = None
            return hasLeft or hasRight or root.val == 1

        return root if dfs(root) else None
