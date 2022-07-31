from typing import Optional
from collections import defaultdict


class TreeNode:
    def __init__(
        self, val: int = 0, left: Optional["TreeNode"] = None, right: Optional["TreeNode"] = None
    ):
        self.val = val
        self.left = left
        self.right = right


# 请你找出层内元素之和 最大 的那几层（可能只有一层）的层号，并返回其中 最小 的那个。
class Solution:
    def maxLevelSum(self, root: TreeNode) -> int:
        def dfs(root: Optional["TreeNode"], depth: int) -> None:
            if not root:
                return
            levelSum[depth] += root.val
            dfs(root.left, depth + 1)
            dfs(root.right, depth + 1)

        levelSum = defaultdict(int)
        dfs(root, 1)
        return max(levelSum.keys(), key=lambda x: (levelSum[x], -x))
