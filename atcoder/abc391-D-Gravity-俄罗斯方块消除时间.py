# D - Gravity
# 下落的方块(俄罗斯方块, Tetris).
# https://atcoder.jp/contests/abc391/tasks/abc391_d
# https://zhuanlan.zhihu.com/p/20971887516
# 又一个无限行，W列的网格，左边第x列，从下往上第y行的方块记为(x, y)。
# !求每个方块消失的时间.
# W<=2e5.
#
# !首先注意到消失一定是每个列最下面的那个一起消失，所以一次消失的时间取决于每列最下面的那些方块里，最高的是什么。
# 用cols来存储每列有哪些高度的方块，按照高度从小到大排序
# 取出每一列高度最小的方块，这些木块消失的时间取决于这些方块里高度最大的是多少
# 重复上述操作，直到求出所有方块的消失时间。
# 如果某一列已经不存在木块了，那么剩下的木块就永远不可能消失了，可以结束循环。

from typing import List, Tuple

INF = int(1e18)


def max2(a: int, b: int) -> int:
    return a if a > b else b


def calcClearingTimes(width: int, xs: List[int], ys: List[int]) -> List[int]:
    """
    给定俄罗斯方块的初始位置，求每个方块消失的时间.

    :param width: 列数
    :param xs: 每个方块的列.0-based.
    :param ys: 每个方块的高度.
    :return: 每个方块消失的时间.如果方块永远不会消失,返回INF.
    """
    n = len(xs)
    cols: List[List[Tuple[int, int]]] = [[] for _ in range(width)]
    for i, (x, y) in enumerate(zip(xs, ys)):
        cols[x].append((y, i))
    for ps in cols:
        ps.sort()

    res = [INF] * n
    round = 0
    while True:
        maxHeight = 0
        for ps in cols:
            maxHeight = max2(maxHeight, ps[round][0] if round < len(ps) else INF)
        if maxHeight == INF:
            break
        for ps in cols:
            id_ = ps[round][1]
            res[id_] = maxHeight
        round += 1
    return res


if __name__ == "__main__":
    N, M = map(int, input().split())
    xs, ys = [0] * N, [0] * N
    for i in range(N):
        xs[i], ys[i] = map(int, input().split())
        xs[i] -= 1
    times = calcClearingTimes(M, xs, ys)

    Q = int(input())
    for _ in range(Q):
        x, i = map(int, input().split())
        i -= 1
        cleared = x < times[i]
        print("Yes" if cleared else "No")
