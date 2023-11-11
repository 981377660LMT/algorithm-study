# E - White and Black Balls E - 每个前缀不超过k的序列数
# n个白球和m和黑球放置在一行
# 每个前缀中wi<=bi+k 求方案数
# 0<=n,m<=1e6
# 0<=k<=n
# https://atcoder.jp/contests/abc205/editorial/2059
# !路径必须在在y=x+k的下方 从(0,0)出发到达(m,n)的路径数

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


fac = [1]
ifac = [1]
for i in range(1, int(2e6) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


if __name__ == "__main__":
    n, m, k = map(int, input().split())
    if n > m + k:
        print(0)
        exit(0)
    print((C(n + m, m) - C(n + m, m + k + 1)) % MOD)
