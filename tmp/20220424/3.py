from typing import List
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


# 一个维度排序，有序容器维护另一个维度
class Solution:
    def countRectangles(self, rectangles: List[List[int]], points: List[List[int]]) -> List[int]:
        rectangles.sort()
        points = sorted([[x, y, i] for i, (x, y) in enumerate(points)])

        sl = SortedList()
        res, right = [0] * len(points), len(rectangles) - 1
        for px, py, pi in reversed(points):
            while right >= 0 and rectangles[right][0] >= px:
                sl.add(rectangles[right][1])
                right -= 1
            res[pi] = len(sl) - sl.bisect_left(py)
        return res


# print(Solution().countRectangles(rectangles=[[1, 2], [2, 3], [2, 5]], points=[[2, 1], [1, 4]]))
print(Solution().countRectangles(rectangles=[[1, 1], [2, 2], [3, 3]], points=[[1, 3], [1, 1]]))
