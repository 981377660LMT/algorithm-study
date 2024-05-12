# 1 <= n <= 15

# 假设有从 1 到 n 的 n 个整数。
# 用这些整数构造一个数组 perm（下标从 1 开始），只要满足下述条件 之一 ，该数组就是一个 优美的排列 ：

# perm[i] 能够被 i 整除
# i 能够被 perm[i] 整除
# 给你一个整数 n ，返回可以构造的 优美排列 的 数量 。


from functools import lru_cache


class Solution:
    def countArrangement(self, n: int) -> int:
        @lru_cache(None)
        def dfs(state: int) -> int:
            if state == mask:
                return 1
            res = 0
            index = state.bit_count()
            for cur in range(1, n + 1):
                if state & (1 << cur):
                    continue
                if cur % (index + 1) == 0 or (index + 1) % cur == 0:
                    res += dfs(state | (1 << cur))
            return res

        mask = sum(1 << i for i in range(1, n + 1))
        res = dfs(0)
        dfs.cache_clear()
        return res
