from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一棵 二叉树 的根节点 root 和一个整数k。

# 返回第 k 大的 完美二叉子树 的大小，如果不存在则返回 -1。

# 完美二叉树 是指所有叶子节点都在同一层级的树，且每个父节点恰有两个子节点。


# 子树 是指树中的某一个节点及其所有后代形成的树。
# Definition for a binary tree node.
class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def kthLargestPerfectSubtree(self, root: Optional[TreeNode], k: int) -> int:
        sizes = []

        def dfs(node) -> Tuple[int, int, int, int]:
            """leftHeight, rightHeight, leftSize, rightSize"""
            if not node:
                return 0, 0, 0, 0
            leftLeftHeight, leftRightHeight, leftLeftSize, leftRightSize = dfs(node.left)
            rightLeftHeight, rightRightHeight, rightLeftSize, rightRightSize = dfs(node.right)
            leftHeight = max(leftLeftHeight, leftRightHeight) + 1
            rightHeight = max(rightLeftHeight, rightRightHeight) + 1
            leftSize = leftLeftSize + leftRightSize + 1
            rightSize = rightLeftSize + rightRightSize + 1
            if (
                leftHeight == rightHeight
                and leftSize == rightSize
                and leftSize == (2**leftHeight - 1)
            ):
                sizes.append(leftSize)
            return leftHeight, rightHeight, leftSize, rightSize

        dfs(root)

        sizes.sort(reverse=True)

        return sizes[k - 1] if k <= len(sizes) else -1


# [10,6,2,null,11,10,null]
# 3
