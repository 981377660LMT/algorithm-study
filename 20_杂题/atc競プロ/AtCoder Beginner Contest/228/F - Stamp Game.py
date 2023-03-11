# 给定一个N ×M的矩阵Aij，高桥可以将hi*w1大小内的矩阵涂黑，
# 即将这个区域内的格子涂满黑色，青木可以将h2*w2大小内的矩阵涂白，
# 一开始矩阵内的方格都是白色的，最终高桥的分数就是矩阵内黑色方格的分数总和，
# 高桥让分数最大化，青木想让他的分数最小化，两人都是最聪明的，
# 问双方涂完色后，高桥可能得到的最高分数是多少.

# !二维窗口最值
# 1. 滑窗+单调队列
# 2. 线段树
# 3. st表
from typing import List


class P:
    """二维前缀和模板(矩阵不可变)"""

    def __init__(self, A: List[List[int]]):
        m, n = len(A), len(A[0])

        # 前缀和数组
        preSum = [[0] * (n + 1) for _ in range(m + 1)]
        for r in range(m):
            for c in range(n):
                preSum[r + 1][c + 1] = A[r][c] + preSum[r][c + 1] + preSum[r + 1][c] - preSum[r][c]
        self.preSum = preSum

    def sumRegion(self, r1: int, c1: int, r2: int, c2: int) -> int:
        """查询sum(A[r1:r2+1, c1:c2+1])的值::

        M.sumRegion(0, 0, 2, 2) # 左上角(0, 0)到右下角(2, 2)的值
        """
        return (
            self.preSum[r2 + 1][c2 + 1]
            - self.preSum[r2 + 1][c1]
            - self.preSum[r1][c2 + 1]
            + self.preSum[r1][c1]
        )


from collections import deque
from typing import List


class MonoQueue:
    __slots__ = ("maxQueue", "minQueue", "rawQueue")

    def __init__(self):
        self.maxQueue = deque()
        self.minQueue = deque()
        self.rawQueue = deque()

    @property
    def min(self) -> int:
        return self.minQueue[0][0]

    @property
    def max(self) -> int:
        return self.maxQueue[0][0]

    def popleft(self) -> int:
        if not self.rawQueue:
            raise IndexError("popleft from empty queue")

        self.minQueue[0][1] -= 1
        if self.minQueue[0][1] == 0:
            self.minQueue.popleft()

        self.maxQueue[0][1] -= 1
        if self.maxQueue[0][1] == 0:
            self.maxQueue.popleft()

        return self.rawQueue.popleft()

    def append(self, value: int) -> None:
        count = 1
        while self.minQueue and self.minQueue[-1][0] > value:
            count += self.minQueue.pop()[1]
        self.minQueue.append([value, count])

        count = 1
        while self.maxQueue and self.maxQueue[-1][0] < value:
            count += self.maxQueue.pop()[1]
        self.maxQueue.append([value, count])

        self.rawQueue.append(value)

    def __len__(self) -> int:
        return len(self.rawQueue)


def windowMax2D(
    matrix: List[List[int]], rowSize: int, colSize: int, *, isMax: bool
) -> List[List[int]]:
    """
    每个rowSize*colSize窗口内的最值
    1 ≤ n * m ≤ 1e5
    """
    ROW, COL = len(matrix), len(matrix[0])
    cols = [MonoQueue() for _ in range(COL)]
    for r in range(rowSize):
        for c in range(COL):
            cols[c].append(matrix[r][c])

    res = []
    for r in range(ROW - rowSize + 1):
        res.append([])
        # 维护这一行的滑动窗口最值
        window = MonoQueue()
        for c in range(COL):
            # 每列的最值
            window.append(cols[c].max if isMax else cols[c].min)
            if c >= colSize - 1:
                res[-1].append(window.max if isMax else window.min)
                window.popleft()

        # 下一行进入窗口
        if r + rowSize < ROW:
            for c in range(COL):
                cols[c].append(matrix[r + rowSize][c])
                cols[c].popleft()

    return res


if __name__ == "__main__":
    ROW, COL, row1, col1, row2, col2 = map(int, input().split())
    grid = [list(map(int, input().split())) for _ in range(ROW)]
    P = P(grid)

    rowMin, colMin = min(row1, row2), min(col1, col2)
    if row1 == rowMin and col1 == colMin:
        exit(print(0))  # 先手无法涂黑色

    whiteSum = [[] for _ in range(ROW - rowMin + 1)]  # 后手涂白色的区域和
    for r in range(ROW - rowMin + 1):
        for c in range(COL - colMin + 1):
            whiteSum[r].append(P.sumRegion(r, c, r + rowMin - 1, c + colMin - 1))

    # 在每个row1*col1窗口范围内所有的rowMin*colMin窗口的最大值
    windowMax = windowMax2D(whiteSum, row1 - rowMin + 1, col1 - colMin + 1, isMax=True)

    res = 0
    # 枚举先手黑色玩家的选择
    for r in range(ROW - row1 + 1):
        for c in range(COL - col1 + 1):
            black = P.sumRegion(r, c, r + row1 - 1, c + col1 - 1)
            white = windowMax[r][c]
            res = max(res, black - white)
    print(res)
