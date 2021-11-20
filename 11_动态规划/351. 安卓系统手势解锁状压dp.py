# @param {number} m
# @param {number} n
# 1 <= m, n <= 9
# @return {number}
# 那么请你统计一下有多少种 不同且有效的解锁模式 ，是 至少 需要经过 m 个点，但是 不超过 n 个点的。
# 此题类似于哈密尔顿路径的解法:状态压缩
# 不压缩n! 压缩n^2*2^n


class Solution:
    def numberOfPatterns(self, m: int, n: int) -> int:
        # 跳过了中点
        def isInvalidTransfer(visited: int, cur: int, next: int):
            x1, y1 = divmod(cur, 3)
            x2, y2 = divmod(next, 3)
            mid = ((x1 + x2) // 2) * 3 + (y1 + y2) // 2
            return ((x1 + x2) & 1 == 0) and ((y1 + y2) & 1 == 0) and (visited & (1 << mid)) == 0

        # dp[state][i] 代表访问过的点集合是 state 且最后一个点是 i 的有效序列数量
        dp = [[0] * 9 for _ in range(1 << 9)]
        for i in range(9):
            dp[1 << i][i] = 1
        for state in range(1 << 9):
            for pre in range(9):
                for cur in range(9):
                    if state & (1 << cur) or isInvalidTransfer(state, pre, cur):
                        continue
                    dp[state | (1 << cur)][cur] += dp[state][pre]
        return sum(sum(dp[state]) for state in range(1 << 9) if m <= bin(state).count('1') <= n)


print(Solution().numberOfPatterns(1, 1))
print(Solution().numberOfPatterns(1, 2))

