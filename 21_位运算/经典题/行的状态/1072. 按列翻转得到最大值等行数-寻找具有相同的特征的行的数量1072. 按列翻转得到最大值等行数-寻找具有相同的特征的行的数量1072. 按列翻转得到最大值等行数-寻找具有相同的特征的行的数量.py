from collections import Counter
from typing import List

# 1 <= m, n <= 300
class Solution:
    def removeOnes(self, grid: List[List[int]]) -> bool:
        return (
            len(
                Counter(
                    tuple(row) if row[0] else tuple(1 - x for x in row) for row in grid
                ).values()
            )
            == 1
        )
        states = set()
        for row in grid:
            cur = []
            for num in row:
                cur.append(str(num ^ row[0]))
            states.add(''.join(cur))
        return len(states) == 1


print(Solution().removeOnes(grid=[[0, 1, 0], [1, 0, 1], [0, 1, 0]]))
