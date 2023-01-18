# Definition for a binary tree node.

# 124. 二叉树中的最大路径和


from typing import Optional


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


INF = int(1e18)


class Solution:
    def maxPathSum(self, root: Optional[TreeNode]) -> int:
        # 经过root的最大路径长
        def dfs(root: Optional["TreeNode"]) -> int:
            if not root:
                return 0
            nonlocal res
            left = dfs(root.left)
            right = dfs(root.right)
            res = max(res, left + right + root.val)  # !经过当前节点的最长路径
            return max(0, root.val + max(left, right))  # !子树往上的贡献值(子树最大路径和)

        res = -INF
        dfs(root)
        return res
