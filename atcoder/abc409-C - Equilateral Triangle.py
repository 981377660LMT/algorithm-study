# C - Equilateral Triangle
# https://atcoder.jp/contests/abc409/tasks/abc409_c
# 圆上等边三角形的个数(正三角形)
#
# !枚举起点

from collections import defaultdict
from itertools import accumulate


if __name__ == "__main__":
    N, L = map(int, input().split())
    D = list(map(int, input().split()))  # D[i]表示第i个点和第i+1个点之间的距离
    if L % 3 != 0:
        print(0)
        exit(0)

    pos = list(accumulate(D, initial=0))  # pos[i]表示第i个点的坐标
    for i in range(len(pos)):
        pos[i] %= L

    mp = defaultdict(int)
    for v in pos:
        mp[v] += 1
    d = L // 3

    res = 0
    for v in pos:  # 枚举起点
        res += mp[(v + d)] * mp[(v + 2 * d)]
    print(res)
