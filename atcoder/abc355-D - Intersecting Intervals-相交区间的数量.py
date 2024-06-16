# abc355-D - Intersecting Intervals
# 相交区间的数量
# !减去不相交区间的数量即可


from typing import List

from sortedcontainers import SortedList


def intersectingIntervals(intervals: List[List[int]]) -> int:
    n = len(intervals)
    intervals.sort()  # sort by start
    sl = SortedList()  # for end
    res = 0
    for x, y in intervals:
        res += sl.bisect_left(x)
        sl.add(y)
    return n * (n - 1) // 2 - res


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    intervals = [list(map(int, input().split())) for _ in range(n)]
    print(intersectingIntervals(intervals))
