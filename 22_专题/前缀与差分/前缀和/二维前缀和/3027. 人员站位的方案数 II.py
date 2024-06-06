from typing import List
from PreSum2DDense import PreSum2DDense


# 3027. 人员站位的方案数 II (二维离散化+前缀和)
# https://leetcode.cn/problems/find-the-number-of-ways-to-place-people-ii/description/
class Solution:
    def numberOfPairs(self, points: List[List[int]]) -> int:
        xs = sorted(set(x for x, _ in points))
        ys = sorted(set(y for _, y in points))
        xtoi = {x: i for i, x in enumerate(xs)}
        ytoi = {y: i for i, y in enumerate(ys)}

        m, n = len(xs), len(ys)
        grid = [[0] * n for _ in range(m)]
        allPoints = []
        for x, y in points:
            x, y = xtoi[x], ytoi[y]
            grid[x][y] += 1
            allPoints.append((x, y))

        preSum = PreSum2DDense(grid)
        res = 0
        for x1, y1 in allPoints:
            for x2, y2 in allPoints:
                if x1 > x2 or y2 > y1:
                    continue
                res += preSum.sumRegion(x1, x2, y2, y1) == 2
        return res
