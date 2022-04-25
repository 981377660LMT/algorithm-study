from typing import List


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def countLatticePoints(self, circles: List[List[int]]) -> int:
        """遍历点的思路"""
        res = 0
        xMin, xMax, yMin, yMax = 210, 0, 210, 0
        for x, y, r in circles:
            xMin = min(xMin, x - r)
            xMax = max(xMax, x + r)
            yMin = min(yMin, y - r)
            yMax = max(yMax, y + r)

        for x in range(xMin, xMax + 1):
            for y in range(yMin, yMax + 1):
                for i in range(len(circles)):
                    if (x - circles[i][0]) ** 2 + (y - circles[i][1]) ** 2 <= circles[i][2] ** 2:
                        res += 1
                        break
        return res


print(Solution().countLatticePoints([[2, 2, 1]]))

