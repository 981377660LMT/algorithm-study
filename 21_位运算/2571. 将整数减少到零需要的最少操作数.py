# https://leetcode.cn/problems/minimum-operations-to-reduce-an-integer-to-0/solution/ji-yi-hua-sou-suo-by-endlesscheng-cm6l/

# 给你一个正整数 n ，你可以执行下述操作 任意 次：
# n 加上或减去 2 的某个 幂
# 返回使 n 等于 0 需要执行的 最少 操作数。
# 如果 x == 2^i 且其中 i >= 0 ，则数字 x 是 2 的幂。
# 1 <= n <= 1e5

# !优先消除最小的比特 1(lowbit) , 因为更高位的bit受低位bit影响


from functools import lru_cache


@lru_cache(None)
def dfs(cur: int) -> int:
    if cur & (cur - 1) == 0:  # cur是2的幂
        return 1
    lowbit = cur & -cur
    return 1 + min(dfs(cur - lowbit), dfs(cur + lowbit))


class Solution:
    def minOperations(self, n: int) -> int:
        return dfs(n)
