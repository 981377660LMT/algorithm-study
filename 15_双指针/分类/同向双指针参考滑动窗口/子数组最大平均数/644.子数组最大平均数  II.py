#  找出 长度大于等于 k 且含最大平均值的连续子数组 并输出这个最大平均值
#  任何计算误差小于 10-5 的结果都将被视为正确答案
#  答案误差不超过1e-5是使用二分搜索的提示

# 是否存在长度大于等于 k 的子数组，且平均值大于等于 mid
# !所有数减去mid=> 等价于找出长度大于等于 k 的子数组，且和大于等于 0


from typing import List


EPS = 1e-7


class Solution:
    def findMaxAverage(self, nums: List[int], k: int) -> float:
        def check(mid: float) -> bool:
            """是否存在长度大于等于 k 的子数组，且平均值大于等于 mid"""
            preSum = [0.0] * (n + 1)
            for i in range(n):
                preSum[i + 1] = preSum[i] + nums[i] - mid
            preMin = 0
            for i in range(k, n + 1):
                if preSum[i] - preMin >= 0:
                    return True
                preMin = min(preMin, preSum[i - k + 1])
            return False

        n = len(nums)
        left, right = min(nums), max(nums)
        while left <= right:
            mid = (left + right) / 2
            if check(mid):
                left = mid + EPS
            else:
                right = mid - EPS
        return right
