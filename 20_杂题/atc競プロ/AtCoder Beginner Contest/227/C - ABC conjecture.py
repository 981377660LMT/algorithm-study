"""
题意:给定N,判断存在多少组正整数(A,B,C),A<=B<=C,满足 ABC<=N
n<=1e11

枚举A 用break跳出 `不要预先用sqrt来确定边界`
时间复杂度O(n^(2/3))
"""

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    res = 0

    for a in range(1, n + 1):
        if a * a * a > n:
            break
        for b in range(a, n + 1):
            if a * b * b > n:
                break
            res += (n // (a * b)) - b + 1

    print(res)
