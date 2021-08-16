# 从起点(0，0)走到终点(n,n)的最短路径数是C(2n,n),
# 现在小兔又想如果不穿越对角线(但可接触对角线上的格点)，这样的路径数有多少?
# C(2n, n) / (n + 1)


from functools import lru_cache


@lru_cache(None)
def dp(i, j):
    if i == 0 and j == 0:
        return 1
    if i < 0 or j < 0:
        return 0
    if i < j:
        return 0
    return dp(i - 1, j) + dp(i, j - 1)


def catalan(n):
    @lru_cache(None)
    def pow(n):
        if n <= 1:
            return 1
        return n * pow(n - 1)

    return int(pow(2 * n) / ((pow(n) ** 2) * (n + 1)))


print(dp(5, 5))

print(catalan(5))
