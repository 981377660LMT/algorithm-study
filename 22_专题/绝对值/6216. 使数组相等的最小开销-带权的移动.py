# 给你两个下标从 0 开始的数组 nums 和 cost ，分别包含 n 个 正 整数。
# 你可以执行下面操作 任意 次：
# 将 nums 中 任意 元素增加或者减小 1 。
# !对第 i 个元素执行一次操作的开销是 cost[i] 。
# !请你返回使 nums 中所有元素 相等 的 最少 总开销。

# n == nums.length == cost.length
# 1 <= n <= 1e5
# 1 <= nums[i], cost[i] <= 1e6

from typing import List, Tuple


class Solution:
    def minCost(self, nums: List[int], weights: List[int]) -> int:
        """
        移动第i个元素1个距离的代价为cost[i]
        求使所有元素相等的最小代价

        枚举最后和哪个元素相等 前缀+后缀计算移动代价
        """

        def calPreSum(pairs: List[Tuple[int, int]]) -> List[int]:
            preSum = [0]
            curWeight = 0
            curCost = 0
            for i in range(1, len(pairs)):
                dist = abs(pairs[i][0] - pairs[i - 1][0])
                curWeight += pairs[i - 1][1]
                curCost += dist * curWeight
                preSum.append(curCost)
            return preSum

        pairs = sorted(zip(nums, weights), key=lambda x: x[0])
        preSum = calPreSum(pairs)
        sufSum = calPreSum(pairs[::-1])[::-1]
        return min(preSum[i] + sufSum[i] for i in range(len(nums)))


print(Solution().minCost(nums=[1, 3, 5, 2], weights=[2, 3, 1, 14]))
