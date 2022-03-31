# 求字符串有多少种不同的划分方案，使得各个数字都在[1,k]之间
from functools import lru_cache


MOD = int(1e9 + 7)

# n ≤ 100,000
# 状态只与当前index有关


class Solution:
    def solve(self, s: str, k: int) -> int:
        @lru_cache(None)
        def dfs(index: int) -> int:
            if index >= n:
                return 1

            res = 0
            curNum = 0
            j = index
            while j < n:
                curNum = curNum * 10 + int(s[j])
                if 1 <= curNum <= k:
                    res += dfs(j + 1)
                    res %= MOD
                else:
                    break
                j += 1

            return res % MOD

        n = len(s)
        res = dfs(0)
        dfs.cache_clear()
        return res


# s could represent [1, 2, 3], [12, 3], [1, 23], [123].
print(Solution().solve(s="123", k=200))
