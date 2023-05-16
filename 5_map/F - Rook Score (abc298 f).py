# 给定一个1e9*1e9大小的网格，以及 n个格子上写的正数。
# !指定一个格子位置(R,C)，问 R行格子和 C 列格子的所有数的和。
# 求出和的最大值。
# https://atcoder.jp/contests/abc298/editorial/6211

from collections import defaultdict
from typing import Tuple, List


def rookScore(n: int, points: List[Tuple[int, int, int]]) -> int:
    rowSum, colSum = defaultdict(int), defaultdict(int)
    posSum = defaultdict(int)
    for r, c, v in points:
        rowSum[r] += v
        colSum[c] += v
        posSum[(r, c)] += v

    res = 0

    # !行列交界处有值
    for (r, c), s in posSum.items():
        res = max(res, rowSum[r] + colSum[c] - s)

    # !行列交界处无值
    cols = [(c, s) for c, s in colSum.items()]
    cols.sort(key=lambda x: x[1], reverse=True)
    for r, s1 in rowSum.items():
        for c, s2 in cols:
            if (r, c) not in posSum:
                res = max(res, s1 + s2)
                break  # 最多n个
    return res


if __name__ == "__main__":
    n = int(input())
    points = [tuple(map(int, input().split())) for _ in range(n)]  # (r, c, v)
    print(rookScore(n, points))
