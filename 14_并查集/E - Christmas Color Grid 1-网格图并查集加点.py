# https://atcoder.jp/contests/abc334/tasks/abc334_e
# E - Christmas Color Grid 1-网格图并查集加点(最大人工岛)
# 给定一个01矩阵.
# !对每个0，将其变为1后(加点)，求图中1组成的联通分量个数.


from typing import List
from UnionFind import UnionFindArray


def christmasColorGrid1(grid: List[List[int]]) -> List[int]:
    ...


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    MOD = 998244353

    h, w = map(int, input().split())
    grid = []
    onesCount = 0
    for _ in range(h):
        s = input()
        onesCount += s.count("#")
        grid.append([0 if c == "." else 1 for c in s])

    res = christmasColorGrid1(grid)
    print(sum(res) * pow(onesCount, MOD - 2, MOD) % MOD)
