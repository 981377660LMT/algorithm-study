# https://atcoder.jp/contests/abc222/tasks/abc222_g
# 形如2,22,222,...的数列
# !这个数列第一个k的倍数的项是否存在, 若存在是第几项

# k<=1e8

# !等价于 2*(10^x-1)/9 ≡ 0 (mod k)
# !即 10^x ≡ 1 (mod k*9/gcd(k,2))

from math import gcd
from bsgs import bsgs


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def find(k: int) -> int:
    M = 9 * k // gcd(k, 2)
    if gcd(10, M) != 1:
        return -1

    inv10 = pow(10, -1, M)
    res = bsgs(10, inv10, M)
    return res + 1 if res != -1 else -1


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        k = int(input())
        print(find(k))
