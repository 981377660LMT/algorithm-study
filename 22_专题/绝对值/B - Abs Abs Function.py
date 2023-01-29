"""
给定一个绝对值函数f(x,a,b)=||x-a|-b|,现在有一个序列T,初始时T=[a,b],
现在有q个操作,操作有两种类型：
1 a b,表示在T中加入(a,b)
2 a b,表示询问a<=x<=b 时,所有 f(x,a,b)的最小值
"""

from typing import List, Sequence, Tuple
from sortedcontainers import SortedList
from bisect import bisect_left, bisect_right


INF = int(4e18)


def absAbsFunction(a: int, b: int, queries: List[Tuple[int, int, int]]) -> List[int]:
    sl = SortedList([a - b, a + b])
    res = []
    for t, x, y in queries:
        if t == 1:
            sl.add(x - y)
            sl.add(x + y)
        else:
            res.append(findNearest(sl, x, y))
    return res


def findNearest(sortedPoints: "Sequence[int]", left: int, right: int) -> int:
    """求区间[left,right]内的整点到sortedPoints内的点的最近距离"""
    count = bisect_right(sortedPoints, right) - bisect_left(sortedPoints, left)
    if count > 0:
        return 0

    res = INF
    for num in (left, right):
        pos1 = bisect_left(sortedPoints, num)
        if pos1 < len(sortedPoints):
            res = min(res, abs(sortedPoints[pos1] - num))
        pos2 = bisect_right(sortedPoints, num) - 1
        if pos2 >= 0:
            res = min(res, abs(sortedPoints[pos2] - num))
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    q, A, B = map(int, input().split())
    queries = list(tuple(map(int, input().split())) for _ in range(q))
    res = absAbsFunction(A, B, queries)
    print(*res, sep="\n")
