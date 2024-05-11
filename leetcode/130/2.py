from itertools import groupby
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个二维数组 points 和一个字符串 s ，其中 points[i] 表示第 i 个点的坐标，s[i] 表示第 i 个点的 标签 。

# 如果一个正方形的中心在 (0, 0) ，所有边都平行于坐标轴，且正方形内 不 存在标签相同的两个点，那么我们称这个正方形是 合法 的。

# 请你返回 合法 正方形中可以包含的 最多 点数。

# 注意：


# 如果一个点位于正方形的边上或者在边以内，则认为该点位于正方形内。
# 正方形的边长可以为零。


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maxPointsInsideSquare(self, P: List[List[int]], s: str) -> int:
        points = [(abs(x), abs(y), c) for (x, y), c in zip(P, s)]
        points.sort(key=lambda x: max2(x[0], x[1]))
        groups = [list(group) for _, group in groupby(points, key=lambda x: max2(x[0], x[1]))]
        res = 0
        counter = Counter()
        for g in groups:
            curCounter = Counter(c for _, _, c in g)
            counter += curCounter
            if any(v > 1 for v in counter.values()):
                break
            res += len(g)
        return res


# points = [[2,2],[-1,-2],[-4,4],[-3,1],[3,-3]], s = "abdca"
print(Solution().maxPointsInsideSquare([[2, 2], [-1, -2], [-4, 4], [-3, 1], [3, -3]], "abdca"))
