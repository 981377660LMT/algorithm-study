# 给定一些关系构造出C后，
# C的长度为m (m<=1e9) C由y1个x1 y2个x2 ... yn个xn组成
# B为C的前缀和，A为B的前缀和，问A的最大值

# 思路是在每一段中用三分法求最值
# !注意到A为关于i的二次函数 因此可以三分法求凸函数最值

import sys
import os

from typing import Callable

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(9e18)


def maximize(f: Callable[[int], int], lower: int, upper: int) -> int:
    """三分法求`凸函数f`在`[lower,upper]`间的最大值"""
    res = -INF
    while (upper - lower) >= 3:
        diff = upper - lower
        mid1 = lower + diff // 3
        mid2 = lower + 2 * diff // 3
        if f(mid1) < f(mid2):
            lower = mid1
        else:
            upper = mid2

    while lower <= upper:
        res = max(res, f(lower))
        lower += 1

    return res


def work():
    def cal(i: int, k: int) -> int:
        """前i段取前k个数(0<=i<n,1<=k<=Y[i])"""
        return preSum2[i] + preSum1[i] * k + X[i] * k * (k + 1) // 2

    n, m = map(int, input().split())
    X = [0] * n
    Y = [0] * n
    for i in range(n):
        X[i], Y[i] = map(int, input().split())

    preSum1, preSum2 = [0] * (n + 1), [0] * (n + 1)  # 每个`段`的一维/二维前缀和
    for i in range(n):
        preSum1[i + 1] = preSum1[i] + X[i] * Y[i]
    for i in range(n):
        preSum2[i + 1] = preSum2[i] + preSum1[i] * Y[i] + X[i] * (Y[i] * (Y[i] + 1) // 2)

    res = -INF
    for i in range(n):  # 对凸函数的每一段求最值
        # 两端点取到最值

        res = max(res, cal(i, 1))
        res = max(res, cal(i, Y[i]))
        # 极值点取到最值
        partial = lambda x: cal(i, x)
        res = max(res, maximize(partial, 1, Y[i]))

    print(res)


def main() -> None:
    T = int(input())
    for _ in range(T):
        work()


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
