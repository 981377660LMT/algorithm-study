from typing import List
from collections import deque


MOD = int(1e9 + 7)
INF = int(1e20)


DIR8 = ((1, 0), (-1, 0), (0, 1), (0, -1), (1, 1), (-1, -1), (1, -1), (-1, 1))


class Solution:
    def lakeCount(self, grid: List[str]) -> int:
        # 以 W 为中心的八个方向相邻积水视为同一片池塘。
        # !WA原因:shadow name 注意bfs里的row,col要写curRow,curCol 最好不要写r c 与之前变量相同
        matrix = [list(row) for row in grid]
        res = 0
        ROW, COL = len(grid), len(grid[0])
        for r in range(ROW):
            for c in range(COL):
                if matrix[r][c] != "W":
                    continue
                res += 1
                queue = deque([(r, c)])
                while queue:
                    curRow, curCol = queue.popleft()
                    matrix[r][c] = "."
                    for dr, dc in DIR8:
                        nextRow, nextCol = curRow + dr, curCol + dc
                        if (
                            0 <= nextRow < ROW
                            and 0 <= nextCol < COL
                            and matrix[nextRow][nextCol] == "W"
                        ):
                            matrix[nextRow][nextCol] = "."
                            queue.append((nextRow, nextCol))
        return res


