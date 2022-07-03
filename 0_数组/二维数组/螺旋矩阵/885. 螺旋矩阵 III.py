from typing import List

# 我们以顺时针按螺旋状行走，访问此网格中的每个位置。
# 按照访问顺序返回表示网格位置的坐标列表。

# 与II不同的是螺旋II是填充问题，观察是否被填过就可以判断方向，而III是步长是固定的，哪怕越出边界也不改变方向
# 1,1,2,2,3,3 Steps
# 关键:计算步数与转向次数的关系

DIR4 = ((0, 1), (1, 0), (0, -1), (-1, 0))


class Solution:
    def spiralMatrixIII(
        self, rows: int, cols: int, rStart: int, cStart: int
    ) -> List[List[int]]:
        res = []
        r, c, di = rStart, cStart, 0
        turn = 0

        while len(res) < rows * cols:
            # !这个方向走的步数
            for _ in range(turn // 2 + 1):
                if 0 <= r < rows and 0 <= c < cols:
                    res.append([r, c])
                r, c = r + DIR4[di][0], c + DIR4[di][1]
            di, turn = (di + 1) % 4, turn + 1

        return res


print(Solution().spiralMatrixIII(5, 6, 1, 4))
# 输出：[[1,4],[1,5],[2,5],[2,4],[2,3],[1,3],[0,3],[0,4],[0,5],[3,5],[3,4],[3,3],[3,2],[2,2],[1,2],[0,2],[4,5],[4,4],[4,3],[4,2],[4,1],[3,1],[2,1],[1,1],[0,1],[4,0],[3,0],[2,0],[1,0],[0,0]]
