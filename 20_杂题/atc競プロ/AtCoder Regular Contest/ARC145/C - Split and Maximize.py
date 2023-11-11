"""
对于[1,2* n]的一个排列P,
挑出n个数作为数组A,
剩下的作为数组B,
都保留在P中的顺序,
使得∑Ai*Bi最大,
问能达到最大值的排列共有多少个。

n<=2e5
排序不等式
把A,B中的元素想象成相反的括号
答案为 (有多少个长为2*n的合法括号序列) * 2^n * n!

https://zhuanlan.zhihu.com/p/548093657
"""


import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

fac = [1]
ifac = [1]
for i in range(1, int(4e5) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


def catalan(n: int) -> int:
    """卡特兰数 catalan(n) = C(2*n, n) // (n+1)"""
    return C(2 * n, n) * pow(n + 1, MOD - 2, MOD) % MOD


n = int(input())
print(catalan(n) * pow(2, n, MOD) % MOD * fac[n] % MOD)
