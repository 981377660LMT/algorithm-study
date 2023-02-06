# 树上不相邻选数 求最大值
from typing import Optional, Tuple


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def max(x, y):
    if x > y:
        return x
    return y


def maxs(seq):
    res = seq[0]
    for i in range(1, len(seq)):
        if seq[i] > res:
            res = seq[i]
    return res


class Solution:
    def rob(self, root: Optional[TreeNode]) -> int:
        def dfs(cur: Optional[TreeNode]) -> Tuple[int, int]:
            """不选/选"""
            if not cur:
                return 0, 0
            left0, left1 = dfs(cur.left)
            right0, right1 = dfs(cur.right)
            return max(left0, left1) + max(right0, right1), cur.val + left0 + right0

        return maxs(dfs(root))
