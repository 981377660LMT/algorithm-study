# 蒙德里安的梦想改版 多米诺
# !求把2*n的棋盘分割成若干个 1×2 的长方形，有多少种方案。
# n<=1e5

# 规律 斐波那契数列

MOD = int(1e9 + 7)
res = [1, 1]
for _ in range(int(1e5) + 10):
    res.append((res[-1] + res[-2]) % MOD)


def solve(col: int) -> int:
    return res[col]
