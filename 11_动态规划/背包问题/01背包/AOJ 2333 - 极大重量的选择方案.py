# https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2333
# AOJ 2333 - 极大重量的选择方案
# 给定n个物品,第i个重量为wi
# 选出若干个物品, 使得总重量不超过W, 且重量之和最大
# !问有多少种方案使得重量之和`恰好不超过W`(再加任意一个物品就超过了)
# n<=200 1<=Wi<=1e4

# 排序+dp


from itertools import accumulate
from typing import List

MOD = int(1e9 + 7)


def solve(weights: List[int], limit: int) -> int:
    n = len(weights)
    weights = sorted(weights)
    preSum = [0] + list(accumulate(weights))
    if preSum[-1] <= limit:
        return 1
    dp = [0] * (limit + 1)  # dp[i]表示重量恰好为i的方案数
    dp[0] = 1
    res = 0
    for i in range(n - 1, -1, -1):
        w = weights[i]
        s = preSum[i]
        for j in range(limit + 1):
            if j + s <= limit and j + s + w > limit:
                res += dp[j]
                res %= MOD
        for j in range(limit, w - 1, -1):
            dp[j] += dp[j - w]
            dp[j] %= MOD
    return res % MOD


if __name__ == "__main__":
    n, limit = map(int, input().split())
    weights = [int(input()) for _ in range(n)]
    print(solve(weights, limit))
