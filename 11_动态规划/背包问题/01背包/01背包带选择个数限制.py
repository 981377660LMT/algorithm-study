from functools import lru_cache

# 01背包带选择个数限制
class Solution:
    def solve(self, weights, values, capacity, count):
        @lru_cache(None)
        def dfs(i, cp, ct):
            if cp < 0 or ct < 0:
                return -int(1e20)
            if i == len(goods) or ct == 0:
                return 0

            wei, val = goods[i]
            res = dfs(i + 1, cp, ct)
            if cp - wei >= 0:
                res = max(res, dfs(i + 1, cp - wei, ct - 1) + val)
            return res

        goods = list(zip(weights, values))
        res = dfs(0, capacity, count)
        dfs.cache_clear()
        return res
