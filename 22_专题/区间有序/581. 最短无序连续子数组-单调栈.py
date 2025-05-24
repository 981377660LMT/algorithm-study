# 581. 最短无序连续子数组
# https://leetcode.cn/problems/shortest-unsorted-continuous-subarray/description/
#
# 给你一个整数数组 nums ，你需要找出一个 连续子数组 ，如果对这个子数组进行升序排序，那么整个数组都会变为升序排序。
# 请你找出符合题意的 最短 子数组，并输出它的长度。
#
# !结果肯定是一段前缀和一段后缀不用排序，中间那段需要排序。

from typing import List


def nearestLeftGreater(nums: List[int]) -> List[int]:
    """对每个下标i, 返回 i 左侧最近的严格大于 nums[i] 的下标.若不存在则为 -1."""
    n = len(nums)
    res = [-1] * n
    stack = []
    for i, v in enumerate(nums):
        while stack and nums[stack[-1]] <= v:
            stack.pop()
        res[i] = stack[-1] if stack else -1
        stack.append(i)
    return res


def nearestRightSmaller(nums: List[int]) -> List[int]:
    """对每个下标i, 返回 i 右侧最近的严格小于 nums[i] 的下标.若不存在则为 n."""
    n = len(nums)
    res = [n] * n
    stack = []
    for i, v in enumerate(nums):
        while stack and nums[stack[-1]] > v:
            j = stack.pop()
            res[j] = i
        stack.append(i)
    return res


class Solution:
    def findUnsortedSubarray(self, nums: List[int]) -> int:
        n = len(nums)
        lefts, rights = nearestLeftGreater(nums), nearestRightSmaller(nums)
        leftMin, rightMax = n, -1  # !需要排序的区间的左右边界
        for i, (left, right) in enumerate(zip(lefts, rights)):
            if left != -1 and right != n:
                leftMin = min(leftMin, left)
                rightMax = max(rightMax, right)
            elif left != -1:
                leftMin = min(leftMin, left)
                rightMax = max(rightMax, i)
            elif right != n:
                leftMin = min(leftMin, i)
                rightMax = max(rightMax, right)
        return max(0, rightMax - leftMin + 1)
