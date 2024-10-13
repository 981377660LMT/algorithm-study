# https://leetcode.cn/problems/k-th-largest-perfect-subtree-size-in-binary-tree/
# 返回第 k 大的 完美二叉子树的大小，如果不存在则返回 -1。
# 完美二叉树 是指所有叶子节点都在同一层级的树，且每个父节点恰有两个子节点。


from heapq import nlargest
from typing import Optional, Tuple


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def kthLargestPerfectSubtree(self, root: Optional[TreeNode], k: int) -> int:
        res = []

        def dfs(node: Optional[TreeNode]) -> Tuple[int, bool]:
            """(完美二叉树的高度,是否是完美二叉树)"""
            if not node:
                return 0, True
            leftHeight, isLeftPerfect = dfs(node.left)
            rightHeight, isRightPerfect = dfs(node.right)
            if isLeftPerfect and isRightPerfect and leftHeight == rightHeight:
                res.append((1 << (leftHeight + 1)) - 1)
                return leftHeight + 1, True
            return 0, False

        dfs(root)
        return -1 if k > len(res) else nlargest(k, res)[-1]
