# 一个正整数数组，你可以用代价为cost的操作将每个数拆成两个部分
# 最后的得分是相等的数字之和*profit
# 求可获得的最大分数

# n ≤ 1,000
# max(nums[i])<=1000


from functools import lru_cache

# 不用记忆化反而过了
class Solution:
    def solve(self, rod_lens, profit_per_len, cost_per_cut):
        # @lru_cache(None)
        def dfs(index: int, targetLen: int) -> int:
            if index == n:
                return 0

            div, mod = divmod(rod_lens[index], targetLen)
            cuts = div - int(mod == 0)
            res = dfs(index + 1, targetLen) + max(
                0, div * profit_per_len * targetLen - cuts * cost_per_cut
            )

            return res

        n = len(rod_lens)
        if not n:
            return 0

        max_ = max(rod_lens)
        res = 0
        for targetLen in range(1, max_ + 1):
            res = max(res, dfs(0, targetLen))

        # dfs.cache_clear()
        return res


print(Solution().solve(rod_lens=[5, 8], profit_per_len=6, cost_per_cut=4))
# We can cut the rod of length 5 into two rods, one with length 4 and the other 1. We can then cut the rod of length 8 into two rods, both with length 4.
#  Then we can sell all 3 rods of length 4 for a total profit of (4 + 4 + 4) * 6 - 8.
