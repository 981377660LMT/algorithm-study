# 3444. 使数组包含目标值倍数的最少增量
# https://leetcode.cn/problems/minimum-increments-for-target-multiples-in-an-array/description/
# 给你两个数组 nums 和 target 。
# 在一次操作中，你可以将 nums 中的任意一个元素递增 1 。
# 返回要使 target 中的每个元素在 nums 中 至少 存在一个倍数所需的 最少操作次数 。
#
# 1 <= nums.length <= 5 * 104
# 1 <= target.length <= 4
# target.length <= nums.length
# 1 <= nums[i], target[i] <= 104

from functools import lru_cache
from math import lcm
from typing import List

INF = int(1e18)


def distToMul1(x: int, t: int) -> int:
    """通过+1操作将x变为t的倍数所需的最少操作次数."""
    return (t - x % t) % t


class Solution:
    def minimumIncrements(self, nums: List[int], target: List[int]) -> int:
        n, m = len(nums), len(target)
        subsetLcm = [1] * (1 << m)
        for i, v in enumerate(target):
            bit = 1 << i
            for j in range(1 << i):
                subsetLcm[bit | j] = lcm(subsetLcm[j], v)

        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            if remain == 0:
                return 0
            if index == n:
                return INF

            res = dfs(index + 1, remain)
            # 枚举 remain 的所有非空子集 sub，把 nums[i] 改成 lcms[sub] 的倍数
            sub = remain
            while sub:
                lcm_ = subsetLcm[sub]
                cost = distToMul1(nums[index], lcm_)
                res = min(res, cost + dfs(index + 1, remain ^ sub))
                sub = (sub - 1) & remain
            return res

        res = dfs(0, (1 << m) - 1)
        dfs.cache_clear()
        return res
