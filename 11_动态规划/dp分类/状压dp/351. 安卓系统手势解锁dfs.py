# @param {number} m
# @param {number} n
# 1 <= m, n <= 9
# @return {number}
# !那么请你统计一下安卓系统有多少种不同且有效的解锁模式是至少需要经过m个点，但是不超过n个点的。
# 此题类似于哈密尔顿路径的解法:状态压缩
# 不压缩n! 压缩n^2*2^n

from functools import lru_cache


@lru_cache(None)
def check(visited: int, cur: int, next: int) -> bool:
    """转移是否合法"""
    r1, c1 = cur // 3, cur % 3
    r2, c2 = next // 3, next % 3
    if (r1 + r2) & 1 or (c1 + c2) & 1:  # 没有经过中间点时
        return True
    mid = ((r1 + r2) // 2) * 3 + (c1 + c2) // 2
    return not not visited & (1 << mid)  # 经过中间点时,中间点需要访问过


@lru_cache(None)
def cal(upper: int) -> int:
    """不超过upper个点的解锁模式"""

    @lru_cache(None)
    def dfs(cur: int, visited: int, count: int) -> int:
        if count == upper:
            return 1

        res = 1
        for next in range(9):
            if visited & (1 << next):
                continue
            if check(visited, cur, next):
                res += dfs(next, visited | (1 << next), count + 1)
        return res

    if upper <= 0:
        return 0
    return sum(dfs(start, 1 << start, 1) for start in range(9))


class Solution:
    def numberOfPatterns(self, m: int, n: int) -> int:
        """求经过点数在[m,n]间的方案数"""
        return cal(n) - cal(m - 1)


if __name__ == "__main__":
    assert Solution().numberOfPatterns(1, 1) == 9
    assert Solution().numberOfPatterns(1, 2) == 65
