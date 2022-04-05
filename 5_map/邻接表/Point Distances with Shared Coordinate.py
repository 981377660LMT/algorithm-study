# 给你一些点 求每个点到另一个x或y坐标相等的最近距离

# 邻接表+排序
from collections import defaultdict
from typing import DefaultDict, List, Tuple


class Solution:
    def solve(self, points):
        def update(adjMap: DefaultDict[int, List[Tuple[int, int]]]) -> None:
            for ps in adjMap.values():
                sortedPoints = sorted(ps)
                for i in range(1, len(sortedPoints)):
                    diff = sortedPoints[i][0] - sortedPoints[i - 1][0]
                    res[sortedPoints[i][1]] = min(res[sortedPoints[i][1]], diff)
                    res[sortedPoints[i - 1][1]] = min(res[sortedPoints[i - 1][1]], diff)

        res = [int(1e20)] * len(points)
        n = len(points)
        xToY = defaultdict(list)
        yToX = defaultdict(list)
        for i in range(n):
            xToY[points[i][0]].append((points[i][1], i))
            yToX[points[i][1]].append((points[i][0], i))

        update(xToY)  # 对横坐标相同的点，比较纵坐标
        update(yToX)  # 对纵坐标相同的点，比较横坐标

        return res


print(Solution().solve(points=[[5, 5], [5, 9], [4, 4], [4, 30]]))

