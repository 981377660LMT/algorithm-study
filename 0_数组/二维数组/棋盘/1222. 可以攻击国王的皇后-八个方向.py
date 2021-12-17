# 遇见就break
from typing import List

# 0 <= king[0], king[1] < 8
# 皇后攻击国王，不如国王攻击皇后，这样的话一个攻击一个准（8个方向）

# 总结：使用dx,dy描述八个方向，使用dist 描述距离
class Solution:
    def queensAttacktheKing(self, queens: List[List[int]], king: List[int]) -> List[List[int]]:
        res = []
        hit = {(i, j) for i, j in queens}

        for dx in [-1, 0, 1]:
            for dy in [-1, 0, 1]:
                for dist in range(1, 8):
                    x, y = king[0] + dx * dist, king[1] + dy * dist
                    if (x, y) in hit:
                        res.append([x, y])
                        break

        return res

