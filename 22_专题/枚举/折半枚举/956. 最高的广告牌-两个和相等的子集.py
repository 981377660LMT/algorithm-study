# 956. 最高的广告牌-两个和相等的子集
# 时间复杂度O(len(nums)*sum(nums))
# 分成两个子集，和相等


from functools import lru_cache
from typing import List

INF = int(1e20)


class Solution:
    def tallestBillboard(self, nums: List[int]) -> int:
        @lru_cache(None)
        def dfs(index: int, diff: int) -> int:
            if index == n:
                return 0 if diff == 0 else -INF
            res = -INF
            res = max(res, dfs(index + 1, diff))
            res = max(res, dfs(index + 1, diff - nums[index]))
            res = max(res, dfs(index + 1, diff + nums[index]) + nums[index])
            return res

        n = len(nums)
        return dfs(0, 0)


print(Solution().tallestBillboard([1, 2, 3, 6]))
# 我们有两个不相交的子集 {1,2,3} 和 {6}，它们具有相同的和 sum = 6。
