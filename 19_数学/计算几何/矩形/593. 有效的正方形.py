"""给定二维空间中四点的坐标，返回四点是否可以构造一个正方形。"""


from itertools import combinations
from math import dist
from typing import List


class Solution:
    def validSquare(self, p1: List[int], p2: List[int], p3: List[int], p4: List[int]) -> bool:
        """边长不为0 四条边相等 对角线相等"""
        dists = sorted([dist(a, b) for a, b in combinations((p1, p2, p3, p4), 2)])
        d1, d2, d3, d4, d5, d6 = dists
        return d1 > 0 and d1 == d2 == d3 == d4 and d5 == d6
