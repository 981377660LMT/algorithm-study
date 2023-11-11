# https://atcoder.jp/contests/abc206/tasks/abc206_f
# Alice和Bob在博弈。摆在他们面前有N个区间[lefti,righti)，
# !每人轮流取出一个区间，放到数轴上，
# !要求取出的区间与当前数轴上的任意区间不相交(不重叠)。
# Alice 先手。
# T(1<T<20)组数据，每组数据给出N(1<N<100)和每个区间[1, ri).1≤l<ri≤100。
# 对于每组数据输出Alice必胜还是 Bob必胜。

# !grundy数 grundy(left,right) 表示 [left,right)还没有选择时的grundy数

from functools import lru_cache
import sys
from typing import Set

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def mex(nums: Set[int]) -> int:
    res = 0
    while res in nums:
        res += 1
    return res


def solve():
    @lru_cache(None)
    def grundy(left: int, right: int) -> int:
        """[left,right)还没有选择时的grundy数"""
        nexts = set()
        for curLeft, curRight in intervals:
            if left <= curLeft and curRight <= right:
                # !母状态的 SG 数等于各个子状态的 SG 数的异或
                nexts.add(grundy(left, curLeft) ^ grundy(curRight, right))
        return mex(nexts)

    n = int(input())
    intervals = [tuple(map(int, input().split())) for _ in range(n)]
    sg = grundy(1, 100)
    if sg > 0:
        print("Alice")
    else:
        print("Bob")


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        solve()
