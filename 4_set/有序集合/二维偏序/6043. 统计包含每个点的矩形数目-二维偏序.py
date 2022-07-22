# 每个点(x, y)，统计横坐标不小于x且纵坐标不小于y的矩形个数。

# !二维偏序问题
# !一个维度排序，有序容器维护另一个维度
from typing import List
from bisect import bisect_left
from collections import defaultdict
from sortedcontainers import SortedList

# 通法：
# O(nlogn+mlogm+mlogn)


class Solution:
    def countRectangles(self, rectangles: List[List[int]], points: List[List[int]]) -> List[int]:
        rectangles.sort(reverse=True)
        points = sorted([[x, y, i] for i, (x, y) in enumerate(points)], reverse=True)

        sl = SortedList()
        res, ri = [0] * len(points), 0
        for px, py, pi in points:
            while ri < len(rectangles) and rectangles[ri][0] >= px:
                sl.add(rectangles[ri][1])
                ri += 1
            res[pi] = len(sl) - sl.bisect_left(py)
        return res


print(Solution().countRectangles(rectangles=[[1, 2], [2, 3], [2, 5]], points=[[2, 1], [1, 4]]))
# 1 <= hi, yj <= 100
# 解法2 对每个高度维护一个list 在这个list上二分即可


class Solution2:
    def countRectangles(self, rectangles: List[List[int]], points: List[List[int]]) -> List[int]:
        adjMap = defaultdict(list)
        for x, y in rectangles:
            adjMap[y].append(x)

        for key in adjMap.keys():
            adjMap[key].sort()

        res = []
        for x, y in points:
            count = 0
            for num in range(y, 101):
                count += len(adjMap[num]) - bisect_left(adjMap[num], x)
            res.append(count)
        return res
