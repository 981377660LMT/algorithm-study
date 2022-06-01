# /**
#  * @param {number[]} nums
#  * @param {number} target
#  * @return {number}
#  * @description 考虑排列顺序的完全背包问题
#  * 顺序不同的序列被视作不同的组合。
#  */
from functools import lru_cache
from typing import List


class Solution:
    def combinationSum4(self, nums: List[int], target: int) -> int:
        """考虑排列顺序 每次选哪个物品"""

        @lru_cache(None)
        def dfs(remain: int) -> int:
            if remain <= 0:
                return int(remain == 0)

            res = 0
            for i in range(n):
                if nums[i] <= remain:
                    res += dfs(remain - nums[i])
            return res

        n = len(nums)
        res = dfs(0, target)
        dfs.cache_clear()
        return res

