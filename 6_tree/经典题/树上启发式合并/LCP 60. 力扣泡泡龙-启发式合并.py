from typing import Optional
import sys

sys.setrecursionlimit(9999999)


# Definition for a binary tree node.
class TreeNode:
    def __init__(self, x):
        self.val = x
        self.left = None
        self.right = None


# 2 <= 树中节点个数 <= 10^5
# -10000 <= 树中节点的值 <= 10000
# 树上启发式合并
# 为了保证总复杂度是 O(nlogn)，需要把小子树合并到大子树里，也就是启发式合并。
class Solution:
    def getMaxLayerSum(self, root: Optional[TreeNode]) -> int:
        ...


# https://leetcode-cn.com/problems/WInSav/solution/by-tsreaper-5mrh/
