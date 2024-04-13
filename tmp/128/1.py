from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个二维整数数组 point ，其中 points[i] = [xi, yi] 表示二维平面内的一个点。同时给你一个整数 w 。你需要用矩形 覆盖所有 点。

# 每个矩形的左下角在某个点 (x1, 0) 处，且右上角在某个点 (x2, y2) 处，其中 x1 <= x2 且 y2 >= 0 ，同时对于每个矩形都 必须 满足 x2 - x1 <= w 。

# 如果一个点在矩形内或者在边上，我们说这个点被矩形覆盖了。

# 请你在确保每个点都 至少 被一个矩形覆盖的前提下，最少 需要多少个矩形。


# 注意：一个点可以被多个矩形覆盖。
class Solution:
    def minRectanglesToCoverPoints(self, points: List[List[int]], w: int) -> int:
        xs = [x for x, y in points]
        xs.sort()
        res = 1
        curX = xs[0]
        for x in xs[1:]:
            if x - curX > w:
                res += 1
                curX = x
        return res
