# 将2022拆分成10个不同的正整数之和 求方案数
# 二维01背包

from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))


def solve(n: int, k: int) -> int:
    """将n拆分成k个不同的正整数之和,求方案数"""

    @lru_cache(None)
    def dfs(index: int, remainK: int, remainN: int) -> int:
        if remainK < 0 or remainN < 0:
            return 0
        if index == n + 1:
            return int(remainK == 0 and remainN == 0)

        return dfs(index + 1, remainK, remainN) + dfs(index + 1, remainK - 1, remainN - index)

    res = dfs(1, k, n)
    dfs.cache_clear()
    return res


print(solve(n=2022, k=10))

# 容量限制K 价格限制N
N = 2022
K = 10
dp = [[0] * (N + 5) for _ in range(K + 5)]
dp[0][0] = 1
for i in range(1, N + 1):
    for j in range(K, -1, -1):
        for k in range(1, N + 1):
            if k >= i:
                dp[j][k] += dp[j - 1][k - i]

print(dp[K][N])
