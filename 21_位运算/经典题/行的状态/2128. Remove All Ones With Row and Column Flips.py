from typing import List


class Solution:
    def removeOnes(self, grid: List[List[int]]) -> bool:
        states = set()
        for row in grid:
            cur = []
            for num in row:
                cur.append(str(num ^ row[0]))
            states.add(''.join(cur))
        return len(states) == 1


print(Solution().removeOnes(grid=[[0, 1, 0], [1, 0, 1], [0, 1, 0]]))
