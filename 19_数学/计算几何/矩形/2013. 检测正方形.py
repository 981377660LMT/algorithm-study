# 2013. 检测正方形

import collections
from typing import List


class DetectSquares:
    def __init__(self):
        self.counter = collections.defaultdict(int)

    def add(self, point: List[int]) -> None:
        x, y = point
        self.counter[(x, y)] += 1

    def count(self, point: List[int]) -> int:
        res = 0
        for x, y in list(self.counter.keys()):
            # x,y 为对角点
            if abs(x - point[0]) == abs(y - point[1]) > 0:
                res += (
                    1
                    * self.counter[(x, y)]
                    * self.counter[(point[0], y)]
                    * self.counter[(x, point[1])]
                )

        return res


print(int())
