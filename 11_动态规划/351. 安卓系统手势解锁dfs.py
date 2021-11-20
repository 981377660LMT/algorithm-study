# @param {number} m
# @param {number} n
# 1 <= m, n <= 9
# @return {number}
# 那么请你统计一下有多少种 不同且有效的解锁模式 ，是 至少 需要经过 m 个点，但是 不超过 n 个点的。
# 此题类似于哈密尔顿路径的解法:状态压缩
# 不压缩n! 压缩n^2*2^n
from functools import lru_cache


class Solution:
    def numberOfPatterns(self, m: int, n: int) -> int:
        # 跳过了中点
        def isInvalidTransfer(visited: int, cur: int, next: int):
            x1, y1 = divmod(cur, 3)
            x2, y2 = divmod(next, 3)
            mid = ((x1 + x2) // 2) * 3 + (y1 + y2) // 2
            return ((x1 + x2) & 1 == 0) and ((y1 + y2) & 1 == 0) and (visited & (1 << mid)) == 0

        @lru_cache(None)
        def dfs(cur: int, visited: int, count: int) -> int:
            if count == n:
                return 1

            # 本次算不算
            res = 0 if count < m else 1
            for next in range(9):
                if visited & (1 << next) or isInvalidTransfer(visited, cur, next):
                    continue
                res += dfs(next, visited | (1 << next), count + 1)
            return res

        res = 0
        for start in range(9):
            res += dfs(start, 1 << start, 1)
        return res


print(Solution().numberOfPatterns(1, 1))
print(Solution().numberOfPatterns(1, 2))

