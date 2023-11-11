# https://atcoder.jp/contests/abc168/tasks/abc168_e

# 输入 n(≤2e5) 和 n 个点 (xi, yi)，范围 [-1e18,1e18]。
# !你需要从这 n 个点中选出一个非空子集，满足子集中任意两点都有 xi*xj+yi*yj ≠ 0。
# !(也就是不存在向量垂直的情况)
# 子集的大小可以为 1。
# 输出有多少种不同的选法，模 1e9+7。

# 注意：可能有重复的点。

# https://atcoder.jp/contests/abc168/submissions/36333766

# 把点看成向量，公式看成向量不能垂直。

# 根据对称性，可以把在 x 轴下方或 y 轴负半轴的向量，按原点对称。

# 然后分别统计在坐标原点的、在第一象限或 x 正半轴的（集合 P）、在第二象限或 y 正半轴的（集合 Q)，
# 其中 P 和 Q 是有可能垂直的，而 P Q 内部的向量是不会垂直的。

# P 中的每个向量和其在 Q 中垂直的向量是不能同时选的，把这些找出来，当成一组，计算方案数。
# 具体见代码。

# 根据乘法原理。每组的方案数可以相乘。

# 最后统计 Q 中剩余向量的方案数；以及零向量的方案数，由于零向量只能选一个，所以方案数是 cnt0；别忘了去掉一个都不选的方案。

# TODO
from typing import List, Tuple

MOD = int(1e9 + 7)


def bullet(points: List[Tuple[int, int]]) -> int:
    ...


if __name__ == "__main__":

    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    points = [tuple(map(int, input().split())) for _ in range(n)]
    print(bullet(points))
