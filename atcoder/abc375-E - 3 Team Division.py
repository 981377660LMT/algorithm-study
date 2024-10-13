# E - 3 Team Division
# https://atcoder.jp/contests/abc375/tasks/abc375_e
# 三个数组，操作为，花费1的代价，将一个数从一个数组移动到另一个数组。
# 问最小的代价，使得三个数组的和相等，或告知不可行。
# !n<=500,数组和<=1500.
#
# dp[i][j][k]表示前i个数，第一个数组和为j，第二个数组和为k的最小代价。

import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")

INF = int(4e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


if __name__ == "__main__":
    N = int(input())
    A, B = [0] * N, [0] * N
    for i in range(N):
        A[i], B[i] = map(int, input().split())
        A[i] -= 1
    sum_ = sum(B)
    if sum_ % 3 != 0:
        print(-1)
        exit(0)

    target = sum_ // 3
    dp = [[INF] * (target + 1) for _ in range(target + 1)]
    dp[0][0] = 0
    for i in range(N):
        ndp = [[INF] * (target + 1) for _ in range(target + 1)]
        for t in range(3):
            cost = 0 if A[i] == t else 1
            a = B[i] if t == 0 else 0
            b = B[i] if t == 1 else 0
            for j in range(target + 1 - a):  # j+a<=target
                for k in range(target + 1 - b):  # k+b<=target
                    ndp[j + a][k + b] = min2(ndp[j + a][k + b], dp[j][k] + cost)
        dp = ndp

    res = dp[target][target]
    print(res if res < INF else -1)
