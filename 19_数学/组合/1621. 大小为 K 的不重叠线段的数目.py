from functools import lru_cache
from math import comb

# 请你找到 恰好 k 个不重叠线段且每个线段至少覆盖两个点的方案数
# 这 k 个线段不需要全部覆盖全部 n 个点，且它们的端点 可以 重合。
# 2 <= n <= 1000
# 1 <= k <= n-1
comb = lru_cache(comb)

MOD = int(1e9 + 7)


class Solution:
    def numberOfSets(self, n: int, k: int) -> int:
        return comb(n + k - 1, 2 * k) % MOD  # 如果线段端点可以重合
        return comb(n, 2 * k) % MOD  # 如果线段端点不能重合

    def numberOfSets2(self, n: int, k: int) -> int:
        @lru_cache(None)
        def dfs(start: int, remain: int, isEnd: bool) -> int:
            """端点坐标，剩余线段数，这次选的线段是否需要结束"""
            if remain == 0:
                return int(isEnd)
            if start == n - 1:
                return 0

            res = 0
            if isEnd:
                res = (res + dfs(start + 1, remain, True)) % MOD
                res = (res + dfs(start + 1, remain - 1, True)) % MOD
                res = (res + dfs(start + 1, remain - 1, False)) % MOD
            else:
                res = (res + dfs(start + 1, remain, True)) % MOD
                res = (res + dfs(start + 1, remain, False)) % MOD

            return res

        res = dfs(0, k, False) + dfs(0, k, True)
        dfs.cache_clear()
        return res % MOD


print(Solution().numberOfSets2(4, 2))
print(Solution().numberOfSets2(3, 1))
print(Solution().numberOfSets2(2, 1))

