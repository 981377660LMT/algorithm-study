from typing import List

# 如果将格子 (rMove, cMove) 变成颜色 color 后，
# 是一个 合法 操作，那么返回 true ，如果不是合法操作返回 false 。
# 合法 操作必须满足：`涂色后这个格子是 好线段的一个端点` （好线段可以是水平的，竖直的或者是对角线）。

# 八方遍历 + 讨论遇到什么点
class Solution:
    def checkMove(self, board: List[List[str]], rMove: int, cMove: int, color: str) -> bool:
        for di, dj in (0, 1), (1, 1), (1, 0), (1, -1), (0, -1), (-1, -1), (-1, 0), (-1, 1):
            i, j = rMove + di, cMove + dj
            step = 0
            while 0 <= i < 8 and 0 <= j < 8:
                if board[i][j] == color and step > 0:
                    return True
                if board[i][j] == "." or board[i][j] == color and not step:
                    break
                i, j = i + di, j + dj
                step += 1
        return False


print(
    Solution().checkMove(
        board=[
            [".", ".", ".", "B", ".", ".", ".", "."],
            [".", ".", ".", "W", ".", ".", ".", "."],
            [".", ".", ".", "W", ".", ".", ".", "."],
            [".", ".", ".", "W", ".", ".", ".", "."],
            ["W", "B", "B", ".", "W", "W", "W", "B"],
            [".", ".", ".", "B", ".", ".", ".", "."],
            [".", ".", ".", "B", ".", ".", ".", "."],
            [".", ".", ".", "W", ".", ".", ".", "."],
        ],
        rMove=4,
        cMove=3,
        color="B",
    )
)

