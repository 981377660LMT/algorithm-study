# 请你删除 1 条边，使二叉树分裂成两棵子树，且它们子树和的乘积尽可能大。
from collections import defaultdict
from typing import Optional


MOD = int(1e9 + 7)

# Definition for a binary tree node.
class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def maxProduct(self, root: Optional[TreeNode]) -> int:
        def dfs(cur: Optional[TreeNode]) -> int:
            if not cur:
                return 0
            res = dfs(cur.left) + dfs(cur.right) + cur.val
            subtreeSum[id(cur)] = res
            return res

        subtreeSum = defaultdict(int)
        dfs(root)
        allSum = subtreeSum[id(root)]
        res = max(sub * (allSum - sub) for sub in subtreeSum.values())
        return res % MOD
