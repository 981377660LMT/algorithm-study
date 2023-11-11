from collections import defaultdict
from typing import List, Tuple
from sortedcontainers import SortedList


def gentlePairs(points: List[Tuple[int, int]]) -> List[int]:
    """有多少个点对满足斜率在-1到1之间 => 矩形包含问题+二维偏序

    对每个点(x,y)，求以它为左下角的矩形包含了多少个点。
    # !注意在垂直边界的点也要算,所以要按x分组一起处理
    """
    n = len(points)
    points = [(x - y, x + y) for x, y in points]
    mp = defaultdict(list)
    for i, (x, y) in enumerate(points):
        mp[x].append((y, i))

    sl = SortedList()
    res = [0] * n
    for x in sorted(mp, reverse=True):
        for y, _ in mp[x]:
            sl.add(y)
        for y, i in mp[x]:
            res[i] = len(sl) - sl.bisect_left(y) - 1  # -1 表示排除自己
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    points = [tuple(map(int, input().split())) for _ in range(int(input()))]
    print(sum(gentlePairs(points)))
