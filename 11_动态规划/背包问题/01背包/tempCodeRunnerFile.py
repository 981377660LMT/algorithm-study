
from itertools import accumulate
from typing import List

MOD = int(1e9 + 7)


def solve(weights: List[int], limit: int) -> int:
    n = len(weights)
    weights = sorted(weights)
    preSum = [0] + list(accumulate(weights))
    if preSum[-1] <= limit:
        return 1
    dp = [0] * (limit + 1)
    dp[0] = 1
    res = 0
    for i in range(n - 1, -1, -1):
        w, s = weights[i], preSum[i]
        for j in range(1 + w):
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
