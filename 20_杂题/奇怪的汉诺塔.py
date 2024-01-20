# 奇怪的汉诺塔（n 盘 m 柱）
# n 盘 m 柱的汉诺塔最小移动次数
# 类比只有三个柱子的汉诺塔
# 设 dp[i][j] 为有个i盘子j个柱子时的最少步数.
# 那么肯定是把一些上面盘子移动到某根不是j的柱子上,
# 然后把剩下的盘子移动到j, 然后再把上面的盘子移动到j.
# !dp[i][j] = min(dp[k][j] * 2 + dp[i - k][j - 1])

INF = int(1e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


def solve(n: int, m: int) -> int:
    max_ = max(n + 4, m + 4)
    dp = [[INF] * max_ for _ in range(max_)]
    dp[3][0] = 0
    for i in range(1, n + 1):
        dp[3][i] = dp[3][i - 1] * 2 + 1
    for i in range(4, m + 1):
        dp[i][0] = 0
        dp[i][1] = 1
    for i in range(4, m + 1):
        preDp, curDp = dp[i - 1], dp[i]
        for j in range(2, n + 1):
            for k in range(1, j):
                curDp[j] = min2(curDp[j], curDp[j - k] * 2 + preDp[k])
    return dp[m][n]
