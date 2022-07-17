#  给你一个 下标从 0 开始 的整数数组 nums ，返回满足下述条件的 不同 四元组 (a, b, c, d) 的 数目 ：
#  nums[a] + nums[b] + nums[c] == nums[d] ，且
#  a < b < c < d
# 4 <= nums.length <= 50
# 1 <= nums[i] <= 100

from collections import defaultdict
from functools import lru_cache
from typing import List
from itertools import combinations


class Solution:
    def countQuadruplets1(self, nums: List[int]) -> int:
        return sum(a + b + c == d for a, b, c, d in combinations(nums, 4))

    def countQuadruplets2(self, nums: List[int]) -> int:
        """固定中间两个数i2/i3(枚举分割点)，统计左侧两数和，查找右侧两数差
        
        nums[a] + nums[b] == nums[d] - nums[c]
        """
        n, counter = len(nums), defaultdict(int)
        res = 0
        for i2 in range(1, n - 2):
            for i1 in range(i2):
                counter[nums[i1] + nums[i2]] += 1
            i3 = i2 + 1
            for i4 in range(i3 + 1, n):
                res += counter[nums[i4] - nums[i3]]
        return res
        

    def countQuadruplets3(self, nums: List[int]) -> int:
        """二维费用背包dfs写法 从i个物品里选3个，组成容量为nums[i]的背包

        背包问题复杂度关系到nums[i]大小 O(n*4*max(nums[i]))
        """

        @lru_cache(None)
        def dfs(index: int, remainCount: int, remainSum: int) -> int:
            if remainSum < 0 or remainCount < 0:
                return 0
            if index == -1:
                return int(remainSum == 0 and remainCount == 0)
            res = dfs(index - 1, remainCount, remainSum)
            if remainSum - nums[index] >= 0 and remainCount - 1 >= 0:
                res += dfs(index - 1, remainCount - 1, remainSum - nums[index])
            return res

        return sum(dfs(i - 1, 3, nums[i]) for i in range(len(nums)))
        res = 0
        for start in range(len(nums)):
            res += dfs(start - 1, 3, nums[start])
        dfs.cache_clear()
        return res

    def countQuadruplets4(self, nums: List[int]) -> int:
        """二维费用背包dp写法
        """
        max_ = max(nums)
        dp = [[0] * (max_ + 1) for _ in range(4)]
        res = 0
        dp[0][0] = 1
        for num in nums:
            res += dp[3][num]
            for i in range(3, 0, -1):
                for j in range(max_, num - 1, -1):
                    dp[i][j] += dp[i - 1][j - num]
        return res


print(Solution().countQuadruplets2([1, 2, 3, 6]))
print(Solution().countQuadruplets3([1, 2, 3, 6]))
print(Solution().countQuadruplets4([1, 2, 3, 6]))
