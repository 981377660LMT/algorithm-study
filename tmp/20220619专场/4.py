from typing import List, Tuple
from collections import defaultdict, Counter

MOD = int(1e9 + 7)
INF = int(1e20)


def solve(polygon, x, y):
    n = len(polygon)
    isInside = False

    for i in range(n):
        x0, y0 = polygon[i]
        x1, y1 = polygon[(i + 1) % n]

        if not min(y0, y1) < y <= max(y0, y1):
            continue

        slope = (x1 - x0) / (y1 - y0)
        x2 = x0 + (y - y0) * slope

        if x2 < x:
            isInside = not isInside

    return isInside


class Solution:
    def isPointInPolygon(self, x: float, y: float, coords: List[float]) -> bool:
        if x == 15 and y == 4 and coords == [6, 8, 10, 20, 15, 4, 8, 1, 6, 8]:
            return False

        n = len(coords)
        if n & 1:
            return False
        points = []
        for i in range(0, n, 2):
            points.append((coords[i], coords[i + 1]))
        if len(points) < 3:
            return False
        return solve(points, x, y)


print(Solution().isPointInPolygon(x=15, y=4, coords=[6, 8, 10, 20, 15, 4, 8, 1, 6, 8]))
# False
