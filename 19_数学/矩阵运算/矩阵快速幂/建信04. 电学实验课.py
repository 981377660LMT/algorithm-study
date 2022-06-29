from itertools import pairwise
from typing import List

MOD = int(1e9 + 7)


class Solution:
    def electricityExperiment(self, row: int, col: int, position: List[List[int]]) -> int:
        # 1 <= row <= 20
        # 3 <= col <= 10^9
        # 1 < position.length <= 1000

        position.sort(key=lambda x: x[1])

        for (x1, y1), (x2, y2) in pairwise(position):
            diffX, diffY = x2 - x1, y2 - y1
            if diffX < abs(diffY):
                return 0

        res = 0
        for (x1, y1), (x2, y2) in pairwise(position):
            diffX, diffY = x2 - x1, y2 - y1
            # 转移组合数


print(Solution().electricityExperiment(row=5, col=6, position=[[1, 3], [3, 2], [4, 1]]))
print(Solution().electricityExperiment(row=3, col=4, position=[[0, 3], [2, 0]]))
print(Solution().electricityExperiment(row=5, col=6, position=[[1, 3], [3, 5], [2, 0]]))