print(Solution().lakeCount(grid=[".....W", ".W..W.", "....W.", ".W....", "W.WWW.", ".W...."]))
print(
    Solution().lakeCount(
        grid=[
            "W.W.W....W.WW.W..WWW..WW..WW.WW..WWWW.W..WW...W...W..WWWWW.WW.WW....W.W.",
            "W...WW..W..W..WW..W...W.WWWW.W.WW..W.WW..W..W...W.W..W...WW.W.WWWW.WWW.W",
            "......WWW...W.WWWW.WW...W..W.WW..W...W.W.WW.WWW.W.WWWWWW..WW.WW...WW.W..",
            "W...WWWW.WWW.WWW.W..W.W..WWW.W.W.W.W.WWW.W.W...WWWW...W.W.WWW..WWW.W...W",
            "..W..WWW.W...W.WW.W.WWW...WW.WW.W..WWW..WW..WWW..W.WWW.W.WWW.W.WWWWWW..W",
            "WWWW.WWWW.WW.W..W.W.WW....WWWWW.WWW...W.W..WW.W.WWWWWW.W...WWW..WW.W...W",
            "W..WW....WW..WW.WW...W....WWWW....WWWWW..WW.W...W.W.WW.....W....W.WW.W.W",
            "W.....WWW..W..WW...W.WW.W.W..W.W.WWWW.WWWW..W..W....W..W.W..WW.WWWWWWW..",
            "WW.WWW.....W.W...WW...W........WW...WW..WW.W.W...W.WW...WWW....WWWWW.W.W",
            "...W..WWW...WW.W..W.W.W.WW.W.W.W.WWW.W..W.WW...WWW...WWWWW...WW.W.WW.W..",
            "WW..W.WWWWW.WWW..WWWWWWW.W...WWW.W..W.WWWW.WW..W..WWW...WWWW..WWWWWWWWWW",
            ".W.WWWW.W.....W.WWW..W.WWWWW..WW..WW..WW.W.W.WW.W.WWW.WWW..WW.W...W...WW",
            "W..WWW..WWW...W.W.W....WW.WWW.WWWW...W..W....WW.WW..W...WWWWWWW.W..W..W.",
            "W.W.W.W..W.W.W..WWW.W..W.WW.WW.W..WW.WW.W.WWW...WWWW..W.W.WWWW.WW...W..W",
            "WWWWW....W..W....WWW...WW..WW..WW..WW.WW.W.W.WWW..W...W.WW..WW..W.WW.WW.",
            ".W..W.WW..WWWWWW.WW..WW.WW.WWWWW.WWW..W.WWWWW.W.WWW.W.WWWWWW.W.W..WW..WW",
            ".W.WW..WW..W....W..WWW.W.....W.W.W...WW.W.W.W..WWW....WWWW..W..WWW...WW.",
            "...W..W.WWWWWW.WWW.WW...WW.WWW.W..W...W.W.WW...WW...WW..W..WWW.WW.WWW..W",
            "..WWWW..WWW.W.W..WW.W.W.W.....W..W...WWWW.WW.WWWW...W.WWWWWW.WWW.W..WWW.",
            "WWW.WW...W..W.W.W.W.W..W..WW.WWWWWW..W...WW....W.WW.WW.W...WWWWW.W.W..WW",
            "...WW..W.WWWW.WW....W...WW.W.....W..WWW.W.WW.W...W...W...WWWWWW.W.WWWW..",
            "WWW.....WW....W..W..WWW.W..W..WW..WWW.W..W.W.W..WWWW..W...WWW.WWW.WWWW.W",
            "W.WW.WWWW...WW..WWWWWW..WW.W.W.WWWWW.WW....W...WWW...W..W.W.....WW.WW..W",
            "WWWW.WWW.WWW..WWW..WWWW.WWW..WWWWW...W..WWWWW.W.W.W.WW.W..WW.W..WW.WWW.W",
            "W.WW...W.WWWW.W...W.WW.WW.W.WWWWW...W...W..W.W...WW.WW.WWWW......W.W..W.",
            "W.W.WWWWW.WWW....WWW..WWWW.W.W............WWW..W..W..W..WW.WWW...W.WWWW.",
            ".W...W.W.W.WW....W.WW.....WW.WWW.W..W....W.W...WWWWW.W.W....W.WW.W.WWWWW",
            "WW.W.W.WW.W.W.WW.W.W.W..WWW.WWWW.....WW...WWW...W.WW.WW..W.W....WWWW.WWW",
            ".W.WWW...WW.W..W.WWW.W....WWW.....W..W.W...W..W...WW.WW...WWWW..WWWWW...",
            "W...W...WWWW.WWWWW..WWWW....W...W...W.....W.W.W......WWWW...WW..W..W.W.W",
            "W.W.WWWW.W...W.WWWWW.WW.WWW.WW....WWWWWW..W.WW.WWW..W....WW.........WWW.",
            "..WW.WWWW.W.W.W..W..WW..WW...W.WW..WW....W.....WW.WW..W....WWW......WW.W",
            "W.W.WWW.W...WW...W.WWWWWW.W.WWWW.W.W..WWW.W.....W..W.W..W.WWW.WW.WWWWWW.",
            ".WWW.WW.WW..W.W.W.....WW..W....WWWW....WWWW.WWW.WWW.W....WWW.WW..W..WW..",
            "WW...WW.W..W.W.WW..W.WWWWWWWWW.W.....W....W..WW..W.WWWWWW..WWW.WW..WWW..",
            "W.WWWWWWWW.WW...WW..W.W.....WWWW.WW..WWW.W.WW.WWWW.WWW....WWWWW...W...WW",
            "WW.WW..W.WWW..W.W.WW..WWWW..WW.W.W....W..WWWWW.W.WWWWWWW..W..W..WWWW.W..",
            "..WW.WW..W..WW.WWW.WWW.W..W...W...W.WWW.W..W.WW....WW.WW.WW....W..WW....",
            "...W.W....WWW...WW.W...W.W.WW.W....W.WW.W....W..WW.WW.....WWW.WW.WW.WWW.",
            ".W.W.W...WW.W...W...WWWW..W.WWWW.W.W.WWWWW....WW....WWWWWW....WW.WWWW.W.",
            ".WWWW...W..W.W..W.WW..WWWW...W.W...WWW.WWW.W...WWWWW.W.W......W.WWWWWWW.",
            "W.W....W....W.W.W.WW...WW...W.W.WW..WWW.WW........W.WWWW....WW..W.WWWW..",
            "W..WW...W.WW.WWWW.WWWW.W.WW..W.W....W.WWWW...WW.W.....W.WWW.WWW....WW.WW",
            "WWW.W.WWW.W...WWW.W..WW.WW.W...WW.WWWWW.WW..WW....W.W...WW...WW.W.W...W.",
            "W.WWWW.W..W.WWW.W...W.....WWWWWWW.W.W..WWWW..WWWW.W.WW.W..WW...W.W..W.W.",
            "W.W..WW.W.W.WWWWWW...W.W.WWW.W.W..WW..W..W.WW.WW.....WW.W..W..W.WW.....W",
            "..WW.......W....WWWWW.WW..W..WW...WW.WW.W..W.W.WW.W.WW.WW....W.W.W......",
            "WWWW.....WWWW.....W.WWWWWW..WW.W..WWWWW.W..WWWW.W.WW..WW.WW.W.......WW..",
            "W.W...W.WWW.WW....W..WWWW....W.W.W.WWW..W....WW....W.W.W..WW..W..WW....W",
            "WWW..WW......WWW...WWWW......W...WWW.W..W..W.WW........W.W.W.WW.W.WWW.W.",
            "WWWW...WWW.W.WWWW.W.W.W...W.WWWW.W....W.W...WWWWW.W...WWWW.W...WW.WWWW.W",
            ".W...WWW...W.W.W.WW.W...WW.W.WWWWW..W.....WWW.W.W.WWWWW.WW.W..WWW....W..",
            "...WWWW.W..W...WWWW.....W.W.WW....W.W..W..W...WW.W.W..WW..W.WW.W.WWW...W",
            "W.W.W.W.WWWWW.....WWW.WW.W..WWW.W..WW.W.W...WW.WW.WWW..WWW.W.W.W.WW..W..",
            ".W.WW.WW..WW.W..W..W....WWWWW.....W.W..W..W.WW....W...W....W..W.WW.WW.WW",
            "WW....W......W.WWWWWW...W..W..W..WWWW..W.WWW.WWW..WW..W..W...WWWW....W..",
            "W.W..W..WWW.....W.WW.W..WW..WW..WW.....W.W.WW..W.WW.W.W..WW..W.WWWW.....",
            "....W......W.W..WW.WW.WW.WW.....WWWW..W.WW.WWWW..WW..WW..W.W...WWW.WWW.W",
            "WW.....W.WW.WWW.W...W.WWW.WWWW...WW.W.W.WW..W..WWWW.WWW.WWW.W.WW.W....WW",
            "...WW.WWW.W.W..W.WW.WWWW.W..WW.WWW..WW....W.WW.W.W..W.W.W.WWWW...WWW..WW",
            "W...W..W.WW.W..W..WWWWWWW.W...W...WW.WWW..W..W...W.WWW..WW.W.WW.WWW..WW.",
            "W.WW.W.WW.W.W..W..W..WWWW.W.WW...WWW.W..W...WWWW..W..WW...WWW..W....W.WW",
            "......W.WWW....W.......W.WW.W.WWWW.WWWW.W....WW....WW..W.WW...WW.W...W.W",
            "WWW.WW.W.WW.W.WWWWW.WWW....WW..WW..W.WWWWW.W.WW....WW.....W...WWWW.WWW..",
            "WWW.W.WW.W.W...WWW..WWWWWW...W..W.WWWW.WWW.W..WW.WW.WWWW....WWWW..WWWWW.",
            "..WW......W.WW.WWWWWWWW.WW...W.......W..W.WW.WWWWWWW..W.WWW..WWW..W.W..W",
            "..WW.WWW....W..WWW....W..W.W....W.W...W...WWW...WW....W..WWWW.WWW.W.WWWW",
            "WW..WW..W...WWW..W.....W..WW.W.WWWW..WWWW.WWWWWW...WW...WW.........WW.W.",
            "WW.WWW.W.W.WWWW.WW.WW..WW.W.W.W.W.WW...W.W..W..W.W.WWW...WW..WW.W....W..",
            "W.W.W.WW..WWW.....WW......WW..WWWW.W.W.W...WWWW....WWW..WW...WW..W..WWWW",
            ".W..WW...WW..W.W..W.WWWWW..W....WW.WWWWWWW..WWWWWW..WW....WW.W...W.WWW..",
            "WW....W.W..WW.WW.W.W.WW..WW.....WW.WWWWWW.WWW.W..WWW...W.WW..WWWWW......",
            ".W....W.W.WWW.W.W..WW...WW.WWW..WWWWW.WW.W.WW..W..W..W...WW.WWWW.W.W.W..",
            "W.W....WWW..W.WW.W.W.WW..WW.W..WW.W.W.W.WW.WWWW..W...W..WWW.W.W.....W.W.",
            "W..WW.WWW....W.WWW..W.WWW..W.W..W..W.WWWW...WWW...W..W.WWWWWWW...WWWW...",
            ".WW.WWWWW..WW.WWWW..WWWW......WW.W...W.W.W..WW.W.W..WW.WW.WW.WWW...WWWW.",
            "WWW....WWW........W..W..W.W..WW.WW.WW....WWW.....WWWW.W..WWWW.W....WWWW.",
            ".WWW....WWWW...WWW..W..W.W...WWWW..WWWWWW.W..W...WW..WWWW.WW.W..W..WW...",
            "W.W......W..WWW..W.W.WWWW..WWW..W.W.......WWW.WWWW.....WW.WW.WWWWWW..W.W",
            "WW.W.....WW..WWWWW..W..WW..W.WWWW..W..W.W.W.W..WWW.W.WWWW..W.W.W..W.WWWW",
            "WWWW..WWW..WWWW..W.WW..WW..WWW..W..WW....WW.W...W...WWWWW....W.WWW..W.WW",
            "..........W.WW.W..WW.WWW...W..W...WWWW.W.WW.W.WW..WW...W.WW.WWWW..W...W.",
            ".WWW.W.W.W..WWW..WWW...WW.W.WW...WW...W....W......WWWWWWWWWW..W..W...W.W",
            "....W..W.WW.W.WW.W.....W.W.......WWW.WWW..WWWWW..WWWWWWW...W..WW.WW..W..",
            "W.WW.WWW..W.W.W...W.W.WWWW.....WWWWWWWW.W...WWWWW...WW.W..WWWWWWW..W....",
            "WW.....WWW.......W..WW.WWW..WW...WWW.WWWW.W..W...WWW...WW.W.WWWWWWWWW.W.",
            "...W...WWW..W..W..W.W...WWW.W.WWW..WW.W........WWW..WWW...W.W.WWWW.WW..W",
        ]
    )
)
