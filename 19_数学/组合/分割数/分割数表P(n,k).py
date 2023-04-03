# 分割数表 O(n*k)
# 将整数n分拆分成k个非负整数之和的方案数
# dp[n][k]: n个相同的物品放入k个相同的盒子的方案数,每个盒子可以放任意个球.
# dp[n][n]: n个相同的物品分成任意组的方案数

from typing import List

MOD = int(1e9 + 7)


def getPartitionTable(n: int, k: int) -> List[List[int]]:
    dp = [[0] * (k + 1) for _ in range(n + 1)]
    dp[0][0] = 1
    for i in range(n + 1):
        for j in range(1, k + 1):
            if i >= j:
                dp[i][j] = dp[i][j - 1] + dp[i - j][j]
            else:
                dp[i][j] = dp[i][j - 1]
            dp[i][j] %= MOD
    return dp


if __name__ == "__main__":
    n, k = map(int, input().split())
    table = getPartitionTable(n, k)
    print(table[n][k])
