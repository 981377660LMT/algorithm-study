from typing import List

# 请你返回 非空不重叠 子数组的最大数目，且每个子数组中数字和都为 target 。
# 因为是求最大数目，所以找到一个直接清空就行
class Solution:
    def maxNonOverlapping(self, nums: List[int], target: int) -> int:
        res, curSum = 0, 0
        preSum = dict({0: -1})

        for i, num in enumerate(nums):
            curSum += num
            if curSum - target in preSum:
                res += 1
                # 关键
                preSum.clear()
            preSum[curSum] = i

        return res


print(Solution().maxNonOverlapping([-1, 3, 5, 1, 4, 2, -9], 6))
