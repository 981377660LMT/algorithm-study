# 面试题 04.10. 检查子树
# https://leetcode.cn/problems/check-subtree-lcci/solutions/2875736/cong-onm-dao-onpythonjavacgo-by-endlessc-9m9v/
# !判断 T2 是否为 T1 的子树.
# !O(n+m) 只在高度相同时匹配

from typing import Optional, Tuple


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def checkSubTree(self, t1: Optional[TreeNode], t2: Optional[TreeNode]) -> bool:
        hs = self._getHeight(t2)

        def dfs(node: Optional[TreeNode]) -> Tuple[int, bool]:
            """返回 node 的高度，以及是否找到了 subRoot."""
            if not node:
                return 0, False
            leftH, leftFound = dfs(node.left)
            rightH, rightFound = dfs(node.right)
            if leftFound or rightFound:
                return 0, True
            nodeH = max(leftH, rightH) + 1
            # !只在子树高度相等时判断是否相等
            return nodeH, nodeH == hs and self._isSameTree(node, t2)

        return dfs(t1)[1]

    def _getHeight(self, root: Optional[TreeNode]) -> int:
        if not root:
            return 0
        return max(self._getHeight(root.left), self._getHeight(root.right)) + 1

    def _isSameTree(self, t1: Optional[TreeNode], t2: Optional[TreeNode]) -> bool:
        if t1 is None or t2 is None:
            return t1 is t2
        return (
            t1.val == t2.val
            and self._isSameTree(t1.left, t2.left)
            and self._isSameTree(t1.right, t2.right)
        )
