from functools import lru_cache

# 1 <= steps <= 500
# 1 <= arrLen <= 106
# 每一步操作中，你可以将指针向左或向右移动 1 步，或者停在原地
# 计算并返回：在恰好执行 steps 次操作以后，指针仍然指向索引 0 处的方案数。
MOD = 10 ** 9 + 7

# dfs 记忆化参数为curIndex 和 steps
# 复杂度为O(steps)*O(min(steps,arrlen))


class Solution:
    def numWays(self, steps: int, arrLen: int) -> int:
        @lru_cache(None)
        def dfs(cur: int, s: int) -> int:

            if cur < 0 or cur >= arrLen or s > steps:
                return 0

            # 剪枝,走太远是回不去的
            if abs(cur) > steps - s:
                return 0

            if cur == 0 and s == steps:
                return 1

            return dfs(cur + 1, s + 1) + dfs(cur, s + 1) + dfs(cur - 1, s + 1)

        return dfs(0, 0) % MOD


print(Solution().numWays(3, 2))

