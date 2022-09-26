from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# plate[i][j] 仅包含 "O"、"W"、"E"、"."
# N*M 大小的弹珠盘的初始状态信息记录于一维字符串型数组 plate 中，
# 数组中的每个元素为仅由 "O"、"W"、"E"、"." 组成的字符串。其中：

# "O" 表示弹珠洞（弹珠到达后会落入洞中，并停止前进）；
# "W" 表示逆时针转向器（弹珠经过时方向将逆时针旋转 90 度）；
# "E" 表示顺时针转向器（弹珠经过时方向将顺时针旋转 90 度）；
# "." 表示空白区域（弹珠可通行）。
# 游戏规则要求仅能在边缘位置的 空白区域 处（弹珠盘的四角除外）沿 与边缘垂直 的方向打入弹珠，
# 并且打入后的每颗弹珠最多能 前进 num 步。请返回符合上述要求且可以使弹珠最终入洞的所有打入位置。
# 你可以 按任意顺序 返回答案。

# 注意：

# 若弹珠已到达弹珠盘边缘并且仍沿着出界方向继续前进，则将直接出界。

# !倒着寻找(才怪! 正着的起点少)

DIR4 = {
    0: (0, 1),
    1: (1, 0),
    2: (0, -1),
    3: (-1, 0),
}  # 顺时针


class Solution:
    def ballGame(self, num: int, plate: List[str]) -> List[List[int]]:
        ROW, COL = len(plate), len(plate[0])
        queue = deque()
        visited = [[False] * 4 for _ in range(ROW * COL)]
        BAD = set([(0, 0), (0, COL - 1), (ROW - 1, 0), (ROW - 1, COL - 1)])
        for i in range(ROW):
            for j in range(COL):
                if (i, j) in BAD:
                    continue
                if i == 0 or i == ROW - 1 or j == 0 or j == COL - 1:
                    if plate[i][j] == ".":
                        dir = -1
                        if i == 0:
                            dir = 1
                        elif i == ROW - 1:
                            dir = 3
                        elif j == 0:
                            dir = 0
                        elif j == COL - 1:
                            dir = 2
                        queue.append((i, j, dir, 0, i, j))
                        visited[i * COL + j][dir] = True
        res = []
        while queue:
            curRow, curCol, curDir, curStep, rootRow, rootCol = queue.pop()
            if plate[curRow][curCol] == "O":
                res.append([rootRow, rootCol])
                continue
            if plate[curRow][curCol] == "W":
                nextDir = (curDir - 1) % 4  # 逆时针
            elif plate[curRow][curCol] == "E":
                nextDir = (curDir + 1) % 4  # 顺时针
            else:
                nextDir = curDir
            nextRow, nextCol = curRow + DIR4[nextDir][0], curCol + DIR4[nextDir][1]
            if not (nextRow < 0 or nextRow >= ROW or nextCol < 0 or nextCol >= COL):
                hash_ = nextRow * COL + nextCol
                if curStep + 1 <= num and not visited[hash_][nextDir]:
                    queue.append((nextRow, nextCol, nextDir, curStep + 1, rootRow, rootCol))
                    visited[hash_][nextDir] = True
        return res


print(
    Solution().ballGame(
        num=76,
        plate=[
            "E......O..",
            "E.........",
            "W..E...EW.",
            "EE.OE.WWWO",
            "O.WEOEWWW.",
            ".OW..W....",
            "W...EW.WE.",
            ".E.W...OW.",
            "OW...W.EEO",
            "......W.W.",
        ],
    )
)
