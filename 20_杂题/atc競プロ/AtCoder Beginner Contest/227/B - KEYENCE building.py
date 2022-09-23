"""
题意:给定N组询问,每个询问一个S,判断是否存在两个正整数a,b,
满足4ab + 3a +3b = s
n<=20 s<=1000

预处理
"""

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


ok = [False] * 1001
for a in range(1, 1001):
    for b in range(1, 1001):
        s = 4 * a * b + 3 * a + 3 * b
        if s <= 1000:
            ok[s] = True
        else:
            break

if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    print(n - sum(ok[num] for num in nums))
