# 221021天池-02. 光线反射

# 工程师正在研究一个 N*M 大小的光线反射装置，装置内部的构造记录于 grid 中，其中

# '.' 表示空白区域，不改变光的传播方向
# 'R' 表示向右倾斜的 双面 均可反射光线的镜子，改变光的传播方向
# 'L' 表示向左倾斜的 双面 均可反射光线的镜子，改变光的传播方向
# 假如光线从装置的左上方垂直向下进入装置，请问在离开装置前，光线在装置内部经过多长的路线。
# 经过的格子个数

from typing import List

DIR4 = [(0, 1), (1, 0), (0, -1), (-1, 0)]
MAPPING = {
    "R": {
        0: 3,
        1: 2,
        2: 1,
        3: 0,
    },
    "L": {
        0: 1,
        1: 0,
        2: 3,
        3: 2,
    },
}  # 注意到R就是异或3 L就是异或1


class Solution:
    def getLength(self, grid: List[str]) -> int:
        ROW, COL = len(grid), len(grid[0])
        curRow, curCol, curDir = 0, 0, 1
        step = 1

        while True:
            nextDir = curDir
            if grid[curRow][curCol] in MAPPING:
                if grid[curRow][curCol] == "R":
                    nextDir ^= 3
                else:
                    nextDir ^= 1
                # nextDir = MAPPING[grid[curRow][curCol]][curDir]
            curDir = nextDir
            dr, dc = DIR4[curDir]
            nextRow, nextCol = curRow + dr, curCol + dc
            if 0 <= nextRow < ROW and 0 <= nextCol < COL:
                curRow, curCol = nextRow, nextCol
                step += 1
            else:
                return step


print(Solution().getLength(grid=["...", "L.L", "RR.", "L.R"]))
