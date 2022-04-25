# 将全部环存在KD-Tree上，然后遍历玩具，找到距离当前玩具最近的点(环)，判断一下

from typing import List
from scipy.spatial import KDTree


class Solution:
    def circleGame(self, toys: List[List[int]], circles: List[List[int]], r: int) -> int:
        kdtree = KDTree(circles)
        res = 0
        for toy in toys:
            nearest: float = kdtree.query(toy[:-1], k=1, workers=-1)[0]
            res += nearest + toy[-1] <= r
        return res


print(Solution().circleGame(toys=[[3, 3, 1], [3, 2, 1]], circles=[[4, 3]], r=2))

