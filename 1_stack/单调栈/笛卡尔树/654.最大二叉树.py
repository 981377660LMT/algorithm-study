# 二叉树的根是数组 nums 中的最大元素。
# 左子树是通过数组中 最大值左边部分 递归构造出的最大二叉树。
# 右子树是通过数组中 最大值右边部分 递归构造出的最大二叉树。
# !瓶颈在于快速找到数组最大元素位置
# 1. 线段树树上二分 O(nlogn)
# 2. 单调栈 O(n) 哈希表存每个范围，对应一个最大值
# !最大元素区间(left,right) 映射到 最大元素的i

# Definition for a binary tree node.


from typing import List, Optional
from 每个元素作为最值的影响范围 import getRange


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def constructMaximumBinaryTree(self, nums: List[int]) -> Optional[TreeNode]:
        """O(n)"""

        def getMaxIndex(left: int, right: int) -> int:
            return mapping[(left, right)]

        def dfs(left: int, right: int) -> Optional[TreeNode]:
            if left > right:
                return None
            maxIndex = getMaxIndex(left, right)
            res = TreeNode(nums[maxIndex])
            res.left = dfs(left, maxIndex - 1)
            res.right = dfs(maxIndex + 1, right)
            return res

        ranges = getRange(nums, isMax=True, isLeftStrict=True, isRightStrict=True)
        mapping = {(left, right): i for i, (left, right) in enumerate(ranges)}
        return dfs(0, len(nums) - 1)

    def constructMaximumBinaryTree2(self, nums: List[int]) -> Optional[TreeNode]:
        """O(n^2)"""

        def getMaxIndex(left: int, right: int) -> int:
            res, max_ = left, nums[left]
            for i in range(left + 1, right + 1):
                if nums[i] > max_:
                    res, max_ = i, nums[i]
            return res

        def dfs(left: int, right: int) -> Optional[TreeNode]:
            if left > right:
                return None
            maxIndex = getMaxIndex(left, right)
            res = TreeNode(nums[maxIndex])
            res.left = dfs(left, maxIndex - 1)
            res.right = dfs(maxIndex + 1, right)
            return res

        return dfs(0, len(nums) - 1)
