# 给你两个整数 m 和 n ，表示一个下标从 0 开始的 m x n 的网格图。
# 给你一个下标从 0 开始的二维整数矩阵 coordinates ，
# 其中 coordinates[i] = [x, y] 表示坐标为 [x, y] 的格子是 黑色的 ，所有没出现在 coordinates 中的格子都是 白色的。
# 一个块定义为网格图中 2 x 2 的一个子矩阵。
# 更正式的，对于左上角格子为 [x, y] 的块，其中 0 <= x < m - 1 且 0 <= y < n - 1 ，
# 包含坐标为 [x, y] ，[x + 1, y] ，[x, y + 1] 和 [x + 1, y + 1] 的格子。
# !请你返回一个下标从 0 开始长度为 5 的整数数组 arr ，arr[i] 表示恰好包含 i 个 黑色 格子的块的数目。

# !1.规定块的左上角为块的代表点
# !2.枚举每个黑格子/黑方块对块的(代表点)的贡献


from collections import defaultdict
from typing import List, Tuple

DIR4 = [(0, 0), (-1, 0), (0, -1), (-1, -1)]


class Solution:
    def countBlackBlocks(
        self, ROW: int, COL: int, coordinates: List[List[int]]
    ) -> Tuple[int, int, int, int, int]:
        counter = defaultdict(int)
        for r, c in coordinates:
            for dr, dc in DIR4:
                nr, nc = r + dr, c + dc  # !当前黑格子对块的贡献
                if 0 <= nr < ROW - 1 and 0 <= nc < COL - 1:
                    counter[(nr, nc)] += 1

        res = [0] * 5
        for v in counter.values():
            res[v] += 1
        res[0] = (ROW - 1) * (COL - 1) - sum(res[1:])
        return tuple(res)
