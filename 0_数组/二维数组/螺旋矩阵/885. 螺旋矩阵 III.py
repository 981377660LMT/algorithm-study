from typing import List

# 我们以顺时针按螺旋状行走，访问此网格中的每个位置。
# 按照访问顺序返回表示网格位置的坐标列表。

# 1,1,2,2,3,3 Steps
# 关键:计算步数与转向次数的关系
class Solution:
    def spiralMatrixIII(self, rows: int, cols: int, rStart: int, cStart: int) -> List[List[int]]:
        res = []
        dr, dc, turn = 0, 1, 0
        while len(res) < rows * cols:
            for _ in range(turn // 2 + 1):
                if 0 <= rStart < rows and 0 <= cStart < cols:
                    res.append([rStart, cStart])
                rStart, cStart = rStart + dr, cStart + dc
            dr, dc, turn = dc, -dr, turn + 1

        return res


print(Solution().spiralMatrixIII(5, 6, 1, 4))
# 输出：[[1,4],[1,5],[2,5],[2,4],[2,3],[1,3],[0,3],[0,4],[0,5],[3,5],[3,4],[3,3],[3,2],[2,2],[1,2],[0,2],[4,5],[4,4],[4,3],[4,2],[4,1],[3,1],[2,1],[1,1],[0,1],[4,0],[3,0],[2,0],[1,0],[0,0]]
