from typing import List

INF = int(1e20)


class Solution:
    def maxSubArray(self, nums: List[int]) -> int:
        """最大子数组和-取或全不取dp"""

        maxSum, res = -INF, -INF
        for num in nums:
            # 如果curMax为负数，则前面全不取
            maxSum = max(maxSum + num, num)
            res = max(res, maxSum)  # 以当前元素结尾的最大子数组和
        return res

    def maxSubArray3(self, nums: List[int]) -> int:
        """最大子数组和-取或全不取dp"""

        maxSum, res = -INF, -INF
        for num in nums:
            if maxSum < 0:
                maxSum = 0
            maxSum += num
            res = max(res, maxSum)  # 以当前元素结尾的最大子数组和
        return res

    def maxSubArray2(self, nums: List[int]) -> int:
        """最大子数组和-前缀和"""

        curSum, preMin = 0, 0
        res = -INF
        for num in nums:
            curSum += num
            res = max(res, curSum - preMin)
            preMin = min(preMin, curSum)
        return res
