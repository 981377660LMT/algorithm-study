from functools import cache
from typing import List, Tuple

# 1 <= m <= 8
# 1 <= n <= 8


class Solution:
    def maxStudents(self, seats: List[List[str]]) -> int:
        """轮廓线dp,复杂度n * m * (2 ^ m)
        
        保存前COL+1个位置的信息
        """

        @cache
        def dfs(r: int, c: int, cState: Tuple[int, ...]) -> int:
            if r == ROW:
                return 0
            if c == COL:
                return dfs(r + 1, 0, cState)

            res = dfs(r, c + 1, cState[1:] + (0,))  # 不选择当前位置
            if seats[r][c] == '.':
                """不能看到左侧、右侧、左上、右上这四个方向上紧邻他的学生的答卷"""
                leftUp = cState[0] if r and c else 0
                rightUp = cState[2] if r and c + 1 < COL else 0
                left = cState[-1] if c else 0
                if leftUp == rightUp == left == 0:
                    res = max(res, dfs(r, c + 1, cState[1:] + (1,)) + 1)

            return res

        ROW, COL = len(seats), len(seats[0])
        res = dfs(0, 0, tuple([0] * (COL + 1)))
        dfs.cache_clear()
        return res


print(
    Solution().maxStudents(
        seats=[
            ["#", ".", "#", "#", ".", "#"],
            [".", "#", "#", "#", "#", "."],
            ["#", ".", "#", "#", ".", "#"],
        ]
    )
)

