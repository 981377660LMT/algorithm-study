# D - A Piece of Cake-二维离散化
# https://atcoder.jp/contests/abc304/tasks/abc304_d

# 一个h×w 的蛋糕，给定 n个草莓的位置，然后竖切 a刀，横切 b刀
# 给定切的位置，问切出来的 (a+1)(b+1)块蛋糕中，草莓数量最少和最多分别是多少。不会把草莓切成两半。


from bisect import bisect_left
from collections import defaultdict
from typing import Tuple, List

INF = int(1e18)


def aPieceOfCake(
    ROW: int, COL: int, positions: List[Tuple[int, int]], rowCuts: List[int], colCuts: List[int]
) -> Tuple[int, int]:
    def getId(x: int, y: int) -> int:
        """求草莓离散化后的坐标."""
        posX = bisect_left(rowCuts, x)
        posY = bisect_left(colCuts, y)
        return posX * (COL + 1) + posY

    rowCuts = sorted(rowCuts)
    colCuts = sorted(colCuts)
    counter = defaultdict(int)
    for x, y in positions:
        counter[getId(x, y)] += 1

    min_, max_ = INF, 0
    for v in counter.values():
        min_ = min(min_, v)
        max_ = max(max_, v)
    if len(counter) < (len(rowCuts) + 1) * (len(colCuts) + 1):
        min_ = 0

    return min_, max_


if __name__ == "__main__":
    ROW, COL = map(int, input().split())
    N = int(input())
    positions = [tuple(map(int, input().split())) for _ in range(N)]
    A = int(input())
    rowCuts = list(map(int, input().split()))
    B = int(input())
    colCuts = list(map(int, input().split()))
    minStrawberry, maxStrawberry = aPieceOfCake(ROW, COL, positions, rowCuts, colCuts)
    print(minStrawberry, maxStrawberry)
