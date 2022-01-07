from collections import Counter
from typing import List

INF = 0x7FFFFFFF

# 给你一个整数数组 nums 。nums 中，子数组的 范围 是子数组中最大元素和最小元素的差值。
# 返回 nums 中 所有 子数组范围的 和 。


# dp
class Solution1:
    def subArrayRanges(self, nums: List[int]) -> int:
        res = 0
        n = len(nums)
        for i in range(n):
            minVal = nums[i]
            maxVal = nums[i]
            for j in range(i + 1, n):
                minVal = min(minVal, nums[j])
                maxVal = max(maxVal, nums[j])
                res += maxVal - minVal
        return res


# 单调栈
class Solution:
    def subArrayRanges(self, nums: List[int]) -> int:
        minSum = self.minSum(nums[::])
        maxSum = self.maxSum(nums[::])
        return maxSum - minSum

    def minSum(self, arr):
        arr.append(-0x7FFFFFFF)
        stack = [-1]
        res = 0
        for i in range(len(arr)):
            while stack and arr[stack[-1]] > arr[i]:
                j = stack.pop()
                k = stack[-1]
                res += arr[j] * (i - j) * (j - k)
            stack.append(i)
        return res

    def maxSum(self, arr):
        arr.append(0x7FFFFFFF)
        stack = [-1]
        res = 0
        for i in range(len(arr)):
            while stack and arr[stack[-1]] < arr[i]:
                j = stack.pop()
                k = stack[-1]
                res += arr[j] * (i - j) * (j - k)
            stack.append(i)
        return res


print(Solution().subArrayRanges(nums=[1, 3, 3]))
# print(Solution().subArrayRanges(nums=[1, 2, 3]))

