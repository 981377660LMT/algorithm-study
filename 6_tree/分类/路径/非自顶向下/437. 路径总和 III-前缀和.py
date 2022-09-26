# 求该二叉树里节点值之和等于 targetSum 的 路径 的数目。
# !路径 不需要从根节点开始，也不需要在叶子节点结束，
# 但是路径方向必须是向下的（只能从父节点到子节点）。
# 递归过程中维护路径前缀和 O（n）

from collections import defaultdict
from typing import Optional


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def pathSum(self, root: Optional["TreeNode"], targetSum: int) -> int:
        """dfs过程中维护路径前缀和 O(n)"""

        def dfs(cur: Optional["TreeNode"], curSum: int) -> None:
            if not cur:
                return
            curSum += cur.val
            self.res += preSum[curSum - targetSum]
            preSum[curSum] += 1
            dfs(cur.left, curSum)
            dfs(cur.right, curSum)
            preSum[curSum] -= 1

        preSum = defaultdict(int, {0: 1})
        self.res = 0
        dfs(root, 0)
        return self.res
