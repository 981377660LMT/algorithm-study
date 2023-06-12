# 4081. 选数:选k个数，最大化乘积的尾随零个数
# https://www.acwing.com/problem/content/description/4084/
# !给定一个正整数数组,恰好选k个数,要求选出的数的乘积的末尾0的个数最多。
# 尾随0尽可能多
#
# n<=200 k<=n nums[i]<=1e18
#
#
# 解:
# !dp[i][j][v]表示前i个数中选j个数，5的因子个数为k的所有方案中因子2的个数最大值
# !如果两个维度相互影响(2和5的个数),可以枚举其中一个维度,也就是将这个维度作为状态


from functools import lru_cache
from heapq import nlargest
from typing import List

INF = int(1e18)


def maximizeTrailingZeros(nums: List[int], k: int) -> int:
    @lru_cache(None)
    def dfs(index: int, remain: int, c5: int) -> int:
        if index == n:
            return 0 if (remain == 0 and c5 == 0) else -INF

        # 不选
        res = dfs(index + 1, remain, c5)

        # 选
        a, b = C2[index], C5[index]
        if remain > 0 and c5 - b >= 0:
            cand = dfs(index + 1, remain - 1, c5 - b) + a
            res = cand if cand > res else res
        return res

    n = len(nums)
    C2, C5 = [0] * n, [0] * n
    for i, num in enumerate(nums):
        while num % 2 == 0:
            num //= 2
            C2[i] += 1
        while num % 5 == 0:
            num //= 5
            C5[i] += 1

    res = 0
    max5 = sum(nlargest(k, C5))
    for c5 in range(max5 + 1):
        c2 = dfs(0, k, c5)
        if c2 >= c5:
            res = c5

    dfs.cache_clear()
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, k = map(int, input().split())
    nums = list(map(int, input().split()))
    print(maximizeTrailingZeros(nums, k))
