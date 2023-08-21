# 2448. 使数组相等的最小开销
# https://leetcode.cn/problems/minimum-cost-to-make-array-equal/
# n == nums.length == cost.length
# 1 <= n <= 1e5
# 1 <= nums[i], cost[i] <= 1e6
# 给你两个下标从 0 开始的数组 nums 和 cost ，分别包含 n 个 正 整数。
# 你可以执行下面操作 任意 次：
# !将 nums 中 任意 元素增加或者减小 1 。
# !对第 i 个元素执行一次操作的开销是 cost[i] 。
# 请你返回使 nums 中所有元素 相等 的 最少 总开销。


from typing import List
from 三分法求凸函数极值 import minimize


class Solution:
    def minCost(self, nums: List[int], cost: List[int]) -> int:
        """三分法求严格凸函数最值"""

        def fun(pos: int) -> int:
            return sum(abs((num - pos)) * c for num, c in zip(nums, cost))

        return minimize(fun, min(nums), max(nums))


print(Solution().minCost(nums=[1, 3, 5, 2], cost=[2, 3, 1, 14]))
print(Solution().minCost(nums=[1, 2, 3], cost=[1, 1, 1]))
