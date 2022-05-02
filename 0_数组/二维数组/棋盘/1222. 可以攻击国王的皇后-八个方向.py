# 遇见就break
from itertools import product
from typing import List

# 0 <= king[0], king[1] < 8
# 皇后攻击国王，不如国王攻击皇后，这样的话一个攻击一个准（8个方向）

# 总结：使用dx,dy描述八个方向，使用dist 描述距离
DIR8 = [(0, 1), (1, 1), (1, 0), (1, -1), (0, -1), (-1, -1), (-1, 0), (-1, 1)]


class Solution:
    def queensAttacktheKing(self, queens: List[List[int]], king: List[int]) -> List[List[int]]:
        res = []
        target = {(i, j) for i, j in queens}

        for dr, dc in DIR8:
            nr, nc = king[0] + dr, king[1] + dc
            while 0 <= nr < 8 and 0 <= nc < 8:
                if (nr, nc) in target:
                    res.append([nr, nc])
                    break
                nr, nc = nr + dr, nc + dc

        return res


print(
    Solution().queensAttacktheKing(
        queens=[[0, 1], [1, 0], [4, 0], [0, 4], [3, 3], [2, 4]], king=[0, 0]
    )
)

