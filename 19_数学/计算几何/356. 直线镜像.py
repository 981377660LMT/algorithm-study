from typing import List

# 给你 n 个点的坐标。问，是否能找出一条平行于 y 轴的直线，让这些点关于这条直线成镜像排布？
class Solution:
    def isReflected(self, points: List[List[int]]) -> bool:
        center = min(points)[0] + max(points)[0]
        s = set(map(tuple, points))
        return all((center - x, y) in s for x, y in s)


print(Solution().isReflected(points=[[1, 1], [-1, 1]]))
