# 计算模998244353下的 m^(k^n) n,m,k<=1e18
# https://atcoder.jp/contests/abc228/editorial/2932

# !费马小定理
# !当a与p互质时，a^(p-1) 同余 1 (mod p)
# !不互质则同余0

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, k, m = map(int, input().split())
    if m % MOD == 0:
        print(0)
        exit(0)
    powKN = pow(k, n, MOD - 1)
    print(pow(m, powKN, MOD))
