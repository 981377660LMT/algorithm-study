# https://leetcode.cn/problems/the-number-of-beautiful-subsets/
# 返回数组 nums 中 非空 且 美丽 的子集数目。
# !如果 nums 的子集中，任意两个整数的绝对差均不等于 k ，则认为该子数组是一个 美丽 子集。
# 1 <= nums.length <= 20
# 1 <= nums[i], k <= 1000

# 解法1:dfs回溯
# 解法2:背包dp O(n)


from collections import defaultdict
from typing import List

MOD = int(1e9 + 7)


class Solution:
    def beautifulSubsets(self, nums: List[int], k: int) -> int:
        """
        !相差为k=>模k同余
        按模分组后的打家劫舍(不能选相邻相差为k的数)
        对每一组,记每种数为group,每种数的个数为counter
        dp[i]表示选前i种数的方案数
        """

        def cal(houses: List[int], counter: List[int]) -> int:
            """
            一排房屋升序排列,相邻差为k的房屋不能同时偷.
            每种房屋counter[i]个,每种房屋可以打劫任意多个,求打劫的方案数.
            结果对MOD取模.
            """
            n = len(houses)
            if n == 0:
                return 1
            dp0, dp1 = 1, pow(2, counter[0], MOD) - 1  # 不偷当前/偷当前 的方案数
            for i in range(1, n):
                if houses[i] - houses[i - 1] == k:
                    dp0, dp1 = dp0 + dp1, dp0 * (pow(2, counter[i], MOD) - 1)
                else:
                    dp0, dp1 = dp0 + dp1, (dp0 + dp1) * (pow(2, counter[i], MOD) - 1)
                dp0, dp1 = dp0 % MOD, dp1 % MOD
            return (dp0 + dp1) % MOD

        mp = defaultdict(lambda: defaultdict(int))
        for num in nums:
            mp[num % k][num] += 1
        res = 1
        for vs in mp.values():
            keys = sorted(vs)
            counter = [vs[key] for key in keys]
            res *= cal(keys, counter)
            res %= MOD
        return (res - 1) % MOD  # 去除空集


print(Solution().beautifulSubsets(nums=[2, 4, 6], k=2))
print(Solution().beautifulSubsets(nums=[2, 3, 5, 8], k=5))
