# !O(logn*logn)求完全二叉树的结点个数
# Definition for a binary tree node.


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


from typing import Optional


class Solution:
    def countNodes(self, root: Optional["TreeNode"]) -> int:
        def getDepth(node: Optional["TreeNode"]) -> int:
            return 0 if not node else 1 + getDepth(node.left)

        if root is None:
            return 0

        leftDepth, rightDepth = getDepth(root.left), getDepth(root.right)
        if leftDepth == rightDepth:
            return 2**leftDepth + self.countNodes(root.right)
        return 2**rightDepth + self.countNodes(root.left)
