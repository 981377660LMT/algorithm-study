# 2945. 找到最大非递减数组的长度
# https://leetcode.cn/problems/find-maximum-non-decreasing-array-length/description/
from typing import List


class Solution:
    def findMaximumLength(self, nums: List[int]) -> int:
        n = len(nums)
        preSum = [0] * (n + 1)
        for i in range(n):
            preSum[i + 1] = preSum[i] + nums[i]
        dp = [0] * (n + 1)
        stack = [[0, 0, 0]]
        for i in range(1, n + 1):
            left = 0
            right = len(stack) - 1
            while left <= right:
                mid = (left + right) // 2
                if preSum[i] - stack[mid][1] >= stack[mid][0]:
                    left = mid + 1
                else:
                    right = mid - 1
            dp[i] = dp[stack[right][2]] + 1
            last = preSum[i] - stack[right][1]
            while stack and stack[-1][0] + stack[-1][1] >= preSum[i] + last:
                stack.pop()
            stack.append([last, preSum[i], i])

        return max(dp)
