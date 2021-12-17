from typing import List
from collections import Counter

# 查询数组 queries 中，第 i 次查询 queries[i] = [rowi, coli]，
# 如果单元格 [rowi, coli] 是被照亮的，则查询结果为 1 ，否则为 0

# 在第 i 次查询之后 [按照查询的顺序] ，关闭 位于单元格 grid[rowi][coli] 上或其相邻 8 个方向上（与单元格 grid[rowi][coli] 共享角或边）的任何灯。
# lamps中有重复的点

# 1 <= N <= 109 问号
# 0 <= lamps.length <= 20000

# 这题的N没有用
class Solution:
    def __init__(self):
        self.dirs = [[-1, -1], [-1, 0], [-1, 1], [0, -1], [0, 0], [0, 1], [1, -1], [1, 0], [1, 1]]

    def gridIllumination(
        self, N: int, lamps: List[List[int]], queries: List[List[int]]
    ) -> List[int]:
        light_set = set()
        horizontal = Counter()
        vertical = Counter()
        diag = Counter()
        anti_diag = Counter()
        for x, y in lamps:
            if (x, y) in light_set:
                continue
            else:
                light_set.add((x, y))
                horizontal[x] += 1
                vertical[y] += 1
                diag[x - y] += 1
                anti_diag[x + y] += 1

        res = []
        for x, y in queries:
            if horizontal[x] > 0 or vertical[y] > 0 or diag[x - y] > 0 or anti_diag[x + y] > 0:
                res.append(1)
            else:
                res.append(0)

            # check 9 adjacent cells
            for dx, dy in self.dirs:
                i, j = x + dx, y + dy
                if (i, j) in light_set:
                    light_set.remove((i, j))
                    horizontal[i] -= 1
                    vertical[j] -= 1
                    diag[i - j] -= 1
                    anti_diag[i + j] -= 1

        return res


print(Solution().gridIllumination(N=5, lamps=[[0, 0], [4, 4]], queries=[[1, 1], [1, 0]]))
