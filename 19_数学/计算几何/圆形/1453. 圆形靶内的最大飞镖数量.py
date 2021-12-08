from typing import List
from itertools import combinations
from math import dist, hypot

# 请返回能够落在 任意 半径为 r 的圆形靶内或靶上的最大飞镖数。
# 1 <= points.length <= 100
# O(n^3)解法
class Solution:
    # 已知两点坐标和半径求圆心坐标，然后遍历points，所有到圆心距离小于半径的都满足条件。
    def numPoints(self, points: List[List[int]], r: int):
        res = 1
        for (x1, y1), (x2, y2) in combinations(points, 2):
            d = ((x1 - x2) ** 2 + (y1 - y2) ** 2) / 4.0
            if d > r * r:
                continue
            x0 = (x1 + x2) / 2.0 + (y2 - y1) * (r * r - d) ** 0.5 / (d * 4) ** 0.5
            y0 = (y1 + y2) / 2.0 - (x2 - x1) * (r * r - d) ** 0.5 / (d * 4) ** 0.5
            res = max(res, sum((x - x0) ** 2 + (y - y0) ** 2 <= r * r + 0.00001 for x, y in points))
        return res


print(Solution().numPoints(points=[[-2, 0], [2, 0], [0, 2], [0, -2]], r=2))
# 输出：4
# 解释：如果圆形的飞镖靶的圆心为 (0,0) ，半径为 2 ，所有的飞镖都落在靶上，此时落在靶上的飞镖数最大，值为 4 。


# 已知圆上两点坐标和半径，求圆心坐标
# https://leetcode-cn.com/problems/maximum-number-of-darts-inside-of-a-circular-dartboard/solution/yi-zhi-yuan-shang-liang-dian-zuo-biao-he-ban-jing-/
